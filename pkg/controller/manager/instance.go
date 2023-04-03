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

package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/apis/internals"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/constants"
	"github.com/kappital/kappital/pkg/controller/utils"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/resource"
	"github.com/kappital/kappital/pkg/utils/errors"
	"github.com/kappital/kappital/pkg/utils/uuid"
	"github.com/kappital/kappital/pkg/utils/version"
)

// InstanceController the controller of the instance which deploy, search, and delete the instance
// in database or cluster
type InstanceController struct {
	web.Controller
	instance resource.InstanceResource
	binding  resource.ServiceBindingResource
}

// CreateInstance deploy the custom resource into cluster
func (i *InstanceController) CreateInstance() {
	serviceBinding := i.GetString(constants.ServiceBindingPathParam)
	clusterName := i.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	var err error
	var resourceName string
	defer utils.AuditLog(i.Ctx, "CreateInstance", utils.DeployAction, &resourceName, &err)
	if !utils.ValidString(serviceBinding) || !utils.ValidString(clusterName) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	instanceCreation, err := getAndResolveServiceParam(i.Ctx)
	if err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, errors.ErrServiceParam.WrapErrorReasonWith(err.Error()))
		return
	}
	resourceName = fmt.Sprintf("Deploy Instance for Service Package [%s] to Cluster [%s]",
		serviceBinding, clusterName)
	if err = resource.ValidationInstance(instanceCreation); err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, errors.ErrServiceParam.WrapErrorReasonWith(err.Error()))
		return
	}
	if !i.binding.IsServiceBindingDeployed(serviceBinding, clusterName) {
		subRes, err := transCreationToServiceBinding(instanceCreation)
		if err != nil {
			utils.ReplyJSON(i.Ctx, http.StatusBadRequest, errors.ErrServiceParam.WrapErrorReasonWith(err.Error()))
			return
		}
		if err = i.binding.CreateServiceBinding(*subRes); err != nil {
			utils.ReplyJSON(i.Ctx, http.StatusInternalServerError,
				errors.ErrServiceInstall.WrapErrorReasonWith(err.Error()))
			return
		}
	}
	binding, err := i.binding.GetInternalServiceBinding(serviceBinding, clusterName)
	if err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusInternalServerError,
			errors.ErrDataUnmarshal.WrapErrorReasonWith(err.Error()))
		return
	}
	instances, err := transCreationToServiceInstance(*binding, instanceCreation)
	if err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusInternalServerError,
			errors.ErrDataUnmarshal.WrapErrorReasonWith(err.Error()))
		return
	}
	if err = i.instance.CreateInstance(instances, map[string]string{"name": serviceBinding, "cluster_name": clusterName}); err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusInternalServerError,
			errors.ErrServiceInstanceCreate.WrapErrorReasonWith(err.Error()))
		return
	}
	for _, instance := range instances {
		klog.Infof("create service instance %s success", instance.Name)
	}
	utils.ReplyJSON(i.Ctx, http.StatusOK, "success")
}

// GetInstances get the instance information from database and check does the instance is existed in cluster
func (i *InstanceController) GetInstances() {
	serviceBinding := i.GetString(constants.ServiceBindingPathParam)
	clusterName := i.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	namespace := i.GetString(constants.NamespaceQueryParam, apis.DefaultNamespace)
	var err error
	var resourceName string
	defer utils.AuditLog(i.Ctx, "GetInstances", utils.QueryAction, &resourceName, &err)
	if (!utils.ValidString(serviceBinding) && len(serviceBinding) > 0) || !utils.ValidString(clusterName) || !utils.ValidString(namespace) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Get Instance List of Service Binding [%s] from Namespace [%s] in Cluster [%s]",
		serviceBinding, namespace, clusterName)
	resp, err := i.instance.GetInstances(serviceBinding, clusterName, namespace)
	if err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusInternalServerError, err)
		return
	}
	utils.ReplyJSON(i.Ctx, http.StatusOK, resp)
}

