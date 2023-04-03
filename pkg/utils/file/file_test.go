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

package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
)

func TestIsFileExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test IsFileExist (not exist)", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileExist(tt.args.path); got != tt.want {
				t.Errorf("IsFileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDirExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test IsDirExist (not exist)", want: false},
		{name: "Test IsDirExist (exist)", args: args{path: "."}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDirExist(tt.args.path); got != tt.want {
				t.Errorf("IsDirExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCompressedFile(t *testing.T) {
	type args struct {
		path      string
		fileExist bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test IsCompressedFile (tgzPostfix)", args: args{path: "xxx.tgz", fileExist: true}, want: true},
		{name: "Test IsCompressedFile (tarGzipSuffix)", args: args{path: "xxx.tar.gz", fileExist: true}, want: true},
		{name: "Test IsCompressedFile (zipSuffix)", args: args{path: "xxx.zip", fileExist: true}, want: true},
		{name: "Test IsCompressedFile (illegal)", args: args{path: "xxx", fileExist: true}},
		{name: "Test IsCompressedFile (file is not exist)", args: args{path: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := gomonkey.ApplyFunc(IsFileExist, func(_ string) bool { return tt.args.fileExist })
			defer p.Reset()
			if got := IsCompressedFile(tt.args.path); got != tt.want {
				t.Errorf("IsCompressedFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFileToBase64(t *testing.T) {
	convey.Convey("Test ReadFileToBase64", t, func() {
		p := gomonkey.ApplyFuncSeq(filepath.Abs, []gomonkey.OutputCell{
			{Values: gomonkey.Params{"", fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{"", nil}},
			{Values: gomonkey.Params{"", nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
			{Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]byte{}, nil}},
		})
		res := ReadFileToBase64("")
		convey.So(res, convey.ShouldBeEmpty)
		res = ReadFileToBase64("")
		convey.So(res, convey.ShouldBeEmpty)
		res = ReadFileToBase64("")
		convey.So(res, convey.ShouldBeEmpty)
	})
}
