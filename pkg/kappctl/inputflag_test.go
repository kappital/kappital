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

package kappctl

import (
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func Test_inputFlag_AddBoolFlag(t *testing.T) {
	var res bool
	tmp := &cobra.Command{}
	i := inputFlag{name: "test name", shortHand: "x", defaultValue: true, usage: "test"}
	i.AddBoolFlag(&res, tmp)
	if res != true {
		t.Errorf("AddBoolFlag() = %v, want true", res)
	}

	res = false
	j := inputFlag{name: "test name", shortHand: "x", defaultValue: "true", usage: "test"}
	j.AddBoolFlag(&res, tmp)
	if res != false {
		t.Errorf("AddBoolFlag() = %v, want false", res)
	}
}

func Test_inputFlag_AddStringFlag(t *testing.T) {
	var res string
	tmp := &cobra.Command{}
	i := inputFlag{name: "test name", shortHand: "x", defaultValue: "false", usage: "test"}
	i.AddStringFlag(&res, tmp)
	if !reflect.DeepEqual(res, "false") {
		t.Errorf("AddStringFlag() = %v, want false", res)
	}

	res = ""
	j := inputFlag{name: "test name", shortHand: "x", defaultValue: false, usage: "test"}
	j.AddStringFlag(&res, tmp)
	if !reflect.DeepEqual(res, "") {
		t.Errorf("AddStringFlag() = %v, want ''", res)
	}
}

func Test_inputFlag_GetFlagName(t *testing.T) {
	type fields struct {
		name         string
		shortHand    string
		defaultValue interface{}
		usage        string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Test inputFlag GetFlagName",
			fields: fields{name: "test name", shortHand: "x", defaultValue: false, usage: "test"},
			want:   "test name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inputFlag{
				name:         tt.fields.name,
				shortHand:    tt.fields.shortHand,
				defaultValue: tt.fields.defaultValue,
				usage:        tt.fields.usage,
			}
			if got := i.GetFlagName(); got != tt.want {
				t.Errorf("GetFlagName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inputFlag_MarkFlagRequired(t *testing.T) {
	i := inputFlag{name: "test name", shortHand: "x", defaultValue: false, usage: "test"}
	i.MarkFlagRequired(&cobra.Command{})
}

func Test_newInputFlag(t *testing.T) {
	type args struct {
		name         string
		shortHand    string
		defaultValue interface{}
		usage        string
	}
	tests := []struct {
		name string
		args args
		want inputFlag
	}{
		{
			name: "Test newInputFlag",
			args: args{name: "test name", shortHand: "x", defaultValue: false, usage: "test"},
			want: inputFlag{name: "test name", shortHand: "x", defaultValue: false, usage: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newInputFlag(tt.args.name, tt.args.shortHand, tt.args.defaultValue, tt.args.usage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInputFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
