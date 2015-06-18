package generator

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
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

		k := "Comment"
		schema := definitions[k]
		//for k, schema := range definitions {
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)

		if assert.NoError(t, err) {
			//b, _ := json.MarshalIndent(genModel, "", "  ")
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
						//break
					}

					//assert.EqualValues(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(formatted)))
				}
			}
			//}
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
	{GenSchema{resolvedType: resolvedType{GoType: "map[string]int32", IsMap: true}}, "map[string]int32"},
	{GenSchema{resolvedType: resolvedType{GoType: "models.Task", IsComplexObject: true, IsNullable: true, IsAnonymous: false}}, "*models.Task"},
}

func TestGenSchemaType(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schemaType")}
	for _, v := range schTypeGenDataSimple {
		tt.assertRender(v.Value, v.Expected)
	}
}
func TestGenerateModel_Primitives(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	for _, v := range schTypeGenDataSimple {
		val := v.Value
		if val.IsComplexObject {
			continue
		}
		val.Name = "theType"
		exp := v.Expected
		if val.IsNullable {
			exp = exp[1:]
		}
		tt.assertRender(val, "type TheType "+exp+"\n\n")
	}
}

func TestGenerateModel_Nota(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "Nota"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("type Nota map[string]int32")), res)
			}
		}
	}
}

func TestGenerateModel_NotaWithName(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "NotaWithName"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsAdditionalProperties)
			assert.False(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsMap)
			assert.False(t, genModel.IsAnonymous)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile("type "+k+" struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("AdditionalProperties map\\[string\\]int32 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("Name string `json:\"name\"`"), res)
				assert.Regexp(t, regexp.MustCompile(k+"\\) UnmarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(k+"\\) MarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(m)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(m.AdditionalProperties)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(data, &stage1)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(data, &stage2)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(v, &toadd)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("result[k] = toadd")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("m.AdditionalProperties = result")), res)
				for _, p := range genModel.Properties {
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("delete(stage2, \""+p.Name+"\")")), res)
				}

			}
		}
	}
}

func TestGenerateModel_MapRef(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithMap"]
		genModel, err := makeGenDefinition("WithMap", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.False(t, genModel.HasAdditionalProperties)
			prop := getDefinitionProperty(genModel, "data")
			assert.True(t, prop.HasAdditionalProperties)
			assert.True(t, prop.IsMap)
			assert.False(t, prop.IsComplexObject)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("type WithMap struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("Data map\\[string\\]string `json:\"data\"`"), res)
		}
	}
}

func TestGenerateModel_WithAdditional(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "WithAdditional"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.ExtraSchemas) {
			assert.False(t, genModel.HasAdditionalProperties)
			assert.False(t, genModel.IsMap)
			assert.False(t, genModel.IsAdditionalProperties)
			assert.True(t, genModel.IsComplexObject)

			sch := genModel.ExtraSchemas[0]
			assert.True(t, sch.HasAdditionalProperties)
			assert.False(t, sch.IsMap)
			assert.True(t, sch.IsAdditionalProperties)
			assert.False(t, sch.IsComplexObject)

			if assert.NotNil(t, sch.AdditionalProperties) {
				prop := findProperty(genModel.Properties, "data")
				assert.False(t, prop.HasAdditionalProperties)
				assert.False(t, prop.IsMap)
				assert.False(t, prop.IsAdditionalProperties)
				assert.True(t, prop.IsComplexObject)
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, genModel)
				if assert.NoError(t, err) {
					res := buf.String()
					assert.Regexp(t, regexp.MustCompile("type "+k+" struct\\s*{"), res)
					assert.Regexp(t, regexp.MustCompile("Data "+k+"DataAddedProps0 `json:\"data\"`"), res)
					assert.Regexp(t, regexp.MustCompile("type "+k+"DataAddedProps0 struct\\s*{"), res)
					assert.Regexp(t, regexp.MustCompile("AdditionalProperties map\\[string\\]string `json:\"-\"`"), res)
					assert.Regexp(t, regexp.MustCompile("Name string `json:\"name\"`"), res)
					assert.Regexp(t, regexp.MustCompile(k+"DataAddedProps0\\) UnmarshalJSON"), res)
					assert.Regexp(t, regexp.MustCompile(k+"DataAddedProps0\\) MarshalJSON"), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(m)")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(m.AdditionalProperties)")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(data, &stage1)")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(data, &stage2)")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(v, &toadd)")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("result[k] = toadd")), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("m.AdditionalProperties = result")), res)
					for _, p := range sch.Properties {
						assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("delete(stage2, \""+p.Name+"\")")), res)
					}
				}
			}
		}
	}
}

func TestGenerateModel_JustRef(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["JustRef"]
		genModel, err := makeGenDefinition("JustRef", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, genModel.AllOf)
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "JustRef", genModel.Name)
			assert.Equal(t, "JustRef", genModel.GoType)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("type JustRef struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("Notable"), res)
		}
	}
}

