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
	"crypto/tls"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/kappital/kappital/pkg/utils/cryption"
)

func TestHealthAndReadinessProvider(t *testing.T) {
	cert, key, err := cryption.GetSelfCertAndKey()
	if err != nil {
		t.Errorf("cannot create the self certificate and key, err: %s", err)
	}
	ip, err := GetLocalIP()
	if err != nil {
		t.Errorf("cannot get localhost ip")
	}
	p := gomonkey.ApplyFunc(os.Exit, func(_ int) {})
	defer p.Reset()
	HealthAndReadinessProvider(ip, cert, key)
}

func Test_constructTLSConfig(t *testing.T) {
	type args struct {
		cert []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *tls.Config
		wantErr bool
	}{
		{name: "Test constructTLSConfig", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := constructTLSConfig(tt.args.cert, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("constructTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructTLSConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_constructServer(t *testing.T) {
	type args struct {
		address string
		cert    []byte
		key     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Server
		wantErr bool
	}{
		{name: "Test constructServer", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := constructServer(tt.args.address, tt.args.cert, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("constructServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("constructServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
