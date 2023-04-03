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
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/constants"
	"github.com/kappital/kappital/pkg/controller/utils"
	"github.com/kappital/kappital/pkg/resource"
	"github.com/kappital/kappital/pkg/utils/errors"
	"github.com/kappital/kappital/pkg/utils/validation"
)

// ServiceBindingController the controller of the service binding which deploy, search, and delete the service binding
// in database or cluster
type ServiceBindingController struct {
	web.Controller
	resource resource.ServiceBindingResource
}

// CreateServiceBinding deploy the service package into cluster
func (s *ServiceBindingController) CreateServiceBinding() {
	serviceBody, err := getAndResolveServiceParam(s.Ctx)
	var resourceName string
	defer utils.AuditLog(s.Ctx, "CreateServiceBinding", utils.DeployAction, &resourceName, &err)
	if err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, errors.ErrServiceParam.WrapErrorReasonWith(err.Error()))
		return
	}
	resourceName = fmt.Sprintf("Deploy Service [%s] into Cluster [%s]", serviceBody.Service.Name, serviceBody.ClusterID)
	sb, err := s.resource.GetInternalServiceBinding(serviceBody.Service.Spec.Description.Name, serviceBody.ClusterID)
	if err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusInternalServerError,
			errors.ErrServiceInstall.WrapErrorReasonWith(err.Error()))
		return
	}
	if sb != nil {
		utils.ReplyJSON(s.Ctx, http.StatusOK, map[string]string{"Name": sb.Name, "ID": sb.ID})
		return
	}

	subRes, err := transCreationToServiceBinding(serviceBody)
	if err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, errors.ErrServiceParam.WrapErrorReasonWith(err.Error()))
		return
	}

	if err = s.resource.CreateServiceBinding(*subRes); err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusInternalServerError,
			errors.ErrServiceInstall.WrapErrorReasonWith(err.Error()))
		return
	}

	klog.Infof("service %s serviceBinding create in cluster %s success.", subRes.Name, subRes.ClusterID)
	utils.ReplyJSON(s.Ctx, http.StatusOK, map[string]string{"Name": subRes.Name, "ID": subRes.ID})
}

// DeleteServiceBinding destroy the service binding from cluster
func (s *ServiceBindingController) DeleteServiceBinding() {
	serviceBinding := s.GetString(constants.ServiceBindingPathParam)
	clusterName := s.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	var err error
	var resourceName string
	defer utils.AuditLog(s.Ctx, "DeleteServiceBinding", utils.UninstallAction, &resourceName, &err)
	if !utils.ValidString(serviceBinding) || !utils.ValidString(clusterName) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Uninstall Service Binding [%s] in Cluster [%s]", serviceBinding, clusterName)
	if err = s.resource.DeleteServiceBinding(serviceBinding, clusterName); err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusInternalServerError, errors.ErrServiceDelete.WrapErrorReasonWith(err.Error()))
		return
	}
	utils.ReplyJSON(s.Ctx, http.StatusOK, nil)
}

// GetServiceBindings get the service bindings' information from database and cluster
func (s *ServiceBindingController) GetServiceBindings() {
	clusterName := s.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	var err error
	var resourceName string
	defer utils.AuditLog(s.Ctx, "GetServiceBindings", utils.QueryAction, &resourceName, &err)
	if !utils.ValidString(clusterName) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Get Service Binding List from Cluster [%s]", clusterName)
	sis, err := s.resource.GetServiceBindings(clusterName)
	if err != nil {
		utils.ReplyJSON(s.Ctx, http.StatusInternalServerError, err)
		return
	}
	utils.ReplyJSON(s.Ctx, http.StatusOK, sis)
}

// GetServiceBindingDetail get the service binding detail information
func (s *ServiceBindingController) GetServiceBindingDetail() {
	serviceBinding := s.GetString(constants.ServiceBindingPathParam)
	clusterName := s.GetString(constants.ClusterNameQueryParam, apis.DefaultCluster)
	var err error
	var resourceName string
	defer utils.AuditLog(s.Ctx, "GetServiceBindingDetail", utils.QueryAction, &resourceName, &err)
	if !utils.ValidString(serviceBinding) || !utils.ValidString(clusterName) {
		err = utils.ErrIllegalParameters
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, utils.ErrIllegalParameters)
		return
	}
	resourceName = fmt.Sprintf("Get Service Binding [%s] from Cluster [%s]", serviceBinding, clusterName)
	detail, err := validation.ValidBool(s.Ctx.Input.Query(constants.Detail))
	if err != nil {
		klog.Warningf("cannot get parameter value of detail, will set it to the false")
	}
	si, err := s.resource.GetServiceBinding(serviceBinding, clusterName, detail)
	if err != nil || si == nil {
		err = fmt.Errorf("cannot get the service binding")
		utils.ReplyJSON(s.Ctx, http.StatusBadRequest, err)
		return
	}
	utils.ReplyJSON(s.Ctx, http.StatusOK, si)
}
