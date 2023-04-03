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

package routers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"k8s.io/klog/v2"

	"github.com/kappital/kappital/pkg/constants"
	"github.com/kappital/kappital/pkg/routers/flowcontroller"
)

var validMethodSet = map[string]struct{}{http.MethodGet: {}, http.MethodDelete: {}, http.MethodPost: {}}

const (
	checkIdentityEnv            = "CHECK_IDENTITY"
	acceptCertificateCommonName = "Kappital - Client"
)

// InitFilters for url, and pre-check the requests
func InitFilters() {
	web.InsertFilter("/api/*", web.BeforeStatic, formatFilter)
	web.InsertFilter("/api/*", web.BeforeStatic, beforeStaticFilter)
	web.InsertFilter("/api/*", web.BeforeExec, flowControlFilter)

	web.InsertFilter("/*", web.BeforeStatic, formatFilter)
	web.InsertFilter("/*", web.BeforeStatic, beforeStaticFilter)
	web.InsertFilter("/*", web.BeforeExec, flowControlFilter)

	if check, err := strconv.ParseBool(os.Getenv(checkIdentityEnv)); check && err == nil {
		klog.Info("Open the Identity Check")
		web.InsertFilter("/api/*", web.BeforeExec, checkIdentity)
		web.InsertFilter("/*", web.BeforeStatic, checkIdentity)
	}
}

func beforeStaticFilter(ctx *context.Context) {
	if !strings.HasPrefix(ctx.Request.Header.Get("Content-Type"), "multipart/form-data") {
		return
	}

	if ctx.Request.ContentLength > constants.FileSize {
		errMsg := fmt.Sprintf("Uploaded Content Length Exceed %dM", constants.FileSize/1024/1024)
		klog.Errorf(errMsg)
		setFilterErrorMsg(ctx, http.StatusBadRequest, errMsg)
	}
}

func formatFilter(ctx *context.Context) {
	rawURL := ctx.Input.URL()
	method := ctx.Input.Method()
	_, find := validMethodSet[method]
	if !find {
		klog.Errorf("[FormatURL] illegal Method! method: %s", method)
		setFilterErrorMsg(ctx, http.StatusNotFound, "method is illegal")
		return
	}
	if int64(len(ctx.Input.RequestBody)) > constants.FileSize {
		klog.Errorf("[FormatURL] illegal Request Body Size! current size %d, max size %d", len(ctx.Input.RequestBody), constants.FileSize)
		setFilterErrorMsg(ctx, http.StatusBadRequest, "request body size is too large")
		return
	}
	formatURL := path.Clean(rawURL)
	// note: only formatURL allowed ( trailing slash allowed for compatibility)
	if (rawURL != formatURL && rawURL != formatURL+"/") || !path.IsAbs(formatURL) {
		klog.Errorf("[FormatURL] illegal URL format! rawURL:%s %s(%s)", method, rawURL, formatURL)
		setFilterErrorMsg(ctx, http.StatusNotFound, "url is illegal")
	}
}

func setFilterErrorMsg(ctx *context.Context, code int, msg string) {
	ctx.ResponseWriter.WriteHeader(code)
	ctx.WriteString(msg)
}

// Check whether the server could handle the request
func flowControlFilter(ctx *context.Context) {
	if !flowcontroller.TryToPassReq() {
		klog.Warningf("request(%s %s) is rejected by flow controller, too many request", ctx.Input.Method(), ctx.Input.URI())
		setFilterErrorMsg(ctx, http.StatusTooManyRequests, "too many request")
		return
	}
}

func checkIdentity(ctx *context.Context) {
	var find bool
	for _, cert := range ctx.Request.TLS.PeerCertificates {
		if cert.Subject.CommonName == acceptCertificateCommonName {
			find = true
			break
		}
	}
	if !find {
		setFilterErrorMsg(ctx, http.StatusUnauthorized, "reject certificates")
	}
}
