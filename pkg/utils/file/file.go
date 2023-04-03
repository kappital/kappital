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
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	tgzPostfix    = ".tgz"
	tarGzipSuffix = ".tar.gz"
	zipSuffix     = ".zip"
)

// IsFileExist does the file is existed in this path
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsCompressedFile does the file is exist and is valid compressed file
func IsCompressedFile(path string) bool {
	if !IsFileExist(path) {
		return false
	}
	if strings.HasSuffix(path, tgzPostfix) ||
		strings.HasSuffix(path, tarGzipSuffix) ||
		strings.HasSuffix(path, zipSuffix) {
		return true
	}
	return false
}

// IsDirExist does the directory is existed in this path
func IsDirExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// ReadFileToBase64 read a file and trans the context to base64
func ReadFileToBase64(path string) string {
	return base64.StdEncoding.EncodeToString(ReadFileToBytes(path))
}

// ReadFileToBytes read a file and trans the context to bytes
func ReadFileToBytes(path string) []byte {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	bytes, err := ioutil.ReadFile(filepath.Clean(absPath))
	if err != nil {
		return nil
	}
	return bytes
}
