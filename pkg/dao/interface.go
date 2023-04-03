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

package dao

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

// Dao layer interface
type Dao interface {
	Create(obj interface{}, params map[string]string) error
	Get(cols map[string]string) (interface{}, error)
	GetByPrimaryKey(id string) (interface{}, error)
	GetList(cols map[string]string) (interface{}, error)
	GetListByStatusSets(status sets.String) (interface{}, error)
	Update(obj interface{}, cols ...string) error
	UpdateStatusMsg(obj interface{}, status, msg string) error
	Delete(obj interface{}) error
}
