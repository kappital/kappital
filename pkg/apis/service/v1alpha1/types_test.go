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
	"github.com/kappital/kappital/pkg/apis"
	"testing"
)

func TestCloudNativeService_Validation(t *testing.T) {
	type fields struct {
		Spec CloudNativeServiceSpec
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Test CloudNativeService Validation (no deployment resource)",
			fields: fields{Spec: CloudNativeServiceSpec{}},
		},
		{
			name:   "Test CloudNativeService Validation (not valid descriptor)",
			fields: fields{Spec: CloudNativeServiceSpec{Operator: &OperatorSpec{}}},
		},
		{
			name: "Test CloudNativeService Validation",
			fields: fields{
				Spec: CloudNativeServiceSpec{
					Operator:    &OperatorSpec{},
					Description: apis.Descriptor{Name: "name", Version: "version", Type: "operator"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CloudNativeService{Spec: tt.fields.Spec}
			if got := c.Validation(); got != tt.want {
				t.Errorf("Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}
