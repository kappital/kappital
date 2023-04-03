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
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/beego/beego/v2/client/orm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/dao/instance"
	"github.com/kappital/kappital/pkg/dao/servicebinding"
	"github.com/kappital/kappital/pkg/models"
	mo "github.com/kappital/kappital/pkg/models/operation"
	"github.com/kappital/kappital/pkg/utils/operations"
	"github.com/kappital/kappital/pkg/watcher"
)

// ServiceBindingResource operate service information in database and/or cluster
type ServiceBindingResource struct {
	bindingDao  servicebinding.ServiceBinding
	instanceDao instance.Instance
	mo.ServiceBindingOperation
}

// DeleteServiceInCluster delete the service in cluster
func (s *ServiceBindingResource) DeleteServiceInCluster(_, _ string) (string, error) {
	return "", nil
}

// GetResourceType get the resource type, and it is service binding
func (s *ServiceBindingResource) GetResourceType() Type {
	return ServiceBindingType
}

// GetCommonDBObject get the service binding information from the database
func (s *ServiceBindingResource) GetCommonDBObject(pk string) (interface{}, error) {
	obj, err := s.bindingDao.GetByPrimaryKey(pk)
	if err != nil {
		klog.Errorf("failed to get service binding, error: %v", err)
		return nil, err
	}
	ins, ok := obj.(internals.ServiceBinding)
	if !ok {
		klog.Errorf("failed to trans service binding, error: %v", err)
		return nil, err
	}

	return &ins, nil
}

// GetObjectProcessTime get the service binding process timestamp
func (s *ServiceBindingResource) GetObjectProcessTime(obj interface{}) time.Time {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.servicebinding, actual :%s", reflect.TypeOf(obj).Name())
		return time.Time{}
	}
	return binding.ProcessTime
}

// GetObjectStatus get the ServiceBinding status
func (s *ServiceBindingResource) GetObjectStatus(obj interface{}) string {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.servicebinding, actual :%s", reflect.TypeOf(obj).Name())
		return ""
	}
	return binding.Status
}

// GetObjectID get the service binding's id
func (s *ServiceBindingResource) GetObjectID(obj interface{}) string {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.servicebinding, actual :%s", reflect.TypeOf(obj).Name())
		return ""
	}
	return binding.ID
}

// GetObjectListByStatusSets get the object list
func (s *ServiceBindingResource) GetObjectListByStatusSets(status sets.String) ([]interface{}, error) {
	tmp, err := s.bindingDao.GetListByStatusSets(status)
	if err != nil {
		klog.Errorf("get status %s by instance failed, error: %s", status, err)
		return nil, err
	}
	bindings, ok := tmp.([]internals.ServiceBinding)
	if !ok {
		return nil, fmt.Errorf("cannot get service binding list beacuase not same type")
	}
	result := make([]interface{}, 0, len(bindings))
	for i := range bindings {
		result = append(result, &bindings[i])
	}
	return result, nil
}

// UpdateProcessFailed update the object to failed status
func (s *ServiceBindingResource) UpdateProcessFailed(obj interface{}, status string, msg string) error {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.servicebinding, actual :%s", reflect.TypeOf(obj).Name())
		return fmt.Errorf("invalid object type, expected: internals.servicebinding, actual :%s",
			reflect.TypeOf(obj).Name())
	}

	return s.bindingDao.UpdateStatusMsg(*binding, status, msg)
}

// UpdateObjProcessTime update the process timestamp
func (s *ServiceBindingResource) UpdateObjProcessTime(obj interface{}, processTime time.Time) error {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.servicebinding, actual :%s", reflect.TypeOf(obj).Name())
		return fmt.Errorf("invalid object type, expected: internals.servicebinding, actual :%s",
			reflect.TypeOf(obj).Name())
	}

	binding.ProcessTime = processTime
	return s.bindingDao.Update(*binding)
}

// GetObjUpdateTime get the object update timestamp
func (s *ServiceBindingResource) GetObjUpdateTime(obj interface{}) time.Time {
	binding, ok := obj.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type, expected: internals.operator, actual :%s", reflect.TypeOf(obj).Name())
		return time.Time{}
	}

	return binding.UpdateTime
}

