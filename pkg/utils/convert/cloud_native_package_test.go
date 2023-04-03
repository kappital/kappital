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

package convert

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	"github.com/kappital/kappital/pkg/utils/file"
)

func Test_cloudNativePackage_TransferToCloudNativeService(t *testing.T) {
	k := cloudNativePackage{}
	_, err := k.TransferToCloudNativeService("../../kappctl/cmd/initiate/kappital-demo")
	if err != nil {
		t.Errorf("cannot get nil form TransferToCloudNativeService, err: %v", err)
	}
}

func Test_cloudNativePackage_fileResolution(t *testing.T) {
	convey.Convey("Test cloudNativePackage fileResolution", t, func() {
		k := cloudNativePackage{}
		err := k.fileResolution(file.FakeFileInfo{}, 0)
		convey.So(err, convey.ShouldBeNil)
		k.fileCounter = 0
		t.Setenv("UNIT-TEST-FAKE-FILE", "x")
		err = k.fileResolution(file.FakeFileInfo{}, 0)
		convey.So(err, convey.ShouldNotBeNil)
		spRequiredDirSet.Insert("fake")
		k.fileCounter = maxFileCount
		err = k.fileResolution(file.FakeFileInfo{}, 10)
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func Test_cloudNativePackage_processGVK(t *testing.T) {
	type args struct {
		decoder *yaml.YAMLOrJSONDecoder
		gvk     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test cloudNativePackage processGVK (ClusterRole)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testCR.Name]), decodeBufferSize),
				gvk:     clusterRoleGVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (ClusterRoleBinding)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testCRB.Name]), decodeBufferSize),
				gvk:     clusterRoleBindingGVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (Deployment)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testDeploy.Name]), decodeBufferSize),
				gvk:     deploymentGVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (ServiceAccount)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testSA.Name]), decodeBufferSize),
				gvk:     serviceAccountGVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (CustomResourceDefinition v1)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testCRDV1.Name]), decodeBufferSize),
				gvk:     crdV1GVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (CustomResourceDefinition v1beta1)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testCRDV1Beta1.Name]), decodeBufferSize),
				gvk:     crdV1Beta1GVK,
			},
		},
		{
			name: "Test cloudNativePackage processGVK (CustomServiceDefinition)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testCSD.Name]), decodeBufferSize),
				gvk:     csdGVK,
			},
		},
		{
			name:    "Test cloudNativePackage processGVK (not include situation)",
			args:    args{gvk: "clusterRoleGVK"},
			wantErr: true,
		},
		{
			name: "Test cloudNativePackage processGVK (has error)",
			args: args{
				decoder: yaml.NewYAMLOrJSONDecoder(bytes.NewReader(contentMap[testSA.Name]), decodeBufferSize),
				gvk:     csdGVK,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &cloudNativePackage{}
			if tt.wantErr {
				p := gomonkey.ApplyMethodSeq(reflect.TypeOf(tt.args.decoder), "Decode", []gomonkey.OutputCell{
					{Values: gomonkey.Params{fmt.Errorf("mock error")}},
				})
				defer p.Reset()
			}
			if err := k.processGVK(tt.args.decoder, tt.args.gvk); (err != nil) != tt.wantErr {
				t.Errorf("processGVK() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cloudNativePackage_processGVKDirectory(t *testing.T) {
	convey.Convey("Test cloudNativePackage processGVKDirectory", t, func() {
		k := cloudNativePackage{}
		err := k.processGVKDirectory("", 10, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldNotBeNil)
		p := gomonkey.ApplyFuncSeq(ioutil.ReadDir, []gomonkey.OutputCell{
			{Values: gomonkey.Params{[]fs.FileInfo{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]fs.FileInfo{file.FakeFileInfo{}}, nil}},
			{Values: gomonkey.Params{[]fs.FileInfo{}, nil}},
			{Values: gomonkey.Params{[]fs.FileInfo{file.FakeFileInfo{}}, nil}},
			{Values: gomonkey.Params{[]fs.FileInfo{file.FakeFileInfo{}}, nil}},
			{Values: gomonkey.Params{[]fs.FileInfo{}, fmt.Errorf("mock error")}},
		})
		defer p.Reset()
		err = k.processGVKDirectory("", 0, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldNotBeNil)
		err = k.processGVKDirectory("", 2, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldNotBeNil)
		err = k.processGVKDirectory("", 1, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldBeNil)
		t.Setenv("UNIT-TEST-FAKE-FILE", "x")
		err = k.processGVKDirectory("", 3, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldBeNil)
		err = k.processGVKDirectory("", 1, dirAcceptGVKSet[operatorDirName])
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func Test_cloudNativePackage_processGVKFile(t *testing.T) {
	convey.Convey("Test cloudNativePackage processGVKFile", t, func() {
		convey.Convey("case 1: get the errors before for loop", func() {
			k := cloudNativePackage{}
			err := k.processGVKFile("fake path", 10, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processGVKFile("fake path", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			k.fileCounter = maxFileCount + 1
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			k.fileCounter = 0
			p := gomonkey.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
				{Values: gomonkey.Params{nil, fmt.Errorf("fack error")}},
				{Values: gomonkey.Params{make([]byte, singleFileSize+2), nil}},
			})
			defer p.Reset()
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
		})
		k := cloudNativePackage{}
		convey.Convey("case 2: illegal and not accept GVK", func() {
			p := gomonkey.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
				{Values: gomonkey.Params{[]byte{}, nil}},
				{Values: gomonkey.Params{[]byte{'x'}, nil}},
				{Values: gomonkey.Params{contentMap[testEmptyGVK.Name], nil}},
				{Values: gomonkey.Params{contentMap[testCSD.Name], nil}},
				{Values: gomonkey.Params{contentMap[testDeploy.Name], nil}},
			})
			defer p.Reset()
			err := k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldBeNil)
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processGVKFile("fake.json", 0, dirAcceptGVKSet[operatorDirName])
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_cloudNativePackage_processMetadata(t *testing.T) {
	convey.Convey("Test cloudNativePackage processMetadata", t, func() {
		convey.Convey("case 1: get errOverMaxFileCount", func() {
			k := cloudNativePackage{fileCounter: maxFileCount}
			err := k.processMetadata("fake path")
			convey.So(err, convey.ShouldEqual, errOverMaxFileCount)
		})
		k := cloudNativePackage{}
		convey.Convey("case 2: cannot read the file", func() {
			p := gomonkey.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
				{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
				{Values: gomonkey.Params{make([]byte, singleFileSize+2), nil}},
				{Values: gomonkey.Params{[]byte{}, nil}},
			})
			defer p.Reset()
			err := k.processMetadata("fake path")
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processMetadata("fake path")
			convey.So(err, convey.ShouldNotBeNil)
			err = k.processMetadata("fake path")
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: other cases", func() {
			p := gomonkey.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
				{Values: gomonkey.Params{[]byte(testMetadatas[0]), nil}},
				{Values: gomonkey.Params{[]byte(testMetadatas[1]), nil}},
			})
			defer p.Reset()
			err := k.processMetadata("fake path")
			convey.So(err, convey.ShouldBeNil)
			err = k.processMetadata("fake path")
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func Test_newCloudNativePackageLoader(t *testing.T) {
	tests := []struct {
		name string
		want Convert
	}{
		{
			name: "Test newCloudNativePackageLoader",
			want: &cloudNativePackage{
				operator: svcv1alpha1.OperatorSpec{
					Deployments:         make([]appsv1.Deployment, 0, initResourceSize),
					ServiceAccounts:     make([]corev1.ServiceAccount, 0, initResourceSize),
					ClusterRoles:        make([]rbacv1.ClusterRole, 0, initResourceSize),
					ClusterRoleBindings: make([]rbacv1.ClusterRoleBinding, 0, initResourceSize),
				},
				manifest: manifest{
					CRDs: make([]apis.AbstractResource, 0, initResourceSize),
					CSDs: make([]svcv1alpha1.CustomServiceDefinition, 0, initResourceSize),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCloudNativePackageLoader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCloudNativePackageLoader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cloudNativePackage_convert2CloudNativeService(t *testing.T) {
	type fields struct {
		manifest manifest
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			name: "Test cloudNativePackage convert2CloudNativeService",
			fields: fields{manifest: manifest{
				CRDs: []apis.AbstractResource{{ObjectMeta: metav1.ObjectMeta{Name: "test-xxx"}}},
				CSDs: []svcv1alpha1.CustomServiceDefinition{
					{Spec: svcv1alpha1.CustomServiceDefinitionSpec{CRDName: "test-xxx"}},
				},
			}},
		},
		{
			name: "Test cloudNativePackage convert2CloudNativeService (not found)",
			fields: fields{manifest: manifest{
				CRDs: []apis.AbstractResource{{ObjectMeta: metav1.ObjectMeta{Name: "test-xxx-x"}}},
				CSDs: []svcv1alpha1.CustomServiceDefinition{
					{Spec: svcv1alpha1.CustomServiceDefinitionSpec{CRDName: "test-xxx-y"}},
				},
			}},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &cloudNativePackage{manifest: tt.fields.manifest}
			got, err := k.convert2CloudNativeService()
			if (err != nil) != tt.wantErr {
				t.Errorf("convert2CloudNativeService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantNil {
				t.Errorf("convert2CloudNativeService() got = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
