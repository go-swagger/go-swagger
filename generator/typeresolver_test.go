package generator

import (
	"fmt"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
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
	{"Comment", "", "models.Comment"},
	{"UserCard", "UserItem", "models.UserItem"},
}

func TestTypeResolver_AdditionalItems(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	tpe := spec.StringProperty()
	if assert.NoError(t, err) {
		// arrays of primitives and string formats with additional formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			var coll spec.Schema
			coll.Type = []string{"array"}
			coll.Items = new(spec.SchemaOrArray)
			coll.Items.Schema = tpe
			coll.AdditionalItems = new(spec.SchemaOrBool)
			coll.AdditionalItems.Schema = &sch

			rt, err := resolver.ResolveSchema(&coll, true)
			if assert.NoError(t, err) && assert.True(t, rt.IsArray) {
				assert.True(t, rt.HasAdditionalItems)
				assert.False(t, rt.IsNullable)
				//if assert.NotNil(t, rt.ElementType) {
				//assertPrimitiveResolve(t, "string", "", "string", *rt.ElementType)
				//}
			}
		}
	}
}

func TestTypeResolver_BasicTypes(t *testing.T) {

	_, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {

		// primitives and string formats
		for _, val := range schTypeVals {
			sch := new(spec.Schema)
			sch.Typed(val.Type, val.Format)

			rt, err := resolver.ResolveSchema(sch, true)
			if assert.NoError(t, err) {
				assert.False(t, rt.IsNullable)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}
		}

		// arrays of primitives and string formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
			}
		}

		// primitives and string formats
		for _, val := range schTypeVals {
			sch := new(spec.Schema)
			sch.Typed(val.Type, val.Format)
			sch.Extensions = make(spec.Extensions)
			sch.Extensions["x-isnullable"] = true

			rt, err := resolver.ResolveSchema(sch, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}
		}

		// arrays of primitives and string formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			sch.AddExtension("x-isnullable", true)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
			}
		}

	}

}

func TestTypeResolver_Refs(t *testing.T) {

	_, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {

		// referenced objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(sch, true)
			if assert.NoError(t, err) {
				assert.Equal(t, val.Expected, rt.GoType)
				assert.False(t, rt.IsAnonymous)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// referenced array objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
			}
		}
		// for named objects
		// referenced objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(sch, false)
			if assert.NoError(t, err) {
				assert.Equal(t, val.Expected, rt.GoType)
				assert.False(t, rt.IsAnonymous)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// referenced array objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), false)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
			}
		}
	}
}

