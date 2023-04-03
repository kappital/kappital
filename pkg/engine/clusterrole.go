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

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

// reconcileClusterRoleAndBinding will process the cluster role and cluster role binding resources
// such as create, update, or delete. During the create and update  processes, service engine will check does
// the cluster has the same name cluster role and cluster role binding.
// If the desired cluster role and cluster role binding exist, but the cluster does not have, engine will create it.
// If the desired cluster role and cluster role binding does not exist, but the cluster exist, engine will delete this
// cluster role and cluster role binding in cluster.
// If the desired and cluster both have the same name cluster role and cluster role binding, engine will upgrade
// the cluster one to the desired.
// In addition, if the service package have the delete signal, will delete the cluster role and cluster role binding
// with the owner reference.
func (r *ServicePackageReconciler) reconcileClusterRoleAndBinding(ctx context.Context,
	pack *enginev1alpha1.ServicePackage, permissions []enginev1alpha1.Permission) error {
	if pack.IsDeleting() {
		klog.Infof("get the delete signal, engine will delete all cluster role and cluster role binding with " +
			"the owner reference")
		return r.deleteClusterRoleAndBinding(ctx, pack)
	}

	crMap, crbMap := getClusterRoleAndBindingMap(permissions, pack.Name, pack.Namespace)
	// compare desired cluster roles in bundle with roles in cluster
	if err := r.deleteOrUpdateClusterRole(ctx, pack, crMap); err != nil {
		return err
	}
	// create new cluster role
	if err := r.createOrUpdateClusterRole(ctx, pack, crMap); err != nil {
		return err
	}
	// compare desired roles in bundle with roles in cluster
	if err := r.deleteOrUpdateClusterRoleBinding(ctx, pack, crbMap); err != nil {
		return err
	}
	// create new cluster role binding
	return r.createOrUpdateClusterRoleBinding(ctx, pack, crbMap)
}

