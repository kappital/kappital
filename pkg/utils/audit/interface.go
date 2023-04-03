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
	"time"
)

var cmd AuditLog

// AuditLog interface for expand
type AuditLog interface {
	initConfig(config AuditLogConfig) error
	info(info AuditLogInfo)
	error(info AuditLogInfo)
	fault(info AuditLogInfo)
}

func init() {
	cmd = &defaultLog{}
}

// SetAuditLog with specific logs
func SetAuditLog(l AuditLog) {
	cmd = l
}

// InitAuditLog from config
func InitAuditLog(config AuditLogConfig) error {
	return cmd.initConfig(config)
}

// Info log of audit
func Info(info AuditLogInfo) {
	if len(info.SourceIP) == 0 { // if ip is empty, means the inner actions
		info.SourceIP = "localhost"
	}
	info.Timestamp = time.Now().Unix()
	info.TraceRating = NormalRating
	cmd.info(info)
}

// Error log of audit
func Error(info AuditLogInfo) {
	if len(info.SourceIP) == 0 { // if ip is empty, means the inner actions
		info.SourceIP = "localhost"
	}
	info.Timestamp = time.Now().Unix()
	info.TraceRating = WarningRating
	cmd.info(info)
}

// Fault log of audit
func Fault(info AuditLogInfo) {
	if len(info.SourceIP) == 0 { // if ip is empty, means the inner actions
		info.SourceIP = "localhost"
	}
	info.Timestamp = time.Now().Unix()
	info.TraceRating = IncidentRating
	cmd.info(info)
}
