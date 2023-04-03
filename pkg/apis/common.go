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

package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceType defines the package type, can be operator or helm
type ServiceType string

const (
	operatorServiceType ServiceType = "operator"
	helmServiceType     ServiceType = "helm"

	maxStringLen = 64

	defaultSource = "OpenSource"
)

var serviceTypeSet = map[ServiceType]struct{}{
	operatorServiceType: {},
	helmServiceType:     {},
}

// AbstractResource defines a generic format of a resource to facilitate unmarshalling of resource files
type AbstractResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              Generic `json:"spec,omitempty"`
	Status            Generic `json:"status,omitempty"`
}

// Generic defines a generic form of resource file
type Generic map[string]interface{}

// RawServiceType represents all types of CloudNativeService
type RawServiceType string

// Phase of the runtime service instance
type Phase string