func TestGenerateModel_WithRef(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithRef"]
		genModel, err := makeGenDefinition("WithRef", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithRef", genModel.Name)
			assert.Equal(t, "WithRef", genModel.GoType)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("type WithRef struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("Notes Notable `json:\"notes\"`"), res)
		}
	}
}

func TestGenerateModel_WithNullableRef(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithNullableRef"]
		genModel, err := makeGenDefinition("WithNullableRef", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithNullableRef", genModel.Name)
			assert.Equal(t, "WithNullableRef", genModel.GoType)
			prop := getDefinitionProperty(genModel, "notes")
			assert.True(t, prop.IsNullable)
			assert.True(t, prop.IsComplexObject)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("type WithNullableRef struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("Notes \\*Notable `json:\"notes\"`"), res)
		}
	}
}

func TestGenerateModel_WithItems(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithItems"]
		genModel, err := makeGenDefinition("WithItems", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.Empty(t, genModel.Items)
			assert.True(t, genModel.IsComplexObject)
			prop := getDefinitionProperty(genModel, "tags")
			assert.NotEmpty(t, prop.Items)
			assert.True(t, prop.IsArray)
			assert.False(t, prop.IsComplexObject)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("type WithItems struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("Tags \\[\\]string `json:\"tags\"`"), res)
		}
	}
}

func TestGenerateModel_WithItemsAndAdditional(t *testing.T) {
	tt := templateTest{t, modelTemplate.Lookup("schema")}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		for _, k := range []string{"WithItemsAndAdditional", "WithItemsAndAdditional2"} {
			schema := definitions[k]
			genModel, err := makeGenDefinition(k, "models", schema, specDoc)
			if assert.NoError(t, err) {
				assert.Empty(t, genModel.Items)
				assert.True(t, genModel.IsComplexObject)
				prop := getDefinitionProperty(genModel, "tags")
				assert.NotEmpty(t, prop.Items)
				assert.True(t, prop.IsArray)
				assert.False(t, prop.IsComplexObject)
				buf := bytes.NewBuffer(nil)
				tt.template.Execute(buf, genModel)
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile("type "+k+" struct\\s*{"), res)
				// this would fail if it accepts additionalItems because it would come out as []interface{}
				assert.Regexp(t, regexp.MustCompile("Tags \\[\\]string `json:\"tags\"`"), res)
			}
		}
	}
}

func TestGenerateModel_SimpleTuple(t *testing.T) {
	tt := templateTest{t, modelTemplate}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "SimpleTuple"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) && assert.Empty(t, genModel.ExtraSchemas) {
			assert.True(t, genModel.IsTuple)
			assert.False(t, genModel.IsComplexObject)
			assert.False(t, genModel.IsArray)
			assert.False(t, genModel.IsAnonymous)
			assert.Equal(t, k, genModel.Name)
			assert.Equal(t, k, genModel.GoType)
			assert.Len(t, genModel.Properties, 5)
			buf := bytes.NewBuffer(nil)
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("swagger:model "+k), res)
			assert.Regexp(t, regexp.MustCompile("type "+k+" struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("P0 int64 `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P1 string `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P2 strfmt.DateTime `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P3 Notable `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P4 \\*Notable `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile(k+"\\) UnmarshalJSON"), res)
			assert.Regexp(t, regexp.MustCompile(k+"\\) MarshalJSON"), res)
			assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(data)")), res)

			for i, p := range genModel.Properties {
				r := "m.P" + strconv.Itoa(i)
				if !p.IsNullable {
					r = "&" + r
				}
				assert.Regexp(t, regexp.MustCompile("json.Unmarshal\\(stage1\\["+strconv.Itoa(i)+"\\], "+r+"\\)"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("P"+strconv.Itoa(i)+",")), res)
			}
		}
	}
}

func TestGenerateModel_TupleWithExtra(t *testing.T) {
	tt := templateTest{t, modelTemplate}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "TupleWithExtra"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) && assert.Empty(t, genModel.ExtraSchemas) {
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
			tt.template.Execute(buf, genModel)
			res := buf.String()
			assert.Regexp(t, regexp.MustCompile("swagger:model "+k), res)
			assert.Regexp(t, regexp.MustCompile("type "+k+" struct\\s*{"), res)
			assert.Regexp(t, regexp.MustCompile("P0 int64 `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P1 string `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P2 strfmt.DateTime `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("P3 Notable `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile("AdditionalItems \\[\\]float64 `json:\"-\"`"), res)
			assert.Regexp(t, regexp.MustCompile(k+"\\) UnmarshalJSON"), res)
			assert.Regexp(t, regexp.MustCompile(k+"\\) MarshalJSON"), res)

			for i, p := range genModel.Properties {
				r := "m.P" + strconv.Itoa(i)
				if !p.IsNullable {
					r = "&" + r
				}
				assert.Regexp(t, regexp.MustCompile("json.Unmarshal\\(stage1\\["+strconv.Itoa(i)+"\\], "+r+"\\)"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("P"+strconv.Itoa(i)+",")), res)
			}
			assert.Regexp(t, regexp.MustCompile("var lastIndex int"), res)
			assert.Regexp(t, regexp.MustCompile("var toadd float64"), res)
			assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("for _, val := range stage1[lastIndex+1:]")), res)
			assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(val, &toadd)")), res)
			assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(data)")), res)
			assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("data = append(data, m.AdditionalItems...)")), res)
		}
	}
}

