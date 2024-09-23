package generator

import (
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeResolver_NestedAliasedSlice(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Statix"
	schema := definitions[k]

	tr := newTypeResolver("models", "", specDoc)
	specDoc.Spec().Definitions["StatixItems0"] = *schema.Items.Schema.Items.Schema.Items.Schema
	schema.Items.Schema.Items.Schema.Items.Schema = spec.RefProperty("#/definitions/StatixItems0")
	tr.KnownDefs["StatixItems0"] = struct{}{}
	tr.ModelName = k
	rt, err := tr.ResolveSchema(&schema, false, false)
	require.NoError(t, err)

	assert.Equal(t, "[][][]*models.StatixItems0", rt.AliasedType)
}

func TestTypeResolver_PointerLifting(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	require.NoError(t, err)

	testPointToPrimitives(t, *resolver, false /* not aliased */)
	testPointToPrimitives(t, *resolver, true /* aliased */)
	testPointToSliceElements(t, *resolver, false /* not aliased */)
	testPointToSliceElements(t, *resolver, true /* aliased */)
	testPointToAdditionalPropertiesElements(t, *resolver, false /* not aliased */)
	testPointToAdditionalPropertiesElements(t, *resolver, true /* aliased */)
}

type builtinVal struct {
	Type, Format, Expected, AliasedType string

	Nullable, Aliased bool

	Default interface{}

	Extensions spec.Extensions

	Required         bool
	ReadOnly         bool
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
	MaxLength        *int64
	MinLength        *int64
	Pattern          string
	MaxItems         *int64
	MinItems         *int64
	UniqueItems      bool
	MultipleOf       *float64
	Enum             []interface{}
}

func nullableExt() spec.Extensions {
	return spec.Extensions{"x-nullable": true}
}

func isNullableExt() spec.Extensions {
	return spec.Extensions{"x-isnullable": true}
}

func notNullableExt() spec.Extensions {
	return spec.Extensions{"x-nullable": false}
}

func isNotNullableExt() spec.Extensions {
	return spec.Extensions{"x-isnullable": false}
}

var boolPointerVals = []builtinVal{
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: nullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: nullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: false, ReadOnly: false, Extensions: nullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: isNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: isNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: false, ReadOnly: false, Extensions: isNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: false, ReadOnly: false, Extensions: notNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: false, Extensions: notNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false, Extensions: notNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true, Extensions: notNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true, Extensions: notNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: false, ReadOnly: false, Extensions: isNotNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: false, Extensions: isNotNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false, Extensions: isNotNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true, Extensions: isNotNullableExt()},
	{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true, Extensions: isNotNullableExt()},
}

