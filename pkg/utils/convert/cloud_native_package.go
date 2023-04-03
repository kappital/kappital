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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
)

const (
	operatorDirName  = "operator"
	manifestsDirName = "manifests"
)

var (
	spRequiredDirSet = sets.NewString(manifestsDirName, operatorDirName)

	dirAcceptGVKSet = map[string]sets.String{
		operatorDirName:  sets.NewString(clusterRoleGVK, clusterRoleBindingGVK, deploymentGVK, serviceAccountGVK),
		manifestsDirName: sets.NewString(csdGVK, crdV1GVK, crdV1Beta1GVK),
	}

	metadataFileNameSet = sets.NewString("metadata.yaml", "metadata.json", "metadata.yml")
)

type cloudNativePackage struct {
	metadata apis.Descriptor
	operator svcv1alpha1.OperatorSpec
	manifest manifest

	fileCounter   int
	totalFileSize int64
	rootDir       string
}

type manifest struct {
	CRDs []apis.AbstractResource
	CSDs []svcv1alpha1.CustomServiceDefinition
}

func newCloudNativePackageLoader() Convert {
	return &cloudNativePackage{
		operator: svcv1alpha1.OperatorSpec{
			Deployments:         make([]appsv1.Deployment, 0, initResourceSize),
			ServiceAccounts:     make([]corev1.ServiceAccount, 0, initResourceSize),
			ClusterRoles:        make([]rbacv1.ClusterRole, 0, initResourceSize),
			ClusterRoleBindings: make([]rbacv1.ClusterRoleBinding, 0, initResourceSize),
		},
		manifest: manifest{
			CRDs: make([]apis.AbstractResource, 0, initResourceSize),
			CSDs: make([]svcv1alpha1.CustomServiceDefinition, 0, initResourceSize),
		},
	}
}

func (k *cloudNativePackage) TransferToCloudNativeService(path string) (*svcv1alpha1.CloudNativeService, error) {
	// generate the pass-in path, and get the absolute path
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("cannot get the absolute file path, err: %v", err)
	}
	// set the absolute path to the rootDir path
	k.rootDir = absPath
	// k.rootDir depth is 0, but its info will add one
	fileInfos, err := ioutil.ReadDir(k.rootDir)
	if err != nil {
		return nil, fmt.Errorf("unable to read package %s, err: %s", path, err)
	}
	for _, info := range fileInfos {
		if err = k.fileResolution(info, 1); err != nil {
			if errors.Is(err, filepath.SkipDir) {
				klog.Warningf("the directory is not include the necessary resolution directory, "+
					"skip %s dir during the %s", info.Name(), path)
				continue
			}
			return nil, err
		}
	}
	return k.convert2CloudNativeService()
}

func (k *cloudNativePackage) fileResolution(info os.FileInfo, depth int) error {
	if info.IsDir() {
		if !spRequiredDirSet.Has(info.Name()) {
			return filepath.SkipDir
		}
		return k.processGVKDirectory(info.Name(), depth+1, dirAcceptGVKSet[info.Name()])
	}
	if metadataFileNameSet.Has(info.Name()) {
		return k.processMetadata(filepath.Join(k.rootDir, info.Name()))
	}
	return nil
}

func (k *cloudNativePackage) processMetadata(path string) error {
	// does the file count greater than max file count
	if k.fileCounter+1 > maxFileCount {
		return errOverMaxFileCount
	}
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("error reading file %s, err: %v", path, err)
	}
	// does the single file size is greater than the max single file size
	// or, the total read file size is greater than the max total file size
	if int64(len(raw)) > singleFileSize || k.totalFileSize+int64(len(raw)) > totalFileSize {
		return errOverFileSize
	}
	k.fileCounter++
	k.totalFileSize += int64(len(raw))

	klog.Infof("convert file %s", path)
	if err = yaml.NewYAMLToJSONDecoder(bytes.NewReader(raw)).Decode(&k.metadata); err != nil {
		return fmt.Errorf("cannot decode the metadata file, err: %s", err)
	}
	if !k.metadata.Validation() {
		return fmt.Errorf("invalid metadata file content")
	}
	return nil
}

func (k *cloudNativePackage) processGVKDirectory(name string, depth int, acceptGVKSet sets.String) error {
	if depth > maxDepth {
		return errOverMaxDirDepth
	}
	dir := filepath.Join(k.rootDir, name)
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read package %s, err: %s", dir, err)
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			if err = k.processGVKDirectory(filepath.Join(dir, info.Name()), depth+1, acceptGVKSet); err != nil && !errors.Is(err, errOverMaxDirDepth) {
				return err
			}
			continue
		}
		if err = k.processGVKFile(filepath.Join(dir, info.Name()), depth+1, acceptGVKSet); err != nil {
			return err
		}
	}
	return nil
}

