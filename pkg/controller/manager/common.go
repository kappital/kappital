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
	"time"

	"github.com/beego/beego/v2/server/web/context"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/utils/uuid"
)

func getAndResolveServiceParam(ctx *context.Context) (*instancev1alpha1.ServiceInstanceCreation, error) {
	var creation = &instancev1alpha1.ServiceInstanceCreation{}
	if ctx.Input.RequestBody != nil && len(ctx.Input.RequestBody) != 0 {
		if err := json.Unmarshal(ctx.Input.RequestBody, creation); err != nil {
			return nil, err
		}
	}

	if err := creation.Validate(); err != nil {
		return creation, err
	}
	return creation, nil
}

func transCreationToServiceBinding(sbReq *instancev1alpha1.ServiceInstanceCreation) (*internals.ServiceBinding, error) {
	if sbReq == nil {
		return nil, fmt.Errorf("the pass in variable is empty, cannot translate to the Service Binding")
	}
	var err error
	now := time.Now()
	serviceBinding := &internals.ServiceBinding{
		ID:          uuid.NewUUID(),
		Name:        sbReq.Service.Spec.Description.Name,
		Version:     sbReq.Service.Spec.Version,
		Namespace:   apis.KappitalSystemNamespace,
		ServiceName: sbReq.Service.Spec.Description.Name,
		ServiceID:   string(sbReq.Service.UID),
		ClusterID:   sbReq.ClusterID,
		ClusterName: sbReq.ClusterID,
		ProcessTime: now.UTC(),
		UpdateTime:  now.UTC(),
		Permissions: servicePermissionBuilder(sbReq.Service.Spec.Operator.ClusterRoles,
			sbReq.Service.Spec.Operator.ClusterRoleBindings),
		CapabilityPlugin: serviceCapabilityPluginBuilder(sbReq.Service.Spec.Manifests),
	}
	serviceBinding.CRD, err = serviceCRDBuilder(sbReq.Service.Spec.Manifests)
	if err != nil {
		return nil, err
	}
	serviceBinding.Workload.Deployments = deploymentWorkloadBuilder(sbReq.Service.Spec.Operator.Deployments)

	return serviceBinding, nil
}

func serviceCRDBuilder(csds []svcv1alpha1.CustomServiceDefinition) ([]string, error) {
	var customResource []string
	for _, csd := range csds {
		crdString, err := json.Marshal(csd.Spec.CRD)
		if err != nil {
			klog.Errorf("json Marshal service %s crd to string failed", csd.Name)
			return nil, err
		}
		customResource = append(customResource, string(crdString))
	}
	return customResource, nil
}

func servicePermissionBuilder(crs []rbacv1.ClusterRole, crbs []rbacv1.ClusterRoleBinding) []enginev1alpha1.Permission {
	crMap := getClusterRoleMap(crs)
	permissions := make([]enginev1alpha1.Permission, 0, len(crs))
	for _, binding := range crbs {
		cr, find := crMap[binding.RoleRef.Name]
		if !find {
			continue
		}
		permissions = append(permissions, enginev1alpha1.Permission{
			ServiceAccountName: binding.Subjects[0].Name,
			Rules:              cr.Rules,
		})
	}
	return permissions
}

func serviceCapabilityPluginBuilder(_ []svcv1alpha1.CustomServiceDefinition) enginev1alpha1.CapabilityPlugin {
	return enginev1alpha1.CapabilityPlugin{}
}

func deploymentWorkloadBuilder(deployments []appsv1.Deployment) []enginev1alpha1.ServiceDeploymentSpec {
	var deploymentSpecs []enginev1alpha1.ServiceDeploymentSpec
	for _, deployment := range deployments {
		deploymentSpec := enginev1alpha1.ServiceDeploymentSpec{
			Name: deployment.Name,
			Spec: deployment.Spec,
		}
		deploymentSpecs = append(deploymentSpecs, deploymentSpec)
	}
	return deploymentSpecs
}

func getClusterRoleMap(crs []rbacv1.ClusterRole) map[string]rbacv1.ClusterRole {
	crMap := make(map[string]rbacv1.ClusterRole, len(crs))
	for _, cr := range crs {
		crMap[cr.Name] = cr
	}
	return crMap
}
