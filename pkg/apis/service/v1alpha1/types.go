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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
)

// CloudNativeService constructs service information
// in a common way that eliminates differences between different implementations.
type CloudNativeService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudNativeServiceSpec `json:"spec,omitempty"`
}

// CloudNativeServiceSpec defines the specification for a CloudNativeService.
type CloudNativeServiceSpec struct {
	Operator    *OperatorSpec             `json:"operator,omitempty"`
	RawResource *RawResource              `json:"rawResource,omitempty"`
	Description apis.Descriptor           `json:"description"`
	Manifests   []CustomServiceDefinition `json:"manifests,omitempty"`
	Version     string                    `json:"version"`
}

// CustomServiceDefinition defines additional information of a CRD
type CustomServiceDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CustomServiceDefinitionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// CRVersion gives default values of a specific version defined by CRD
type CRVersion struct {
	Name          string `json:"name"`
	CRName        string `json:"CRName"`
	DefaultValues string `json:"defaultValues,omitempty"`
}

// CustomServiceDefinitionSpec is the spec of CustomServiceDefinition
type CustomServiceDefinitionSpec struct {
	CRD                    *apis.AbstractResource `json:"CRD,omitempty"`
	CRDName                string                 `json:"CRDName,omitempty"`
	CRVersions             []CRVersion            `json:"CRVersions,omitempty"`
	Description            string                 `json:"description,omitempty"`
	Role                   ResourceRole           `json:"role,omitempty"`
	CapabilityRequirements []GVKAndName           `json:"capabilityRequirements,omitempty"`
}

// GVKAndName uniquely identifies a GVK+Name
type GVKAndName struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// ResourceRole is the role of a CRD, can be ServiceEntity or Attribute
type ResourceRole string

// OperatorSpec represents ServicePack lifecycle.
type OperatorSpec struct {
	Deployments         []appsv1.Deployment       `json:"deployments"`
	ServiceAccounts     []corev1.ServiceAccount   `json:"serviceAccounts,omitempty"`
	Roles               []rbac.Role               `json:"roles,omitempty"`
	RoleBindings        []rbac.RoleBinding        `json:"roleBindings,omitempty"`
	ClusterRoles        []rbac.ClusterRole        `json:"clusterRoles,omitempty"`
	ClusterRoleBindings []rbac.ClusterRoleBinding `json:"clusterRoleBindings,omitempty"`
}

// RawResource define CloudNativeService include 3rd service, one of
type RawResource struct {
	Type apis.RawServiceType   `json:"type,omitempty"`
	Spec apis.AbstractResource `json:"spec,omitempty"`
}

// CloudNativeServiceResponse is the http response of PushService
type CloudNativeServiceResponse struct {
	Repo    string `json:"repo"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Repository is the http request of CreateRepo
type Repository struct {
	Name   string `json:"project_name"`
	Public bool   `json:"public"`
}

// RepoResp is the http response of GetRepos
type RepoResp struct {
	Name            string
	Public          bool
	ServiceCount    int64
	CreateTimestamp time.Time
}

// Validation does the CloudNativeService is valid or not
func (c *CloudNativeService) Validation() bool {
	return c.Spec.Validation()
}

// Validation does the CloudNativeServiceSpec is valid or not, and generate some attribute
func (s *CloudNativeServiceSpec) Validation() bool {
	if s.Operator == nil && s.RawResource == nil {
		klog.Errorf("the package cannot be empty for the deployment information")
		return false
	}
	if !s.Description.Validation() {
		klog.Errorf("the description information is invalid")
		return false
	}
	s.Version = s.Description.Version
	return true
}
