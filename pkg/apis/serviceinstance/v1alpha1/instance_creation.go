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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
)

// ServiceInstanceCreation of the manager which the request body
type ServiceInstanceCreation struct {
	InstanceName            string                         `json:"instanceName,omitempty"`            // CloudNativeServiceInstance metadata.name
	ClusterID               string                         `json:"clusterID,omitempty"`               // default if clusterID is null
	Service                 svcv1alpha1.CloudNativeService `json:"service,omitempty"`                 // service CloudNativeService
	InstanceCustomResources []InstanceCustomResource       `json:"instanceCustomResources,omitempty"` // cr list
}

// InstanceCustomResource user's custom resource of the instance
type InstanceCustomResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              json.RawMessage `json:"spec,omitempty"`
}

// Validate validate the struct variables
func (s *ServiceInstanceCreation) Validate() error {
	if len(s.ClusterID) == 0 {
		s.ClusterID = apis.DefaultCluster
	}

	for i := range s.InstanceCustomResources {
		if len(s.InstanceCustomResources[i].Namespace) == 0 {
			s.InstanceCustomResources[i].Namespace = apis.DefaultNamespace
		}
	}
	return nil
}
