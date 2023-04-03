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
	"fmt"
)

// Config the address and port of the manager
type Config struct {
	ManagerHTTPSServer           string `json:"manager-https-server,omitempty"`
	ManagerClientCertificateData string `json:"manager-client-certificate-data,omitempty"`
	ManagerClientKeyData         string `json:"manager-client-key-data,omitempty"`
	ManagerCA                    string `json:"manager-ca,omitempty"`
	ManagerSkipVerify            bool   `json:"manager-skip-verify,omitempty"`
}

// BuildManagerURL build the manager http/https request
func (c Config) BuildManagerURL(urlFormat string, args []interface{}) string {
	parameters := make([]interface{}, 0, len(args)+1)
	parameters = append(parameters, c.ManagerHTTPSServer)
	parameters = append(parameters, args...)
	return fmt.Sprintf(urlFormat, parameters...)
}
