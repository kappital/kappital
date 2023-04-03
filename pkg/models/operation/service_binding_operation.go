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

// ServiceBindingOperation to manager the service binding data in database
type ServiceBindingOperation struct{}

// Insert service binding information to database
func (s ServiceBindingOperation) Insert(obj interface{}) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	sb.Generate(time.Now().UTC(), false)
	tx := models.NewTransaction(models.GetNewOrm())
	err := tx.BeginTransaction()
	if err != nil {
		return err
	}
	defer models.Handler(&err, tx)

	for i := range sb.Resources {
		err = ResourceOperation{}.InsertWithRelFk(*sb.Resources[i], sb, tx.GetTransaction())
		if err != nil {
			return err
		}
	}
	_, err = tx.GetTransaction().Insert(&sb)
	return models.IgnoreDBInsertIDError(err)
}

// InsertTx service binding information to database with transaction
func (s ServiceBindingOperation) InsertTx(obj interface{}, tx orm.TxOrmer) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	sb.Generate(time.Now().UTC(), false)

	for i := range sb.Resources {
		err := ResourceOperation{}.InsertWithRelFk(*sb.Resources[i], sb, tx)
		if err != nil {
			return err
		}
	}
	_, err := tx.Insert(&sb)
	return models.IgnoreDBInsertIDError(err)
}

// InsertWithRelFk service binding does not need to implement this method
func (s ServiceBindingOperation) InsertWithRelFk(interface{}, interface{}, orm.TxOrmer) error {
	return fmt.Errorf("ServiceBindingModel do not have InsertWithRelFk method, " +
		"because the ServiceBindingModel do not have fk")
}

// Get the service binding from the database and filter by cols
func (s ServiceBindingOperation) Get(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ServiceBindingModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var item models.ServiceBindingModel
	err := seter.One(&item)
	return item, err
}

// GetByPrimaryKey get the service binding with its primary key (id)
func (s ServiceBindingOperation) GetByPrimaryKey(id string) (interface{}, error) {
	binding := models.ServiceBindingModel{ID: id}
	err := models.GetNewOrm().Read(&binding)
	return binding, err
}

// GetDetail of service binding
func (s ServiceBindingOperation) GetDetail(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ServiceBindingModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var item models.ServiceBindingModel
	if err := seter.One(&item); err != nil {
		return nil, err
	}
	if _, err := models.GetNewOrm().LoadRelated(&item, "resources"); err != nil {
		return nil, err
	}
	for i := range item.Resources {
		if _, err := models.GetNewOrm().LoadRelated(item.Resources[i], "instances"); err != nil {
			return nil, err
		}
	}
	return item, nil
}

// GetList of service binding
func (s ServiceBindingOperation) GetList(cols map[string]string) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ServiceBindingModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	var items []models.ServiceBindingModel
	_, err := seter.All(&items)
	return items, err
}

// GetListByFilter get the service binding information by filter
func (s ServiceBindingOperation) GetListByFilter(filter map[string][]interface{}) (interface{}, error) {
	seter := models.GetNewOrm().QueryTable(models.ServiceBindingModel{})
	for k, v := range filter {
		if v == nil {
			seter = seter.Filter(k+"__isnull", true)
		} else {
			seter = seter.Filter(k, v...)
		}
	}
	var items []models.ServiceBindingModel
	_, err := seter.All(&items)
	return items, err
}

// IsExist does the service binding information is existed in database with cols filter
func (s ServiceBindingOperation) IsExist(cols map[string]string) bool {
	seter := models.GetNewOrm().QueryTable(models.ServiceBindingModel{})
	for k, v := range cols {
		seter = seter.Filter(k, v)
	}
	return seter.Exist()
}

// Update service binding information
func (s ServiceBindingOperation) Update(obj interface{}, cols ...string) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	sb.Generate(time.Now().UTC(), true)
	sql := models.GetNewOrm()
	old := models.ServiceBindingModel{ID: sb.ID}
	if err := sql.Read(&old); err != nil {
		return err
	}
	sb.CreateTime = old.CreateTime
	_, err := sql.Update(&sb, cols...)
	return err
}

// UpdateTx update service binding information with transaction
func (s ServiceBindingOperation) UpdateTx(obj interface{}, tx orm.TxOrmer, cols ...string) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	for i := range sb.Resources {
		err := ResourceOperation{}.UpdateTx(*sb.Resources[i], tx)
		if err != nil {
			return err
		}
	}

	old := models.ServiceBindingModel{ID: sb.ID}
	if err := tx.Read(&old); err != nil {
		return err
	}
	sb.CreateTime = old.CreateTime
	_, err := tx.Update(&sb, cols...)
	return err
}

// Delete the service binding
func (s ServiceBindingOperation) Delete(obj interface{}) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	sql := models.GetNewOrm()
	if len(sb.ID) == 0 {
		if err := sql.Read(&sb, "name", "cluster_name"); err != nil {
			if errors.Is(err, orm.ErrNoRows) {
				return nil
			}
			return err
		}
	}
	_, err := sql.Delete(&sb)
	return err
}

// DeleteTx the service binding with transaction
func (s ServiceBindingOperation) DeleteTx(obj interface{}, tx orm.TxOrmer) error {
	sb, ok := obj.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("obj type is not ServiceBindingModel")
	}
	if tx == nil {
		return fmt.Errorf("transaction should not be nil")
	}
	if len(sb.ID) == 0 {
		if err := tx.Read(&sb, "name", "cluster_name"); err != nil {
			if errors.Is(err, orm.ErrNoRows) {
				return nil
			}
			return err
		}
	}
	_, err := tx.Delete(&sb)
	return err
}
