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
	"os"
	"time"
)

// DatabaseConfig the basic config for the database
type DatabaseConfig struct {
	SQLDriver   string
	MaxIdle     int
	MaxConn     int
	MaxLifetime int
	SslEnable   string
}

// DefaultDatabaseConfiguration get the default database configuration
func DefaultDatabaseConfiguration() *DatabaseConfig {
	return &DatabaseConfig{
		SQLDriver:   DefaultSQLDriverName,
		MaxIdle:     10,
		MaxConn:     256,
		MaxLifetime: 1800,
		SslEnable:   os.Getenv(SslEnableKey),
	}
}

// DatabaseWatcherConfig configuration for watch the database connection
type DatabaseWatcherConfig struct {
	// database uri
	Connection string
	// min interval seconds for database reconnection
	ListenerMinReconnectInterval time.Duration
	// max interval seconds for database reconnection
	ListenerMaxReconnectInterval time.Duration
}

// DefaultDatabaseWatcherConfig get the default database watcher config
func DefaultDatabaseWatcherConfig() *DatabaseWatcherConfig {
	return &DatabaseWatcherConfig{
		ListenerMinReconnectInterval: 1 * time.Second,
		ListenerMaxReconnectInterval: 20 * time.Second,
	}
}
