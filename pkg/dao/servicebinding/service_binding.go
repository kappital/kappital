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
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/models"
	mo "github.com/kappital/kappital/pkg/models/operation"
	"github.com/kappital/kappital/pkg/utils/uuid"
	"github.com/kappital/kappital/pkg/utils/version"
)

// ServiceBinding the dao layer of service binding for database CRUD
type ServiceBinding struct {
	db mo.ServiceBindingOperation
}

// Create insert a data record to the database
func (s ServiceBinding) Create(obj interface{}, _ map[string]string) error {
	serviceBinding, ok := obj.(internals.ServiceBinding)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	binding, err := transServiceBinding2Model(serviceBinding)
	if err != nil {
		return err
	}
	var resources []*models.ResourceModel
	v1CRDs, v1beta1CRDs := version.GetCrdV1AndBeta1Slice(serviceBinding.CRD)
	if len(v1CRDs) != 0 {
		for _, v1CRD := range v1CRDs {
			resource := &models.ResourceModel{
				ID:              uuid.NewUUID(),
				Kind:            v1CRD.Spec.Names.Kind,
				Group:           v1CRD.Spec.Group,
				APIVersion:      fmt.Sprintf("%s/%s", v1CRD.Spec.Group, v1CRD.Spec.Versions[0].Name),
				Resource:        v1CRD.Spec.Names.Plural,
				CreateTimestamp: binding.CreateTime,
				UpdateTimestamp: binding.UpdateTime,
			}
			resources = append(resources, resource)
		}
	}

	if len(v1beta1CRDs) != 0 {
		for _, v1beta1CRD := range v1beta1CRDs {
			resource := &models.ResourceModel{
				ID:              uuid.NewUUID(),
				Kind:            v1beta1CRD.Spec.Names.Kind,
				Group:           v1beta1CRD.Spec.Group,
				APIVersion:      fmt.Sprintf("%s/%s", v1beta1CRD.Spec.Group, v1beta1CRD.Spec.Versions[0].Name),
				Resource:        v1beta1CRD.Spec.Names.Plural,
				CreateTimestamp: binding.CreateTime,
				UpdateTimestamp: binding.UpdateTime,
			}
			resources = append(resources, resource)
		}
	}

	binding.Resources = resources
	return s.db.Insert(binding)
}

// Get the service binding from database and filter by cols
func (s ServiceBinding) Get(cols map[string]string) (interface{}, error) {
	result, err := s.db.Get(cols)
	if err != nil {
		return nil, err
	}
	return transModel2ServiceBinding(result.(models.ServiceBindingModel))
}

// GetByPrimaryKey get the service binding by the primary key (id)
func (s ServiceBinding) GetByPrimaryKey(id string) (interface{}, error) {
	obj, err := s.db.GetByPrimaryKey(id)
	if err != nil {
		return nil, err
	}
	binding, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return nil, fmt.Errorf("can not trans model trans to binding")
	}
	return transModel2ServiceBinding(binding)
}

