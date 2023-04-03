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

package models

import (
	"errors"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type serviceType string

const (
	// AliasName of the database
	AliasName = "default"
	// SslEnableKey the env key of the database ssl
	SslEnableKey = "SSL_ENABLE"

	// DefaultSQLDriverName the default sql driver
	DefaultSQLDriverName = "sqlite"
)

const (
	// StatusInstalling for operator
	StatusInstalling = "Installing"
	// StatusSuccess for operator
	StatusSuccess = "Success"
	// StatusAbnormal for operator
	StatusAbnormal = "Abnormal"

	// StatusUnknown for instance
	StatusUnknown = "Unknown"
	// StatusInitializing for instance
	StatusInitializing = "Initializing"
	// StatusInitialized for instance
	StatusInitialized = "Initialized"
	// StatusInitFailed for instance
	StatusInitFailed = "InitFailed"
	// StatusFailed for operator and instance
	StatusFailed = "Failed"

	// StatusUpgrading for instance operator
	StatusUpgrading = "Upgrading"
	// StatusRollingBack for instance operator
	StatusRollingBack = "RollingBack"
	// StatusDeleting for instance operator
	StatusDeleting = "Deleting"
	// StatusUpgradeFailed for instance operator
	StatusUpgradeFailed = "UpgradeFailed"
	// StatusDeleteFailed for instance operator
	StatusDeleteFailed = "DeleteFailed"
	// StatusRollBackFailed for instance operator
	StatusRollBackFailed = "RollBackFailed"
)

// FailedStatusList of instance and operator
var FailedStatusList = []string{StatusUpgradeFailed, StatusDeleteFailed, StatusRollBackFailed, StatusFailed, StatusInitFailed}

const (
	// Manager the service type name of manager service
	Manager serviceType = "Manager"
)

const (
	noLastInsertIDAvailableErr  = "no LastInsertId available"
	lastInsertIDNotSupportedErr = "LastInsertId is not supported by this driver"
)

// NotDBInsertIDError because insert to the database it may have some id errors, we need to ignore this kind of errors
func NotDBInsertIDError(err error) bool {
	return err != nil && !strings.Contains(err.Error(), noLastInsertIDAvailableErr) &&
		!strings.Contains(err.Error(), lastInsertIDNotSupportedErr) &&
		!errors.Is(err, orm.ErrLastInsertIdUnavailable)
}

// IgnoreDBInsertIDError because insert to the database it may have some id errors, we need to ignore this kind of errors
func IgnoreDBInsertIDError(err error) error {
	if NotDBInsertIDError(err) {
		return err
	}
	return nil
}
