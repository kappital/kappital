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
	"reflect"
	"testing"
)

var identifierTest = CNSIdentifier{
	Repo:    "repo",
	Name:    "name",
	Version: "version",
}

func TestCNSIdentifier_GetFilterMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "Test CNSIdentifier GetFilterMap",
			want: map[string]string{"repository": "repo", "name": "name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := identifierTest.GetFilterMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilterMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCNSIdentifier_GetTmpDirName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test CNSIdentifier GetTmpDirName",
			want: "repo-name-version",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := identifierTest.GetTmpDirName(); got != tt.want {
				t.Errorf("GetTmpDirName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCNSIdentifier_GetTmpYamlFileName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test CNSIdentifier GetYamlFileName",
			want: "repo-name-version.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := identifierTest.GetYamlFileName(); got != tt.want {
				t.Errorf("GetYamlFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