func TestGenerateModel_WithTuple(t *testing.T) {
	tt := templateTest{t, modelTemplate}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "WithTuple"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.ExtraSchemas) && assert.NotEmpty(t, genModel.Properties) {
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
			err := tt.template.Execute(buf, genModel)
			if assert.NoError(t, err) {
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile("swagger:model "+k+"Flags"), res)
				assert.Regexp(t, regexp.MustCompile("type "+k+"FlagsTuple0 struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("P0 int64 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P1 string `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile(k+"FlagsTuple0\\) UnmarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(k+"FlagsTuple0\\) MarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(data)")), res)

				for i, p := range sch.Properties {
					r := "m.P" + strconv.Itoa(i)
					if !p.IsNullable {
						r = "&" + r
					}
					assert.Regexp(t, regexp.MustCompile("json.Unmarshal\\(stage1\\["+strconv.Itoa(i)+"\\], "+r+"\\)"), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("P"+strconv.Itoa(i)+",")), res)
				}
			}
		}
	}
}

func TestGenerateModel_WithTupleWithExtra(t *testing.T) {
	tt := templateTest{t, modelTemplate}
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		k := "WithTupleWithExtra"
		schema := definitions[k]
		genModel, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) && assert.NotEmpty(t, genModel.ExtraSchemas) && assert.NotEmpty(t, genModel.Properties) {
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
			err := tt.template.Execute(buf, genModel)
			if assert.NoError(t, err) {
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile("swagger:model "+k+"Flags"), res)
				assert.Regexp(t, regexp.MustCompile("type "+k+"FlagsTuple0 struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("P0 int64 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P1 string `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("AdditionalItems \\[\\]float32 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile(k+"FlagsTuple0\\) UnmarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(k+"FlagsTuple0\\) MarshalJSON"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(data)")), res)

				for i, p := range sch.Properties {
					r := "m.P" + strconv.Itoa(i)
					if !p.IsNullable {
						r = "&" + r
					}
					assert.Regexp(t, regexp.MustCompile("json.Unmarshal\\(stage1\\["+strconv.Itoa(i)+"\\], "+r+"\\)"), res)
					assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("P"+strconv.Itoa(i)+",")), res)
				}

				assert.Regexp(t, regexp.MustCompile("var lastIndex int"), res)
				assert.Regexp(t, regexp.MustCompile("var toadd float32"), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("for _, val := range stage1[lastIndex+1:]")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Unmarshal(val, &toadd)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("json.Marshal(data)")), res)
				assert.Regexp(t, regexp.MustCompile(regexp.QuoteMeta("data = append(data, m.AdditionalItems...)")), res)
			}
		}
	}
}

func TestGenerateModel_WithAllOf(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.models.yml")
	if assert.NoError(t, err) {
		definitions := specDoc.Spec().Definitions
		schema := definitions["WithAllOf"]
		genModel, err := makeGenDefinition("WithAllOf", "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.Len(t, genModel.AllOf, 6)
			assert.True(t, genModel.AllOf[1].HasAdditionalProperties)
			assert.True(t, genModel.IsComplexObject)
			assert.Equal(t, "WithAllOf", genModel.Name)
			assert.Equal(t, "WithAllOf", genModel.GoType)
			buf := bytes.NewBuffer(nil)
			err := modelTemplate.Execute(buf, genModel)
			if assert.NoError(t, err) {
				res := buf.String()
				assert.Regexp(t, regexp.MustCompile("type WithAllOf struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("type WithAllOfAO2AddedProps2 struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("type WithAllOfAO3Tuple3 struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("type WithAllOfAO4Tuple4 struct\\s*{"), res)
				assert.Regexp(t, regexp.MustCompile("Notable"), res)
				assert.Regexp(t, regexp.MustCompile("Title string `json:\"title\"`"), res)
				assert.Regexp(t, regexp.MustCompile("Body string `json:\"body\"`"), res)
				assert.Regexp(t, regexp.MustCompile("Name string `json:\"name\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P0 float32 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P0 float64 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P1 strfmt.DateTime `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("P1 strfmt.Date `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("AdditionalItems \\[\\]strfmt.Password `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("AdditionalProperties map\\[string\\]int32 `json:\"-\"`"), res)
				assert.Regexp(t, regexp.MustCompile("AdditionalProperties map\\[string\\]int64 `json:\"-\"`"), res)
			}
		}
	}
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
