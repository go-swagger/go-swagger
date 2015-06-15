package generator

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

type templateTest struct {
	t        testing.TB
	template *template.Template
}

func (tt *templateTest) assertRender(data interface{}, expected string) bool {
	buf := bytes.NewBuffer(nil)
	err := tt.template.Execute(buf, data)
	if !assert.NoError(tt.t, err) {
		return false
	}
	return assert.Equal(tt.t, expected, buf.String())
}

func TestGenerateModel_Primitives(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions

		schema := definitions["AllPrimitives"]
		genModel, err := makeCodegenModel("AllPrimitives", "models", schema, specDoc)
		if assert.NoError(t, err) {
			rendered := bytes.NewBuffer(nil)

			err := modelTemplate.Execute(rendered, genModel)
			if assert.NoError(t, err) {
				//fmt.Println(rendered.String())
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("all_primitives.go", rendered.Bytes())
					if assert.NoError(t, err) {
						fmt.Println(string(formatted))
					}
					//assert.EqualValues(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(formatted)))
				}
			}
		}
	}
}

func TestGenerateModel_DocString(t *testing.T) {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	tt := templateTest{t, templ}

	var gmp genModelProperty
	gmp.Title = "The title of the property"
	gmp.Description = "The description of the property"
	var expected = `The title of the property

The description of the property
`
	tt.assertRender(gmp, expected)

	gmp.Title = ""
	expected = `The description of the property
`
	tt.assertRender(gmp, expected)

	gmp.Description = ""
	gmp.Name = "theModel"
	expected = `TheModel the model`
	tt.assertRender(gmp, expected)
}

func TestGenerateModel_PropertyValidation(t *testing.T) {
	templ := template.Must(template.New("propertyValidationDocString").Funcs(FuncMap).Parse(string(assetStuctFieldValidation)))
	tt := templateTest{t, templ}

	var gmp genModelProperty
	gmp.Required = true
	tt.assertRender(gmp, `
Required: true
`)
	gmp.Maximum = 10
	gmp.ExclusiveMaximum = true
	gmp.Minimum = 10
	gmp.ExclusiveMinimum = true
	gmp.MaxLength = 20
	gmp.MinLength = 20
	gmp.Pattern = "\\w[\\w- ]+"
	gmp.MaxItems = 30
	gmp.MinItems = 30
	gmp.UniqueItems = true

	tt.assertRender(gmp, `
Required: true
Maximum: < 10
Minimum: > 10
Max Length: 20
Min Length: 20
Pattern: \w[\w- ]+
Max Items: 30
Min Items: 30
Unique: true
`)

	gmp.Required = false
	gmp.ExclusiveMaximum = false
	gmp.ExclusiveMinimum = false
	tt.assertRender(gmp, `
Maximum: 10
Minimum: 10
Max Length: 20
Min Length: 20
Pattern: \w[\w- ]+
Max Items: 30
Min Items: 30
Unique: true
`)

}

func TestGenerateModel_SchemaField(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("structfield")}

	var gmp genModelProperty
	gmp.Name = "some name"
	gmp.T = resolvedType{GoType: "string", IsPrimitive: true}
	gmp.Title = "The title of the property"

	tt.assertRender(gmp, `/* The title of the property */
`+"SomeName string `json:\"some name\"`\n")

	gmp.Description = "The description of the property"
	gmp.Required = true
	gmp.Maximum = 10
	gmp.ExclusiveMaximum = true
	gmp.Minimum = 10
	gmp.ExclusiveMinimum = true
	gmp.MaxLength = 20
	gmp.MinLength = 20
	gmp.Pattern = "\\w[\\w- ]+"
	gmp.MaxItems = 30
	gmp.MinItems = 30
	gmp.UniqueItems = true
	tt.assertRender(gmp, `/* The title of the property

The description of the property

Required: true
Maximum: < 10
Minimum: > 10
Max Length: 20
Min Length: 20
Pattern: \w[\w- ]+
Max Items: 30
Min Items: 30
Unique: true
 */
`+"SomeName string `json:\"some name\"`\n")
}

