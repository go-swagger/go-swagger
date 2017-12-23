package validate

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringValidator_Validate_Panic(t *testing.T) {
	var schemaJSON = `
{
    "properties": {
        "name": {
            "type": "string",
            "pattern": "^[A-Za-z]+$",
            "minLength": 1
        },
        "place": {
            "type": "string",
            "pattern": "^[A-Za-z]+$",
            "minLength": 1
        }
    },
    "required": [
        "name"
    ]
}`
	var inputJSON = `{"name": "Ivan"}`
	schema := new(spec.Schema)
	require.NoError(t, json.Unmarshal([]byte(schemaJSON), schema))
	var input map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(inputJSON), &input))
	input["place"] = json.Number("10")

	assert.Error(t, AgainstSchema(schema, input, strfmt.Default))
}

func TestNumberValidator_ConvertToFloatEdgeCases(t *testing.T) {
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
