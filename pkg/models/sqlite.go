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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

const (
	sqliteRootPath   = "/opt/kappital/database"
	sqliteDriverName = "sqlite3"
)

type sqlite struct {
	name string
}

// InitSQLDriver for sqlite
func (s *sqlite) InitSQLDriver(cfg *DatabaseConfig, serviceType serviceType) error {
	s.initDBConfig(cfg, serviceType)
	if err := orm.RegisterDriver(sqliteDriverName, orm.DRSqlite); err != nil {
		return err
	}
	if err := orm.RegisterDataBase(AliasName, sqliteDriverName, s.name); err != nil {
		return err
	}
	s.registerModels(serviceType)
	dbInstance, _ := orm.GetDB()
	dbInstance.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	if err := orm.RunSyncdb(AliasName, false, false); err != nil {
		return err
	}
	if err := os.Chmod(s.name, os.FileMode(0600)); err != nil {
		return err
	}
	return nil
}

// GetDBConnection not using for the sqlite
func (s sqlite) GetDBConnection() string {
	return ""
}

func (s sqlite) registerModels(serviceType serviceType) {
	switch serviceType {
	case Manager:
		orm.RegisterModel(new(ServiceBindingModel), new(ResourceModel), new(InstanceModel))
	}
}

func (s *sqlite) initDBConfig(cfg *DatabaseConfig, serviceType serviceType) {
	s.name = filepath.Join(sqliteRootPath, fmt.Sprintf("%v-%v.db", serviceType, cfg.SQLDriver))
}