// TODO:
// * Tuples, Tuples with AdditionalItems
// * Slices with additional items
// * Embedded Structs
// * Schemas for simple types

var schTypeGenDataSimple = []struct {
	Value    resolvedType
	Expected string
}{
	{resolvedType{GoType: "string", IsPrimitive: true}, "string"},
	{resolvedType{GoType: "string", IsPrimitive: true, IsNullable: true}, "*string"},
	{resolvedType{GoType: "bool", IsPrimitive: true}, "bool"},
	{resolvedType{GoType: "int32", IsPrimitive: true}, "int32"},
	{resolvedType{GoType: "int64", IsPrimitive: true}, "int64"},
	{resolvedType{GoType: "float32", IsPrimitive: true}, "float32"},
	{resolvedType{GoType: "float64", IsPrimitive: true}, "float64"},
	{resolvedType{GoType: "strfmt.Base64", IsPrimitive: true}, "strfmt.Base64"},
	{resolvedType{GoType: "strfmt.Date", IsPrimitive: true}, "strfmt.Date"},
	{resolvedType{GoType: "strfmt.DateTime", IsPrimitive: true}, "strfmt.DateTime"},
	{resolvedType{GoType: "strfmt.URI", IsPrimitive: true}, "strfmt.URI"},
	{resolvedType{GoType: "strfmt.Email", IsPrimitive: true}, "strfmt.Email"},
	{resolvedType{GoType: "strfmt.Hostname", IsPrimitive: true}, "strfmt.Hostname"},
	{resolvedType{GoType: "strfmt.IPv4", IsPrimitive: true}, "strfmt.IPv4"},
	{resolvedType{GoType: "strfmt.IPv6", IsPrimitive: true}, "strfmt.IPv6"},
	{resolvedType{GoType: "strfmt.UUID", IsPrimitive: true}, "strfmt.UUID"},
	{resolvedType{GoType: "strfmt.UUID3", IsPrimitive: true}, "strfmt.UUID3"},
	{resolvedType{GoType: "strfmt.UUID4", IsPrimitive: true}, "strfmt.UUID4"},
	{resolvedType{GoType: "strfmt.UUID5", IsPrimitive: true}, "strfmt.UUID5"},
	{resolvedType{GoType: "strfmt.ISBN", IsPrimitive: true}, "strfmt.ISBN"},
	{resolvedType{GoType: "strfmt.ISBN10", IsPrimitive: true}, "strfmt.ISBN10"},
	{resolvedType{GoType: "strfmt.ISBN13", IsPrimitive: true}, "strfmt.ISBN13"},
	{resolvedType{GoType: "strfmt.CreditCard", IsPrimitive: true}, "strfmt.CreditCard"},
	{resolvedType{GoType: "strfmt.SSN", IsPrimitive: true}, "strfmt.SSN"},
	{resolvedType{GoType: "strfmt.HexColor", IsPrimitive: true}, "strfmt.HexColor"},
	{resolvedType{GoType: "strfmt.RGBColor", IsPrimitive: true}, "strfmt.RGBColor"},
	{resolvedType{GoType: "strfmt.Duration", IsPrimitive: true}, "strfmt.Duration"},
	{resolvedType{GoType: "strfmt.Password", IsPrimitive: true}, "strfmt.Password"},
	{resolvedType{GoType: "interface{}", IsInterface: true}, "interface{}"},
	{resolvedType{GoType: "[]int32", IsArray: true}, "[]int32"},
	{resolvedType{GoType: "[]string", IsArray: true}, "[]string"},
	{resolvedType{GoType: "models.Task", IsComplexObject: true, IsNullable: true, IsAnonymous: false}, "*models.Task"},
}

func TestGenSchemaType(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schemaType")}
	for _, v := range schTypeGenDataSimple {
		tt.assertRender(v.Value, v.Expected)
	}
}
