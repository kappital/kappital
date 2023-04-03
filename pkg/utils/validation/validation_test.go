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

package validation

import (
	"testing"
)

func TestValidBool(t *testing.T) {
	type args struct {
		passIn string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "Test ValidBool (1)", args: args{passIn: "1"}, want: true},
		{name: "Test ValidBool (t)", args: args{passIn: "t"}, want: true},
		{name: "Test ValidBool (T)", args: args{passIn: "T"}, want: true},
		{name: "Test ValidBool (true)", args: args{passIn: "true"}, want: true},
		{name: "Test ValidBool (TRUE)", args: args{passIn: "TRUE"}, want: true},
		{name: "Test ValidBool (True)", args: args{passIn: "True"}, want: true},
		{name: "Test ValidBool (0)", args: args{passIn: "0"}},
		{name: "Test ValidBool (f)", args: args{passIn: "f"}},
		{name: "Test ValidBool (F)", args: args{passIn: "F"}},
		{name: "Test ValidBool (false)", args: args{passIn: "false"}},
		{name: "Test ValidBool (FALSE)", args: args{passIn: "FALSE"}},
		{name: "Test ValidBool (false)", args: args{passIn: "False"}},
		{name: "Test ValidBool (has error)", args: args{passIn: "xxxx"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidBool(tt.args.passIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidCommonString(t *testing.T) {
	type args struct {
		passIn       string
		min          int
		max          int
		canEmpty     bool
		defaultValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test ValidCommonString (passIn is valid)", args: args{passIn: "test", min: 0, max: 10}, want: "test"},
		{name: "Test ValidCommonString (passIn is invalid, but can be empty)", args: args{passIn: "t", min: 10, max: 10, canEmpty: true}},
		{name: "Test ValidCommonString (passIn is invalid, but default values is valid)", args: args{passIn: "t", min: 3, max: 10, defaultValue: "test"}, want: "test"},
		{name: "Test ValidCommonString (all in valid)", args: args{passIn: "test", min: 10, max: 20}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidCommonString(tt.args.passIn, tt.args.min, tt.args.max, tt.args.canEmpty, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidCommonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidCommonString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkString(t *testing.T) {
	type args struct {
		str string
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test checkString (not valid)", args: args{str: "test", min: 100, max: 200}},
		{name: "Test checkString (valid case 1)", args: args{str: "namespace", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 2)", args: args{str: "name_space", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 3)", args: args{str: "namespace12", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 4)", args: args{str: "NAMESPACE", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 5)", args: args{str: "NAME_SPACE", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 6)", args: args{str: "Name_space_123", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 7)", args: args{str: "Name123_space_123", min: 0, max: 100}, want: true},
		{name: "Test checkString (valid case 8)", args: args{str: "test-123", min: 0, max: 100}, want: true},
		{name: "Test checkString (invalid case 1)", args: args{str: "_", min: 0, max: 100}},
		{name: "Test checkString (invalid case 2)", args: args{str: "_123", min: 0, max: 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkString(tt.args.str, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("checkString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validLength(t *testing.T) {
	type args struct {
		str string
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test validateLength (valid)", args: args{str: "test", min: 0, max: 100}, want: true},
		{name: "Test validateLength (not valid)", args: args{str: "test", min: 100, max: 200}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLength(tt.args.str, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("validateLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
