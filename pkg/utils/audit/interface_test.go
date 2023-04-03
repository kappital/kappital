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
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func TestAuditLog(t *testing.T) {
	config := DefaultAuditLogConfig()
	config.AppName = "Unit Test"
	config.Filename = ""
	if err := InitAuditLog(config); err != nil {
		t.Errorf("InitAuditLog error should be nil, but: %v", err)
		return
	}
	p := gomonkey.ApplyFunc(os.Exit, func(_ int) {})
	defer p.Reset()
	Info(AuditLogInfo{})
	Error(AuditLogInfo{})
	Fault(AuditLogInfo{})

	SetAuditLog(&FakeAuditLogger)
	if err := InitAuditLog(config); err != nil {
		t.Errorf("InitAuditLog error should be nil, but: %v", err)
		return
	}
	Info(AuditLogInfo{})
	Error(AuditLogInfo{})
	Fault(AuditLogInfo{})
}
