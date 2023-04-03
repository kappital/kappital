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

package servicebinding

import (
	"reflect"
	"time"

	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/dao/servicebinding"
)

func getTypedObj(param interface{}) *internals.ServiceBinding {
	binding, ok := param.(*internals.ServiceBinding)
	if !ok {
		klog.Errorf("invalid object type for install servicebinding handler, expected:models.servicebinding, "+
			"actual :%s", reflect.TypeOf(param).Name())
		return nil
	}
	return binding
}

func updateSuccessStatus(binding *internals.ServiceBinding) error {
	binding.Status = enginev1alpha1.SucceededPhase
	binding.ProcessTime = time.Time{}
	binding.Message = ""
	dbStore := servicebinding.ServiceBinding{}
	return dbStore.Update(*binding)
}

func updateProcessTimeout(binding *internals.ServiceBinding, timeout time.Duration) error {
	if !binding.ProcessTime.IsZero() {
		return nil
	}
	binding.ProcessTime = time.Now().Add(timeout)
	dbStore := servicebinding.ServiceBinding{}
	err := dbStore.Update(*binding)
	if err != nil {
		return err
	}
	return nil
}
