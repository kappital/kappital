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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
)

const (
	subResourceOwnerKey = "metadata.ownerReferences"

	deletedPeriod = time.Hour
)

// ServicePackageReconciler reconciles a Operator object
type ServicePackageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ServicePackageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// filter the namespace which the reconciler to watch. now, we only accept for the namespace kappital-system.
	// if the log shows something like "but the current namespace is ." means the *default* namespace.
	if req.Namespace != apis.KappitalSystemNamespace {
		klog.Warningf("pack engine only accept the custom resource with namespace *%s*, "+
			"but the current namespace is %s.", apis.KappitalSystemNamespace, req.Namespace)
		return ctrl.Result{}, nil
	}
	var err error
	pack := &enginev1alpha1.ServicePackage{}

	if err = r.Get(ctx, req.NamespacedName, pack); err != nil {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		klog.Errorf("failed to get pack package [%s] in namespace [%s], err: %v", req.Name, req.Namespace, err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// if the service package is already deleted, do not reconcile it, and wait for 1 hour to re-check this service
	// package. during this 1-hour period, the manager will delete this package
	if pack.IsDeleted() {
		return ctrl.Result{RequeueAfter: deletedPeriod}, nil
	}
	pack.VerifyStatus()

	if updateErr := r.Status().Update(ctx, pack); updateErr != nil {
		klog.Errorf("update pack package [%s:%s] status failed, because: %s",
			pack.Spec.Name, pack.Spec.Version, updateErr)
		return ctrl.Result{}, updateErr
	}

	if err = r.reconcileSubResources(ctx, pack); err != nil {
		klog.Errorf("unable to reconcile, err: %s", err)
	}
	return ctrl.Result{}, r.analysisFinalError(ctx, pack, err)
}

// reconcileSubResources will reconcile the sub resources for the service package. It will include:
// CustomResourceDefinitions, Permission-Related resources (ServiceAccount, ClusterRole, and ClusterRoleBinding),
// and Application Objects (DaemonSet, Deployment, and StatefulSet). Finally, the Service-Engine will check does the
// Application Objects is running correctly.
func (r *ServicePackageReconciler) reconcileSubResources(ctx context.Context,
	pack *enginev1alpha1.ServicePackage) error {
	resource, err := enginev1alpha1.TranslateBase64CodeToResource(pack.Spec.Resources)
	if err != nil {
		klog.Errorf("cannot analysis the binary code to struct, err: %s", err)
		return err
	}
	if err = r.reconcileCRD(ctx, pack, resource.CustomResourceDefinitions); err != nil {
		klog.Warningf("some problem for reconcile crds in cluster ([%s-%s]) error: %v, there may already exist "+
			"the old crds in the cluster, operator will continue use the old crds. if the operator cannot running,"+
			"please manual cleanup crds", pack.Name, pack.Spec.Version, err)
	}

	if err = r.reconcilePermission(ctx, pack, resource.Permissions); err != nil {
		pack.SetToFailed("Unable to reconcile permissions")
		klog.Errorf("failed to reconcile permissions for [%s-%s]: %v", pack.Name, pack.Spec.Version, err)
		return err
	}

	if err = r.reconcileDeployment(ctx, pack, resource.Workload.Deployments); err != nil {
		pack.SetToFailed("Unable to reconcile deployment")
		klog.Errorf("failed to reconcile deployment for [%s-%s]: %v", pack.Name, pack.Spec.Version, err)
		return err
	}

	if err = r.reconcileDaemonSet(ctx, pack, resource.Workload.DaemonSets); err != nil {
		pack.SetToFailed("Unable to reconcile daemon set")
		klog.Errorf("failed to reconcile daemon set for [%s-%s]: %v", pack.Name, pack.Spec.Version, err)
		return err
	}

	if err = r.reconcileStatefulSet(ctx, pack, resource.Workload.StatefulSets); err != nil {
		pack.SetToFailed("Unable to reconcile stateful set")
		klog.Errorf("failed to reconcile stateful set for [%s-%s]: %v", pack.Name, pack.Spec.Version, err)
		return err
	}

	if pack.NeedCheckRuntime() {
		ok, err := checkRuntimeStatus(ctx, r, pack, resource.Workload)
		if err != nil {
			return err
		}
		if ok {
			pack.SetToRunning()
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServicePackageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addOwnerIndex := func(rawObject client.Object) []string {
		return genOwnerKey(metav1.GetControllerOf(rawObject.(metav1.Object)))
	}

	//  support sub-resources field-selector
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support deployment field selector. err: %s", err)
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.DaemonSet{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support daemonSet field selector. err: %s", err)
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.StatefulSet{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support statefulSet field selector. err: %s", err)
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &rbacv1.ClusterRole{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support clusterRole field selector. err: %s", err)
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &rbacv1.ClusterRoleBinding{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support clusterRoleBinding field selector. err: %s", err)
		return err
	}
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ServiceAccount{},
		subResourceOwnerKey, addOwnerIndex); err != nil {
		klog.Errorf("failed to support ServiceAccount field selector. err: %s", err)
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&enginev1alpha1.ServicePackage{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Complete(r)
}

func genOwnerKey(owner *metav1.OwnerReference) []string {
	if owner == nil {
		return nil
	}
	if owner.APIVersion != enginev1alpha1.GroupVersion.Version || owner.Kind != enginev1alpha1.ServicePackagesKind {
		return nil
	}
	klog.Infof(owner.Name)
	return []string{owner.Name}
}

// analysisFinalError will get the final error, and update the service package status.
func (r *ServicePackageReconciler) analysisFinalError(ctx context.Context, pack *enginev1alpha1.ServicePackage,
	err error) error {
	pack.UpdateStatus(err)
	if updateErr := r.Status().Update(ctx, pack); updateErr != nil {
		klog.Errorf("update pack package [%s:%s] status failed, because: %s",
			pack.Spec.Name, pack.Spec.Version, updateErr)
		return updateErr
	}
	return err
}

func belongToOwnerReference(name string, references []metav1.OwnerReference) bool {
	for _, reference := range references {
		if isEngineResource(name, reference) {
			return true
		}
	}
	return false
}

func isEngineResource(name string, reference metav1.OwnerReference) bool {
	return reference.APIVersion == enginev1alpha1.GroupVersion.String() &&
		reference.Kind == enginev1alpha1.ServicePackageKind &&
		reference.Name == name
}
