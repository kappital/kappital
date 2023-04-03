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

package flowcontroller

import (
	"sync"

	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"
)

var rateLimiter flowcontrol.RateLimiter
var once sync.Once

type Config struct {
	QPS    float64
	Burst  int
	Enable bool
}

func DefaultFlowControllerConfig() *Config {
	return &Config{
		QPS:    10,
		Burst:  30,
		Enable: true,
	}
}

func Init(fcInit *Config) {
	if rateLimiter == nil {
		// init flow controller bucket
		once.Do(func() {
			rateLimiter = flowcontrol.NewTokenBucketRateLimiter(float32(fcInit.QPS), fcInit.Burst)
		})
	}
}

func TryToPassReq() bool {
	if !rateLimiter.TryAccept() {
		klog.Errorf("flow controller tokens exhaust, please try later")
		return false
	}
	return true
}
