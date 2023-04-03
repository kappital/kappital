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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
)

func TestCommonUtilRequest(t *testing.T) {
	convey.Convey("Test CommonUtilRequest_old", t, func() {
		convey.Convey("case 1: json.Marshal has error", func() {
			p := gomonkey.ApplyFunc(json.Marshal, func(_ interface{}) ([]byte, error) {
				return nil, fmt.Errorf("mock error")
			})
			defer p.Reset()
			code, buf, err := CommonUtilRequest(&RequestInfo{Body: "x"})
			convey.So(code, convey.ShouldEqual, 0)
			convey.So(buf, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: http.NewRequest has error", func() {
			p := gomonkey.ApplyFunc(http.NewRequestWithContext, func(_ context.Context, _, _ string, _ io.Reader) (*http.Request, error) {
				return nil, fmt.Errorf("mock error")
			})
			defer p.Reset()
			code, buf, err := CommonUtilRequest(&RequestInfo{Method: "xxx", Body: "x"})
			convey.So(code, convey.ShouldEqual, 0)
			convey.So(buf, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: mock the client.Do", func() {
			code, buf, err := CommonUtilRequest(&RequestInfo{Method: "xxx", HeaderAdder: map[string]string{"a": "b"}, HeaderSetter: map[string]string{"b": "c"}})
			convey.So(code, convey.ShouldEqual, 0)
			convey.So(buf, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 4:", func() {
			var ts *httptest.Server
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("output"))
				if err != nil {
					t.Errorf("unable to write output into server, err: %v", err)
					return
				}
			}))
			defer ts.Close()
			code, buf, err := CommonUtilRequest(&RequestInfo{Method: "GET", Path: ts.URL, HeaderAdder: map[string]string{"a": "b"}, HeaderSetter: map[string]string{"b": "c"}})
			convey.So(code, convey.ShouldEqual, http.StatusOK)
			convey.So(buf, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func TestRequestInfo_getClient(t *testing.T) {
	convey.Convey("Test RequestInfo getClient", t, func() {
		r := RequestInfo{Path: "https://x.x.x.x", CaCrt: "xxx"}
		p := gomonkey.ApplyMethodSeq(reflect.TypeOf(base64.StdEncoding), "DecodeString", []gomonkey.OutputCell{
			// first call get error
			{Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			// second call get error
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			// third call get error
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			// three time called does not get error but X509KeyPair gets error
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}},
			// three time called does not get error and X509KeyPair does not get error
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(tls.X509KeyPair, []gomonkey.OutputCell{
			{Values: gomonkey.Params{tls.Certificate{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{tls.Certificate{}, nil}},
		})

		client, err := r.getClient()
		convey.So(client, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)
		client, err = r.getClient()
		convey.So(client, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)
		client, err = r.getClient()
		convey.So(client, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)
		client, err = r.getClient()
		convey.So(client, convey.ShouldBeNil)
		convey.So(err, convey.ShouldNotBeNil)
		client, err = r.getClient()
		convey.So(client, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGetLocalIP(t *testing.T) {
	convey.Convey("Test GetLocalIP", t, func() {
		p := gomonkey.ApplyFuncSeq(net.InterfaceAddrs, []gomonkey.OutputCell{
			{Values: gomonkey.Params{nil, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]net.Addr{&net.IPNet{IP: []byte("xxxx")}}, nil}},
			{Values: gomonkey.Params{[]net.Addr{}, nil}},
		})
		defer p.Reset()

		ip, err := GetLocalIP()
		convey.So(ip, convey.ShouldBeEmpty)
		convey.So(err, convey.ShouldNotBeNil)

		ip, err = GetLocalIP()
		convey.So(ip, convey.ShouldNotBeEmpty)
		convey.So(err, convey.ShouldBeNil)

		ip, err = GetLocalIP()
		convey.So(ip, convey.ShouldBeEmpty)
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func Test_getHTTPRequestWithBody(t *testing.T) {
	type args struct {
		method string
		path   string
		body   interface{}
		isFile bool
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name:    "Test getHTTPRequestWithBody (is file but not accept)",
			args:    args{isFile: true},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "Test getHTTPRequestWithBody (is file)",
			args: args{method: "x", path: "x", isFile: true, body: &bytes.Buffer{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHTTPRequestWithBody(tt.args.method, tt.args.path, tt.args.body, tt.args.isFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHTTPRequestWithBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantNil {
				t.Errorf("getHTTPRequestWithBody() got = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func TestReplaceIP(t *testing.T) {
	type args struct {
		org         string
		newIP       string
		defaultPort string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test ReplaceIP (has port)",
			args: args{org: "x.x.x.x:x", newIP: "y.y.y.y"},
			want: "y.y.y.y:x",
		},
		{
			name: "Test ReplaceIP (does not have port)",
			args: args{newIP: "x.x.x.x", defaultPort: "x"},
			want: "x.x.x.x:x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceIP(tt.args.org, tt.args.newIP, tt.args.defaultPort); got != tt.want {
				t.Errorf("ReplaceIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
