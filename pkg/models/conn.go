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
	"os"

	"github.com/beego/beego/v2/client/orm"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/constants"
)

var dbAliasName = "default"

// Handler functions for error and transaction
var Handler = func(err *error, tx Transaction) {
	if IgnoreDBInsertIDError(*err) != nil {
		tx.RollBackTransaction()
		return
	}
	*err = tx.CommitTransaction()
}

// Transaction of the database actions
type Transaction interface {
	BeginTransaction() error
	CommitTransaction() error
	RollBackTransaction()

	GetTransaction() orm.TxOrmer
}

// TransactionImpl which implement the sqlSession
type TransactionImpl struct {
	sqlSession  orm.Ormer
	transaction orm.TxOrmer
}

// GetNewOrm get the orm operator
func GetNewOrm() orm.Ormer {
	return orm.NewOrmUsingDB(dbAliasName)
}

// NewTransaction create a new transaction
func NewTransaction(sqlSession orm.Ormer) Transaction {
	if os.Getenv(constants.UsingFakeEnv) == "true" {
		return &FakeTransactionImpl{}
	}
	return &TransactionImpl{sqlSession: sqlSession}
}

// BeginTransaction begins using the transaction
func (t *TransactionImpl) BeginTransaction() error {
	var err error
	t.transaction, err = t.sqlSession.Begin()
	return err
}

// CommitTransaction commit the transaction if it does not have error
func (t *TransactionImpl) CommitTransaction() error {
	return t.transaction.Commit()
}

// RollBackTransaction not return err for meaningless golint and codex
func (t *TransactionImpl) RollBackTransaction() {
	if err := t.transaction.Rollback(); err != nil {
		klog.Errorf("failed to rollback transaction, error: %v", err)
	}
}

// GetTransaction from implement
func (t *TransactionImpl) GetTransaction() orm.TxOrmer {
	return t.transaction
}
