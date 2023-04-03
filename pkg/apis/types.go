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

import "k8s.io/apimachinery/pkg/runtime/schema"

const (
	// CloudNativeServiceInstanceKind name
	CloudNativeServiceInstanceKind = "CloudNativeServiceInstance"
	// CloudNativeServiceKind name
	CloudNativeServiceKind = "CloudNativeService"
	// CloudNativePackageKind name
	CloudNativePackageKind = "CloudNativePackage"
	// CloudNativeAPIVersionV1Alpha1 name
	CloudNativeAPIVersionV1Alpha1 = "core.kappital.io/v1alpha1"

	// KappitalSystemNamespace default manager and engine deployed namespace
	KappitalSystemNamespace = "kappital-system"
	// DefaultNamespace of kubernetes
	DefaultNamespace = "default"
	// DefaultCluster of manager to deploy the service instance
	DefaultCluster = "default"
	// DefaultVersion of the service package
	DefaultVersion = "latest"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "core.kappital.io", Version: "v1alpha1"}
)

const (
	// InstanceProcessor asynchronizing name
	InstanceProcessor string = "instance"
	// OperatorProcessor asynchronizing name
	OperatorProcessor string = "operator"
)
