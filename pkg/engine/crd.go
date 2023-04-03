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

package engine

import (
	"context"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	"github.com/kappital/kappital/pkg/utils/version"
)

// reconcileCRD will process the custom resource definition resources. It will only try to create the CRD.
// In other words, if the cluster have the same name CRD, it will not create or update crd.
func (r *ServicePackageReconciler) reconcileCRD(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	crdStrings []string) error {
	if pack.IsDeleting() {
		klog.Warningf("the service %s is deleting or already deleted. the service crd will not delete.", pack.Name)
		return nil
	}
	if len(crdStrings) == 0 {
		klog.Warning("no crd defined, ignore reconcile crd")
		return nil
	}
	klog.Infof("reconcile crd, has %d crds", len(crdStrings))
	crdV1s, crdV1Beta1s := version.GetCrdV1AndBeta1Slice(crdStrings)
	errV1, errV1Beta1 := r.createCRDV1(ctx, crdV1s), r.createCRDV1Beta1(ctx, crdV1Beta1s)
	return dealWithErrorCount(errV1, errV1Beta1, len(crdV1s), len(crdV1Beta1s))
}

func (r *ServicePackageReconciler) createCRDV1(ctx context.Context, crds []apiextensionsv1.CustomResourceDefinition) int {
	if len(crds) != 0 && !version.SupportCRDUseV1() {
		klog.Warning("the crd slice is empty or this cluster not support for apiextensionsv1 crd")
		return 0
	}
	errCount := 0
	for i := range crds {
		crd := crds[i]
		if err := r.Create(ctx, &crd); err != nil {
			klog.Warningf("cannot create crd(v1) [%s] because %s", crd.Name, err)
			errCount++
		}
		klog.Infof("create crd [%s] for v1", crd.Name)
	}
	return errCount
}

func (r *ServicePackageReconciler) createCRDV1Beta1(ctx context.Context, crds []apiextensionsv1beta1.CustomResourceDefinition) int {
	if len(crds) != 0 && !version.SupportCRDUseV1Beta1() {
		klog.Warning("the crd slice is empty or this cluster not support for v1beta1 crd")
		return 0
	}
	errCount := 0
	for i := range crds {
		crd := crds[i]
		if err := r.Create(ctx, &crd); err != nil {
			klog.Warningf("cannot create crd (v1beta1) [%s] because %s", crd.Name, err)
			errCount++
		}
		klog.Infof("create crd [%s] for v1beta1", crd.Name)
	}
	return errCount
}

func dealWithErrorCount(errCountV1, errCountV1Beta1, totalV1, totalV1Beta1 int) error {
	if errCountV1 == 0 && errCountV1Beta1 == 0 {
		return nil
	}
	if errCountV1 == totalV1 && errCountV1Beta1 == totalV1Beta1 {
		return fmt.Errorf("all crds developed failed")
	}
	if errCountV1 <= totalV1 && errCountV1Beta1 <= totalV1Beta1 {
		return fmt.Errorf("partial crds developed failed")
	}
	return nil
}
