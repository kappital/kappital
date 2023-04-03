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

package manager

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/mock"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	"github.com/kappital/kappital/pkg/resource"
	"github.com/kappital/kappital/pkg/utils/audit"
)

func TestMain(m *testing.M) {
	audit.SetAuditLog(&audit.FakeAuditLogger)
	ctx, _ := mock.NewMockContext(&http.Request{})

	testInstanceController = &InstanceController{
		Controller: web.Controller{Ctx: ctx},
		instance:   resource.InstanceResource{},
		binding:    resource.ServiceBindingResource{},
	}
	testInstanceController.Ctx.Request = &http.Request{}

	testServiceBindingController = &ServiceBindingController{
		Controller: web.Controller{Ctx: ctx},
		resource:   resource.ServiceBindingResource{},
	}
	testServiceBindingController.Ctx.Request = &http.Request{}

	m.Run()
}

func Test_deploymentWorkloadBuilder(t *testing.T) {
	tests := []struct {
		name        string
		deployments []appsv1.Deployment
		want        []enginev1alpha1.ServiceDeploymentSpec
	}{
		{
			name:        "test deploymentWorkloadBuilder",
			deployments: []appsv1.Deployment{{ObjectMeta: metav1.ObjectMeta{Name: "test-deployment"}}},
			want:        []enginev1alpha1.ServiceDeploymentSpec{{Name: "test-deployment"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deploymentWorkloadBuilder(tt.deployments); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deploymentWorkloadBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceCRDBuilder(t *testing.T) {
	tests := []struct {
		name    string
		csds    []svcv1alpha1.CustomServiceDefinition
		want    []string
		wantErr bool
	}{
		{
			name: "Test serviceCRDBuilder (without error)",
			csds: []svcv1alpha1.CustomServiceDefinition{
				{Spec: svcv1alpha1.CustomServiceDefinitionSpec{CRD: &apis.AbstractResource{}}},
			},
			want:    []string{"{\"metadata\":{\"creationTimestamp\":null}}"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serviceCRDBuilder(tt.csds)
			if (err != nil) != tt.wantErr {
				t.Errorf("serviceCRDBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serviceCRDBuilder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_servicePermissionBuilder(t *testing.T) {
	type args struct {
		crs  []rbacv1.ClusterRole
		crbs []rbacv1.ClusterRoleBinding
	}
	tests := []struct {
		name string
		args args
		want []enginev1alpha1.Permission
	}{
		{
			name: "Test servicePermissionBuilder",
			args: args{
				crs: []rbacv1.ClusterRole{{ObjectMeta: metav1.ObjectMeta{Name: "test-cluster-role"}}},
				crbs: []rbacv1.ClusterRoleBinding{
					{
						Subjects: []rbacv1.Subject{{Name: "test-service-account-name"}},
						RoleRef:  rbacv1.RoleRef{Name: "test-cluster-role"},
					},
					{RoleRef: rbacv1.RoleRef{Name: "not exist name"}},
				},
			},
			want: []enginev1alpha1.Permission{{ServiceAccountName: "test-service-account-name"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := servicePermissionBuilder(tt.args.crs, tt.args.crbs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("servicePermissionBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getClusterRoleMap(t *testing.T) {
	tests := []struct {
		name string
		crs  []rbacv1.ClusterRole
		want map[string]rbacv1.ClusterRole
	}{
		{
			name: "Test getClusterRoleMap",
			crs:  []rbacv1.ClusterRole{{ObjectMeta: metav1.ObjectMeta{Name: "test cluster role"}}},
			want: map[string]rbacv1.ClusterRole{
				"test cluster role": {ObjectMeta: metav1.ObjectMeta{Name: "test cluster role"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClusterRoleMap(tt.crs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClusterRoleMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceCapabilityPluginBuilder(t *testing.T) {
	tests := []struct {
		name string
		want enginev1alpha1.CapabilityPlugin
	}{
		{name: "Test serviceCapabilityPluginBuilder"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceCapabilityPluginBuilder(nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serviceCapabilityPluginBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
