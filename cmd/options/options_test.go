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

package options

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"k8s.io/client-go/util/cert"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/beego/beego/v2/server/web"
	"github.com/smartystreets/goconvey/convey"

	"github.com/kappital/kappital/pkg/utils/gateway"
	"github.com/kappital/kappital/pkg/utils/version"
)

func TestNewServerRunOptions(t *testing.T) {
	convey.Convey("Test NewServerRunOptions", t, func() {
		convey.Convey("case 1: mock error situation", func() {
			p := gomonkey.ApplyFuncSeq(version.SetupClusterVersion, []gomonkey.OutputCell{
				{Values: gomonkey.Params{fmt.Errorf("mock error")}},
			})
			defer p.Reset()
			p.ApplyFuncSeq(gateway.GetLocalIP, []gomonkey.OutputCell{
				{Values: gomonkey.Params{"", fmt.Errorf("mock error")}},
			})

			sro, err := NewServerRunOptions("")
			convey.So(sro, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
			sro, err = NewServerRunOptions(version.ServiceNameManager)
			convey.So(sro, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: does not mock error", func() {
			sro, err := NewServerRunOptions(version.ServiceNameManager)
			convey.So(sro, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)

			os.Args = []string{"x1", "x2"}
			sro, err = NewServerRunOptions(version.ServiceNameManager)
			convey.So(sro, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestServerRunOptions_generateConfig(t *testing.T) {
	tests := []struct {
		name    string
		Listen  web.Listen
		setEnvs map[string]string
		wantErr bool
	}{
		{name: "Test ServerRunOptions generateConfig (disable http and https)", wantErr: true},
		{
			name:    "Test ServerRunOptions generateConfig (enable https, missing cert file)",
			setEnvs: map[string]string{enableHTTPSEnvKey: "true"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerRunOptions{Config: web.Config{Listen: tt.Listen}}
			for k, v := range tt.setEnvs {
				t.Setenv(k, v)
			}
			if err := s.generateConfig(); (err != nil) != tt.wantErr {
				t.Errorf("generateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerRunOptions_getFlagSetValue(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test ServerRunOptions getFlagSetValue", args: args{prefix: "xx"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = []string{"x1", "x2"}
			fs := flag.NewFlagSet("xx", flag.ContinueOnError)
			var unusedInt int
			fs.IntVar(&unusedInt, "x3", 0, "")
			s := &ServerRunOptions{fs: fs}
			if err := s.getFlagSetValue(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("getFlagSetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test parseBool (1)", args: args{str: "1"}, want: true},
		{name: "Test parseBool (t)", args: args{str: "t"}, want: true},
		{name: "Test parseBool (T)", args: args{str: "T"}, want: true},
		{name: "Test parseBool (true)", args: args{str: "true"}, want: true},
		{name: "Test parseBool (TRUE)", args: args{str: "TRUE"}, want: true},
		{name: "Test parseBool (True)", args: args{str: "True"}, want: true},
		{name: "Test parseBool (0)", args: args{str: "0"}},
		{name: "Test parseBool (f)", args: args{str: "f"}},
		{name: "Test parseBool (F)", args: args{str: "F"}},
		{name: "Test parseBool (false)", args: args{str: "false"}},
		{name: "Test parseBool (FALSE)", args: args{str: "FALSE"}},
		{name: "Test parseBool (false)", args: args{str: "False"}},
		{name: "Test parseBool (has error, use default value)", args: args{str: "xxxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.args.str); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerRunOptions_secureServer(t *testing.T) {
	convey.Convey("Test ServerRunOptions secureServer", t, func() {
		s := &ServerRunOptions{}
		err := s.secureServer()
		convey.So(err, convey.ShouldBeNil)

		s.Listen.EnableHTTPS = true
		p := gomonkey.ApplyFuncSeq(os.ReadFile, []gomonkey.OutputCell{
			{Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]byte{}, nil}}, {Values: gomonkey.Params{[]byte{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{[]byte{}, nil}}, {Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}}, {Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}}, {Values: gomonkey.Params{[]byte{}, nil}},
			{Values: gomonkey.Params{[]byte{}, nil}}, {Values: gomonkey.Params{[]byte{}, nil}},
		})
		defer p.Reset()
		p.ApplyFuncSeq(tls.X509KeyPair, []gomonkey.OutputCell{
			{Values: gomonkey.Params{tls.Certificate{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{tls.Certificate{}, nil}},
			{Values: gomonkey.Params{tls.Certificate{}, nil}},
			{Values: gomonkey.Params{tls.Certificate{}, nil}},
		})
		p.ApplyFuncSeq(cert.NewPool, []gomonkey.OutputCell{
			{Values: gomonkey.Params{&x509.CertPool{}, fmt.Errorf("mock error")}},
			{Values: gomonkey.Params{&x509.CertPool{}, nil}},
		})
		err = s.secureServer()
		convey.So(err, convey.ShouldNotBeNil)
		err = s.secureServer()
		convey.So(err, convey.ShouldNotBeNil)
		err = s.secureServer()
		convey.So(err, convey.ShouldNotBeNil)
		err = s.secureServer()
		convey.So(err, convey.ShouldBeNil)

		s.Listen.TrustCaFile = "x"
		err = s.secureServer()
		convey.So(err, convey.ShouldNotBeNil)
		err = s.secureServer()
		convey.So(err, convey.ShouldBeNil)
	})
}
