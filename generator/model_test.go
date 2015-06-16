package generator

import (
	"bytes"
	"fmt"
	"strings"
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

func TestGenerateModel_Sanity(t *testing.T) {
	// just checks if it can render and format these things
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions

		//k := "Comment"
		//schema := definitions[k]
		for k, schema := range definitions {
			genModel, err := makeGenDefinition(k, "models", schema, specDoc)

			if assert.NoError(t, err) {
				//b, _ := json.MarshalIndent(genModel.IsCo, "", "  ")
				//fmt.Println(string(b))
				rendered := bytes.NewBuffer(nil)

				err := modelTemplate.Execute(rendered, genModel)
				if assert.NoError(t, err) {
					if assert.NoError(t, err) {
						formatted, err := formatGoFile(strings.ToLower(k)+".go", rendered.Bytes())
						if assert.NoError(t, err) {
							fmt.Println(string(formatted))
						} else {
							fmt.Println(rendered.String())
							break
						}

						//assert.EqualValues(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(formatted)))
					}
				}
			}
		}
	}
}

func TestGenerateModel_DocString(t *testing.T) {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	tt := templateTest{t, templ}

	var gmp GenSchema
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
	expected = `TheModel the model
`
	tt.assertRender(gmp, expected)
}

func TestGenerateModel_PropertyValidation(t *testing.T) {
	templ := template.Must(template.New("propertyValidationDocString").Funcs(FuncMap).Parse(string(assetStuctFieldValidation)))
	tt := templateTest{t, templ}

	var gmp GenSchema
	gmp.Required = true
	tt.assertRender(gmp, `
Required: true
`)
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

	var gmp GenSchema
	gmp.Name = "some name"
	gmp.resolvedType = resolvedType{GoType: "string", IsPrimitive: true}
	gmp.Title = "The title of the property"

	tt.assertRender(gmp, `/* The title of the property
 */
`+"SomeName string `json:\"some name\"`\n")

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
	tt.assertRender(gmp, `/* The title of the property

The description of the property

Required: true
Read Only: true
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
	{GenSchema{resolvedType: resolvedType{GoType: "interface{}", IsInterface: true}}, "interface{}"},
	{GenSchema{resolvedType: resolvedType{GoType: "[]int32", IsArray: true}}, "[]int32"},
	{GenSchema{resolvedType: resolvedType{GoType: "[]string", IsArray: true}}, "[]string"},
	{GenSchema{resolvedType: resolvedType{GoType: "models.Task", IsComplexObject: true, IsNullable: true, IsAnonymous: false}}, "*models.Task"},
}

func TestGenSchemaType(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schemaType")}
	for _, v := range schTypeGenDataSimple {
		tt.assertRender(v.Value, v.Expected)
	}
}
