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

package resource

import (
	errs "errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/apis/internals"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/dao/instance"
	"github.com/kappital/kappital/pkg/models"
	mo "github.com/kappital/kappital/pkg/models/operation"
	co "github.com/kappital/kappital/pkg/utils/operations"
	"github.com/kappital/kappital/pkg/watcher"
)

// InstanceResource operate instance information in database and/or cluster
type InstanceResource struct {
	instanceStore instance.Instance
	binding       mo.ServiceBindingOperation
}

// GetResourceType of the service instance
func (i *InstanceResource) GetResourceType() Type {
	return ServiceInstanceType
}

// GetCommonDBObject get the data from database and return the inner structure will use in manager
func (i *InstanceResource) GetCommonDBObject(pk string) (interface{}, error) {
	obj, err := i.instanceStore.GetByPrimaryKey(pk)
	if err != nil {
		klog.Errorf("failed to get instance, error: %v", err)
		return nil, err
	}
	ins, ok := obj.(internals.ServiceInstance)
	if !ok {
		klog.Errorf("failed to trans service binding, error: %v", err)
		return nil, err
	}
	return &ins, nil
}

// GetObjectProcessTime get the process time for the service instance (using for synchronizing)
func (i *InstanceResource) GetObjectProcessTime(obj interface{}) time.Time {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
		return time.Time{}
	}
	return ins.ProcessTime
}

// GetObjectStatus get the service instance status
func (i *InstanceResource) GetObjectStatus(obj interface{}) string {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
		return ""
	}
	return ins.Status
}

// GetObjectID get the service instance id (primary key in database )
func (i *InstanceResource) GetObjectID(obj interface{}) string {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
		return ""
	}
	return ins.ID
}

// GetObjectListByStatusSets get the service instance list with the status
func (i *InstanceResource) GetObjectListByStatusSets(status sets.String) ([]interface{}, error) {
	tmp, err := i.instanceStore.GetListByStatusSets(status)
	if err != nil {
		klog.Errorf("get status %s by instance failed, error: %s", status, err)
		return nil, err
	}
	instances, ok := tmp.([]internals.ServiceInstance)
	if !ok {
		return nil, fmt.Errorf("cannot get the service instance from internal")
	}
	result := make([]interface{}, 0, len(instances))
	for i := range instances {
		result = append(result, &instances[i])
	}
	return result, nil
}

// UpdateProcessFailed update the service instance status to Failed
func (i *InstanceResource) UpdateProcessFailed(obj interface{}, status string, msg string) error {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
		return fmt.Errorf("invalid object type,expected:models.instance, actual: %s", reflect.TypeOf(obj).Name())
	}

	return i.instanceStore.UpdateStatusMsg(*ins, status, msg)
}

// UpdateObjProcessTime update the process time for the service instance (using for synchronizing)
func (i *InstanceResource) UpdateObjProcessTime(obj interface{}, processTime time.Time) error {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		return fmt.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
	}
	ins.ProcessTime = processTime

	return i.instanceStore.Update(*ins, "process_time")
}

// GetObjUpdateTime get the service instance update timestamp
func (i *InstanceResource) GetObjUpdateTime(obj interface{}) time.Time {
	ins, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.instance, actual: %s", reflect.TypeOf(obj).Name())
		return time.Time{}
	}
	return ins.UpdateTime
}

