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
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	serviceName string
	clusterName string
}

// Cmd singleton pattern of delete Whole Service to the cluster
var Cmd operation

func (o operation) getArgumentMap() map[string]interface{} {
	return map[string]interface{}{
		kappctl.ServiceName.GetFlagName(): o.serviceName,
		kappctl.Cluster.GetFlagName():     o.clusterName,
	}
}

// NewCommand create the new command for DELETE Whole Cloud Native Service Instance
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service service-name",
		Short: "Delete a Cloud Native Service in a cluster",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.Cluster.AddStringFlag(&o.clusterName, cmd)
	return cmd
}

// PreRunE run before deleting the service, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	o.serviceName = args[0]
	var err error
	o.config, err = kappctl.GetConfig()
	if err != nil {
		return err
	}
	return kappctl.IsInputValidate(o.getArgumentMap())
}

// RunE delete the service to cluster
func (o *operation) RunE() error {
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodDelete,
		Path:      o.config.BuildManagerURL(kappctl.DeleteServiceURL, []interface{}{o.serviceName, o.clusterName}),
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return fmt.Errorf("delete service instance %s failed, err: %s", o.serviceName, err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("delete service %s failed, statusCode: %d, detail: %s",
			o.serviceName, code, string(buf))
	}
	fmt.Printf("delete service %s success.\n", o.serviceName)
	return nil
}
