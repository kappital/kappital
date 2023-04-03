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

package convert

import (
	"fmt"
	"strings"
)

const (
	// the GVK resources for the kubernetes
	clusterRoleGVK        = "rbac.authorization.k8s.io/v1/ClusterRole"
	clusterRoleBindingGVK = "rbac.authorization.k8s.io/v1/ClusterRoleBinding"
	deploymentGVK         = "apps/v1/Deployment"
	serviceAccountGVK     = "v1/ServiceAccount"
	crdV1GVK              = "apiextensions.k8s.io/v1/CustomResourceDefinition"
	crdV1Beta1GVK         = "apiextensions.k8s.io/v1beta1/CustomResourceDefinition"
	// the GVK resource for kappital
	csdGVK = "core.kappital.io/v1alpha1/CustomServiceDefinition"

	// directory limitation for the directory resolution
	maxFileCount           = 1024
	maxDepth               = 3
	decodeBufferSize       = 30
	initResourceSize       = 10
	singleFileSize   int64 = 1024 * 1024 * 1
	totalFileSize    int64 = 1024 * 1024 * 10
)

var (
	errOverMaxFileCount = fmt.Errorf("the file count is over the max file count %d", maxFileCount)
	errOverMaxDirDepth  = fmt.Errorf("the directory depth is over the max deep level %d", maxDepth)
	errOverFileSize     = fmt.Errorf("the file size is over the max file size %d or total file size %d", singleFileSize, totalFileSize)
)

func isValidFileType(path string) bool {
	tmp := strings.ToLower(path)
	return strings.HasSuffix(tmp, ".json") || strings.HasSuffix(tmp, ".yaml") || strings.HasSuffix(tmp, ".yml")
}
