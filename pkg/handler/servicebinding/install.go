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

package servicebinding

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/apis/internals"
	co "github.com/kappital/kappital/pkg/utils/operations"
)

const bindingProcessTimeout = 3 * time.Minute

// Handler singleton pattern of install, delete, or upgrade the service binding
type Handler struct{}

// BeforeInstall do some processes before install the service binding
func (h *Handler) BeforeInstall(_ interface{}) (bool, error) {
	return false, nil
}

// Install the cloud native service instance in cluster
func (h *Handler) Install(obj interface{}) (retry bool, err error) {
	serviceBinding := getTypedObj(obj)
	// renew the timeout period at the first time
	if err = updateProcessTimeout(serviceBinding, bindingProcessTimeout); err != nil {
		klog.Errorf("update process time for binding %s error: %v", serviceBinding.Name, err)
		return true, err
	}

	exist, ready, er := h.checkBindingStatus(serviceBinding)
	if er != nil {
		klog.Errorf("failed to check binding status, error: %v", er)
		return true, er
	}
	if ready {
		// binding ready, exist the asynchronizing and will not retry
		klog.Infof("[install binding] binding %s is already succeed, stop retry", serviceBinding.Name)
		return false, nil
	}
	if exist {
		return true, fmt.Errorf("[install binding] binding is already exist but not succeed, rechecking")
	}

	err = h.createServiceBinding(serviceBinding)
	if err != nil {
		klog.Errorf("failed to create binding servicepackage, error: %v", err)
		return true, err
	}

	klog.Infof("[install binding] binding %s servicepackage is created, check for next loop", serviceBinding.Name)
	return true, nil
}

// AfterInstall install service binding does not need to implement this method
func (h *Handler) AfterInstall(_ interface{}) (bool, error) {
	return false, nil
}

func (h *Handler) createServiceBinding(binding *internals.ServiceBinding) error {
	servicePackage := enginev1alpha1.ServicePackage{
		TypeMeta: metav1.TypeMeta{
			Kind:       enginev1alpha1.ServicePackageKind,
			APIVersion: enginev1alpha1.ServicePackageAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      binding.Name,
			Namespace: apis.KappitalSystemNamespace,
		},
		Spec: enginev1alpha1.ServicePackageSpec{
			ServiceID: binding.ServiceID,
			Name:      binding.ServiceName,
			Version:   binding.Version,
		},
	}

	serviceResource := enginev1alpha1.ServiceResource{
		CustomResourceDefinitions: binding.CRD,
		Permissions:               binding.Permissions,
		CapabilityPlugin:          binding.CapabilityPlugin,
		Workload:                  binding.Workload,
	}

	var err error
	servicePackage.Spec.Resources, err = enginev1alpha1.TranslateResourcesToBase64(serviceResource)
	if err != nil {
		klog.Errorf("trans binding %s service resource to base64 failed", binding.Name)
		return err
	}

	if err = co.GetClusterOperation().DeployCustomResource(enginev1alpha1.ServicePackageGroupVersionResource,
		apis.KappitalSystemNamespace, servicePackage); err != nil {
		klog.Errorf("create service binding %s failed.", binding.Name)
	}

	return err
}

func (h *Handler) checkBindingStatus(binding *internals.ServiceBinding) (exist, ready bool, err error) {
	subExist, ready := checkBindingReady(binding)
	if ready {
		if err = updateSuccessStatus(binding); err != nil {
			klog.Errorf("[binding handler] update binding status failed, error: %v", err)
			return true, true, err
		}
	}
	return subExist, ready, nil
}

// checkBindingReady return retry and errMsg
func checkBindingReady(binding *internals.ServiceBinding) (bool, bool) {
	sp, found, err := co.GetClusterOperation().GetServicePackageByName(binding.Name, apis.KappitalSystemNamespace)
	if err != nil {
		return true, false
	}

	if !found {
		return false, false
	}

	return true, sp.Status.Phase == enginev1alpha1.RunningPhase
}
