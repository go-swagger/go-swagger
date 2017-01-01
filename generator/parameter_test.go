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
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestBodyParams(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.bodyparams.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, _, op, ok := b.Analyzed.OperationForName("updateTask")
	if assert.True(t, ok) && assert.NotNil(t, op) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		resolver.KnownDefs = make(map[string]struct{})
		for k := range b.Doc.Spec().Definitions {
			resolver.KnownDefs[k] = struct{}{}
		}
		for _, param := range op.Parameters {
			if param.Name == "body" {
				gp, err := b.MakeParameter("a", resolver, param, nil)
				if assert.NoError(t, err) {
					assert.True(t, gp.IsBodyParam())
					if assert.NotNil(t, gp.Schema) {
						assert.True(t, gp.Schema.IsComplexObject)
						assert.False(t, gp.Schema.IsAnonymous)
						assert.Equal(t, "models.Task", gp.Schema.GoType)
					}
				}
			}
		}
	}

	b, err = opBuilder("createTask", "../fixtures/codegen/todolist.bodyparams.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, _, op, ok = b.Analyzed.OperationForName("createTask")
	if assert.True(t, ok) && assert.NotNil(t, op) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		resolver.KnownDefs = make(map[string]struct{})
		for k := range b.Doc.Spec().Definitions {
			resolver.KnownDefs[k] = struct{}{}
		}
		for _, param := range op.Parameters {
			if param.Name == "body" {
				gp, err := b.MakeParameter("a", resolver, param, nil)
				if assert.NoError(t, err) {
					assert.True(t, gp.IsBodyParam())
					if assert.NotNil(t, gp.Schema) {
						assert.True(t, gp.Schema.IsComplexObject)
						assert.False(t, gp.Schema.IsAnonymous)
						assert.Equal(t, "CreateTaskBody", gp.Schema.GoType)

						gpe, ok := b.ExtraSchemas["CreateTaskBody"]
						assert.True(t, ok)
						assert.True(t, gpe.IsComplexObject)
						assert.False(t, gpe.IsAnonymous)
						assert.Equal(t, "CreateTaskBody", gpe.GoType)
					}
				}
			}
		}
	}
}

