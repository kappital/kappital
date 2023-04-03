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

package uuid

import (
	"sync"

	"github.com/pborman/uuid"
)

var uidLock sync.Mutex
var lastUUID uuid.UUID

const truncatedSize = 8

// NewUUID create a new uuid
func NewUUID() string {
	uidLock.Lock()
	defer uidLock.Unlock()
	id := uuid.NewUUID()
	for uuid.Equal(lastUUID, id) {
		id = uuid.NewUUID()
	}
	lastUUID = id
	return id.String()
}

// NewUUID8 get the first 8 digit if the uuid
func NewUUID8() string {
	id := NewUUID()
	return id[0:truncatedSize]
}
