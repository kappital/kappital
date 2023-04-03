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
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type defaultLog struct {
	log *logrus.Logger
}

type defaultHook struct {
	appName string
}

// Levels of the hook which watch
func (d *defaultHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel}
}

// Fire the entry.Data constance attribute
func (d *defaultHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = d.appName
	return nil
}

func (d *defaultLog) initConfig(config AuditLogConfig) error {
	if len(config.AppName) == 0 {
		return fmt.Errorf("missing the app name, will not use the audit log, and stop the running")
	}
	hook := &defaultHook{appName: config.AppName}
	d.log = logrus.New()
	// set the output format as json, and not use logrus timestamp format
	d.log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
	d.log.SetLevel(logrus.InfoLevel)
	d.log.AddHook(hook)

	if len(config.Filename) == 0 {
		d.log.SetOutput(os.Stdout)
		return nil
	}
	file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	d.log.SetOutput(file)
	return nil
}

func (d defaultLog) info(info AuditLogInfo) {
	d.log.WithFields(info.getLogrusFields()).Infoln(info.Message)
}

func (d defaultLog) error(info AuditLogInfo) {
	d.log.WithFields(info.getLogrusFields()).Errorln(info.Message)
}

func (d defaultLog) fault(info AuditLogInfo) {
	d.log.WithFields(info.getLogrusFields()).Fatalln(info.Message)
}
