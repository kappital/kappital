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

package servicebinding

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/dao/instance"
	"github.com/kappital/kappital/pkg/dao/servicebinding"
	co "github.com/kappital/kappital/pkg/utils/operations"
)

// BeforeDelete do some processes before delete the service binding
func (h *Handler) BeforeDelete(obj interface{}) (retry bool, err error) {
	binding := getTypedObj(obj)
	// renew the timeout period at the first time
	if err = updateProcessTimeout(binding, bindingProcessTimeout); err != nil {
		return true, err
	}

	insExists, err := instanceExists(binding)
	if err != nil {
		klog.Errorf("failed to check instance for binding[%s], error: %v", binding.Name, err)
		return true, err
	}

	if insExists {
		return true, fmt.Errorf("[delete binding] instance exists, waiting for next loop")
	}

	return false, nil
}

// Delete the cloud native service instance in cluster
func (h *Handler) Delete(obj interface{}) (retry bool, err error) {
	binding := getTypedObj(obj)
	if err = updateProcessTimeout(binding, bindingProcessTimeout); err != nil {
		return true, err
	}
	if err = deleteResources(binding); err != nil {
		klog.Errorf("failed to delete resource for binding[%s]", binding.Name)
		return true, err
	}

	return deleteRecord(binding)
}

// AfterDelete delete service binding does not need to implement this method
func (h *Handler) AfterDelete(_ interface{}) (bool, error) {
	return false, nil
}

func instanceExists(binding *internals.ServiceBinding) (bool, error) {
	filter := map[string]string{"service_binding_id": binding.ID}

	dbStore := instance.Instance{}
	obj, err := dbStore.GetList(filter)
	if err != nil {
		klog.Errorf("failed to query instances for binding[%s] from db, error: %s", binding.Name, err)
		return true, err
	}

	instances, ok := obj.([]internals.ServiceInstance)
	if !ok {
		klog.Errorf("failed to trans instances for binding[%s] from db, error: %s", binding.Name, err)
		return true, err
	}

	if len(instances) > 0 {
		return true, nil
	}
	return false, nil
}

func deleteResources(binding *internals.ServiceBinding) error {
	gvr := schema.GroupVersionResource{
		Group:    enginev1alpha1.ServicePackageGroupVersionResource.Group,
		Version:  enginev1alpha1.ServicePackageGroupVersionResource.Version,
		Resource: enginev1alpha1.ServicePackageGroupVersionResource.Resource,
	}
	sp, found, err := co.GetClusterOperation().GetServicePackageByName(binding.Name, binding.Namespace)
	if err != nil {
		klog.Errorf("[delete binding] get binding %s resource failed, err: %s", binding.Name, err)
		return err
	}
	if !found {
		return nil
	}
	// delete the service package resources, such as cluster role, service account, and etc.
	sp.Spec.Version = ""
	if err = co.GetClusterOperation().UpdateCustomResource(gvr, binding.Namespace, sp); err != nil {
		return err
	}
	// when all resources have been deleted, delete the service package cr in cluster
	if sp.Status.Phase == enginev1alpha1.DeletingPhase {
		return fmt.Errorf("waiting for resource delete")
	}
	if err := co.GetClusterOperation().DeleteCustomResource(gvr, binding.Name, binding.Namespace); err != nil {
		klog.Errorf("[delete binding] delete binding %s resource failed, err: %s", binding.Name, err)
		return err
	}
	return fmt.Errorf("waiting for resource delete")
}

func deleteRecord(binding *internals.ServiceBinding) (bool, error) {
	dbStore := servicebinding.ServiceBinding{}
	if err := dbStore.Delete(*binding); err != nil {
		return true, err
	}
	return false, nil
}
