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
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/handler"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/resource"
	"github.com/kappital/kappital/pkg/watcher"
)

const (
	retryInterval = 5 * time.Second
	// LongestDurationForProcess defines the timeout interval for asynchronous event processing
	LongestDurationForProcess = 20 * time.Minute
)

// Processor the struct of process attribute
type Processor struct {
	workQueue       workqueue.RateLimitingInterface
	workers         int
	resource        resource.IResource
	handler         handler.IHandler
	processName     string
	processorObject interface{}
	careStatusSet   sets.String
}

// NewProcessor return the entity object of process
func NewProcessor(name string, processorObj interface{}, handlerType handler.Type, resourceType resource.Type,
	workerCounts int, actionSets sets.String) *Processor {
	return &Processor{
		workers: workerCounts,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(),
			fmt.Sprintf("%s-processor", name)),
		resource:        resource.GetResourceByType(resourceType),
		handler:         handler.GetHandlerByType(handlerType),
		processName:     name,
		processorObject: processorObj,
		careStatusSet:   actionSets,
	}
}

// ProcessType return the process name
func (p *Processor) ProcessType() string {
	return p.processName
}

// ProcessObject return the object of process
func (p *Processor) ProcessObject() interface{} {
	return p.processorObject
}

// Run start to run process
func (p *Processor) Run(stopCh <-chan struct{}) {
	defer p.workQueue.ShutDown()

	klog.Infof("Starting %s processor", p.processName)
	defer klog.Infof("Shutting down %s processor", p.processName)

	for i := 0; i < p.workers; i++ {
		go wait.Until(p.sync, 0, stopCh)
	}

	<-stopCh
}

func (p *Processor) sync() {
	key, quit := p.workQueue.Get()
	if quit {
		klog.V(3).Infof("quit syncing work queue")
		return
	}
	defer p.workQueue.Done(key)

	retry, err := p.syncHandler(key.(string))
	if err != nil {
		klog.Errorf("Error processing %s %s: %s", p.processName, key, err.Error())
	}
	if retry {
		p.workQueue.AddAfter(key, retryInterval)
	}
}

func (p *Processor) syncHandler(name string) (retry bool, err error) {
	obj, err := p.resource.GetCommonDBObject(name)
	if err != nil {
		return !errors.Is(err, orm.ErrNoRows), err
	}
	status := p.resource.GetObjectStatus(obj)
	defer func() {
		if !retry {
			// reset process timeout when retry false, error nil
			innerErr := p.resource.UpdateObjProcessTime(obj, time.Time{})
			if innerErr != nil {
				retry = true
				err = innerErr
				return
			}
			return
		}
		// check handle timeout
		if !p.subProcessTimeout(obj) {
			return
		}
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		innerErr := p.resource.UpdateProcessFailed(obj, getFailedStatus(status),
			fmt.Sprintf("timeout to handle resource %s, status %s, last error: %s",
				p.resource.GetResourceType(), p.resource.GetObjectStatus(obj), errMsg))
		if innerErr != nil {
			retry = true
			err = innerErr
			return
		}
		retry = false
		err = fmt.Errorf("timeout to handle resource %s, status %s, last error:%s",
			p.resource.GetResourceType(), p.resource.GetObjectStatus(obj), errMsg)
	}()
	if p.careStatusSet.Has(status) {
		switch status {
		case models.StatusDeleting:
			return processStatus(obj, []func(interface{}) (bool, error){p.handler.BeforeDelete,
				p.handler.Delete, p.handler.AfterDelete})
		case models.StatusUpgrading, models.StatusRollingBack:
			return processStatus(obj, []func(interface{}) (bool, error){p.handler.BeforeUpgrade,
				p.handler.Upgrade, p.handler.AfterUpgrade})
		case models.StatusInstalling, models.StatusInitializing:
			return processStatus(obj, []func(interface{}) (bool, error){p.handler.BeforeInstall,
				p.handler.Install, p.handler.AfterInstall})
		}
	}
	return false, nil
}

func (p *Processor) subProcessTimeout(obj interface{}) bool {
	if time.Now().UTC().After(p.resource.GetObjUpdateTime(obj).Add(LongestDurationForProcess)) {
		return true
	}
	if p.resource.GetObjectProcessTime(obj).IsZero() {
		return false
	}
	return time.Now().After(p.resource.GetObjectProcessTime(obj))
}

func getFailedStatus(status string) string {
	switch status {
	case models.StatusDeleting:
		return models.StatusDeleteFailed
	case models.StatusUpgrading:
		return models.StatusUpgradeFailed
	case models.StatusRollingBack:
		return models.StatusRollBackFailed
	default:
		return models.StatusFailed
	}
}

// Process implements watcher.Processor interface
func (p *Processor) Process(opType string, obj interface{}) {
	switch opType {
	case watcher.OPList, watcher.OPCreate:
		fallthrough
	case watcher.OPUpdate, watcher.OPDelete:
		if p.careStatusSet.Has(p.resource.GetObjectStatus(obj)) {
			p.workQueue.Add(p.resource.GetObjectID(obj))
		}
	}
}

// List implements watcher.Processor interface
func (p *Processor) List() ([]interface{}, error) {
	res, err := p.resource.GetObjectListByStatusSets(p.careStatusSet)
	if err != nil {
		klog.Errorf("list %s objects failed, error: %v", p.processName, err)
	}
	return res, nil
}

func processStatus(obj interface{}, funcs []func(interface{}) (bool, error)) (bool, error) {
	for _, f := range funcs {
		retry, err := f(obj)
		if err != nil || retry {
			return retry, err
		}
	}
	return false, nil
}
