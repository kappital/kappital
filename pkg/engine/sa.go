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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

// reconcileServiceAccount will process the service account resources such as create, update, or delete. During
// the create and update  processes, service engine will check does the cluster has the same name service account.
// If the desired service account exist, but the cluster does not have, engine will create it.
// If the desired service account does not exist, but the cluster exist, engine will delete this service account in
// cluster.
// If the desired and cluster both have the same name service account, engine will upgrade the cluster one to the
// desired.
// In addition, if the service package have the delete signal, will delete the service account with the owner reference.
func (r *ServicePackageReconciler) reconcileServiceAccount(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	permissions []enginev1alpha1.Permission) error {
	if pack.IsDeleting() {
		klog.Infof("get the delete signal, engine will delete all service account with the owner reference")
		return r.deleteServiceAccount(ctx, pack)
	}

	saMap := getServiceAccountMap(permissions, pack.Namespace)
	if err := r.deleteOrUpdateServiceAccount(ctx, pack, saMap); err != nil {
		return err
	}
	return r.createOrUpdateServiceAccount(ctx, pack, saMap)
}

func (r *ServicePackageReconciler) deleteServiceAccount(ctx context.Context, pack *enginev1alpha1.ServicePackage) error {
	saList, err := r.getServiceAccountList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range saList {
		item := saList[i]
		if err = r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete service account [%s], because: %s", item.Name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateServiceAccount(ctx context.Context,
	pack *enginev1alpha1.ServicePackage, saMap map[string]corev1.ServiceAccount) error {
	saList, err := r.getServiceAccountList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range saList {
		item := saList[i]
		desired, find := saMap[item.Name]
		if !find {
			klog.Infof("cannot find the service account [%s] in namespace [%s], "+
				"will delete this service account", item.Name, item.Namespace)
			if err = r.Delete(ctx, &item); err != nil {
				klog.Errorf("failed to delete service account [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the service account [%s]", item.Name)
		} else if pack.IsUpgrading() {
			if err = ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for service account [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("find the service account [%s] in namespace [%s]", item.Name, item.Namespace)
			if err = r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update service account [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade service account [%s]", desired.Namespace)
		}
		delete(saMap, item.Name)
	}
	return nil
}

func (r *ServicePackageReconciler) createOrUpdateServiceAccount(ctx context.Context,
	pack *enginev1alpha1.ServicePackage, saMap map[string]corev1.ServiceAccount) error {
	for name, item := range saMap {
		desired := item.DeepCopy()
		if err := ctrl.SetControllerReference(pack, desired, r.Scheme); err != nil {
			klog.Errorf("failed to set controller reference for service account [%s], because: %s", name, err)
			return err
		}

		tmp := corev1.ServiceAccount{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: pack.Namespace}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the service account [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the service account [%s] in namespace [%s] is exist, will update this service account",
				name, desired.Namespace)
			if err = r.Update(ctx, desired); err != nil {
				klog.Errorf("cannot update service account [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update service account [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("service account [%s] is not exist in namespace [%s], will create this service account",
				name, desired.Namespace)
			if err = r.Create(ctx, desired); err != nil {
				klog.Errorf("cannot create the service account [%s], because: %s", name, err)
			}
			klog.Infof("create service account [%s]", name)
		} else {
			klog.Errorf("fail to check the service account [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) getServiceAccountList(ctx context.Context,
	name, namespace string) ([]corev1.ServiceAccount, error) {
	saList := corev1.ServiceAccountList{}
	if err := r.List(ctx, &saList, client.InNamespace(namespace)); err != nil {
		klog.Errorf("unable to list service account, because: %s", err)
		return []corev1.ServiceAccount{}, err
	}
	services := make([]corev1.ServiceAccount, 0, len(saList.Items))
	for _, item := range saList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			services = append(services, item)
		}
	}
	return services, nil
}

func getServiceAccountMap(permissions []enginev1alpha1.Permission, namespace string) map[string]corev1.ServiceAccount {
	saMap := make(map[string]corev1.ServiceAccount, len(permissions))
	for _, permission := range permissions {
		saMap[permission.ServiceAccountName] = corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      permission.ServiceAccountName,
				Namespace: namespace,
			},
		}
	}
	return saMap
}