func (k *cloudNativePackage) processGVKFile(path string, depth int, acceptGVKSet sets.String) error {
	if depth > maxDepth {
		return errOverMaxDirDepth
	}
	if !isValidFileType(path) {
		return fmt.Errorf("the file type is not accept")
	}
	if k.fileCounter+1 > maxFileCount {
		return errOverMaxFileCount
	}
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	if int64(len(raw)) > singleFileSize || k.totalFileSize+int64(len(raw)) > totalFileSize {
		return errOverFileSize
	}
	k.fileCounter++
	k.totalFileSize += int64(len(raw))
	klog.Infof("convert file %s", path)
	metaDecoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), decodeBufferSize)
	detailDecoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(raw), decodeBufferSize)
	for {
		meta := metav1.TypeMeta{}
		if err = metaDecoder.Decode(&meta); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error decoding GVK in file %s: %s", path, err)
		}
		if meta.APIVersion == "" || meta.Kind == "" {
			return fmt.Errorf("%s is invalid, GVK should be provided", path)
		}
		currGVK := fmt.Sprintf("%s/%s", meta.APIVersion, meta.Kind)
		if !acceptGVKSet.Has(currGVK) {
			return fmt.Errorf("current GVK %s cannot accept to resolution", currGVK)
		}
		if err = k.processGVK(detailDecoder, currGVK); err != nil {
			return fmt.Errorf("cannot resolution the file %s, err: %v", path, err)
		}
	}
	return nil
}

func (k *cloudNativePackage) processGVK(decoder *yaml.YAMLOrJSONDecoder, gvk string) error {
	var err error
	switch gvk {
	case crdV1Beta1GVK, crdV1GVK:
		crd := apis.AbstractResource{}
		if err = decoder.Decode(&crd); err == nil {
			k.manifest.CRDs = append(k.manifest.CRDs, crd)
		}
	case csdGVK:
		csd := svcv1alpha1.CustomServiceDefinition{}
		if err = decoder.Decode(&csd); err == nil {
			k.manifest.CSDs = append(k.manifest.CSDs, csd)
		}
	case serviceAccountGVK:
		sa := corev1.ServiceAccount{}
		if err = decoder.Decode(&sa); err == nil {
			k.operator.ServiceAccounts = append(k.operator.ServiceAccounts, sa)
		}
	case clusterRoleGVK:
		cr := rbacv1.ClusterRole{}
		if err = decoder.Decode(&cr); err == nil {
			k.operator.ClusterRoles = append(k.operator.ClusterRoles, cr)
		}
	case clusterRoleBindingGVK:
		crb := rbacv1.ClusterRoleBinding{}
		if err = decoder.Decode(&crb); err == nil {
			k.operator.ClusterRoleBindings = append(k.operator.ClusterRoleBindings, crb)
		}
	case deploymentGVK:
		deploy := appsv1.Deployment{}
		if err = decoder.Decode(&deploy); err == nil {
			k.operator.Deployments = append(k.operator.Deployments, deploy)
		}
	default:
		return fmt.Errorf("the GVK %s is not accept", gvk)
	}
	if err != nil {
		return fmt.Errorf("cannot decode into resource %s, err: %s", gvk, err)
	}
	return nil
}

func (k *cloudNativePackage) convert2CloudNativeService() (*svcv1alpha1.CloudNativeService, error) {
	cns := svcv1alpha1.CloudNativeService{}
	cns.Spec.Manifests = make([]svcv1alpha1.CustomServiceDefinition, 0, len(k.manifest.CSDs))
	for _, csd := range k.manifest.CSDs {
		var find bool
		for i, crd := range k.manifest.CRDs {
			if crd.Name == csd.Spec.CRDName {
				csd.Spec.CRD = &k.manifest.CRDs[i]
				cns.Spec.Manifests = append(cns.Spec.Manifests, csd)
				find = true
				break
			}
		}
		if !find {
			return nil, fmt.Errorf("cannot find the crd to the csd %s", csd.Spec.CRDName)
		}
	}
	cns.APIVersion = apis.CloudNativeAPIVersionV1Alpha1
	cns.Kind = apis.CloudNativeServiceKind
	cns.Name = k.metadata.Name
	cns.Spec.Description = k.metadata
	cns.Spec.Operator = &k.operator
	return &cns, nil
}
