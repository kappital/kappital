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
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/models"
	mo "github.com/kappital/kappital/pkg/models/operation"
)

// InstallConditionType of operator or custom resource
type InstallConditionType string

// ConditionStatus of resource status
type ConditionStatus string

const (
	// InstallOperator the condition type of operator
	InstallOperator InstallConditionType = "InstallOperator"
	// CreateResource the condition type of custom resource
	CreateResource InstallConditionType = "CreateResource"

	// Waiting condition status for resources which waiting for install
	Waiting ConditionStatus = "Waiting"
	// Running condition status for resources
	Running ConditionStatus = "Running"
	// Success condition status for resources which has deployed succeeded
	Success ConditionStatus = "Success"
	// Failed condition status for resources which deploy failed or upgrade failed
	Failed ConditionStatus = "Failed"
)

// Instance the dao layer of instance for database CRUD
type Instance struct {
	binding  mo.ServiceBindingOperation
	instance mo.InstanceOperation
}

// Create Insert the ServiceInstance into database
// params is a map of the necessary values, such as ServiceBinding's name, and the ClusterId
func (i Instance) Create(obj interface{}, params map[string]string) error {
	// get the whole service binding object
	tmp, err := i.binding.GetDetail(params)
	if err != nil {
		return err
	}
	binding, ok := tmp.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("the obj is not ServiceBindingModel")
	}
	// get this service binding's resourceId map, using for check does the instance cr has the crd in database
	// and use resource id add into instance model as foreign key
	resourceIDMap := binding.GetResourceIDMap()
	// begin the transaction for database
	now := time.Now().UTC()
	tx := models.NewTransaction(models.GetNewOrm())
	if err = tx.BeginTransaction(); err != nil {
		return err
	}
	defer models.Handler(&err, tx)
	instances, ok := obj.([]internals.ServiceInstance)
	if !ok {
		return fmt.Errorf("obj type is not Slice of ServiceInstance")
	}
	for _, instance := range instances {
		key := fmt.Sprintf("%s;%s;%s", instance.Kind, instance.APIVersion, instance.Resource)
		id, find := resourceIDMap[key]
		if !find {
			return fmt.Errorf("the instance [%s] does not exist its CRD in database", instance.Name)
		}
		model := models.InstanceModel{
			ID:               instance.ID,
			Kind:             instance.Kind,
			APIVersion:       instance.APIVersion,
			Name:             instance.Name,
			Namespace:        instance.Namespace,
			ServiceBindingID: instance.ServiceBindingID,
			RawResource:      instance.RawResource,
			Status:           models.StatusInitializing,
			ErrorMessage:     "",
			ServiceID:        instance.ServiceID,
			ServiceName:      instance.ServiceName,
			ClusterID:        instance.ClusterID,
			ClusterName:      instance.ClusterName,
			CreateTimestamp:  now,
			UpdateTime:       now,
			Resource:         &models.ResourceModel{ID: id},
		}
		if err = i.instance.InsertTx(model, tx.GetTransaction()); err != nil {
			return err
		}
	}
	return nil
}

// Get the instance from database and filter by cols
func (i Instance) Get(cols map[string]string) (interface{}, error) {
	obj, err := i.instance.GetDetail(cols)
	if err != nil {
		return nil, err
	}
	return transModel2Instance(obj.(models.InstanceModel))
}

// GetByPrimaryKey get the instance by the primary key (id)
func (i Instance) GetByPrimaryKey(id string) (interface{}, error) {
	obj, err := i.instance.GetByPrimaryKey(id)
	if err != nil {
		return nil, err
	}
	return transModel2Instance(obj.(models.InstanceModel))
}

