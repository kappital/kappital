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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"

	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

const (
	// DaemonSetService name of the ds resources
	DaemonSetService = "DaemonSet"
	// DeploymentService name of the deploy resource
	DeploymentService = "Deployment"
	// StatefulSetService name of the sts resource
	StatefulSetService = "StatefulSet"
)

// checkRuntimeStatus check the application objects' runtime status, which the aim replicas is equal to the current
// replicas. In addition, if application objects have 0 (zero) replicas, which may because of the namespace or node
// pod limitation.
func checkRuntimeStatus(ctx context.Context, r *ServicePackageReconciler, pack *enginev1alpha1.ServicePackage,
	workload enginev1alpha1.Workload) (bool, error) {
	deployOk, deployErr := checkDeploymentsRuntime(ctx, r, pack, workload.Deployments)
	if deployErr != nil {
		return false, deployErr
	}
	dsOk, dsErr := checkDaemonSetsRuntime(ctx, r, pack, workload.DaemonSets)
	if dsErr != nil {
		return false, dsErr
	}
	stsOk, stsErr := checkStatefulSetsRuntime(ctx, r, pack, workload.StatefulSets)
	if stsErr != nil {
		return false, stsErr
	}
	return deployOk && dsOk && stsOk, nil
}

func checkDeploymentsRuntime(ctx context.Context, r *ServicePackageReconciler, pack *enginev1alpha1.ServicePackage,
	deploySpecs []enginev1alpha1.ServiceDeploymentSpec) (bool, error) {
	deployMap := getDeployMap(deploySpecs)
	namespace := pack.Namespace
	ok := true
	for name := range deployMap {
		deploy := appsv1.Deployment{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &deploy); err != nil {
			if errors.IsNotFound(err) {
				ok = false
				pack.SetToFailed(notFoundReason(DeploymentService, name, namespace))
				continue
			}
			klog.Errorf("failed to get development [%s], because: %s", name, err)
			return false, err
		}
		if deploy.Status.AvailableReplicas == 0 {
			ok = false
			pack.SetToUnknown(unknownReason(DeploymentService, name, namespace))
		} else if deploy.Status.UnavailableReplicas > 0 {
			ok = false
			pack.SetToFailed(failedReason(DeploymentService, name, namespace, deploy.Status.Replicas,
				deploy.Status.AvailableReplicas))
		}
	}
	return ok, nil
}

func checkDaemonSetsRuntime(ctx context.Context, r *ServicePackageReconciler, pack *enginev1alpha1.ServicePackage,
	daemonSpecs []enginev1alpha1.ServiceDaemonSetSpec) (bool, error) {
	daemonSetMap := getDaemonSetMap(daemonSpecs)
	namespace := pack.Namespace
	ok := true
	for name := range daemonSetMap {
		ds := appsv1.DaemonSet{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &ds); err != nil {
			if errors.IsNotFound(err) {
				ok = false
				pack.SetToFailed(notFoundReason(DaemonSetService, name, namespace))
				continue
			}
			klog.Errorf("failed to get daemon set [%s], because: %s", name, err)
			return false, err
		}
		if ds.Status.NumberAvailable == 0 {
			ok = false
			pack.SetToUnknown(unknownReason(DaemonSetService, name, namespace))
		} else if ds.Status.DesiredNumberScheduled != ds.Status.CurrentNumberScheduled {
			ok = false
			pack.SetToFailed(failedReason(DaemonSetService, name, namespace,
				ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled))
		}
	}
	return ok, nil
}

func checkStatefulSetsRuntime(ctx context.Context, r *ServicePackageReconciler, pack *enginev1alpha1.ServicePackage,
	statefulSpecs []enginev1alpha1.ServiceStatefulSetSpec) (bool, error) {
	stsMap := getStatefulSetMap(statefulSpecs)
	namespace := pack.Namespace
	ok := true
	for name := range stsMap {
		sts := appsv1.StatefulSet{}
		if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &sts); err != nil {
			if errors.IsNotFound(err) {
				ok = false
				pack.SetToFailed(notFoundReason(StatefulSetService, name, namespace))
				continue
			}
			klog.Errorf("failed to get stateful set [%s], because: %s", name, err)
			return false, err
		}
		if sts.Status.CurrentReplicas == 0 {
			ok = false
			pack.SetToUnknown(unknownReason(StatefulSetService, name, namespace))
		} else if sts.Status.Replicas != sts.Status.CurrentReplicas {
			ok = false
			pack.SetToFailed(failedReason(StatefulSetService, name, namespace, sts.Status.Replicas,
				sts.Status.CurrentReplicas))
		}
	}
	return ok, nil
}

func failedReason(serviceType, name, namespace string, desired, actual int32) string {
	return fmt.Sprintf("%s: %s in namespace %s is running failed, it may not provide normal service. "+
		"This service want %d replica(s), but current only have %d replica(s).",
		serviceType, name, namespace, desired, actual)
}

func notFoundReason(serviceType, name, namespace string) string {
	return fmt.Sprintf("%s: %s in namespace %s is not found", serviceType, name, namespace)
}

func unknownReason(serviceType, name, namespace string) string {
	return fmt.Sprintf("%s: %s in namespace %s has unknown reasons to run 0 replicas in cluster. It may because "+
		"of the pod limitation in this namespace, please check the cluster resources and limitation.",
		serviceType, name, namespace)
}