func (r *ServicePackageReconciler) deleteClusterRoleAndBinding(ctx context.Context,
	pack *enginev1alpha1.ServicePackage) error {
	crList, err := r.getClusterRoleList(ctx, pack.Name)
	if err != nil {
		return err
	}
	for i := range crList {
		item := crList[i]
		if err = r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete cluster role [%s], because: %s", item.Name, err)
			return err
		}
	}
	crbList, err := r.getClusterRoleBindingList(ctx, pack.Name)
	if err != nil {
		return err
	}
	for i := range crbList {
		item := crList[i]
		if err = r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete cluster role binding [%s], because: %s", item.Name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateClusterRole(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	crMap map[string]rbacv1.ClusterRole) error {
	crList, err := r.getClusterRoleList(ctx, pack.Name)
	if err != nil {
		return err
	}
	for i := range crList {
		item := crList[i]
		desired, find := crMap[item.Name]
		if !find {
			klog.Infof("cannot find the cluster role [%s] in namespace [%s], "+
				"will delete this cluster role", item.Name, item.Namespace)
			if err := r.Delete(ctx, &item); err != nil {
				klog.Errorf("failed to delete cluster role [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the cluster role [%s]", item.Name)
		} else if pack.IsUpgrading() {
			if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for cluster role [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("find the cluster role [%s] in namespace [%s]", item.Name, item.Namespace)
			if err := r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update cluster role [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade cluster role [%s]", desired.Namespace)
		}
		delete(crMap, item.Name)
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateClusterRoleBinding(ctx context.Context,
	pack *enginev1alpha1.ServicePackage, crMap map[string]rbacv1.ClusterRoleBinding) error {
	crbList, err := r.getClusterRoleBindingList(ctx, pack.Name)
	if err != nil {
		return err
	}
	for i := range crbList {
		item := crbList[i]
		desired, find := crMap[item.Name]
		if !find {
			klog.Infof("cannot find the cluster role binding [%s] in namespace [%s], "+
				"will delete this cluster role binding", item.Name, item.Namespace)
			if err := r.Delete(ctx, &item); err != nil {
				klog.Errorf("failed to delete cluster role binding [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the cluster role binding [%s]", item.Name)
		} else if pack.IsUpgrading() {
			if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for cluster role binding [%s], because: %s",
					item.Name, err)
				return err
			}
			klog.Infof("find the cluster role binding [%s] in namespace [%s]", item.Name, item.Namespace)
			if err := r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update cluster role binding [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade cluster role binding [%s]", desired.Namespace)
		}
		delete(crMap, item.Name)
	}
	return nil
}

func (r *ServicePackageReconciler) createOrUpdateClusterRole(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	crMap map[string]rbacv1.ClusterRole) error {
	for name, item := range crMap {
		desired := item.DeepCopy()
		if err := ctrl.SetControllerReference(pack, desired, r.Scheme); err != nil { //nolint:gosec
			klog.Errorf("failed to set controller reference for cluster role [%s], because: %s", name, err)
			return err
		}

		tmp := rbacv1.ClusterRole{}
		if err := r.Get(ctx, types.NamespacedName{Name: name}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the cluster role [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the cluster role [%s] in namespace [%s] is exist, will update this cluster role",
				name, desired.Namespace)
			if err = r.Update(ctx, desired); err != nil {
				klog.Errorf("cannot update cluster role [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update cluster role [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("cluster role [%s] is not exist in namespace [%s], will create this cluster role",
				name, desired.Namespace)
			if err = r.Create(ctx, desired); err != nil {
				klog.Errorf("cannot create the cluster role [%s], because: %s", name, err)
			}
			klog.Infof("create cluster role [%s]", name)
		} else {
			klog.Errorf("fail to check the cluster role [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) createOrUpdateClusterRoleBinding(ctx context.Context,
	pack *enginev1alpha1.ServicePackage, crbMap map[string]rbacv1.ClusterRoleBinding) error {
	for name, item := range crbMap {
		desired := item.DeepCopy()
		if err := ctrl.SetControllerReference(pack, desired, r.Scheme); err != nil {
			klog.Errorf("failed to set controller reference for cluster role binding [%s], because: %s", name, err)
			return err
		}

		tmp := rbacv1.ClusterRoleBinding{}
		if err := r.Get(ctx, types.NamespacedName{Name: name}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the cluster role binding [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the cluster role binding [%s] in namespace [%s] is exist, will update this cluster "+
				"role binding", name, desired.Namespace)
			if err = r.Update(ctx, desired); err != nil {
				klog.Errorf("cannot update cluster role binding [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update cluster role binding [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("cluster role binding [%s] is not exist in namespace [%s], will create this cluster "+
				"role binding", name, desired.Namespace)
			if err = r.Create(ctx, desired); err != nil {
				klog.Errorf("cannot create the cluster role binding [%s], because: %s", name, err)
			}
			klog.Infof("create cluster role binding [%s]", name)
		} else {
			klog.Errorf("fail to check the cluster role binding [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) getClusterRoleList(ctx context.Context, name string) ([]rbacv1.ClusterRole, error) {
	crList := rbacv1.ClusterRoleList{}
	if err := r.List(ctx, &crList); err != nil {
		klog.Errorf("unable to list cluster role, because: %s", err)
		return []rbacv1.ClusterRole{}, err
	}

	crs := make([]rbacv1.ClusterRole, 0, len(crList.Items))
	for _, item := range crList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			crs = append(crs, item)
		}
	}
	return crs, nil
}

func (r *ServicePackageReconciler) getClusterRoleBindingList(ctx context.Context,
	name string) ([]rbacv1.ClusterRoleBinding, error) {
	crbList := rbacv1.ClusterRoleBindingList{}
	if err := r.List(ctx, &crbList); err != nil {
		klog.Errorf("unable to list cluster role binding, because: %s", err)
		return []rbacv1.ClusterRoleBinding{}, err
	}

	crbs := make([]rbacv1.ClusterRoleBinding, 0, len(crbList.Items))
	for _, item := range crbList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			crbs = append(crbs, item)
		}
	}
	return crbs, nil
}

func getClusterRoleAndBindingMap(permissions []enginev1alpha1.Permission,
	name, namespace string) (map[string]rbacv1.ClusterRole, map[string]rbacv1.ClusterRoleBinding) {
	crMap := make(map[string]rbacv1.ClusterRole, len(permissions))
	crbMap := make(map[string]rbacv1.ClusterRoleBinding, len(permissions))
	for _, permission := range permissions {
		crName := fmt.Sprintf("%s-cr-%s", name, permission.ServiceAccountName)
		crbName := fmt.Sprintf("%s-crb-%s", name, permission.ServiceAccountName)
		crMap[crName] = rbacv1.ClusterRole{
			TypeMeta:   metav1.TypeMeta{Kind: "ClusterRole", APIVersion: rbacv1.SchemeGroupVersion.String()},
			ObjectMeta: metav1.ObjectMeta{Name: crName, Namespace: namespace},
			Rules:      permission.Rules,
		}
		crbMap[crbName] = rbacv1.ClusterRoleBinding{
			TypeMeta:   metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: rbacv1.SchemeGroupVersion.String()},
			ObjectMeta: metav1.ObjectMeta{Name: crbName, Namespace: namespace},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					APIGroup:  corev1.SchemeGroupVersion.Group,
					Name:      permission.ServiceAccountName,
					Namespace: namespace,
				},
			},
			RoleRef: rbacv1.RoleRef{Kind: "ClusterRole", APIGroup: rbacv1.SchemeGroupVersion.Group, Name: crName},
		}
	}
	return crMap, crbMap
}
