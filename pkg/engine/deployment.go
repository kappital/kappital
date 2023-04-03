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

// reconcileDeployment will process the deployment resources such as create, update, or delete. During
// the create and update  processes, service engine will check does the cluster has the same name deployment.
// If the desired deployment exist, but the cluster does not have, engine will create it.
// If the desired deployment does not exist, but the cluster exist, engine will delete this deployment in
// cluster.
// If the desired and cluster both have the same name deployment, engine will upgrade the cluster one to the
// desired.
// In addition, if the service package have the delete signal, will delete the deployment with the owner reference.
func (r *ServicePackageReconciler) reconcileDeployment(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	deploySpecs []enginev1alpha1.ServiceDeploymentSpec) error {
	if pack.IsDeleting() {
		klog.Infof("get the delete signal, engine will delete all deployments with the owner reference")
		return r.deleteDeployment(ctx, pack)
	}
	if len(deploySpecs) == 0 {
		// helm service package with no deployment config, do nothing
		return nil
	}
	deployMap := getDeployMap(deploySpecs)
	if err := r.deleteOrUpdateDeployment(ctx, pack, deployMap); err != nil {
		return err
	}
	return r.createOrUpdateDeployment(ctx, pack, deployMap)
}

func (r *ServicePackageReconciler) deleteDeployment(ctx context.Context, pack *enginev1alpha1.ServicePackage) error {
	deployList, err := r.getDeploymentList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range deployList {
		item := deployList[i]
		if err := r.Delete(ctx, &item); err != nil {
			klog.Errorf("failed to delete deployments [%s], because: %s", item.Name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) deleteOrUpdateDeployment(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	deployMap map[string]appsv1.DeploymentSpec) error {
	deployList, err := r.getDeploymentList(ctx, pack.Name, pack.Namespace)
	if err != nil {
		return err
	}
	for i := range deployList {
		item := deployList[i]
		spec, find := deployMap[item.Name]
		if !find {
			klog.Infof("cannot find the deployment [%s] in namespace [%s], "+
				"will delete this deployment", item.Name, item.Namespace)
			if err := r.Delete(ctx, &item); err != nil {
				klog.Errorf("failed to delete deployment [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("delete the deployment [%s]", item.Name)
		} else if pack.IsUpgrading() {
			desired := constructDeployment(item.Name, pack.Namespace, spec)
			if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
				klog.Errorf("failed to set controller reference for deployment [%s], because: %s", item.Name, err)
				return err
			}
			klog.Infof("find the deployment [%s] in namespace [%s]", item.Name, item.Namespace)
			if err := r.Update(ctx, &desired); err != nil {
				klog.Errorf("failed to update deployment [%s], because: %s", desired.Name, err)
				return err
			}
			klog.Infof("upgrade deployment [%s]", desired.Namespace)
		}
		delete(deployMap, item.Name)
	}
	return nil
}

func (r ServicePackageReconciler) createOrUpdateDeployment(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	deployMap map[string]appsv1.DeploymentSpec) error {
	for name, spec := range deployMap {
		desired := constructDeployment(name, pack.Namespace, spec)
		if err := ctrl.SetControllerReference(pack, &desired, r.Scheme); err != nil {
			klog.Errorf("failed to set controller reference for deployment [%s], because: %s", name, err)
			return err
		}
		tmp := appsv1.Deployment{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: pack.Namespace}, &tmp); err == nil {
			if !pack.IsUpgrading() {
				klog.Infof("the deployment [%s] is already exist, not need update", name)
				continue
			}
			klog.Infof("the deployment [%s] in namespace [%s] is exist, will update this deployment",
				name, desired.Namespace)
			if err = r.Update(ctx, &desired); err != nil {
				klog.Errorf("cannot update deployment [%s], because: %s", name, err)
				return err
			}
			klog.Infof("update deployment [%s]", name)
		} else if errors.IsNotFound(err) {
			klog.Infof("deployment [%s] is not exist in namespace [%s], will create this deployment",
				name, desired.Namespace)
			if err = r.Create(ctx, &desired); err != nil {
				klog.Errorf("cannot create the deployment [%s], because: %s", name, err)
			}
			klog.Infof("create deployment [%s]", name)
		} else {
			klog.Errorf("fail to check the deployment [%s] in cluster, because: %s", name, err)
			return err
		}
	}
	return nil
}

func (r *ServicePackageReconciler) getDeploymentList(ctx context.Context,
	name, namespace string) ([]appsv1.Deployment, error) {
	deployList := appsv1.DeploymentList{}
	if err := r.List(ctx, &deployList, client.InNamespace(namespace)); err != nil {
		klog.Errorf("unable to list deployment, because: %s", err)
		return []appsv1.Deployment{}, err
	}

	deploys := make([]appsv1.Deployment, 0, len(deployList.Items))
	for _, item := range deployList.Items {
		if belongToOwnerReference(name, item.OwnerReferences) {
			deploys = append(deploys, item)
		}
	}
	return deploys, nil
}

// getDeployMap transforms desired deployments from slice to map
func getDeployMap(deploySpecs []enginev1alpha1.ServiceDeploymentSpec) map[string]appsv1.DeploymentSpec {
	tmp := make([]enginev1alpha1.ServiceDeploymentSpec, len(deploySpecs))
	copy(tmp, deploySpecs)
	m := make(map[string]appsv1.DeploymentSpec, len(tmp))
	for _, d := range tmp {
		m[d.Name] = d.Spec
	}
	return m
}

// constructDeployment constructs deployment
func constructDeployment(name, namespace string, deploySpec appsv1.DeploymentSpec) appsv1.Deployment {
	return appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: deploySpec,
	}
}