// CreateServiceBinding create the service binding into cluster and insert the record to the database
func (s *ServiceBindingResource) CreateServiceBinding(serviceBinding internals.ServiceBinding) error {
	klog.Infof("create service binding %s", serviceBinding.Name)
	_, err := s.bindingDao.Get(map[string]string{"name": serviceBinding.Name,
		"cluster_name": serviceBinding.ClusterName})
	if err == nil {
		klog.Infof("service binding %s has been created", serviceBinding.Name)
		return nil
	}

	if !errors.Is(err, orm.ErrNoRows) {
		klog.Infof("get binding %s failed", serviceBinding.Name)
		return err
	}

	serviceBinding.Status = models.StatusInstalling
	if err = s.bindingDao.Create(serviceBinding, map[string]string{}); err != nil {
		return err
	}

	if err = watcher.AddEvent(serviceBinding, watcher.OPCreate, apis.OperatorProcessor); err != nil {
		klog.Infof("service binding %s add watcher event failed", serviceBinding.Name)
		return err
	}

	klog.Infof("service binding %s has been created", serviceBinding.Name)
	return nil
}

// IsServiceBindingDeployed does the service binding has already deployed to the target cluster
func (s *ServiceBindingResource) IsServiceBindingDeployed(name, cluster string) bool {
	return s.IsExist(map[string]string{"name": name, "cluster_name": cluster})
}

// DeleteServiceBinding use the service binding name and cluster name to delete the service binding
func (s *ServiceBindingResource) DeleteServiceBinding(bindingName, clusterName string) error {
	filter := map[string]string{
		"name":         bindingName,
		"cluster_name": clusterName,
	}
	obj, err := s.bindingDao.Get(filter)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil
		}
		return err
	}
	binding, ok := obj.(internals.ServiceBinding)
	if !ok {
		return fmt.Errorf("get binding %s cluster %s to binding failed", bindingName, clusterName)
	}

	objInstances, err := s.instanceDao.GetList(map[string]string{"service_binding_id": binding.ID})
	if err != nil {
		return err
	}

	if objInstances != nil {
		instances, ok := objInstances.([]internals.ServiceInstance)
		if !ok {
			return fmt.Errorf("get binding %s cluster %s instances failed", bindingName, clusterName)
		}
		for _, instanceObj := range instances {
			instanceObj.Status = models.StatusDeleting
			if err = s.instanceDao.Update(instanceObj, "status"); err != nil {
				return err
			}

			if err = watcher.AddEvent(instanceObj, watcher.OPDelete, apis.InstanceProcessor); err != nil {
				return err
			}
		}
	}

	binding.Status = models.StatusDeleting
	if err = s.bindingDao.Update(binding, "status"); err != nil {
		return err
	}
	return watcher.AddEvent(binding, watcher.OPDelete, apis.OperatorProcessor)
}

// GetInternalServiceBinding get the ServiceBinding as the internal format
func (s *ServiceBindingResource) GetInternalServiceBinding(serviceBindingName string,
	clusterName string) (*internals.ServiceBinding, error) {
	tmp, err := s.bindingDao.Get(map[string]string{"name": serviceBindingName, "cluster_name": clusterName})
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	sb, ok := tmp.(internals.ServiceBinding)
	if !ok {
		return nil, fmt.Errorf("cannot get ServiceBindingModel because not this type")
	}
	return &sb, nil
}

// GetServiceBindings get service bindings of this cluster
func (s *ServiceBindingResource) GetServiceBindings(clusterName string) ([]instancev1alpha1.CloudNativeServiceInstance, error) {
	sbs, err := s.GetList(map[string]string{"cluster_name": clusterName})
	if err != nil {
		return nil, err
	}
	return s.transModelSliceToResponse(sbs.([]models.ServiceBindingModel)), nil
}

// GetServiceBinding use name, clusterId to get the service binding information
func (s *ServiceBindingResource) GetServiceBinding(name, clusterName string,
	detail bool) (*instancev1alpha1.CloudNativeServiceInstance, error) {
	sb, err := s.Get(map[string]string{"name": name, "cluster_name": clusterName})
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	engine, _, err := operations.GetClusterOperation().GetServicePackageByName(name, apis.KappitalSystemNamespace)
	if err != nil {
		return nil, err
	}
	result := s.transferServiceBindingModelToResponse(sb.(models.ServiceBindingModel), &engine, detail)
	return &result, nil
}

