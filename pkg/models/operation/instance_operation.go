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
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"

	"github.com/kappital/kappital/pkg/models"
)

// InstanceOperation to manager the instance data in database
type InstanceOperation struct{}

// Insert instance information to database
func (i InstanceOperation) Insert(interface{}) error {
	return fmt.Errorf("InstanceModel do not have Insert method (without transaction), because it does not have this situation")
}

// InsertTx instance information to database with transaction
func (i InstanceOperation) InsertTx(obj interface{}, tx orm.TxOrmer) error {
	repo, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not InstanceModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	repo.Generate(time.Now().UTC(), false)
	_, err := tx.Insert(&repo)
	return models.IgnoreDBInsertIDError(err)
}

// InsertWithRelFk insert the instance with its relation foreign key
func (i InstanceOperation) InsertWithRelFk(obj interface{}, fk interface{}, tx orm.TxOrmer) error {
	instance, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not InstanceModel")
	}
	instance.Generate(time.Now().UTC(), false)
	resource, ok := fk.(models.ResourceModel)
	if !ok {
		return fmt.Errorf("instance.InsertWithRelFk relFk is not ResourceModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	instance.Resource = &resource
	return i.InsertTx(instance, tx)
}

// Get the instance from the database and filter by cols
func (i InstanceOperation) Get(cols map[string]string) (interface{}, error) {
	setter := models.GetNewOrm().QueryTable(models.InstanceModel{})
	for k, v := range cols {
		setter = setter.Filter(k, v)
	}
	var item models.InstanceModel
	err := setter.One(&item)
	return item, err
}

// GetByPrimaryKey get the instance with its primary key (id)
func (i InstanceOperation) GetByPrimaryKey(id string) (interface{}, error) {
	instance := models.InstanceModel{ID: id}
	err := models.GetNewOrm().Read(&instance)
	return instance, err
}

// GetDetail of instance
func (i InstanceOperation) GetDetail(cols map[string]string) (interface{}, error) {
	return i.Get(cols)
}

// GetList of instance
func (i InstanceOperation) GetList(cols map[string]string) (interface{}, error) {
	setter := models.GetNewOrm().QueryTable(models.InstanceModel{})
	for k, v := range cols {
		setter = setter.Filter(k, v)
	}
	var items []models.InstanceModel
	_, err := setter.All(&items)
	return items, err
}

// GetListByFilter get the instance information by filter
func (i InstanceOperation) GetListByFilter(filter map[string][]interface{}) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.InstanceModel{})
	for k, v := range filter {
		if v == nil {
			seter = seter.Filter(k+"__isnull", true)
		} else {
			seter = seter.Filter(k, v...)
		}
	}
	var items []models.InstanceModel
	_, err := seter.All(&items)
	return items, err
}

// IsExist does the instance information is existed in database with cols filter
func (i InstanceOperation) IsExist(cols map[string]string) bool {
	seter := models.GetNewOrm().QueryTable(models.InstanceModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	return seter.Exist()
}

// Update instance information
func (i InstanceOperation) Update(obj interface{}, cols ...string) error {
	instance, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not InstanceModel")
	}
	instance.Generate(time.Now().UTC(), true)
	old := models.InstanceModel{ID: instance.ID}
	sql := models.GetNewOrm()
	if err := sql.Read(&old); err != nil {
		return err
	}
	instance.CreateTimestamp = old.CreateTimestamp
	_, err := sql.Update(&instance, cols...)
	return err
}

// UpdateTx update instance information with transaction
func (i InstanceOperation) UpdateTx(obj interface{}, tx orm.TxOrmer, cols ...string) error {
	instance, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not RepositoryModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	instance.Generate(time.Now().UTC(), true)
	old := models.InstanceModel{ID: instance.ID}
	if err := tx.Read(&old); err != nil {
		return err
	}
	instance.CreateTimestamp = old.CreateTimestamp
	_, err := tx.Update(&instance, cols...)
	return err
}

// Delete the instance
func (i InstanceOperation) Delete(obj interface{}) error {
	sql := models.GetNewOrm()
	instance, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not InstanceModel")
	}
	if len(instance.ID) == 0 {
		if err := sql.Read(&instance, "name", "namespace", "cluster_name"); err != nil {
			if errors.Is(err, orm.ErrNoRows) {
				return nil
			}
			return err
		}
	}
	_, err := sql.Delete(&instance)
	return err
}

// DeleteTx the instance with transaction
func (i InstanceOperation) DeleteTx(obj interface{}, tx orm.TxOrmer) error {
	instance, ok := obj.(models.InstanceModel)
	if !ok {
		return fmt.Errorf("obj type is not InstanceModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	if len(instance.ID) == 0 {
		if err := tx.Read(&instance, "name", "namespace", "cluster_name"); err != nil {
			if errors.Is(err, orm.ErrNoRows) {
				return nil
			}
			return err
		}
	}
	_, err := tx.Delete(&instance)
	return err
}