func TestTypeResolver_AdditionalProperties(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {

		// primitives as additional properties
		for _, val := range schTypeVals {
			sch := new(spec.Schema)

			sch.Typed(val.Type, val.Format)
			parent := new(spec.Schema)
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = sch

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// array of primitives as additional properties
		for _, val := range schTypeVals {
			sch := new(spec.Schema)

			sch.Typed(val.Type, val.Format)
			parent := new(spec.Schema)
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = new(spec.Schema).CollectionOf(*sch)

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "map[string][]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// refs as additional properties
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)
			parent := new(spec.Schema)
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = sch

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// when additional properties and properties present, it's a complex object

		// primitives as additional properties
		for _, val := range schTypeVals {
			sch := new(spec.Schema)

			sch.Typed(val.Type, val.Format)
			parent := new(spec.Schema)
			parent.Properties = make(map[string]spec.Schema)
			parent.Properties["id"] = *spec.Int32Property()
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = sch

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// array of primitives as additional properties
		for _, val := range schTypeVals {
			sch := new(spec.Schema)

			sch.Typed(val.Type, val.Format)
			parent := new(spec.Schema)
			parent.Properties = make(map[string]spec.Schema)
			parent.Properties["id"] = *spec.Int32Property()
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = new(spec.Schema).CollectionOf(*sch)

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "map[string][]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// refs as additional properties
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)
			parent := new(spec.Schema)
			parent.Properties = make(map[string]spec.Schema)
			parent.Properties["id"] = *spec.Int32Property()
			parent.AdditionalProperties = new(spec.SchemaOrBool)
			parent.AdditionalProperties.Schema = sch

			rt, err := resolver.ResolveSchema(parent, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

	}
}

func basicTaskListResolver(t testing.TB) (*spec.Document, *typeResolver, error) {
	tlb, err := spec.Load("../fixtures/codegen/tasklist.basic.yml")
	if err != nil {
		return nil, nil, err
	}
	swsp := tlb.Spec()
	uc := swsp.Definitions["UserCard"]
	uc.AddExtension("x-go-name", "UserItem")
	swsp.Definitions["UserCard"] = uc
	return tlb, &typeResolver{
		Doc:           tlb,
		ModelsPackage: "models",
	}, nil
}

func TestTypeResolver_TupleTypes(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {
		// tuple type (items with multiple schemas)
		parent := new(spec.Schema)
		parent.Typed("array", "")
		parent.Items = new(spec.SchemaOrArray)
		parent.Items.Schemas = append(
			parent.Items.Schemas,
			*spec.StringProperty(),
			*spec.Int64Property(),
			*spec.Float64Property(),
			*spec.BoolProperty(),
			*spec.ArrayProperty(spec.StringProperty()),
			*spec.RefProperty("#/definitions/Comment"),
		)

		rt, err := resolver.ResolveSchema(parent, true)
		if assert.NoError(t, err) {
			assert.False(t, rt.IsArray)
			assert.True(t, rt.IsTuple)
		}
	}
}
func TestTypeResolver_AnonymousStructs(t *testing.T) {

	_, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {
		// anonymous structs should be accounted for
		parent := new(spec.Schema)
		parent.Typed("object", "")
		parent.Properties = make(map[string]spec.Schema)
		parent.Properties["name"] = *spec.StringProperty()
		parent.Properties["age"] = *spec.Int32Property()

		rt, err := resolver.ResolveSchema(parent, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsNullable)
			assert.True(t, rt.IsAnonymous)
			assert.True(t, rt.IsComplexObject)
		}

		parent.Extensions = make(spec.Extensions)
		parent.Extensions["x-isnullable"] = true

		rt, err = resolver.ResolveSchema(parent, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsNullable)
			assert.True(t, rt.IsAnonymous)
			assert.True(t, rt.IsComplexObject)
		}
	}
}
func TestTypeResolver_ObjectType(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	resolver.ModelName = "TheModel"
	defer func() { resolver.ModelName = "" }()

	if assert.NoError(t, err) {
		//very poor schema definitions (as in none)
		types := []string{"object", ""}
		for _, tpe := range types {
			sch := new(spec.Schema)
			sch.Typed(tpe, "")
			rt, err := resolver.ResolveSchema(sch, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "map[string]interface{}", rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}

			sch.Properties = make(map[string]spec.Schema)
			var ss spec.Schema
			sch.Properties["tags"] = *(&ss).CollectionOf(*spec.StringProperty())
			rt, err = resolver.ResolveSchema(sch, false)
			assert.True(t, rt.IsComplexObject)
			assert.False(t, rt.IsMap)
			assert.Equal(t, "models.TheModel", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)

			sch.Properties = nil
			nsch := new(spec.Schema)
			nsch.Typed(tpe, "")
			nsch.AllOf = []spec.Schema{*sch}
			rt, err = resolver.ResolveSchema(nsch, false)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "models.TheModel", rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}
		sch := new(spec.Schema)
		rt, err := resolver.ResolveSchema(sch, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsMap)
			assert.False(t, rt.IsComplexObject)
			assert.Equal(t, "map[string]interface{}", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)

		}
		sch = new(spec.Schema)
		var sp spec.Schema
		sp.Typed("object", "")
		sch.AllOf = []spec.Schema{sp}
		rt, err = resolver.ResolveSchema(sch, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsComplexObject)
			assert.False(t, rt.IsMap)
			assert.Equal(t, "models.TheModel", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)
		}
	}
}

func assertPrimitiveResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType) {
	assert.Equal(t, tpe, tr.SwaggerType, fmt.Sprintf("expected %q (%q, %q) to for the swagger type but got %q", tpe, tfmt, exp, tr.SwaggerType))
	assert.Equal(t, tfmt, tr.SwaggerFormat, fmt.Sprintf("expected %q (%q, %q) to for the swagger format but got %q", tfmt, tpe, exp, tr.SwaggerFormat))
	assert.Equal(t, exp, tr.GoType, fmt.Sprintf("expected %q (%q, %q) to for the go type but got %q", exp, tpe, tfmt, tr.GoType))
}
