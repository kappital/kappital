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
	"reflect"
	"testing"
)

func TestCloudNativePackageSlice_Len(t *testing.T) {
	tests := []struct {
		name string
		c    CloudNativePackageSlice
		want int
	}{
		{name: "Test CloudNativePackageSlice Len", c: CloudNativePackageSlice{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudNativePackageSlice_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    CloudNativePackageSlice
		args args
		want bool
	}{
		{
			name: "Test CloudNativePackageSlice Less",
			c:    CloudNativePackageSlice([]CloudNativePackage{{}, {}}),
			args: args{0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloudNativePackageSlice_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    CloudNativePackageSlice
		args args
	}{
		{
			name: "Test CloudNativePackageSlice Swap",
			c:    CloudNativePackageSlice([]CloudNativePackage{{}, {}}),
			args: args{0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Swap(tt.args.i, tt.args.j)
		})
	}
}

func TestCloudNativePackage_Replicate(t *testing.T) {
	tests := []struct {
		name string
		want CloudNativePackage
	}{
		{name: "Test CloudNativePackage Replicate", want: CloudNativePackage{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CloudNativePackage{}
			if got := c.Replicate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
