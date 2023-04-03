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

package handler

import (
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/handler/instance"
	"github.com/kappital/kappital/pkg/handler/servicebinding"
)

// Type of handler
type Type string

const (
	// ServiceHandler the handler type of service
	ServiceHandler Type = "service"
	// InstanceHandler the handler type of instance
	InstanceHandler Type = "instance"
)

// IHandler the handler of synchronizing
type IHandler interface {
	BeforeInstall(obj interface{}) (bool, error)
	Install(obj interface{}) (bool, error)
	AfterInstall(obj interface{}) (bool, error)
	BeforeUpgrade(obj interface{}) (bool, error)
	Upgrade(obj interface{}) (bool, error)
	AfterUpgrade(obj interface{}) (bool, error)
	BeforeDelete(obj interface{}) (bool, error)
	Delete(obj interface{}) (bool, error)
	AfterDelete(obj interface{}) (bool, error)
}

var handlerMap map[Type]IHandler

func init() {
	handlerMap = map[Type]IHandler{
		ServiceHandler:  &servicebinding.Handler{},
		InstanceHandler: &instance.Handler{},
	}
}

// GetHandlerByType get the handler (instance or service binding) by its type
func GetHandlerByType(t Type) IHandler {
	handler, ok := handlerMap[t]
	if !ok {
		klog.Errorf("unknown handler type %s", t)
		return nil
	}
	if handler == nil {
		klog.Errorf("unknown handler type %s, get nil handler", t)
		return nil
	}
	return handler
}
