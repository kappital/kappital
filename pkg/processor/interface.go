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

package processor

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/handler"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/resource"
)

const (
	// WorkersPerProcessor indicates the number of asynchronous threads for each type of event processing
	WorkersPerProcessor = 5
)

// IProcessor is asynchronous thread interface
type IProcessor interface {
	List() ([]interface{}, error)
	Process(opType string, obj interface{})
	Run(stopCh <-chan struct{})
	ProcessType() string
	ProcessObject() interface{}
}

var processors map[string]IProcessor

func init() {
	processors = map[string]IProcessor{}
	registerProcess(NewProcessor(apis.OperatorProcessor, &internals.ServiceBinding{},
		handler.ServiceHandler, resource.ServiceBindingType, WorkersPerProcessor,
		sets.NewString(models.StatusInstalling, models.StatusUpgrading, models.StatusRollingBack,
			models.StatusDeleting)))
	registerProcess(NewProcessor(apis.InstanceProcessor, &internals.ServiceInstance{},
		handler.InstanceHandler, resource.ServiceInstanceType, WorkersPerProcessor,
		sets.NewString(models.StatusInitializing, models.StatusUpgrading, models.StatusDeleting)))
}
func registerProcess(process IProcessor) {
	if processors == nil {
		processors = make(map[string]IProcessor)
	}
	processors[process.ProcessType()] = process
}

// GetProcesses returns the type of process
func GetProcesses() map[string]IProcessor {
	return processors
}

// StartAllProcessors start all register process
func StartAllProcessors(stopCh <-chan struct{}) {
	for _, process := range processors {
		go process.Run(stopCh)
	}
}
