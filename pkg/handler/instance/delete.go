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
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/dao/instance"
	co "github.com/kappital/kappital/pkg/utils/operations"
)

// BeforeDelete do some processes before delete the service instance
func (h *Handler) BeforeDelete(_ interface{}) (bool, error) {
	return false, nil
}

// Delete is to delete the whole cloud native service instance
func (h *Handler) Delete(obj interface{}) (bool, error) {
	si, ok := obj.(*internals.ServiceInstance)
	if !ok {
		klog.Errorf("invalid object type for upgrade instance handler, expected: models.ServiceInstance, "+
			"actual: %s", reflect.TypeOf(obj).Name())
		return false, nil
	}
	// Delete the instance cr in cluster
	repeat, err := h.deleteInstanceCR(*si)
	if err != nil {
		return true, err
	}
	// Delete the instance in database
	instanceDao := instance.Instance{}
	if err = instanceDao.Delete(*si); err != nil {
		return true, err
	}
	klog.Infof("delete the instance %s during the service %s", si.Name, si.ServiceName)
	return repeat, nil
}

func (h *Handler) deleteInstanceCR(si internals.ServiceInstance) (bool, error) {
	apiVersionSplit := strings.Split(si.APIVersion, "/")
	if len(apiVersionSplit) < 2 {
		return false, fmt.Errorf("the ServiceInstance's APIVersion is illeagle")
	}
	gvr := schema.GroupVersionResource{
		Group:    apiVersionSplit[0],
		Version:  apiVersionSplit[1],
		Resource: si.Resource,
	}
	err := co.GetClusterOperation().DeleteCustomResource(gvr, si.Name, si.Namespace)
	if err != nil {
		// this custom resource is deleting, and wait for it already deleted, in other words, the error is not found
		return true, nil
	}
	if !errors.IsNotFound(err) {
		klog.Errorf("delete custom resource %s in cluster failed, err: %s", si.Name, err)
		return true, err
	}
	return false, nil
}

// AfterDelete delete service instance does not need to implement this method
func (h *Handler) AfterDelete(_ interface{}) (bool, error) {
	return false, nil
}
