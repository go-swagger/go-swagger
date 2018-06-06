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
	"math"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestNumberValidator_EdgeCases(t *testing.T) {
	// Apply
	var min float64 = float64(math.MinInt32 - 1)
	var max float64 = float64(math.MaxInt32 + 1)

	v := numberValidator{
		Path: "path",
		In:   "in",
		//Default:
		//MultipleOf:
		Maximum:          &max, // *float64
		ExclusiveMaximum: false,
		Minimum:          &min, // *float64
		ExclusiveMinimum: false,
		// Allows for more accurate behavior regarding integers
		Type:   "integer",
		Format: "int32",
	}

	// numberValidator applies to: Parameter,Schema,Items,Header

	sources := []interface{}{
		new(spec.Parameter),
		new(spec.Schema),
		new(spec.Items),
		new(spec.Header),
	}

	testNumberApply(t, &v, sources)

	assert.False(t, v.Applies(float64(32), reflect.Float64))

	// Now for different scenarios on Minimum, Maximum
	// - The Maximum value does not respect the Type|Format specification
	// - Value is checked as float64 with Maximum as float64 and fails
	res := new(Result)

	res = v.Validate(int64(math.MaxInt32 + 2))
	assert.True(t, res.HasErrors())
	// - The Minimum value does not respect the Type|Format specification
	// - Value is checked as float64 with Maximum as float64 and fails
	res = v.Validate(int64(math.MinInt32 - 2))
	assert.True(t, res.HasErrors())
}

func testNumberApply(t *testing.T, v *numberValidator, sources []interface{}) {
	for _, source := range sources {
		// numberValidator does not applies to:
		assert.False(t, v.Applies(source, reflect.String))
		assert.False(t, v.Applies(source, reflect.Struct))
		// numberValidator applies to:
		assert.True(t, v.Applies(source, reflect.Int))
		assert.True(t, v.Applies(source, reflect.Int8))
		assert.True(t, v.Applies(source, reflect.Uint16))
		assert.True(t, v.Applies(source, reflect.Uint32))
		assert.True(t, v.Applies(source, reflect.Uint64))
		assert.True(t, v.Applies(source, reflect.Uint))
		assert.True(t, v.Applies(source, reflect.Uint8))
		assert.True(t, v.Applies(source, reflect.Uint16))
		assert.True(t, v.Applies(source, reflect.Uint32))
		assert.True(t, v.Applies(source, reflect.Uint64))
		assert.True(t, v.Applies(source, reflect.Float32))
		assert.True(t, v.Applies(source, reflect.Float64))
	}
}

func TestStringValidator_EdgeCases(t *testing.T) {
	// Apply

	v := stringValidator{}

	// stringValidator applies to: Parameter,Schema,Items,Header

	sources := []interface{}{
		new(spec.Parameter),
		new(spec.Schema),
		new(spec.Items),
		new(spec.Header),
	}

	testStringApply(t, &v, sources)

	assert.False(t, v.Applies("A string", reflect.String))

}

func testStringApply(t *testing.T, v *stringValidator, sources []interface{}) {
	for _, source := range sources {
		// numberValidator does not applies to:
		assert.False(t, v.Applies(source, reflect.Struct))
		assert.False(t, v.Applies(source, reflect.Int))
		// numberValidator applies to:
		assert.True(t, v.Applies(source, reflect.String))
	}
}

func TestBasicCommonValidator_EdgeCases(t *testing.T) {
	// Apply

	v := basicCommonValidator{}

	// basicCommonValidator applies to: Parameter,Schema,Header

	sources := []interface{}{
		new(spec.Parameter),
		new(spec.Schema),
		new(spec.Header),
	}

	testCommonApply(t, &v, sources)

	assert.False(t, v.Applies("A string", reflect.String))

}

func testCommonApply(t *testing.T, v *basicCommonValidator, sources []interface{}) {
	for _, source := range sources {
		assert.True(t, v.Applies(source, reflect.String))
	}
}

func TestBasicSliceValidator_EdgeCases(t *testing.T) {
	// Apply

	v := basicSliceValidator{}

	// basicCommonValidator applies to: Parameter,Schema,Header

	sources := []interface{}{
		new(spec.Parameter),
		new(spec.Items),
		new(spec.Header),
	}

	testSliceApply(t, &v, sources)

	assert.False(t, v.Applies(new(spec.Schema), reflect.Slice))
	assert.False(t, v.Applies(new(spec.Parameter), reflect.String))

}

func testSliceApply(t *testing.T, v *basicSliceValidator, sources []interface{}) {
	for _, source := range sources {
		assert.True(t, v.Applies(source, reflect.Slice))
	}
}

type anything struct {
	anyProperty int
}

// hasDuplicates() is currently not exercised by common spec testcases
// (this method is not used by the validator atm)
// Here is a unit exerciser
// NOTE: this method is probably obsolete and superseeded by values.go:UniqueItems()
// which is superior in every respect to this one.
func TestBasicSliceValidator_HasDuplicates(t *testing.T) {
	s := basicSliceValidator{}
	// hasDuplicates() makes no hypothesis about the underlying object,
	// save being an array, slice or string (same constraint as reflect.Value.Index())
	// it also comes without safeguard or anything.
	vi := []int{1, 2, 3}
	vs := []string{"a", "b", "c"}
	vt := []anything{
		anything{anyProperty: 1},
		anything{anyProperty: 2},
		anything{anyProperty: 3},
	}
	assert.False(t, s.hasDuplicates(reflect.ValueOf(vi), len(vi)))
	// how UniqueItems() is superior? Look:   err := uniqueItems("path","body", vi)
	assert.False(t, s.hasDuplicates(reflect.ValueOf(vs), len(vs)))
	assert.False(t, s.hasDuplicates(reflect.ValueOf(vt), len(vt)))

	di := []int{1, 1, 3}
	ds := []string{"a", "b", "a"}
	dt := []anything{
		anything{anyProperty: 1},
		anything{anyProperty: 2},
		anything{anyProperty: 2},
	}
	assert.True(t, s.hasDuplicates(reflect.ValueOf(di), len(di)))
	assert.True(t, s.hasDuplicates(reflect.ValueOf(ds), len(ds)))
	assert.True(t, s.hasDuplicates(reflect.ValueOf(dt), len(dt)))
}
