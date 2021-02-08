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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type templateTest struct {
	t        testing.TB
	template *template.Template
}

func (tt *templateTest) assertRender(data interface{}, expected string) (success bool) {
	buf := bytes.NewBuffer(nil)
	defer func() {
		success = !tt.t.Failed()
	}()

	require.NoError(tt.t, tt.template.Execute(buf, data))

	trimmed := strings.TrimLeft(buf.String(), "\n\t ")
	exp := strings.TrimLeft(expected, "\n\t ")
	assert.Equal(tt.t, exp, trimmed)

	return
}

func TestGenerateModel_Sanity(t *testing.T) {
	t.Parallel()

	// just checks if it can render and format these things
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions

	t.Run("mode sanity check", func(t *testing.T) {
		for k, schema := range definitions {
			opts := opts()
			genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
			require.NoError(t, err)

			rendered := bytes.NewBuffer(nil)

			require.NoErrorf(t, opts.templates.MustGet("model").Execute(rendered, genModel),
				"Unexpected error while rendering models for fixtures/codegen/todolist.models.yml: %v", err)

			_, err = opts.LanguageOpts.FormatContent(strings.ToLower(k)+".go", rendered.Bytes())
			require.NoError(t, err)
		}
	})
}

func TestGenerateModel_DocString(t *testing.T) {
	funcMap := FuncMapFunc(DefaultLanguageFunc())
	templ := template.Must(template.New("docstring").Funcs(funcMap).Parse(string(assets["docstring.gotmpl"])))
	tt := templateTest{t, templ}

	var gmp GenSchema
	gmp.Title = "The title of the property"
	gmp.Description = "The description of the property"
	var expected = `The title of the property
//
// The description of the property`
	tt.assertRender(gmp, expected)

	gmp.Title = ""
	expected = `The description of the property`
	tt.assertRender(gmp, expected)

	gmp.Description = ""
	gmp.Name = "theModel"
	expected = `the model`
	tt.assertRender(gmp, expected)
}

func TestGenerateModel_PropertyValidation(t *testing.T) {
	funcMap := FuncMapFunc(DefaultLanguageFunc())
	templ := template.Must(template.New("propertyValidationDocString").Funcs(funcMap).Parse(string(assets["validation/structfield.gotmpl"])))
	tt := templateTest{t, templ}

	var gmp GenSchema
	gmp.Required = true
	tt.assertRender(gmp, `
// Required: true`)
	var fl float64 = 10
	var in1 int64 = 20
	var in2 int64 = 30
	gmp.Maximum = &fl
	gmp.ExclusiveMaximum = true
	gmp.Minimum = &fl
	gmp.ExclusiveMinimum = true
	gmp.MaxLength = &in1
	gmp.MinLength = &in1
	gmp.Pattern = "\\w[\\w- ]+"
	gmp.MaxItems = &in2
	gmp.MinItems = &in2
	gmp.UniqueItems = true

	tt.assertRender(gmp, `
// Required: true
// Maximum: < 10
// Minimum: > 10
// Max Length: 20
// Min Length: 20
// Pattern: \w[\w- ]+
// Max Items: 30
// Min Items: 30
// Unique: true`)

	gmp.Required = false
	gmp.ExclusiveMaximum = false
	gmp.ExclusiveMinimum = false
	tt.assertRender(gmp, `
// Maximum: 10
// Minimum: 10
// Max Length: 20
// Min Length: 20
// Pattern: \w[\w- ]+
// Max Items: 30
// Min Items: 30
// Unique: true`)

}

func TestGenerateModel_SchemaField(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("structfield")}

	var gmp GenSchema
	gmp.Name = "some name"
	gmp.OriginalName = "some name"
	gmp.resolvedType = resolvedType{GoType: "string", IsPrimitive: true, IsEmptyOmitted: true}
	gmp.Title = "The title of the property"
	gmp.CustomTag = "mytag:\"foobar,foobaz\""

	tt.assertRender(&gmp, `// The title of the property
`+"SomeName string `json:\"some name,omitempty\" mytag:\"foobar,foobaz\"`\n")

	var fl float64 = 10
	var in1 int64 = 20
	var in2 int64 = 30

	gmp.Description = "The description of the property"
	gmp.Required = true
	gmp.Maximum = &fl
	gmp.ExclusiveMaximum = true
	gmp.Minimum = &fl
	gmp.ExclusiveMinimum = true
	gmp.MaxLength = &in1
	gmp.MinLength = &in1
	gmp.Pattern = "\\w[\\w- ]+"
	gmp.MaxItems = &in2
	gmp.MinItems = &in2
	gmp.UniqueItems = true
	gmp.ReadOnly = true
	gmp.StructTags = []string{"json", "db", "example"}
	gmp.Example = "some example\""
	tt.assertRender(&gmp, `// The title of the property
//
// The description of the property
// Example: some example"
// Required: true
// Read Only: true
// Maximum: < 10
// Minimum: > 10
// Max Length: 20
// Min Length: 20
// Pattern: \w[\w- ]+
// Max Items: 30
// Min Items: 30
// Unique: true
`+"SomeName string `json:\"some name\" db:\"some name\" example:\"some example\\\"\" mytag:\"foobar,foobaz\"`\n")

	gmp.Example = "some example``"
	tt.assertRender(&gmp, `// The title of the property
//
// The description of the property
// Example: some example`+"``"+`
// Required: true
// Read Only: true
// Maximum: < 10
// Minimum: > 10
// Max Length: 20
// Min Length: 20
// Pattern: \w[\w- ]+
// Max Items: 30
// Min Items: 30
// Unique: true
`+"SomeName string \"json:\\\"some name\\\" db:\\\"some name\\\" example:\\\"some example``\\\" mytag:\\\"foobar,foobaz\\\"\"\n")
}

var schTypeGenDataSimple = []struct {
	Value    GenSchema
	Expected string
}{
	{GenSchema{resolvedType: resolvedType{GoType: "string", IsPrimitive: true}}, "string"},
	{GenSchema{resolvedType: resolvedType{GoType: "string", IsPrimitive: true, IsNullable: true}}, "*string"},
	{GenSchema{resolvedType: resolvedType{GoType: "bool", IsPrimitive: true}}, "bool"},
	{GenSchema{resolvedType: resolvedType{GoType: "int32", IsPrimitive: true}}, "int32"},
	{GenSchema{resolvedType: resolvedType{GoType: "int64", IsPrimitive: true}}, "int64"},
	{GenSchema{resolvedType: resolvedType{GoType: "float32", IsPrimitive: true}}, "float32"},
	{GenSchema{resolvedType: resolvedType{GoType: "float64", IsPrimitive: true}}, "float64"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Base64", IsPrimitive: true}}, "strfmt.Base64"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Date", IsPrimitive: true}}, "strfmt.Date"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.DateTime", IsPrimitive: true}}, "strfmt.DateTime"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.URI", IsPrimitive: true}}, "strfmt.URI"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Email", IsPrimitive: true}}, "strfmt.Email"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Hostname", IsPrimitive: true}}, "strfmt.Hostname"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.IPv4", IsPrimitive: true}}, "strfmt.IPv4"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.IPv6", IsPrimitive: true}}, "strfmt.IPv6"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.UUID", IsPrimitive: true}}, "strfmt.UUID"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.UUID3", IsPrimitive: true}}, "strfmt.UUID3"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.UUID4", IsPrimitive: true}}, "strfmt.UUID4"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.UUID5", IsPrimitive: true}}, "strfmt.UUID5"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.ISBN", IsPrimitive: true}}, "strfmt.ISBN"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.ISBN10", IsPrimitive: true}}, "strfmt.ISBN10"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.ISBN13", IsPrimitive: true}}, "strfmt.ISBN13"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.CreditCard", IsPrimitive: true}}, "strfmt.CreditCard"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.SSN", IsPrimitive: true}}, "strfmt.SSN"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.HexColor", IsPrimitive: true}}, "strfmt.HexColor"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.RGBColor", IsPrimitive: true}}, "strfmt.RGBColor"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Duration", IsPrimitive: true}}, "strfmt.Duration"},
	{GenSchema{resolvedType: resolvedType{GoType: "strfmt.Password", IsPrimitive: true}}, "strfmt.Password"},
	{GenSchema{resolvedType: resolvedType{GoType: "io.ReadCloser", IsStream: true}}, "io.ReadCloser"},
	{GenSchema{resolvedType: resolvedType{GoType: "interface{}", IsInterface: true}}, "interface{}"},
	{GenSchema{resolvedType: resolvedType{GoType: "[]int32", IsArray: true}}, "[]int32"},
	{GenSchema{resolvedType: resolvedType{GoType: "[]string", IsArray: true}}, "[]string"},
	{GenSchema{resolvedType: resolvedType{GoType: "map[string]int32", IsMap: true}}, "map[string]int32"},
	{GenSchema{resolvedType: resolvedType{GoType: "models.Task", IsComplexObject: true, IsNullable: true, IsAnonymous: false}}, "*models.Task"},
}

