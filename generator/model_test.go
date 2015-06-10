package generator

import (
	"fmt"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

var schTypeVals = []struct{ Type, Format, Expected string }{
	{"boolean", "", "bool"},
	{"string", "", "string"},
	{"integer", "int8", "int8"},
	{"integer", "int16", "int16"},
	{"integer", "int32", "int32"},
	{"integer", "int64", "int64"},
	{"integer", "", "int64"},
	{"integer", "uint8", "uint8"},
	{"integer", "uint16", "uint16"},
	{"integer", "uint32", "uint32"},
	{"integer", "uint64", "uint64"},
	{"number", "float", "float32"},
	{"number", "double", "float64"},
	{"number", "", "float64"},
	{"string", "byte", "strfmt.Base64"},
	{"string", "date", "strfmt.Date"},
	{"string", "date-time", "strfmt.DateTime"},
	{"string", "uri", "strfmt.URI"},
	{"string", "email", "strfmt.Email"},
	{"string", "hostname", "strfmt.Hostname"},
	{"string", "ipv4", "strfmt.IPv4"},
	{"string", "ipv6", "strfmt.IPv6"},
	{"string", "uuid", "strfmt.UUID"},
	{"string", "uuid3", "strfmt.UUID3"},
	{"string", "uuid4", "strfmt.UUID4"},
	{"string", "uuid5", "strfmt.UUID5"},
	{"string", "isbn", "strfmt.ISBN"},
	{"string", "isbn10", "strfmt.ISBN10"},
	{"string", "isbn13", "strfmt.ISBN13"},
	{"string", "creditcard", "strfmt.CreditCard"},
	{"string", "ssn", "strfmt.SSN"},
	{"string", "hexcolor", "strfmt.HexColor"},
	{"string", "rgbcolor", "strfmt.RGBColor"},
	{"string", "duration", "strfmt.Duration"},
	{"string", "password", "strfmt.Password"},
	{"file", "", "httpkit.File"},
}

var schRefVals = []struct{ Type, GoType, Expected string }{
	{"Pet", "", "models.Pet"},
	{"pet", "", "models.Pet"},
}

func TestTypeForSchema_Primitives(t *testing.T) {
	for _, val := range schTypeVals {
		sch := new(spec.Schema)
		sch.Typed(val.Type, val.Format)
		actual := typeForSchema(sch, "models")
		assert.Equal(t, val.Expected, actual)
	}
}

func TestTypeForSchema_Objects(t *testing.T) {
	sch := new(spec.Schema)
	sch.Ref, _ = spec.NewRef("#/definitions/Pet")
	actual := typeForSchema(sch, "models")
	assert.Equal(t, "models.Pet", actual)

}

func TestTypeForSchema_Arrays(t *testing.T) {
	for _, val := range schTypeVals {
		sch := new(spec.Schema)
		sch.Typed(val.Type, val.Format)
		actual := typeForSchema(new(spec.Schema).CollectionOf(*sch), "models")
		assert.Equal(t, "[]"+val.Expected, actual)
	}
}

func TestTypeResolver(t *testing.T) {
	tlb, err := spec.Load("../fixtures/codegen/tasklist.basic.yml")
	if assert.NoError(t, err) {
		resolver := &typeResolver{
			Doc:           tlb,
			ModelsPackage: "models",
		}

		// primitives and string formats
		for _, val := range schTypeVals {
			sch := new(spec.Schema)
			sch.Typed(val.Type, val.Format)

			rt, err := resolver.ResolveSchema(sch)
			if assert.NoError(t, err) {
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}
		}

		// arrays of primitives and string formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch))
			if assert.NoError(t, err) && assert.True(t, rt.IsArray) && assert.NotNil(t, rt.ElementType) {
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, *rt.ElementType)
			}
		}

		// referenced objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(sch)
			if assert.NoError(t, err) {
				assert.Equal(t, val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// referenced array objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch))
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
				assert.Equal(t, val.Expected, rt.ElementType.GoType)
				assert.Equal(t, "object", rt.ElementType.SwaggerType)
			}
		}

		//for _, val := range schTypeVals {
		//}

		// object definitions with additional properties
		sch := new(spec.Schema)
		sch.Typed("object", "")
		rt, err := resolver.ResolveSchema(sch)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsMap)
			assert.Equal(t, "map[string]interface{}", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)

			if assert.NotNil(t, rt.ElementType) {
				assert.True(t, rt.ElementType.IsInterface)
				assert.Equal(t, "interface{}", rt.ElementType.GoType)
			}
		}
	}
}

func assertPrimitiveResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType) {
	assert.Equal(t, tpe, tr.SwaggerType, fmt.Sprintf("expected %q (%q, %q) to for the swagger type but got %q", tpe, tfmt, exp, tr.SwaggerType))
	assert.Equal(t, tfmt, tr.SwaggerFormat, fmt.Sprintf("expected %q (%q, %q) to for the swagger format but got %q", tfmt, tpe, exp, tr.SwaggerFormat))
	assert.Equal(t, exp, tr.GoType, fmt.Sprintf("expected %q (%q, %q) to for the go type but got %q", exp, tpe, tfmt, tr.GoType))
}
