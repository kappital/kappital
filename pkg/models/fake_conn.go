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
	"context"
	"database/sql"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/utils"
)

// FakeTransactionImpl for unit test
type FakeTransactionImpl struct {
	transaction orm.TxOrmer
}

// BeginTransaction fake method for unit test
func (t *FakeTransactionImpl) BeginTransaction() error {
	t.transaction = &doNothingTxOrm{}
	return nil
}

// CommitTransaction fake method for unit test
func (t *FakeTransactionImpl) CommitTransaction() error {
	return nil
}

// RollBackTransaction fake method for unit test
func (t *FakeTransactionImpl) RollBackTransaction() {}

// GetTransaction from implement
func (t *FakeTransactionImpl) GetTransaction() orm.TxOrmer {
	return t.transaction
}

type doNothingTxOrm struct{}

// Read fake method
func (d doNothingTxOrm) Read(_ interface{}, _ ...string) error {
	return nil
}

// ReadWithCtx fake method
func (d doNothingTxOrm) ReadWithCtx(_ context.Context, _ interface{}, _ ...string) error {
	return nil
}

// ReadForUpdate fake method
func (d doNothingTxOrm) ReadForUpdate(_ interface{}, _ ...string) error {
	return nil
}

// ReadForUpdateWithCtx fake method
func (d doNothingTxOrm) ReadForUpdateWithCtx(_ context.Context, _ interface{}, _ ...string) error {
	return nil
}

// ReadOrCreate fake method
func (d doNothingTxOrm) ReadOrCreate(_ interface{}, _ string, _ ...string) (bool, int64, error) {
	return false, 0, nil
}

// ReadOrCreateWithCtx fake method
func (d doNothingTxOrm) ReadOrCreateWithCtx(_ context.Context, _ interface{}, _ string, _ ...string) (bool, int64, error) {
	return false, 0, nil
}

// LoadRelated fake method
func (d doNothingTxOrm) LoadRelated(_ interface{}, _ string, _ ...utils.KV) (int64, error) {
	return 0, nil
}

// LoadRelatedWithCtx fake method
func (d doNothingTxOrm) LoadRelatedWithCtx(_ context.Context, _ interface{}, _ string, _ ...utils.KV) (int64, error) {
	return 0, nil
}

// QueryM2M fake method
func (d doNothingTxOrm) QueryM2M(_ interface{}, _ string) orm.QueryM2Mer {
	return nil
}

// QueryM2MWithCtx fake method
func (d doNothingTxOrm) QueryM2MWithCtx(_ context.Context, _ interface{}, _ string) orm.QueryM2Mer {
	return nil
}

// QueryTable fake method
func (d doNothingTxOrm) QueryTable(_ interface{}) orm.QuerySeter {
	return nil
}

// QueryTableWithCtx fake method
func (d doNothingTxOrm) QueryTableWithCtx(_ context.Context, _ interface{}) orm.QuerySeter {
	return nil
}

// DBStats fake method
func (d doNothingTxOrm) DBStats() *sql.DBStats {
	return nil
}

// Insert fake method
func (d doNothingTxOrm) Insert(_ interface{}) (int64, error) {
	return 0, nil
}

// InsertWithCtx fake method
func (d doNothingTxOrm) InsertWithCtx(_ context.Context, _ interface{}) (int64, error) {
	return 0, nil
}

// InsertOrUpdate fake method
func (d doNothingTxOrm) InsertOrUpdate(_ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// InsertOrUpdateWithCtx fake method
func (d doNothingTxOrm) InsertOrUpdateWithCtx(_ context.Context, _ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// InsertMulti fake method
func (d doNothingTxOrm) InsertMulti(_ int, _ interface{}) (int64, error) {
	return 0, nil
}

// InsertMultiWithCtx fake method
func (d doNothingTxOrm) InsertMultiWithCtx(_ context.Context, _ int, _ interface{}) (int64, error) {
	return 0, nil
}

// Update fake method
func (d doNothingTxOrm) Update(_ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// UpdateWithCtx fake method
func (d doNothingTxOrm) UpdateWithCtx(_ context.Context, _ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// Delete fake method
func (d doNothingTxOrm) Delete(_ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// DeleteWithCtx fake method
func (d doNothingTxOrm) DeleteWithCtx(_ context.Context, _ interface{}, _ ...string) (int64, error) {
	return 0, nil
}

// Raw fake method
func (d doNothingTxOrm) Raw(_ string, _ ...interface{}) orm.RawSeter {
	return nil
}

// RawWithCtx fake method
func (d doNothingTxOrm) RawWithCtx(_ context.Context, _ string, _ ...interface{}) orm.RawSeter {
	return nil
}

// Driver fake method
func (d doNothingTxOrm) Driver() orm.Driver {
	return nil
}

// Commit fake method
func (d doNothingTxOrm) Commit() error {
	return nil
}

// Rollback fake method
func (d doNothingTxOrm) Rollback() error {
	return nil
}

// RollbackUnlessCommit fake method
func (d doNothingTxOrm) RollbackUnlessCommit() error {
	return nil
}
