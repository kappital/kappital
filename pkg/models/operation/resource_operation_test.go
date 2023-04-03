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
	"github.com/kappital/kappital/pkg/models"
	"reflect"
	"testing"
)

var resource = ResourceOperation{}

func TestResourceOperation_Delete(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test ResourceOperation Delete (obj is not ResourceModel)",
			args:    args{models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation Delete (missing the pk)",
			args:    args{models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test ResourceOperation Delete",
			args: args{models.ResourceModel{ID: "delete-test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resource.Delete(tt.args.obj); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_DeleteTx(t *testing.T) {
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
			name:    "Test ResourceOperation DeleteTx (obj is not ResourceModel)",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation DeleteTx (missing the pk)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation DeleteTx (tx is nil)",
			args:    args{obj: models.ResourceModel{ID: "delete-test"}},
			wantErr: true,
		},
		{
			name: "Test ResourceOperation DeleteTx",
			args: args{obj: models.ResourceModel{ID: "delete-test"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = resource.DeleteTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("DeleteTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_Get(t *testing.T) {
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
			name: "Test ResourceOperation Get",
			args: args{cols: map[string]string{"id": "binding-id-1-resource-1"}},
			want: models.ResourceModel{
				ID:              "binding-id-1-resource-1",
				Kind:            "kind-1",
				Group:           "group-1",
				APIVersion:      "api-1",
				Resource:        "resource-1",
				CreateTimestamp: now,
				ServiceBinding:  &models.ServiceBindingModel{ID: "binding-id-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resource.Get(tt.args.cols)
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

func TestResourceOperation_GetByPrimaryKey(t *testing.T) {
	tests := []testForPrimaryKey{
		{
			name: "Test ResourceOperation GetByPrimaryKey",
			id:   "binding-id-1-resource-1",
			want: models.ResourceModel{
				ID:              "binding-id-1-resource-1",
				Kind:            "kind-1",
				Group:           "group-1",
				APIVersion:      "api-1",
				Resource:        "resource-1",
				CreateTimestamp: now,
				ServiceBinding:  &models.ServiceBindingModel{ID: "binding-id-1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resource.GetByPrimaryKey(tt.id)
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

func TestResourceOperation_GetDetail(t *testing.T) {
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
			name: "Test ResourceOperation GetDetail",
			args: args{cols: map[string]string{"id": "binding-id-1-resource-1"}},
			want: *(testServiceBindingModels[0].Resources[0]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resource.GetDetail(tt.args.cols)
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

func TestResourceOperation_GetList(t *testing.T) {
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
			name: "Test ResourceOperation GetList",
			args: args{cols: map[string]string{"id": "binding-id-1-resource-1"}},
			want: []models.ResourceModel{
				{
					ID:              "binding-id-1-resource-1",
					Kind:            "kind-1",
					Group:           "group-1",
					APIVersion:      "api-1",
					Resource:        "resource-1",
					CreateTimestamp: now,
					ServiceBinding:  &models.ServiceBindingModel{ID: "binding-id-1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resource.GetList(tt.args.cols)
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

func TestResourceOperation_GetListByFilter(t *testing.T) {
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
			name: "Test ResourceOperation GetListByFilter",
			args: args{map[string][]interface{}{"instances": nil, "resource": {"v1"}}},
			want: []models.ResourceModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resource.GetListByFilter(tt.args.filter)
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

func TestResourceOperation_Insert(t *testing.T) {
	type args struct {
		in0 interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test ResourceOperation Insert", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resource.Insert(tt.args.in0); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_InsertTx(t *testing.T) {
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
			name:    "Test ResourceOperation Insert (obj is not ResourceModel)",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation Insert (tx is nil)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test ResourceOperation Insert",
			args: args{
				obj: models.ResourceModel{
					ID: "binding-id-1-resource-1xxx",
					Instances: []*models.InstanceModel{
						{Resource: &models.ResourceModel{ID: "binding-id-1-resource-1xxx"}},
					},
				},
				tx: tx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = resource.InsertTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_InsertWithRelFk(t *testing.T) {
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
		{
			name:    "Test ResourceOperation InsertWithRelFk (obj is not ResourceModel)",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation InsertWithRelFk (fk is not ServiceBindingModel)",
			args:    args{obj: models.ResourceModel{}, fk: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation InsertWithRelFk (tx is nil)",
			args:    args{obj: models.ResourceModel{}, fk: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name: "Test ResourceOperation InsertWithRelFk",
			args: args{obj: models.ResourceModel{}, fk: models.ServiceBindingModel{}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = resource.InsertWithRelFk(tt.args.obj, tt.args.fk, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertWithRelFk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_IsExist(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test ResourceOperation IsExist",
			args: args{cols: map[string]string{"id": "binding-id-1-resource-1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resource.IsExist(tt.args.cols); got != tt.want {
				t.Errorf("IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceOperation_Update(t *testing.T) {
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
			name:    "Test ResourceOperation Update",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resource.Update(tt.args.obj, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceOperation_UpdateTx(t *testing.T) {
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
			name:    "Test ResourceOperation Update (tx is nil)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name:    "Test ResourceOperation Update (does not has id)",
			args:    args{obj: models.ResourceModel{}, tx: tx},
			wantErr: true,
		},
		{
			name: "Test ResourceOperation Update",
			args: args{obj: models.ResourceModel{ID: "binding-id-1-resource-1"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = resource.UpdateTx(tt.args.obj, tt.args.tx, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("UpdateTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