func generateNumberPointerVals(t, v string) (result []builtinVal) {
	vv := v
	if vv == "" || vv == "int" {
		if t == "integer" {
			vv = "int64"
		} else {
			vv = "float64"
		}
	}
	if vv == "uint" && t == "integer" {
		vv = "uint64"
	}
	if t == "number" {
		if v == "float" {
			vv = "float32"
		} else {
			vv = "float64"
		}
	}
	return []builtinVal{
		// plain vanilla
		{Type: t, Format: v, Expected: vv},
		{Type: t, Format: v, Expected: vv, Nullable: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Extensions: isNullableExt()}, // 2

		// plain vanilla readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 9

		// required
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Extensions: isNullableExt()}, // 12

		// required, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Default: 3, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 19

		// minimum validation
		{Type: t, Format: v, Expected: vv, Minimum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 23

		// minimum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Minimum: swag.Float64(2), Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true, Minimum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 31

		// required, minimum validation
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 35

		// required, minimum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Minimum: swag.Float64(2), ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), ReadOnly: true, Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 42

		// maximum validation
		{Type: t, Format: v, Expected: vv, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 46

		// maximum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 54

		// required, maximum validation
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 58

		// required, maximum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Maximum: swag.Float64(2), ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), ReadOnly: true, Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 65

		// minimum and maximum validation
		{Type: t, Format: v, Expected: vv, Minimum: swag.Float64(2), Maximum: swag.Float64(5)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(0), Maximum: swag.Float64(1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: isNullableExt()}, // 72

		// minimum and maximum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(0), Maximum: swag.Float64(3)},
		{Type: t, Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(-1), ReadOnly: true, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, Minimum: swag.Float64(-1), Maximum: swag.Float64(6)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(1), Maximum: swag.Float64(6)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(-6), Maximum: swag.Float64(-1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 83

		// required, minimum and maximum validation
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(5)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(0), Maximum: swag.Float64(1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: isNullableExt()}, // 89

		// required, minimum and maximum validation, readonly and defaults
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(0), Maximum: swag.Float64(3)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		{Type: t, Format: v, Expected: vv, Nullable: false, Required: true, Default: 3, Minimum: swag.Float64(-1), ReadOnly: true, Maximum: swag.Float64(2)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(-1), Maximum: swag.Float64(6)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(1), Maximum: swag.Float64(6)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(-6), Maximum: swag.Float64(-1)},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: nullableExt()},
		{Type: t, Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 99
	}
}

var stringPointerVals = []builtinVal{
	{Type: "string", Format: "", Expected: "string", Nullable: false},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Extensions: isNullableExt()}, // 2

	// plain vanilla readonly and defaults
	{Type: "string", Format: "", Expected: "string", Nullable: false, ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Default: 3},
	{Type: "string", Format: "", Expected: "string", Nullable: false, Default: 3, ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 9

	// required
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, Extensions: isNullableExt()}, // 12

	// required, readonly and defaults
	{Type: "string", Format: "", Expected: "string", Nullable: false, Required: true, ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, Default: 3},
	{Type: "string", Format: "", Expected: "string", Nullable: false, Required: true, Default: 3, ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 19

	// minLength validation
	{Type: "string", Format: "", Expected: "string", Nullable: false, MinLength: swag.Int64(2)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, MinLength: swag.Int64(0)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, MinLength: swag.Int64(2), Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, MinLength: swag.Int64(2), Extensions: isNullableExt()}, // 23

	// minLength validation, readonly and defaults
	{Type: "string", Format: "", Expected: "string", Nullable: false, ReadOnly: true, MinLength: swag.Int64(2)},
	{Type: "string", Format: "", Expected: "string", Nullable: false, ReadOnly: true, MinLength: swag.Int64(0)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, ReadOnly: true, MinLength: swag.Int64(2), Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, ReadOnly: true, MinLength: swag.Int64(2), Extensions: isNullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: false, Default: 3, MinLength: swag.Int64(2)},
	{Type: "string", Format: "", Expected: "string", Nullable: false, Default: 3, ReadOnly: true, MinLength: swag.Int64(2)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Default: 3, ReadOnly: true, MinLength: swag.Int64(2), Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Default: 3, ReadOnly: true, MinLength: swag.Int64(2), Extensions: isNullableExt()}, // 31

	// required, minLength validation
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(0)},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), Extensions: isNullableExt()}, // 35

	// required, minLength validation, readonly and defaults
	{Type: "string", Format: "", Expected: "string", Nullable: false, Required: true, MinLength: swag.Int64(2), ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), Default: 3},
	{Type: "string", Format: "", Expected: "string", Nullable: false, Required: true, MinLength: swag.Int64(2), Default: 3, ReadOnly: true},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "", Expected: "string", Nullable: true, Required: true, MinLength: swag.Int64(2), Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 42
}

var strfmtValues = []builtinVal{
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: false},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Extensions: isNullableExt()}, // 2

	// plain vanilla readonly and defaults
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: false, ReadOnly: true},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Default: 3},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: false, Default: 3, ReadOnly: true},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 9

	// required
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, Extensions: isNullableExt()}, // 12

	// required, readonly and defaults
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: false, Required: true, ReadOnly: true},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, Default: 3},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: false, Required: true, Default: 3, ReadOnly: true},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "password", Expected: "strfmt.Password", Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 19

	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Extensions: isNullableExt()}, // 22

	// plain vanilla readonly and defaults
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, ReadOnly: true},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Default: 3},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Default: 3, ReadOnly: true},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 29

	// required
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Extensions: isNullableExt()}, // 32

	// required, readonly and defaults
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, ReadOnly: true},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Default: 3},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Default: 3, ReadOnly: true},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
	{Type: "string", Format: "binary", Expected: "io.ReadCloser", Nullable: false, Required: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 39
}

