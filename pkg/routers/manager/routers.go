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
	"github.com/beego/beego/v2/server/web"

	"github.com/kappital/kappital/pkg/controller/manager"
	"github.com/kappital/kappital/pkg/routers"
)

// InitRouters init the routers for manager
func InitRouters() {
	registerServiceBindingAPI()
	registerInstanceAPI()

	routers.InitFilters()
}

func registerServiceBindingAPI() {
	// servicebinding
	web.Router("/api/v1alpha1/servicebinding", &manager.ServiceBindingController{},
		"post:CreateServiceBinding")
	web.Router("/api/v1alpha1/servicebinding/:service_binding", &manager.ServiceBindingController{},
		"delete:DeleteServiceBinding")
	web.Router("/api/v1alpha1/servicebinding", &manager.ServiceBindingController{},
		"get:GetServiceBindings")
	web.Router("/api/v1alpha1/servicebinding/:service_binding", &manager.ServiceBindingController{},
		"get:GetServiceBindingDetail")
}

func registerInstanceAPI() {
	web.Router("/api/v1alpha1/servicebinding/:service_binding/instance", &manager.InstanceController{},
		"post:CreateInstance")
	web.Router("/api/v1alpha1/servicebinding/:service_binding/instance/:instance", &manager.InstanceController{},
		"delete:DeleteInstance")
	web.Router("/api/v1alpha1/servicebinding/:service_binding/instance", &manager.InstanceController{},
		"get:GetInstances")
	web.Router("/api/v1alpha1/servicebinding/:service_binding/instance/:instance", &manager.InstanceController{},
		"get:GetInstanceDetail")
}
