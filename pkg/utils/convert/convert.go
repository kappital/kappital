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
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
)

type Convert interface {
	TransferToCloudNativeService(filePath string) (*svcv1alpha1.CloudNativeService, error)
}

func GetLoader() Convert {
	return newCloudNativePackageLoader()
}
