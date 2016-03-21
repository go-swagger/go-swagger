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

package generator

import (
	"fmt"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
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

			rt, err := resolver.ResolveSchema(&coll, true, true)
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

			rt, err := resolver.ResolveSchema(sch, true, false)
			if assert.NoError(t, err) {
				assert.False(t, rt.IsNullable, "expected %s with format %q to not be nullable", val.Type, val.Format)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}
		}

		// arrays of primitives and string formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
			}
		}

		// primitives and string formats
		for _, val := range schTypeVals {
			sch := new(spec.Schema)
			sch.Typed(val.Type, val.Format)
			sch.Extensions = make(spec.Extensions)
			sch.Extensions[xIsNullable] = true

			rt, err := resolver.ResolveSchema(sch, true, false)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}

			// Test x-nullable overrides x-isnullable
			sch.Extensions[xIsNullable] = false
			sch.Extensions[xNullable] = true
			rt, err = resolver.ResolveSchema(sch, true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}

			// Test x-nullable without x-isnullable
			delete(sch.Extensions, xIsNullable)
			sch.Extensions[xNullable] = true
			rt, err = resolver.ResolveSchema(sch, true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
				assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
			}
		}

		// arrays of primitives and string formats
		for _, val := range schTypeVals {
			var sch spec.Schema
			sch.Typed(val.Type, val.Format)
			sch.AddExtension(xIsNullable, true)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true, true)
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

			rt, err := resolver.ResolveSchema(sch, true, true)
			if assert.NoError(t, err) {
				assert.Equal(t, val.Expected, rt.GoType)
				assert.False(t, rt.IsAnonymous)
				assert.True(t, rt.IsNullable)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// referenced array objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
				assert.Equal(t, "[]*"+val.Expected, rt.GoType)
			}
		}
		// for named objects
		// referenced objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(sch, false, true)
			if assert.NoError(t, err) {
				assert.Equal(t, val.Expected, rt.GoType)
				assert.False(t, rt.IsAnonymous)
				assert.True(t, rt.IsNullable)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

		// referenced array objects
		for _, val := range schRefVals {
			sch := new(spec.Schema)
			sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

			rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), false, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsArray)
				assert.Equal(t, "[]*"+val.Expected, rt.GoType)
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

			rt, err := resolver.ResolveSchema(parent, true, false)
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

			rt, err := resolver.ResolveSchema(parent, true, false)
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

			rt, err := resolver.ResolveSchema(parent, true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "map[string]*"+val.Expected, rt.GoType)
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

			rt, err := resolver.ResolveSchema(parent, true, true)
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

			rt, err := resolver.ResolveSchema(parent, true, true)
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

			rt, err := resolver.ResolveSchema(parent, true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "map[string]*"+val.Expected, rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}

	}
}

func TestTypeResolver_Notables(t *testing.T) {
	doc, resolver, err := specResolver(t, "../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		def := doc.Spec().Definitions["Notables"]
		rest, err := resolver.ResolveSchema(&def, false, true)
		if assert.NoError(t, err) {
			assert.True(t, rest.IsArray)
			assert.False(t, rest.IsAnonymous)
			assert.False(t, rest.IsNullable)
			assert.Equal(t, "[]*models.Notable", rest.GoType)
		}
	}
}

