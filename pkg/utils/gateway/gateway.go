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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// RequestInfo of the http or https request
type RequestInfo struct {
	Method       string
	Path         string
	Body         interface{}
	IsFile       bool
	HeaderAdder  map[string]string
	HeaderSetter map[string]string
	CaCrt        string
	ClientCrt    string
	ClientKey    string
	Skip         bool
}

func (r RequestInfo) getRequest() (*http.Request, error) {
	var req *http.Request
	var err error
	if r.Body != nil {
		req, err = getHTTPRequestWithBody(r.Method, r.Path, r.Body, r.IsFile)
	} else {
		req, err = http.NewRequest(r.Method, r.Path, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("create http client failed, err: %w", err)
	}
	for k, v := range r.HeaderAdder {
		req.Header.Add(k, v)
	}
	for k, v := range r.HeaderSetter {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (r RequestInfo) getClient() (*http.Client, error) {
	if len(r.CaCrt) == 0 || strings.HasPrefix(r.Path, "http://") {
		return &http.Client{}, nil
	}
	pool := x509.NewCertPool()
	caCrt, err := base64.StdEncoding.DecodeString(r.CaCrt)
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)
	clientCrt, err := base64.StdEncoding.DecodeString(r.ClientCrt)
	if err != nil {
		return nil, err
	}
	clientKey, err := base64.StdEncoding.DecodeString(r.ClientKey)
	if err != nil {
		return nil, err
	}
	cliCrt, err := tls.X509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				Certificates:       []tls.Certificate{cliCrt},
				InsecureSkipVerify: r.Skip, //nolint:gosec
			},
		},
	}, nil
}

// CommonUtilRequest common get http/https request
func CommonUtilRequest(info *RequestInfo) (int, []byte, error) {
	req, err := info.getRequest()
	if err != nil {
		return 0, nil, err
	}
	client, err := info.getClient()
	if err != nil {
		return 0, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("perform request failed, err: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck
	code := resp.StatusCode
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("read response body failed, err: %s", err)
	}
	return code, buf, nil
}

// GetLocalIP returns the non loop back local IP of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("can't find Interface IP")
}

func getHTTPRequestWithBody(method, path string, body interface{}, isFile bool) (*http.Request, error) {
	if !isFile {
		var byteBody []byte
		byteBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body failed, err: %w", err)
		}
		return http.NewRequest(method, path, bytes.NewReader(byteBody))
	}

	buf, ok := body.(*bytes.Buffer)
	if !ok {
		return nil, fmt.Errorf("the body struct only accetp bytes.Buffer")
	}
	return http.NewRequest(method, path, buf)
}

// ReplaceIP address of org
func ReplaceIP(org, newIP, defaultPort string) string {
	arr := strings.Split(org, ":")
	if len(arr) == 2 {
		return fmt.Sprintf("%s:%s", newIP, arr[1])
	}
	return fmt.Sprintf("%s:%s", newIP, defaultPort)
}
