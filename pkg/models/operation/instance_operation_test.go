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
	"reflect"
	"testing"

	"github.com/beego/beego/v2/client/orm"

	"github.com/kappital/kappital/pkg/models"
)

var instance = InstanceOperation{}

func TestInstanceOperation_Delete(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InstanceOperation Delete (obj is not VersionModel)",
			args:    args{models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation Delete (cannot get the instance)",
			args: args{models.InstanceModel{}},
		},
		{
			name: "Test InstanceOperation Delete",
			args: args{models.InstanceModel{ID: "delete-test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := instance.Delete(tt.args.obj); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceOperation_DeleteTx(t *testing.T) {
	tx, err := models.GetNewOrm().Begin()
	if err != nil {
		t.Errorf("cannot get the tx, err: %s", err)
	}
	type args struct {
		obj interface{}
		tx  orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InstanceOperation DeleteTx (obj is not InstanceModel)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation DeleteTx (tx is nil)",
			args:    args{obj: models.InstanceModel{}, tx: nil},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation DeleteTx (cannot get the instance)",
			args: args{obj: models.InstanceModel{}, tx: tx},
		},
		{
			name: "Test InstanceOperation DeleteTx",
			args: args{obj: models.InstanceModel{ID: "delete-test"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = instance.DeleteTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("DeleteTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceOperation_Get(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test InstanceOperation Get",
			args: args{cols: map[string]string{"id": "instance-1"}},
			want: models.InstanceModel{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Get(tt.args.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_GetByPrimaryKey(t *testing.T) {
	tests := []testForPrimaryKey{
		{
			name: "Test InstanceOperation GetByPrimaryKey",
			id:   "instance-1",
			want: models.InstanceModel{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.GetByPrimaryKey(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByPrimaryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByPrimaryKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_GetDetail(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test InstanceOperation GetDetail",
			args: args{cols: map[string]string{"id": "instance-1"}},
			want: models.InstanceModel{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.GetDetail(tt.args.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_GetList(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test InstanceOperation GetDetail",
			args: args{cols: map[string]string{"id": "instance-1"}},
			want: []models.InstanceModel{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.GetList(tt.args.cols)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_GetListByFilter(t *testing.T) {
	type args struct {
		filter map[string][]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test InstanceOperation GetListByFilter",
			args: args{map[string][]interface{}{"namespace": nil, "status": {"v1"}}},
			want: []models.InstanceModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.GetListByFilter(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListByFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetListByFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_Insert(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Test InstanceOperation Insert", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := instance.Insert(nil); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceOperation_InsertTx(t *testing.T) {
	tx, err := models.GetNewOrm().Begin()
	if err != nil {
		t.Errorf("cannot get the tx, err: %s", err)
	}
	type args struct {
		obj interface{}
		tx  orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InstanceOperation InsertTx (obj is not InstanceModel)",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation InsertTx (tx is nil)",
			args:    args{obj: models.InstanceModel{}},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation InsertTx",
			args: args{obj: models.InstanceModel{ID: "insert2"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = instance.InsertTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertTx() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.tx != nil {
				if err = tt.args.tx.Rollback(); err != nil {
					t.Errorf("cannot rollback the insert, err: %s", err)
				}
			}
		})
	}
}

func TestInstanceOperation_InsertWithRelFk(t *testing.T) {
	tx, err := models.GetNewOrm().Begin()
	if err != nil {
		t.Errorf("cannot get the tx, err: %s", err)
	}
	type args struct {
		obj interface{}
		fk  interface{}
		tx  orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test InstanceOperation InsertWithRelFK (invalid input obj)", wantErr: true},
		{
			name:    "Test InstanceOperation InsertWithRelFK (invalid input fk obj)",
			args:    args{obj: models.InstanceModel{ID: "test-instance-1"}, fk: nil},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation InsertWithRelFK (tx is nil)",
			args:    args{obj: models.InstanceModel{ID: "test-instance-1"}, fk: models.ResourceModel{ID: "test-resource"}},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation InsertWithRelFK (valid input)",
			args: args{obj: models.InstanceModel{ID: "test-instance-2"}, fk: models.ResourceModel{ID: "test-resource"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = instance.InsertWithRelFk(tt.args.obj, tt.args.fk, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertWithRelFk() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.tx != nil {
				if err = tt.args.tx.Rollback(); err != nil {
					t.Errorf("cannot rollback the insert, err: %s", err)
				}
			}
		})
	}
}

func TestInstanceOperation_IsExist(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test InstanceOperation IsExist",
			args: args{cols: map[string]string{"id": "instance-1"}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := instance.IsExist(tt.args.cols); got != tt.want {
				t.Errorf("IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOperation_Update(t *testing.T) {
	type args struct {
		obj  interface{}
		cols []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InstanceOperation Update (invalid input)",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation Update (cannot find the obj)",
			args:    args{obj: models.InstanceModel{ID: "x"}},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation Update",
			args: args{
				obj: models.InstanceModel{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := instance.Update(tt.args.obj, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceOperation_UpdateTx(t *testing.T) {
	tx, err := models.GetNewOrm().Begin()
	if err != nil {
		t.Errorf("cannot get the tx, err: %s", err)
	}
	type args struct {
		obj  interface{}
		tx   orm.TxOrmer
		cols []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test InstanceOperation Update (invalid input)",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation Update (tx is nil)",
			args:    args{obj: models.InstanceModel{ID: "x"}, tx: nil},
			wantErr: true,
		},
		{
			name:    "Test InstanceOperation Update (not found)",
			args:    args{obj: models.InstanceModel{ID: "id-1"}, tx: tx},
			wantErr: true,
		},
		{
			name: "Test InstanceOperation Update",
			args: args{
				obj: models.InstanceModel{
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
				tx: tx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = instance.UpdateTx(tt.args.obj, tt.args.tx, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("UpdateTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