func TestGenSchemaType(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("schemaType")}
	for _, v := range schTypeGenDataSimple {
		tt.assertRender(v.Value, v.Expected)
	}
}
func TestGenerateModel_Primitives(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("schema")}
	for _, v := range schTypeGenDataSimple {
		v.Value.IncludeValidator = true
		v.Value.IncludeModel = true
		val := v.Value
		val.ReceiverName = "o"
		if val.IsComplexObject {
			continue
		}
		val.Name = "theType"
		exp := v.Expected
		if val.IsInterface || val.IsStream {
			tt.assertRender(&val, "type TheType "+exp+"\n  \n")
			continue
		}
		tt.assertRender(&val, "type TheType "+exp+"\n  \n// Validate validates this the type\nfunc (o theType) Validate(formats strfmt.Registry) error {\n  return nil\n}\n// ContextValidate validates this the type based on context it is used \nfunc (o theType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {\n  return nil\n}\n")
	}
}

func TestGenerateModel_Zeroes(t *testing.T) {
	for _, v := range schTypeGenDataSimple {
		switch v.Value.GoType {
		// verifying Zero for primitive
		case "string":
			assert.Equal(t, `""`, v.Value.Zero())
		case "bool":
			assert.Equal(t, `false`, v.Value.Zero())
		case "int32", "int64", "float32", "float64":
			assert.Equal(t, `0`, v.Value.Zero())
		// verifying Zero for primitive formatters
		case "strfmt.Date", "strfmt.DateTime", "strfmt.OjbectId": // akin to structs
			rex := regexp.MustCompile(regexp.QuoteMeta(v.Value.GoType) + `{}`)
			assert.True(t, rex.MatchString(v.Value.Zero()))
			k := v.Value
			k.IsAliased = true
			k.AliasedType = k.GoType
			k.GoType = "myAliasedType"
			rex = regexp.MustCompile(regexp.QuoteMeta(k.GoType+"("+k.AliasedType) + `{}` + `\)`)
			assert.True(t, rex.MatchString(k.Zero()))
		case "strfmt.Duration": // akin to integer
			rex := regexp.MustCompile(regexp.QuoteMeta(v.Value.GoType) + `\(\d*\)`)
			assert.True(t, rex.MatchString(v.Value.Zero()))
			k := v.Value
			k.IsAliased = true
			k.AliasedType = k.GoType
			k.GoType = "myAliasedType"
			rex = regexp.MustCompile(regexp.QuoteMeta(k.GoType+"("+k.AliasedType) + `\(\d*\)` + `\)`)
			assert.True(t, rex.MatchString(k.Zero()))
		case "strfmt.Base64": // akin to []byte
			rex := regexp.MustCompile(regexp.QuoteMeta(v.Value.GoType) + `\(\[\]byte.*\)`)
			assert.True(t, rex.MatchString(v.Value.Zero()))
			k := v.Value
			k.IsAliased = true
			k.AliasedType = k.GoType
			k.GoType = "myAliasedType"
			rex = regexp.MustCompile(regexp.QuoteMeta(k.GoType+"("+k.AliasedType) + `\(\[\]byte.*\)` + `\)`)
			assert.True(t, rex.MatchString(k.Zero()))
		case "interface{}":
			assert.Equal(t, `nil`, v.Value.Zero())
		case "io.ReadCloser":
			continue
		default:
			switch {
			case strings.HasPrefix(v.Value.GoType, "[]") || strings.HasPrefix(v.Value.GoType, "map["): // akin to slice or map
				assert.True(t, strings.HasPrefix(v.Value.Zero(), "make("))

			case strings.HasPrefix(v.Value.GoType, "models."):
				assert.True(t, strings.HasPrefix(v.Value.Zero(), "new("))

			default: // akin to string
				rex := regexp.MustCompile(regexp.QuoteMeta(v.Value.GoType) + `\(".*"\)`)
				assert.True(t, rex.MatchString(v.Value.Zero()))
				k := v.Value
				k.IsAliased = true
				k.AliasedType = k.GoType
				k.GoType = "myAliasedType"
				rex = regexp.MustCompile(regexp.QuoteMeta(k.GoType+"("+k.AliasedType) + `\(".*"\)` + `\)`)
				assert.True(t, rex.MatchString(k.Zero()))
			}
		}
	}
}
func TestGenerateModel_Nota(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Nota"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type Nota map[string]int32", res)
}

func TestGenerateModel_NotaWithRef(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NotaWithRef"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("nota_with_ref.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type NotaWithRef map[string]Notable", res)
}

func TestGenerateModel_NotaWithMeta(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NotaWithMeta"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("nota_with_meta.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type NotaWithMeta map[string]NotaWithMetaAnon", res)
	assertInCode(t, "type NotaWithMetaAnon struct {", res)
	assertInCode(t, "Comment *string `json:\"comment\"`", res)
	assertInCode(t, "Count int32 `json:\"count,omitempty\"`", res)
}

func TestGenerateModel_RunParameters(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "RunParameters"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.IsAdditionalProperties)
	assert.True(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsMap)
	assert.False(t, genModel.IsAnonymous)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "BranchName string `json:\"branch_name,omitempty\"`", res)
	assertInCode(t, "CommitSha string `json:\"commit_sha,omitempty\"`", res)
	assertInCode(t, "Refs interface{} `json:\"refs,omitempty\"`", res)
}

func TestGenerateModel_NotaWithName(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NotaWithName"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.True(t, genModel.IsAdditionalProperties)
	assert.False(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsMap)
	assert.False(t, genModel.IsAnonymous)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, k+" map[string]int32 `json:\"-\"`", res)
	assertInCode(t, "Name *string `json:\"name\"`", res)
	assertInCode(t, k+") UnmarshalJSON", res)
	assertInCode(t, k+") MarshalJSON", res)
	assertInCode(t, "json.Marshal(stage1)", res)
	assertInCode(t, "stage1.Name = m.Name", res)
	assertInCode(t, "json.Marshal(m."+k+")", res)
	assertInCode(t, "json.Unmarshal(data, &stage1)", res)
	assertInCode(t, "json.Unmarshal(data, &stage2)", res)
	assertInCode(t, "json.Unmarshal(v, &toadd)", res)
	assertInCode(t, "result[k] = toadd", res)
	assertInCode(t, "m."+k+" = result", res)
	for _, p := range genModel.Properties {
		assertInCode(t, "delete(stage2, \""+p.Name+"\")", res)
	}
}

func TestGenerateModel_NotaWithRefRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NotaWithRefRegistry"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("nota_with_ref_registry.go", buf.Bytes())
	require.NoError(t, err)

	assertInCode(t, "type "+k+" map[string]map[string]map[string]Notable", string(ff))
}

func TestGenerateModel_WithCustomTag(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithCustomTag"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	assertInCode(t, "mytag:\"foo,bar\"", buf.String())
}

func TestGenerateModel_NotaWithMetaRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NotaWithMetaRegistry"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("nota_with_meta_registry.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type "+k+" map[string]map[string]map[string]NotaWithMetaRegistryAnon", res)
	assertInCode(t, "type NotaWithMetaRegistryAnon struct {", res)
	assertInCode(t, "Comment *string `json:\"comment\"`", res)
	assertInCode(t, "Count int32 `json:\"count,omitempty\"`", res)
}

