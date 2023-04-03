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

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/kappctl/cmd/config"
	"github.com/kappital/kappital/pkg/kappctl/cmd/create"
	"github.com/kappital/kappital/pkg/kappctl/cmd/delete"
	"github.com/kappital/kappital/pkg/kappctl/cmd/get"
	"github.com/kappital/kappital/pkg/kappctl/cmd/initiate"
)

// NewKappctlCmd create the kappctl root command
func NewKappctlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kappctl",
		Short: "kappctl manages the full lifecycle of Kappital Resources",
	}

	cmd.AddCommand(initiate.NewCommand())
	cmd.AddCommand(config.NewCommand())
	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(create.NewCommand())
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewKappctlCmd()
	cobra.CheckErr(cmd.Execute())
}
