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
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"

	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

func Test_operation_NewCommand(t *testing.T) {
	o := &operation{}
	if got := o.NewCommand(); got == nil {
		t.Errorf("delete instance operation NewCommand() = %v, do not want nil", got)
	}
}

func Test_operation_PreRunE(t *testing.T) {
	convey.Convey("test delete instance operation PreRunE", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: cannot get the kappctl config for delete instance", func() {
			err := o.PreRunE([]string{""})
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: can get the kappctl config and valid args slice length for delete instance", func() {
			p := gomonkey.ApplyFunc(kappctl.GetConfig, func() (*kappctl.Config, error) {
				return &kappctl.Config{}, nil
			})
			defer p.Reset()
			err := o.PreRunE([]string{"instance"})
			convey.So(err, convey.ShouldBeNil)
			err = o.PreRunE([]string{"1ee23456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"})
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("test delete instance operation RunE", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: gateway.CommonUtilRequest get error", func() {
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: gateway.CommonUtilRequest get failed http code", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return -1, nil, nil
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: gateway.CommonUtilRequest get valid data for delete instance", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return http.StatusOK, []byte{}, nil
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
