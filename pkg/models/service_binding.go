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

package models

import (
	"fmt"
	"time"

	"github.com/kappital/kappital/pkg/apis/internals"
	"github.com/kappital/kappital/pkg/utils/uuid"
)

// ServiceBindingModel defines the table fields of service_binding_model in database
type ServiceBindingModel struct {
	ID                       string    `orm:"size(40);pk;column(id)"`
	Name                     string    `orm:"size(64);column(name)"`
	Version                  string    `orm:"size(64);column(service_version)"`
	Namespace                string    `orm:"size(64);default(kappital-system);column(namespace)"`
	ServiceName              string    `orm:"size(64);column(service_name)"`
	ClusterName              string    `orm:"size(64);default(default);column(cluster_name)"`
	ServiceID                string    `orm:"size(40);column(service_id)"`
	VersionID                string    `orm:"size(40);column(version_id)"`
	Status                   string    `orm:"size(64);default(Pending);column(status)"`
	ErrorMessage             string    `orm:"type(text);null;column(error_message)"`
	Workloads                string    `orm:"type(text);null;column(workloads)"`
	Permissions              string    `orm:"type(text);null;column(permissions)"`
	CustomResourceDefinition string    `orm:"type(text);null;column(crd)"`
	CapabilityPlugin         string    `orm:"type(text);null;column(capability_plugin)"`
	CreateTime               time.Time `orm:"type(datetime);auto_now_add;column(create_timestamp)"`
	UpdateTime               time.Time `orm:"type(datetime);null;column(update_timestamp)"`
	ProcessTime              time.Time `orm:"type(datetime);null;column(process_timestamp)"`

	Resources []*ResourceModel `json:"resources" orm:"null;reverse(many)"`
}

// TableUnique makes combined columns unique
func (s *ServiceBindingModel) TableUnique() [][]string {
	return [][]string{{"name", "cluster_name"}}
}

// Generate fills a service_binding_model record with id and timestamps
func (s *ServiceBindingModel) Generate(currTimestamp time.Time, isUpdate bool) {
	if len(s.ID) == 0 {
		s.ID = uuid.NewUUID()
	}
	if isUpdate {
		s.UpdateTime = currTimestamp
	} else {
		if s.CreateTime.Equal(time.Time{}) {
			s.CreateTime = currTimestamp
		}
	}
	for i := range s.Resources {
		s.Resources[i].Generate(currTimestamp, isUpdate)
	}
}

// GetResourceIDMap get the ServiceBinding's ResourceModel ID map
// Key is the format of ResourceModel.Kind;ResourceModel.APIVersion;ResourceModel.Resource
func (s *ServiceBindingModel) GetResourceIDMap() map[string]string {
	resourceMap := make(map[string]string, len(s.Resources))
	for _, resource := range s.Resources {
		key := fmt.Sprintf("%s;%s;%s", resource.Kind, resource.APIVersion, resource.Resource)
		resourceMap[key] = resource.ID
	}
	return resourceMap
}

// ResourceModel defines the table fields of resource_model in database
type ResourceModel struct {
	ID              string    `orm:"size(40);pk;column(id)"`
	Kind            string    `orm:"size(64);column(kind)"`
	Group           string    `orm:"size(64);column(group)"`
	APIVersion      string    `orm:"size(64);column(api_version)"`
	Resource        string    `orm:"size(64);column(resource)"`
	CreateTimestamp time.Time `orm:"type(datetime);auto_now_add;column(create_timestamp)"`
	UpdateTimestamp time.Time `orm:"type(datetime);null;column(update_timestamp)"`

	Instances []*InstanceModel `json:"instances" orm:"null;reverse(many)"`

	ServiceBinding *ServiceBindingModel `orm:"null;rel(fk)"`
}

// Generate fills a resource_model record with id and timestamps
func (r *ResourceModel) Generate(currTimestamp time.Time, isUpdate bool) {
	if len(r.ID) == 0 {
		r.ID = uuid.NewUUID()
	}
	if isUpdate {
		r.UpdateTimestamp = currTimestamp
	} else {
		if r.CreateTimestamp.Equal(time.Time{}) {
			r.CreateTimestamp = currTimestamp
		}
	}

	for i := range r.Instances {
		r.Instances[i].Generate(currTimestamp, isUpdate)
	}
}

// GetInstanceToResourceMap get the (instance's key) - (resource) map,
// Key is the format of InstanceModel.Name;InstanceModel.Namespace;InstanceModel.ClusterName
func (r *ResourceModel) GetInstanceToResourceMap() (map[string]ResourceModel, map[string]InstanceModel) {
	resourceMap := make(map[string]ResourceModel, len(r.Instances))
	instanceMap := make(map[string]InstanceModel, len(r.Instances))
	for _, instance := range r.Instances {
		key := fmt.Sprintf("%s;%s;%s", instance.Name, instance.Namespace, instance.ClusterName)
		resourceMap[key] = *r
		instanceMap[key] = *instance
	}
	return resourceMap, instanceMap
}

// InstanceModel defines the table fields of instance_model in database
type InstanceModel struct {
	ID                  string                 `orm:"size(40);pk;column(id)"`
	Kind                string                 `orm:"size(64);column(kind)"`
	APIVersion          string                 `orm:"size(64);column(api_version)"`
	Name                string                 `orm:"size(64);column(name)"`
	Namespace           string                 `orm:"size(64);column(namespace)"`
	ServiceBindingID    string                 `orm:"size(40);column(service_binding_id)"`
	RawResource         string                 `orm:"type(text);null;column(raw_resource)"`
	Status              string                 `orm:"size(64);column(status)"`
	ErrorMessage        string                 `orm:"type(text);null;column(error_message)"`
	ServiceID           string                 `orm:"size(40);column(service_id)"`
	InstanceType        internals.InstanceType `orm:"size(64);column(instance_type)"`
	ServiceName         string                 `orm:"size(64);column(service_name)"`
	ClusterID           string                 `orm:"size(64);column(cluster_id)"`
	ClusterName         string                 `orm:"size(64);column(cluster_name)"`
	PackageDependencies string                 `orm:"type(text);null;column(pkg_dependencies)"`
	CreateTimestamp     time.Time              `orm:"type(datetime);auto_now_add;column(create_timestamp)"`
	ProcessTime         time.Time              `orm:"type(datetime);null;column(process_time)"`
	UpdateTime          time.Time              `orm:"type(datetime);null;column(update_timestamp)"`
	InstallState        string                 `orm:"type(text);column(install_state)"`

	Resource *ResourceModel `orm:"null;rel(fk)"`
}

// Generate fills a instance_model record with id and timestamps
func (i *InstanceModel) Generate(currTimestamp time.Time, isUpdate bool) {
	if len(i.ID) == 0 {
		i.ID = uuid.NewUUID()
	}
	if isUpdate {
		i.UpdateTime = currTimestamp
	} else {
		if i.CreateTimestamp.Equal(time.Time{}) {
			i.CreateTimestamp = currTimestamp
		}
	}
}
