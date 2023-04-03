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

	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	out "github.com/kappital/kappital/pkg/apis/view"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	serviceName  string
	clusterName  string
	outputFormat string
	allResult    bool
}

// Cmd singleton pattern of get Service from cluster
var Cmd operation

func (o operation) getArgumentMap() map[string]interface{} {
	return map[string]interface{}{
		"service-name":                     o.serviceName,
		kappctl.Cluster.GetFlagName():      o.clusterName,
		kappctl.OutputFormat.GetFlagName(): o.outputFormat,
		kappctl.GetAll.GetFlagName():       o.allResult,
	}
}

// NewCommand create the new command for GET runtime Cloud Native Service
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service [service-name]",
		Aliases: []string{"services", "svc"},
		Short:   "Query one or many Cloud Native Services",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.OutputFormat.AddStringFlag(&o.outputFormat, cmd)
	kappctl.Cluster.AddStringFlag(&o.clusterName, cmd)
	kappctl.GetAll.AddBoolFlag(&o.allResult, cmd)
	return cmd
}

// PreRunE run before getting the runtime service, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	if len(args) == 1 {
		o.serviceName = args[0]
	}
	var err error
	if o.config, err = kappctl.GetConfig(); err != nil {
		return err
	}
	return kappctl.IsInputValidate(o.getArgumentMap())
}

// RunE get the runtime service
func (o *operation) RunE() error {
	getDetail := false
	if len(o.serviceName) > 0 {
		o.serviceName = fmt.Sprintf("/%s", o.serviceName)
		getDetail = true
	}
	url := o.config.BuildManagerURL(kappctl.GetServiceURL, []interface{}{o.serviceName, o.clusterName})
	if getDetail {
		url += fmt.Sprintf("&detail=%v", getDetail)
	}
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodGet,
		Path:      url,
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("cannot get the service binding, http code: %d, msg: %s", code, string(buf))
	}
	return outputResult(buf, len(o.serviceName) == 0, o.outputFormat)
}

func outputResult(buf []byte, isMultiple bool, format string) error {
	if len(format) > 0 {
		return kappctl.OutputYAMLOrJSONString(buf, format)
	}

	var itfs []interface{}
	if isMultiple {
		var svcs []instancev1alpha1.CloudNativeServiceInstance
		if err := json.Unmarshal(buf, &svcs); err != nil {
			return fmt.Errorf("failed to unmarshal http response: %s", err)
		}
		if len(svcs) == 0 {
			fmt.Println("No resources found")
			return nil
		}
		for _, svc := range svcs {
			itfs = append(itfs, convertServiceToTable(svc))
		}
		kappctl.TableFormatter(itfs)
		return nil
	}
	var svc instancev1alpha1.CloudNativeServiceInstance
	if err := json.Unmarshal(buf, &svc); err != nil {
		return fmt.Errorf("failed to unmarshal http response: %s", err)
	}
	itfs = append(itfs, convertServiceToTable(svc))
	kappctl.TableFormatter(itfs)
	return nil
}

func convertServiceToTable(svc instancev1alpha1.CloudNativeServiceInstance) out.Service {
	res := out.Service{
		Name:      svc.Name,
		Cluster:   svc.Spec.ClusterName,
		Namespace: svc.Namespace,
		Phase:     string(svc.Status.Phase),
		Message:   svc.Status.Message,
		Created:   kappctl.GetAgeOutput(svc.CreationTimestamp.Time),
	}
	return res
}
