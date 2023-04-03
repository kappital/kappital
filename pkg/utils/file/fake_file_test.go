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
	"io/fs"
	"reflect"
	"testing"
	"time"
)

func TestFakeFileInfo_IsDir(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "Test FakeFileInfo IsDir"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.IsDir(); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFileInfo_ModTime(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{name: "Test FakeFileInfo ModTime"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.ModTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFileInfo_Mode(t *testing.T) {
	tests := []struct {
		name string
		want fs.FileMode
	}{
		{name: "Test FakeFileInfo Mode", want: fs.ModePerm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.Mode(); got != tt.want {
				t.Errorf("Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFileInfo_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "Test FakeFileInfo Name", want: "fake"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFileInfo_Size(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		{name: "Test FakeFileInfo Size"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFileInfo_Sys(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{name: "Test FakeFileInfo Sys"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFileInfo{}
			if got := f.Sys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFile_Read(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "Test FakeFile Read"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFile{}
			got, err := f.Read(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFile_ReadAt(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "Test FakeFile ReadAt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFile{}
			got, err := f.ReadAt(nil, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadAt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFile_Seek(t *testing.T) {
	tests := []struct {
		name    string
		want    int64
		wantErr bool
	}{
		{name: "Test FakeFile Seek"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFile{}
			got, err := f.Seek(0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Seek() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeFile_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Test FakeFile Close"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeFile{}
			if err := f.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
