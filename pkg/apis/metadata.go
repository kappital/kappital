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

package apis

import (
	"k8s.io/klog/v2"
)

// Descriptor is a list of textual information of a Cloud Native Service
type Descriptor struct {
	// Name of the service, services with a same name can be distinguished by their version
	Name string `json:"name,omitempty"`
	// Version of the service.
	Version string `json:"version,omitempty"`
	// Service type declaration, can be either helm or operator, now only offer the operator
	Type ServiceType `json:"type,omitempty"`
	// BriefDescription describes the Cloud Native Service in one statement.
	BriefDescription string `json:"briefDescription,omitempty"`
	// Detail description of the Cloud Native Service, in format of markdown.
	Detail string `json:"detail,omitempty"`
	// Logo contains service logo content/URL
	Logo Logo `json:"logo,omitempty"`
	// Maintainers contains necessary information to get in touch with maintainers of the Cloud Native Service.
	Maintainers []Maintainer `json:"maintainers,omitempty"`
	// Provider contains necessary information to get in touch with provider of the Cloud Native Service.
	Provider AppLink `json:"provider,omitempty"`
	// Source indicates the source of the Cloud Native Service, open source or from some companies.
	Source string `json:"source,omitempty"`
	// Industries declares appropriate industries of the Cloud Native Service.
	Industries []string `json:"industries,omitempty"`
	// Links contains links of information related to the Cloud Native Service itself.
	Links []AppLink `json:"links,omitempty"`
	// Keywords are keywords of the Cloud Native Service
	Keywords []string `json:"keywords,omitempty"`

	// MinKubeVersion is the minimum version of kubernetes needed for this service to deploy successfully
	MinKubeVersion string `json:"minKubeVersion,omitempty"`
	// Architectures indicates nodes of which architectures its instances can run on.
	// e.g. x86_64, arm, etc.
	Architectures []string `json:"architecture,omitempty"`
	// Capabilities declares typical usages of the Cloud Native Service.
	Capabilities []string `json:"capabilities,omitempty"`
	// Categories declares appropriate scenarios of the Cloud Native Service.
	Categories []string `json:"categories,omitempty"`
	// Devices declares required devices on target nodes to ensure its instance work well.
	// e.g. GPU, NPU, etc.
	Devices []string `json:"devices,omitempty"`
	// Scenes indicates in which scenes its instances can survive.
	Scenes []string `json:"scenes,omitempty"`
}

// AppLink contains name and URL to connect to certain app.
type AppLink struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Maintainer includes name and email to know a maintainer.
type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Logo contains necessary information to load a logo.
type Logo struct {
	// Base64 formatted content of the logo.
	// Optional, can skip if URL provided.
	Data string `json:"base64data,omitempty"`
	// Media type of the logo.
	MediaType string `json:"mediaType,omitempty"`
}

// Validation does the Descriptor is valid or not, and generate some attribute
func (d *Descriptor) Validation() bool {
	// 1. check the key attribute is empty or larger than the max length
	if len(d.Name) == 0 || len(d.Name) > maxStringLen {
		klog.Errorf("the Service Package NAME cannot be empty or longer than the %d bytes", maxStringLen)
		return false
	}
	if len(d.Version) == 0 {
		klog.Warningf("because does not have service package version, set the version as %s", DefaultVersion)
		d.Version = DefaultVersion
	}
	if len(d.Version) > maxStringLen {
		klog.Errorf("the Service Package VERSION cannot longer than the %d bytes", maxStringLen)
		return false
	}
	// 2. check the "Source" attribute is empty or not, if it is empty, set to the default value
	if len(d.Source) == 0 {
		d.Source = defaultSource
	}
	// 3. check the "Type" is valid or not
	if _, ok := serviceTypeSet[d.Type]; !ok {
		klog.Errorf("the service type is invalid")
		return false
	}
	d.GenerateSlice()
	return true
}

// GenerateSlice generate the slice attributes, if the slice is nil, set it to the empty slice
// because when insert the data into database will have "null" string for the related attribute
func (d *Descriptor) GenerateSlice() {
	if d.Maintainers == nil {
		d.Maintainers = []Maintainer{}
	}
	if d.Industries == nil {
		d.Industries = []string{}
	}
	if d.Links == nil {
		d.Links = []AppLink{}
	}
	if d.Keywords == nil {
		d.Keywords = []string{}
	}
	if d.Architectures == nil {
		d.Architectures = []string{}
	}
	if d.Capabilities == nil {
		d.Capabilities = []string{}
	}
	if d.Categories == nil {
		d.Categories = []string{}
	}
	if d.Devices == nil {
		d.Devices = []string{}
	}
	if d.Scenes == nil {
		d.Scenes = []string{}
	}
}
