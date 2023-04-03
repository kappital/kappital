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

package errors

import (
	"fmt"
	"net/http"
)

type modulePrefix uint

const (
	errPrefix                           = "KAPPITAL."
	commonErrCode          modulePrefix = 100
	serviceErrCode         modulePrefix = 101
	serviceInstanceErrCode modulePrefix = 102
)

var (
	errArray []*errorImpl

	// ErrDataUnmarshal will happen when translate the string to the struct, or cannot translate structure from one type to the other.
	ErrDataUnmarshal = newKappError(commonErrCode, http.StatusBadRequest, 1, "Data unmarshal error.")

	// ErrServiceInstall has some problem for CloudNativeService deploying failed. May because of cluster disconnection, or cluster limitation problems.
	ErrServiceInstall = newKappError(serviceErrCode, http.StatusBadRequest, 2, "Service install error.")
	// ErrServiceParam cannot analysis or get the param from request
	ErrServiceParam = newKappError(serviceErrCode, http.StatusBadRequest, 3, "Parameter is invalid.")
	// ErrServiceDelete cannot delete the service in cluster
	ErrServiceDelete = newKappError(serviceErrCode, http.StatusBadRequest, 4, "ServiceBinding delete error.")

	// ErrServiceInstanceCreate cannot deploy the user's instance into cluster
	ErrServiceInstanceCreate = newKappError(serviceInstanceErrCode, http.StatusInternalServerError, 1, "Service Instance create error.")
)

// KappError the error that will be used in the manager
type KappError interface {
	error
	GetHTTPCode() int
	GetErrorCode() string
	GetResp() ErrorResp
	WrapErrorReasonWith(format string, args ...interface{}) KappError
	TypeEqual(err KappError) bool
}

// ErrorResp error response
type ErrorResp struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"errorMsg"`
	Reason    string `json:"reason,omitempty"`
}

type errorImpl struct {
	module        modulePrefix
	httpErrorCode int
	errorIndex    int
	errMsg        string
	reason        string
}

// newKappError create the new error for manager
func newKappError(module modulePrefix, httpCode, index int, msg string) KappError {
	err := &errorImpl{
		module:        module,
		httpErrorCode: httpCode,
		errorIndex:    index,
		errMsg:        msg,
	}
	errArray = append(errArray, err)
	return err
}

// Error output the error message to string
func (e *errorImpl) Error() string {
	if e.reason == "" {
		return e.errMsg
	}
	return fmt.Sprintf("%s reason:%s", e.errMsg, e.reason)
}

// GetHTTPCode get the current error http code
func (e *errorImpl) GetHTTPCode() int {
	return e.httpErrorCode
}

// GetErrorCode get the error code
func (e *errorImpl) GetErrorCode() string {
	return fmt.Sprintf("%s%04d%04d", errPrefix, e.module, e.errorIndex)
}

// WrapErrorReasonWith add the error reasons to the current error
func (e *errorImpl) WrapErrorReasonWith(format string, args ...interface{}) KappError {
	e.reason = fmt.Sprintf(format, args...)
	return e
}

// GetResp get the response message and error
func (e *errorImpl) GetResp() ErrorResp {
	return ErrorResp{
		ErrorCode: fmt.Sprintf("%s%04d%04d", errPrefix, e.module, e.errorIndex),
		Message:   e.errMsg,
		Reason:    e.reason,
	}
}

// TypeEqual does the error type is the same
func (e *errorImpl) TypeEqual(err KappError) bool {
	return e.GetErrorCode() == err.GetErrorCode()
}
