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

package models

import (
	"testing"
)

func TestFakeTransactionImpl_BeginTransaction(t1 *testing.T) {
	t := &FakeTransactionImpl{}
	if err := t.BeginTransaction(); err != nil {
		t1.Errorf("BeginTransaction get err = %v, want nil", err)
	}
}

func TestFakeTransactionImpl_CommitTransaction(t1 *testing.T) {
	t := &FakeTransactionImpl{}
	if err := t.CommitTransaction(); err != nil {
		t1.Errorf("CommitTransaction get err = %v, want nil", err)
	}
}

func TestFakeTransactionImpl_GetTransaction(t1 *testing.T) {
	t := &FakeTransactionImpl{}
	if got := t.GetTransaction(); got != nil {
		t1.Errorf("GetTransaction got = %v, want nil", got)
	}
}

func TestFakeTransactionImpl_RollBackTransaction(t1 *testing.T) {
	t := &FakeTransactionImpl{}
	t.RollBackTransaction()
}

func Test_doNothingTxOrm_Commit(t *testing.T) {
	d := doNothingTxOrm{}
	if err := d.Commit(); err != nil {
		t.Errorf("GetTransaction get err = %v, want nil", err)
	}
}

func Test_doNothingTxOrm_DBStats(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.DBStats(); got != nil {
		t.Errorf("DBStats() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_Delete(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.Delete(nil)
	if err != nil {
		t.Errorf("Delete() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("Delete() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_DeleteWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.DeleteWithCtx(nil, nil)
	if err != nil {
		t.Errorf("DeleteWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("DeleteWithCtx() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_Driver(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.Driver(); got != nil {
		t.Errorf("Driver() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_Insert(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.Insert(nil)
	if err != nil {
		t.Errorf("Insert() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("Insert() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_InsertMulti(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.InsertMulti(0, nil)
	if err != nil {
		t.Errorf("InsertMulti() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("InsertMulti() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_InsertMultiWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.InsertMultiWithCtx(nil, 0, nil)
	if err != nil {
		t.Errorf("InsertMultiWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("InsertMultiWithCtx() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_InsertOrUpdate(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.InsertOrUpdate(nil)
	if err != nil {
		t.Errorf("InsertOrUpdate() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("InsertOrUpdate() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_InsertOrUpdateWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.InsertOrUpdateWithCtx(nil, nil)
	if err != nil {
		t.Errorf("InsertOrUpdateWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("InsertOrUpdateWithCtx() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_InsertWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.InsertWithCtx(nil, nil)
	if err != nil {
		t.Errorf("InsertWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("InsertWithCtx() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_LoadRelated(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.LoadRelated(nil, "")
	if err != nil {
		t.Errorf("LoadRelated() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("LoadRelated() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_LoadRelatedWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.LoadRelatedWithCtx(nil, nil, "")
	if err != nil {
		t.Errorf("LoadRelatedWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("LoadRelatedWithCtx() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_QueryM2M(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.QueryM2M(nil, ""); got != nil {
		t.Errorf("QueryM2M() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_QueryM2MWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.QueryM2MWithCtx(nil, nil, ""); got != nil {
		t.Errorf("QueryM2MWithCtx() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_QueryTable(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.QueryTable(nil); got != nil {
		t.Errorf("QueryTable() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_QueryTableWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.QueryTableWithCtx(nil, nil); got != nil {
		t.Errorf("QueryTableWithCtx() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_Raw(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.Raw(""); got != nil {
		t.Errorf("Raw() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_RawWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.RawWithCtx(nil, ""); got != nil {
		t.Errorf("RawWithCtx() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_Read(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.Read(""); got != nil {
		t.Errorf("Read() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_ReadForUpdate(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.ReadForUpdate(""); got != nil {
		t.Errorf("ReadForUpdate() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_ReadForUpdateWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	if got := d.ReadForUpdateWithCtx(nil, nil); got != nil {
		t.Errorf("ReadForUpdateWithCtx() = %v, want nil", got)
	}
}

func Test_doNothingTxOrm_ReadOrCreate(t *testing.T) {
	d := doNothingTxOrm{}
	got, num, err := d.ReadOrCreate(nil, "")
	if err != nil {
		t.Errorf("ReadOrCreate() err = %v, want nil", err)
	}
	if got != false {
		t.Errorf("ReadOrCreate() = %v, want false", got)
	}
	if num != 0 {
		t.Errorf("ReadOrCreate() err = %v, want 0", num)
	}
}

func Test_doNothingTxOrm_ReadOrCreateWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, num, err := d.ReadOrCreateWithCtx(nil, nil, "")
	if err != nil {
		t.Errorf("ReadOrCreateWithCtx() err = %v, want nil", err)
	}
	if got != false {
		t.Errorf("ReadOrCreateWithCtx() = %v, want false", got)
	}
	if num != 0 {
		t.Errorf("ReadOrCreateWithCtx() err = %v, want 0", num)
	}
}

func Test_doNothingTxOrm_ReadWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	if err := d.ReadWithCtx(nil, nil); err != nil {
		t.Errorf("ReadWithCtx get err = %v, want nil", err)
	}
}

func Test_doNothingTxOrm_Rollback(t *testing.T) {
	d := doNothingTxOrm{}
	if err := d.Rollback(); err != nil {
		t.Errorf("Rollback get err = %v, want nil", err)
	}
}

func Test_doNothingTxOrm_RollbackUnlessCommit(t *testing.T) {
	d := doNothingTxOrm{}
	if err := d.RollbackUnlessCommit(); err != nil {
		t.Errorf("RollbackUnlessCommit get err = %v, want nil", err)
	}
}

func Test_doNothingTxOrm_Update(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.Update(nil)
	if err != nil {
		t.Errorf("Update() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("Update() got = %v, want 0", got)
	}
}

func Test_doNothingTxOrm_UpdateWithCtx(t *testing.T) {
	d := doNothingTxOrm{}
	got, err := d.UpdateWithCtx(nil, nil)
	if err != nil {
		t.Errorf("UpdateWithCtx() error = %v, wantErr nil", err)
		return
	}
	if got != 0 {
		t.Errorf("UpdateWithCtx() got = %v, want 0", got)
	}
}
