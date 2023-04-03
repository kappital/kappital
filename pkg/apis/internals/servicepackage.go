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

package internals

import (
	"fmt"
	"strings"
)

// CNSIdentifier is the unique identifier for a Cloud Native Service
type CNSIdentifier struct {
	Repo    string
	Name    string
	Version string
}

// GetYamlFileName get the name of yaml file
func (c CNSIdentifier) GetYamlFileName() string {
	return fmt.Sprintf("%s-%s-%s.yaml", c.Repo, c.Name, c.Version)
}

// GetTmpDirName get the temporary directory name
func (c CNSIdentifier) GetTmpDirName() string {
	return fmt.Sprintf("%s-%s-%s", c.Repo, c.Name, c.Version)
}

// GetFilterMap for getting result from database
func (c CNSIdentifier) GetFilterMap() map[string]string {
	filter := make(map[string]string, 2)
	if len(c.Repo) != 0 {
		filter["repository"] = c.Repo
	}
	if len(c.Name) != 0 {
		filter["name"] = c.Name
	}
	return filter
}

// GetResourceName for the audit log
func (c CNSIdentifier) GetResourceName() string {
	strs := make([]string, 0, 3)
	if len(c.Name) > 0 {
		strs = append(strs, fmt.Sprintf("Service Package [%s]", c.Name))
	}
	if len(c.Repo) > 0 {
		strs = append(strs, fmt.Sprintf("Repository [%s]", c.Repo))
	}
	if len(c.Version) > 0 {
		strs = append(strs, fmt.Sprintf("Version [%s]", c.Version))
	}
	return strings.Join(strs, "; ")
}
