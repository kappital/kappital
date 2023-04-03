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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

// reconcileStatefulSet will process the stateful set resources such as create, update, or delete. During
// the create and update  processes, service engine will check does the cluster has the same name stateful set.
// If the desired stateful set exist, but the cluster does not have, engine will create it.
// If the desired stateful set does not exist, but the cluster exist, engine will delete this stateful set in
// cluster.
// If the desired and cluster both have the same name stateful set, engine will upgrade the cluster one to the
// desired.
// In addition, if the service package have the delete signal, will delete the stateful set with the owner reference.
func (r *ServicePackageReconciler) reconcileStatefulSet(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	stsSpecs []enginev1alpha1.ServiceStatefulSetSpec) error {
	if pack.IsDeleting() {
		klog.Infof("get the delete signal, engine will delete all deployments with the owner reference")
		return r.deleteStatefulSet(ctx, pack)
	}
	if len(stsSpecs) == 0 {
		// helm service package with no deployment config, do nothing
		return nil
	}
	stsMap := getStatefulSetMap(stsSpecs)
	if err := r.deleteOrUpdateStatefulSet(ctx, pack, stsMap); err != nil {
		return err
	}
	return r.createOrUpdateStatefulSet(ctx, pack, stsMap)
}

func (r *ServicePackageReconciler) deleteStatefulSet(ctx context.Context, pack *enginev1alpha1.ServicePackage) error {
	stsList, err := r.getStatefulSetList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range stsList {
		item := stsList[i]
		if err := r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete stateful set [%s], because: %s", item.Name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateStatefulSet(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	stsMap map[string]appsv1.StatefulSetSpec) error {
	stsList, err := r.getStatefulSetList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range stsList {
		item := stsList[i]
		spec, find := stsMap[item.Name]
		if !find {
			klog.Infof("cannot find the stateful set [%s] in namespace [%s], "+
				"will delete this stateful set", item.Name, item.Namespace)
			if err := r.Delete(ctx, &item); err != nil {
				klog.Errorf("failed to delete stateful set [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the stateful set [%s]", item.Name)
		} else if pack.IsUpgrading() {
			desired := constructStatefulSet(item.Name, pack.Namespace, spec)
			if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for stateful set [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("find the stateful set [%s] in namespace [%s]", item.Name, item.Namespace)
			if err := r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update stateful set [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade stateful set [%s]", desired.Namespace)
		}
		delete(stsMap, item.Name)
	}
	return nil
}

func (r ServicePackageReconciler) createOrUpdateStatefulSet(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	stsMap map[string]appsv1.StatefulSetSpec) error {
	for name, spec := range stsMap {
		desired := constructStatefulSet(name, pack.Namespace, spec)
		if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
			klog.Errorf("failed to set controller reference for stateful set [%s], because: %s", name, err)
			return err
		}

		tmp := appsv1.StatefulSet{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: pack.Namespace}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the stateful set [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the stateful set [%s] in namespace [%s] is exist, will update this stateful set",
				name, desired.Namespace)
			if err = r.Update(ctx, &desired); err != nil {
				klog.Errorf("cannot update stateful set [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update stateful set [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("stateful set [%s] is not exist in namespace [%s], will create this stateful set",
				name, desired.Namespace)
			if err = r.Create(ctx, &desired); err != nil {
				klog.Errorf("cannot create the stateful set [%s], because: %s", name, err)
			}
			klog.Infof("create stateful set [%s]", name)
		} else {
			klog.Errorf("fail to check the stateful set [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) getStatefulSetList(ctx context.Context,
	name, namespace string) ([]appsv1.StatefulSet, error) {
	stsList := appsv1.StatefulSetList{}
	if err := r.List(ctx, &stsList, client.InNamespace(namespace)); err != nil {
		klog.Errorf("unable to list stateful set, because: %s", err)
		return []appsv1.StatefulSet{}, err
	}

	sts := make([]appsv1.StatefulSet, 0, len(stsList.Items))
	for _, item := range stsList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			sts = append(sts, item)
		}
	}
	return sts, nil
}

// getDeployMap transforms desired deployments from slice to map
func getStatefulSetMap(stsSpecs []enginev1alpha1.ServiceStatefulSetSpec) map[string]appsv1.StatefulSetSpec {
	tmp := make([]enginev1alpha1.ServiceStatefulSetSpec, len(stsSpecs))
	copy(tmp, stsSpecs)
	m := make(map[string]appsv1.StatefulSetSpec, len(tmp))
	for _, d := range tmp {
		m[d.Name] = d.Spec
	}
	return m
}

// constructDeployment constructs deployment
func constructStatefulSet(name, namespace string, stsSpec appsv1.StatefulSetSpec) appsv1.StatefulSet {
	return appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "StatefulSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: stsSpec,
	}
}
