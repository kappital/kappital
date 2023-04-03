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
	"bytes"
	"fmt"
	"strconv"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/utils/client"
)

const (
	// clusterVersion15 enum kubernetes version for 1.15
	clusterVersion15 = "v1.15"
	// clusterVersion16To22 enum kubernetes version between 1.16 to 1.22
	clusterVersion16To22 = "v1.16-v1.22"
	// clusterVersion22Plus enum kubernetes version for 1.22 and above
	clusterVersion22Plus = "v1.22+"

	clusterVersion15Num = 15
	clusterVersion22Num = 22
)

const (
	apiExtensionsV1      = "apiextensions.k8s.io/v1"
	apiExtensionsV1Beta1 = "apiextensions.k8s.io/v1beta1"
)

var clusterVersion = clusterVersion15

// GetClusterVersion get the version of cluster, it will use for the crd
func GetClusterVersion() string {
	return clusterVersion
}

// SetClusterVersion set the version of cluster, it will be used for the unit test to set a fake cluster version
func SetClusterVersion(version string) {
	clusterVersion = version
}

// SupportCRDUseV1Beta1 the cluster version 1.22 and below will support v1beta1 version crd,
// but, 1.22 above cluster will discard v1beta1 version
func SupportCRDUseV1Beta1() bool {
	return clusterVersion != clusterVersion22Plus
}

// SupportCRDUseV1 the cluster version 1.15 above will support apiextensionsv1 for the crd
func SupportCRDUseV1() bool {
	return clusterVersion != clusterVersion15
}

// SetupClusterVersion get the kubernetes cluster version, and save it in the runtime memory
func SetupClusterVersion() error {
	// get the client
	cli := client.GetCRDClient()
	versionInfo, err := cli.ServerVersion()
	if err != nil {
		klog.Errorf("unable to get service version, err: %s", err)
		return err
	}
	versionNum, err := getSecondLevelVersion(versionInfo.GitVersion)
	if err != nil {
		klog.Errorf("unable to get current server version, err: %s", err)
		return err
	}
	if versionNum <= clusterVersion15Num {
		clusterVersion = clusterVersion15
	} else if clusterVersion15Num < versionNum && versionNum <= clusterVersion22Num {
		clusterVersion = clusterVersion16To22
	} else {
		clusterVersion = clusterVersion22Plus
	}
	return nil
}

func getSecondLevelVersion(version string) (int, error) {
	versionLevels := strings.Split(version, ".")
	if len(versionLevels) < 3 {
		return -1, fmt.Errorf("illeagle version string")
	}
	return strconv.Atoi(versionLevels[1])
}

// GetCrdV1 get and analysis the crd to apiextensionsv1
func GetCrdV1(crdString string) (apiextensionsv1.CustomResourceDefinition, error) {
	crd := apiextensionsv1.CustomResourceDefinition{}
	err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(crdString)), 30).Decode(&crd)
	if err != nil || crd.APIVersion != apiExtensionsV1 {
		return crd, fmt.Errorf("not v1 crd")
	}
	return crd, nil
}

// GetCrdV1Beta1 get and analysis the crd to v1beta1
func GetCrdV1Beta1(crdString string) (apiextensionsv1beta1.CustomResourceDefinition, error) {
	crd := apiextensionsv1beta1.CustomResourceDefinition{}
	err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(crdString)), 30).Decode(&crd)
	if err != nil || crd.APIVersion != apiExtensionsV1Beta1 {
		return crd, fmt.Errorf("not v1beta1 crd")
	}
	return crd, nil
}

// GetCrdV1AndBeta1Slice get the CRD to slices for apiextensionsv1 and v1beta1
func GetCrdV1AndBeta1Slice(crdStrings []string) ([]apiextensionsv1.CustomResourceDefinition,
	[]apiextensionsv1beta1.CustomResourceDefinition) {
	crdV1s := make([]apiextensionsv1.CustomResourceDefinition, 0, len(crdStrings))
	crdV1Beta1s := make([]apiextensionsv1beta1.CustomResourceDefinition, 0, len(crdStrings))
	for _, s := range crdStrings {
		if SupportCRDUseV1() {
			if crd, err := GetCrdV1(s); err != nil {
				klog.Warningf("failed to unmarshall crd to v1. because: %v", err)
			} else {
				crdV1s = append(crdV1s, crd)
			}
		}
		if SupportCRDUseV1Beta1() {
			if crd, err := GetCrdV1Beta1(s); err != nil {
				klog.Warningf("failed to unmarshall crd to v1beta1. because: %v", err)
			} else {
				crdV1Beta1s = append(crdV1Beta1s, crd)
			}
		}
	}
	return crdV1s, crdV1Beta1s
}
