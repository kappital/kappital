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

var binding = ServiceBindingOperation{}

func TestServiceBindingOperation_Delete(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test ServiceBindingOperation Delete (obj is not ServiceBindingOperation)",
			args:    args{models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation Delete (cannot get the binding)",
			args: args{models.ServiceBindingModel{}},
		},
		{
			name: "Test ServiceBindingOperation Delete",
			args: args{models.ServiceBindingModel{ID: "binding-id-1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := binding.Delete(tt.args.obj); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_DeleteTx(t *testing.T) {
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
			name:    "Test ServiceBindingOperation DeleteTx (obj is not ServiceBindingModel)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name:    "Test ServiceBindingOperation DeleteTx (tx is nil)",
			args:    args{obj: models.ServiceBindingModel{}, tx: nil},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation DeleteTx (cannot get the binding)",
			args: args{obj: models.ServiceBindingModel{}, tx: tx},
		},
		{
			name: "Test ServiceBindingOperation DeleteTx",
			args: args{obj: models.ServiceBindingModel{ID: "binding-id-1"}, tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = binding.DeleteTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("DeleteTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_Get(t *testing.T) {
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
			name: "Test ServiceBindingOperation Get",
			args: args{map[string]string{"id": "binding-id-1"}},
			want: testNonDetailServiceBinding,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := binding.Get(tt.args.cols)
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

func TestServiceBindingOperation_GetByPrimaryKey(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Test ServiceBindingOperation GetByPrimaryKey",
			args: args{"binding-id-1"},
			want: testNonDetailServiceBinding,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := binding.GetByPrimaryKey(tt.args.id)
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

func TestServiceBindingOperation_GetDetail(t *testing.T) {
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
			name: "Test ServiceBindingOperation GetDetail",
			args: args{map[string]string{"id": "binding-id-1"}},
			want: testServiceBindingModels[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := binding.GetDetail(tt.args.cols)
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

func TestServiceBindingOperation_GetList(t *testing.T) {
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
			name: "Test ServiceBindingOperation GetList",
			args: args{map[string]string{"id": "binding-id-1"}},
			want: []models.ServiceBindingModel{testNonDetailServiceBinding},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := binding.GetList(tt.args.cols)
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

func TestServiceBindingOperation_GetListByFilter(t *testing.T) {
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
			name: "Test ServiceBindingOperation  GetListByFilter",
			args: args{map[string][]interface{}{"name": nil, "version": {"v1"}}},
			want: []models.ServiceBindingModel{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := binding.GetListByFilter(tt.args.filter)
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

func TestServiceBindingOperation_Insert(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test ServiceBindingOperation Insert (obj is not ServiceBindingOperation)",
			args:    args{models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation Insert (obj is not ServiceBindingOperation)",
			args: args{models.ServiceBindingModel{Resources: []*models.ResourceModel{{}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := binding.Insert(tt.args.obj); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_InsertTx(t *testing.T) {
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
			name:    "Test ServiceBindingOperation Insert (obj is not ServiceBindingOperation)",
			args:    args{obj: models.ResourceModel{}},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation Insert (first time insert)",
			args: args{obj: models.ServiceBindingModel{}, tx: tx},
		},
		{
			name: "Test ServiceBindingOperation Insert (second time insert)",
			args: args{obj: testServiceBindingModels[0], tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = binding.InsertTx(tt.args.obj, tt.args.tx); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_InsertWithRelFk(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Test ServiceBindingOperation InsertWithRelFk", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := binding.InsertWithRelFk(nil, nil, nil); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("InsertWithRelFk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_IsExist(t *testing.T) {
	type args struct {
		cols map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test ServiceBindingOperation IsExist",
			args: args{cols: map[string]string{"id": "binding-id-1"}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binding.IsExist(tt.args.cols); got != tt.want {
				t.Errorf("IsExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceBindingOperation_Update(t *testing.T) {
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
			name:    "Test ServiceBindingOperation Update (invalid input)",
			wantErr: true,
		},
		{
			name:    "Test ServiceBindingOperation Update (cannot get the old input)",
			args:    args{obj: models.ServiceBindingModel{}},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation Update (cannot get the old input)",
			args: args{obj: models.ServiceBindingModel{ID: "binding-id-1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := binding.Update(tt.args.obj, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBindingOperation_UpdateTx(t *testing.T) {
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
			name:    "Test ServiceBindingOperation Update (invalid input)",
			wantErr: true,
		},
		{
			name:    "Test ServiceBindingOperation Update (tx is nil)",
			args:    args{obj: models.ServiceBindingModel{ID: "x"}, tx: nil},
			wantErr: true,
		},
		{
			name:    "Test ServiceBindingOperation Update (not found)",
			args:    args{obj: models.ServiceBindingModel{ID: "id-1"}, tx: tx},
			wantErr: true,
		},
		{
			name: "Test ServiceBindingOperation Update ",
			args: args{obj: testServiceBindingModels[0], tx: tx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = binding.UpdateTx(tt.args.obj, tt.args.tx, tt.args.cols...); (ignoreDBLockError(err) != nil) != tt.wantErr {
				t.Errorf("UpdateTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
