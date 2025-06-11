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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

func schTypeVals() []struct{ Type, Format, Expected string } {
	return []struct{ Type, Format, Expected string }{
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
		{"string", "mac", "strfmt.MAC"},
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
		{"string", "ObjectId", "strfmt.ObjectId"},
		{"string", "password", "strfmt.Password"},
		{"string", "uint8", "string"},
		{"string", "uint16", "string"},
		{"string", "uint32", "string"},
		{"string", "uint64", "string"},
		{"string", "int8", "string"},
		{"string", "int16", "string"},
		{"string", "int32", "string"},
		{"string", "int64", "string"},
		{"file", "", "io.ReadCloser"},
	}
}

func schRefVals() []struct{ Type, GoType, Expected string } {
	return []struct{ Type, GoType, Expected string }{
		{"Comment", "", "models.Comment"},
		{"UserCard", "UserItem", "models.UserItem"},
	}
}

func TestTypeResolver_AdditionalItems(t *testing.T) {
	_, resolver, e := basicTaskListResolver(t)
	require.NoError(t, e)

	tpe := spec.StringProperty()

	// arrays of primitives and string formats with additional formats
	for _, val := range schTypeVals() {
		var sch spec.Schema
		sch.Typed(val.Type, val.Format)
		var coll spec.Schema
		coll.Type = []string{"array"}
		coll.Items = new(spec.SchemaOrArray)
		coll.Items.Schema = tpe
		coll.AdditionalItems = new(spec.SchemaOrBool)
		coll.AdditionalItems.Schema = &sch

		rt, err := resolver.ResolveSchema(&coll, true, true)
		require.NoError(t, err)
		require.True(t, rt.IsArray)

		assert.True(t, rt.HasAdditionalItems)
		assert.False(t, rt.IsNullable)
	}
}

func TestTypeResolver_BasicTypes(t *testing.T) {
	_, resolver, e := basicTaskListResolver(t)
	require.NoError(t, e)

	// primitives and string formats
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)
		sch.Typed(val.Type, val.Format)

		rt, err := resolver.ResolveSchema(sch, true, false)
		require.NoError(t, err)

		assert.False(t, rt.IsNullable, "expected %s with format %q to not be nullable", val.Type, val.Format)
		assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
	}

	// arrays of primitives and string formats
	for _, val := range schTypeVals() {
		var sch spec.Schema
		sch.Typed(val.Type, val.Format)
		rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
		assert.False(t, rt.IsEmptyOmitted)

		s := new(spec.Schema).CollectionOf(sch)
		s.AddExtension(xOmitEmpty, false)
		rt, err = resolver.ResolveSchema(s, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
		assert.False(t, rt.IsEmptyOmitted)

		s = new(spec.Schema).CollectionOf(sch)
		s.AddExtension(xOmitEmpty, true)
		rt, err = resolver.ResolveSchema(s, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
		assert.True(t, rt.IsEmptyOmitted)
	}

	// primitives and string formats
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)
		sch.Typed(val.Type, val.Format)
		sch.Extensions = make(spec.Extensions)
		sch.Extensions[xIsNullable] = true

		rt, err := resolver.ResolveSchema(sch, true, false)
		require.NoError(t, err)

		if val.Type == "file" {
			assert.False(t, rt.IsNullable, "expected %q (%q) to not be nullable", val.Type, val.Format)
		} else {
			assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
		}
		assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)

		// Test x-nullable overrides x-isnullable
		sch.Extensions[xIsNullable] = false
		sch.Extensions[xNullable] = true
		rt, err = resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		if val.Type == "file" {
			assert.False(t, rt.IsNullable, "expected %q (%q) to not be nullable", val.Type, val.Format)
		} else {
			assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
		}
		assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)

		// Test x-nullable without x-isnullable
		delete(sch.Extensions, xIsNullable)
		sch.Extensions[xNullable] = true
		rt, err = resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		if val.Type == "file" {
			assert.False(t, rt.IsNullable, "expected %q (%q) to not be nullable", val.Type, val.Format)
		} else {
			assert.True(t, rt.IsNullable, "expected %q (%q) to be nullable", val.Type, val.Format)
		}
		assertPrimitiveResolve(t, val.Type, val.Format, val.Expected, rt)
	}

	// arrays of primitives and string formats
	for _, val := range schTypeVals() {
		var sch spec.Schema
		sch.Typed(val.Type, val.Format)
		sch.AddExtension(xIsNullable, true)

		rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(sch), true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
	}
}