// CreateInstance into database, and add event to the synchronizing list
// which for deploying the service instance into cluster
func (i *InstanceResource) CreateInstance(instances []internals.ServiceInstance, param map[string]string) error {
	var needAddInstances []internals.ServiceInstance
	for _, instanceTemp := range instances {
		indb, err := i.instanceStore.Get(map[string]string{
			"name":       instanceTemp.Name,
			"namespace":  instanceTemp.Namespace,
			"cluster_id": instanceTemp.ClusterID,
		})
		if err != nil {
			if errs.Is(err, orm.ErrNoRows) {
				needAddInstances = append(needAddInstances, instanceTemp)
				continue
			}
			return err
		}
		klog.Infof("service instance %s has been created", indb.(internals.ServiceInstance).Name)
	}

	if len(needAddInstances) == 0 {
		klog.Infof("all instance has created, no need to create")
		return nil
	}

	if err := i.instanceStore.Create(needAddInstances, param); err != nil {
		klog.Infof("create instance failed, error: %s", err)
		return err
	}

	for _, needCreateInstance := range needAddInstances {
		if err := watcher.AddEvent(needCreateInstance, watcher.OPCreate, apis.InstanceProcessor); err != nil {
			klog.Errorf("[ADD EVENT] add instance %s created event failed, err: %s", needCreateInstance.Name, err)
			return err
		}
		klog.Infof("[ADD EVENT] add instance %s created event success", needCreateInstance.Name)
	}

	return nil
}

// UpdateInstallCondition of the instance
func (i *InstanceResource) UpdateInstallCondition(ins *internals.ServiceInstance,
	conType instance.InstallConditionType, status instance.ConditionStatus, msg string) error {
	if ins.InstallState.SubPhase == nil || len(ins.InstallState.SubPhase) == 0 {
		// init the parameters if not init
		ins.InstallState.SubPhase = i.GetInstanceInitialCondition()
		ins.InstallState.Phase = models.StatusInitializing
		ins.Status = models.StatusInitializing
	}
	for idx, cond := range ins.InstallState.SubPhase {
		if cond.Type == string(conType) {
			if cond.Status == string(status) && !needUpdateMsg(cond.Message, msg) {
				// status is not change
				return nil
			}
			if cond.Status == string(instance.Success) {
				// status has already changed to success, ignore
				return nil
			}
			// if all status is success, will not refresh the message
			cond.Message = msg
			if cond.Status != string(status) {
				cond.LastTransitionTime = metav1.Now()
			}
			cond.Status = string(status)
			ins.InstallState.SubPhase[idx] = cond
			if cond.Status == string(instance.Success) && idx == len(ins.InstallState.SubPhase)-1 {
				// if all status are success
				ins.ProcessTime = time.Time{}
				ins.Status = models.StatusInitialized
				ins.InstallState.Phase = models.StatusInitialized
			}
			break
		}
		if cond.Status != string(instance.Success) {
			cond.Status = string(instance.Success)
			cond.Message = fmt.Sprintf("successful to %s", cond.Type)
			cond.LastTransitionTime = metav1.Now()
			ins.InstallState.SubPhase[idx] = cond
		}
	}

	err := i.instanceStore.Update(*ins, "install_state", "status")
	if err != nil {
		klog.Errorf("failed to update instance %s in cluster %s to db, error: %v", ins.Name, ins.ClusterID, err)
		return err
	}

	return nil
}

// GetInstanceInitialCondition get the instance condition which using for the synchronizing
func (i *InstanceResource) GetInstanceInitialCondition() []internals.Condition {
	return []internals.Condition{
		{
			Type:               string(instance.InstallOperator),
			Status:             string(instance.Waiting),
			Message:            "",
			LastTransitionTime: metav1.Now(),
			RetryCount:         0,
		},
		{
			Type:               string(instance.CreateResource),
			Status:             string(instance.Waiting),
			Message:            "",
			LastTransitionTime: metav1.Now(),
			RetryCount:         0,
		},
	}
}

func needUpdateMsg(old, new string) bool {
	return old != new && new != ""
}

// ValidationInstance check the Instance data in memory is valid or not
func ValidationInstance(instanceCreation *instancev1alpha1.ServiceInstanceCreation) error {
	// check cr namespace is existed
	for _, cr := range instanceCreation.InstanceCustomResources {
		if cr.Name == "" {
			return fmt.Errorf("deploy service %s instance cr metadata.name must be not null",
				instanceCreation.InstanceName)
		}
		// check namespace is existed
		isExist, err := co.GetClusterOperation().IsNamespaceExist(cr.Namespace)
		if err != nil {
			return fmt.Errorf("check cr namespace %s failed, err: %s", cr.Namespace, err)
		}
		if !isExist {
			return fmt.Errorf("check cr namespace %s not found", cr.Namespace)
		}
	}
	return nil
}

