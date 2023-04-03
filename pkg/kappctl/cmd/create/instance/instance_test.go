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

package instance

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/convert"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

func Test_getCRContent(t *testing.T) {
	convey.Convey("Test getCRContent", t, func() {
		p := gomonkey.ApplyFuncSeq(filepath.Abs, []gomonkey.OutputCell{
			{Values: gomonkey.Params{"", fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{"", nil}},
			{Values: gomonkey.Params{"", nil}},
			{Values: gomonkey.Params{"", nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
			{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]byte("xx:xx"), nil}},
			{Values: gomonkey.Params{[]byte("{}"), nil}},
		})

		res, err := getCRContent("1")
		convey.So(res, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)

		res, err = getCRContent("2")
		convey.So(res, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)

		res, err = getCRContent("3")
		convey.So(res, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)

		res, err = getCRContent("4")
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_getInstanceCustomResource(t *testing.T) {
	type args struct {
		spec svcv1alpha1.CustomServiceDefinitionSpec
	}
	tests := []struct {
		name string
		args args
		want instancev1alpha1.InstanceCustomResource
	}{
		{
			name: "Test getInstanceCustomResource (cannot get group)",
			args: args{spec: svcv1alpha1.CustomServiceDefinitionSpec{CRD: &apis.AbstractResource{}}},
		},
		{
			name: "Test getInstanceCustomResource (cannot get names' kind as string)",
			args: args{spec: svcv1alpha1.CustomServiceDefinitionSpec{
				CRD: &apis.AbstractResource{Spec: map[string]interface{}{
					"group": "",
					"names": map[string]interface{}{"kind": 12},
				}},
			}},
		},
		{
			name: "Test getInstanceCustomResource",
			args: args{
				spec: svcv1alpha1.CustomServiceDefinitionSpec{
					CRD: &apis.AbstractResource{Spec: map[string]interface{}{
						"group": "",
						"names": map[string]interface{}{"kind": "12"},
					}},
					CRVersions: []svcv1alpha1.CRVersion{{}},
				}},
			want: instancev1alpha1.InstanceCustomResource{
				TypeMeta:   metav1.TypeMeta{Kind: "12", APIVersion: "/"},
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInstanceCustomResource(tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInstanceCustomResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_operation_NewCommand(t *testing.T) {
	o := &operation{}
	if got := o.NewCommand(); got == nil {
		t.Errorf("create instance operation NewCommand() = %v, do not want nil", got)
	}
}

func Test_operation_PreRunE(t *testing.T) {
	convey.Convey("Test operation PreRun", t, func() {
		o := operation{}
		convey.Convey("case 1: cannot get kappital config", func() {
			p := gomonkey.ApplyFunc(kappctl.GetConfig, func() (*kappctl.Config, error) {
				return nil, fmt.Errorf("mock error")
			})
			defer p.Reset()
			err := o.PreRunE([]string{})
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: other case", func() {
			p := gomonkey.ApplyFunc(kappctl.GetConfig, func() (*kappctl.Config, error) {
				return &kappctl.Config{}, nil
			})
			defer p.Reset()
			err := o.PreRunE([]string{"x"})
			convey.So(err, convey.ShouldBeNil)
			err = o.PreRunE([]string{})
			convey.So(err, convey.ShouldNotBeNil)

			o.dirPath = "absxxx"
			p.ApplyFunc(filepath.Abs, func(s string) (string, error) {
				if s != "abs" {
					return "", fmt.Errorf("mock error")
				}
				return s, nil
			})
			err = o.PreRunE([]string{})
			convey.So(err, convey.ShouldNotBeNil)

			o.dirPath = "abs"
			p.ApplyMethodSeq(reflect.TypeOf(convert.GetLoader()), "TransferToCloudNativeService", []gomonkey.OutputCell{
				{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
				{Values: gomonkey.Params{nil, nil}},
				{Values: gomonkey.Params{&svcv1alpha1.CloudNativeService{}, nil}},
			})
			err = o.PreRunE([]string{})
			convey.So(err, convey.ShouldNotBeNil)
			err = o.PreRunE([]string{})
			convey.So(err, convey.ShouldNotBeNil)
			err = o.PreRunE([]string{})
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("Test operation RunE", t, func() {
		o := operation{config: &kappctl.Config{}}
		err := o.RunE()
		convey.So(err, convey.ShouldNotBeNil)

		o.cns = &svcv1alpha1.CloudNativeService{
			Spec: svcv1alpha1.CloudNativeServiceSpec{
				Manifests: []svcv1alpha1.CustomServiceDefinition{{
					Spec: svcv1alpha1.CustomServiceDefinitionSpec{
						CRD: &apis.AbstractResource{Spec: map[string]interface{}{
							"group": "",
							"names": map[string]interface{}{"kind": "12"},
						}},
						CRVersions: []svcv1alpha1.CRVersion{{}},
					},
				}},
			},
		}
		p := gomonkey.ApplyFuncSeq(gateway.CommonUtilRequest, []gomonkey.OutputCell{
			{Values: gomonkey.Params{0, nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{0, nil, nil}},
			{Values: gomonkey.Params{http.StatusOK, nil, nil}},
			{Values: gomonkey.Params{http.StatusOK, nil, nil}},
		})
		defer p.Reset()
		err = o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldBeNil)

		o.resourcePath = "xx"
		p.ApplyFuncSeq(getCRContent, []gomonkey.OutputCell{
			{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil, nil}},
		})
		err = o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldBeNil)
	})
}
