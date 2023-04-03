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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/kappital/kappital/pkg/apis"
	svcv1alpha1 "github.com/kappital/kappital/pkg/apis/service/v1alpha1"
	instancev1alpha1 "github.com/kappital/kappital/pkg/apis/serviceinstance/v1alpha1"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/convert"
	"github.com/kappital/kappital/pkg/utils/gateway"
)

type operation struct {
	config *kappctl.Config

	serviceName string
	dirPath     string

	cns          *svcv1alpha1.CloudNativeService
	resourcePath string
}

// Cmd singleton pattern of create Instance to the cluster
var Cmd operation

// NewCommand create the new command for CREATE Cloud Native Service Instance
func (o *operation) NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance [service-name]",
		Short: "Install an instance of a specific Cloud Native Service",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.FilePath.AddStringFlag(&o.resourcePath, cmd)
	kappctl.PackageDir.AddStringFlag(&o.dirPath, cmd)
	return cmd
}

// PreRunE run before creating the service instance, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	var err error
	o.config, err = kappctl.GetConfig()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		o.serviceName = args[0]
		return nil
	}
	if len(o.dirPath) == 0 {
		return fmt.Errorf("service name and package path must have one parameter")
	}
	dir, err := filepath.Abs(o.dirPath)
	if err != nil {
		return err
	}
	cns, err := convert.GetLoader().TransferToCloudNativeService(dir)
	if err != nil {
		return err
	}
	if cns == nil {
		return fmt.Errorf("the service package cannot get the correct content")
	}
	o.cns = cns
	o.serviceName = cns.Name
	return err
}

// RunE create the service instance to cluster
func (o *operation) RunE() error {
	sic := instancev1alpha1.ServiceInstanceCreation{
		ClusterID: apis.DefaultCluster,
	}
	if o.cns != nil {
		sic.Service = *(o.cns)
	}
	var err error
	if len(o.resourcePath) == 0 {
		if o.cns == nil {
			return fmt.Errorf("cannot deploy instance without service package, err: only pass in the instance " +
				"name, but does not have the resource file")
		}
		for _, manifest := range o.cns.Spec.Manifests {
			sic.InstanceCustomResources = append(sic.InstanceCustomResources, getInstanceCustomResource(manifest.Spec))
		}
	} else {
		if sic.InstanceCustomResources, err = getCRContent(o.resourcePath); err != nil {
			return err
		}
	}
	code, buf, err := gateway.CommonUtilRequest(&gateway.RequestInfo{
		Method:    http.MethodPost,
		Path:      o.config.BuildManagerURL(kappctl.DeployInstanceURL, []interface{}{o.serviceName, apis.DefaultCluster}),
		Body:      sic,
		CaCrt:     o.config.ManagerCA,
		ClientCrt: o.config.ManagerClientCertificateData,
		ClientKey: o.config.ManagerClientKeyData,
		Skip:      o.config.ManagerSkipVerify,
	})
	if err != nil {
		return fmt.Errorf("deploy service %s failed, err: %s", o.serviceName, err)
	}
	if code != http.StatusOK {
		return fmt.Errorf("deploy service %s failed, statusCode: %d, detail: %s", o.serviceName, code, string(buf))
	}
	fmt.Printf("deploy service %s success.\n", o.serviceName)
	return nil
}

func getCRContent(path string) ([]instancev1alpha1.InstanceCustomResource, error) {
	crFilePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get cr path failed, err: %s", err)
	}
	crByte, err := ioutil.ReadFile(filepath.Clean(crFilePath))
	if err != nil {
		return nil, fmt.Errorf("read cr file failed, err: %s", err)
	}

	decoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(crByte))
	var crs []instancev1alpha1.InstanceCustomResource
	for {
		var instanceCR instancev1alpha1.InstanceCustomResource
		if err = decoder.Decode(&instanceCR); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		crs = append(crs, instanceCR)
	}
	return crs, nil
}
func getInstanceCustomResource(spec svcv1alpha1.CustomServiceDefinitionSpec) instancev1alpha1.InstanceCustomResource {
	group, ok := spec.CRD.Spec["group"].(string)
	if !ok {
		return instancev1alpha1.InstanceCustomResource{}
	}
	kind, ok := spec.CRD.Spec["names"].(map[string]interface{})["kind"].(string)
	if !ok {
		return instancev1alpha1.InstanceCustomResource{}
	}
	return instancev1alpha1.InstanceCustomResource{
		TypeMeta:   metav1.TypeMeta{Kind: kind, APIVersion: fmt.Sprintf("%s/%s", group, spec.CRVersions[0].Name)},
		ObjectMeta: metav1.ObjectMeta{Name: spec.CRVersions[0].CRName},
		Spec:       []byte(spec.CRVersions[0].DefaultValues),
	}
}
