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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/file"
)

const (
	minPort = 0
	maxPort = 65535
)

type operation struct {
	managerIP           string
	managerHTTPSPort    string
	managerCertFilePath string
	managerKeyFilePath  string
	managerCAFilePath   string
	managerSkipVerify   bool
}

func (o operation) isValid() bool {
	if len(o.managerIP) != 0 && !validPort(o.managerHTTPSPort) {
		return false
	}
	return len(o.managerIP) != 0
}

func (o *operation) cleanSpaceCharacter() {
	o.managerIP = strings.ReplaceAll(o.managerIP, " ", "")
	o.managerHTTPSPort = strings.ReplaceAll(o.managerHTTPSPort, " ", "")
	o.managerCertFilePath = strings.ReplaceAll(o.managerCertFilePath, " ", "")
	o.managerKeyFilePath = strings.ReplaceAll(o.managerKeyFilePath, " ", "")
}

func (o operation) constructNewConfig() kappctl.Config {
	cfg := kappctl.Config{}
	cfg.ManagerHTTPSServer = fmt.Sprintf("https://%s:%s", o.managerIP, o.managerHTTPSPort)
	cfg.ManagerClientCertificateData = file.ReadFileToBase64(o.managerCertFilePath)
	cfg.ManagerClientKeyData = file.ReadFileToBase64(o.managerKeyFilePath)
	cfg.ManagerCA = file.ReadFileToBase64(o.managerCAFilePath)
	cfg.ManagerSkipVerify = o.managerSkipVerify
	return cfg
}

// NewCommand create config command
func NewCommand() *cobra.Command {
	o := operation{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config kappctl with Kappital-Manager",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunCheck(&o)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}

	kappctl.ManagerIPAddr.AddStringFlag(&o.managerIP, cmd)
	kappctl.ManagerHTTPSPort.AddStringFlag(&o.managerHTTPSPort, cmd)
	kappctl.ManagerClientCertFile.AddStringFlag(&o.managerCertFilePath, cmd)
	kappctl.ManagerClientKeyFile.AddStringFlag(&o.managerKeyFilePath, cmd)
	kappctl.ManagerCAFile.AddStringFlag(&o.managerCAFilePath, cmd)
	kappctl.ManagerSkipVerify.AddBoolFlag(&o.managerSkipVerify, cmd)
	return cmd
}

func preRunCheck(o *operation) error {
	o.cleanSpaceCharacter()
	if !o.isValid() {
		return fmt.Errorf("the parameter is invalid")
	}
	return nil
}

func run(o operation) error {
	return writeConfigFile(o.constructNewConfig())
}

func writeConfigFile(cfg kappctl.Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get user home directory: %s", err)
	}
	path := filepath.Join(home, kappctl.ConfigFile)
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create direcory %s to store config file: %s", filepath.Dir(path), err)
	}
	if file.IsFileExist(path) {
		if err = os.RemoveAll(path); err != nil {
			return fmt.Errorf("cannot remove old file, err: %v", err)
		}
	}
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, jsonData, 0600); err != nil {
		return fmt.Errorf("write data into file %v: %v", path, err)
	}
	return nil
}

func validPort(port string) bool {
	if len(port) == 0 {
		return false
	}
	num, err := strconv.ParseInt(port, 10, strconv.IntSize)
	if err != nil {
		return false
	}
	return minPort <= num && num <= maxPort
}
