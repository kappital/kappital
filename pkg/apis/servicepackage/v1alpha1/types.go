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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kappital/kappital/pkg/apis"
)

// CloudNativePackage the summary of the static service package
type CloudNativePackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec"`
	Status PackageStatus `json:"status,omitempty"`
}

// PackageSpec the static information of the service instance
type PackageSpec struct {
	Repository string          `json:"repository"`
	Version    string          `json:"version,omitempty"`
	Descriptor apis.Descriptor `json:"descriptor,omitempty"`
}

// PackageStatus the status of this package at remote
type PackageStatus struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// Replicate replicate the current struct
func (c CloudNativePackage) Replicate() CloudNativePackage {
	return c
}

// CloudNativePackageSlice of the CloudNativePackage
type CloudNativePackageSlice []CloudNativePackage

// Len of the slice
func (c CloudNativePackageSlice) Len() int {
	return len(c)
}

// Swap the i and j index
func (c CloudNativePackageSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less compare the i and j index, sort the slice from new to old
func (c CloudNativePackageSlice) Less(i, j int) bool {
	return c[i].CreationTimestamp.After(c[j].CreationTimestamp.Time)
}
