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
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/beego/beego/v2/client/orm"
	"github.com/smartystreets/goconvey/convey"
)

func Test_sqlite_GetDBConnection(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Test sqlite GetDBConnection"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sqlite{}
			if got := s.GetDBConnection(); got != tt.want {
				t.Errorf("GetDBConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sqlite_InitSQLDriver(t *testing.T) {
	convey.Convey("Test sqlite InitSQLDriver", t, func() {
		cfg := DefaultDatabaseConfiguration()
		s := sqlite{}
		p := gomonkey.ApplyFuncSeq(orm.RegisterDriver, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil}},
			{Values: gomonkey.Params{nil}},
			{Values: gomonkey.Params{nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(orm.RegisterDataBase, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil}},
			{Values: gomonkey.Params{nil}},
		})
		p.ApplyFunc(orm.GetDB, func(...string) (*sql.DB, error) { return &sql.DB{}, nil })
		p.ApplyFuncSeq(orm.RunSyncdb, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil}},
		})
		p.ApplyFuncSeq(os.Chmod, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
		})
		defer os.Remove("./Manager-sqlite.db")

		err := s.InitSQLDriver(cfg, Manager)
		convey.So(err, convey.ShouldNotBeNil)
		err = s.InitSQLDriver(cfg, Manager)
		convey.So(err, convey.ShouldNotBeNil)
		err = s.InitSQLDriver(cfg, Manager)
		convey.So(err, convey.ShouldNotBeNil)
		err = s.InitSQLDriver(cfg, serviceType("xx"))
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func Test_sqlite_registerModels(t *testing.T) {
	type args struct {
		serviceType serviceType
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test sqlite registerModels (Manager)", args: args{serviceType: Manager}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("cannot recover painc, r: %v", r)
				}
			}()
			s := sqlite{}
			s.registerModels(tt.args.serviceType)
		})
	}
}
