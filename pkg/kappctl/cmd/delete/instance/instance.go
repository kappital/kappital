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

package instance

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	serviceName  string
	instanceName string
	clusterName  string
}

// Cmd singleton pattern of delete Instance to the cluster
var Cmd operation

func (o operation) getArgumentMap() map[string]interface{} {
	return map[string]interface{}{
		kappctl.ServiceName.GetFlagName():  o.serviceName,
		kappctl.InstanceName.GetFlagName(): o.instanceName,
		kappctl.Cluster.GetFlagName():      o.clusterName,
	}
}

// NewCommand create the new command for DELETE Cloud Native Service Instance
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance instance-name",
		Short: "Delete a Cloud Native Service Instance in a cluster",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.Cluster.AddStringFlag(&o.clusterName, cmd)
	kappctl.ServiceName.AddStringFlag(&o.serviceName, cmd)
	kappctl.ServiceName.MarkFlagRequired(cmd)
	return cmd
}

// PreRunE run before deleting the service instance, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	o.instanceName = args[0]
	var err error
	o.config, err = kappctl.GetConfig()
	if err != nil {
		return err
	}
	return kappctl.IsInputValidate(o.getArgumentMap())
}

// RunE delete the service instance to cluster
func (o *operation) RunE() error {
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodDelete,
		Path:      o.config.BuildManagerURL(kappctl.DeleteInstanceURL, []interface{}{o.serviceName, o.instanceName, o.clusterName}),
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return fmt.Errorf("delete service instance %s failed, err: %s", o.instanceName, err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("delete service instance %s failed, statusCode: %d, detail: %s",
			o.instanceName, code, string(buf))
	}
	fmt.Printf("delete service instance %s success.\n", o.instanceName)
	return nil
}
