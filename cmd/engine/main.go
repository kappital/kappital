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

package main

import (
	"flag"
	"fmt"
	"github.com/kappital/kappital/pkg/utils/cryption"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kappital/kappital/pkg/apis"
	enginev1alpha1 "github.com/kappital/kappital/pkg/apis/engine/v1alpha1"
	controller "github.com/kappital/kappital/pkg/engine"
	"github.com/kappital/kappital/pkg/utils/gateway"
	"github.com/kappital/kappital/pkg/utils/version"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(enginev1alpha1.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1beta1.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
}

func main() {
	if len(os.Args) > 1 && version.Flags.Has(os.Args[1]) {
		fmt.Println(version.Get(version.ServiceNameEngine).String())
		os.Exit(0)
	}
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controllers kappital-manager. "+
			"Enabling this will ensure there is only one active controllers kappital-engine.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ip, err := gateway.GetLocalIP()
	if err != nil {
		klog.Fatalf("cannot get local ip address, error: %v.", err)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "83317d75.kappital.io",
		NewCache:           cache.MultiNamespacedCacheBuilder([]string{apis.KappitalSystemNamespace}),
	})
	if err != nil {
		klog.Fatalf("unable to start kappital-engine, error: %v.", err)
	}
	cert, key, err := cryption.GetSelfCertAndKey()
	if err != nil {
		klog.Fatalf("cannot get the self https certificate, err: %s", err)
	}
	go gateway.HealthAndReadinessProvider(gateway.ReplaceIP(probeAddr, ip, "8081"), cert, key)

	if err = (&controller.ServicePackageReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		klog.Fatalf("unable to create controller, error: %v.", err)
	}

	if err = version.SetupClusterVersion(); err != nil {
		klog.Fatalf("cannot get the cluster version, error: %v.", err)
	}

	setupLog.Info("starting kappital-engine")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Fatalf("problem running kappital-engine, error: %v.", err)
	}
}
