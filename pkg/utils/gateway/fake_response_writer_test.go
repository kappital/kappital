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

package gateway

import (
	"net/http"
	"reflect"
	"testing"
)

func TestFakeResponseWriter_Header(t *testing.T) {
	tests := []struct {
		name string
		want http.Header
	}{
		{name: "Test FakeResponseWriter Header"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeResponseWriter{}
			if got := f.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeResponseWriter_Write(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{name: "Test FakeResponseWriter Write"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeResponseWriter{}
			got, err := f.Write(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFakeResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test FakeResponseWriter Write"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FakeResponseWriter{}
			f.WriteHeader(0)
		})
	}
}
