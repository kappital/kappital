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

package view

// Repository defines the fields of the table as default output of `kappctl get repo`
type Repository struct {
	Name         string
	Public       string
	ServiceCount string
	Created      string
}

// Service defines the fields of the table as default output of `kappctl get service`
type Service struct {
	Name      string
	Cluster   string
	Namespace string
	Phase     string
	Message   string
	Created   string
}

// Instance defines the fields of the table as default output of `kappctl get instance`
type Instance struct {
	InstanceName string
	Namespace    string
	ServiceName  string
	ClusterName  string
	Status       string
	Created      string
}

// Package defines the fields of the table as default output of `kappctl get package`
type Package struct {
	Repository string
	Name       string
	Type       string
	Created    string
}

// Version defines the fields of the table as default output of `kappctl get package`
type Version struct {
	Repository string
	Name       string
	Type       string
	Version    string
	Keywords   string
	Status     string
	Created    string
}
