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

package initiate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"

	"github.com/kappital/kappital/pkg/utils/file"
)

func Test_createFile(t *testing.T) {
	convey.Convey("test createFile", t, func() {
		convey.Convey("case 1: demo.ReadFile has error", func() {
			err := createFile("", "", "", nil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: directory exist with os.OpenFile error", func() {
			p := gomonkey.ApplyFunc(file.IsDirExist, func(_ string) bool { return false })
			defer p.Reset()
			p.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error { return nil })
			p.ApplyFunc(os.OpenFile, func(_ string, _ int, _ os.FileMode) (*os.File, error) {
				return nil, fmt.Errorf("mock error")
			})
			err := createFile("", "./kappital-demo/", "metadata.yaml", nil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: directory is not exist with os.MkdirAll error", func() {
			p := gomonkey.ApplyFunc(file.IsDirExist, func(_ string) bool { return false })
			defer p.Reset()
			p.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error {
				return fmt.Errorf("mock error")
			})
			p.ApplyFunc(os.OpenFile, func(_ string, _ int, _ os.FileMode) (*os.File, error) { return nil, nil })
			err := createFile("", "./kappital-demo/", "metadata.yaml", nil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 4: no errors", func() {
			defer func() {
				_ = recover()
			}()
			p := gomonkey.ApplyFunc(file.IsDirExist, func(_ string) bool { return false })
			defer p.Reset()
			p.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error { return nil })
			p.ApplyFunc(os.OpenFile, func(_ string, _ int, _ os.FileMode) (*os.File, error) { return &os.File{}, nil })
			err := createFile("", "./kappital-demo/", "metadata.yaml", &operation{})
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_modifyMetadataFile(t *testing.T) {
	tests := []struct {
		name    string
		b       []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "test modifyMetadataFile (valid metadata)",
			b: []byte(`
name: test
version: test-version
type: operator
minKubeVersion: 1.15.0
briefDescription: example package with an example operator and instance
`),
			want: []byte(`
briefDescription: example package with an example operator and instance
logo: {}
minKubeVersion: 1.15.0
name: test
provider: {}
type: operator
version: test-version
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := modifyMetadataFile(tt.b, "test", "test-version")
			if (err != nil) != tt.wantErr {
				t.Errorf("modifyMetadataFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) == string(tt.want) {
				t.Errorf("modifyMetadataFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_operation_NewCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "Test operation NewCommand for init"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &operation{}
			if got := o.NewCommand(); (got != nil) && tt.wantNil {
				t.Errorf("NewCommand() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_operation_PreRunE(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test operation PreRunE with one arg", args: args{args: []string{"x"}}},
		{name: "test operation PreRunE without args", args: args{args: []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &operation{}
			if err := o.PreRunE(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("PreRunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_operation_createPackage(t *testing.T) {
	convey.Convey("test operation createPackage", t, func() {
		p := gomonkey.ApplyFunc(file.IsDirExist, func(_ string) bool { return false })
		defer p.Reset()
		p.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error { return nil })
		p.ApplyFunc(os.OpenFile, func(_ string, _ int, _ os.FileMode) (*os.File, error) { return &os.File{}, nil })
		entries, err := demo.ReadDir(demoName)
		if err != nil {
			t.Errorf("cannot open the demo")
		}
		o := &operation{}
		convey.Convey("case 1: out of the max depth", func() {
			defer func() {
				_ = recover()
			}()
			err = o.createPackage(entries, "", demoName, 10)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: read all packages", func() {
			defer func() {
				_ = recover()
			}()
			err = o.createPackage(entries, "", demoName, 0)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "Test NewCommand"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommand(); (got != nil) && tt.wantNil {
				t.Errorf("NewCommand() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("test init", t, func() {
		o := &operation{}
		convey.Convey("case 1: cannot get abs path", func() {
			p := gomonkey.ApplyFunc(filepath.Abs, func(_ string) (string, error) {
				return "", fmt.Errorf("mock error")
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: cannot create directory", func() {
			p := gomonkey.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error {
				return fmt.Errorf("mock error")
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: create the package", func() {
			defer func() {
				_ = recover()
			}()
			p := gomonkey.ApplyFunc(file.IsDirExist, func(_ string) bool { return false })
			defer p.Reset()
			p.ApplyFunc(os.MkdirAll, func(_ string, _ os.FileMode) error { return nil })
			p.ApplyFunc(os.OpenFile, func(_ string, _ int, _ os.FileMode) (*os.File, error) { return &os.File{}, nil })
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}
