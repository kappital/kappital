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

package kappctl

import (
	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/apis"
)

var (
	// ManagerIPAddr the manager service ip address
	ManagerIPAddr = newInputFlag("manager-addr", "", "", "the ip Addr of kappital-manager")
	// ManagerHTTPSPort the manager https service port
	ManagerHTTPSPort = newInputFlag("manager-https-port", "", "", "the HTTPS Port of kappital-manager")
	// ManagerClientCertFile the manager certificate file path
	ManagerClientCertFile = newInputFlag("manager-client-cert", "", "", "the HTTPS client certificate file of kappital-manager")
	// ManagerClientKeyFile the manager key file path
	ManagerClientKeyFile = newInputFlag("manager-client-key", "", "", "the HTTPS client key file of kappital-manager")
	// ManagerCAFile the manager ca file path
	ManagerCAFile = newInputFlag("manager-ca", "", "", "the HTTPS ca file for kappital-manager")
	// ManagerSkipVerify the manager skip verify controller
	ManagerSkipVerify = newInputFlag("manager-skip-verify", "", false, "connect to kappital-manager need to skip verify")
	// OutputFormat the output result format of the result, now only offer yaml and json
	OutputFormat = newInputFlag("output", "o", "", "the output format of the queried resource, can be yaml or json")
	// Cluster name of service and service instance deployed
	Cluster = newInputFlag("cluster", "c", apis.DefaultCluster, "the cluster scope of the Cloud Native Service")
	// GetAll service or instance from cluster
	GetAll = newInputFlag("all", "A", false, `query resources across all repos/clusters. When this flag is set, --repo or --cluster flag will HAS NO EFFECT`)
	// Namespace of the service or service instance deployed
	Namespace = newInputFlag("namespace", "n", apis.DefaultNamespace, "the namespace of the specified instance")
	// ServiceName the name of deployed or will deploy service
	ServiceName = newInputFlag("service", "s", "", "the cloud native service name")
	// InstanceName the name of deployed or will deploy instance
	InstanceName = newInputFlag("instance", "i", "", "the cloud native service instance name")
	// FilePath of the instance file
	FilePath = newInputFlag("file", "f", "", "the custom resource file path")
	// FileName of the directory name of the kappital package
	FileName = newInputFlag("name", "", "kappital-demo", "the kappital package name")
	// PackageVersion of the kappital
	PackageVersion = newInputFlag("version", "v", "0.1.0", "the kappital package version")
	// PackageDir of Cloud Native Package
	PackageDir = newInputFlag("dir", "d", "", "the Cloud Native Package Path")
)

type inputFlag struct {
	name         string
	shortHand    string
	defaultValue interface{}
	usage        string
}

func newInputFlag(name, shortHand string, defaultValue interface{}, usage string) inputFlag {
	return inputFlag{
		name:         name,
		shortHand:    shortHand,
		defaultValue: defaultValue,
		usage:        usage,
	}
}

// GetFlagName for check the key of check argument is valid or not
func (i inputFlag) GetFlagName() string {
	return i.name
}

// AddStringFlag add the flag to the cobra.Command as a string format
func (i inputFlag) AddStringFlag(p *string, cmd *cobra.Command) {
	value, ok := i.defaultValue.(string)
	if !ok {
		return
	}
	cmd.Flags().StringVarP(p, i.name, i.shortHand, value, i.usage)
}

// AddBoolFlag add the flag to the cobra.Command as a boolean format
func (i inputFlag) AddBoolFlag(p *bool, cmd *cobra.Command) {
	value, ok := i.defaultValue.(bool)
	if !ok {
		return
	}
	cmd.Flags().BoolVarP(p, i.name, i.shortHand, value, i.usage)
}

// MarkFlagRequired set the flag to the requirement
func (i inputFlag) MarkFlagRequired(cmd *cobra.Command) {
	_ = cmd.MarkFlagRequired(i.name)
}
