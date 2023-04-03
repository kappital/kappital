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

package audit

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_defaultLog_initConfig(t *testing.T) {
	p := gomonkey.ApplyFuncSeq(os.OpenFile, []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
		{Values: gomonkey.Params{&os.File{}, nil}},
	})
	defer p.Reset()
	type fields struct {
		log *logrus.Logger
	}
	type args struct {
		config AuditLogConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test defaultLog initConfig (err)",
			fields:  fields{log: logrus.New()},
			args:    args{config: AuditLogConfig{AppName: "x", Filename: "xx"}},
			wantErr: true,
		},
		{
			name:   "Test defaultLog initConfig",
			fields: fields{log: logrus.New()},
			args:   args{config: AuditLogConfig{AppName: "x", Filename: "xx"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultLog{
				log: tt.fields.log,
			}
			if err := d.initConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("initConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
