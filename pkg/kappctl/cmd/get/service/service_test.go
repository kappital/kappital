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

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"

	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

var (
	multipleBytes []byte
	emptySlice    []byte
	singleBytes   []byte
)

func Test_operation_NewCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "Test get service operation NewCommand"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &operation{}
			if got := o.NewCommand(); (got == nil) != tt.wantNil {
				t.Errorf("get service operation NewCommand() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_operation_PreRunE(t *testing.T) {
	convey.Convey("test get service operation PreRunE", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: cannot get the kappctl config for get service", func() {
			err := o.PreRunE([]string{""})
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: can get the kappctl config and valid args slice length for get service", func() {
			p := gomonkey.ApplyFunc(kappctl.GetConfig, func() (*kappctl.Config, error) {
				return &kappctl.Config{}, nil
			})
			defer p.Reset()
			err := o.PreRunE([]string{""})
			convey.So(err, convey.ShouldBeNil)
			err = o.PreRunE([]string{"123456789012345678901234567890123456789012345678901ee234567890123456789012345678901234567890"})
			convey.So(err, convey.ShouldNotBeNil)
			err = o.PreRunE([]string{"", ""})
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("test get service operation RunE", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: gateway.CommonUtilRequest get error for get service", func() {
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: gateway.CommonUtilRequest get failed http code for get service", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return 0, nil, nil
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: gateway.CommonUtilRequest get valid data", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return http.StatusOK, singleBytes, nil
			})
			defer p.Reset()
			o.serviceName = "xxxx"
			err := o.RunE()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_outputResult(t *testing.T) {
	type args struct {
		buf        []byte
		isMultiple bool
		format     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test outputResult (has format)", args: args{buf: []byte{}, isMultiple: false, format: "json"}},
		{name: "Test outputResult (multiple, but Unmarshal error)", args: args{isMultiple: true}, wantErr: true},
		{name: "Test outputResult (multiple, but no result)", args: args{buf: emptySlice, isMultiple: true}},
		{name: "Test outputResult (multiple without error)", args: args{buf: multipleBytes, isMultiple: true}},
		{name: "Test outputResult (single, with Unmarshal error)", args: args{isMultiple: false}, wantErr: true},
		{name: "Test outputResult (single without error)", args: args{buf: singleBytes, isMultiple: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(tt.args.buf, tt.args.isMultiple, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func init() {
	var err error
	multipleBytes, err = json.Marshal([]instancev1alpha1.CloudNativeServiceInstance{{}})
	if err != nil {
		fmt.Println("cannot get multiple CloudNativeServiceInstance bytes")
	}
	emptySlice, err = json.Marshal([]instancev1alpha1.CloudNativeServiceInstance{})
	if err != nil {
		fmt.Println("cannot get empty slice CloudNativeServiceInstance bytes")
	}
	singleBytes, err = json.Marshal(instancev1alpha1.CloudNativeServiceInstance{})
	if err != nil {
		fmt.Println("cannot get single CloudNativeServiceInstance bytes")
	}
}
