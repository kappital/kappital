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

package resource

import (
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
)

// Type of the service instance and binding resource
type Type string

const (
	// ServiceInstanceType the service instance resource type name
	ServiceInstanceType Type = "instance"
	// ServiceBindingType the service binding resource type name
	ServiceBindingType Type = "servicebinding"
)

// IResource of the interface to use the database and/or cluster
type IResource interface {
	GetResourceType() Type
	GetCommonDBObject(key string) (interface{}, error)
	GetObjectProcessTime(obj interface{}) time.Time
	GetObjectStatus(obj interface{}) string
	GetObjectID(obj interface{}) string
	GetObjectListByStatusSets(status sets.String) ([]interface{}, error)
	UpdateProcessFailed(obj interface{}, status string, msg string) error
	UpdateObjProcessTime(obj interface{}, processTime time.Time) error
	GetObjUpdateTime(obj interface{}) time.Time
}

var insResource InstanceResource
var serviceBindingResource ServiceBindingResource

// GetResourceByType get the resource by type name
func GetResourceByType(rType Type) IResource {
	switch rType {
	case ServiceInstanceType:
		return &insResource
	case ServiceBindingType:
		return &serviceBindingResource
	default:
		return nil
	}
}
