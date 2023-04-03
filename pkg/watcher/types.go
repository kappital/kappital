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

package watcher

type channelConfig struct {
	// the channel name listening on
	name string
	// target object
	target interface{}
	// processor contains a list function used to get all target objects
	// when startup and a process function to handle notification
	processor Processor
}

// Notification structure containing event notifications.
// OPType means the type of event, RawData means the object content of event
type Notification struct {
	OPType  string `json:"op_type"`
	RawData string `json:"raw_data"`
}
