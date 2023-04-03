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

package version

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	crdV1      = apiextensionsv1.CustomResourceDefinition{TypeMeta: metav1.TypeMeta{APIVersion: apiExtensionsV1}, ObjectMeta: metav1.ObjectMeta{Name: "crdv1"}}
	crdV1Beta1 = apiextensionsv1beta1.CustomResourceDefinition{TypeMeta: metav1.TypeMeta{APIVersion: apiExtensionsV1Beta1}, ObjectMeta: metav1.ObjectMeta{Name: "crdv1beta1"}}

	crdV1String      = ""
	crdV1Beta1String = ""
)

func TestGetClusterVersion(t *testing.T) {
	tests := []struct {
		name            string
		changeToVersion string
		want            string
	}{
		{name: "Test GetClusterVersion(1.15)", changeToVersion: clusterVersion15, want: clusterVersion15},
		{name: "Test GetClusterVersion(1.16-1.22)", changeToVersion: clusterVersion16To22, want: clusterVersion16To22},
		{name: "Test GetClusterVersion(1.22+)", changeToVersion: clusterVersion22Plus, want: clusterVersion22Plus},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.changeToVersion) != 0 {
				clusterVersion = tt.changeToVersion
			}
			if got := GetClusterVersion(); got != tt.want {
				t.Errorf("GetClusterVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrdV1(t *testing.T) {
	type args struct {
		crdString string
	}
	tests := []struct {
		name    string
		args    args
		want    apiextensionsv1.CustomResourceDefinition
		wantErr bool
	}{
		{name: "Test GetCrdV1", args: args{crdString: crdV1String}, want: crdV1, wantErr: false},
		{name: "Test GetCrdV1 (err)", args: args{crdString: crdV1Beta1String}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrdV1(tt.args.crdString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrdV1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("GetCrdV1() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrdV1AndBeta1Slice(t *testing.T) {
	type args struct {
		crdStrings []string
	}
	tests := []struct {
		name  string
		args  args
		want  []apiextensionsv1.CustomResourceDefinition
		want1 []apiextensionsv1beta1.CustomResourceDefinition
	}{
		{
			name:  "Test GetCrdV1AndBeta1Slice",
			args:  args{crdStrings: []string{crdV1String, crdV1Beta1String}},
			want:  []apiextensionsv1.CustomResourceDefinition{crdV1},
			want1: []apiextensionsv1beta1.CustomResourceDefinition{crdV1Beta1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterVersion = clusterVersion16To22
			got, got1 := GetCrdV1AndBeta1Slice(tt.args.crdStrings)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrdV1AndBeta1Slice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetCrdV1AndBeta1Slice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetCrdV1Beta1(t *testing.T) {
	type args struct {
		crdString string
	}
	tests := []struct {
		name    string
		args    args
		want    apiextensionsv1beta1.CustomResourceDefinition
		wantErr bool
	}{
		{name: "Test GetCrdV1Beta1", args: args{crdString: crdV1Beta1String}, want: crdV1Beta1, wantErr: false},
		{name: "Test GetCrdV1Beta1 (err)", args: args{crdString: crdV1String}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrdV1Beta1(tt.args.crdString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrdV1Beta1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("GetCrdV1Beta1() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetClusterVersion(t *testing.T) {
	tests := []struct {
		name            string
		changeToVersion string
	}{
		{name: "Test SetClusterVersion (apiextensionsv1.15)", changeToVersion: "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetClusterVersion(tt.changeToVersion)
			if !reflect.DeepEqual(clusterVersion, tt.changeToVersion) {
				t.Errorf("SetClusterVersion() = %v, want %v", clusterVersion, tt.changeToVersion)
			}
		})
	}
}

func TestSupportCRDUseV1(t *testing.T) {
	tests := []struct {
		name            string
		changeToVersion string
		want            bool
	}{
		{name: "Test SupportCRDUseV1(1.15)", changeToVersion: clusterVersion15, want: false},
		{name: "Test SupportCRDUseV1(1.16-1.22)", changeToVersion: clusterVersion16To22, want: true},
		{name: "Test SupportCRDUseV1(1.22+)", changeToVersion: clusterVersion22Plus, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.changeToVersion) != 0 {
				clusterVersion = tt.changeToVersion
			}
			if got := SupportCRDUseV1(); got != tt.want {
				t.Errorf("SupportCRDUseV1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupportCRDUseV1Beta1(t *testing.T) {
	tests := []struct {
		name            string
		changeToVersion string
		want            bool
	}{
		{name: "Test SupportCRDUseV1Beta1(1.15)", changeToVersion: clusterVersion15, want: true},
		{name: "Test SupportCRDUseV1Beta1(1.16-1.22)", changeToVersion: clusterVersion16To22, want: true},
		{name: "Test SupportCRDUseV1Beta1(1.22+)", changeToVersion: clusterVersion22Plus, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.changeToVersion) != 0 {
				clusterVersion = tt.changeToVersion
			}
			if got := SupportCRDUseV1Beta1(); got != tt.want {
				t.Errorf("SupportCRDUseV1Beta1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSecondLevelVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    int
		wantErr bool
	}{
		{
			name:    "Test getSecondLevelVersion(15)",
			version: "apiextensionsv1.15.0",
			want:    15,
			wantErr: false,
		},
		{
			name:    "Test getSecondLevelVersion(17)",
			version: "apiextensionsv1.17.0",
			want:    17,
			wantErr: false,
		},
		{
			name:    "Test getSecondLevelVersion(23)",
			version: "apiextensionsv1.23.0",
			want:    23,
			wantErr: false,
		},
		{
			name:    "Test getSecondLevelVersion(illegal version levels length less than 3)",
			version: "apiextensionsv1.0",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "Test getSecondLevelVersion(illegal second level string)",
			version: "v.v.v.v",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecondLevelVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSecondLevelVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getSecondLevelVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func init() {
	crdV1Bytes, err := json.Marshal(crdV1)
	if err != nil {
		fmt.Println("cannot marshal crdV1, ", err)
		os.Exit(1)
	}
	crdV1String = string(crdV1Bytes)

	crdV1Beta1Bytes, err := json.Marshal(crdV1Beta1)
	if err != nil {
		fmt.Println("cannot marshal crdV1Beta1, ", err)
		os.Exit(1)
	}
	crdV1Beta1String = string(crdV1Beta1Bytes)
}
