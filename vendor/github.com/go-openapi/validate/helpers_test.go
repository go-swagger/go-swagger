// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpers_addPointerError(t *testing.T) {
	res := new(Result)
	r := errorHelp.addPointerError(res, fmt.Errorf("my error"), "my ref", "path")
	msg := r.Errors[0].Error()
	assert.Contains(t, msg, "could not resolve reference in path to $ref my ref: my error")
}

// Test cases in private method asInt64()
// including expected panic() cases
func TestHelpers_asInt64(t *testing.T) {
	var r int64
	r = valueHelp.asInt64(int(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(uint(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(int8(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(uint8(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(int16(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(uint16(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(int32(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(uint32(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(int64(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(uint64(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(float32(3))
	assert.Equal(t, int64(3), r)
	r = valueHelp.asInt64(float64(3))
	assert.Equal(t, int64(3), r)

	// Non numeric
	//assert.PanicsWithValue(t, "Non numeric value in asInt64()", func() {
	//	valueHelp.asInt64("123")
	//})
	if assert.NotPanics(t, func() {
		valueHelp.asInt64("123")
	}) {
		assert.Equal(t, valueHelp.asInt64("123"), (int64)(0))
	}
}

// Test cases in private method asUint64()
// including expected panic() cases
func TestHelpers_asUint64(t *testing.T) {
	var r uint64
	r = valueHelp.asUint64(int(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(uint(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(int8(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(uint8(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(int16(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(uint16(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(int32(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(uint32(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(int64(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(uint64(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(float32(3))
	assert.Equal(t, uint64(3), r)
	r = valueHelp.asUint64(float64(3))
	assert.Equal(t, uint64(3), r)

	// Non numeric
	//assert.PanicsWithValue(t, "Non numeric value in asUint64()", func() {
	//	valueHelp.asUint64("123")
	//})
	if assert.NotPanics(t, func() {
		valueHelp.asUint64("123")
	}) {
		assert.Equal(t, valueHelp.asUint64("123"), (uint64)(0))
	}
}

// Test cases in private method asFloat64()
// including expected panic() cases
func TestHelpers_asFloat64(t *testing.T) {
	var r float64
	r = valueHelp.asFloat64(int(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(uint(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(int8(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(uint8(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(int16(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(uint16(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(int32(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(uint32(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(int64(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(uint64(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(float32(3))
	assert.Equal(t, float64(3), r)
	r = valueHelp.asFloat64(float64(3))
	assert.Equal(t, float64(3), r)

	// Non numeric
	//assert.PanicsWithValue(t, "Non numeric value in asFloat64()", func() {
	//	valueHelp.asFloat64("123")
	//})
	if assert.NotPanics(t, func() {
		valueHelp.asFloat64("123")
	}) {
		assert.Equal(t, valueHelp.asFloat64("123"), (float64)(0))
	}
}

/* Deprecated helper method:
func TestHelpers_ConvertToFloatEdgeCases(t *testing.T) {
	v := numberValidator{}
	// convert
	assert.Equal(t, float64(12.5), v.convertToFloat(float32(12.5)))
	assert.Equal(t, float64(12.5), v.convertToFloat(float64(12.5)))
	assert.Equal(t, float64(12), v.convertToFloat(int(12)))
	assert.Equal(t, float64(12), v.convertToFloat(int32(12)))
	assert.Equal(t, float64(12), v.convertToFloat(int64(12)))

	// does not convert
	assert.Equal(t, float64(0), v.convertToFloat("12"))
	// overflow : silent loss of info - ok (9.223372036854776e+18)
	assert.NotEqual(t, float64(0), v.convertToFloat(int64(math.MaxInt64)))
}
*/
