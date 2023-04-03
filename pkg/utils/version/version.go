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

package version

import (
	"fmt"
	"runtime"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	// Flags of get the version of module
	Flags = sets.NewString("version", "--version", "-v")

	gitVersion   = "v0.0.0-master"
	gitCommit    = "unknown" // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = "unknown" // state of git tree, either "clean" or "dirty"

	buildDate = "unknown" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

const (
	// ServiceNameManager defines name of current manager service component
	ServiceNameManager = "kappital-manager"
	// ServiceNameEngine defines name of current engine service component
	ServiceNameEngine = "kappital-engine"
)

// Info contains versioning information.
type Info struct {
	ServiceName  string `json:"ServiceName"`
	GitVersion   string `json:"GitVersion"`
	GitCommit    string `json:"GitCommit"`
	GitTreeState string `json:"GitTreeState"`
	BuildDate    string `json:"BuildDate"`
	GoVersion    string `json:"GoVersion"`
	Compiler     string `json:"Compiler"`
	Platform     string `json:"Platform"`
}

// String returns a Go-syntax representation of the Info.
func (info Info) String() string {
	return fmt.Sprintf("%#v", info)
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get(serviceName string) Info {
	return Info{
		ServiceName:  serviceName,
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
