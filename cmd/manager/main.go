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

package main

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/cmd/options"
	"github.com/kappital/kappital/pkg/models"
	"github.com/kappital/kappital/pkg/processor"
	"github.com/kappital/kappital/pkg/routers/flowcontroller"
	"github.com/kappital/kappital/pkg/routers/manager"
	"github.com/kappital/kappital/pkg/utils/audit"
	"github.com/kappital/kappital/pkg/utils/version"
	"github.com/kappital/kappital/pkg/watcher"
)

func main() {
	if len(os.Args) > 1 && version.Flags.Has(os.Args[1]) {
		fmt.Println(version.Get(version.ServiceNameManager).String())
		os.Exit(0)
	}
	if err := audit.InitAuditLog(audit.DefaultAuditLogConfig()); err != nil {
		klog.Fatalf("cannot init the audit log, err: %s", err)
	}

	cfg, err := options.NewServerRunOptions(version.ServiceNameManager)
	if err != nil {
		klog.Fatalf("failed to configure server running options: %s", err)
	}
	manager.InitRouters()
	flowcontroller.Init(cfg.FlowControllerConfig)
	// init sql driver
	if err = models.GetDatabase().InitSQLDriver(cfg.DBConfig, models.Manager); err != nil {
		klog.Fatalf("failed to initialize sql driver, error: %v", err)
	}

	notifyWatcher := watcher.NewWatcher()
	for _, proc := range processor.GetProcesses() {
		if err = notifyWatcher.Watch(proc.ProcessObject(), proc.ProcessType(), proc); err != nil {
			klog.Fatalf("create watch for processor %s failed, err: %s", proc.ProcessType(), err)
		}
	}
	// start modules
	jobStopCh := make(chan struct{})
	// processor modules
	processor.StartAllProcessors(jobStopCh)
	if err = notifyWatcher.StartProcessor(); err != nil {
		klog.Fatalf("start processor failed, error: %s", err)
	}
	web.Run()
	close(jobStopCh)
	klog.Info("kappital-manager server stopped")
}
