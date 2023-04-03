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

package operation

import (
	"github.com/beego/beego/v2/client/orm"
)

// DatabaseOperation for service instance, service package, and etc.
type DatabaseOperation interface {
	Insert(obj interface{}) error
	InsertTx(obj interface{}, tx orm.TxOrmer) error
	InsertWithRelFk(obj interface{}, fk interface{}, tx orm.TxOrmer) error

	Get(cols map[string]string) (interface{}, error)
	GetByPrimaryKey(id string) (interface{}, error)
	GetDetail(cols map[string]string) (interface{}, error)
	GetList(cols map[string]string) (interface{}, error)
	GetListByFilter(filter map[string][]interface{}) (interface{}, error)
	IsExist(cols map[string]string) bool

	Update(obj interface{}, cols ...string) error
	UpdateTx(obj interface{}, tx orm.TxOrmer, cols ...string) error

	Delete(obj interface{}) error
	DeleteTx(obj interface{}, tx orm.TxOrmer) error
}
