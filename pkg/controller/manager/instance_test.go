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

var testInstanceController *InstanceController

func TestInstanceController_GetInstances(t *testing.T) {
	convey.Convey("Test InstanceController GetInstances", t, func() {
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "_")
		testInstanceController.GetInstances()
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "")
		testInstanceController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testInstanceController.GetInstances()
	})
}

func TestInstanceController_GetInstanceDetail(t *testing.T) {
	convey.Convey("Test InstanceController GetInstanceDetail", t, func() {
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "_")
		testInstanceController.GetInstanceDetail()
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "")
		testInstanceController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testInstanceController.GetInstanceDetail()
	})
}

func TestInstanceController_DeleteInstance(t *testing.T) {
	convey.Convey("Test InstanceController DeleteInstance", t, func() {
		testInstanceController.Ctx.Input.SetParam(constants.InstancePathParam, "_")
		testInstanceController.DeleteInstance()
		testInstanceController.Ctx.Input.SetParam(constants.InstancePathParam, "xx")
		testInstanceController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testInstanceController.DeleteInstance()
	})
}

func TestInstanceController_CreateInstance(t *testing.T) {
	convey.Convey("Test InstanceController DeleteInstance", t, func() {
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "_")
		testInstanceController.CreateInstance()
		testInstanceController.Ctx.Input.SetParam(constants.ServiceBindingPathParam, "xx")
		testInstanceController.Ctx.Input.SetParam(constants.ClusterNameQueryParam, "_")
		testInstanceController.CreateInstance()
	})
}
