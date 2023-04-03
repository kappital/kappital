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
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"

	"github.com/kappital/kappital/pkg/models"
)

// ResourceOperation to manager the resource data in database
type ResourceOperation struct{}

// Insert resource information to database
func (r ResourceOperation) Insert(interface{}) error {
	return fmt.Errorf("ResourceModel do not have Insert method (without transaction), because it does not have this situation")
}

// InsertTx insert the resource with transaction
func (r ResourceOperation) InsertTx(obj interface{}, tx orm.TxOrmer) error {
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	resource.Generate(time.Now().UTC(), false)
	for i := range resource.Instances {
		err := InstanceOperation{}.InsertWithRelFk(*resource.Instances[i], resource, tx)
		if err != nil {
			return err
		}
	}
	_, err := tx.Insert(&resource)
	return models.IgnoreDBInsertIDError(err)
}

// InsertWithRelFk insert the resource with its relation foreign key
func (r ResourceOperation) InsertWithRelFk(obj interface{}, fk interface{}, tx orm.TxOrmer) error {
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	resource.Generate(time.Now().UTC(), false)
	serviceBinding, ok := fk.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	resource.Generate(time.Now().UTC(), false)
	resource.ServiceBinding = &serviceBinding
	return r.InsertTx(resource, tx)
}

// Get the resource from the database and filter by cols
func (r ResourceOperation) Get(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ResourceModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var item models.ResourceModel
	err := seter.One(&item)
	return item, err
}

// GetByPrimaryKey get the resource with its primary key (id)
func (r ResourceOperation) GetByPrimaryKey(id string) (interface{}, error) {
	resource := models.ResourceModel{ID: id}
	err := models.GetNewOrm().Read(&resource)
	return resource, err
}

// GetDetail of resource
func (r ResourceOperation) GetDetail(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ResourceModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var item models.ResourceModel
	if err := seter.One(&item); err != nil {
		return nil, err
	}
	_, err := models.GetNewOrm().LoadRelated(&item, "instances")
	return item, err
}

// GetList of resource
func (r ResourceOperation) GetList(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ResourceModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var items []models.ResourceModel
	_, err := seter.All(&items)
	return items, err
}

// GetListByFilter get the resource information by filter
func (r ResourceOperation) GetListByFilter(filter map[string][]interface{}) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ResourceModel{})
	for k, v := range filter {
		if v == nil {
			seter = seter.Filter(k+"__isnull", true)
		} else {
			seter = seter.Filter(k, v...)
		}
	}
	var items []models.ResourceModel
	_, err := seter.All(&items)
	return items, err
}

// IsExist does the resource information is existed in database with cols filter
func (r ResourceOperation) IsExist(cols map[string]string) bool {
	seter := models.GetNewOrm().QueryTable(models.InstanceModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	return seter.Exist()
}

// Update resource information
func (r ResourceOperation) Update(obj interface{}, cols ...string) error {
	sql := models.GetNewOrm()
	tx := models.NewTransaction(sql)
	err := tx.BeginTransaction()
	if err != nil {
		return err
	}
	defer models.Handler(&err, tx)
	return r.UpdateTx(obj, tx.GetTransaction(), cols...)
}

// UpdateTx update resource information with transaction
func (r ResourceOperation) UpdateTx(obj interface{}, tx orm.TxOrmer, cols ...string) error {
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	resource.Generate(time.Now().UTC(), true)
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	for i := range resource.Instances {
		err := InstanceOperation{}.UpdateTx(*resource.Instances[i], tx)
		if err != nil {
			return err
		}
	}
	old := models.ResourceModel{ID: resource.ID}
	if err := tx.Read(&old); err != nil {
		return err
	}
	resource.CreateTimestamp = old.CreateTimestamp
	_, err := tx.Update(&resource, cols...)
	return err
}

// Delete the resource
func (r ResourceOperation) Delete(obj interface{}) error {
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	sql := models.GetNewOrm()
	if len(resource.ID) == 0 {
		return fmt.Errorf("delete the resource missing the pk")
	}
	_, err := sql.Delete(&resource)
	return err
}

// DeleteTx the resource
func (r ResourceOperation) DeleteTx(obj interface{}, tx orm.TxOrmer) error {
	resource, ok := obj.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	if len(resource.ID) == 0 {
		return fmt.Errorf("delete the resource missing the pk")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	_, err := tx.Delete(&resource)
	return err
}
