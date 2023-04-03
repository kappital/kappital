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

package service

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"

	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/convert"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

func Test_operation_NewCommand(t *testing.T) {
	o := &operation{}
	if got := o.NewCommand(); got == nil {
		t.Errorf("create service operation NewCommand() = %v, do not want nil", got)
	}
}

func Test_operation_PreRunE(t *testing.T) {
	convey.Convey("Test operation PreRunE", t, func() {
		p := gomonkey.ApplyMethodSeq(reflect.TypeOf(convert.GetLoader()), "TransferToCloudNativeService", []gomonkey.OutputCell{
			{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil, nil}},
			{Values: gomonkey.Params{&svcv1alpha1.CloudNativeService{}, nil}},
		})
		defer p.Reset()

		o := operation{}
		err := o.PreRunE([]string{""})
		convey.So(err, convey.ShouldNotBeNil)
		err = o.PreRunE([]string{""})
		convey.So(err, convey.ShouldNotBeNil)
		p.ApplyFunc(kappctl.GetConfig, func() (*kappctl.Config, error) {
			return &kappctl.Config{}, nil
		})
		err = o.PreRunE([]string{""})
		convey.So(err, convey.ShouldBeNil)
	})
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("Test operation RunE", t, func() {
		o := &operation{config: &kappctl.Config{}, cns: &svcv1alpha1.CloudNativeService{}}
		p := gomonkey.ApplyFuncSeq(gateway.CommonUtilRequest, []gomonkey.OutputCell{
			{Values: gomonkey.Params{-1, nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{-1, nil, nil}},
			{Values: gomonkey.Params{http.StatusOK, nil, nil}},
			{Values: gomonkey.Params{http.StatusOK, []byte("{\"ID\":\"x\",\"name\":\"xxx\"}"), nil}},
		})
		defer p.Reset()

		err := o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldNotBeNil)
		err = o.RunE()
		convey.So(err, convey.ShouldBeNil)
	})
}