func specResolver(t testing.TB, path string) (*spec.Document, *typeResolver, error) {
	tlb, err := spec.Load(path)
	if err != nil {
		return nil, nil, err
	}
	resolver := &typeResolver{
		Doc:           tlb,
		ModelsPackage: "models",
	}
	resolver.KnownDefs = make(map[string]struct{})
	for k := range tlb.Spec().Definitions {
		resolver.KnownDefs[k] = struct{}{}
	}
	return tlb, resolver, nil
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
	resolver := &typeResolver{
		Doc:           tlb,
		ModelsPackage: "models",
	}

	resolver.KnownDefs = make(map[string]struct{})
	for k, sch := range swsp.Definitions {
		resolver.KnownDefs[k] = struct{}{}
		if nm, ok := sch.Extensions["x-go-name"]; ok {
			resolver.KnownDefs[nm.(string)] = struct{}{}
		}
	}
	return tlb, resolver, nil
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

		rt, err := resolver.ResolveSchema(parent, true, true)
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

		rt, err := resolver.ResolveSchema(parent, true, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsNullable)
			assert.True(t, rt.IsAnonymous)
			assert.True(t, rt.IsComplexObject)
		}

		parent.Extensions = make(spec.Extensions)
		parent.Extensions[xIsNullable] = true

		rt, err = resolver.ResolveSchema(parent, true, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsNullable)
			assert.True(t, rt.IsAnonymous)
			assert.True(t, rt.IsComplexObject)
		}

		// Also test that it's nullable with just x-nullable
		parent.Extensions[xIsNullable] = false
		parent.Extensions[xNullable] = false

		rt, err = resolver.ResolveSchema(parent, true, true)
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
	resolver.KnownDefs["TheModel"] = struct{}{}
	defer func() { resolver.ModelName = "" }()

	if assert.NoError(t, err) {
		//very poor schema definitions (as in none)
		types := []string{"object", ""}
		for _, tpe := range types {
			sch := new(spec.Schema)
			sch.Typed(tpe, "")
			rt, err := resolver.ResolveSchema(sch, true, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsMap)
				assert.False(t, rt.IsComplexObject)
				assert.Equal(t, "interface{}", rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}

			sch.Properties = make(map[string]spec.Schema)
			var ss spec.Schema
			sch.Properties["tags"] = *(&ss).CollectionOf(*spec.StringProperty())
			rt, err = resolver.ResolveSchema(sch, false, true)
			assert.True(t, rt.IsComplexObject)
			assert.False(t, rt.IsMap)
			assert.Equal(t, "models.TheModel", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)

			sch.Properties = nil
			nsch := new(spec.Schema)
			nsch.Typed(tpe, "")
			nsch.AllOf = []spec.Schema{*sch}
			rt, err = resolver.ResolveSchema(nsch, false, true)
			if assert.NoError(t, err) {
				assert.True(t, rt.IsComplexObject)
				assert.False(t, rt.IsMap)
				assert.Equal(t, "models.TheModel", rt.GoType)
				assert.Equal(t, "object", rt.SwaggerType)
			}
		}
		sch := new(spec.Schema)
		rt, err := resolver.ResolveSchema(sch, true, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsMap)
			assert.False(t, rt.IsComplexObject)
			assert.Equal(t, "interface{}", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)

		}
		sch = new(spec.Schema)
		var sp spec.Schema
		sp.Typed("object", "")
		sch.AllOf = []spec.Schema{sp}
		rt, err = resolver.ResolveSchema(sch, true, true)
		if assert.NoError(t, err) {
			assert.True(t, rt.IsComplexObject)
			assert.False(t, rt.IsMap)
			assert.Equal(t, "models.TheModel", rt.GoType)
			assert.Equal(t, "object", rt.SwaggerType)
		}
	}
}

func TestTypeResolver_AliasTypes(t *testing.T) {
	doc, resolver, err := basicTaskListResolver(t)
	if assert.NoError(t, err) {
		resolver.ModelsPackage = ""
		resolver.ModelName = "Currency"
		defer func() {
			resolver.ModelName = ""
			resolver.ModelsPackage = "models"
		}()
		defs := doc.Spec().Definitions[resolver.ModelName]
		rt, err := resolver.ResolveSchema(&defs, false, true)
		if assert.NoError(t, err) {
			assert.False(t, rt.IsAnonymous)
			assert.True(t, rt.IsAliased)
			assert.True(t, rt.IsPrimitive)
			assert.Equal(t, "Currency", rt.GoType)
			assert.Equal(t, "string", rt.AliasedType)
		}
	}
}

func assertPrimitiveResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType) {
	assert.Equal(t, tpe, tr.SwaggerType, fmt.Sprintf("expected %q (%q, %q) to for the swagger type but got %q", tpe, tfmt, exp, tr.SwaggerType))
	assert.Equal(t, tfmt, tr.SwaggerFormat, fmt.Sprintf("expected %q (%q, %q) to for the swagger format but got %q", tfmt, tpe, exp, tr.SwaggerFormat))
	assert.Equal(t, exp, tr.GoType, fmt.Sprintf("expected %q (%q, %q) to for the go type but got %q", exp, tpe, tfmt, tr.GoType))
}

func assertBuiltinResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType, i int) bool {
	return assert.Equal(t, tpe, tr.SwaggerType, fmt.Sprintf("expected %q (%q, %q) at %d for the swagger type but got %q", tpe, tfmt, exp, i, tr.SwaggerType)) &&
		assert.Equal(t, tfmt, tr.SwaggerFormat, fmt.Sprintf("expected %q (%q, %q) at %d for the swagger format but got %q", tfmt, tpe, exp, i, tr.SwaggerFormat)) &&
		assert.Equal(t, exp, tr.GoType, fmt.Sprintf("expected %q (%q, %q) at %d for the go type but got %q", exp, tpe, tfmt, i, tr.GoType))
}

type builtinVal struct {
	Type, Format, Expected string
	Default                interface{}
	Required               bool
	ReadOnly               bool
	Maximum                *float64
	ExclusiveMaximum       bool
	Minimum                *float64
	ExclusiveMinimum       bool
	MaxLength              *int64
	MinLength              *int64
	Pattern                string
	MaxItems               *int64
	MinItems               *int64
	UniqueItems            bool
	MultipleOf             *float64
	Enum                   []interface{}
	Nullable               bool
	Extensions             spec.Extensions
}

func nullableExt() spec.Extensions {
	return map[string]interface{}{"x-nullable": true}
}
func isNullableExt() spec.Extensions {
	return map[string]interface{}{"x-isnullable": true}
}
func notNullableExt() spec.Extensions {
	return map[string]interface{}{"x-nullable": false}
}
func isNotNullableExt() spec.Extensions {
	return map[string]interface{}{"x-isnullable": false}
}

