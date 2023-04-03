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

package instance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/brahma-adshonor/gohook"
	"github.com/smartystreets/goconvey/convey"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/view"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

func Test_convertInstanceToTable(t *testing.T) {
	type args struct {
		ins         models.InstanceModel
		serviceName string
		phase       instancev1alpha1.Phase
	}
	tests := []struct {
		name string
		args args
		want view.Instance
	}{
		{
			name: "Test convertInstanceToTable (SucceededPhase)",
			args: args{
				ins:         models.InstanceModel{Name: "test", Namespace: "test", Status: "x1", ClusterName: "x1"},
				serviceName: "test",
				phase:       instancev1alpha1.SucceededPhase,
			},
			want: view.Instance{InstanceName: "test", Namespace: "test", ServiceName: "test", ClusterName: "x1", Status: "x1"},
		},
		{
			name: "Test convertInstanceToTable (PendingPhase)",
			args: args{
				ins:         models.InstanceModel{Name: "test", Namespace: "test", Status: "x1", ClusterName: "x1"},
				serviceName: "test",
				phase:       instancev1alpha1.PendingPhase,
			},
			want: view.Instance{InstanceName: "test", Namespace: "test", ServiceName: "test", ClusterName: "x1", Status: "Pending"},
		},
		{
			name: "Test convertInstanceToTable (UnknownPhase)",
			args: args{
				ins:         models.InstanceModel{Name: "test", Namespace: "test", Status: "x1", ClusterName: "x1"},
				serviceName: "test",
				phase:       instancev1alpha1.UnknownPhase,
			},
			want: view.Instance{InstanceName: "test", Namespace: "test", ServiceName: "test", ClusterName: "x1", Status: "Unknown"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertInstanceToTable(tt.args.ins, tt.args.serviceName, tt.args.phase)
			got.Created = ""
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertInstanceToTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_operation_NewCommand(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{name: "Test operation NewCommand for get instance"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &operation{}
			if got := o.NewCommand(); (got == nil) != tt.wantNil {
				t.Errorf("get instance operation NewCommand() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_operation_PreRunE(t *testing.T) {
	if err := gohook.Hook(kappctl.GetConfig, func() (*kappctl.Config, error) {
		return &kappctl.Config{ManagerHTTPSServer: "https://x.x.x.x:x"}, nil
	}, nil); err != nil {
		t.Errorf("hook err:%v", err)
		return
	}
	defer gohook.UnHook(kappctl.GetConfig) //nolint:errcheck

	type fields struct {
		serviceName string
		clusterName string
		allResult   bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    []string
		wantErr bool
	}{
		{name: "Test operation PreRunE (no args, not get All, no service name)", wantErr: true},
		{name: "Test operation PreRunE (get all)", fields: fields{allResult: true}},
		{name: "Test operation PreRunE (only has one args)", args: []string{"test"}, wantErr: true},
		{
			name:    "Test operation PreRunE (only has one args and error cluster name)",
			fields:  fields{serviceName: "test", clusterName: "abc-"},
			args:    []string{"test"},
			wantErr: true,
		},
		{
			name:   "Test operation PreRunE (has one args and service name)",
			fields: fields{serviceName: "test"},
			args:   []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &operation{
				serviceName: tt.fields.serviceName,
				clusterName: tt.fields.clusterName,
				allResult:   tt.fields.allResult,
			}
			if err := o.PreRunE(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("PreRunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_operation_RunE(t *testing.T) {
	convey.Convey("test operation RunE", t, func() {
		o := &operation{config: &kappctl.Config{}}
		globalMock := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
			return http.StatusOK, []byte(""), nil
		})
		defer globalMock.Reset()
		convey.Convey("case 1: getAll and getAllServiceInstances has error", func() {
			p := gomonkey.ApplyFunc(o.getAllServiceInstances, func() ([]interface{}, error) {
				return nil, fmt.Errorf("mock error")
			})
			defer p.Reset()
			o.allResult = true
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2: getAll and getAllServiceInstances without error", func() {
			p := gomonkey.ApplyFunc(o.getAllServiceInstances, func() ([]interface{}, error) {
				return []interface{}{}, nil
			})
			defer p.Reset()
			p.ApplyFunc(json.Unmarshal, func(_ []byte, _ interface{}) error {
				return nil
			})
			o.allResult = true
			err := o.RunE()
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("case 3: getAll and getServiceInstances with error", func() {
			p := gomonkey.ApplyFunc(o.getServiceInstances, func() ([]interface{}, error) {
				return nil, fmt.Errorf("mock error")
			})
			defer p.Reset()
			err := o.RunE()
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 4: getAll and getServiceInstances has nil res", func() {
			p := gomonkey.ApplyFunc(o.getServiceInstances, func() ([]interface{}, error) {
				return nil, nil
			})
			defer p.Reset()
			p.ApplyFunc(json.Unmarshal, func(_ []byte, _ interface{}) error {
				return nil
			})
			err := o.RunE()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_operation_getAllServiceInstances_withHttpRequestError(t *testing.T) {
	convey.Convey("test operation getAllServiceInstances withHttpRequestError", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: error url for the http request", func() {
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2; http request cannot get the http.StatusOK", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return 0, nil, nil
			})
			defer p.Reset()
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: http request get invalid buf", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return http.StatusOK, []byte("invalid"), nil
			})
			defer p.Reset()
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 4: http request get valid buf but no CloudNativeServiceInstance", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				var svcs []instancev1alpha1.CloudNativeServiceInstance
				jsonBytes, err := json.Marshal(svcs)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_operation_getAllServiceInstances(t *testing.T) {
	convey.Convey("test operation getAllServiceInstances", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 2: getInstanceListServiceName with error", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				svcs := []instancev1alpha1.CloudNativeServiceInstance{{}}
				jsonBytes, err := json.Marshal(svcs)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: getInstanceListServiceName without error", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(info *gateway.RequestInfo) (int, []byte, error) {
				if strings.Contains(info.Path, "instance") {
					ins := []models.InstanceModel{{Name: "1"}}
					jsonBytes, err := json.Marshal(ins)
					return http.StatusOK, jsonBytes, err
				}
				svcs := []instancev1alpha1.CloudNativeServiceInstance{{ObjectMeta: metav1.ObjectMeta{Name: "1"}}}
				jsonBytes, err := json.Marshal(svcs)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			got, err := o.getAllServiceInstances()
			convey.So(got, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_operation_getInstanceListServiceName_withError(t *testing.T) {
	convey.Convey("test operation getInstanceListServiceName withError", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: error url for the http request", func() {
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2; http request cannot get the http.StatusOK", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return 0, nil, nil
			})
			defer p.Reset()
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: http request get invalid buf", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return http.StatusOK, []byte("invalid"), nil
			})
			defer p.Reset()
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func Test_operation_getInstanceListServiceName_withoutError(t *testing.T) {
	convey.Convey("test operation getInstanceListServiceName withoutError", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: http request get correct result (w/ outputFormat)", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				ins := []models.InstanceModel{{}}
				jsonBytes, err := json.Marshal(ins)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			o.outputFormat = "json"
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("case 2: http request get correct result (w/o outputFormat and w/o instanceName)", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				ins := []models.InstanceModel{{}}
				jsonBytes, err := json.Marshal(ins)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			o.outputFormat = ""
			o.instanceName = ""
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("case 3: http request get correct result (w/o outputFormat and w/ instanceName)", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				ins := []models.InstanceModel{{Name: "1"}}
				jsonBytes, err := json.Marshal(ins)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			o.outputFormat = ""
			o.instanceName = "1"
			got, err := o.getInstanceListServiceName(instancev1alpha1.CloudNativeServiceInstance{})
			convey.So(got, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_operation_getServiceInstances(t *testing.T) {
	convey.Convey("test operation getServiceInstances", t, func() {
		o := &operation{config: &kappctl.Config{}}
		convey.Convey("case 1: error url for the http request", func() {
			got, err := o.getServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 2; http request cannot get the http.StatusOK", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return 0, nil, nil
			})
			defer p.Reset()
			got, err := o.getServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 3: http request get invalid buf", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				return http.StatusOK, []byte("invalid"), nil
			})
			defer p.Reset()
			got, err := o.getServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("case 4: http request get correct result", func() {
			p := gomonkey.ApplyFunc(gateway.CommonUtilRequest, func(_ *gateway.RequestInfo) (int, []byte, error) {
				svc := instancev1alpha1.CloudNativeServiceInstance{}
				jsonBytes, err := json.Marshal(svc)
				return http.StatusOK, jsonBytes, err
			})
			defer p.Reset()
			got, err := o.getServiceInstances()
			convey.So(got, convey.ShouldBeNil)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func Test_outputYamlOrJSON(t *testing.T) {
	type args struct {
		svc    instancev1alpha1.CloudNativeServiceInstance
		ins    []models.InstanceModel
		format string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test outputYamlOrJSON (pending)",
			args: args{
				svc: instancev1alpha1.CloudNativeServiceInstance{
					Status: instancev1alpha1.CloudNativeServiceInstanceStatus{Phase: instancev1alpha1.PendingPhase},
				},
				ins: []models.InstanceModel{{}},
			},
		},
		{
			name: "Test outputYamlOrJSON (not succeeded and pending)",
			args: args{
				svc: instancev1alpha1.CloudNativeServiceInstance{
					Status: instancev1alpha1.CloudNativeServiceInstanceStatus{Phase: instancev1alpha1.FailedPhase},
				},
				ins: []models.InstanceModel{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputYamlOrJSON(tt.args.svc, tt.args.ins, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("outputYamlOrJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
