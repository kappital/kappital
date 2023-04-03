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

package internals

import (
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceInstance the service instance struct which using in the program internal
type ServiceInstance struct {
	ID                 string
	Name               string
	Namespace          string
	ClusterID          string
	ClusterName        string
	RawResource        string
	CreateTime         time.Time
	Status             string
	Kind               string
	APIVersion         string
	ProcessTime        time.Time
	Resource           string
	Message            string
	ServiceBindingName string
	ServiceBindingID   string
	ServiceName        string
	ServiceID          string
	InstanceType       InstanceType
	UpdateTime         time.Time
	InstallState       InstallState
	RuntimeState       RuntimeState
}

// InstallState of the cloud native service instance
type InstallState struct {
	Phase    string
	SubPhase []Condition
}

// RuntimeState of the cloud native service instance
type RuntimeState struct {
	Phase    string
	RawState json.RawMessage
}

// Condition of the cloud native service instance
type Condition struct {
	Type               string
	Status             string
	Message            string
	LastTransitionTime metav1.Time
	RetryCount         int32
}

// InstanceFilter of the cloud native service instance
type InstanceFilter struct {
	ServiceBindingName string
	Instance           string
	Namespace          string
	ClusterID          string
}

// InstanceType the instance type of the service instance
type InstanceType string
