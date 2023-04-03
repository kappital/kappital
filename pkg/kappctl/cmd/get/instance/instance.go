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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	out "github.com/kappital/kappital/pkg/apis/view"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	instanceName string

	serviceName string
	namespace   string
	clusterName string

	outputFormat string
	allResult    bool
}

// Cmd singleton pattern of get Service Instance from cluster
var Cmd operation

func (o operation) getArgumentMap() map[string]interface{} {
	return map[string]interface{}{
		kappctl.InstanceName.GetFlagName(): o.instanceName,
		kappctl.ServiceName.GetFlagName():  o.serviceName,
		kappctl.Namespace.GetFlagName():    o.namespace,
		kappctl.Cluster.GetFlagName():      o.clusterName,
		kappctl.GetAll.GetFlagName():       o.allResult,
		kappctl.OutputFormat.GetFlagName(): o.outputFormat,
	}
}

// NewCommand create the new command for GET Cloud Native Service Instance
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance instance-name",
		Aliases: []string{"instances"},
		Short:   "Query instances of a Cloud Native Instance.",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.ServiceName.AddStringFlag(&o.serviceName, cmd)
	kappctl.Namespace.AddStringFlag(&o.namespace, cmd)
	kappctl.Cluster.AddStringFlag(&o.clusterName, cmd)

	kappctl.OutputFormat.AddStringFlag(&o.outputFormat, cmd)
	kappctl.GetAll.AddBoolFlag(&o.allResult, cmd)
	return cmd
}

// PreRunE run before getting the service instance(s), check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	var err error
	if o.config, err = kappctl.GetConfig(); err != nil {
		return err
	}
	if len(args) == 1 {
		o.instanceName = args[0]
	}
	if len(args) != 1 && !o.allResult && len(o.serviceName) == 0 {
		return fmt.Errorf("please specify the instance and service name")
	}
	if o.allResult {
		return nil
	}
	if len(o.serviceName) == 0 {
		return fmt.Errorf("the instance is managered by a service, please specify the service name")
	}
	return kappctl.IsInputValidate(o.getArgumentMap())
}

// RunE get the service instance(s)
func (o *operation) RunE() error {
	if o.allResult {
		res, err := o.getAllServiceInstances()
		if err != nil {
			return err
		}
		kappctl.TableFormatter(res)
		return nil
	}
	res, err := o.getServiceInstances()
	if err != nil {
		return err
	}
	kappctl.TableFormatter(res)
	return nil
}

func (o *operation) getAllServiceInstances() ([]interface{}, error) {
	// 1. Get deployed service list
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodGet,
		Path:      o.config.BuildManagerURL(kappctl.GetServicesURL, []interface{}{o.clusterName}),
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("cannot get the service list, because get the http code: %d", code)
	}
	var svcs []instancev1alpha1.CloudNativeServiceInstance
	if err = json.Unmarshal(buf, &svcs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal http response: %s", err)
	}
	if len(svcs) == 0 {
		fmt.Println("No service deployed into cluster.")
		return nil, nil
	}
	// 2. Get each instance from the service list, and add all instance to a slice
	var itfs []interface{}
	for _, svc := range svcs {
		instances, err := o.getInstanceListServiceName(svc)
		if err != nil {
			return nil, err
		}
		if instances != nil {
			itfs = append(itfs, instances...)
		}
	}
	return itfs, nil
}

func (o *operation) getServiceInstances() ([]interface{}, error) {
	o.serviceName = fmt.Sprintf("/%s", o.serviceName)
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodGet,
		Path:      o.config.BuildManagerURL(kappctl.GetServiceURL+"&detail=true", []interface{}{o.serviceName, o.clusterName}),
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("cannot get the service list, because get the http code: %d", code)
	}

	var svc instancev1alpha1.CloudNativeServiceInstance
	if err := json.Unmarshal(buf, &svc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal http response: %s", err)
	}
	return o.getInstanceListServiceName(svc)
}

func (o *operation) getInstanceListServiceName(svc instancev1alpha1.CloudNativeServiceInstance) ([]interface{}, error) {
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodGet,
		Path:      o.config.BuildManagerURL(kappctl.GetInstancesURL, []interface{}{svc.Name}),
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("cannot get the instance list, because get the http code: %d", code)
	}
	var ins []models.InstanceModel
	if err = json.Unmarshal(buf, &ins); err != nil {
		return nil, fmt.Errorf("failed to unmarshal http response: %s", err)
	}

	if len(o.outputFormat) > 0 {
		return nil, outputYamlOrJSON(svc, ins, o.outputFormat)
	}

	var itfs []interface{}
	for _, item := range ins {
		if len(o.instanceName) == 0 {
			itfs = append(itfs, convertInstanceToTable(item, svc.Name, svc.Status.Phase))
		} else if o.instanceName == item.Name {
			itfs = append(itfs, convertInstanceToTable(item, svc.Name, svc.Status.Phase))
		}
	}
	return itfs, nil
}

func outputYamlOrJSON(svc instancev1alpha1.CloudNativeServiceInstance, ins []models.InstanceModel, format string) error {
	svc.Spec.CustomResources = make([]instancev1alpha1.Resource, 0, len(ins))
	for _, item := range ins {
		addIn := instancev1alpha1.Resource{
			TypeMeta: metav1.TypeMeta{
				Kind:       item.Kind,
				APIVersion: item.APIVersion,
			},
			Name:       item.Name,
			Namespace:  item.Namespace,
			UID:        item.ServiceID,
			Status:     item.Status,
			RawMessage: item.RawResource,
		}
		if svc.Status.Phase == instancev1alpha1.PendingPhase {
			addIn.Status = string(instancev1alpha1.PendingPhase)
		} else if svc.Status.Phase != instancev1alpha1.SucceededPhase {
			addIn.Status = string(instancev1alpha1.UnknownPhase)
		}
		svc.Spec.CustomResources = append(svc.Spec.CustomResources, addIn)
	}
	jsonBytes, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	return kappctl.OutputYAMLOrJSONString(jsonBytes, format)
}

func convertInstanceToTable(ins models.InstanceModel, serviceName string, phase instancev1alpha1.Phase) out.Instance {
	res := out.Instance{
		ServiceName:  serviceName,
		ClusterName:  ins.ClusterName,
		InstanceName: ins.Name,
		Namespace:    ins.Namespace,
		Status:       ins.Status,
		Created:      kappctl.GetAgeOutput(ins.CreateTimestamp),
	}
	switch phase {
	case instancev1alpha1.SucceededPhase:
		res.Status = ins.Status
	case instancev1alpha1.PendingPhase:
		res.Status = string(phase)
	default:
		res.Status = string(instancev1alpha1.UnknownPhase)
	}
	return res
}
