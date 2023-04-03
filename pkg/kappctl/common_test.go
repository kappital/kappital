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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/brahma-adshonor/gohook"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetAgeOutput(t *testing.T) {
	err := gohook.Hook(time.Since, func(_ time.Time) time.Duration {
		return 0
	}, nil)
	if err != nil {
		t.Errorf("hook err, :%v", err)
		return
	}
	defer gohook.UnHook(time.Since) //nolint:errcheck
	type args struct {
		timestamp time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GetAgeOutput",
			args: args{timestamp: time.Now().UTC()},
			want: "0s ago",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAgeOutput(tt.args.timestamp); got != tt.want {
				t.Errorf("GetAgeOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	convey.Convey("test GetConfig", t, func() {
		convey.Convey("case 1: cannot get config file", func() {
			_, err := GetConfig()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("case 2: cannot ReadFile", func() {
			p := gomonkey.ApplyFunc(os.Getenv, func(_ string) string {
				return "monkey"
			})
			defer p.Reset()
			_, err := GetConfig()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("case 3: cannot Unmarshal", func() {
			p1 := gomonkey.ApplyFunc(os.Getenv, func(_ string) string {
				return "monkey"
			})
			defer p1.Reset()
			p2 := gomonkey.ApplyFunc(os.ReadFile, func(_ string) ([]byte, error) {
				return nil, nil
			})
			defer p2.Reset()
			p3 := gomonkey.ApplyFunc(json.Unmarshal, func(_ []byte, _ interface{}) error {
				return fmt.Errorf("monkey")
			})
			defer p3.Reset()
			_, err := GetConfig()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("case 4: no errors", func() {
			p1 := gomonkey.ApplyFunc(os.Getenv, func(_ string) string {
				return "monkey"
			})
			defer p1.Reset()
			p2 := gomonkey.ApplyFunc(os.ReadFile, func(_ string) ([]byte, error) {
				return []byte("{}"), nil
			})
			defer p2.Reset()
			_, err := GetConfig()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestOutputYAMLOrJSONString(t *testing.T) {
	type args struct {
		buf    []byte
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test OutputYAMLOrJSONString (no format)"},
		{name: "Test OutputYAMLOrJSONString (format is invalid)", args: args{format: "xxx"}, wantErr: true},
		{
			name: "Test OutputYAMLOrJSONString (json)",
			args: args{buf: []byte("{\"key\":\"value\"}"), format: jsonOutputFormat},
		},
		{
			name: "Test OutputYAMLOrJSONString (yaml)",
			args: args{buf: []byte("{\"key\":\"value\"}"), format: yamlOutputFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OutputYAMLOrJSONString(tt.args.buf, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("OutputYAMLOrJSONString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test validFormat (json)", args: args{format: "json"}, wantErr: false},
		{name: "Test validFormat (JSON)", args: args{format: "JSON"}, wantErr: false},
		{name: "Test validFormat (yaml)", args: args{format: "yaml"}, wantErr: false},
		{name: "Test validFormat (YAML)", args: args{format: "YAML"}, wantErr: false},
		{name: "Test validFormat (empty)", args: args{format: ""}, wantErr: false},
		{name: "Test validFormat (invalid)", args: args{format: "false"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("validFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getKappitalConfigPath(t *testing.T) {
	convey.Convey("test getKappitalConfigPath", t, func() {
		convey.Convey("case 1: cannot use UserHomeDir", func() {
			p := gomonkey.ApplyFunc(os.UserHomeDir, func() (string, error) {
				return "", fmt.Errorf("")
			})
			defer p.Reset()
			_, err := getKappitalConfigPath()
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("case 2: the home dir has config file", func() {
			p := gomonkey.ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
				return nil, nil
			})
			defer p.Reset()
			got, err := getKappitalConfigPath()
			convey.So(err, convey.ShouldBeNil)
			convey.So(got, convey.ShouldNotBeEmpty)
		})

		convey.Convey("case 3: can get the config file from env", func() {
			p := gomonkey.ApplyFunc(os.Getenv, func(_ string) string {
				return "monkey"
			})
			defer p.Reset()
			got, err := getKappitalConfigPath()
			convey.So(err, convey.ShouldBeNil)
			convey.So(got, convey.ShouldNotBeEmpty)
		})

		convey.Convey("case 4: no mock for this method", func() {
			got, err := getKappitalConfigPath()
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(got, convey.ShouldBeEmpty)
		})
	})
}

func TestIsInputValidate(t *testing.T) {
	type args struct {
		inputs map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test IsInputValidate (with maxStringLength error)",
			args: args{map[string]interface{}{
				"0": false,
				"1": "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			}},
			wantErr: true,
		},
		{
			name:    "Test IsInputValidate (with invalid cluster)",
			args:    args{map[string]interface{}{"cluster": "123"}},
			wantErr: true,
		},
		{
			name: "Test IsInputValidate",
			args: args{map[string]interface{}{"1": "123"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsInputValidate(tt.args.inputs); (err != nil) != tt.wantErr {
				t.Errorf("IsInputValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValidClusterName(t *testing.T) {
	type args struct {
		clusterName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test isValidClusterName (invalid first block for number)", args: args{clusterName: "1234"}, wantErr: true},
		{name: "Test isValidClusterName (invalid first block for special character)", args: args{clusterName: "a`b/c?d=xfe"}, wantErr: true},
		{name: "Test isValidClusterName (empty string after -)", args: args{clusterName: "abc-"}, wantErr: true},
		{name: "Test isValidClusterName (invalid after first block for special character)", args: args{clusterName: "abc-a`b/c?d=xfe"}, wantErr: true},
		{name: "Test isValidClusterName (valid clusterName case 1)", args: args{clusterName: "abc"}},
		{name: "Test isValidClusterName (valid clusterName case 2)", args: args{clusterName: "abc123"}},
		{name: "Test isValidClusterName (valid clusterName case 3)", args: args{clusterName: "abc-1"}},
		{name: "Test isValidClusterName (valid clusterName case 4)", args: args{clusterName: "abc-abc"}},
		{name: "Test isValidClusterName (valid clusterName case 5)", args: args{clusterName: "abc-abc1"}},
		{name: "Test isValidClusterName (valid clusterName case 6)", args: args{clusterName: "abc1-abc1"}},
		{name: "Test isValidClusterName (valid clusterName case 7)", args: args{clusterName: "abc-abc1abc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValidClusterName(tt.args.clusterName); (err != nil) != tt.wantErr {
				t.Errorf("IsInputValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