// GetList get service binding list, and filter by cols
func (s ServiceBinding) GetList(cols map[string]string) (interface{}, error) {
	items, err := s.db.GetList(cols)
	if err != nil {
		return nil, err
	}
	result, err := transModelSlice2ServiceBindingSlice(items.([]models.ServiceBindingModel))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetListByStatusSets get service binding list, and filters by status
func (s ServiceBinding) GetListByStatusSets(status sets.String) (interface{}, error) {
	filter := map[string][]interface{}{
		"status__in": {status.UnsortedList()},
	}
	obj, err := s.db.GetListByFilter(filter)
	if err != nil {
		return nil, err
	}

	bindings, ok := obj.([]models.ServiceBindingModel)
	if !ok {
		return nil, fmt.Errorf("can not trans model trans to binding")
	}

	serviceBindings := make([]internals.ServiceBinding, 0, len(bindings))
	for _, binding := range bindings {
		serviceBinding, err := transModel2ServiceBinding(binding)
		if err != nil {
			return nil, err
		}
		serviceBindings = append(serviceBindings, serviceBinding)
	}

	return serviceBindings, nil
}

// Update the service binding to the database with cols
func (s ServiceBinding) Update(obj interface{}, cols ...string) error {
	binding, ok := obj.(internals.ServiceBinding)
	if !ok {
		klog.Errorf("obj type is not ServiceBindingModel, actual: %s", reflect.TypeOf(obj).Name())
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	bindingModel, err := transServiceBinding2Model(binding)
	if err != nil {
		return err
	}

	return s.db.Update(bindingModel, cols...)
}

// UpdateStatusMsg update the status massage for service binding
func (s ServiceBinding) UpdateStatusMsg(obj interface{}, status, msg string) error {
	binding, ok := obj.(internals.ServiceBinding)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	binding.Status = status
	binding.Message = msg
	bindingModel, err := transServiceBinding2Model(binding)
	if err != nil {
		return err
	}

	return s.db.Update(bindingModel)
}

// Delete the service binding
func (s ServiceBinding) Delete(obj interface{}) error {
	binding, ok := obj.(internals.ServiceBinding)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	bindingModel, err := transServiceBinding2Model(binding)
	if err != nil {
		return err
	}

	return s.db.Delete(bindingModel)
}

func transServiceBinding2Model(serviceBinding internals.ServiceBinding) (models.ServiceBindingModel, error) {
	now := time.Now().UTC()
	binding := models.ServiceBindingModel{
		ID:           serviceBinding.ID,
		Name:         serviceBinding.Name,
		Namespace:    serviceBinding.Namespace,
		ServiceName:  serviceBinding.ServiceName,
		Version:      serviceBinding.Version,
		ClusterName:  serviceBinding.ClusterName,
		ServiceID:    serviceBinding.ServiceID,
		Status:       serviceBinding.Status,
		ErrorMessage: serviceBinding.Message,
		CreateTime:   now,
		UpdateTime:   now,
	}

	if len(serviceBinding.Status) == 0 {
		binding.Status = models.StatusInstalling
	}

	workloadByte, err := json.Marshal(serviceBinding.Workload)
	if err != nil {
		return models.ServiceBindingModel{}, err
	}

	permissions, err := json.Marshal(serviceBinding.Permissions)
	if err != nil {
		return models.ServiceBindingModel{}, err
	}

	capabilityPluginByte, err := json.Marshal(serviceBinding.CapabilityPlugin)
	if err != nil {
		return models.ServiceBindingModel{}, err
	}

	crdsByte, err := json.Marshal(serviceBinding.CRD)
	if err != nil {
		return models.ServiceBindingModel{}, err
	}

	binding.Workloads = string(workloadByte)
	binding.Permissions = string(permissions)
	binding.CapabilityPlugin = string(capabilityPluginByte)
	binding.CustomResourceDefinition = string(crdsByte)

	return binding, nil
}

func transModel2ServiceBinding(model models.ServiceBindingModel) (internals.ServiceBinding, error) {
	serviceBinding := internals.ServiceBinding{
		ID:          model.ID,
		Name:        model.Name,
		Version:     model.Version,
		Namespace:   model.Namespace,
		ServiceName: model.ServiceName,
		ServiceID:   model.ServiceID,
		ClusterName: model.ClusterName,
		Status:      model.Status,
		Message:     model.ErrorMessage,
		ProcessTime: model.ProcessTime,
		UpdateTime:  model.UpdateTime,
	}

	var crd []string
	if err := json.Unmarshal([]byte(model.CustomResourceDefinition), &crd); err != nil {
		klog.Errorf("json Unmarshal crd string to []string failed, err: %s", err)
		return internals.ServiceBinding{}, err
	}
	serviceBinding.CRD = crd

	var permissions []enginev1alpha1.Permission
	if err := json.Unmarshal([]byte(model.Permissions), &permissions); err != nil {
		klog.Errorf("json Unmarshal string to permissions struct failed, err: %s", err)
		return internals.ServiceBinding{}, err
	}
	serviceBinding.Permissions = permissions

	var workload enginev1alpha1.Workload
	if err := json.Unmarshal([]byte(model.Workloads), &workload); err != nil {
		klog.Errorf("json Unmarshal string to workload struct failed, err: %s", err)
		return internals.ServiceBinding{}, err
	}
	serviceBinding.Workload = workload

	var capabilityPlugin enginev1alpha1.CapabilityPlugin
	if err := json.Unmarshal([]byte(model.CapabilityPlugin), &capabilityPlugin); err != nil {
		klog.Errorf("json Unmarshal string to capabilityPlugin struct failed, err: %s", err)
		return internals.ServiceBinding{}, err
	}
	serviceBinding.CapabilityPlugin = capabilityPlugin

	return serviceBinding, nil
}

func transModelSlice2ServiceBindingSlice(bindings []models.ServiceBindingModel) ([]internals.ServiceBinding, error) {
	result := make([]internals.ServiceBinding, 0, len(bindings))
	for _, binding := range bindings {
		curr, err := transModel2ServiceBinding(binding)
		if err != nil {
			return nil, err
		}
		result = append(result, curr)
	}
	return result, nil
}