var boolPointerVals = []builtinVal{
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: nullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: nullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: false, ReadOnly: false, Extensions: nullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: true, Extensions: nullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: isNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: isNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: false, ReadOnly: false, Extensions: isNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: notNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: notNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false, Extensions: notNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true, Extensions: notNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true, Extensions: notNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: true, Required: false, ReadOnly: false, Extensions: isNotNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: true, Default: nil, Required: true, ReadOnly: false, Extensions: isNotNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: false, ReadOnly: false, Extensions: isNotNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: nil, Required: true, ReadOnly: true, Extensions: isNotNullableExt()},
	builtinVal{Type: "boolean", Format: "", Expected: "bool", Nullable: false, Default: true, Required: true, ReadOnly: true, Extensions: isNotNullableExt()},
}

func generateIntPointerVals(v string) (result []builtinVal) {

	vv := v
	if vv == "" || vv == "int" {
		vv = "int64"
	}
	if vv == "uint" {
		vv = "uint64"
	}
	return []builtinVal{
		// plain vanilla
		builtinVal{Type: "integer", Format: v, Expected: vv},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Extensions: isNullableExt()}, // 2

		// plain vanilla readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 9

		// required
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Extensions: isNullableExt()}, // 12

		// required, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Default: 3, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 19

		// minimum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Minimum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 23

		// minimum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Minimum: swag.Float64(2), Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true, Minimum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 31

		// required, minimum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Extensions: isNullableExt()}, // 35

		// required, minimum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Minimum: swag.Float64(2), ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), ReadOnly: true, Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 42

		// maximum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 46

		// maximum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, ReadOnly: true, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 54

		// required, maximum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 58

		// required, maximum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Maximum: swag.Float64(2), ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), ReadOnly: true, Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Maximum: swag.Float64(2), Default: 3, ReadOnly: true, Extensions: isNullableExt()}, // 65

		// minimum and maximum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Minimum: swag.Float64(2), Maximum: swag.Float64(5)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(0), Maximum: swag.Float64(1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: isNullableExt()}, // 72

		// minimum and maximum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(0), Maximum: swag.Float64(3)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(-1), ReadOnly: true, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, Minimum: swag.Float64(-1), Maximum: swag.Float64(6)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(1), Maximum: swag.Float64(6)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Default: 3, Minimum: swag.Float64(-6), Maximum: swag.Float64(-1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 83

		// required, minimum and maximum validation
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(5)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(0), Maximum: swag.Float64(1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Minimum: swag.Float64(2), Maximum: swag.Float64(6), Extensions: isNullableExt()}, // 89

		// required, minimum and maximum validation, readonly and defaults
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(0), Maximum: swag.Float64(3)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(0)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: false, Required: true, Default: 3, Minimum: swag.Float64(-1), ReadOnly: true, Maximum: swag.Float64(2)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, ReadOnly: true, Maximum: swag.Float64(2), Extensions: isNullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(-1), Maximum: swag.Float64(6)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(1), Maximum: swag.Float64(6)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, Minimum: swag.Float64(-6), Maximum: swag.Float64(-1)},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: nullableExt()},
		builtinVal{Type: "integer", Format: v, Expected: vv, Nullable: true, Required: true, Default: 3, ReadOnly: true, Minimum: swag.Float64(-1), Maximum: swag.Float64(2), Extensions: isNullableExt()}, // 99
	}
}

func TestTypeResolver_PointerLifting(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)

	if assert.NoError(t, err) {
		// primitives and string formats
		for i, val := range boolPointerVals {
			assertBuiltinVal(t, resolver, i, val)
		}
		for _, v := range []string{"", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
			passed := true
			for i, val := range generateIntPointerVals(v) {
				//fmt.Println("trying", i)
				if !assertBuiltinVal(t, resolver, i, val) {
					passed = false
				}
			}
			if !passed {
				break
			}
		}
	}
}

func assertBuiltinVal(t *testing.T, resolver *typeResolver, i int, val builtinVal) bool {
	sch := new(spec.Schema)
	sch.Typed(val.Type, val.Format)
	sch.Default = val.Default
	sch.ReadOnly = val.ReadOnly
	sch.Extensions = val.Extensions
	sch.Minimum = val.Minimum
	sch.Maximum = val.Maximum
	sch.MultipleOf = val.MultipleOf

	rt, err := resolver.ResolveSchema(sch, true, val.Required)
	if assert.NoError(t, err) {
		if val.Nullable {
			if !assert.True(t, rt.IsNullable, "expected nullable for item at: %d", i) {
				fmt.Println("isRequired:", val.Required)
				// pretty.Println(sch)
				return false
			}
		} else {
			if !assert.False(t, rt.IsNullable, "expected not nullable for item at: %d", i) {
				fmt.Println("isRequired:", val.Required)
				// pretty.Println(sch)
				return false
			}
		}
		if !assertBuiltinResolve(t, val.Type, val.Format, val.Expected, rt, i) {
			return false
		}
	}
	return true
}