var arrayFormParams = []paramTestContext{
	{"siBool", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatBool", "swag.ConvertBool", nil}},
	{"siString", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"", "", nil}},
	{"siNested", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"", "", &paramItemsTestContext{"", "", &paramItemsTestContext{"", "", nil}}}},
	{"siInt", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt64", "swag.ConvertInt64", nil}},
	{"siInt32", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt32", "swag.ConvertInt32", nil}},
	{"siInt64", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt64", "swag.ConvertInt64", nil}},
	{"siFloat", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat64", "swag.ConvertFloat64", nil}},
	{"siFloat32", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat32", "swag.ConvertFloat32", nil}},
	{"siFloat64", "arrayFormParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat64", "swag.ConvertFloat64", nil}},
}

func TestFormArrayParams(t *testing.T) {
	b, err := opBuilder("arrayFormParams", "../fixtures/codegen/todolist.arrayform.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	for _, v := range arrayFormParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

var arrayQueryParams = []paramTestContext{
	{"siBool", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatBool", "swag.ConvertBool", nil}},
	{"siString", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"", "", nil}},
	{"siNested", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"", "", &paramItemsTestContext{"", "", &paramItemsTestContext{"", "", nil}}}},
	{"siInt", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt64", "swag.ConvertInt64", nil}},
	{"siInt32", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt32", "swag.ConvertInt32", nil}},
	{"siInt64", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatInt64", "swag.ConvertInt64", nil}},
	{"siFloat", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat64", "swag.ConvertFloat64", nil}},
	{"siFloat32", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat32", "swag.ConvertFloat32", nil}},
	{"siFloat64", "arrayQueryParams", "", "", codeGenOpBuilder{}, &paramItemsTestContext{"swag.FormatFloat64", "swag.ConvertFloat64", nil}},
}

func TestQueryArrayParams(t *testing.T) {
	b, err := opBuilder("arrayQueryParams", "../fixtures/codegen/todolist.arrayquery.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	for _, v := range arrayQueryParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

var simplePathParams = []paramTestContext{
	{"siBool", "simplePathParams", "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}, nil},
	{"siString", "simplePathParams", "", "", codeGenOpBuilder{}, nil},
	{"siInt", "simplePathParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siInt32", "simplePathParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siInt64", "simplePathParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siFloat", "simplePathParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
	{"siFloat32", "simplePathParams", "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}, nil},
	{"siFloat64", "simplePathParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
}

func TestSimplePathParams(t *testing.T) {
	b, err := opBuilder("simplePathParams", "../fixtures/codegen/todolist.simplepath.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simplePathParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

var simpleHeaderParams = []paramTestContext{
	{"id", "simpleHeaderParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siBool", "simpleHeaderParams", "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}, nil},
	{"siString", "simpleHeaderParams", "", "", codeGenOpBuilder{}, nil},
	{"siInt", "simpleHeaderParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siInt32", "simpleHeaderParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siInt64", "simpleHeaderParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siFloat", "simpleHeaderParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
	{"siFloat32", "simpleHeaderParams", "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}, nil},
	{"siFloat64", "simpleHeaderParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
}

func TestSimpleHeaderParams(t *testing.T) {
	b, err := opBuilder("simpleHeaderParams", "../fixtures/codegen/todolist.simpleheader.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleHeaderParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

var simpleFormParams = []paramTestContext{
	{"id", "simpleFormParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siBool", "simpleFormParams", "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}, nil},
	{"siString", "simpleFormParams", "", "", codeGenOpBuilder{}, nil},
	{"siInt", "simpleFormParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siInt32", "simpleFormParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siInt64", "simpleFormParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siFloat", "simpleFormParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
	{"siFloat32", "simpleFormParams", "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}, nil},
	{"siFloat64", "simpleFormParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
}

func TestSimpleFormParams(t *testing.T) {
	b, err := opBuilder("simpleFormParams", "../fixtures/codegen/todolist.simpleform.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleFormParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

var simpleQueryParams = []paramTestContext{
	{"id", "simpleQueryParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siBool", "simpleQueryParams", "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}, nil},
	{"siString", "simpleQueryParams", "", "", codeGenOpBuilder{}, nil},
	{"siInt", "simpleQueryParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siInt32", "simpleQueryParams", "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}, nil},
	{"siInt64", "simpleQueryParams", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"siFloat", "simpleQueryParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
	{"siFloat32", "simpleQueryParams", "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}, nil},
	{"siFloat64", "simpleQueryParams", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
}

func TestSimpleQueryParamsAST(t *testing.T) {
	b, err := opBuilder("simpleQueryParams", "../fixtures/codegen/todolist.simplequery.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleQueryParams {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

type paramItemsTestContext struct {
	Formatter string
	Converter string
	Items     *paramItemsTestContext
}

type paramTestContext struct {
	Name      string
	OpID      string
	Formatter string
	Converter string
	B         codeGenOpBuilder
	Items     *paramItemsTestContext
}

func (ctx *paramTestContext) assertParameter(t testing.TB) bool {
	_, _, op, err := ctx.B.Analyzed.OperationForName(ctx.OpID)
	if assert.True(t, err) && assert.NotNil(t, op) {
		resolver := &typeResolver{ModelsPackage: ctx.B.ModelsPackage, Doc: ctx.B.Doc}
		resolver.KnownDefs = make(map[string]struct{})
		for k := range ctx.B.Doc.Spec().Definitions {
			resolver.KnownDefs[k] = struct{}{}
		}
		for _, param := range op.Parameters {
			if param.Name == ctx.Name {
				gp, err := ctx.B.MakeParameter("a", resolver, param, nil)
				if assert.NoError(t, err) {
					return assert.True(t, ctx.assertGenParam(t, param, gp))
				}
			}
		}
		return false
	}
	return false
}

func (ctx *paramTestContext) assertGenParam(t testing.TB, param spec.Parameter, gp GenParameter) bool {
	// went with the verbose option here, easier to debug
	if !assert.Equal(t, param.In, gp.Location) {
		return false
	}
	if !assert.Equal(t, param.Name, gp.Name) {
		return false
	}
	if !assert.Equal(t, fmt.Sprintf("%q", param.Name), gp.Path) {
		return false
	}
	if !assert.Equal(t, "i", gp.IndexVar) {
		return false
	}
	if !assert.Equal(t, "a", gp.ReceiverName) {
		return false
	}
	if !assert.Equal(t, "a."+swag.ToGoName(param.Name), gp.ValueExpression) {
		return false
	}
	if !assert.Equal(t, ctx.Formatter, gp.Formatter) {
		return false
	}
	if !assert.Equal(t, ctx.Converter, gp.Converter) {
		return false
	}
	if !assert.Equal(t, param.Description, gp.Description) {
		return false
	}
	if !assert.Equal(t, param.CollectionFormat, gp.CollectionFormat) {
		return false
	}
	if !assert.Equal(t, param.Required, gp.Required) {
		return false
	}
	if !assert.Equal(t, param.Minimum, gp.Minimum) || !assert.Equal(t, param.ExclusiveMinimum, gp.ExclusiveMinimum) {
		return false
	}
	if !assert.Equal(t, param.Maximum, gp.Maximum) || !assert.Equal(t, param.ExclusiveMaximum, gp.ExclusiveMaximum) {
		return false
	}
	if !assert.Equal(t, param.MinLength, gp.MinLength) {
		return false
	}
	if !assert.Equal(t, param.MaxLength, gp.MaxLength) {
		return false
	}
	if !assert.Equal(t, param.Pattern, gp.Pattern) {
		return false
	}
	if !assert.Equal(t, param.MaxItems, gp.MaxItems) {
		return false
	}
	if !assert.Equal(t, param.MinItems, gp.MinItems) {
		return false
	}
	if !assert.Equal(t, param.UniqueItems, gp.UniqueItems) {
		return false
	}
	if !assert.Equal(t, param.MultipleOf, gp.MultipleOf) {
		return false
	}
	if !assert.EqualValues(t, param.Enum, gp.Enum) {
		return false
	}
	if !assert.Equal(t, param.Type, gp.SwaggerType) {
		return false
	}
	if !assert.Equal(t, param.Format, gp.SwaggerFormat) {
		return false
	}
	if _, ok := primitives[gp.GoType]; ok {
		if !assert.True(t, gp.IsPrimitive) {
			return false
		}
	} else {
		if !assert.False(t, gp.IsPrimitive) {
			return false
		}
	}
	// verify rendered template
	if param.In == "body" {
		if !assertBodyParam(t, param, gp) {
			return false
		}
		return true
	}

	if ctx.Items != nil {
		return ctx.Items.Assert(t, param.Items, gp.Child)
	}

	return true
}

func assertBodyParam(t testing.TB, param spec.Parameter, gp GenParameter) bool {
	if !assert.Equal(t, "body", param.In) || !assert.Equal(t, "body", gp.Location) {
		return false
	}
	if !assert.NotNil(t, gp.Schema) {
		return false
	}
	return true
}

func (ctx *paramItemsTestContext) Assert(t testing.TB, pItems *spec.Items, gpItems *GenItems) bool {
	if !assert.NotNil(t, pItems) || !assert.NotNil(t, gpItems) {
		return false
	}
	// went with the verbose option here, easier to debug
	if !assert.Equal(t, ctx.Formatter, gpItems.Formatter) {
		return false
	}
	if !assert.Equal(t, ctx.Converter, gpItems.Converter) {
		return false
	}
	if !assert.Equal(t, pItems.CollectionFormat, gpItems.CollectionFormat) {
		return false
	}
	if !assert.Equal(t, pItems.Minimum, gpItems.Minimum) || !assert.Equal(t, pItems.ExclusiveMinimum, gpItems.ExclusiveMinimum) {
		return false
	}
	if !assert.Equal(t, pItems.Maximum, gpItems.Maximum) || !assert.Equal(t, pItems.ExclusiveMaximum, gpItems.ExclusiveMaximum) {
		return false
	}
	if !assert.Equal(t, pItems.MinLength, gpItems.MinLength) {
		return false
	}
	if !assert.Equal(t, pItems.MaxLength, gpItems.MaxLength) {
		return false
	}
	if !assert.Equal(t, pItems.Pattern, gpItems.Pattern) {
		return false
	}
	if !assert.Equal(t, pItems.MaxItems, gpItems.MaxItems) {
		return false
	}
	if !assert.Equal(t, pItems.MinItems, gpItems.MinItems) {
		return false
	}
	if !assert.Equal(t, pItems.UniqueItems, gpItems.UniqueItems) {
		return false
	}
	if !assert.Equal(t, pItems.MultipleOf, gpItems.MultipleOf) {
		return false
	}
	if !assert.EqualValues(t, pItems.Enum, gpItems.Enum) {
		return false
	}
	if !assert.Equal(t, pItems.Type, gpItems.SwaggerType) {
		return false
	}
	if !assert.Equal(t, pItems.Format, gpItems.SwaggerFormat) {
		return false
	}
	if ctx.Items != nil {
		return ctx.Items.Assert(t, pItems.Items, gpItems.Child)
	}
	return true

}

var bug163Properties = []paramTestContext{
	{"stringTypeInQuery", "getSearch", "", "", codeGenOpBuilder{}, nil},
	{"numberTypeInQuery", "getSearch", "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}, nil},
	{"integerTypeInQuery", "getSearch", "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}, nil},
	{"booleanTypeInQuery", "getSearch", "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}, nil},
}

func TestGenParameters_Simple(t *testing.T) {
	b, err := opBuilder("getSearch", "../fixtures/bugs/163/swagger.yml")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range bug163Properties {
		v.B = b
		if !v.assertParameter(t) {
			t.FailNow()
		}
	}
}

func TestGenParameter_Issue163(t *testing.T) {
	b, err := opBuilder("getSearch", "../fixtures/bugs/163/swagger.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("get_search_parameters.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "var stringTypeInQueryDefault string = string(\"qsValue\")", res)
					assertInCode(t, "o.StringTypeInQuery = &stringTypeInQueryDefault", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue195(t *testing.T) {
	b, err := opBuilder("getTesting", "../fixtures/bugs/195/swagger.json")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("get_testing.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "TestingThis *int64", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue196(t *testing.T) {
	b, err := opBuilder("postEvents", "../fixtures/bugs/196/swagger.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("post_events.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "body.Validate", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue217(t *testing.T) {
	// Check for string

	assertNoValidator(t, "postEcho", "../fixtures/bugs/217/string.yml")
	assertNoValidator(t, "postEcho", "../fixtures/bugs/217/interface.yml")
	assertNoValidator(t, "postEcho", "../fixtures/bugs/217/map.yml")
	assertNoValidator(t, "postEcho", "../fixtures/bugs/217/array.yml")
}

func assertNoValidator(t testing.TB, opName, path string) {
	b, err := opBuilder(opName, path)
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(&buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("post_echo.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertNotInCode(t, "body.Validate", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue249(t *testing.T) {
	b, err := opBuilder("putTesting", "../fixtures/bugs/249/swagger.json")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("put_testing.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertNotInCode(t, "valuesTestingThis := o.TestingThis", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue248(t *testing.T) {
	b, err := opBuilder("CreateThing", "../fixtures/bugs/248/swagger.json")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("create_thing.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, ", *o.OptionalQueryEnum", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue350(t *testing.T) {
	b, err := opBuilder("withBoolDefault", "../fixtures/codegen/todolist.allparams.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("with_bool_default.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "Verbose: &verboseDefault", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue351(t *testing.T) {
	b, err := opBuilder("withArray", "../fixtures/codegen/todolist.allparams.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(t, err) {
				ff, err := opts.LanguageOpts.FormatContent("with_array.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "validate.MinLength(fmt.Sprintf(\"%s.%v\", \"sha256\", i), \"query\", sha256I, 64)", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue511(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("postModels", "../fixtures/bugs/511/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertNotInCode(t, "fds := runtime.Values(r.Form)", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue628_Collection(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("collection", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `workspaceIDI, err := formats.Parse(workspaceIDIV)`, res)
					assertInCode(t, `workspaceIDIR = append(workspaceIDIR, workspaceIDI)`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue628_Single(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("single", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `value, err := formats.Parse("uuid", raw)`, res)
					assertInCode(t, `o.WorkspaceID = *(value.(*strfmt.UUID))`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue628_Details(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("details", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `value, err := formats.Parse("uuid", raw)`, res)
					assertInCode(t, `o.ID = *(value.(*strfmt.UUID))`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue731_Collection(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("collection", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `for _, v := range o.WorkspaceID`, res)
					assertInCode(t, `valuesWorkspaceID = append(valuesWorkspaceID, v.String())`, res)
					assertInCode(t, `joinedWorkspaceID := swag.JoinByFormat(valuesWorkspaceID, "")`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue731_Single(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("single", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `qWorkspaceID := qrWorkspaceID.String()`, res)
					assertInCode(t, `r.SetQueryParam("workspace_id", qWorkspaceID)`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue731_Details(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("details", "../fixtures/bugs/628/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `r.SetPathParam("id", o.ID.String())`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue809_Client(t *testing.T) {
	assert := assert.New(t)

	gen, err := methodPathOpBuilder("get", "/foo", "../fixtures/bugs/809/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, "valuesGroups", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue809_Server(t *testing.T) {
	assert := assert.New(t)

	gen, err := methodPathOpBuilder("get", "/foo", "../fixtures/bugs/809/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_models.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, "groupsIC := rawData", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue710(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("createTask", "../fixtures/codegen/todolist.allparams.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("create_task_parameter.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, "(typeVar", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue776_LocalFileRef(t *testing.T) {
	b, err := opBuilder("GetItem", "../fixtures/bugs/776/param.yaml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverParameter").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "Body *GetItemParamsBody", string(ff))
					assertNotInCode(t, "type GetItemParamsBody struct", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
			var buf2 bytes.Buffer
			if assert.NoError(t, templates.MustGet("serverOperation").Execute(&buf2, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf2.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "type GetItemParamsBody struct", string(ff))
				} else {
					fmt.Println(buf2.String())
				}
			}
		}
	}

}

func TestGenParameter_ArrayQueryParameters(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("arrayQueryParams", "../fixtures/codegen/todolist.arrayquery.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("array_query_params.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `siBoolIC := swag.SplitByFormat(qvSiBool, "ssv")`, res)
					assertInCode(t, `var siBoolIR []bool`, res)
					assertInCode(t, `for i, siBoolIV := range siBoolIC`, res)
					assertInCode(t, `siBoolI, err := swag.ConvertBool(siBoolIV)`, res)
					assertInCode(t, `siBoolIR = append(siBoolIR, siBoolI)`, res)
					assertInCode(t, `o.SiBool = siBoolIR`, res)
					assertInCode(t, `siBoolSize := int64(len(o.SiBool))`, res)
					assertInCode(t, `err := validate.MinItems("siBool", "query", siBoolSize, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siBool", "query", siBoolSize, 50)`, res)

					assertInCode(t, `siFloatIC := rawData`, res)
					assertInCode(t, `var siFloatIR []float64`, res)
					assertInCode(t, `for i, siFloatIV := range siFloatIC`, res)
					assertInCode(t, `siFloatI, err := swag.ConvertFloat64(siFloatIV)`, res)
					assertInCode(t, `return errors.InvalidType(fmt.Sprintf("%s.%v", "siFloat", i), "query", "float64", siFloatI)`, res)
					assertInCode(t, `err := validate.Minimum(fmt.Sprintf("%s.%v", "siFloat", i), "query", float64(siFloatI), 3, true)`, res)
					assertInCode(t, `err := validate.Maximum(fmt.Sprintf("%s.%v", "siFloat", i), "query", float64(siFloatI), 100, true); err != nil`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siFloat", i), "query", float64(siFloatI), 1.5)`, res)
					assertInCode(t, `siFloatIR = append(siFloatIR, siFloatI)`, res)
					assertInCode(t, `o.SiFloat = siFloatIR`, res)
					assertInCode(t, `siFloatSize := int64(len(o.SiFloat))`, res)
					assertInCode(t, `err := validate.MinItems("siFloat", "query", siFloatSize, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siFloat", "query", siFloatSize, 50)`, res)

					assertInCode(t, `siFloat32IC := swag.SplitByFormat(qvSiFloat32, "")`, res)
					assertInCode(t, `var siFloat32IR []float32`, res)
					assertInCode(t, `for i, siFloat32IV := range siFloat32IC`, res)
					assertInCode(t, `siFloat32I, err := swag.ConvertFloat32(siFloat32IV)`, res)
					assertInCode(t, `err := validate.Minimum(fmt.Sprintf("%s.%v", "siFloat32", i), "query", float64(siFloat32I), 3, true)`, res)
					assertInCode(t, `err := validate.Maximum(fmt.Sprintf("%s.%v", "siFloat32", i), "query", float64(siFloat32I), 100, true)`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siFloat32", i), "query", float64(siFloat32I), 1.5)`, res)
					assertInCode(t, `siFloat32IR = append(siFloat32IR, siFloat32I)`, res)
					assertInCode(t, `o.SiFloat32 = siFloat32IR`, res)

					assertInCode(t, `siFloat64IC := swag.SplitByFormat(qvSiFloat64, "pipes")`, res)
					assertInCode(t, `var siFloat64IR []float64`, res)
					assertInCode(t, `for i, siFloat64IV := range siFloat64IC`, res)
					assertInCode(t, `siFloat64I, err := swag.ConvertFloat64(siFloat64IV)`, res)
					assertInCode(t, `err := validate.Minimum(fmt.Sprintf("%s.%v", "siFloat64", i), "query", float64(siFloat64I), 3, true)`, res)
					assertInCode(t, `err := validate.Maximum(fmt.Sprintf("%s.%v", "siFloat64", i), "query", float64(siFloat64I), 100, true)`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siFloat64", i), "query", float64(siFloat64I), 1.5)`, res)
					assertInCode(t, `siFloat64IR = append(siFloat64IR, siFloat64I)`, res)
					assertInCode(t, `o.SiFloat64 = siFloat64IR`, res)
					assertInCode(t, `siFloat64Size := int64(len(o.SiFloat64))`, res)
					assertInCode(t, `err := validate.MinItems("siFloat64", "query", siFloat64Size, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siFloat64", "query", siFloat64Size, 50)`, res)

					assertInCode(t, `siIntIC := swag.SplitByFormat(qvSiInt, "pipes")`, res)
					assertInCode(t, `var siIntIR []int64`, res)
					assertInCode(t, `for i, siIntIV := range siIntIC`, res)
					assertInCode(t, `siIntI, err := swag.ConvertInt64(siIntIV)`, res)
					assertInCode(t, `err := validate.MinimumInt(fmt.Sprintf("%s.%v", "siInt", i), "query", int64(siIntI), 8, true)`, res)
					assertInCode(t, `err := validate.MaximumInt(fmt.Sprintf("%s.%v", "siInt", i), "query", int64(siIntI), 100, true)`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siInt", i), "query", float64(siIntI), 2)`, res)
					assertInCode(t, `siIntIR = append(siIntIR, siIntI)`, res)
					assertInCode(t, `o.SiInt = siIntIR`, res)
					assertInCode(t, `siIntSize := int64(len(o.SiInt))`, res)
					assertInCode(t, `err := validate.MinItems("siInt", "query", siIntSize, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siInt", "query", siIntSize, 50)`, res)

					assertInCode(t, `siInt32IC := swag.SplitByFormat(qvSiInt32, "tsv")`, res)
					assertInCode(t, `var siInt32IR []int32`, res)
					assertInCode(t, `for i, siInt32IV := range siInt32IC`, res)
					assertInCode(t, `siInt32I, err := swag.ConvertInt32(siInt32IV)`, res)
					assertInCode(t, `err := validate.MinimumInt(fmt.Sprintf("%s.%v", "siInt32", i), "query", int64(siInt32I), 8, true)`, res)
					assertInCode(t, `err := validate.MaximumInt(fmt.Sprintf("%s.%v", "siInt32", i), "query", int64(siInt32I), 100, true)`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siInt32", i), "query", float64(siInt32I), 2)`, res)
					assertInCode(t, `siInt32IR = append(siInt32IR, siInt32I)`, res)
					assertInCode(t, `o.SiInt32 = siInt32IR`, res)
					assertInCode(t, `siFloat32Size := int64(len(o.SiFloat32))`, res)
					assertInCode(t, `err := validate.MinItems("siFloat32", "query", siFloat32Size, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siFloat32", "query", siFloat32Size, 50)`, res)
					assertInCode(t, `siInt32Size := int64(len(o.SiInt32))`, res)
					assertInCode(t, `err := validate.MinItems("siInt32", "query", siInt32Size, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siInt32", "query", siInt32Size, 50)`, res)

					assertInCode(t, `siInt64IC := swag.SplitByFormat(qvSiInt64, "ssv")`, res)
					assertInCode(t, `var siInt64IR []int64`, res)
					assertInCode(t, `for i, siInt64IV := range siInt64IC`, res)
					assertInCode(t, `siInt64I, err := swag.ConvertInt64(siInt64IV)`, res)
					assertInCode(t, `err := validate.MinimumInt(fmt.Sprintf("%s.%v", "siInt64", i), "query", int64(siInt64I), 8, true)`, res)
					assertInCode(t, `err := validate.MaximumInt(fmt.Sprintf("%s.%v", "siInt64", i), "query", int64(siInt64I), 100, true)`, res)
					assertInCode(t, `err := validate.MultipleOf(fmt.Sprintf("%s.%v", "siInt64", i), "query", float64(siInt64I), 2)`, res)
					assertInCode(t, `siInt64IR = append(siInt64IR, siInt64I)`, res)
					assertInCode(t, `o.SiInt64 = siInt64IR`, res)
					assertInCode(t, `siInt64Size := int64(len(o.SiInt64))`, res)
					assertInCode(t, `err := validate.MinItems("siInt64", "query", siInt64Size, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siInt64", "query", siInt64Size, 50)`, res)

					assertInCode(t, `siStringIC := swag.SplitByFormat(qvSiString, "csv")`, res)
					assertInCode(t, `var siStringIR []string`, res)
					assertInCode(t, `for i, siStringIV := range siStringIC`, res)
					assertInCode(t, `siStringI := siStringIV`, res)
					assertInCode(t, `err := validate.MinLength(fmt.Sprintf("%s.%v", "siString", i), "query", siStringI, 5)`, res)
					assertInCode(t, `err := validate.MaxLength(fmt.Sprintf("%s.%v", "siString", i), "query", siStringI, 50)`, res)
					assertInCode(t, `err := validate.Pattern(fmt.Sprintf("%s.%v", "siString", i), "query", siStringI, `+"`"+`[A-Z][\w-]+`+"`"+`)`, res)
					assertInCode(t, `siStringIR = append(siStringIR, siStringI)`, res)
					assertInCode(t, `o.SiString = siStringIR`, res)
					assertInCode(t, `siStringSize := int64(len(o.SiString))`, res)
					assertInCode(t, `err := validate.MinItems("siString", "query", siStringSize, 5)`, res)
					assertInCode(t, `err := validate.MaxItems("siString", "query", siStringSize, 50)`, res)

					assertInCode(t, `siNestedIC := rawData`, res)
					assertInCode(t, `var siNestedIR [][][]string`, res)
					assertInCode(t, `for i, siNestedIV := range siNestedIC`, res)
					assertInCode(t, `siNestedIIC := swag.SplitByFormat(siNestedIV, "pipes")`, res)
					assertInCode(t, `var siNestedIIR [][]string`, res)
					assertInCode(t, `for ii, siNestedIIV := range siNestedIIC {`, res)
					assertInCode(t, `siNestedIIIC := swag.SplitByFormat(siNestedIIV, "csv")`, res)
					assertInCode(t, `var siNestedIIIR []string`, res)
					assertInCode(t, `for iii, siNestedIIIV := range siNestedIIIC`, res)
					assertInCode(t, `siNestedIII := siNestedIIIV`, res)
					assertInCode(t, `err := validate.MinLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), iii), "query", siNestedIII, 5)`, res)
					assertInCode(t, `err := validate.MaxLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), iii), "query", siNestedIII, 50)`, res)
					assertInCode(t, `err := validate.Pattern(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), iii), "query", siNestedIII, `+"`"+`[A-Z][\w-]+`+"`"+`)`, res)
					assertInCode(t, `siNestedIIIR = append(siNestedIIIR, siNestedIII)`, res)
					assertInCode(t, `siNestedIIiSize := int64(len(siNestedII))`, res)
					assertInCode(t, `err := validate.MinItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), "query", siNestedIIiSize, 3)`, res)
					assertInCode(t, `err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), "query", siNestedIIiSize, 30)`, res)
					assertInCode(t, `siNestedIIR = append(siNestedIIR, siNestedII)`, res)
					assertInCode(t, `siNestedISize := int64(len(siNestedI))`, res)
					assertInCode(t, `err := validate.MinItems(fmt.Sprintf("%s.%v", "siNested", i), "query", siNestedISize, 2)`, res)
					assertInCode(t, `err := validate.MaxItems(fmt.Sprintf("%s.%v", "siNested", i), "query", siNestedISize, 20)`, res)
					assertInCode(t, `siNestedIR = append(siNestedIR, siNestedI)`, res)
					assertInCode(t, `o.SiNested = siNestedIR`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}