func TestGenerateModel_WithMap(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithMap"]
	opts := opts()
	genModel, err := makeGenDefinition("WithMap", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithMap struct {", res)
	assertInCode(t, "Data map[string]string `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithMapInterface(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithMapInterface"]
	opts := opts()
	genModel, err := makeGenDefinition("WithMapInterface", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "extraInfo")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	assert.Equal(t, "map[string]interface{}", prop.GoType)
	assert.True(t, prop.Required)
	assert.True(t, prop.HasValidations)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithMapInterface struct {", res)
	assertInCode(t, "ExtraInfo map[string]interface{} `json:\"extraInfo\"`", res)
}

func TestGenerateModel_WithMapRef(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithMapRef"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "Data map[string]Notable `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithMapComplex(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithMapComplex"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "Data map[string]"+k+"DataAnon `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithMapRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithMapRegistry"]
	opts := opts()
	genModel, err := makeGenDefinition("WithMap", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithMap struct {", res)
	assertInCode(t, "Data map[string]map[string]map[string]string `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithMapRegistryRef(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithMapRegistryRef"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "Data map[string]map[string]map[string]Notable `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithMapComplexRegistry(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithMapComplexRegistry"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.False(t, genModel.HasAdditionalProperties)
	prop := getDefinitionProperty(genModel, "data")
	assert.True(t, prop.HasAdditionalProperties)
	assert.True(t, prop.IsMap)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "Data map[string]map[string]map[string]"+k+"DataAnon `json:\"data,omitempty\"`", res)
}

func TestGenerateModel_WithAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithAdditional"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.NotEmpty(t, genModel.ExtraSchemas)

	assert.False(t, genModel.HasAdditionalProperties)
	assert.False(t, genModel.IsMap)
	assert.False(t, genModel.IsAdditionalProperties)
	assert.True(t, genModel.IsComplexObject)

	sch := genModel.ExtraSchemas[0]
	assert.True(t, sch.HasAdditionalProperties)
	assert.False(t, sch.IsMap)
	assert.True(t, sch.IsAdditionalProperties)
	assert.False(t, sch.IsComplexObject)

	require.NotNil(t, sch.AdditionalProperties)

	prop := findProperty(genModel.Properties, "data")
	assert.False(t, prop.HasAdditionalProperties)
	assert.False(t, prop.IsMap)
	assert.False(t, prop.IsAdditionalProperties)
	assert.True(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "Data *"+k+"Data `json:\"data,omitempty\"`", res)
	assertInCode(t, "type "+k+"Data struct {", res)
	assertInCode(t, k+"Data map[string]string `json:\"-\"`", res)
	assertInCode(t, "Name *string `json:\"name\"`", res)
	assertInCode(t, k+"Data) UnmarshalJSON", res)
	assertInCode(t, k+"Data) MarshalJSON", res)
	assertInCode(t, "json.Marshal(stage1)", res)
	assertInCode(t, "stage1.Name = m.Name", res)
	assertInCode(t, "json.Marshal(m."+k+"Data)", res)
	assertInCode(t, "json.Unmarshal(data, &stage1)", res)
	assertInCode(t, "json.Unmarshal(data, &stage2)", res)
	assertInCode(t, "json.Unmarshal(v, &toadd)", res)
	assertInCode(t, "result[k] = toadd", res)
	assertInCode(t, "m."+k+"Data = result", res)

	for _, p := range sch.Properties {
		assertInCode(t, "delete(stage2, \""+p.Name+"\")", res)
	}
}

func TestGenerateModel_JustRef(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("schema")}
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["JustRef"]
	opts := opts()
	genModel, err := makeGenDefinition("JustRef", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.NotEmpty(t, genModel.AllOf)
	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "JustRef", genModel.Name)
	assert.Equal(t, "JustRef", genModel.GoType)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type JustRef struct {", res)
	assertInCode(t, "Notable", res)
}

func TestGenerateModel_WithRef(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("schema")}
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithRef"]
	opts := opts()
	genModel, err := makeGenDefinition("WithRef", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "WithRef", genModel.Name)
	assert.Equal(t, "WithRef", genModel.GoType)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithRef struct {", res)
	assertInCode(t, "Notes *Notable `json:\"notes,omitempty\"`", res)
}

func TestGenerateModel_WithNullableRef(t *testing.T) {
	tt := templateTest{t, templates.MustGet("model").Lookup("schema")}
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithNullableRef"]
	opts := opts()
	genModel, err := makeGenDefinition("WithNullableRef", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "WithNullableRef", genModel.Name)
	assert.Equal(t, "WithNullableRef", genModel.GoType)
	prop := getDefinitionProperty(genModel, "notes")
	assert.True(t, prop.IsNullable)
	assert.True(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithNullableRef struct {", res)
	assertInCode(t, "Notes *Notable `json:\"notes,omitempty\"`", res)
}

func TestGenerateModel_Scores(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Scores"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("scores.go", buf.Bytes())
	require.NoError(t, err)

	assertInCode(t, "type Scores []float32", string(ff))
}

func TestGenerateModel_JaggedScores(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "JaggedScores"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("jagged_scores.go", buf.Bytes())
	require.NoError(t, err)

	assertInCode(t, "type JaggedScores [][][]float32", string(ff))
}

func TestGenerateModel_Notables(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Notables"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Equal(t, "[]*Notable", genModel.GoType)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("notables.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type Notables []*Notable", res)
}

func TestGenerateModel_Notablix(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Notablix"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("notablix.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type Notablix [][][]*Notable", res)
}

func TestGenerateModel_Stats(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Stats"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("stats.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type Stats []*StatsItems0", res)
	assertInCode(t, "type StatsItems0 struct {", res)
	assertInCode(t, "Points []int64 `json:\"points\"`", res)
}

func TestGenerateModel_Statix(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Statix"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("statix.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "type Statix [][][]*StatixItems0", res)
	assertInCode(t, "type StatixItems0 struct {", res)
	assertInCode(t, "Points []int64 `json:\"points\"`", res)
}

func TestGenerateModel_WithItems(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithItems"]
	opts := opts()
	tt := templateTest{t, opts.templates.MustGet("model").Lookup("schema")}

	genModel, err := makeGenDefinition("WithItems", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Nil(t, genModel.Items)
	assert.True(t, genModel.IsComplexObject)
	prop := getDefinitionProperty(genModel, "tags")

	assert.NotNil(t, prop.Items)
	assert.True(t, prop.IsArray)
	assert.False(t, prop.IsComplexObject)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "type WithItems struct {", res)
	assertInCode(t, "Tags []string `json:\"tags\"`", res)
}

func TestGenerateModel_WithComplexItems(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithComplexItems"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Nil(t, genModel.Items)
	assert.True(t, genModel.IsComplexObject)
	prop := getDefinitionProperty(genModel, "tags")
	assert.NotNil(t, prop.Items)
	assert.True(t, prop.IsArray)
	assert.False(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)

	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	b, err := opts.LanguageOpts.FormatContent("with_complex_items.go", buf.Bytes())
	require.NoError(t, err)

	res := string(b)
	assertInCode(t, "type WithComplexItems struct {", res)
	assertInCode(t, "type WithComplexItemsTagsItems0 struct {", res)
	assertInCode(t, "Tags []*WithComplexItemsTagsItems0 `json:\"tags\"`", res)
}

func TestGenerateModel_WithItemsAndAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithItemsAndAdditional"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Nil(t, genModel.Items)
	assert.True(t, genModel.IsComplexObject)
	prop := getDefinitionProperty(genModel, "tags")
	assert.True(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	b, err := opts.LanguageOpts.FormatContent("with_complex_items.go", buf.Bytes())
	require.NoError(t, err)

	res := string(b)
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "type "+k+"TagsTuple0 struct {", res)
	// this would fail if it accepts additionalItems because it would come out as []interface{}
	assertInCode(t, "Tags *"+k+"TagsTuple0 `json:\"tags,omitempty\"`", res)
	assertInCode(t, "P0 *string `json:\"-\"`", res)
	assertInCode(t, k+"TagsTuple0Items []interface{} `json:\"-\"`", res)
}

func TestGenerateModel_WithItemsAndAdditional2(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithItemsAndAdditional2"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Nil(t, genModel.Items)
	assert.True(t, genModel.IsComplexObject)
	prop := getDefinitionProperty(genModel, "tags")
	assert.True(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	b, err := opts.LanguageOpts.FormatContent("with_complex_items.go", buf.Bytes())
	require.NoError(t, err)

	res := string(b)
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "type "+k+"TagsTuple0 struct {", res)
	// this would fail if it accepts additionalItems because it would come out as []interface{}
	assertInCode(t, "P0 *string `json:\"-\"`", res)
	assertInCode(t, "Tags *"+k+"TagsTuple0 `json:\"tags,omitempty\"`", res)
	assertInCode(t, k+"TagsTuple0Items []int32 `json:\"-\"`", res)
}

func TestGenerateModel_WithComplexAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithComplexAdditional"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Nil(t, genModel.Items)
	assert.True(t, genModel.IsComplexObject)
	prop := getDefinitionProperty(genModel, "tags")
	assert.True(t, prop.IsComplexObject)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	b, err := opts.LanguageOpts.FormatContent("with_complex_additional.go", buf.Bytes())
	require.NoError(t, err)

	res := string(b)
	assertInCode(t, "type WithComplexAdditional struct {", res)
	assertInCode(t, "type WithComplexAdditionalTagsTuple0 struct {", res)
	assertInCode(t, "Tags *WithComplexAdditionalTagsTuple0 `json:\"tags,omitempty\"`", res)
	assertInCode(t, "P0 *string `json:\"-\"`", res)
	assertInCode(t, "WithComplexAdditionalTagsTuple0Items []*WithComplexAdditionalTagsItems `json:\"-\"`", res)
}

func TestGenerateModel_SimpleTuple(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SimpleTuple"
	schema := definitions[k]
	opts := opts()

	tt := templateTest{t, opts.templates.MustGet("model")}
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Len(t, genModel.ExtraSchemas, 1)

	// NOTE: with PR#1592, an extra schema is added here because of the allOf tuple element.
	// This uncovers another issue with special AllOfs (e.g. allOf [ ..., x-nullable:true ])
	// TODO(fredbi): fix liftSpecialAllOf() to revert to: assert.Empty(t, genModel.ExtraSchemas)
	assert.True(t, genModel.IsTuple)
	assert.False(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsArray)
	assert.False(t, genModel.IsAnonymous)
	assert.Equal(t, k, genModel.Name)
	assert.Equal(t, k, genModel.GoType)
	assert.Len(t, genModel.Properties, 5)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	res := buf.String()
	assertInCode(t, "swagger:model "+k, res)
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "P0 *int64 `json:\"-\"`", res)
	assertInCode(t, "P1 *string `json:\"-\"`", res)
	assertInCode(t, "P2 *strfmt.DateTime `json:\"-\"`", res)
	assertInCode(t, "P3 *Notable `json:\"-\"`", res)
	// NOTE: with PR#1592, an extra schema is added here because of the allOf tuple element.
	// This uncovers another issue with special AllOfs (e.g. allOf [ ..., x-nullable:true ])
	// TODO(fredbi): fix liftSpecialAllOf() to revert to: assert.Empty(t, genModel.ExtraSchemas)
	// assertInCode(t, "P4 *Notable `json:\"-\"`", res)
	assertInCode(t, "P4 *SimpleTupleItems4 `json:\"-\"`", res)
	assertInCode(t, k+") UnmarshalJSON", res)
	assertInCode(t, k+") MarshalJSON", res)
	assertInCode(t, "json.Marshal(data)", res)
	assert.NotRegexp(t, regexp.MustCompile("lastIndex"), res)

	for i, p := range genModel.Properties {
		m := "m.P" + strconv.Itoa(i)
		r := "&dataP" + strconv.Itoa(i)
		var rr string
		if !p.IsNullable {
			rr = "dataP" + strconv.Itoa(i)
		} else {
			rr = r
		}
		assertInCode(t, fmt.Sprintf("buf = bytes.NewBuffer(stage1[%d])", i), res)
		assertInCode(t, fmt.Sprintf("dec.Decode(%s)", r), res)
		assertInCode(t, "P"+strconv.Itoa(i)+",", res)
		assertInCode(t, fmt.Sprintf("%s = %s", m, rr), res)
	}
}

func TestGenerateModel_TupleWithExtra(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "TupleWithExtra"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Empty(t, genModel.ExtraSchemas)

	assert.True(t, genModel.IsTuple)
	assert.False(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsArray)
	assert.False(t, genModel.IsAnonymous)
	assert.True(t, genModel.HasAdditionalItems)
	assert.NotNil(t, genModel.AdditionalItems)
	assert.Equal(t, k, genModel.Name)
	assert.Equal(t, k, genModel.GoType)
	assert.Len(t, genModel.Properties, 4)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("tuple_with_extra.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "swagger:model "+k, res)
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "P0 *int64 `json:\"-\"`", res)
	assertInCode(t, "P1 *string `json:\"-\"`", res)
	assertInCode(t, "P2 *strfmt.DateTime `json:\"-\"`", res)
	assertInCode(t, "P3 *Notable `json:\"-\"`", res)
	assertInCode(t, k+"Items []float64 `json:\"-\"`", res)
	assertInCode(t, k+") UnmarshalJSON", res)
	assertInCode(t, k+") MarshalJSON", res)

	for i, p := range genModel.Properties {
		m := "m.P" + strconv.Itoa(i)
		r := "&dataP" + strconv.Itoa(i)
		var rr string
		if !p.IsNullable {
			rr = "dataP" + strconv.Itoa(i)
		} else {
			rr = r
		}
		assertInCode(t, fmt.Sprintf("lastIndex = %d", i), res)
		assertInCode(t, fmt.Sprintf("buf = bytes.NewBuffer(stage1[%d])", i), res)
		assertInCode(t, "dec := json.NewDecoder(buf)", res)
		assertInCode(t, fmt.Sprintf("dec.Decode(%s)", r), res)
		assertInCode(t, "P"+strconv.Itoa(i)+",", res)
		assertInCode(t, fmt.Sprintf("%s = %s", m, rr), res)
	}

	assertInCode(t, "var lastIndex int", res)
	assertInCode(t, "var toadd float64", res)
	assertInCode(t, "for _, val := range stage1[lastIndex+1:]", res)
	assertInCode(t, "buf = bytes.NewBuffer(val)", res)
	assertInCode(t, "dec := json.NewDecoder(buf)", res)
	assertInCode(t, "dec.Decode(&toadd)", res)
	assertInCode(t, "json.Marshal(data)", res)
	assertInCode(t, "for _, v := range m."+k+"Items", res)
}

func TestGenerateModel_TupleWithComplex(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "TupleWithComplex"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.True(t, genModel.IsTuple)
	assert.False(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsArray)
	assert.False(t, genModel.IsAnonymous)
	assert.True(t, genModel.HasAdditionalItems)
	assert.NotNil(t, genModel.AdditionalItems)
	assert.Equal(t, k, genModel.Name)
	assert.Equal(t, k, genModel.GoType)
	assert.Len(t, genModel.Properties, 4)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("tuple_with_extra.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "swagger:model "+k, res)
	assertInCode(t, "type "+k+" struct {", res)
	assertInCode(t, "P0 *int64 `json:\"-\"`", res)
	assertInCode(t, "P1 *string `json:\"-\"`", res)
	assertInCode(t, "P2 *strfmt.DateTime `json:\"-\"`", res)
	assertInCode(t, "P3 *Notable `json:\"-\"`", res)
	assertInCode(t, k+"Items []*TupleWithComplexItems `json:\"-\"`", res)
	assertInCode(t, k+") UnmarshalJSON", res)
	assertInCode(t, k+") MarshalJSON", res)

	for i, p := range genModel.Properties {
		m := "m.P" + strconv.Itoa(i)
		r := "&dataP" + strconv.Itoa(i)
		var rr string
		if !p.IsNullable {
			rr = "dataP" + strconv.Itoa(i)
		} else {
			rr = r
		}
		assertInCode(t, fmt.Sprintf("lastIndex = %d", i), res)
		assertInCode(t, fmt.Sprintf("buf = bytes.NewBuffer(stage1[%d])", i), res)
		assertInCode(t, "dec := json.NewDecoder(buf)", res)
		assertInCode(t, fmt.Sprintf("dec.Decode(%s)", r), res)
		assertInCode(t, "P"+strconv.Itoa(i)+",", res)
		assertInCode(t, fmt.Sprintf("%s = %s", m, rr), res)
	}

	assertInCode(t, "var lastIndex int", res)
	assertInCode(t, "var toadd *TupleWithComplexItems", res)
	assertInCode(t, "for _, val := range stage1[lastIndex+1:]", res)
	assertInCode(t, "buf = bytes.NewBuffer(val)", res)
	assertInCode(t, "dec := json.NewDecoder(buf)", res)
	assertInCode(t, "dec.Decode(toadd)", res)
	assertInCode(t, "json.Marshal(data)", res)
	assertInCode(t, "for _, v := range m."+k+"Items", res)
}

func TestGenerateModel_WithTuple(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithTuple"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.NotEmpty(t, genModel.ExtraSchemas)
	require.NotEmpty(t, genModel.Properties)

	assert.False(t, genModel.IsTuple)
	assert.True(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsArray)
	assert.False(t, genModel.IsAnonymous)

	sch := genModel.ExtraSchemas[0]
	assert.True(t, sch.IsTuple)
	assert.False(t, sch.IsComplexObject)
	assert.False(t, sch.IsArray)
	assert.False(t, sch.IsAnonymous)
	assert.Equal(t, k+"FlagsTuple0", sch.Name)
	assert.False(t, sch.HasAdditionalItems)
	assert.Nil(t, sch.AdditionalItems)

	prop := genModel.Properties[0]
	assert.False(t, genModel.IsTuple)
	assert.True(t, genModel.IsComplexObject)
	assert.False(t, prop.IsArray)
	assert.False(t, prop.IsAnonymous)
	assert.Equal(t, k+"FlagsTuple0", prop.GoType)
	assert.Equal(t, "flags", prop.Name)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("with_tuple.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "swagger:model "+k+"Flags", res)
	assertInCode(t, "type "+k+"FlagsTuple0 struct {", res)
	assertInCode(t, "P0 *int64 `json:\"-\"`", res)
	assertInCode(t, "P1 *string `json:\"-\"`", res)
	assertInCode(t, k+"FlagsTuple0) UnmarshalJSON", res)
	assertInCode(t, k+"FlagsTuple0) MarshalJSON", res)
	assertInCode(t, "json.Marshal(data)", res)
	assert.NotRegexp(t, regexp.MustCompile("lastIndex"), res)

	for i, p := range sch.Properties {
		m := "m.P" + strconv.Itoa(i)
		r := "&dataP" + strconv.Itoa(i)
		var rr string
		if !p.IsNullable {
			rr = "dataP" + strconv.Itoa(i)
		} else {
			rr = r
		}
		assertInCode(t, fmt.Sprintf("buf = bytes.NewBuffer(stage1[%d])", i), res)
		assertInCode(t, "dec := json.NewDecoder(buf)", res)
		assertInCode(t, fmt.Sprintf("dec.Decode(%s)", r), res)
		assertInCode(t, "P"+strconv.Itoa(i)+",", res)
		assertInCode(t, fmt.Sprintf("%s = %s", m, rr), res)
	}
}

func TestGenerateModel_WithTupleWithExtra(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "WithTupleWithExtra"
	schema := definitions[k]
	opts := opts()
	tt := templateTest{t, opts.templates.MustGet("model")}
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.NotEmpty(t, genModel.ExtraSchemas)
	require.NotEmpty(t, genModel.Properties)

	assert.False(t, genModel.IsTuple)
	assert.True(t, genModel.IsComplexObject)
	assert.False(t, genModel.IsArray)
	assert.False(t, genModel.IsAnonymous)

	sch := genModel.ExtraSchemas[0]
	assert.True(t, sch.IsTuple)
	assert.False(t, sch.IsComplexObject)
	assert.False(t, sch.IsArray)
	assert.False(t, sch.IsAnonymous)
	assert.Equal(t, k+"FlagsTuple0", sch.Name)
	assert.True(t, sch.HasAdditionalItems)
	assert.NotEmpty(t, sch.AdditionalItems)

	prop := genModel.Properties[0]
	assert.False(t, genModel.IsTuple)
	assert.True(t, genModel.IsComplexObject)
	assert.False(t, prop.IsArray)
	assert.False(t, prop.IsAnonymous)
	assert.Equal(t, k+"FlagsTuple0", prop.GoType)
	assert.Equal(t, "flags", prop.Name)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, tt.template.Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("with_tuple.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "swagger:model "+k+"Flags", res)
	assertInCode(t, "type "+k+"FlagsTuple0 struct {", res)
	assertInCode(t, "P0 *int64 `json:\"-\"`", res)
	assertInCode(t, "P1 *string `json:\"-\"`", res)
	assertInCode(t, k+"FlagsTuple0Items []float32 `json:\"-\"`", res)
	assertInCode(t, k+"FlagsTuple0) UnmarshalJSON", res)
	assertInCode(t, k+"FlagsTuple0) MarshalJSON", res)
	assertInCode(t, "json.Marshal(data)", res)

	for i, p := range sch.Properties {
		m := "m.P" + strconv.Itoa(i)
		r := "&dataP" + strconv.Itoa(i)
		var rr string
		if !p.IsNullable {
			rr = "dataP" + strconv.Itoa(i)
		} else {
			rr = r
		}
		assertInCode(t, fmt.Sprintf("lastIndex = %d", i), res)
		assertInCode(t, fmt.Sprintf("buf = bytes.NewBuffer(stage1[%d])", i), res)
		assertInCode(t, "dec := json.NewDecoder(buf)", res)
		assertInCode(t, fmt.Sprintf("dec.Decode(%s)", r), res)
		assertInCode(t, "P"+strconv.Itoa(i)+",", res)
		assertInCode(t, fmt.Sprintf("%s = %s", m, rr), res)
	}

	assertInCode(t, "var lastIndex int", res)
	assertInCode(t, "var toadd float32", res)
	assertInCode(t, "for _, val := range stage1[lastIndex+1:]", res)
	assertInCode(t, "buf = bytes.NewBuffer(val)", res)
	assertInCode(t, "dec := json.NewDecoder(buf)", res)
	assertInCode(t, "dec.Decode(&toadd)", res)
	assertInCode(t, "json.Marshal(data)", res)
	assertInCode(t, "for _, v := range m."+k+"FlagsTuple0Items", res)
}

func TestGenerateModel_WithAllOfAndDiscriminator(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["Cat"]
	opts := opts()
	genModel, err := makeGenDefinition("Cat", "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Len(t, genModel.AllOf, 2)

	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "Cat", genModel.Name)
	assert.Equal(t, "Cat", genModel.GoType)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("cat.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "type Cat struct {", res)
	assertInCode(t, "Pet", res)
	assertInCode(t, "HuntingSkill *string `json:\"huntingSkill\"`", res)
}

func TestGenerateModel_WithAllOfAndDiscriminatorAndArrayOfPolymorphs(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["PetWithPets"]
	opts := opts()
	genModel, err := makeGenDefinition("PetWithPets", "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Len(t, genModel.AllOf, 2)

	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "PetWithPets", genModel.Name)
	assert.Equal(t, "PetWithPets", genModel.GoType)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("PetWithPets.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "type PetWithPets struct {", res)
	assertInCode(t, "UnmarshalPetSlice", res)
}

func TestGenerateModel_WithAllOf(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["WithAllOf"]
	opts := opts()
	genModel, err := makeGenDefinition("WithAllOf", "models", schema, specDoc, opts)
	require.NoError(t, err)

	assert.Len(t, genModel.AllOf, 7)
	assert.True(t, genModel.AllOf[1].HasAdditionalProperties)
	assert.True(t, genModel.IsComplexObject)
	assert.Equal(t, "WithAllOf", genModel.Name)
	assert.Equal(t, "WithAllOf", genModel.GoType)
	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("all_of_schema.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "type WithAllOf struct {", res)
	assertInCode(t, "type WithAllOfAO2P2 struct {", res)
	assertInCode(t, "type WithAllOfAO3P3 struct {", res)
	assertInCode(t, "type WithAllOfParamsAnon struct {", res)
	assertInCode(t, "type WithAllOfAO4Tuple4 struct {", res)
	assertInCode(t, "type WithAllOfAO5Tuple5 struct {", res)
	assertInCode(t, "Notable", res)
	assertInCode(t, "Title string `json:\"title,omitempty\"`", res)
	assertInCode(t, "Body string `json:\"body,omitempty\"`", res)
	assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
	assertInCode(t, "P0 *float32 `json:\"-\"`", res)
	assertInCode(t, "P0 *float64 `json:\"-\"`", res)
	assertInCode(t, "P1 *strfmt.DateTime `json:\"-\"`", res)
	assertInCode(t, "P1 *strfmt.Date `json:\"-\"`", res)
	assertInCode(t, "Opinion string `json:\"opinion,omitempty\"`", res)
	assertInCode(t, "WithAllOfAO5Tuple5Items []strfmt.Password `json:\"-\"`", res)
	assertInCode(t, "AO1 map[string]int32 `json:\"-\"`", res)
	assertInCode(t, "WithAllOfAO2P2 map[string]int64 `json:\"-\"`", res)
}

func findProperty(properties []GenSchema, name string) *GenSchema {
	for _, p := range properties {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

func getDefinitionProperty(genModel *GenDefinition, name string) *GenSchema {
	return findProperty(genModel.Properties, name)
}

func TestNumericKeys(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/162/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["AvatarUrls"]
	opts := opts()
	genModel, err := makeGenDefinition("AvatarUrls", "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("all_of_schema.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Nr16x16 string `json:\"16x16,omitempty\"`", res)
}

func TestGenModel_Issue196(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/196/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["Event"]
	opts := opts()
	genModel, err := makeGenDefinition("Event", "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("primitive_event.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Event) Validate(formats strfmt.Registry) error", res)
}

func TestGenModel_Issue222(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/tasklist.basic.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Price"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)
	require.True(t, genModel.HasValidations)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("price.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Price) Validate(formats strfmt.Registry) error", res)
	assertInCode(t, "Currency Currency `json:\"currency,omitempty\"`", res)
	assertInCode(t, "m.Currency.Validate(formats); err != nil", res)
}

func TestGenModel_Issue243(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "HasDynMeta"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("has_dyn_meta.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Metadata DynamicMetaData `json:\"metadata,omitempty\"`", res)
}

func TestGenModel_Issue252(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/252/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SodaBrand"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)
	require.False(t, genModel.IsNullable)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("soda_brand.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "type "+k+" string", res)
	assertInCode(t, "(m "+k+") validateSodaBrand", res)
	assertInCode(t, "(m "+k+") Validate", res)
}

func TestGenModel_Issue251(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/251/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "example"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("example.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)

	assertInCode(t, "type "+swag.ToGoName(k)+" struct", res)
	assertInCode(t, "Begin *strfmt.DateTime `json:\"begin\"`", res)
	assertInCode(t, "End strfmt.DateTime `json:\"end,omitempty\"`", res)
	assertInCode(t, "Name string `json:\"name,omitempty\"`", res)
	assertInCode(t, "(m *"+swag.ToGoName(k)+") validateBegin", res)
	assertInCode(t, "(m *"+swag.ToGoName(k)+") Validate", res)
}

func TestGenModel_Issue257(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "HasSpecialCharProp"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("example.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)

	assertInCode(t, "type "+swag.ToGoName(k)+" struct", res)
	assertInCode(t, "AtType string `json:\"@type,omitempty\"`", res)
	assertInCode(t, "Type string `json:\"type,omitempty\"`", res)
}

func TestGenModel_Issue340(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ImageTar"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("image_tar.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)

	assertInCode(t, "type "+swag.ToGoName(k)+" io.ReadCloser", res)
	assertNotInCode(t, "func (m ImageTar) Validate(formats strfmt.Registry) error", res)
}

func TestGenModel_Issue381(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "flags_list"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("flags_list.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertNotInCode(t, "m[i] != nil", res)
}

func TestGenModel_Issue300(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ActionItem"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("action_item.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Name ActionName `json:\"name\"`", res)
}

func TestGenModel_Issue398(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Property"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("action_item.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Computed bool `json:\"computed,omitempty\"`", res)
	assertInCode(t, "Intval *int64 `json:\"intval\"`", res)
	assertInCode(t, "PropType *string `json:\"propType\"`", res)
	assertInCode(t, "Strval *string `json:\"strval\"`", res)
}

func TestGenModel_Issue454(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/454/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["genericResource"]
	opts := opts()
	genModel, err := makeGenDefinition("genericResource", "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("generic_resource.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "rcv.Meta = stage1.Meta", res)
	assertInCode(t, "json.Marshal(stage1)", res)
	assertInCode(t, "stage1.Meta = m.Meta", res)
	assertInCode(t, "json.Marshal(m.GenericResource)", res)
}

func TestGenModel_Issue423(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/423/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["SRN"]
	opts := opts()
	genModel, err := makeGenDefinition("SRN", "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("SRN.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "propSite, err := UnmarshalSite(bytes.NewBuffer(data.Site), runtime.JSONConsumer())", res)
	assertInCode(t, "result.siteField = propSite", res)
}

func TestGenModel_Issue453(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/453/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "out_obj"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("out_obj.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `func (m *OutObj) validateFld3(formats strfmt.Registry)`, res)
}

func TestGenModel_Issue455(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/455/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "out_obj"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("out_obj.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `if err := validate.Required("fld2", "body", m.Fld2); err != nil {`, res)
}

func TestGenModel_Issue763(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/763/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "test_list"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("test_list.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "TheArray []*int32 `json:\"the_array\"`", res)
	assertInCode(t, `validate.MinimumInt("the_array"+"."+strconv.Itoa(i), "body", int64(*m.TheArray[i]), 0, false)`, res)
	assertInCode(t, `validate.MaximumInt("the_array"+"."+strconv.Itoa(i), "body", int64(*m.TheArray[i]), 10, false)`, res)
}

func TestGenModel_Issue811_NullType(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/811/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "teamRepos"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("team_repos.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Language interface{} `json:\"language,omitempty\"`", res)
}

func TestGenModel_Issue811_Emojis(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/811/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "emojis"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("team_repos.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Plus1 string `json:\"+1,omitempty\"`", res)
	assertInCode(t, "Minus1 string `json:\"-1,omitempty\"`", res)
}

func TestGenModel_Issue752_EOFErr(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/azure-text-analyis.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "OperationResult"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("out_obj.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `&& err != io.EOF`, res)
}

func TestImports_ExistingModel(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/existing-model.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	opts := opts()

	k := "JsonWebKeySet"
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	require.NotNil(t, genModel)
	require.NotNil(t, genModel.Imports)
	assert.Equal(t, "github.com/user/package", genModel.Imports["jwk"])

	k = "JsonWebKey"
	genModel, err = makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	require.NotNil(t, genModel)
	require.NotNil(t, genModel.Imports)
	assert.Equal(t, "github.com/user/package", genModel.Imports["jwk"])
}

func TestGenModel_Issue786(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/786/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "MyFirstObject"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	require.False(t, genModel.Properties[0].AdditionalProperties.IsNullable)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("MyFirstObject.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `m.validateEntreeChoiceValueEnum("entree_choice"+"."+k, "body", m.EntreeChoice[k])`, res)
}

func TestGenModel_Issue822(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/822/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Pet"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	ap := genModel.AdditionalProperties
	require.True(t, genModel.HasAdditionalProperties)
	require.NotNil(t, ap)
	require.False(t, ap.IsNullable)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("pet.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `PetAdditionalProperties map[string]interface{}`, res)
	assertInCode(t, `m.PetAdditionalProperties = result`, res)
	assertInCode(t, `additional, err := json.Marshal(m.PetAdditionalProperties)`, res)
}

func TestGenModel_Issue981(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/981/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "User"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("user.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "FirstName string `json:\"first_name,omitempty\"`", res)
	assertInCode(t, "LastName string `json:\"last_name,omitempty\"`", res)
	assertInCode(t, "if swag.IsZero(m.Type)", res)
	assertInCode(t, `validate.MinimumInt("user_type", "body", m.Type, 1, false)`, res)
	assertInCode(t, `validate.MaximumInt("user_type", "body", m.Type, 5, false)`, res)
}

func TestGenModel_Issue1341(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1341/swagger.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ExecutableValueString"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("executable_value_string.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, `return errors.New(422, "invalid ValueType value: %q", base.ValueType`, res)
	assertInCode(t, "result.testField = base.Test", res)
	assertInCode(t, "Test *string `json:\"Test\"`", res)
	assertInCode(t, "Test: m.Test(),", res)
}

// This tests to check that format validation is performed on non required schema properties
func TestGenModel_Issue1347(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1347/fixture-1347.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	schema := definitions["ContainerConfig"]
	opts := opts()
	genModel, err := makeGenDefinition("ContainerConfig", "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("Foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	// Just verify that the validation call is generated even though we add a non-required property
	assertInCode(t, `validate.FormatOf("config1", "body", "hostname", m.Config1.String(), formats)`, res)
}

// This tests to check that format validation is performed on MAC format
func TestGenModel_Issue1348(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1348/fixture-1348.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ContainerConfig"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	// Just verify that the validation call is generated with proper format
	assertInCode(t, `if err := validate.FormatOf("config1", "body", "mac", m.Config1.String(), formats)`, res)
}

// This tests that additionalProperties with validation is generated properly.
func TestGenModel_Issue1198(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1198/fixture-1198.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "pet"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	// Just verify that the validation call is generated with proper format
	assertInCode(t, `if err := m.validateDate(formats); err != nil {`, res)
}

// This tests that additionalProperties with validation is generated properly.
func TestGenModel_Issue1397a(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1397/fixture-1397a.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ContainerConfig"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	// Just verify that the validation call is generated with proper format
	assertInCode(t, `if swag.IsZero(m[k]) { // not required`, res)
}

// This tests that an enum of object values validates properly.
func TestGenModel_Issue1397b(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1397/fixture-1397b.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ContainerConfig"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	// Just verify that the validation call is generated with proper format
	assertInCode(t, `if err := m.validateContainerConfigEnum("", "body", m); err != nil {`, res)
}

// This tests that additionalProperties with an array of polymorphic objects is generated properly.
func TestGenModel_Issue1409(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/1409/fixture-1409.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Graph"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	// Just verify that the validation call is generated with proper format
	assertInCode(t, `nodes, err := UnmarshalNodeSlice(bytes.NewBuffer(data.Nodes), runtime.JSONConsumer())`, res)
	assertInCode(t, `if err := json.Unmarshal(raw, &rawProps); err != nil {`, res)
	assertInCode(t, `m.GraphAdditionalProperties[k] = toadd`, res)
	assertInCode(t, `b3, err = json.Marshal(m.GraphAdditionalProperties)`, res)
}

// This tests makes sure model definitions from inline schema in response are properly flattened and get validation
func TestGenModel_Issue866(t *testing.T) {
	defer discardOutput()()

	specDoc, err := loads.Spec("../fixtures/bugs/866/fixture-866.yaml")
	require.NoError(t, err)

	p, ok := specDoc.Spec().Paths.Paths["/"]
	require.True(t, ok)

	op := p.Get
	responses := op.Responses.StatusCodeResponses
	for k, r := range responses {
		t.Logf("Response: %d", k)
		schema := *r.Schema
		opts := opts()
		genModel, err := makeGenDefinition("GetOKBody", "models", schema, specDoc, opts)
		require.NoError(t, err)

		buf := bytes.NewBuffer(nil)
		require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

		ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
		require.NoError(t, err)

		res := string(ct)
		assertInCode(t, `if err := validate.Required(`, res)
		assertInCode(t, `if err := validate.MaxLength(`, res)
		assertInCode(t, `if err := m.validateAccessToken(formats); err != nil {`, res)
		assertInCode(t, `if err := m.validateAccountID(formats); err != nil {`, res)
	}
}

// This tests makes sure marshalling and validation is generated in aliased formatted definitions
func TestGenModel_Issue946(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/946/fixture-946.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "mydate"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, `type Mydate strfmt.Date`, res)
	assertInCode(t, `func (m *Mydate) UnmarshalJSON(b []byte) error {`, res)
	assertInCode(t, `return ((*strfmt.Date)(m)).UnmarshalJSON(b)`, res)
	assertInCode(t, `func (m Mydate) MarshalJSON() ([]byte, error) {`, res)
	assertInCode(t, `return (strfmt.Date(m)).MarshalJSON()`, res)
	assertInCode(t, `if err := validate.FormatOf("", "body", "date", strfmt.Date(m).String(), formats); err != nil {`, res)
}

// This tests makes sure that docstring in inline schema in response properly reflect the Required property
func TestGenModel_Issue910(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/910/fixture-910.yaml")
	require.NoError(t, err)

	p, ok := specDoc.Spec().Paths.Paths["/mytest"]
	require.True(t, ok)

	op := p.Get
	responses := op.Responses.StatusCodeResponses
	for k, r := range responses {
		t.Logf("Response: %d", k)
		schema := *r.Schema
		opts := opts()
		genModel, err := makeGenDefinition("GetMyTestOKBody", "models", schema, specDoc, opts)
		require.NoError(t, err)

		buf := bytes.NewBuffer(nil)
		require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

		ct, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
		require.NoError(t, err)

		res := string(ct)
		assertInCode(t, "// bar\n	// Required: true\n	Bar *int64 `json:\"bar\"`", res)
		assertInCode(t, "// foo\n	// Required: true\n	Foo interface{} `json:\"foo\"`", res)
		assertInCode(t, "// baz\n	Baz int64 `json:\"baz,omitempty\"`", res)
		assertInCode(t, "// quux\n	Quux []string `json:\"quux\"`", res)
		assertInCode(t, `if err := validate.Required("bar", "body", m.Bar); err != nil {`, res)
		assertInCode(t, `if m.Foo == nil {`, res) // interface{} now checked against nil (validate.Required fails on any zero value)
		assertNotInCode(t, `if err := validate.Required("baz", "body", m.Baz); err != nil {`, res)
		assertNotInCode(t, `if err := validate.Required("quux", "body", m.Quux); err != nil {`, res)
		// NOTE(fredbi); fixed Required in slices. This property has actually no validation
	}
}

func TestGenerateModel_Xorder(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/x-order.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "sessionData"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("x-order.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	// if no x-order then alphabetical order, like DeviceID, SessionID, UMain.
	// There is x-order
	//   sessionId
	//     x-order: 0
	//   deviceId:
	//     x-order: 1
	//   uMain:
	//     x-order: 2
	// This is need for msgpack-array.
	foundDeviceID := strings.Index(res, "DeviceID")
	foundSessionID := strings.Index(res, "SessionID")
	foundUMain := strings.Index(res, "UMain")
	foundAaa := strings.Index(res, "Aaa")
	foundBbb := strings.Index(res, "Bbb")
	foundZzz := strings.Index(res, "Zzz")
	assert.True(t, foundSessionID < foundDeviceID)
	assert.True(t, foundSessionID < foundUMain)
	assert.True(t, foundUMain < foundAaa)
	assert.True(t, foundAaa < foundBbb)
	assert.True(t, foundBbb < foundZzz)
}

func TestGenModel_Issue1623(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/enhancements/1623/swagger.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Foo"
	schema := definitions[k]
	opts := opts()

	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("Foo.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	assertInCode(t, "ArrayHasOmitEmptyFalse []string `json:\"arrayHasOmitEmptyFalse\"`", res)
	assertInCode(t, "ArrayHasOmitEmptyTrue []string `json:\"arrayHasOmitEmptyTrue,omitempty\"`", res)
	assertInCode(t, "ArrayNoOmitEmpty []string `json:\"arrayNoOmitEmpty\"`", res)
	assertInCode(t, "GeneralHasOmitEmptyFalse string `json:\"generalHasOmitEmptyFalse\"`", res)
	assertInCode(t, "GeneralHasOmitEmptyTrue string `json:\"generalHasOmitEmptyTrue,omitempty\"`", res)
	assertInCode(t, "GeneralNoOmitEmpty string `json:\"generalNoOmitEmpty,omitempty\"`", res)
	assertInCode(t, "RefHasOmitEmptyFalse Bar `json:\"refHasOmitEmptyFalse,omitempty\"`", res)
	assertInCode(t, "RefHasOmitEmptyTrue Bar `json:\"refHasOmitEmptyTrue,omitempty\"`", res)
	assertInCode(t, "RefNoOmitEmpty Bar `json:\"refNoOmitEmpty,omitempty\"`", res)
	assertInCode(t, "IntHasJSONString int64 `json:\"intHasJsonString,omitempty,string\"`", res)
	assertInCode(t, "BoolHasJSONString bool `json:\"boolHasJsonString,omitempty,string\"`", res)

}

func TestGenerateModel_Issue2457(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/bugs/2457/fixture-2457.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ObjWithCustomTag"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	assertInCode(t, "myobjtag:\"foo,bar\"", buf.String())
}

func TestGenModel_Pr2464(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/enhancements/2464/fixture-2464.yaml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions

	tests := map[string]struct {
		model    string
		expected []string
	}{
		"withDot": {
			model: "band",
			expected: []string{
				`BandNr2Dot4Ghz Band = "2.4Ghz"`,
				`BandNr24Ghz Band = "24Ghz"`,
			},
		},
		"withSigns": {
			model: "change",
			expected: []string{
				`ChangeDash1 Change = "-1"`,
				`ChangeNr0 Change = "0"`,
				`ChangePlus1 Change = "+1"`,
			},
		},
		"hashtag": {
			model: "topic",
			expected: []string{
				`TopicHashtagOne Topic = "#one"`,
				`TopicTwoHashtagTwo Topic = "two#two"`,
			},
		},
	}

	for name, spec := range tests {
		t.Run(name, func(t *testing.T) {
			k := spec.model
			schema := definitions[k]
			opts := opts()
			genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
			require.NoError(t, err)

			buf := bytes.NewBuffer(nil)
			require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

			for _, expected := range spec.expected {
				assertInCode(t, expected, buf.String())
			}
		})
	}
}

func TestGenModel_KeepSpecPropertiesOrder(t *testing.T) {
	ymlFile := "../fixtures/codegen/keep-spec-order.yml"
	opts := opts()
	abcType := "abctype"

	specDoc, err := loads.Spec(ymlFile)
	require.NoError(t, err)

	orderedSpecDoc, err := loads.Spec(WithAutoXOrder(ymlFile))
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	orderedDefinitions := orderedSpecDoc.Spec().Definitions

	genModel, err := makeGenDefinition(abcType, "models", definitions[abcType], specDoc, opts)
	assert.NoError(t, err)
	orderGenModel, err := makeGenDefinition(abcType, "models", orderedDefinitions[abcType], orderedSpecDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	orderBuf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(orderBuf, orderGenModel))

	ff, err := opts.LanguageOpts.FormatContent("keepSpecOrder.go", buf.Bytes())
	require.NoError(t, err)

	modelCode := string(ff)
	ff, err = opts.LanguageOpts.FormatContent("keepSpecOrder-ordered.go", orderBuf.Bytes())
	require.NoError(t, err)

	orderModelCode := string(ff)

	// without auto order, properties sorted by alphanumeric
	foundA := strings.Index(modelCode, "Aaa")
	foundB := strings.Index(modelCode, "Bbb")
	foundC := strings.Index(modelCode, "Ccc")
	assert.True(t, foundA < foundB)
	assert.True(t, foundB < foundC)

	foundOrderA := strings.Index(orderModelCode, "Aaa")
	foundOrderB := strings.Index(orderModelCode, "Bbb")
	foundOrderC := strings.Index(orderModelCode, "Ccc")

	assert.True(t, foundOrderC < foundOrderB)
	assert.True(t, foundOrderB < foundOrderA)

	foundInnerA := strings.Index(modelCode, "InnerAaa")
	foundInnerB := strings.Index(modelCode, "InnerBbb")
	foundInnerC := strings.Index(modelCode, "InnerCcc")
	assert.True(t, foundInnerA < foundInnerB)
	assert.True(t, foundInnerB < foundInnerC)

	foundOrderInnerA := strings.Index(orderModelCode, "InnerAaa")
	foundOrderInnerB := strings.Index(orderModelCode, "InnerBbb")
	foundOrderInnerC := strings.Index(orderModelCode, "InnerCcc")

	assert.True(t, foundOrderInnerC < foundOrderInnerB)
	assert.True(t, foundOrderInnerB < foundOrderInnerA)
}

func TestGenModel_StrictAdditionalProperties(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/strict-additional-properties.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Body"
	schema := definitions[k]
	opts := opts()

	opts.StrictAdditionalProperties = true

	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ff, err := opts.LanguageOpts.FormatContent("strictAdditionalProperties.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ff)
	for _, tt := range []struct {
		name      string
		assertion func(testing.TB, string, string) bool
	}{
		{k, assertInCode},
		{k + "Explicit", assertInCode},
		{k + "Implicit", assertInCode},
		{k + "Disabled", assertNotInCode},
	} {
		fn := funcBody(res, "*"+tt.name+") UnmarshalJSON(data []byte) error")
		if assert.NotEmpty(t, fn, "Method UnmarshalJSON should be defined for type *"+tt.name) {
			tt.assertion(t, "dec.DisallowUnknownFields()", fn)
		}
	}
}

func TestGenModel_XMLStructTags_WithXML(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/xml-model.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "XmlWithAttribute"
	opts := opts()
	opts.WithXML = true

	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("xml_with_attribute.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "Author *string `json:\"author\" xml:\"author\"`", res)
	assertInCode(t, "Children []*XMLChild `json:\"children\" xml:\"children\"`", res)
	assertInCode(t, "ID int64 `json:\"id,omitempty\" xml:\"id,attr,omitempty\"`", res)
	assertInCode(t, "IsPublished *bool `json:\"isPublished\" xml:\"published,attr\"`", res)
	assertInCode(t, "SingleChild *XMLChild `json:\"singleChild,omitempty\" xml:\"singleChild,omitempty\"`", res)
	assertInCode(t, "Title string `json:\"title,omitempty\" xml:\"xml-title,omitempty\"`", res)
}

func TestGenModel_XMLStructTags_Explicit(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/xml-model.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "XmlWithAttribute"
	opts := opts()

	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("xml_with_attribute.go", buf.Bytes())
	require.NoErrorf(t, err, "format error: %v\n%s", err, buf.String())

	res := string(ct)
	assertInCode(t, "Author *string `json:\"author\"`", res)
	assertInCode(t, "Children []*XMLChild `json:\"children\"`", res)
	assertInCode(t, "ID int64 `json:\"id,omitempty\" xml:\"id,attr,omitempty\"`", res)
	assertInCode(t, "IsPublished *bool `json:\"isPublished\" xml:\"published,attr\"`", res)
	assertInCode(t, "SingleChild *XMLChild `json:\"singleChild,omitempty\"`", res)
	assertInCode(t, "Title string `json:\"title,omitempty\" xml:\"xml-title,omitempty\"`", res)
}

func TestGenerateModels(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	cwd := testCwd(t)
	const root = "generated_models"
	defer func() {
		_ = os.RemoveAll(filepath.Join(cwd, root))
	}()

	t.Run("generate models", func(t *testing.T) {
		cases := map[string]generateFixture{
			"allDefinitions": {
				spec:   "../fixtures/bugs/1042/fixture-1042.yaml",
				target: "../fixtures/bugs/1042",
				verify: func(t testing.TB, target string) {
					target = filepath.Join(target, defaultModelsTarget)
					require.True(t, fileExists(target, ""))
					assert.True(t, fileExists(target, "a.go"))
					assert.True(t, fileExists(target, "b.go"))
				},
			},
			"acceptDefinitions": {
				spec:   "../fixtures/enhancements/2333/fixture-definitions.yaml",
				target: "../fixtures/enhancements/2333",
				prepare: func(_ testing.TB, opts *GenOpts) {
					opts.AcceptDefinitionsOnly = true
				},
				verify: func(t testing.TB, target string) {
					target = filepath.Join(target, defaultModelsTarget)
					require.True(t, fileExists(target, ""))
					assert.True(t, fileExists(target, "model_interface.go"))
					assert.True(t, fileExists(target, "records_model.go"))
					assert.True(t, fileExists(target, "records_model_with_max.go"))
					assert.False(t, fileExists(target, "restapi"))
				},
			},
		}
		for k, cas := range cases {
			name := k
			thisCas := cas

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				defer thisCas.warnFailed(t)

				opts := testGenOpts()
				defer thisCas.prepareTarget(t, name, "model_test", root, opts)()

				if thisCas.prepare != nil {
					thisCas.prepare(t, opts)
				}

				t.Logf("generating test models at: %s", opts.Target)

				err := GenerateModels([]string{"", ""}, opts) // NOTE: generate all models, ignore ""
				if thisCas.wantError {
					require.Errorf(t, err, "expected an error for models build fixture: %s", opts.Spec)
				} else {
					require.NoError(t, err, "unexpected error for models build fixture: %s", opts.Spec)
				}

				if thisCas.verify != nil {
					thisCas.verify(t, opts.Target)
				}
			})
		}
	})
}

func Test_Issue2482(t *testing.T) {
	// generation of a pointer converter for enum vars
	specDoc, err := loads.Spec("../fixtures/bugs/252/swagger.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SodaBrand"
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", definitions[k], specDoc, opts)
	require.NoError(t, err)
	require.False(t, genModel.IsNullable)

	buf := bytes.NewBuffer(nil)
	require.NoError(t, opts.templates.MustGet("model").Execute(buf, genModel))

	ct, err := opts.LanguageOpts.FormatContent("soda_brand.go", buf.Bytes())
	require.NoError(t, err)

	res := string(ct)
	assertInCode(t, "func NewSodaBrand(value SodaBrand) *SodaBrand {", res)
	assertInCode(t, "v := value", res)
	assertInCode(t, "return &v", res)
}
