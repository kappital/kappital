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
	"os"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kappital/kappital/pkg/models"
)

var (
	now = time.Now().UTC()

	testServiceBindingModels = []models.ServiceBindingModel{
		{
			ID:          "binding-id-1",
			Name:        "test-name-1",
			Version:     "v1.1.1",
			Namespace:   "kappital-system",
			ServiceName: "test-name-1",
			ClusterName: "default",
			ServiceID:   "service-id-1",
			VersionID:   "service-id-1-id-v-1-1",
			Status:      "Success",
			CreateTime:  now,
			Resources: []*models.ResourceModel{
				{
					ID:              "binding-id-1-resource-1",
					Kind:            "kind-1",
					Group:           "group-1",
					APIVersion:      "api-1",
					Resource:        "resource-1",
					CreateTimestamp: now,
					Instances: []*models.InstanceModel{
						{
							ID:               "instance-1",
							Kind:             "kind-1",
							APIVersion:       "api-1",
							Name:             "instance-1",
							Namespace:        "default",
							ServiceBindingID: "binding-id-1",
							Status:           "Success",
							ServiceID:        "service-id-1",
							ServiceName:      "test-name-1",
							ClusterName:      "default",
							CreateTimestamp:  now,
							Resource:         &models.ResourceModel{ID: "binding-id-1-resource-1"},
						},
					},
					ServiceBinding: &models.ServiceBindingModel{ID: "binding-id-1"},
				},
			},
		},
	}
	testNonDetailServiceBinding = models.ServiceBindingModel{
		ID:          "binding-id-1",
		Name:        "test-name-1",
		Version:     "v1.1.1",
		Namespace:   "kappital-system",
		ServiceName: "test-name-1",
		ClusterName: "default",
		ServiceID:   "service-id-1",
		VersionID:   "service-id-1-id-v-1-1",
		Status:      "Success",
		CreateTime:  now,
	}
)

func TestMain(m *testing.M) {
	p := gomonkey.ApplyFunc(time.Now, func() time.Time {
		return time.Time{}
	})
	defer p.Reset()
	var err error
	if err = orm.RegisterDriver("sqlite3", orm.DRSqlite); err != nil {
		fmt.Printf("cannot register orm driver for operation, err: %v\n", err)
		return
	}
	if err = orm.RegisterDataBase("default", "sqlite3", "./test-operation.db"); err != nil {
		fmt.Printf("cannot register database for operation, err: %v\n", err)
		return
	}
	orm.RegisterModel(new(models.ServiceBindingModel), new(models.ResourceModel), new(models.InstanceModel))
	if err = orm.RunSyncdb("default", false, true); err != nil {
		fmt.Printf("run sync db error %v", err)
		return
	}
	preInsertServiceBindingValues()
	m.Run()
	if err = os.Remove("./test-operation.db"); err != nil {
		fmt.Printf("cannot remove unit test db for operation, err: %v\n", err)
	}
	if err = os.Remove("./test-operation.db-journal"); err != nil {
		fmt.Printf("cannot remove unit test db journal for operation, err: %v\n", err)
	}
}

func preInsertServiceBindingValues() {
	for _, item := range testServiceBindingModels {
		for _, resource := range item.Resources {
			if _, err := models.GetNewOrm().Insert(resource); models.IgnoreDBInsertIDError(err) != nil {
				fmt.Printf("cannot Insert resource with id: %s for service binding id %s, because: %v", resource.ID, item.ID, err)
				return
			}
			if err := insertInstances(resource.Instances); err != nil {
				fmt.Printf("cannot Insert resource id [%s] instances, because: %v", resource.ID, err)
				return
			}
		}
		if _, err := models.GetNewOrm().Insert(&item); models.IgnoreDBInsertIDError(err) != nil {
			fmt.Printf("cannot Insert service with id: %s, because: %v\n", item.ID, err)
			return
		}
	}
}

func insertInstances(instances []*models.InstanceModel) error {
	for _, instance := range instances {
		if _, err := models.GetNewOrm().Insert(instance); err != nil {
			return err
		}
	}
	return nil
}

func ignoreDBLockError(err error) error {
	if err != nil && err.Error() == "database is locked" {
		return nil
	}
	return err
}

// Test structs to reduce the duplicate code

type testForPrimaryKey struct {
	name    string
	id      string
	want    interface{}
	wantErr bool
}
