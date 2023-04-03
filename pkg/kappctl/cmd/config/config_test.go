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

package config

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/kappital/kappital/pkg/utils/file"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"reflect"
	"testing"

	"github.com/kappital/kappital/pkg/kappctl"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "Test config command NewCommand"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommand(); (got == nil) != tt.wantNil {
				t.Errorf("config command NewCommand() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_operation_constructNewConfig(t *testing.T) {
	type fields struct {
		managerIP        string
		managerHTTPSPort string
	}
	tests := []struct {
		name   string
		fields fields
		want   kappctl.Config
	}{
		{
			name:   "Test_operation_constructNewConfig (manager config only)",
			fields: fields{managerIP: "x.x.x.x", managerHTTPSPort: "xx"},
			want:   kappctl.Config{ManagerHTTPSServer: "https://x.x.x.x:xx"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operation{
				managerIP:        tt.fields.managerIP,
				managerHTTPSPort: tt.fields.managerHTTPSPort,
			}
			if got := o.constructNewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructNewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_operation_valid(t *testing.T) {
	type fields struct {
		managerIP        string
		managerHTTPSPort string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Test operation isValid (valid)",
			fields: fields{managerIP: "x.x.x.x", managerHTTPSPort: "12"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := operation{
				managerIP:        tt.fields.managerIP,
				managerHTTPSPort: tt.fields.managerHTTPSPort,
			}
			if got := o.isValid(); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	convey.Convey("Test run", t, func() {
		p := gomonkey.ApplyFuncSeq(os.UserHomeDir, []gomonkey.OutputCell{
			{Values: gomonkey.Params{"", fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{"", nil}},
			{Values: gomonkey.Params{"", nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(os.MkdirAll, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil}},
		})
		p.ApplyFuncSeq(file.IsFileExist, []gomonkey.OutputCell{
			{Values: gomonkey.Params{true}},
			{Values: gomonkey.Params{true}},
		})
		p.ApplyFuncSeq(os.RemoveAll, []gomonkey.OutputCell{
			{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{nil}},
		})

		err := run(operation{})
		convey.So(err, convey.ShouldNotBeNil)
		err = run(operation{})
		convey.So(err, convey.ShouldNotBeNil)
		err = run(operation{})
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func Test_validPort(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test validPort (empty)", args: args{port: ""}},
		{name: "Test validPort (-1)", args: args{port: "-1"}},
		{name: "Test validPort (65536)", args: args{port: "65536"}},
		{name: "Test validPort (65535)", args: args{port: "65535"}, want: true},
		{name: "Test validPort (6553xx5)", args: args{port: "6553xx5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPort(tt.args.port); got != tt.want {
				t.Errorf("validPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
