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

package audit

import (
	"os"

	"github.com/sirupsen/logrus"
)

type traceRating string
type TraceType string

const (
	// NormalRating of trace
	NormalRating traceRating = "Normal"
	// WarningRating of trace
	WarningRating traceRating = "Warning"
	// IncidentRating of trace
	IncidentRating traceRating = "Incident"

	// ConsoleActionType of trace
	ConsoleActionType TraceType = "ConsoleAction"
	// APICallType of trace
	APICallType TraceType = "ApiCall"
	// SystemAction of trace
	SystemAction TraceType = "SystemAction"
)

// AuditLogInfo of each audit basic messages
type AuditLogInfo struct {
	Timestamp    int64
	SourceIP     string
	ResourceType string
	ResourceName string
	TraceName    string
	TraceRating  traceRating
	TraceType    TraceType
	Message      string
}

func (a AuditLogInfo) getLogrusFields() logrus.Fields {
	return logrus.Fields{
		"timestamp":     a.Timestamp,
		"source_ip":     a.SourceIP,
		"resource_type": a.ResourceType,
		"resource_name": a.ResourceName,
		"trace_name":    a.TraceName,
		"trace_rating":  a.TraceRating,
		"trace_type":    a.TraceType,
	}
}

// AuditLogConfig of audit
type AuditLogConfig struct {
	AppName  string
	Filename string
	MaxSize  int64
}

// DefaultAuditLogConfig get the default audit log config
func DefaultAuditLogConfig() AuditLogConfig {
	return AuditLogConfig{
		AppName:  os.Getenv("APP_NAME"),
		Filename: "/opt/kappital/audit/audit.log", // Default audit log location
		MaxSize:  1024 * 1024,                     // Default the audit will be 1M
	}
}
