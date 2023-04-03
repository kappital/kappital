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
	"testing"

	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
)

func TestServiceInstanceCreation_Validate(t *testing.T) {
	type fields struct {
		InstanceName            string
		ClusterID               string
		Service                 svcv1alpha1.CloudNativeService
		InstanceCustomResources []InstanceCustomResource
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ServiceInstanceCreation Validate",
			fields: fields{
				ClusterID:               "",
				InstanceCustomResources: []InstanceCustomResource{{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceInstanceCreation{
				InstanceName:            tt.fields.InstanceName,
				ClusterID:               tt.fields.ClusterID,
				Service:                 tt.fields.Service,
				InstanceCustomResources: tt.fields.InstanceCustomResources,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
