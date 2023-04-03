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

import "testing"

var invalidTestStr = "1234567890123456789012345678901234567890123456789012345678901234567890"

func TestDescriptor_Validation(t *testing.T) {
	type fields struct {
		Name    string
		Version string
		Type    ServiceType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test Descriptor Validation (name is empty)",
			want: false,
		},
		{
			name:   "Test Descriptor Validation (name len is greater than the max)",
			fields: fields{Name: invalidTestStr},
			want:   false,
		},
		{
			name: "Test Descriptor Validation (version is invalid)",
			fields: fields{
				Name:    "name",
				Version: invalidTestStr,
			},
			want: false,
		},
		{
			name: "Test Descriptor Validation (service type is invalid)",
			fields: fields{
				Name: "name",
				Type: "operatorServiceType",
			},
			want: false,
		},
		{
			name: "Test Descriptor Validation",
			fields: fields{
				Name:    "name",
				Version: "version",
				Type:    operatorServiceType,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Descriptor{
				Name:    tt.fields.Name,
				Version: tt.fields.Version,
				Type:    tt.fields.Type,
			}
			if got := d.Validation(); got != tt.want {
				t.Errorf("Validation() = %v, want %v", got, tt.want)
			}
		})
	}
}
