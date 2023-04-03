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
	"reflect"
	"testing"

	"github.com/brahma-adshonor/gohook"
	extendclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestGetCRDClient(t *testing.T) {
	defer gohook.UnHook(newCRDClient) //nolint:errcheck
	err := gohook.Hook(newCRDClient, func() (*extendclient.Clientset, error) {
		return &extendclient.Clientset{}, nil
	}, nil)
	if err != nil {
		t.Errorf("hook err, :%v", err)
		return
	}
	GetCRDClient()
}

func Test_newCRDClient(t *testing.T) {
	err := gohook.Hook(ctrl.GetConfigOrDie, func() *rest.Config {
		return &rest.Config{}
	}, nil)
	if err != nil {
		t.Errorf("hook err, :%v", err)
		return
	}
	defer gohook.UnHook(ctrl.GetConfigOrDie) //nolint:errcheck

	err = gohook.Hook(extendclient.NewForConfig, func(_ *rest.Config) (*extendclient.Clientset, error) {
		return &extendclient.Clientset{}, nil
	}, nil)
	if err != nil {
		t.Errorf("hook err, :%v", err)
		return
	}
	defer gohook.UnHook(extendclient.NewForConfig) //nolint:errcheck

	cli, err := newCRDClient()
	if !reflect.DeepEqual(cli, &extendclient.Clientset{}) || err != nil {
		t.Errorf("newCRDClient got= %v, getErr = %v, want %v, wantErr nil", cli, err, &extendclient.Clientset{})
	}
}
