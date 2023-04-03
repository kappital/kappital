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

package client

import (
	"fmt"
	"os"
	"sync"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	crdOnce sync.Once

	crdClient *clientset.Clientset
)

// GetCRDClient returns with initialized CRD client,
// it tries to init if the client not exist yet.
func GetCRDClient() *clientset.Clientset {
	crdOnce.Do(func() {
		var err error
		crdClient, err = newCRDClient()
		if err != nil {
			klog.Errorf("cannot get crd client, err: %v", err)
			os.Exit(1)
		}
	})

	return crdClient
}

func newCRDClient() (*clientset.Clientset, error) {
	cli, err := clientset.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		return nil, fmt.Errorf("failed to generate CRD client with kubeconfig certs, err: %w", err)
	}
	return cli, nil
}