// GetList get instance list, and filter by cols
func (i Instance) GetList(cols map[string]string) (interface{}, error) {
	items, err := i.instance.GetList(cols)
	if err != nil {
		return nil, err
	}
	result, err := transModelSlice2InstanceSlice(items.([]models.InstanceModel))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetListByStatusSets get instance list, and filters by status
func (i Instance) GetListByStatusSets(status sets.String) (interface{}, error) {
	filter := map[string][]interface{}{
		"status__in": {status.UnsortedList()},
	}

	obj, err := i.instance.GetListByFilter(filter)
	if err != nil {
		return nil, err
	}

	instances, ok := obj.([]models.InstanceModel)
	if !ok {
		return nil, fmt.Errorf("can not trans model trans to instance")
	}

	serviceInstances := make([]internals.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		serviceInstance, err := transModel2Instance(instance)
		if err != nil {
			return nil, err
		}
		serviceInstances = append(serviceInstances, serviceInstance)
	}

	return serviceInstances, nil
}

// Update the instance to the database with cols
func (i Instance) Update(obj interface{}, cols ...string) error {
	internal, ok := obj.(internals.ServiceInstance)
	if !ok {
		return fmt.Errorf("obj type is not ServiceInstance")
	}
	instance, err := transformInstanceToModel(internal)
	if err != nil {
		return err
	}
	return i.instance.Update(instance, cols...)
}

// UpdateStatusMsg update the status massage for instance
func (i Instance) UpdateStatusMsg(obj interface{}, status, msg string) error {
	instance, ok := obj.(internals.ServiceInstance)
	if !ok {
		return fmt.Errorf("obj type is not ServiceInstance")
	}
	var innerMsg string
	for idx, cond := range instance.InstallState.SubPhase {
		if cond.Status == string(Running) {
			cond.Status = string(Failed)
			cond.LastTransitionTime = metav1.Now()
			instance.InstallState.SubPhase[idx] = cond
			if cond.Message == "" {
				cond.Message = fmt.Sprintf("timed out to do %s", cond.Type)
			}
			innerMsg = cond.Message
		}
	}
	instance.InstallState.Phase = status
	instance.Status = status
	instance.ProcessTime = time.Time{}
	instance.Message = msg
	if innerMsg != "" {
		instance.Message += fmt.Sprintf(":%s", innerMsg)
	}

	ins, err := transformInstanceToModel(instance)
	if err != nil {
		return err
	}

	return i.instance.Update(ins, "install_state", "status", "process_time", "error_message")
}

// Delete the instance
func (i Instance) Delete(obj interface{}) error {
	instance, ok := obj.(internals.ServiceInstance)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	instanceModel, err := transformInstanceToModel(instance)
	if err != nil {
		return err
	}

	return i.instance.Delete(instanceModel)
}

func transformInstanceToModel(ins internals.ServiceInstance) (models.InstanceModel, error) {
	installPhase, err := json.Marshal(ins.InstallState)
	if err != nil {
		return models.InstanceModel{}, err
	}

	return models.InstanceModel{
		ID:                  ins.ID,
		Kind:                ins.Kind,
		APIVersion:          ins.APIVersion,
		Name:                ins.Name,
		Namespace:           ins.Namespace,
		ServiceBindingID:    ins.ServiceBindingID,
		RawResource:         ins.RawResource,
		Status:              ins.Status,
		ErrorMessage:        ins.Message,
		ServiceID:           ins.ServiceID,
		InstanceType:        ins.InstanceType,
		ServiceName:         ins.ServiceName,
		ClusterID:           ins.ClusterID,
		ClusterName:         ins.ClusterName,
		PackageDependencies: "",
		CreateTimestamp:     ins.CreateTime,
		ProcessTime:         ins.ProcessTime,
		UpdateTime:          ins.UpdateTime,
		InstallState:        string(installPhase),
	}, nil
}

func transModel2Instance(instance models.InstanceModel) (internals.ServiceInstance, error) {
	var installPhase internals.InstallState
	if instance.InstallState != "" {
		if err := json.Unmarshal([]byte(instance.InstallState), &installPhase); err != nil {
			return internals.ServiceInstance{}, err
		}
	}

	resourceOperation := mo.ResourceOperation{}
	obj, err := resourceOperation.GetByPrimaryKey(instance.Resource.ID)
	if err != nil {
		return internals.ServiceInstance{}, err
	}
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return internals.ServiceInstance{}, fmt.Errorf("obj type is not Slice of resource")
	}

	return internals.ServiceInstance{
		ID:                 instance.ID,
		Name:               instance.Name,
		Namespace:          instance.Namespace,
		ClusterID:          instance.ClusterID,
		ClusterName:        instance.ClusterName,
		RawResource:        instance.RawResource,
		Resource:           resource.Resource,
		CreateTime:         instance.CreateTimestamp,
		Status:             instance.Status,
		Kind:               instance.Kind,
		APIVersion:         instance.APIVersion,
		Message:            instance.ErrorMessage,
		ServiceBindingName: instance.ServiceName,
		ServiceBindingID:   instance.ServiceBindingID,
		ServiceName:        instance.ServiceName,
		ServiceID:          instance.ServiceID,
		UpdateTime:         instance.UpdateTime,
		ProcessTime:        instance.ProcessTime,
		InstallState:       installPhase,
	}, nil
}

func transModelSlice2InstanceSlice(instances []models.InstanceModel) ([]internals.ServiceInstance, error) {
	result := make([]internals.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		curr, err := transModel2Instance(instance)
		if err != nil {
			return nil, err
		}
		result = append(result, curr)
	}
	return result, nil
}
