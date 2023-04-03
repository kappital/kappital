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
	"os"
	"time"
)

// FakeFileInfo for unit test
type FakeFileInfo struct{}

// Name for fake method
func (f FakeFileInfo) Name() string {
	return "fake"
}

// Size for fake method
func (f FakeFileInfo) Size() int64 {
	return 0
}

// Mode for fake method
func (f FakeFileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

// ModTime for fake method
func (f FakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

// IsDir for fake method
func (f FakeFileInfo) IsDir() bool {
	return len(os.Getenv("UNIT-TEST-FAKE-FILE")) > 0
}

// Sys for fake method
func (f FakeFileInfo) Sys() interface{} {
	return nil
}

// FakeFile for unit test
type FakeFile struct{}

// Read for fake method
func (f FakeFile) Read(_ []byte) (int, error) {
	return 0, nil
}

// ReadAt for fake method
func (f FakeFile) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, nil
}

// Seek for fake method
func (f FakeFile) Seek(_ int64, _ int) (int64, error) {
	return 0, nil
}

// Close for fake method
func (f FakeFile) Close() error {
	return nil
}