// GetInstanceDetail get the detail information of the instance
func (i *InstanceController) GetInstanceDetail() {
	serviceBinding := i.GetString(constants.ServiceBindingPathParam)
	instanceName := i.GetString(constants.InstancePathParam)
	clusterName := i.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	namespace := i.GetString(constants.NamespaceQueryParam, apis.DefaultNamespace)
	var err error
	var resourceName string
	defer utils.AuditLog(i.Ctx, "GetInstanceDetail", utils.QueryAction, &resourceName, &err)
	if !utils.ValidString(serviceBinding) || !utils.ValidString(instanceName) || !utils.ValidString(clusterName) || !utils.ValidString(namespace) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Get Instance [%s] if Service Binding [%s] from Namespace [%s] in Cluster [%s]",
		instanceName, serviceBinding, namespace, clusterName)
	resp, err := i.instance.GetInstance(serviceBinding, clusterName, namespace, instanceName)
	if err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusInternalServerError, err)
		return
	}
	utils.ReplyJSON(i.Ctx, http.StatusOK, resp)
}

// DeleteInstance delete the instance from cluster and database
func (i *InstanceController) DeleteInstance() {
	serviceBinding := i.GetString(constants.ServiceBindingPathParam)
	instanceName := i.GetString(constants.InstancePathParam)
	clusterName := i.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	namespace := i.GetString(constants.NamespaceQueryParam, apis.DefaultNamespace)
	var err error
	var resourceName string
	defer utils.AuditLog(i.Ctx, "DeleteInstance", utils.UninstallAction, &resourceName, &err)
	if !utils.ValidString(instanceName) || !utils.ValidString(clusterName) || !utils.ValidString(namespace) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Uninstall Service Instance [%s] of Service Binding [%s] from Namespace [%s] in Cluster [%s]",
		instanceName, serviceBinding, namespace, clusterName)
	if err = i.instance.DeleteInstance(clusterName, instanceName, namespace); err != nil {
		utils.ReplyJSON(i.Ctx, http.StatusBadRequest, err)
		return
	}
	utils.ReplyJSON(i.Ctx, http.StatusOK, nil)
}

func transCreationToServiceInstance(binding internals.ServiceBinding,
	serviceBindingReq *instancev1alpha1.ServiceInstanceCreation) ([]internals.ServiceInstance, error) {
	now := time.Now()
	var instances []internals.ServiceInstance
	for _, cr := range serviceBindingReq.InstanceCustomResources {
		crByte, err := json.Marshal(cr)
		if err != nil {
			return nil, err
		}
		instance := internals.ServiceInstance{
			ID:                 uuid.NewUUID(),
			Name:               cr.Name,
			Namespace:          cr.Namespace,
			RawResource:        string(crByte),
			ClusterName:        binding.ClusterName,
			CreateTime:         now,
			Status:             models.StatusInitializing,
			Kind:               cr.Kind,
			APIVersion:         cr.APIVersion,
			ServiceBindingName: binding.Name,
			ServiceBindingID:   binding.ID,
			ServiceName:        binding.ServiceName,
			ServiceID:          binding.ServiceID,
			UpdateTime:         now,
		}
		plural, err := getResourceFromCRD(binding.CRD, instance.Kind, cr.GroupVersionKind().Group)
		if err != nil {
			return nil, err
		}
		instance.Resource = plural
		instances = append(instances, instance)
	}

	return instances, nil
}

func getResourceFromCRD(crds []string, kind, group string) (string, error) {
	v1CRDs, v1beta1CRDs := version.GetCrdV1AndBeta1Slice(crds)
	for _, v1CRD := range v1CRDs {
		if v1CRD.Spec.Names.Kind == kind && v1CRD.Spec.Group == group {
			return v1CRD.Spec.Names.Plural, nil
		}
	}

	for _, v1beta1CRD := range v1beta1CRDs {
		if v1beta1CRD.Spec.Names.Kind == kind && v1beta1CRD.Spec.Group == group {
			return v1beta1CRD.Spec.Names.Plural, nil
		}
	}

	return "", fmt.Errorf("can not find cr plural because the cr is not fit CRD")
}
