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

package gateway

import "net/http"

// FakeResponseWriter for the unit test
type FakeResponseWriter struct{}

// Header fake header method
func (f FakeResponseWriter) Header() http.Header {
	return nil
}

// Write fake method
func (f FakeResponseWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

// WriteHeader fake method
func (f FakeResponseWriter) WriteHeader(_ int) {}
