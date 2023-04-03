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

package operations

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

const (
	// servicePackageResource resource name
	servicePackageResource = "servicepackages"
)

var kubeConfigPath = ""

var operation ClusterOperation

// ClusterOperation the interface of all cluster actions' operation
type ClusterOperation interface {
	// GetServicePackageByName get the service package CR by its name. This method will return the ServicePackage
	// if existed.
	GetServicePackageByName(name, namespace string) (enginev1alpha1.ServicePackage, bool, error)
	// DoesCustomResourceExist will use resource's schema.GroupVersion to find the cr in this cluster
	DoesCustomResourceExist(gv schema.GroupVersion, plural, name, namespace string) (bool, error)
	// DeployCustomResource will install the custom resource into cluster
	DeployCustomResource(gvr schema.GroupVersionResource, namespace string, resource interface{}) error
	// UpdateCustomResource will update the custom resource into cluster
	UpdateCustomResource(gvr schema.GroupVersionResource, namespace string, resource interface{}) error
	// DeleteCustomResource will delete the custom resource from cluster
	DeleteCustomResource(gvr schema.GroupVersionResource, name, namespace string) error
	// IsNamespaceExist will check cluster namespace is existed, true means is exists
	IsNamespaceExist(namespace string) (bool, error)
}

// GetClusterOperation return an ClusterOperation of this interface
func GetClusterOperation() ClusterOperation {
	return operation
}

// SetClusterOperation set up a new ClusterOperation
func SetClusterOperation(o ClusterOperation) {
	operation = o
}

// init defaultOperation will be used at beginning
func init() {
	operation = &defaultOperation{}
	kubeConfigPath = os.Getenv("KubeConfig")
}

// getConfig will try to get the rest.Config from the kubeconfig first. If kappital cannot get the config from kubeconfig, it
// will try to get the rest.Config from the service account. If both ways cannot get the correct the rest.Config, will
// return a non-nil error.
func getConfig() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err == nil {
		return config, nil
	}
	config, err = ctrl.GetConfig()
	return config, err
}
