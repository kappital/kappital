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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Phase of the runtime service instance
type Phase string

const (
	// PendingPhase the service instance during installing
	PendingPhase Phase = "Pending"
	// SucceededPhase the service instance has already deployed succeeded
	SucceededPhase Phase = "Succeeded"
	// FailedPhase all service instance does not deployed succeeded
	FailedPhase Phase = "Failed"
	// UnknownPhase partial service instance deployed failed or some service instance cannot find in cluster
	UnknownPhase Phase = "Unknown"
)

// CloudNativeServiceInstance the summary of the runtime service package
type CloudNativeServiceInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CloudNativeServiceInstanceSpec   `json:"spec"`
	Status            CloudNativeServiceInstanceStatus `json:"status,omitempty"`
}

// CloudNativeServiceInstanceSpec the static information of the service instance
type CloudNativeServiceInstanceSpec struct {
	Name                string           `json:"name"`
	Version             string           `json:"version"`
	ServiceName         string           `json:"serviceName,omitempty"`
	ServiceID           string           `json:"serviceID,omitempty"`
	ClusterName         string           `json:"clusterName,omitempty"`
	Description         string           `json:"description,omitempty"`
	ClusterID           string           `json:"clusterId,omitempty"`
	CustomResources     []Resource       `json:"customResources,omitempty"`
	CapabilityResources []Resource       `json:"capabilityResources,omitempty"`
	ServiceReference    ServiceReference `json:"serviceReference"`
	DependentResources  []Resource       `json:"dependentResources,omitempty"`
}

// Resource the user's instance
type Resource struct {
	metav1.TypeMeta `json:",inline"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace,omitempty"`
	UID             string `json:"uid,omitempty"`
	Status          string `json:"status,omitempty"`
	RawMessage      string `json:"rawMessage,omitempty"`
}

// ServiceReference the service package information
type ServiceReference struct {
	metav1.TypeMeta `json:",inline"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace,omitempty"`
	UID             string `json:"uid"`
	Status          string `json:"status,omitempty"`
}

// CloudNativeServiceInstanceStatus the status of the service and its instance
type CloudNativeServiceInstanceStatus struct {
	Phase            Phase              `json:"phase"`
	Message          string             `json:"message"`
	DependencyStatus []DependencyStatus `json:"dependencyStatus,omitempty"`
}

// DependencyStatus the dependency service status
// TODO: implement dependency in feature
type DependencyStatus struct {
	Phase          string `json:"phase,omitempty"`
	Message        string `json:"message,omitempty"`
	ServiceID      string `json:"serviceID"`
	ServiceName    string `json:"serviceName"`
	ServiceVersion string `json:"serviceVersion"`
}

// ClusterInformation which manager managed
type ClusterInformation struct {
	Name            string    `json:"name"`
	KubeConfig      string    `json:"kubeConfig,omitempty"`
	Version         string    `json:"version,omitempty"`
	CreateTimestamp time.Time `json:"createTimestamp,omitempty"`
	UpdateTimestamp time.Time `json:"updateTimestamp,omitempty"`
}
