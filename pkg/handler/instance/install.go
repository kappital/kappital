/*
 * Copyright 2022 Huawei Cloud Computing Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instance

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/dao/instance"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/resource"
	co "github.com/kappital/kappital/pkg/utils/operations"
)

const (
	instanceProcessTimeout     = 3 * time.Minute
	bindingReadyProcessTimeout = 2 * time.Minute
)

// Handler singleton pattern of install, delete, or upgrade the service instance
type Handler struct {
	instance resource.InstanceResource
}

// BeforeInstall do some processes before install the service instance
func (h *Handler) BeforeInstall(obj interface{}) (retry bool, err error) {
	serviceInstance, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type for upgrade instance handler, expected: models.ServiceInstance, "+
			"actual: %s", reflect.TypeOf(obj).Name())
		return false, nil
	}
	var operatorSuccess bool
	defer func() {
		if err == nil {
			err = h.instance.UpdateInstallCondition(serviceInstance, instance.CreateResource, instance.Running, "")
			if err != nil {
				klog.Errorf("failed to update instance %s condition, error: %v", serviceInstance.Name, err)
			}
		} else if operatorSuccess {
			innerErr := h.instance.UpdateInstallCondition(serviceInstance, instance.InstallOperator,
				instance.Success, "all servicebindings status are succeed")
			if innerErr != nil {
				klog.Errorf("failed to update instance %s condition, error: %v", serviceInstance.Name, err)
				retry = true
			}
		}
	}()
	sp, found, err := co.GetClusterOperation().GetServicePackageByName(serviceInstance.ServiceBindingName,
		apis.KappitalSystemNamespace)
	if err != nil {
		klog.Errorf("get service binding %s failed, err: %s", serviceInstance.ServiceBindingName, err)
		return true, err
	}
	if !found {
		return true, nil
	}

	if err = updateProcessTimeout(serviceInstance, bindingReadyProcessTimeout); err != nil {
		return true, err
	}

	if strContains(sp.Status.Phase, models.FailedStatusList) {
		err = h.instance.UpdateProcessFailed(serviceInstance, models.StatusInitFailed,
			fmt.Sprintf("servicebinding %s failed status in cluster %s",
				serviceInstance.ServiceBindingName, serviceInstance.ClusterID))
		if err != nil {
			// only for synchronize with database, other situation will not retry
			klog.Errorf("failed to update failed record of operator %s, error: %v", serviceInstance.ID, err)
			return true, err
		}
		return false, fmt.Errorf("operator process failed, stop install instance %s", serviceInstance.ID)
	}

	operatorSuccess = true

	return false, nil
}

func updateProcessTimeout(ins *internals.ServiceInstance, timeout time.Duration) error {
	if !ins.ProcessTime.IsZero() {
		return nil
	}
	ins.ProcessTime = time.Now().Add(timeout)
	instanceDB := instance.Instance{}
	err := instanceDB.Update(*ins, "process_time")
	if err != nil {
		return err
	}
	return nil
}

// Install deploy/apply the service custom resource into cluster
func (h *Handler) Install(obj interface{}) (retry bool, err error) {
	item, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type for upgrade instance handler, expected: models.ServiceInstance, "+
			"actual: %s", reflect.TypeOf(obj).Name())
		return false, nil
	}
	var exist bool
	defer func() {
		if exist {
			err = h.instance.UpdateInstallCondition(item, instance.CreateResource, instance.Success,
				"instance resource create success")
			if err != nil {
				klog.Errorf("failed to update instance %s condition, error: %v", item.Name, err)
				retry = true
			}
		} else if err != nil {
			innerErr := h.instance.UpdateInstallCondition(item, instance.CreateResource, instance.Running,
				fmt.Sprintf("failed to install instance, error: %v", err))
			if innerErr != nil {
				klog.Errorf("failed to update instance %s condition, error: %v", item.Name, innerErr)
				retry = true
			}
		}
	}()

	// renew the process time
	if err = updateProcessTimeout(item, instanceProcessTimeout); err != nil {
		klog.Errorf("failed to update process time for instance %s", item.Name)
		return true, err
	}

	retry, nextProcess, err := doesCustomResourceExist(item)
	if !nextProcess {
		return retry, err
	}

	retry, nextProcess, err = deployCustomResource(item)
	if !nextProcess {
		return retry, err
	}
	exist = true
	return false, nil
}

// AfterInstall install service instance does not need to implement this method
func (h *Handler) AfterInstall(_ interface{}) (bool, error) {
	return false, nil
}

// strContains check string array contains specific string
func strContains(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func getGroupVersion(apiVersion string) (schema.GroupVersion, error) {
	apiVersionSplit := strings.Split(apiVersion, "/")
	if len(apiVersionSplit) != 2 {
		return schema.GroupVersion{}, fmt.Errorf("invalid api version")
	}
	return schema.GroupVersion{
		Group:   apiVersionSplit[0],
		Version: apiVersionSplit[1],
	}, nil
}

func doesCustomResourceExist(item *internals.ServiceInstance) (bool, bool, error) {
	gv, err := getGroupVersion(item.APIVersion)
	if err != nil {
		klog.Errorf("cannot get group version, err: %s", err)
		return true, false, err
	}
	exist, err := co.GetClusterOperation().DoesCustomResourceExist(gv, item.Resource, item.Name, item.Namespace)
	if err != nil {
		klog.Errorf("query cr %s is exist failed, err: %s", item.Name, err)
		return true, false, err
	}
	if exist {
		return false, false, nil
	}
	return false, true, nil
}

func deployCustomResource(item *internals.ServiceInstance) (bool, bool, error) {
	gv, err := getGroupVersion(item.APIVersion)
	if err != nil {
		klog.Errorf("cannot get group version, err: %s", err)
		return true, false, err
	}
	gvr := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: item.Resource}
	klog.Infof("gvr %s", gvr)
	if err = co.GetClusterOperation().DeployCustomResource(gvr, item.Namespace, item.RawResource); err != nil {
		klog.Errorf("create cr %s failed, err: %s", item.Name, err)
		return true, false, err
	}
	return false, true, nil
}
