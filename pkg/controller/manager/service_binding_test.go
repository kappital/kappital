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

package manager

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/kappital/kappital/pkg/constants"
)

var testServiceBindingController *ServiceBindingController

func TestServiceBindingController_DeleteServiceBinding(t *testing.T) {
	convey.Convey("Test ServiceBindingController DeleteServiceBinding", t, func() {
		testServiceBindingController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "_")
		testServiceBindingController.DeleteServiceBinding()
		testServiceBindingController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "xx")
		testServiceBindingController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testServiceBindingController.DeleteServiceBinding()
	})
}

func TestServiceBindingController_GetServiceBindings(t *testing.T) {
	convey.Convey("Test ServiceBindingController GetServiceBindings", t, func() {
		testServiceBindingController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testServiceBindingController.GetServiceBindings()
	})
}

func TestServiceBindingController_GetServiceBindingDetail(t *testing.T) {
	convey.Convey("Test ServiceBindingController GetServiceBindingDetail", t, func() {
		testServiceBindingController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "_")
		testServiceBindingController.GetServiceBindingDetail()
		testServiceBindingController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "xx")
		testServiceBindingController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testServiceBindingController.GetServiceBindingDetail()
	})
}
