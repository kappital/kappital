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

// reconcileDaemonSet will process the daemon set resources such as create, update, or delete. During
// the create and update  processes, service engine will check does the cluster has the same name daemon set.
// If the desired daemon set exist, but the cluster does not have, engine will create it.
// If the desired daemon set does not exist, but the cluster exist, engine will delete this daemon set in
// cluster.
// If the desired and cluster both have the same name daemon set, engine will upgrade the cluster one to the
// desired.
// In addition, if the service package have the delete signal, will delete the daemon set with the owner reference.
func (r *ServicePackageReconciler) reconcileDaemonSet(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	daemonSetSpecs []enginev1alpha1.ServiceDaemonSetSpec) error {
	if pack.IsDeleting() {
		klog.Infof("get the delete signal, engine will delete all daemon set with the owner reference")
		return r.deleteDaemonSet(ctx, pack)
	}
	if len(daemonSetSpecs) == 0 {
		// helm service package with no deployment config, do nothing
		return nil
	}
	dsMap := getDaemonSetMap(daemonSetSpecs)
	if err := r.deleteOrUpdateDaemonSet(ctx, pack, dsMap); err != nil {
		return err
	}
	return r.createOrUpdateDaemonSet(ctx, pack, dsMap)
}

func (r *ServicePackageReconciler) deleteDaemonSet(ctx context.Context, kappital *enginev1alpha1.ServicePackage) error {
	dsList, err := r.getDaemonSetList(ctx, kappital.Name, kappital.Namespace)
	if err != nil {
		return err
	}
	for i := range dsList {
		item := dsList[i]
		if err := r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete daemon set [%s], because: %s", item.Name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateDaemonSet(ctx context.Context, kappital *enginev1alpha1.ServicePackage,
	dsMap map[string]appsv1.DaemonSetSpec) error {
	dsList, err := r.getDaemonSetList(ctx, kappital.Name, kappital.Namespace)
	if err != nil {
		return err
	}
	for i := range dsList {
		item := dsList[i]
		spec, find := dsMap[item.Name]
		if !find {
			klog.Infof("cannot find the daemon set [%s] in namespace [%s], "+
				"will delete this daemon set", item.Name, item.Namespace)
			if err := r.Delete(ctx, &item); err != nil { //nolint:gosec
				klog.Errorf("failed to delete daemon set [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the daemon set [%s]", item.Name)
		} else if kappital.IsUpgrading() {
			desired := constructDaemonSet(item.Name, kappital.Namespace, spec)
			if err := ctrl.SetControllerReference(kappital, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for daemon set [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("find the daemon set [%s] in namespace [%s]", item.Name, item.Namespace)
			if err := r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update daemon set [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade daemon set [%s]", desired.Namespace)
		}
		delete(dsMap, item.Name)
	}
	return nil
}

func (r ServicePackageReconciler) createOrUpdateDaemonSet(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	dsMap map[string]appsv1.DaemonSetSpec) error {
	for name, spec := range dsMap {
		desired := constructDaemonSet(name, pack.Namespace, spec)
		if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
			klog.Errorf("failed to set controller reference for daemon set [%s], because: %s", name, err)
			return err
		}

		tmp := appsv1.DaemonSet{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: pack.Namespace}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the daemon set [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the daemon set [%s] in namespace [%s] is exist, will update this daemon set",
				name, desired.Namespace)
			if err = r.Update(ctx, &desired); err != nil {
				klog.Errorf("cannot update daemon set [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update daemon set [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("daemon set [%s] is not exist in namespace [%s], will create this daemon set",
				name, desired.Namespace)
			if err = r.Create(ctx, &desired); err != nil {
				klog.Errorf("cannot create the daemon set [%s], because: %s", name, err)
			}
			klog.Infof("create daemon set [%s]", name)
		} else {
			klog.Errorf("fail to check the daemon set [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) getDaemonSetList(ctx context.Context, name, namespace string) ([]appsv1.DaemonSet, error) {
	dsList := appsv1.DaemonSetList{}
	if err := r.List(ctx, &dsList, client.InNamespace(namespace)); err != nil {
		klog.Errorf("unable to list daemon set, because: %s", err)
		return []appsv1.DaemonSet{}, err
	}

	ds := make([]appsv1.DaemonSet, 0, len(dsList.Items))
	for _, item := range dsList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			ds = append(ds, item)
		}
	}
	return ds, nil
}

// getDaemonSetMap transforms desired daemon set from slice to map
func getDaemonSetMap(daemonSetSpecs []enginev1alpha1.ServiceDaemonSetSpec) map[string]appsv1.DaemonSetSpec {
	tmp := make([]enginev1alpha1.ServiceDaemonSetSpec, len(daemonSetSpecs))
	copy(tmp, daemonSetSpecs)
	m := make(map[string]appsv1.DaemonSetSpec, len(tmp))
	for _, d := range tmp {
		m[d.Name] = d.Spec
	}
	return m
}

// constructDeployment constructs deployment
func constructDaemonSet(name, namespace string, daemonSetSpec appsv1.DaemonSetSpec) appsv1.DaemonSet {
	return appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "DaemonSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: daemonSetSpec,
	}
}
