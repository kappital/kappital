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
	"encoding/json"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
)

var (
	testSA = corev1.ServiceAccount{
		TypeMeta:   metav1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-sa"},
	}
	testCR = rbacv1.ClusterRole{
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterRole", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-cr"},
	}
	testCRB = rbacv1.ClusterRoleBinding{
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-crb"},
	}
	testDeploy = appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-deploy"},
	}
	testCSD = svcv1alpha1.CustomServiceDefinition{
		TypeMeta:   metav1.TypeMeta{Kind: "CustomServiceDefinition", APIVersion: "core.kappital.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-csd"},
	}
	testCRDV1 = apiextensionsv1.CustomResourceDefinition{
		TypeMeta:   metav1.TypeMeta{Kind: "CustomResourceDefinition", APIVersion: "apiextensions.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-crd-v1"},
	}
	testCRDV1Beta1 = apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta:   metav1.TypeMeta{Kind: "CustomResourceDefinition", APIVersion: "apiextensions.k8s.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-crd-v1beta1"},
	}
	testEmptyGVK      = apis.AbstractResource{ObjectMeta: metav1.ObjectMeta{Name: "test-empty"}}
	testNotIncludeGVK = appsv1.StatefulSet{
		TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-sts"},
	}

	contentMap = map[string][]byte{}

	testMetadatas = []string{`name: kappital-demo
version: 0.1.0
type: operator
minKubeVersion: 1.15.0
briefDescription: example package with an example operator and instance`,
		`name: xx1234567890123456789012345678901234567890123456789012345678901234567890
version: 0.1.0
type: operator
minKubeVersion: 1.15.0
briefDescription: example package with an example operator and instance`,
	}
)

func TestMain(m *testing.M) {
	createUnitTestJSONFiles()
	m.Run()
}

func createUnitTestJSONFiles() {
	var err error
	contentSA, err := json.MarshalIndent(testSA, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentCR, err := json.MarshalIndent(testCR, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentCRB, err := json.MarshalIndent(testCRB, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentDeploy, err := json.MarshalIndent(testDeploy, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentCSD, err := json.MarshalIndent(testCSD, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentCRDV1, err := json.MarshalIndent(testCRDV1, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentCRDV1Beta1, err := json.MarshalIndent(testCRDV1Beta1, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentIllegal, err := json.MarshalIndent(testEmptyGVK, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentNotInclude, err := json.MarshalIndent(testNotIncludeGVK, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	contentMap = map[string][]byte{
		testSA.Name:            contentSA,
		testCR.Name:            contentCR,
		testCRB.Name:           contentCRB,
		testDeploy.Name:        contentDeploy,
		testCSD.Name:           contentCSD,
		testCRDV1.Name:         contentCRDV1,
		testCRDV1Beta1.Name:    contentCRDV1Beta1,
		testEmptyGVK.Name:      contentIllegal,
		testNotIncludeGVK.Name: contentNotInclude,
	}
}