func (s *ServiceBindingResource) transferServiceBindingModelToResponse(
	item models.ServiceBindingModel, engine *enginev1alpha1.ServicePackage,
	detail bool) instancev1alpha1.CloudNativeServiceInstance {
	binding := s.getServiceBindingMeta(item)
	if engine != nil {
		binding.Spec.ServiceReference = instancev1alpha1.ServiceReference{
			TypeMeta: metav1.TypeMeta{
				Kind:       engine.Kind,
				APIVersion: engine.APIVersion,
			},
			Name:      engine.Name,
			Namespace: engine.Namespace,
			UID:       string(engine.UID),
			Status:    engine.Status.Phase,
		}
	}
	if detail {
		notFound, pending := 0, 0
		binding.Spec.CustomResources, notFound, pending = s.transferResourceModelToResponse(item.Resources)
		binding.Status.Phase, binding.Status.Message = dealWithExceptionCount(notFound, pending,
			len(item.Resources))
	}
	return binding
}

func (s *ServiceBindingResource) transferResourceModelToResponse(
	resources []*models.ResourceModel) ([]instancev1alpha1.Resource, int, int) {
	var resp []instancev1alpha1.Resource
	notFound, pending := 0, 0
	for _, resource := range resources {
		res, currNotFound, currPending := s.transferInstanceModelToResponse(resource.Instances)
		resp = append(resp, res...)
		notFound += currNotFound
		pending += currPending
	}
	return resp, notFound, pending
}

func (s *ServiceBindingResource) transferInstanceModelToResponse(
	instances []*models.InstanceModel) ([]instancev1alpha1.Resource, int, int) {
	resource := make([]instancev1alpha1.Resource, 0, len(instances))
	notFound, pending := 0, 0
	for _, m := range instances {
		resource = append(resource, instancev1alpha1.Resource{
			TypeMeta: metav1.TypeMeta{
				Kind:       m.Kind,
				APIVersion: m.APIVersion,
			},
			Name:      m.Name,
			Namespace: m.Namespace,
			UID:       m.ID,
			Status:    m.Status,
		})
	}
	return resource, notFound, pending
}

func (s *ServiceBindingResource) transModelSliceToResponse(
	items []models.ServiceBindingModel) []instancev1alpha1.CloudNativeServiceInstance {
	sis := make([]instancev1alpha1.CloudNativeServiceInstance, 0, len(items))
	for _, item := range items {
		sis = append(sis, s.transferServiceBindingModelToResponse(item, nil, true))
	}
	return sis
}

func (s *ServiceBindingResource) getServiceBindingMeta(item models.ServiceBindingModel) instancev1alpha1.CloudNativeServiceInstance {
	return instancev1alpha1.CloudNativeServiceInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       apis.CloudNativeServiceInstanceKind,
			APIVersion: apis.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              item.Name,
			Namespace:         item.Namespace,
			CreationTimestamp: metav1.Time{Time: item.CreateTime},
		},
		Spec: instancev1alpha1.CloudNativeServiceInstanceSpec{
			Name:        item.Name,
			Version:     item.Version,
			ServiceName: item.ServiceName,
			ServiceID:   item.ServiceID,
			ClusterName: item.ClusterName,
		},
	}
}

func dealWithExceptionCount(notFoundCount, pendingCount, total int) (instancev1alpha1.Phase, string) {
	if total == 0 {
		return instancev1alpha1.SucceededPhase, ""
	}
	if notFoundCount == total {
		return instancev1alpha1.FailedPhase, fmt.Sprintf("All Custom Resources (%d) are not found in cluster. Please "+
			"delete the Custom Resources by manager, DO NOT manual delete them in cluster. If you want to fix up "+
			"the status for manager, Please clean up the relationship data in database.", notFoundCount)
	}
	if pendingCount == total {
		return instancev1alpha1.PendingPhase, fmt.Sprintf("All Custom Resources (%d) are not deploy into cluster. "+
			"Please wait a few miniutes.", pendingCount)
	}
	phase := instancev1alpha1.SucceededPhase
	message := ""
	if notFoundCount > 0 && pendingCount > 0 {
		phase = instancev1alpha1.UnknownPhase
		message = fmt.Sprintf("Unknown reasons make the Cloud Native Service InstanceModel may not provide the "+
			"service. There have %d custom resource(s) cannot find in cluster, and %d custom resource(s) still "+
			"during the installing status [Total have %d custom resource(s)].", notFoundCount, pendingCount, total)
	} else if notFoundCount > 0 {
		phase = instancev1alpha1.UnknownPhase
		message = fmt.Sprintf("There have %d custom resource(s) cannot find in cluster. [Total have %d custom "+
			"resource(s)]", notFoundCount, total)
	} else if pendingCount > 0 {
		phase = instancev1alpha1.PendingPhase
		message = fmt.Sprintf("There have %d custom resource(s) still installing. [Total have %d custom "+
			"resources(s)]", pendingCount, total)
	}
	return phase, message
}