func testPointToAdditionalPropertiesElements(t testing.TB, tr typeResolver, aliased bool) bool {
	if aliased {
		tr.ModelName = "MyAliasedThing"
	}
	resolver := &tr
	for i, val := range boolPointerVals {
		if !assertBuiltinAdditionalPropertiesElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	for _, v := range []string{"", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
		passed := true
		for i, val := range generateNumberPointerVals("integer", v) {
			if !assertBuiltinAdditionalPropertiesElem(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}
	for _, v := range []string{"", "float", "double"} {
		passed := true
		for i, val := range generateNumberPointerVals("number", v) {
			if !assertBuiltinAdditionalPropertiesElem(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}
	for i, val := range stringPointerVals {
		if !assertBuiltinAdditionalPropertiesElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	for i, val := range strfmtValues {
		if !assertBuiltinAdditionalPropertiesElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	return true
}

func testPointToSliceElements(t testing.TB, tr typeResolver, aliased bool) bool {
	if aliased {
		tr.ModelName = "MyAliasedThing"
	}
	resolver := &tr
	for i, val := range boolPointerVals {
		if !assertBuiltinSliceElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	for _, v := range []string{"", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
		passed := true
		for i, val := range generateNumberPointerVals("integer", v) {
			if !assertBuiltinSliceElem(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}
	for _, v := range []string{"", "float", "double"} {
		passed := true
		for i, val := range generateNumberPointerVals("number", v) {
			if !assertBuiltinSliceElem(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}
	for i, val := range stringPointerVals {
		if !assertBuiltinSliceElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	for i, val := range strfmtValues {
		if !assertBuiltinSliceElem(t, resolver, aliased, i, val) {
			return false
		}
	}
	return true
}

func testPointToPrimitives(t testing.TB, tr typeResolver, aliased bool) bool {
	if aliased {
		tr.ModelName = "MyAliasedThing"
		tr.KnownDefs[tr.ModelName] = struct{}{}
	}
	resolver := &tr
	for i, val := range boolPointerVals {
		if !assertBuiltinVal(t, resolver, aliased, i, val) {
			return false
		}
	}
	for _, v := range []string{"", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
		passed := true
		for i, val := range generateNumberPointerVals("integer", v) {
			if !assertBuiltinVal(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}
	for _, v := range []string{"", "float", "double"} {
		passed := true
		for i, val := range generateNumberPointerVals("number", v) {
			if !assertBuiltinVal(t, resolver, aliased, i, val) {
				passed = false
			}
		}
		if !passed {
			return false
		}
	}

	for i, val := range stringPointerVals {
		if !assertBuiltinVal(t, resolver, aliased, i, val) {
			return false
		}
	}

	for i, val := range strfmtValues {
		if !assertBuiltinVal(t, resolver, aliased, i, val) {
			return false
		}
	}
	return true
}

func assertBuiltinVal(t testing.TB, resolver *typeResolver, aliased bool, i int, val builtinVal) bool {
	val.Aliased = aliased
	if aliased {
		val.AliasedType = val.Expected
		val.Expected = "models.MyAliasedThing"
	}

	sch := new(spec.Schema)
	sch.Typed(val.Type, val.Format)
	sch.Default = val.Default
	sch.ReadOnly = val.ReadOnly
	sch.Extensions = val.Extensions
	sch.Minimum = val.Minimum
	sch.Maximum = val.Maximum
	sch.MultipleOf = val.MultipleOf
	sch.MinLength = val.MinLength
	sch.MaxLength = val.MaxLength

	rt, err := resolver.ResolveSchema(sch, !aliased, val.Required)
	require.NoError(t, err)
	if val.Nullable {
		if !assert.True(t, rt.IsNullable, "expected nullable for item at: %d", i) {
			// fmt.Println("isRequired:", val.Required)
			// pretty.Println(sch)
			return false
		}
	} else {
		if !assert.False(t, rt.IsNullable, "expected not nullable for item at: %d", i) {
			// fmt.Println("isRequired:", val.Required)
			// pretty.Println(sch)
			return false
		}
	}
	if !assert.Equal(t, val.Aliased, rt.IsAliased, "expected (%q, %q) to be an aliased type", val.Type, val.Format) {
		return false
	}
	if val.Aliased {
		if !assert.Equal(t, val.AliasedType, rt.AliasedType, "expected %q (%q, %q) to be aliased as %q, but got %q", val.Expected, val.Type, val.Format, val.AliasedType, rt.AliasedType) {
			return false
		}
	}
	if !assertBuiltinResolve(t, val.Type, val.Format, val.Expected, rt, i) {
		return false
	}
	return true
}

func assertBuiltinSliceElem(t testing.TB, resolver *typeResolver, aliased bool, i int, val builtinVal) bool {
	val.Nullable = false
	if nullableExtension(val.Extensions) != nil {
		val.Nullable = *nullableExtension(val.Extensions)
	}
	sliceType := "[]" + val.Expected
	if val.Nullable {
		sliceType = "[]*" + val.Expected
	}
	val.Expected = sliceType

	val.Aliased = aliased
	if aliased {
		val.AliasedType = val.Expected
		val.Expected = "models.MyAliasedThing"
	}

	items := new(spec.Schema)
	items.Typed(val.Type, val.Format)
	items.Default = val.Default
	items.ReadOnly = val.ReadOnly
	items.Extensions = val.Extensions
	items.Minimum = val.Minimum
	items.Maximum = val.Maximum
	items.MultipleOf = val.MultipleOf
	items.MinLength = val.MinLength
	items.MaxLength = val.MaxLength

	sch := spec.ArrayProperty(items)

	rt, err := resolver.ResolveSchema(sch, !aliased, val.Required)
	require.NoError(t, err)

	if val.Nullable {
		if !assert.True(t, rt.ElemType.IsNullable, "expected nullable for item at: %d", i) {
			return false
		}
	} else {
		if !assert.False(t, rt.ElemType != nil && rt.ElemType.IsNullable, "expected not nullable for item at: %d", i) {
			return false
		}
	}

	if val.Aliased {
		if !assert.Equal(t, val.Aliased, rt.IsAliased, "expected (%q, %q) to be an aliased type at: %d", val.Type, val.Format, i) {
			return false
		}
		if !assert.Equal(t, val.AliasedType, rt.AliasedType, "expected %q (%q, %q) to be aliased as %q, but got %q at %d", val.Expected, val.Type, val.Format, val.AliasedType, rt.AliasedType, i) {
			return false
		}
	}

	if !assertBuiltinSliceElemnResolve(t, val.Type, val.Format, val.Expected, rt, i) {
		return false
	}
	return true
}

func assertBuiltinAdditionalPropertiesElem(t testing.TB, resolver *typeResolver, aliased bool, i int, val builtinVal) bool {
	val.Nullable = false
	if nullableExtension(val.Extensions) != nil {
		val.Nullable = *nullableExtension(val.Extensions)
	}
	sliceType := "map[string]" + val.Expected
	if val.Nullable {
		sliceType = "map[string]*" + val.Expected
	}
	val.Expected = sliceType

	val.Aliased = aliased
	if aliased {
		val.AliasedType = val.Expected
		val.Expected = "models.MyAliasedThing"
	}

	items := new(spec.Schema)
	items.Typed(val.Type, val.Format)
	items.Default = val.Default
	items.ReadOnly = val.ReadOnly
	items.Extensions = val.Extensions
	items.Minimum = val.Minimum
	items.Maximum = val.Maximum
	items.MultipleOf = val.MultipleOf
	items.MinLength = val.MinLength
	items.MaxLength = val.MaxLength

	sch := spec.MapProperty(items)

	rt, err := resolver.ResolveSchema(sch, !aliased, val.Required)
	require.NoError(t, err)

	if val.Nullable {
		if !assert.True(t, rt.ElemType.IsNullable, "expected nullable for item at: %d", i) {
			return false
		}
	} else {
		if !assert.False(t, rt.ElemType != nil && rt.ElemType.IsNullable, "expected not nullable for item at: %d", i) {
			return false
		}
	}

	if !assert.Equal(t, val.Aliased, rt.IsAliased, "expected (%q, %q) to be an aliased type at %d", val.Type, val.Format, i) {
		return false
	}

	if val.Aliased {
		if !assert.Equal(t, val.AliasedType, rt.AliasedType, "expected %q (%q, %q) to be aliased as %q, but got %q at %d", val.Expected, val.Type, val.Format, val.AliasedType, rt.AliasedType, i) {
			return false
		}
	}

	if !assertBuiltinSliceElemnResolve(t, val.Type, val.Format, val.Expected, rt, i) {
		return false
	}
	return true
}

func assertBuiltinResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType, i int) bool {
	return assert.Equal(t, tpe, tr.SwaggerType, "expected %q (%q, %q) at %d for the swagger type but got %q", tpe, tfmt, exp, i, tr.SwaggerType) &&
		assert.Equal(t, tfmt, tr.SwaggerFormat, "expected %q (%q, %q) at %d for the swagger format but got %q", tfmt, tpe, exp, i, tr.SwaggerFormat) &&
		assert.Equal(t, exp, tr.GoType, "expected %q (%q, %q) at %d for the go type but got %q", exp, tpe, tfmt, i, tr.GoType)
}

func assertBuiltinSliceElemnResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType, i int) bool {
	return assert.Equal(t, tpe, tr.ElemType.SwaggerType, "expected %q (%q, %q) at %d for the swagger type but got %q", tpe, tfmt, exp, i, tr.SwaggerType) &&
		assert.Equal(t, tfmt, tr.ElemType.SwaggerFormat, "expected %q (%q, %q) at %d for the swagger format but got %q", tfmt, tpe, exp, i, tr.SwaggerFormat) &&
		assert.Equal(t, exp, tr.GoType, "expected %q (%q, %q) at %d for the go type but got %q", exp, tpe, tfmt, i, tr.GoType)
}
