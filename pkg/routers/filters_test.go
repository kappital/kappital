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

package routers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/beego/beego/v2/server/web/context"

	"github.com/kappital/kappital/pkg/constants"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

var fakeCtx = &context.Context{
	Input:  &context.BeegoInput{},
	Output: &context.BeegoOutput{},
	Request: &http.Request{
		Header: http.Header{"test": []string{}},
	},
	ResponseWriter: &context.Response{ResponseWriter: gateway.FakeResponseWriter{}},
}

func TestInitFilters(t *testing.T) {
	InitFilters()
}

func Test_beforeStaticFilter(t *testing.T) {
	fakeCtx.Request.Header.Set("Content-Type", "")
	beforeStaticFilter(fakeCtx)

	fakeCtx.Request.Header.Set("Content-Type", "multipart/form-data")
	fakeCtx.Request.ContentLength = 100 * 1024 * 1024
	beforeStaticFilter(fakeCtx)
}

func Test_formatFilter(t *testing.T) {
	fakeCtx.Reset(gateway.FakeResponseWriter{}, &http.Request{})

	fakeCtx.Input = &context.BeegoInput{
		Context: &context.Context{
			Request: &http.Request{
				URL:    &url.URL{Path: ""},
				Method: http.MethodConnect,
			},
		},
	}
	formatFilter(fakeCtx)

	fakeCtx.Input = &context.BeegoInput{
		Context: &context.Context{
			Request: &http.Request{
				URL:    &url.URL{Path: "/a/b/../"},
				Method: http.MethodGet,
			},
		},
	}
	formatFilter(fakeCtx)

	fakeCtx.Input = &context.BeegoInput{
		Context: &context.Context{
			Request: &http.Request{
				URL:    &url.URL{Path: "/a/b"},
				Method: http.MethodGet,
			},
		},
		RequestBody: make([]byte, constants.FileSize+1),
	}
	formatFilter(fakeCtx)
}
