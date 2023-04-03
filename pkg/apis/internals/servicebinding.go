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
	"time"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

// ServiceBinding the service binding struct which using in the program internal
type ServiceBinding struct {
	ID               string
	Name             string
	Version          string
	Namespace        string
	ServiceName      string
	ServiceID        string
	ClusterID        string
	ClusterName      string
	Status           string
	Message          string
	ProcessTime      time.Time
	UpdateTime       time.Time
	CRD              []string
	Permissions      []enginev1alpha1.Permission
	Workload         enginev1alpha1.Workload
	CapabilityPlugin enginev1alpha1.CapabilityPlugin
}
