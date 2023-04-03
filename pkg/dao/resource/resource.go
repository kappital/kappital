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
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/models"
	mo "github.com/kappital/kappital/pkg/models/operation"
)

// Resource the dao layer of resource for database CRUD
type Resource struct {
	resource mo.ResourceOperation
	binding  mo.ServiceBindingOperation
}

// Create insert a data record to the database
func (r Resource) Create(obj interface{}, params map[string]string) error {
	// get the whole service binding object
	tmp, err := r.binding.GetDetail(params)
	if err != nil {
		return err
	}
	binding, ok := tmp.(models.ServiceBindingModel)
	if !ok {
		return fmt.Errorf("the obj type is not models.ServiceBindingModel")
	}
	sql := orm.NewOrm()
	tx := models.NewTransaction(sql)
	if err = tx.BeginTransaction(); err != nil {
		return err
	}
	defer models.Handler(&err, tx)

	resources, ok := obj.([]internals.ServiceResource)
	if !ok {
		return fmt.Errorf("obj type is not Slice of ServiceInstance")
	}
	for _, resource := range resources {
		model := transformResourceToModel(resource)
		model.ServiceBinding = &binding
		if err = r.resource.InsertTx(model, tx.GetTransaction()); err != nil {
			return err
		}
	}
	return nil
}

// Get the resource from database and filter by cols
func (r Resource) Get(cols map[string]string) (interface{}, error) {
	obj, err := r.resource.Get(cols)
	if err != nil {
		return nil, err
	}
	return transModel2Resource(obj.(models.ResourceModel)), nil
}

// GetByPrimaryKey get the resource by the primary key (id)
func (r Resource) GetByPrimaryKey(id string) (interface{}, error) {
	obj, err := r.resource.GetByPrimaryKey(id)
	if err != nil {
		return nil, err
	}
	return transModel2Resource(obj.(models.ResourceModel)), nil
}

// GetList get resource list, and filter by cols
func (r Resource) GetList(cols map[string]string) (interface{}, error) {
	items, err := r.resource.GetList(cols)
	if err != nil {
		return nil, err
	}
	return transModelSlice2ResourceSlice(items.([]models.ResourceModel)), nil
}

// GetListByStatusSets Resource does not implement this method
func (r Resource) GetListByStatusSets(_ sets.String) (interface{}, error) {
	return nil, nil
}

// Update the resource to the database with cols
func (r Resource) Update(obj interface{}, cols ...string) error {
	internal, ok := obj.(internals.ServiceResource)
	if !ok {
		return fmt.Errorf("obj type is not ResourceModel")
	}
	old, err := r.resource.GetByPrimaryKey(internal.ID)
	if err != nil {
		return err
	}
	resource := transformResourceToModel(internal)
	resource.ServiceBinding = old.(models.ResourceModel).ServiceBinding
	return r.resource.Update(resource, cols...)
}

// UpdateStatusMsg Resource does not implement this method
func (r Resource) UpdateStatusMsg(_ interface{}, _, _ string) error {
	return nil
}

// Delete Resource does not implement this method
func (r Resource) Delete(_ interface{}) error {
	return nil
}

func transModel2Resource(resource models.ResourceModel) internals.ServiceResource {
	return internals.ServiceResource{
		ID:         resource.ID,
		Kind:       resource.Kind,
		Group:      resource.Group,
		APIVersion: resource.APIVersion,
		Resource:   resource.Resource,
	}
}

func transModelSlice2ResourceSlice(resources []models.ResourceModel) []internals.ServiceResource {
	result := make([]internals.ServiceResource, 0, len(resources))
	for _, resource := range resources {
		result = append(result, transModel2Resource(resource))
	}
	return result
}

func transformResourceToModel(resource internals.ServiceResource) models.ResourceModel {
	return models.ResourceModel{
		ID:         resource.ID,
		Kind:       resource.Kind,
		Group:      resource.Group,
		APIVersion: resource.APIVersion,
		Resource:   resource.Resource,
	}
}
