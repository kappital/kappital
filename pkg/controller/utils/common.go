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

package utils

import (
	"fmt"

	"github.com/beego/beego/v2/server/web/context"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/utils/audit"
	"github.com/kappital/kappital/pkg/utils/errors"
)

var (
	// ErrIllegalParameters pass in from url or body struct
	ErrIllegalParameters = fmt.Errorf("illegal parameters")
)

type action string

const (
	// QueryAction of the manager which search service packages, service instances or service bindings
	QueryAction action = "Query"
	// DeployAction of manager which deploy service binding or service instance
	DeployAction action = "Deploy"
	// UninstallAction of manager which uninstall service binding or service instance
	UninstallAction action = "Uninstall"
)

// AuditLog write the audit log from the defer method, and the detail dependents on the error.
// Because only use the kappctl binary tool or CURL APIs to visit manager,
// thus, it only offers APICallType of trace.
func AuditLog(ctx *context.Context, traceName string, resourceType action, name *string, err *error) {
	if (*err) != nil {
		audit.Error(audit.AuditLogInfo{
			SourceIP:     ctx.Request.RemoteAddr,
			ResourceType: string(resourceType),
			ResourceName: *name,
			TraceName:    traceName,
			TraceType:    audit.APICallType,
			Message:      (*err).Error(),
		})
	} else {
		audit.Info(audit.AuditLogInfo{
			SourceIP:     ctx.Request.RemoteAddr,
			ResourceType: string(resourceType),
			ResourceName: *name,
			TraceName:    traceName,
			TraceType:    audit.APICallType,
			Message:      "success",
		})
	}
}

// ReplyJSON sends json reply to http client
func ReplyJSON(ctx *context.Context, stateCode int, resp interface{}) {
	var msg interface{}
	if resp != nil {
		switch resp.(type) {
		case error, errors.KappError:
			msg = resp.(error).Error()
		default:
			msg = resp
		}
	}

	ctx.Output.SetStatus(stateCode)
	err := ctx.Output.JSON(msg, false, false)
	if err != nil {
		klog.Errorf("failed to write json resp, err: %v", err)
	}
}
