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

package initiate

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/kappital/kappital/pkg/apis"
	"github.com/kappital/kappital/pkg/kappctl"
	"github.com/kappital/kappital/pkg/utils/file"
)

type operation struct {
	dirPath string
	name    string
	version string
}

var (
	// cmd singleton pattern of create the demo kappital service package
	cmd operation

	//go:embed kappital-demo
	demo embed.FS
)

const (
	demoName     = "kappital-demo"
	metadataFile = "metadata.yaml"
	maxDepth     = 5
)

func (o operation) getArgumentMap() map[string]interface{} {
	return map[string]interface{}{
		kappctl.FileName.GetFlagName():       o.name,
		kappctl.PackageVersion.GetFlagName(): o.version,
	}
}

// NewCommand get the "init" command
func NewCommand() *cobra.Command {
	return cmd.NewCommand()
}

// NewCommand create the new command for create kappital demo package for user
func (o *operation) NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "init",
		Short: "Create a Kappital package scaffold from scratch.",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return o.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.RunE()
		},
	}
	kappctl.FileName.AddStringFlag(&o.name, command)
	kappctl.PackageVersion.AddStringFlag(&o.version, command)
	return command
}

// PreRunE run before creating kappital demo package for user, check does the arguments has some problem or not
func (o *operation) PreRunE(args []string) error {
	if len(args) == 1 {
		o.dirPath = args[0]
	} else {
		o.dirPath = "."
	}
	return kappctl.IsInputValidate(o.getArgumentMap())
}

// RunE create kappital demo package for user
func (o *operation) RunE() error {
	// create a dirPath with the o.name at o.dirPath
	absPath, err := filepath.Abs(o.dirPath)
	if err != nil {
		return err
	}
	dirPath := filepath.Clean(filepath.Join(absPath, o.name))
	if err = os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	// using the embed file system to open the current directory
	entries, err := demo.ReadDir(demoName)
	if err != nil {
		return err
	}
	if err = o.createPackage(entries, dirPath, demoName, 0); err != nil {
		return err
	}
	fmt.Printf("init service %s package success.\n", o.name)
	return nil
}

func (o *operation) createPackage(entries []fs.DirEntry, rootPath, pkgPath string, deepLevel int) error {
	if deepLevel >= maxDepth {
		return fmt.Errorf("the directory is too deep")
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if err := createFile(rootPath, pkgPath, entry.Name(), o); err != nil {
				return err
			}
			continue
		}
		newPkgPath := filepath.Clean(filepath.Join(pkgPath, entry.Name()))
		newRootPath := filepath.Clean(filepath.Join(rootPath, entry.Name()))
		newEntries, err := demo.ReadDir(newPkgPath)
		if err != nil {
			return err
		}
		if err = o.createPackage(newEntries, newRootPath, newPkgPath, deepLevel+1); err != nil {
			return err
		}
	}
	return nil
}

func createFile(rootPath, pkgPath, fileName string, o *operation) error {
	b, err := demo.ReadFile(filepath.Clean(filepath.Join(pkgPath, fileName)))
	if err != nil {
		return err
	}
	filePath := filepath.Clean(filepath.Join(rootPath, fileName))
	if !file.IsDirExist(rootPath) {
		if err := os.MkdirAll(rootPath, 0750); err != nil {
			return err
		}
		fmt.Printf("create directory: %s\n", rootPath)
	}

	destFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil || destFile == nil {
		return fmt.Errorf("cannot create the file")
	}
	defer func() {
		err = destFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if fileName == metadataFile {
		b, err = modifyMetadataFile(b, o.name, o.version)
		if err != nil {
			return err
		}
	}
	_, err = destFile.Write(b)
	fmt.Printf("create file: %s\n", fileName)
	return err
}

func modifyMetadataFile(b []byte, name, _ string) ([]byte, error) {
	var metadata apis.Descriptor
	if err := yaml.Unmarshal(b, &metadata); err != nil {
		return nil, err
	}
	metadata.Name = name
	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	body, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return nil, err
	}
	return body, nil
}
