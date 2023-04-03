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

package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/convert"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	cns *svcv1alpha1.CloudNativeService
}

type resp struct {
	ServiceID   string `json:"ID"`
	ServiceName string `json:"name"`
}

// Cmd singleton pattern of create Service to the cluster
var Cmd operation

// NewCommand create the new command for CREATE Cloud Native Service
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service [Cloud Native Package Directory Path]",
		Short: "Use a Kappital package to create a Cloud Native Service in a cluster",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	return cmd
}

// PreRunE run before creating the service, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	cns, err := convert.GetLoader().TransferToCloudNativeService(args[0])
	if err != nil {
		return err
	}
	if cns == nil {
		return fmt.Errorf("the service package %s cannot get the correct content", args[0])
	}
	o.cns = cns
	o.config, err = kappctl.GetConfig()
	return err
}

// RunE create the service to cluster
func (o *operation) RunE() error {
	sic := instancev1alpha1.ServiceInstanceCreation{
		InstanceName: o.cns.Spec.Description.Name,
		ClusterID:    apis.DefaultCluster,
		Service:      *(o.cns),
	}
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodPost,
		Path:      o.config.BuildManagerURL(kappctl.DeployServiceURL, []interface{}{}),
		Body:      sic,
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return fmt.Errorf("deploy service %s failed, err: %s", o.cns.Name, err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("deploy service %s failed, statusCode: %d, detail: %s", o.cns.Name, code, string(buf))
	}

	var result resp
	if err = json.Unmarshal(buf, &result); err != nil {
		return fmt.Errorf("json Unmarshal response %s failed, err: %s", string(buf), err)
	}

	fmt.Printf("deploy service %s success.\n%s\n", o.cns.Name,
		fmt.Sprintf("{'service_name': %s, 'service_id': %s}\n", o.cns.Spec.Description.Name, result.ServiceID))
	return nil
}