// GetInstances get the instance list from database, and filter it by service binding name, cluster name, and namespace
func (i *InstanceResource) GetInstances(sbName, clusterName, ns string) ([]models.InstanceModel, error) {
	// get ServiceBindingModel from the database
	sb, err := i.binding.GetDetail(map[string]string{"name": sbName, "cluster_name": clusterName})
	if err != nil {
		return nil, err
	}
	// get the related resources' instances from the ServiceBindingModel,
	// and check does the instance is existed in cluster, if not exist, update the status in database
	instances, err := i.getAndCheckInstanceInCluster(sb.(models.ServiceBindingModel).Resources, ns)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// GetInstance get the instance from database, and filter it by service binding name, cluster name, and namespace
func (i *InstanceResource) GetInstance(sbName, clusterName, ns, instanceName string) (models.InstanceModel, error) {
	instances, err := i.GetInstances(sbName, clusterName, ns)
	if err != nil {
		return models.InstanceModel{}, err
	}
	for _, item := range instances {
		if item.Name == instanceName {
			return item, nil
		}
	}
	return models.InstanceModel{}, fmt.Errorf("the instance [%s] is not found in cluster or database", instanceName)
}

func (i *InstanceResource) getAndCheckInstanceInCluster(resources []*models.ResourceModel,
	ns string) ([]models.InstanceModel, error) {
	var instances []models.InstanceModel
	for _, item := range resources {
		gv := schema.GroupVersion{
			Group:   item.Group,
			Version: strings.Replace(item.APIVersion, item.Group+"/", "", -1),
		}
		for _, ins := range item.Instances {
			if ins.Namespace != ns {
				continue
			}
			status, err := i.checkInstanceAndUpdate(gv, item.Resource, ins.Name, ns)
			if err != nil {
				return nil, err
			}
			if ins.Status != string(instancev1alpha1.PendingPhase) {
				ins.Status = status
			}
			io := mo.InstanceOperation{}
			if err = io.Update(*ins); err != nil {
				return nil, err
			}
			instances = append(instances, *ins)
		}
	}

	return instances, nil
}

func (i *InstanceResource) checkInstanceAndUpdate(gv schema.GroupVersion, plural, name, ns string) (string, error) {
	find, err := co.GetClusterOperation().DoesCustomResourceExist(gv, plural, name, ns)
	if err != nil {
		return "", err
	}
	if !find {
		return string(instancev1alpha1.FailedPhase), nil
	}
	return string(instancev1alpha1.SucceededPhase), nil
}

// DeleteInstance in database and cluster
func (i *InstanceResource) DeleteInstance(clusterName, instanceName, namespace string) error {
	tmp, err := i.instanceStore.Get(map[string]string{"name": instanceName, "namespace": namespace,
		"cluster_name": clusterName})
	if err != nil {
		return err
	}

	item, ok := tmp.(internals.ServiceInstance)
	if !ok {
		klog.Errorf("obj type is not ServiceInstance, actual: %s", reflect.TypeOf(item).Name())
		return fmt.Errorf("delete instance %s failed, because get data from db failed", instanceName)
	}

	item.Status = models.StatusDeleting
	item.ProcessTime = time.Time{}
	item.UpdateTime = time.Now().UTC()
	if err = i.instanceStore.Update(item, "status", "process_time", "update_timestamp"); err != nil {
		klog.Errorf("failed to update instance[%s] in cluster[%s] into db, error: %s", instanceName, clusterName, err)
		return err
	}

	if err = watcher.AddEvent(item, watcher.OPDelete, apis.InstanceProcessor); err != nil {
		klog.Infof("[ADD EVENT]service binding %s created", item.Name)
		return err
	}
	return nil
}
