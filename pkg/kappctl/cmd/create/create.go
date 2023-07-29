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

package create

import (
	"github.com/spf13/cobra"

	"github.com/kappital/kappital/pkg/kappctl/cmd/create/instance"
	"github.com/kappital/kappital/pkg/kappctl/cmd/create/service"
)

// NewCommand create command
func NewCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Kappital resource",
	}

	createCmd.AddCommand(instance.Cmd.NewCommand())
	createCmd.AddCommand(service.Cmd.NewCommand())

	return createCmd
}
