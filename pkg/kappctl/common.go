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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"

	"github.com/kappital/kappital/pkg/utils/file"
)

const (
	tableCellMaxLen = 64
	dropLen         = 4

	envKey = "KAPPITALCONFIG"
	// ConfigFile the default config file path
	ConfigFile = ".kappital/config"

	// DeployInstanceURL the url format for deploying the instance
	DeployInstanceURL = "%v/api/v1alpha1/servicebinding/%v/instance?cluster_name=%v"
	// DeployServiceURL the url format for deploying the service package into cluster
	DeployServiceURL = "%v/api/v1alpha1/servicebinding"

	// GetInstancesURL the url format for getting the instances from the cluster
	GetInstancesURL = "%v/api/v1alpha1/servicebinding/%v/instance"
	// GetServiceURL the url format for getting the service from the cluster
	GetServiceURL = "%v/api/v1alpha1/servicebinding%v?cluster_name=%v"
	// GetServicesURL the url format for getting the service list from the cluster
	GetServicesURL = "%v/api/v1alpha1/servicebinding?cluster_name=%v"

	// DeleteInstanceURL the url format for deleting the instance from the cluster
	DeleteInstanceURL = "%v/api/v1alpha1/servicebinding/%v/instance/%v?cluster_name=%v"
	// DeleteServiceURL the url format for deleting service
	DeleteServiceURL = "%v/api/v1alpha1/servicebinding/%v?cluster_name=%v"

	// yamlOutputFormat output the format as a yaml structure
	yamlOutputFormat = "yaml"
	// jsonOutputFormat output the format as a json structure
	jsonOutputFormat = "json"
)

// input validation arguments
const (
	maxStringLength = 64
)

// validFormat does the output format is valid
func validFormat(format string) error {
	if !sets.NewString(yamlOutputFormat, jsonOutputFormat, "").Has(strings.ToLower(format)) {
		return fmt.Errorf("output format [%v] is not supported", format)
	}
	return nil
}

// GetConfig get the manager information from config file
func GetConfig() (*Config, error) {
	path, err := getKappitalConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read config file %v: %v", path, err)
	}
	configOptions := Config{}
	if err = json.Unmarshal(config, &configOptions); err != nil {
		return nil, fmt.Errorf("unmarshal config file %v; %v", path, err)
	}
	return &configOptions, nil
}

func getKappitalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDir, ConfigFile)
	if file.IsFileExist(path) {
		return path, nil
	}
	config := os.Getenv(envKey)
	if len(config) == 0 {
		return "", fmt.Errorf("missing the kappctl config file")
	}
	return config, nil
}

// TableFormatter format the table
func TableFormatter(slice []interface{}) {
	if len(slice) == 0 {
		return
	}

	tab := tablewriter.NewWriter(os.Stdout)

	tab.SetAutoFormatHeaders(true)
	tab.SetAutoWrapText(false)

	tab.SetColumnSeparator("")
	tab.SetCenterSeparator("")

	tab.SetAlignment(tablewriter.ALIGN_LEFT)
	tab.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	tab.SetBorder(false)
	tab.SetRowSeparator("")
	tab.SetTablePadding("   ") // pad with tabs
	tab.SetNoWhiteSpace(true)
	tab.SetHeaderLine(false)

	structType := reflect.TypeOf(slice[0])
	row, col := len(slice), structType.NumField()
	header := make([]string, 0, col)
	for i := 0; i < col; i++ {
		field := structType.Field(i)
		header = append(header, field.Name)
	}

	tab.SetHeader(header)

	tableData := make([][]string, row)
	for i := 0; i < row; i++ {
		tableData[i] = make([]string, 0, col)
		fields := reflect.ValueOf(slice[i])
		for j := 0; j < col; j++ {
			field := fields.Field(j).String()
			if field == "null" || field == "{}" {
				field = ""
			}
			if len(field) > tableCellMaxLen {
				field = field[:tableCellMaxLen-dropLen] + " ..."
			}
			tableData[i] = append(tableData[i], field)
		}
	}
	tab.AppendBulk(tableData) // Add Bulk Data

	tab.Render()
}

// GetAgeOutput obtains the creation time with the hour, minute, and second.
func GetAgeOutput(timestamp time.Time) string {
	age := time.Since(timestamp).Round(time.Second).String()
	return age[:strings.IndexAny(age, "hms")+1] + " ago"
}

// OutputYAMLOrJSONString output the buf as json or yaml from the format type
func OutputYAMLOrJSONString(buf []byte, format string) error {
	switch strings.ToLower(format) {
	case yamlOutputFormat:
		body, err := yaml.JSONToYAML(buf)
		if err != nil {
			return fmt.Errorf("failed to convert json response to yaml: %v", err)
		}
		fmt.Println(string(body))
		return nil
	case jsonOutputFormat:
		fmt.Println(string(buf))
		return nil
	case "":
		return nil
	}
	return fmt.Errorf("the format is invalid")
}

// IsInputValidate does the input arguments is valid or not, such as the max length of string should be 64 bytes.
func IsInputValidate(inputs map[string]interface{}) error {
	for k, v := range inputs {
		str, ok := v.(string)
		if !ok {
			continue
		}
		if len(str) > maxStringLength {
			return fmt.Errorf("the argument [%v] length is greather than the max length %d", k, maxStringLength)
		}
		var err error
		switch k {
		case Cluster.GetFlagName():
			err = isValidClusterName(str)
		case OutputFormat.GetFlagName():
			err = validFormat(str)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func isValidClusterName(clusterName string) error {
	if len(clusterName) == 0 {
		return nil
	}
	arr := strings.Split(clusterName, "-")
	regChar, err := regexp.Compile("^[A-Za-z][A-Za-z0-9]+$")
	if err != nil {
		return err
	}
	regCharNum, err := regexp.Compile("^[A-Za-z0-9]+$")
	if err != nil {
		return err
	}
	for i, s := range arr {
		if len(s) == 0 {
			return fmt.Errorf("cannot use empty string after '-'")
		}
		if i == 0 {
			if !regChar.Match([]byte(s)) {
				return fmt.Errorf("contain invalid character(s)")
			}
			continue
		}
		if !regCharNum.Match([]byte(s)) {
			return fmt.Errorf("contain invalid character(s)")
		}
	}
	return nil
}
