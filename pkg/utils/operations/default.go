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
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

type defaultOperation struct{}

// GetServicePackageByName get the service package CR by its name. This method will return the ServicePackage
// if existed.
func (d *defaultOperation) GetServicePackageByName(name, namespace string) (enginev1alpha1.ServicePackage, bool, error) {
	config, err := getConfig()
	if err != nil {
		klog.Errorf("cannot get the client config, err: %v", err)
		return enginev1alpha1.ServicePackage{}, false, err
	}
	config.APIPath = "apis"
	config.GroupVersion = &enginev1alpha1.GroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	cli, err := rest.RESTClientFor(config)
	if err != nil {
		klog.Errorf("cannot get the client, err: %v", err)
		return enginev1alpha1.ServicePackage{}, false, err
	}
	sp := enginev1alpha1.ServicePackage{}
	if err = cli.Get().Resource(servicePackageResource).Namespace(namespace).Name(name).Do(context.TODO()).
		Into(&sp); err != nil {
		if errors.IsNotFound(err) {
			klog.Warning("cannot find the service package")
			return enginev1alpha1.ServicePackage{}, false, nil
		}
		klog.Errorf("cannot find the service package, err: %v", err)
		return enginev1alpha1.ServicePackage{}, false, err
	}
	return sp, true, nil
}

// DoesCustomResourceExist will use resource's schema.GroupVersion to find the cr in this cluster
func (d *defaultOperation) DoesCustomResourceExist(gv schema.GroupVersion,
	plural, name, namespace string) (bool, error) {
	config, err := getConfig()
	if err != nil {
		klog.Errorf("cannot get the client config, err: %v", err)
		return false, err
	}
	config.APIPath = "apis"
	config.GroupVersion = &gv
	config.NegotiatedSerializer = scheme.Codecs
	cli, err := rest.RESTClientFor(config)
	if err != nil {
		klog.Errorf("cannot get the client, err: %v", err)
		return false, err
	}
	err = cli.Get().Resource(plural).Name(name).Namespace(namespace).Do(context.TODO()).Error()
	if errors.IsNotFound(err) {
		return false, nil
	}
	return err == nil, err
}

// DeployCustomResource will install the custom resource into cluster
func (d *defaultOperation) DeployCustomResource(gvr schema.GroupVersionResource, namespace string,
	resource interface{}) error {
	cli, obj, err := getCRClientAndObj(resource)
	if err != nil {
		return err
	}
	// create the custom resource
	_, err = cli.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

// UpdateCustomResource will update the custom resource into cluster
func (d *defaultOperation) UpdateCustomResource(gvr schema.GroupVersionResource, namespace string,
	resource interface{}) error {
	cli, obj, err := getCRClientAndObj(resource)
	if err != nil {
		return err
	}
	// update the custom resource
	_, err = cli.Resource(gvr).Namespace(namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
	return err
}

// DeleteCustomResource will delete the custom resource from cluster
func (d *defaultOperation) DeleteCustomResource(gvr schema.GroupVersionResource, name, namespace string) error {
	cli, err := getCustomResourceClient()
	if err != nil {
		return err
	}
	err = cli.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

// IsNamespaceExist will check cluster namespace is exist, true means is exists
func (d *defaultOperation) IsNamespaceExist(namespace string) (bool, error) {
	config, err := getConfig()
	if err != nil {
		klog.Errorf("cannot get the client config, err: %v", err)
		return false, err
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorf("cannot get the client, err: %v", err)
		return false, err
	}

	_, err = cli.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	switch {
	case errors.IsNotFound(err):
		return false, nil
	case err != nil:
		klog.Warningf("cannot find the namespace %s", namespace)
		return false, nil
	}

	return true, nil
}

func getCRClientAndObj(resource interface{}) (dynamic.Interface, *unstructured.Unstructured, error) {
	cli, err := getCustomResourceClient()
	if err != nil {
		return nil, nil, err
	}
	obj, err := getUnstructuredObject(resource)
	if err != nil {
		return nil, nil, err
	}
	return cli, obj, nil
}

func getCustomResourceClient() (dynamic.Interface, error) {
	config, err := getConfig()
	if err != nil {
		klog.Errorf("cannot get the client config, err: %v", err)
		return nil, err
	}
	cli, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Errorf("cannot get the client, err: %v", err)
		return nil, err
	}
	return cli, err
}

func getUnstructuredObject(resource interface{}) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	var err error
	// if the resource type is string, need to translate string to runtime.Object and then translate to the
	// unstructured.Object
	value, ok := resource.(string)
	if ok {
		// translate the string to runtime.Object
		s, _, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode([]byte(value),
			nil, &unstructured.Unstructured{})
		if err != nil {
			return nil, err
		}
		obj.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(s)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}

	obj.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&resource)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
