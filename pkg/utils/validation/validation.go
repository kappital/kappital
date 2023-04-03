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

package validation

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	expr1 = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	expr2 = "^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"
)

// ValidCommonString generate method to check and return the valid string
func ValidCommonString(passIn string, min, max int, canEmpty bool, defaultValue string) (string, error) {
	if checkString(passIn, min, max) {
		return passIn, nil
	}
	if canEmpty {
		return "", nil
	}
	if checkString(defaultValue, min, max) {
		return defaultValue, nil
	}
	return "", fmt.Errorf("the passed in value and default value both are not valid")
}

// ValidBool generate the pass in string as the true or false
func ValidBool(passIn string) (bool, error) {
	res, err := strconv.ParseBool(passIn)
	if err != nil {
		return false, fmt.Errorf("the pass in valuse %s cannot transfer to boolean, err: %v", passIn, err)
	}
	return res, nil
}

func validateLength(str string, min, max int) bool {
	return min <= len(str) && len(str) <= max
}

func checkString(str string, min, max int) bool {
	if !validateLength(str, min, max) {
		return false
	}
	reg1, err := regexp.Compile(expr1)
	if err != nil {
		return false
	}
	reg2, err := regexp.Compile(expr2)
	if err != nil {
		return false
	}
	return reg1.Match([]byte(str)) || reg2.Match([]byte(str))
}
