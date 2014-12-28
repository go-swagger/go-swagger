// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           sigu-399
// author-github    https://github.com/sigu-399
// author-mail      sigu.399@gmail.com
//
// repository-name  jsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description      Various utility functions.
//
// created          26-02-2013

package jsonschema

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
)

func mustBeInteger(what interface{}) *int {
	var number int
	if isKind(what, reflect.Float64) {
		fnumber := what.(float64)
		if isFloat64AnInteger(fnumber) {
			number = int(fnumber)
			return &number
		} else {
			return nil
		}
	} else if isKind(what, reflect.Int) {
		number = what.(int)
		return &number
	}
	return nil
}

func mustBeNumber(what interface{}) *float64 {
	var number float64

	if isKind(what, reflect.Float64) {
		number = what.(float64)
		return &number
	} else if isKind(what, reflect.Int) {
		number = float64(what.(int))
		return &number
	}
	return nil
}

func isKind(what interface{}, kind reflect.Kind) bool {
	return reflect.ValueOf(what).Kind() == kind
}

func existsMapKey(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

func isStringInSlice(s []string, what string) bool {
	for i := range s {
		if s[i] == what {
			return true
		}
	}
	return false
}

// same as ECMA Number.MAX_SAFE_INTEGER and Number.MIN_SAFE_INTEGER
const (
	max_json_float = float64(1<<53 - 1)  // 9007199254740991.0 	 	 2^53 - 1
	min_json_float = -float64(1<<53 - 1) //-9007199254740991.0	-2^53 - 1
)

// allow for integers [-2^53, 2^53-1] inclusive
func isFloat64AnInteger(f float64) bool {

	if math.IsNaN(f) || math.IsInf(f, 0) || f < min_json_float || f > max_json_float {
		return false
	}

	return f == float64(int64(f)) || f == float64(uint64(f))
}

// formats a number so that it is displayed as the smallest string possible
func validationErrorFormatNumber(n float64) string {

	if isFloat64AnInteger(n) {
		return fmt.Sprintf("%d", int64(n))
	}

	return fmt.Sprintf("%g", n)
}

func marshalToJsonString(value interface{}) (*string, error) {

	mBytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	sBytes := string(mBytes)
	return &sBytes, nil
}

const internalLogEnabled = false

func internalLog(message string) {
	if internalLogEnabled {
		log.Print(message)
	}
}