func TestTypeResolver_Refs(t *testing.T) {
	_, resolver, e := basicTaskListResolver(t)
	require.NoError(t, e)

	// referenced objects
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

		rt, err := resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		assert.Equal(t, val.Expected, rt.GoType)
		assert.False(t, rt.IsAnonymous)
		assert.True(t, rt.IsNullable)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// referenced array objects
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

		rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
		// now this behavior has moved down to the type resolver:
		// * it used to be hidden to the type resolver, but rendered like that eventually
		assert.Equal(t, "[]*"+val.Expected, rt.GoType)
	}
	// for named objects
	// referenced objects
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

		rt, err := resolver.ResolveSchema(sch, false, true)
		require.NoError(t, err)

		assert.Equal(t, val.Expected, rt.GoType)
		assert.False(t, rt.IsAnonymous)
		assert.True(t, rt.IsNullable)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// referenced array objects
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)

		rt, err := resolver.ResolveSchema(new(spec.Schema).CollectionOf(*sch), false, true)
		require.NoError(t, err)

		assert.True(t, rt.IsArray)
		// now this behavior has moved down to the type resolver:
		// * it used to be hidden to the type resolver, but rendered like that eventually
		assert.Equal(t, "[]*"+val.Expected, rt.GoType)
	}
}

func TestTypeResolver_AdditionalProperties(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	require.NoError(t, err)

	// primitives as additional properties
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)

		sch.Typed(val.Type, val.Format)
		parent := new(spec.Schema)
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = sch

		rt, err := resolver.ResolveSchema(parent, true, false)
		require.NoError(t, err)

		assert.True(t, rt.IsMap)
		assert.False(t, rt.IsComplexObject)
		assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// array of primitives as additional properties
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)

		sch.Typed(val.Type, val.Format)
		parent := new(spec.Schema)
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = new(spec.Schema).CollectionOf(*sch)

		rt, err := resolver.ResolveSchema(parent, true, false)
		require.NoError(t, err)

		assert.True(t, rt.IsMap)
		assert.False(t, rt.IsComplexObject)
		assert.Equal(t, "map[string][]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// refs as additional properties
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)
		parent := new(spec.Schema)
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = sch

		rt, err := resolver.ResolveSchema(parent, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsMap)
		assert.False(t, rt.IsComplexObject)
		assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// when additional properties and properties present, it's a complex object

	// primitives as additional properties
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)

		sch.Typed(val.Type, val.Format)
		parent := new(spec.Schema)
		parent.Properties = make(map[string]spec.Schema)
		parent.Properties["id"] = *spec.Int32Property()
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = sch

		rt, err := resolver.ResolveSchema(parent, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// array of primitives as additional properties
	for _, val := range schTypeVals() {
		sch := new(spec.Schema)

		sch.Typed(val.Type, val.Format)
		parent := new(spec.Schema)
		parent.Properties = make(map[string]spec.Schema)
		parent.Properties["id"] = *spec.Int32Property()
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = new(spec.Schema).CollectionOf(*sch)

		rt, err := resolver.ResolveSchema(parent, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "map[string][]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}

	// refs as additional properties
	for _, val := range schRefVals() {
		sch := new(spec.Schema)
		sch.Ref, _ = spec.NewRef("#/definitions/" + val.Type)
		parent := new(spec.Schema)
		parent.Properties = make(map[string]spec.Schema)
		parent.Properties["id"] = *spec.Int32Property()
		parent.AdditionalProperties = new(spec.SchemaOrBool)
		parent.AdditionalProperties.Schema = sch

		rt, err := resolver.ResolveSchema(parent, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "map[string]"+val.Expected, rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}
}

func TestTypeResolver_Notables(t *testing.T) {
	doc, resolver, err := specResolver(t, "../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	def := doc.Spec().Definitions["Notables"]
	rest, err := resolver.ResolveSchema(&def, false, true)
	require.NoError(t, err)

	assert.True(t, rest.IsArray)
	assert.False(t, rest.IsAnonymous)
	assert.False(t, rest.IsNullable)
	assert.Equal(t, "[]*models.Notable", rest.GoType)
}

func specResolver(_ testing.TB, path string) (*loads.Document, *typeResolver, error) {
	tlb, err := loads.Spec(path)
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

func basicTaskListResolver(_ testing.TB) (*loads.Document, *typeResolver, error) {
	tlb, err := loads.Spec("../fixtures/codegen/tasklist.basic.yml")
	if err != nil {
		return nil, nil, err
	}
	swsp := tlb.Spec()
	uc := swsp.Definitions["UserCard"]
	uc.AddExtension(xGoName, "UserItem")
	swsp.Definitions["UserCard"] = uc
	resolver := &typeResolver{
		Doc:           tlb,
		ModelsPackage: "models",
	}

	resolver.KnownDefs = make(map[string]struct{})
	for k, sch := range swsp.Definitions {
		resolver.KnownDefs[k] = struct{}{}
		if nm, ok := sch.Extensions[xGoName]; ok {
			resolver.KnownDefs[nm.(string)] = struct{}{}
		}
	}
	return tlb, resolver, nil
}

func TestTypeResolver_TupleTypes(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	require.NoError(t, err)

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
	require.NoError(t, err)

	assert.False(t, rt.IsArray)
	assert.True(t, rt.IsTuple)
}

func TestTypeResolver_AnonymousStructs(t *testing.T) {
	_, resolver, err := basicTaskListResolver(t)
	require.NoError(t, err)

	// anonymous structs should be accounted for
	parent := new(spec.Schema)
	parent.Typed("object", "")
	parent.Properties = make(map[string]spec.Schema)
	parent.Properties["name"] = *spec.StringProperty()
	parent.Properties["age"] = *spec.Int32Property()

	rt, err := resolver.ResolveSchema(parent, true, true)
	require.NoError(t, err)

	assert.True(t, rt.IsNullable)
	assert.True(t, rt.IsAnonymous)
	assert.True(t, rt.IsComplexObject)

	parent.Extensions = make(spec.Extensions)
	parent.Extensions[xIsNullable] = true

	rt, err = resolver.ResolveSchema(parent, true, true)
	require.NoError(t, err)

	assert.True(t, rt.IsNullable)
	assert.True(t, rt.IsAnonymous)
	assert.True(t, rt.IsComplexObject)

	// Also test that it's nullable with just x-nullable
	parent.Extensions[xIsNullable] = false
	parent.Extensions[xNullable] = false

	rt, err = resolver.ResolveSchema(parent, true, true)
	require.NoError(t, err)

	assert.False(t, rt.IsNullable)
	assert.True(t, rt.IsAnonymous)
	assert.True(t, rt.IsComplexObject)
}

func TestTypeResolver_ObjectType(t *testing.T) {
	_, resolver, e := basicTaskListResolver(t)
	require.NoError(t, e)

	resolver.ModelName = "TheModel"
	resolver.KnownDefs["TheModel"] = struct{}{}
	defer func() { resolver.ModelName = "" }()

	// very poor schema definitions (as in none)
	types := []string{"object", ""}
	for _, tpe := range types {
		sch := new(spec.Schema)
		sch.Typed(tpe, "")
		rt, err := resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsMap)
		assert.False(t, rt.IsComplexObject)
		assert.Equal(t, "interface{}", rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)

		sch.Properties = make(map[string]spec.Schema)
		var ss spec.Schema
		sch.Properties["tags"] = *(&ss).CollectionOf(*spec.StringProperty())
		rt, err = resolver.ResolveSchema(sch, false, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "models.TheModel", rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)

		sch.Properties = nil
		nsch := new(spec.Schema)
		nsch.Typed(tpe, "")
		nsch.AllOf = []spec.Schema{*sch}
		rt, err = resolver.ResolveSchema(nsch, false, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "models.TheModel", rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)

		sch = new(spec.Schema)
		rt, err = resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsMap)
		assert.False(t, rt.IsComplexObject)
		assert.Equal(t, "interface{}", rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)

		sch = new(spec.Schema)
		var sp spec.Schema
		sp.Typed("object", "")
		sch.AllOf = []spec.Schema{sp}
		rt, err = resolver.ResolveSchema(sch, true, true)
		require.NoError(t, err)

		assert.True(t, rt.IsComplexObject)
		assert.False(t, rt.IsMap)
		assert.Equal(t, "models.TheModel", rt.GoType)
		assert.Equal(t, "object", rt.SwaggerType)
	}
}

func TestTypeResolver_AliasTypes(t *testing.T) {
	doc, resolver, err := basicTaskListResolver(t)
	require.NoError(t, err)

	resolver.ModelsPackage = ""
	resolver.ModelName = "Currency"
	defer func() {
		resolver.ModelName = ""
		resolver.ModelsPackage = "models"
	}()

	defs := doc.Spec().Definitions[resolver.ModelName]
	rt, err := resolver.ResolveSchema(&defs, false, true)
	require.NoError(t, err)

	assert.False(t, rt.IsAnonymous)
	assert.True(t, rt.IsAliased)
	assert.True(t, rt.IsPrimitive)
	assert.Equal(t, "Currency", rt.GoType)
	assert.Equal(t, "string", rt.AliasedType)
}

func assertPrimitiveResolve(t testing.TB, tpe, tfmt, exp string, tr resolvedType) {
	assert.Equal(t, tpe, tr.SwaggerType, "expected %q (%q, %q) to for the swagger type but got %q", tpe, tfmt, exp, tr.SwaggerType)
	assert.Equal(t, tfmt, tr.SwaggerFormat, "expected %q (%q, %q) to for the swagger format but got %q", tfmt, tpe, exp, tr.SwaggerFormat)
	assert.Equal(t, exp, tr.GoType, "expected %q (%q, %q) to for the go type but got %q", exp, tpe, tfmt, tr.GoType)
}

func TestTypeResolver_ExistingModel(t *testing.T) {
	doc, err := loads.Spec("../fixtures/codegen/existing-model.yml")
	resolver := newTypeResolver("model", "", doc)
	require.NoError(t, err)

	def := doc.Spec().Definitions["JsonWebKey"]
	tpe, pkg, alias := resolver.knownDefGoType("JsonWebKey", def, nil)
	assert.Equal(t, "jwk.Key", tpe)
	assert.Equal(t, "github.com/user/package", pkg)
	assert.Equal(t, "jwk", alias)
	rest, err := resolver.ResolveSchema(&def, false, true)
	require.NoError(t, err)

	assert.False(t, rest.IsMap)
	assert.False(t, rest.IsArray)
	assert.False(t, rest.IsTuple)
	assert.False(t, rest.IsStream)
	assert.True(t, rest.IsAliased)
	assert.False(t, rest.IsBaseType)
	assert.False(t, rest.IsInterface)
	assert.True(t, rest.IsNullable)
	assert.False(t, rest.IsPrimitive)
	assert.False(t, rest.IsAnonymous)
	assert.True(t, rest.IsComplexObject)
	assert.False(t, rest.IsCustomFormatter)
	assert.Equal(t, "jwk.Key", rest.GoType)
	assert.Equal(t, "github.com/user/package", rest.Pkg)
	assert.Equal(t, "jwk", rest.PkgAlias)

	def = doc.Spec().Definitions["JsonWebKeySet"].Properties["keys"]
	rest, err = resolver.ResolveSchema(&def, false, true)
	require.NoError(t, err)

	assert.False(t, rest.IsMap)
	assert.True(t, rest.IsArray)
	assert.False(t, rest.IsTuple)
	assert.False(t, rest.IsStream)
	assert.False(t, rest.IsAliased)
	assert.False(t, rest.IsBaseType)
	assert.False(t, rest.IsInterface)
	assert.False(t, rest.IsNullable)
	assert.False(t, rest.IsPrimitive)
	assert.False(t, rest.IsAnonymous)
	assert.False(t, rest.IsComplexObject)
	assert.False(t, rest.IsCustomFormatter)
	assert.Equal(t, "[]*jwk.Key", rest.GoType)
	assert.Empty(t, rest.Pkg)
	assert.Empty(t, rest.PkgAlias)
}
