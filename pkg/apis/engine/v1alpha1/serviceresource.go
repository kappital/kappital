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

package v1alpha1

import (
	"encoding/base64"
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
	rbac "k8s.io/api/rbac/v1"
)

// ServiceResource is the resource for the cloud native service instance
type ServiceResource struct {
	CustomResourceDefinitions []string         `json:"customResourceDefinitions"`
	Permissions               []Permission     `json:"permissions"`
	CapabilityPlugin          CapabilityPlugin `json:"capabilityPlugin,omitempty"`
	Workload                  Workload         `json:"workload"`
}

// Permission is the permission for user's operator(s), only support for cluster role and cluster role binding.
// In other words, engine consider the same Service Instance will only have exactly one operator in one cluster.
type Permission struct {
	ServiceAccountName string            `json:"serviceAccountName"`
	Rules              []rbac.PolicyRule `json:"rules"`
}

// CapabilityPlugin is the extra functions for operations
// TODO: Design the CapabilityPlugin structure and apply each functions
type CapabilityPlugin struct {
}

// Workload is the operator workloads, now only support for deployment, daemon set, and stateful set
type Workload struct {
	Deployments  []ServiceDeploymentSpec  `json:"deployments,omitempty"`
	DaemonSets   []ServiceDaemonSetSpec   `json:"daemonSets,omitempty"`
	StatefulSets []ServiceStatefulSetSpec `json:"statefulSets,omitempty"`
}

// ServiceDeploymentSpec the message of Deployment object information.
type ServiceDeploymentSpec struct {
	Name string                `json:"name"`
	Spec appsv1.DeploymentSpec `json:"spec"`
}

// ServiceDaemonSetSpec the message of DaemonSet object information.
type ServiceDaemonSetSpec struct {
	Name string               `json:"name"`
	Spec appsv1.DaemonSetSpec `json:"spec"`
}

// ServiceStatefulSetSpec the message of StatefulSet object information.
type ServiceStatefulSetSpec struct {
	Name string                 `json:"name"`
	Spec appsv1.StatefulSetSpec `json:"spec"`
}

// TranslateResourcesToBase64 translate the ServiceResource object to the base64 code.
func TranslateResourcesToBase64(resource ServiceResource) (string, error) {
	bytes, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// TranslateBase64CodeToResource translate the base64 code to the ServiceResource object.
func TranslateBase64CodeToResource(binary string) (ServiceResource, error) {
	decoded, err := base64.StdEncoding.DecodeString(binary)
	if err != nil {
		return ServiceResource{}, err
	}
	resource := ServiceResource{}
	if err = json.Unmarshal(decoded, &resource); err != nil {
		return ServiceResource{}, err
	}
	return resource, nil
}
