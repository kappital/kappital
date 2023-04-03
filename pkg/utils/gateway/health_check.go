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

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/klog/v2"
)

// HealthAndReadinessProvider get and create the health and readiness provider
func HealthAndReadinessProvider(address string, cert, key []byte) {
	http.HandleFunc("/healthz", setHTTPHeaderToOK)
	http.HandleFunc("/readyz", setHTTPHeaderToOK)
	server, err := constructServer(address, cert, key)
	if err != nil {
		klog.Errorf("cannot get server")
		os.Exit(1)
	}

	if err = server.ListenAndServeTLS("", ""); err != nil {
		klog.Errorf("cannot start health check")
		os.Exit(1)
	}
}

func setHTTPHeaderToOK(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "OK\n"); err != nil {
		klog.Errorf("cannot set http header to ok which host is %s, err: %v", r.Host, err)
	}
}

func constructServer(address string, cert, key []byte) (*http.Server, error) {
	if tlsConfig, err := constructTLSConfig(cert, key); err != nil {
		klog.Errorf("cannot construct tls config, err: %s", err)
		return nil, err
	} else {
		return &http.Server{
			Addr:              address,
			Handler:           http.DefaultServeMux,
			TLSConfig:         tlsConfig,
			ReadHeaderTimeout: time.Minute * 2,
		}, nil
	}
}

func constructTLSConfig(cert, key []byte) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		MinVersion: tls.VersionTLS12,
	}, nil
}
