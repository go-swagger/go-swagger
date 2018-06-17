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
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
				gp, perr := b.MakeParameter("a", resolver, param, nil)
				if assert.NoError(t, perr) {
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
		return assertBodyParam(t, param, gp)
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
					// NOTE(fredbi): removed default values resolution from private details (defaults are resolved in NewXXXParams())
					assertInCode(t, "stringTypeInQueryDefault = string(\"qsValue\")", res)
					assertInCode(t, "StringTypeInQuery: &stringTypeInQueryDefault", res)
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
					assertInCode(t, `value, err := formats.Parse("uuid", workspaceIDIV)`, res) // NOTE(fredbi): added type assertion
					assertInCode(t, `workspaceIDI := *(value.(*strfmt.UUID))`, res)
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

func TestGenParameter_Issue1010_Server(t *testing.T) {
	assert := assert.New(t)

	gen, err := methodPathOpBuilder("get", "/widgets/", "../fixtures/bugs/1010/swagger.yml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("serverParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("get_widgets.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, "validate.Pattern(fmt.Sprintf(\"%s.%v\", \"category_id\", i), \"query\", categoryIDI, `^[0-9abcdefghjkmnpqrtuvwxyz]{29}$`)", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue710(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

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
	spec.Debug = true
	log.SetOutput(ioutil.Discard)
	defer func() {
		spec.Debug = false
		log.SetOutput(os.Stdout)
	}()
	b, err := opBuilderWithFlatten("GetItem", "../fixtures/bugs/776/param.yaml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverParameter").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "Body *models.Item", string(ff))
					assertNotInCode(t, "type GetItemParamsBody struct", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}

}

func TestGenParameter_Issue1111(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("start-es-cluster-instances", "../fixtures/bugs/1111/arrayParam.json")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_clusters_elasticsearch_cluster_id_instances_instance_ids_start_parameters.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `r.SetPathParam("instance_ids", joinedInstanceIds[0])`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue1462(t *testing.T) {
	assert := assert.New(t)

	gen, err := opBuilder("start-es-cluster-instances", "../fixtures/bugs/1462/arrayParam.json")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("post_clusters_elasticsearch_cluster_id_instances_instance_ids_start_parameters.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, `if len(joinedInstanceIds) > 0 {`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue1199(t *testing.T) {
	assert := assert.New(t)
	var assertion = `if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}`
	gen, err := opBuilder("move-clusters", "../fixtures/bugs/1199/nonEmptyBody.json")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("move_clusters_parameters.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, assertion, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenParameter_Issue1325(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	assert := assert.New(t)

	gen, err := opBuilder("uploadFile", "../fixtures/bugs/1325/swagger.yaml")
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
					assertInCode(t, "runtime.NamedReadCloser", res)
				} else {
					fmt.Println(buf.String())
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
					assertInCode(t, `siNestedIIiSize := int64(len(siNestedIIIC))`, res) // NOTE(fredbi): fixed variable (nested arrays)
					assertInCode(t, `err := validate.MinItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), "query", siNestedIIiSize, 3)`, res)
					assertInCode(t, `err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "siNested", i), ii), "query", siNestedIIiSize, 30)`, res)
					assertInCode(t, `siNestedIIR = append(siNestedIIR, siNestedIIIR)`, res) // NOTE(fredbi): fixed variable (nested arrays)
					assertInCode(t, `siNestedISize := int64(len(siNestedIIC))`, res)        //NOTE(fredbi): fixed variable (nested arrays)
					assertInCode(t, `err := validate.MinItems(fmt.Sprintf("%s.%v", "siNested", i), "query", siNestedISize, 2)`, res)
					assertInCode(t, `err := validate.MaxItems(fmt.Sprintf("%s.%v", "siNested", i), "query", siNestedISize, 20)`, res)
					assertInCode(t, `siNestedIR = append(siNestedIR, siNestedIIR)`, res) // NOTE(fredbi): fixed variable (nested arrays)
					assertInCode(t, `o.SiNested = siNestedIR`, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func assertParams(t *testing.T, fixtureConfig map[string]map[string][]string, fixture string) {
	fixtureSpec := path.Base(fixture)
	tassert := assert.New(t)
	for fixtureIndex, fixtureContents := range fixtureConfig {
		gen, err := opBuilderWithFlatten(fixtureIndex, fixture)
		if tassert.NoError(err) {
			op, err := gen.MakeOperation()
			if tassert.NoError(err) {
				opts := opts()
				opts.FlattenSpec = true
				for fixtureTemplate, expectedCode := range fixtureContents {
					buf := bytes.NewBuffer(nil)
					err := templates.MustGet(fixtureTemplate).Execute(buf, op)
					if tassert.NoError(err, "Expected generation to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
						ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
						if tassert.NoError(err, "Expected formatting to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
							res := string(ff)
							for line, codeLine := range expectedCode {
								if !assertInCode(t, strings.TrimSpace(codeLine), res) {
									t.Logf("Code expected did not match for fixture %s at line %d", fixtureSpec, line)
								}
							}
						} else {
							fmt.Println(buf.String())
						}
					}
				}
			}
		}
	}
}
func TestGenParameter_Issue909(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	assert := assert.New(t)
	fixtureConfig := map[string]map[string][]string{
		"1": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`strfmt "github.com/go-openapi/strfmt"`,
				`NotAnOption1 *strfmt.DateTime`,
				`NotAnOption2 *strfmt.UUID`,
				`NotAnOption3 *models.ContainerConfig`,
				`value, err := formats.Parse("date-time", raw)`,
				`o.NotAnOption1 = (value.(*strfmt.DateTime))`,
				`if err := o.validateNotAnOption1(formats); err != nil {`,
				`if err := validate.FormatOf("notAnOption1", "query", "date-time", o.NotAnOption1.String(), formats); err != nil {`,
				`value, err := formats.Parse("uuid", raw)`,
				`o.NotAnOption2 = (value.(*strfmt.UUID))`,
				`if err := o.validateNotAnOption2(formats); err != nil {`,
				`if err := validate.FormatOf("notAnOption2", "query", "uuid", o.NotAnOption2.String(), formats); err != nil {`,
			},
		},
		"2": map[string][]string{
			"serverParameter": []string{
				// expected code lines
				`"github.com/go-openapi/validate"`,
				`IsAnOption2 []strfmt.UUID`,
				`NotAnOption1 []strfmt.DateTime`,
				`NotAnOption3 *models.ContainerConfig`,
				`isAnOption2IC := swag.SplitByFormat(qvIsAnOption2, "csv")`,
				`var isAnOption2IR []strfmt.UUID`,
				`for i, isAnOption2IV := range isAnOption2IC {`,
				`value, err := formats.Parse("uuid", isAnOption2IV)`,
				`isAnOption2I := *(value.(*strfmt.UUID))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", "isAnOption2", i), "query", "uuid", isAnOption2I.String(), formats); err != nil {`,
				`isAnOption2IR = append(isAnOption2IR, isAnOption2I)`,
				`o.IsAnOption2 = isAnOption2IR`,
				`return errors.Required("notAnOption1", "query")`,
				`notAnOption1IC := swag.SplitByFormat(qvNotAnOption1, "csv")`,
				`var notAnOption1IR []strfmt.DateTime`,
				`for i, notAnOption1IV := range notAnOption1IC {`,
				`value, err := formats.Parse("date-time", notAnOption1IV)`,
				`return errors.InvalidType(fmt.Sprintf("%s.%v", "notAnOption1", i), "query", "strfmt.DateTime", value)`,
				`notAnOption1I := *(value.(*strfmt.DateTime))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", "notAnOption1", i), "query", "date-time", notAnOption1I.String(), formats); err != nil {`,
				`notAnOption1IR = append(notAnOption1IR, notAnOption1I)`,
				`o.NotAnOption1 = notAnOption1IR`,
			},
		},
		"3": map[string][]string{
			"serverParameter": []string{
				// expected code lines
				`"github.com/go-openapi/validate"`,
				`strfmt "github.com/go-openapi/strfmt"`,
				`IsAnOption2 [][]strfmt.UUID`,
				`IsAnOption4 [][][]strfmt.UUID`,
				`IsAnOptionalHeader [][]strfmt.UUID`,
				`NotAnOption1 [][]strfmt.DateTime`,
				`NotAnOption3 *models.ContainerConfig`,
				`isAnOption2IC := swag.SplitByFormat(qvIsAnOption2, "pipes")`,
				`var isAnOption2IR [][]strfmt.UUID`,
				`for i, isAnOption2IV := range isAnOption2IC {`,
				`isAnOption2IIC := swag.SplitByFormat(isAnOption2IV, "")`,
				`if len(isAnOption2IIC) > 0 {`,
				`var isAnOption2IIR []strfmt.UUID`,
				`for ii, isAnOption2IIV := range isAnOption2IIC {`,
				`value, err := formats.Parse("uuid", isAnOption2IIV)`,
				`isAnOption2II := *(value.(*strfmt.UUID))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption2", i), ii), "query", "uuid", isAnOption2II.String(), formats); err != nil {`,
				`isAnOption2IIR = append(isAnOption2IIR, isAnOption2II)`,
				`isAnOption2IR = append(isAnOption2IR, isAnOption2IIR)`,
				`o.IsAnOption2 = isAnOption2IR`,
				`isAnOption4IC := swag.SplitByFormat(qvIsAnOption4, "csv")`,
				`var isAnOption4IR [][][]strfmt.UUID`,
				`for i, isAnOption4IV := range isAnOption4IC {`,
				`isAnOption4IIC := swag.SplitByFormat(isAnOption4IV, "")`,
				`if len(isAnOption4IIC) > 0 {`,
				`var isAnOption4IIR [][]strfmt.UUID`,
				`for ii, isAnOption4IIV := range isAnOption4IIC {`,
				`isAnOption4IIIC := swag.SplitByFormat(isAnOption4IIV, "pipes")`,
				`if len(isAnOption4IIIC) > 0 {`,
				`var isAnOption4IIIR []strfmt.UUID`,
				`for iii, isAnOption4IIIV := range isAnOption4IIIC {`,
				`value, err := formats.Parse("uuid", isAnOption4IIIV)`,
				`isAnOption4III := *(value.(*strfmt.UUID))`,
				`if err := validate.Enum(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), iii), "query", isAnOption4III.String(), []interface{}{"a", "b", "c"}); err != nil {`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), iii), "query", "uuid", isAnOption4III.String(), formats); err != nil {`,
				`isAnOption4IIIR = append(isAnOption4IIIR, isAnOption4III)`,
				`isAnOption4IIR = append(isAnOption4IIR, isAnOption4IIIR)`,
				`isAnOption4IIiSize := int64(len(isAnOption4IIIC))`,
				`if err := validate.MinItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), "query", isAnOption4IIiSize, 3); err != nil {`,
				`if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), "query", isAnOption4IIIC); err != nil {`,
				`isAnOption4IR = append(isAnOption4IR, isAnOption4IIR)`,
				`if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "isAnOption4", i), "query", isAnOption4IIC); err != nil {`,
				`o.IsAnOption4 = isAnOption4IR`,
				`if err := o.validateIsAnOption4(formats); err != nil {`,
				`if err := validate.MaxItems("isAnOption4", "query", isAnOption4Size, 4); err != nil {`,
				`isAnOptionalHeaderIC := swag.SplitByFormat(qvIsAnOptionalHeader, "pipes")`,
				`var isAnOptionalHeaderIR [][]strfmt.UUID`,
				`for i, isAnOptionalHeaderIV := range isAnOptionalHeaderIC {`,
				`isAnOptionalHeaderIIC := swag.SplitByFormat(isAnOptionalHeaderIV, "")`,
				`if len(isAnOptionalHeaderIIC) > 0 {`,
				`var isAnOptionalHeaderIIR []strfmt.UUID`,
				`for ii, isAnOptionalHeaderIIV := range isAnOptionalHeaderIIC {`,
				`value, err := formats.Parse("uuid", isAnOptionalHeaderIIV)`,
				`isAnOptionalHeaderII := *(value.(*strfmt.UUID))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOptionalHeader", i), ii), "header", "uuid", isAnOptionalHeaderII.String(), formats); err != nil {`,
				`isAnOptionalHeaderIIR = append(isAnOptionalHeaderIIR, isAnOptionalHeaderII)`,
				`isAnOptionalHeaderIR = append(isAnOptionalHeaderIR, isAnOptionalHeaderIIR)`,
				`o.IsAnOptionalHeader = isAnOptionalHeaderIR`,
				`if err := o.validateIsAnOptionalHeader(formats); err != nil {`,
				`if err := validate.UniqueItems("isAnOptionalHeader", "header", o.IsAnOptionalHeader); err != nil {`,
				`notAnOption1IC := swag.SplitByFormat(qvNotAnOption1, "csv")`,
				`var notAnOption1IR [][]strfmt.DateTime`,
				`for i, notAnOption1IV := range notAnOption1IC {`,
				`notAnOption1IIC := swag.SplitByFormat(notAnOption1IV, "pipes")`,
				`if len(notAnOption1IIC) > 0 {`,
				`var notAnOption1IIR []strfmt.DateTime`,
				`for ii, notAnOption1IIV := range notAnOption1IIC {`,
				`value, err := formats.Parse("date-time", notAnOption1IIV)`,
				`notAnOption1II := *(value.(*strfmt.DateTime))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "notAnOption1", i), ii), "query", "date-time", notAnOption1II.String(), formats); err != nil {`,
				`notAnOption1IIR = append(notAnOption1IIR, notAnOption1II)`,
				`notAnOption1IR = append(notAnOption1IR, notAnOption1IIR)`,
				`o.NotAnOption1 = notAnOption1IR`,
			},
		},
		"4": map[string][]string{
			"serverParameter": []string{
				// expected code lines
				`"github.com/go-openapi/validate"`,
				`strfmt "github.com/go-openapi/strfmt"`,
				`IsAnOption2 [][]strfmt.UUID`,
				`IsAnOption4 [][][]strfmt.UUID`,
				`NotAnOption1 [][]strfmt.DateTime`,
				`NotAnOption3 *models.ContainerConfig`,
				`isAnOption2IC := swag.SplitByFormat(qvIsAnOption2, "")`,
				`var isAnOption2IR [][]strfmt.UUID`,
				`for i, isAnOption2IV := range isAnOption2IC {`,
				`isAnOption2IIC := swag.SplitByFormat(isAnOption2IV, "pipes")`,
				`if len(isAnOption2IIC) > 0 {`,
				`var isAnOption2IIR []strfmt.UUID`,
				`for ii, isAnOption2IIV := range isAnOption2IIC {`,
				`value, err := formats.Parse("uuid", isAnOption2IIV)`,
				`isAnOption2II := *(value.(*strfmt.UUID))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption2", i), ii), "query", "uuid", isAnOption2II.String(), formats); err != nil {`,
				`isAnOption2IIR = append(isAnOption2IIR, isAnOption2II)`,
				`isAnOption2IR = append(isAnOption2IR, isAnOption2IIR)`,
				`o.IsAnOption2 = isAnOption2IR`,
				`isAnOption4IC := swag.SplitByFormat(qvIsAnOption4, "")`,
				`var isAnOption4IR [][][]strfmt.UUID`,
				`for i, isAnOption4IV := range isAnOption4IC {`,
				`isAnOption4IIC := swag.SplitByFormat(isAnOption4IV, "pipes")`,
				`if len(isAnOption4IIC) > 0 {`,
				`var isAnOption4IIR [][]strfmt.UUID`,
				`for ii, isAnOption4IIV := range isAnOption4IIC {`,
				`isAnOption4IIIC := swag.SplitByFormat(isAnOption4IIV, "tsv")`,
				`if len(isAnOption4IIIC) > 0 {`,
				`var isAnOption4IIIR []strfmt.UUID`,
				`for iii, isAnOption4IIIV := range isAnOption4IIIC {`,
				`value, err := formats.Parse("uuid", isAnOption4IIIV)`,
				`isAnOption4III := *(value.(*strfmt.UUID))`,
				`if err := validate.Enum(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), iii), "query", isAnOption4III.String(), []interface{}{"a", "b", "c"}); err != nil {`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), iii), "query", "uuid", isAnOption4III.String(), formats); err != nil {`,
				`isAnOption4IIIR = append(isAnOption4IIIR, isAnOption4III)`,
				`isAnOption4IIR = append(isAnOption4IIR, isAnOption4IIIR)`,
				`isAnOption4IIiSize := int64(len(isAnOption4IIIC))`,
				`if err := validate.MinItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), "query", isAnOption4IIiSize, 3); err != nil {`,
				`if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "isAnOption4", i), ii), "query", isAnOption4IIIC); err != nil {`,
				`isAnOption4IR = append(isAnOption4IR, isAnOption4IIR)`,
				`if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "isAnOption4", i), "query", isAnOption4IIC); err != nil {`,
				`o.IsAnOption4 = isAnOption4IR`,
				`if err := o.validateIsAnOption4(formats); err != nil {`,
				`isAnOption4Size := int64(len(o.IsAnOption4))`,
				`if err := validate.MaxItems("isAnOption4", "query", isAnOption4Size, 4); err != nil {`,
				`return errors.Required("notAnOption1", "query")`,
				`notAnOption1IC := swag.SplitByFormat(qvNotAnOption1, "")`,
				`var notAnOption1IR [][]strfmt.DateTime`,
				`for i, notAnOption1IV := range notAnOption1IC {`,
				`notAnOption1IIC := swag.SplitByFormat(notAnOption1IV, "")`,
				`if len(notAnOption1IIC) > 0 {`,
				`var notAnOption1IIR []strfmt.DateTime`,
				`for ii, notAnOption1IIV := range notAnOption1IIC {`,
				`value, err := formats.Parse("date-time", notAnOption1IIV)`,
				`notAnOption1II := *(value.(*strfmt.DateTime))`,
				`if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "notAnOption1", i), ii), "query", "date-time", notAnOption1II.String(), formats); err != nil {`,
				`notAnOption1IIR = append(notAnOption1IIR, notAnOption1II)`,
				`notAnOption1IR = append(notAnOption1IR, notAnOption1IIR)`,
				`o.NotAnOption1 = notAnOption1IR`,
			},
		},
		"5": map[string][]string{
			"serverResponses": []string{
				// expected code lines
				`"github.com/go-openapi/strfmt"`,
				"XIsAnOptionalHeader0 strfmt.DateTime `json:\"x-isAnOptionalHeader0\"`",
				"XIsAnOptionalHeader1 []strfmt.DateTime `json:\"x-isAnOptionalHeader1\"`",
				"XIsAnOptionalHeader2 [][]int32 `json:\"x-isAnOptionalHeader2\"`",
				"XIsAnOptionalHeader3 [][][]strfmt.UUID `json:\"x-isAnOptionalHeader3\"`",
				`xIsAnOptionalHeader0 := o.XIsAnOptionalHeader0.String()`,
				`rw.Header().Set("x-isAnOptionalHeader0", xIsAnOptionalHeader0)`,
				`var xIsAnOptionalHeader1IR []string`,
				`for _, xIsAnOptionalHeader1I := range o.XIsAnOptionalHeader1 {`,
				`xIsAnOptionalHeader1IS := xIsAnOptionalHeader1I.String()`,
				`if xIsAnOptionalHeader1IS != "" {`,
				`xIsAnOptionalHeader1IR = append(xIsAnOptionalHeader1IR, xIsAnOptionalHeader1IS)`,
				`xIsAnOptionalHeader1 := swag.JoinByFormat(xIsAnOptionalHeader1IR, "tsv")`,
				`hv := xIsAnOptionalHeader1[0]`,
				`rw.Header().Set("x-isAnOptionalHeader1", hv)`,
				`var xIsAnOptionalHeader2IR []string`,
				`for _, xIsAnOptionalHeader2I := range o.XIsAnOptionalHeader2 {`,
				`var xIsAnOptionalHeader2IIR []string`,
				`for _, xIsAnOptionalHeader2II := range xIsAnOptionalHeader2I {`,
				`xIsAnOptionalHeader2IIS := swag.FormatInt32(xIsAnOptionalHeader2II)`,
				`if xIsAnOptionalHeader2IIS != "" {`,
				`xIsAnOptionalHeader2IIR = append(xIsAnOptionalHeader2IIR, xIsAnOptionalHeader2IIS)`,
				`xIsAnOptionalHeader2IS := swag.JoinByFormat(xIsAnOptionalHeader2IIR, "pipes")`,
				`xIsAnOptionalHeader2ISs := xIsAnOptionalHeader2IS[0]`,
				`if xIsAnOptionalHeader2ISs != "" {`,
				`xIsAnOptionalHeader2IR = append(xIsAnOptionalHeader2IR, xIsAnOptionalHeader2ISs)`,
				`xIsAnOptionalHeader2 := swag.JoinByFormat(xIsAnOptionalHeader2IR, "")`,
				`hv := xIsAnOptionalHeader2[0]`,
				`rw.Header().Set("x-isAnOptionalHeader2", hv)`,
				`var xIsAnOptionalHeader3IR []string`,
				`for _, xIsAnOptionalHeader3I := range o.XIsAnOptionalHeader3 {`,
				`var xIsAnOptionalHeader3IIR []string`,
				`for _, xIsAnOptionalHeader3II := range xIsAnOptionalHeader3I {`,
				`var xIsAnOptionalHeader3IIIR []string`,
				`for _, xIsAnOptionalHeader3III := range xIsAnOptionalHeader3II {`,
				`xIsAnOptionalHeader3IIIS := xIsAnOptionalHeader3III.String()`,
				`if xIsAnOptionalHeader3IIIS != "" {`,
				`xIsAnOptionalHeader3IIIR = append(xIsAnOptionalHeader3IIIR, xIsAnOptionalHeader3IIIS)`,
				`xIsAnOptionalHeader3IIS := swag.JoinByFormat(xIsAnOptionalHeader3IIIR, "")`,
				`xIsAnOptionalHeader3IISs := xIsAnOptionalHeader3IIS[0]`,
				`if xIsAnOptionalHeader3IISs != "" {`,
				`xIsAnOptionalHeader3IIR = append(xIsAnOptionalHeader3IIR, xIsAnOptionalHeader3IISs)`,
				`xIsAnOptionalHeader3IS := swag.JoinByFormat(xIsAnOptionalHeader3IIR, "pipes")`,
				`xIsAnOptionalHeader3ISs := xIsAnOptionalHeader3IS[0]`,
				`if xIsAnOptionalHeader3ISs != "" {`,
				`xIsAnOptionalHeader3IR = append(xIsAnOptionalHeader3IR, xIsAnOptionalHeader3ISs)`,
				`xIsAnOptionalHeader3 := swag.JoinByFormat(xIsAnOptionalHeader3IR, "")`,
				`hv := xIsAnOptionalHeader3[0]`,
				`rw.Header().Set("x-isAnOptionalHeader3", hv)`,
			},
		},
	}

	for fixtureIndex, fixtureContents := range fixtureConfig {
		fixtureSpec := strings.Join([]string{"fixture-909-", fixtureIndex, ".yaml"}, "")
		gen, err := opBuilder("getOptional", filepath.Join("..", "fixtures", "bugs", "909", fixtureSpec))
		if assert.NoError(err) {
			op, err := gen.MakeOperation()
			if assert.NoError(err) {
				opts := opts()
				for fixtureTemplate, expectedCode := range fixtureContents {
					buf := bytes.NewBuffer(nil)
					err := templates.MustGet(fixtureTemplate).Execute(buf, op)
					if assert.NoError(err, "Expected generation to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
						ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
						if assert.NoError(err, "Expected formatting to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
							res := string(ff)
							for line, codeLine := range expectedCode {
								if !assertInCode(t, strings.TrimSpace(codeLine), res) {
									t.Logf("Code expected did not match for fixture %s at line %d", fixtureSpec, line)
								}
							}
						} else {
							fmt.Println(buf.String())
						}
					}
				}
			}
		}
	}
}

// verifies that validation method is called on body param with $ref
func TestGenParameter_Issue1237(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	assert := assert.New(t)
	fixtureConfig := map[string]map[string][]string{
		"1": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`var body models.Sg`,
				`if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`if err == io.EOF {`,
				`res = append(res, errors.Required("body", "body"))`,
				`} else {`,
				`res = append(res, errors.NewParseError("body", "body", "", err))`,
				`if err := body.Validate(route.Formats); err != nil {`,
			},
		},
	}
	for _, fixtureContents := range fixtureConfig {
		fixtureSpec := strings.Join([]string{"fixture-1237", ".json"}, "")
		gen, err := opBuilder("add sg", filepath.Join("..", "fixtures", "bugs", "1237", fixtureSpec))
		if assert.NoError(err) {
			op, err := gen.MakeOperation()
			if assert.NoError(err) {
				opts := opts()
				for fixtureTemplate, expectedCode := range fixtureContents {
					buf := bytes.NewBuffer(nil)
					err := templates.MustGet(fixtureTemplate).Execute(buf, op)
					if assert.NoError(err, "Expected generation to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
						ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
						if assert.NoError(err, "Expected formatting to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
							res := string(ff)
							for line, codeLine := range expectedCode {
								if !assertInCode(t, strings.TrimSpace(codeLine), res) {
									t.Logf("Code expected did not match for fixture %s at line %d", fixtureSpec, line)
								}
							}
						} else {
							fmt.Println(buf.String())
						}
					}
				}
			}
		}
	}
}

func TestGenParameter_Issue1392(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	assert := assert.New(t)
	fixtureConfig := map[string]map[string][]string{
		"1": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				`func (o *PatchSomeResourceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	var res []error`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close()`,
				`		var body models.BulkUpdateState`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("massUpdate", "body", "", err))`,
				`		} else {`,
				`			if err := body.Validate(route.Formats); err != nil {`,
				`				res = append(res, err)`,
				`			if len(res) == 0 {`,
				`				o.MassUpdate = body`,
				`	if len(res) > 0 {`,
				`		return errors.CompositeValidationError(res...)`,
			},
		},
		"2": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func (o *PostBodybuilder20Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	var res []error`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close()`,
				`		var body []strfmt.URI`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("myObject", "body", "", err))`,
				`		} else {`,
				`			// validate inline body array`,
				`			o.MyObject = body`,
				`			if err := o.validateMyObjectBody(route.Formats); err != nil {`,
				`				res = append(res, err)`,
				`	if len(res) > 0 {`,
				`		return errors.CompositeValidationError(res...)`,
				`func (o *PostBodybuilder20Params) validateMyObjectBody(formats strfmt.Registry) error {`,
				`	// uniqueItems: true`,
				`	if err := validate.UniqueItems("myObject", "body", o.MyObject); err != nil {`,
				`	myObjectIC := o.MyObject`,
				`	var myObjectIR []strfmt.URI`,
				`	for i, myObjectIV := range myObjectIC {`,
				`		myObjectI := myObjectIV`,
				`		if err := validate.FormatOf(fmt.Sprintf("%s.%v", "myObject", i), "body", "uri", myObjectI.String(), formats); err != nil {`,
				`		myObjectIR = append(myObjectIR, myObjectI)`,
				`	o.MyObject = myObjectIR`,
			},
		},
		"3": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func (o *PostBodybuilder26Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	var res []error`,
				`	o.HTTPRequest = r`,
				`	qs := runtime.Values(r.URL.Query())`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close()`,
				`		var body strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("myObject", "body", "", err))`,
				`		} else {`,
				`			// validate inline body`,
				`			o.MyObject = body`,
				`			if err := o.validateMyObjectBody(route.Formats); err != nil {`,
				`				res = append(res, err)`,
				`	qMyquery, qhkMyquery, _ := qs.GetOK("myquery")`,
				`	if err := o.bindMyquery(qMyquery, qhkMyquery, route.Formats); err != nil {`,
				`		res = append(res, err)`,
				`	if len(res) > 0 {`,
				`		return errors.CompositeValidationError(res...)`,
				`	return nil`,
				`func (o *PostBodybuilder26Params) validateMyObjectBody(formats strfmt.Registry) error {`,
				`	if err := validate.Enum("myObject", "body", o.MyObject.String(), []interface{}{"1992-01-01", "2012-01-01"}); err != nil {`,
				`	if err := validate.FormatOf("myObject", "body", "date", o.MyObject.String(), formats); err != nil {`,
			},
		},
		"4": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func (o *PostBodybuilder27Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	var res []error`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close()`,
				`		var body [][]strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("myObject", "body", "", err))`,
				`		} else {`,
				`			o.MyObject = body`,
				`			if err := o.validateMyObjectBody(route.Formats); err != nil {`,
				`				res = append(res, err)`,
				`	if len(res) > 0 {`,
				`		return errors.CompositeValidationError(res...)`,
				`func (o *PostBodybuilder27Params) validateMyObjectBody(formats strfmt.Registry) error {`,
				`	if err := validate.Enum("myObject", "body", o.MyObject, []interface{}{[]interface{}{[]interface{}{"1992-01-01", "2012-01-01"}}}); err != nil {`,
				`		return err`,
				`	myObjectIC := o.MyObject`,
				`	var myObjectIR [][]strfmt.Date`,
				`	for i, myObjectIV := range myObjectIC {`,
				`		myObjectIIC := myObjectIV`,
				`		if len(myObjectIIC) > 0 {`,
				`			var myObjectIIR []strfmt.Date`,
				`			for ii, myObjectIIV := range myObjectIIC {`,
				`				myObjectII := myObjectIIV`,
				`				if err := validate.Enum(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "myObject", i), ii), "", myObjectII.String(), []interface{}{"1992-01-01", "2012-01-01"}); err != nil {`,
				`					return err`,
				`				if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "myObject", i), ii), "", "date", myObjectII.String(), formats); err != nil {`,
				`					return err`,
				`				myObjectIIR = append(myObjectIIR, myObjectII)`,
				`			myObjectIR = append(myObjectIR, myObjectIIR)`,
				// fixed missing enum validation
				`		if err := validate.Enum(fmt.Sprintf("%s.%v", "myObject", i), "body", myObjectIIC, []interface{}{[]interface{}{"1992-01-01", "2012-01-01"}}); err != nil {`,
				`	o.MyObject = myObjectIR`,
			},
		},
		"5": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				`func (o *Bodybuilder23Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	var res []error`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close()`,
				`		var body []models.ASimpleArray`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("myObject", "body", "", err))`,
				`		} else {`,
				`			o.MyObject = body`,
				`			myObjectSize := int64(len(o.MyObject))`,
				`			if err := validate.MinItems("myObject", "body", myObjectSize, 15); err != nil {`,
				`				return err`,
				// changed index
				`			for i := range body {`,
				`				if err := body[i].Validate(route.Formats); err != nil {`,
				`					res = append(res, err)`,
				`					break`,
				`			if len(res) == 0 {`,
				`				o.MyObject = body`,
				`	if len(res) > 0 {`,
				`		return errors.CompositeValidationError(res...)`,
			},
		},
	}

	for fixtureIndex, fixtureContents := range fixtureConfig {
		fixtureSpec := strings.Join([]string{"fixture-1392-", fixtureIndex, ".yaml"}, "")
		// pick selected operation id in fixture
		operationToTest := ""
		switch fixtureIndex {
		case "1":
			operationToTest = "PatchSomeResource"
		case "2":
			operationToTest = "PostBodybuilder20"
		case "3":
			operationToTest = "PostBodybuilder26"
		case "4":
			operationToTest = "PostBodybuilder27"
		case "5":
			operationToTest = "Bodybuilder23"
		}
		gen, err := opBuilder(operationToTest, filepath.Join("..", "fixtures", "bugs", "1392", fixtureSpec))
		if assert.NoError(err) {
			op, err := gen.MakeOperation()
			if assert.NoError(err) {
				opts := opts()
				for fixtureTemplate, expectedCode := range fixtureContents {
					buf := bytes.NewBuffer(nil)
					err := templates.MustGet(fixtureTemplate).Execute(buf, op)
					if assert.NoError(err, "Expected generation to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
						ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
						if assert.NoError(err, "Expected formatting to go well on %s with template %s", fixtureSpec, fixtureTemplate) {
							res := string(ff)
							for line, codeLine := range expectedCode {
								if !assertInCode(t, strings.TrimSpace(codeLine), res) {
									t.Logf("Code expected did not match for fixture %s at line %d", fixtureSpec, line)
								}
							}
						} else {
							fmt.Println(buf.String())
						}
					}
				}
			}
		}
	}
}

func TestGenParameter_Issue1513(t *testing.T) {
	assert := assert.New(t)
	var assertion = `r.SetBodyParam(o.Something)`
	gen, err := opBuilderWithFlatten("put-enum", "../fixtures/bugs/1513/enums.yaml")
	if assert.NoError(err) {
		op, err := gen.MakeOperation()
		if assert.NoError(err) {
			buf := bytes.NewBuffer(nil)
			opts := opts()
			err := templates.MustGet("clientParameter").Execute(buf, op)
			if assert.NoError(err) {
				ff, err := opts.LanguageOpts.FormatContent("move_clusters_parameters.go", buf.Bytes())
				if assert.NoError(err) {
					res := string(ff)
					assertInCode(t, assertion, res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

// Body param validation on empty objects
func TestGenParameter_Issue1536(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	// testing fixture-1536.yaml with flatten
	// param body with array of empty objects

	fixtureConfig := map[string]map[string][]string{

		// load expectations for parameters in operation get_interface_parameters.go
		"getInterface": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetInterfaceParams() GetInterfaceParams {`,
				`	return GetInterfaceParams{`,
				`type GetInterfaceParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	Generic interface{`,
				`func (o *GetInterfaceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("generic", "body", "", err)`,
				`		} else {`,
				`			o.Generic = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_map_slice_parameters.go
		"getMapSlice": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapSliceParams() GetMapSliceParams {`,
				`	return GetMapSliceParams{`,
				`type GetMapSliceParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericMapSlice []map[string]models.ModelInterface`,
				`func (o *GetMapSliceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []map[string]models.ModelInterface`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericMapSlice", "body", "", err)`,
				`		} else {`,
				`			o.GenericMapSlice = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_with_validations_parameters.go
		"getNestedWithValidations": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedWithValidationsParams() GetNestedWithValidationsParams {`,
				`	return GetNestedWithValidationsParams{`,
				`type GetNestedWithValidationsParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericNestedWithValidations [][][][]interface{`,
				`func (o *GetNestedWithValidationsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericNestedWithValidations", "body", "", err)`,
				`		} else {`,
				`			o.GenericNestedWithValidations = body`,
				`			if err := o.validateGenericNestedWithValidationsBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedWithValidationsParams) validateGenericNestedWithValidationsBody(formats strfmt.Registry) error {`,
				`	genericNestedWithValidationsIC := o.GenericNestedWithValidations`,
				`	var genericNestedWithValidationsIR [][][][]interface{`,
				`	for i, genericNestedWithValidationsIV := range genericNestedWithValidationsIC {`,
				`		genericNestedWithValidationsIIC := genericNestedWithValidationsIV`,
				`		if len(genericNestedWithValidationsIIC) > 0 {`,
				`			var genericNestedWithValidationsIIR [][][]interface{`,
				`			for ii, genericNestedWithValidationsIIV := range genericNestedWithValidationsIIC {`,
				`				genericNestedWithValidationsIIIC := genericNestedWithValidationsIIV`,
				`				if len(genericNestedWithValidationsIIIC) > 0 {`,
				`					var genericNestedWithValidationsIIIR [][]interface{`,
				`					for iii, genericNestedWithValidationsIIIV := range genericNestedWithValidationsIIIC {`,
				`						genericNestedWithValidationsIIIIC := genericNestedWithValidationsIIIV`,
				`						if len(genericNestedWithValidationsIIIIC) > 0 {`,
				`							var genericNestedWithValidationsIIIIR []interface{`,
				`							for _, genericNestedWithValidationsIIIIV := range genericNestedWithValidationsIIIIC {`,
				`								genericNestedWithValidationsIIII := genericNestedWithValidationsIIIIV`,
				`								genericNestedWithValidationsIIIIR = append(genericNestedWithValidationsIIIIR, genericNestedWithValidationsIIII`,
				`							genericNestedWithValidationsIIIR = append(genericNestedWithValidationsIIIR, genericNestedWithValidationsIIIIR`,
				`						genericNestedWithValidationsIiiiiiSize := int64(len(genericNestedWithValidationsIIIIC)`,
				`						if err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "genericNestedWithValidations", i), ii), iii), "", genericNestedWithValidationsIiiiiiSize, 10); err != nil {`,
				`					genericNestedWithValidationsIIR = append(genericNestedWithValidationsIIR, genericNestedWithValidationsIIIR`,
				`			genericNestedWithValidationsIR = append(genericNestedWithValidationsIR, genericNestedWithValidationsIIR`,
				`	o.GenericNestedWithValidations = genericNestedWithValidationsIR`,
			},
		},

		// load expectations for parameters in operation get_another_interface_parameters.go
		"getAnotherInterface": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetAnotherInterfaceParams() GetAnotherInterfaceParams {`,
				`	return GetAnotherInterfaceParams{`,
				`type GetAnotherInterfaceParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	AnotherGeneric interface{`,
				`func (o *GetAnotherInterfaceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("anotherGeneric", "body", "", err)`,
				`		} else {`,
				`			o.AnotherGeneric = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_required_parameters.go
		"getNestedRequired": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedRequiredParams() GetNestedRequiredParams {`,
				`	return GetNestedRequiredParams{`,
				`type GetNestedRequiredParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ObjectNestedRequired [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`func (o *GetNestedRequiredParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("objectNestedRequired", "body", "", err)`,
				`		} else {`,
				`			o.ObjectNestedRequired = body`,
				`			if err := o.validateObjectNestedRequiredBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedRequiredParams) validateObjectNestedRequiredBody(formats strfmt.Registry) error {`,
				`	objectNestedRequiredIC := o.ObjectNestedRequired`,
				`	var objectNestedRequiredIR [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`	for i, objectNestedRequiredIV := range objectNestedRequiredIC {`,
				`		objectNestedRequiredIIC := objectNestedRequiredIV`,
				`		if len(objectNestedRequiredIIC) > 0 {`,
				`			var objectNestedRequiredIIR [][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`			for ii, objectNestedRequiredIIV := range objectNestedRequiredIIC {`,
				`				objectNestedRequiredIIIC := objectNestedRequiredIIV`,
				`				if len(objectNestedRequiredIIIC) > 0 {`,
				`					var objectNestedRequiredIIIR [][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`					for iii, objectNestedRequiredIIIV := range objectNestedRequiredIIIC {`,
				`						objectNestedRequiredIIIIC := objectNestedRequiredIIIV`,
				`						if len(objectNestedRequiredIIIIC) > 0 {`,
				`							var objectNestedRequiredIIIIR []*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`							for iiii, objectNestedRequiredIIIIV := range objectNestedRequiredIIIIC {`,
				`								objectNestedRequiredIIII := objectNestedRequiredIIIIV`,
				`								if err := objectNestedRequiredIIII.Validate(formats); err != nil {`,
				`									if ve, ok := err.(*errors.Validation); ok {`,
				`										return ve.ValidateName(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "objectNestedRequired", i), ii), iii), iiii)`,
				`								objectNestedRequiredIIIIR = append(objectNestedRequiredIIIIR, objectNestedRequiredIIII`,
				`							objectNestedRequiredIIIR = append(objectNestedRequiredIIIR, objectNestedRequiredIIIIR`,
				`					objectNestedRequiredIIR = append(objectNestedRequiredIIR, objectNestedRequiredIIIR`,
				`			objectNestedRequiredIR = append(objectNestedRequiredIR, objectNestedRequiredIIR`,
				`	o.ObjectNestedRequired = objectNestedRequiredIR`,
			},
		},

		// load expectations for parameters in operation get_records_max_parameters.go
		"getRecordsMax": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetRecordsMaxParams() GetRecordsMaxParams {`,
				`	return GetRecordsMaxParams{`,
				`type GetRecordsMaxParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MaxRecords []interface{`,
				`func (o *GetRecordsMaxParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("maxRecords", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("maxRecords", "body", "", err)`,
				`		} else {`,
				`			o.MaxRecords = body`,
				`			if err := o.validateMaxRecordsBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("maxRecords", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetRecordsMaxParams) validateMaxRecordsBody(formats strfmt.Registry) error {`,
				`	maxRecordsSize := int64(len(o.MaxRecords)`,
				`	if err := validate.MaxItems("maxRecords", "body", maxRecordsSize, 10); err != nil {`,
			},
		},

		// load expectations for parameters in operation get_records_parameters.go
		"getRecords": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetRecordsParams() GetRecordsParams {`,
				`	return GetRecordsParams{`,
				`type GetRecordsParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	Records []interface{`,
				`func (o *GetRecordsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("records", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("records", "body", "", err)`,
				`		} else {`,
				`			o.Records = body`,
				// fixed: no validation has to be carried on
				`	} else {`,
				`		res = append(res, errors.Required("records", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},
		// load expectations for parameters in operation get_records_non_required_parameters.go
		"getRecordsNonRequired": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetRecordsNonRequiredParams() GetRecordsNonRequiredParams {`,
				`	return GetRecordsNonRequiredParams{`,
				`type GetRecordsNonRequiredParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	RecordsNonRequired []interface{`,
				`func (o *GetRecordsNonRequiredParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("recordsNonRequired", "body", "", err)`,
				`		} else {`,
				`			o.RecordsNonRequired = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},
		// load expectations for parameters in operation get_map_parameters.go
		"getMap": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapParams() GetMapParams {`,
				`	return GetMapParams{`,
				`type GetMapParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericMap map[string]models.ModelInterface`,
				`func (o *GetMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]models.ModelInterface`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericMap", "body", "", err)`,
				`		} else {`,
				`			o.GenericMap = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_slice_map_parameters.go
		"getSliceMap": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetSliceMapParams() GetSliceMapParams {`,
				`	return GetSliceMapParams{`,
				`type GetSliceMapParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericSliceMap map[string][]models.ModelInterface`,
				`func (o *GetSliceMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]models.ModelInterface`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericSliceMap", "body", "", err)`,
				`		} else {`,
				`			o.GenericSliceMap = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_parameters.go
		"getNested": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedParams() GetNestedParams {`,
				`	return GetNestedParams{`,
				`type GetNestedParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericNested [][][][]interface{`,
				`func (o *GetNestedParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericNested", "body", "", err)`,
				`		} else {`,
				`			o.GenericNested = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},
	}

	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1536", "fixture-1536.yaml"))
}

func TestGenParameter_Issue1536_Arrays(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	fixtureConfig := map[string]map[string][]string{

		// load expectations for parameters in operation get_nested_with_validations_parameters.go
		"getNestedWithValidations": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedWithValidationsParams() GetNestedWithValidationsParams {`,
				`	return GetNestedWithValidationsParams{`,
				`type GetNestedWithValidationsParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericNestedWithValidations [][][][]interface{`,
				`func (o *GetNestedWithValidationsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericNestedWithValidations", "body", "", err)`,
				`		} else {`,
				`			o.GenericNestedWithValidations = body`,
				`			if err := o.validateGenericNestedWithValidationsBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedWithValidationsParams) validateGenericNestedWithValidationsBody(formats strfmt.Registry) error {`,
				`	genericNestedWithValidationsIC := o.GenericNestedWithValidations`,
				`	var genericNestedWithValidationsIR [][][][]interface{`,
				`	for i, genericNestedWithValidationsIV := range genericNestedWithValidationsIC {`,
				`		genericNestedWithValidationsIIC := genericNestedWithValidationsIV`,
				`		if len(genericNestedWithValidationsIIC) > 0 {`,
				`			var genericNestedWithValidationsIIR [][][]interface{`,
				`			for ii, genericNestedWithValidationsIIV := range genericNestedWithValidationsIIC {`,
				`				genericNestedWithValidationsIIIC := genericNestedWithValidationsIIV`,
				`				if len(genericNestedWithValidationsIIIC) > 0 {`,
				`					var genericNestedWithValidationsIIIR [][]interface{`,
				`					for iii, genericNestedWithValidationsIIIV := range genericNestedWithValidationsIIIC {`,
				`						genericNestedWithValidationsIIIIC := genericNestedWithValidationsIIIV`,
				`						if len(genericNestedWithValidationsIIIIC) > 0 {`,
				`							var genericNestedWithValidationsIIIIR []interface{`,
				`							for _, genericNestedWithValidationsIIIIV := range genericNestedWithValidationsIIIIC {`,
				`								genericNestedWithValidationsIIII := genericNestedWithValidationsIIIIV`,
				`								genericNestedWithValidationsIIIIR = append(genericNestedWithValidationsIIIIR, genericNestedWithValidationsIIII`,
				`							genericNestedWithValidationsIIIR = append(genericNestedWithValidationsIIIR, genericNestedWithValidationsIIIIR`,
				`						genericNestedWithValidationsIiiiiiSize := int64(len(genericNestedWithValidationsIIIIC)`,
				`						if err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "genericNestedWithValidations", i), ii), iii), "", genericNestedWithValidationsIiiiiiSize, 10); err != nil {`,
				`					genericNestedWithValidationsIIR = append(genericNestedWithValidationsIIR, genericNestedWithValidationsIIIR`,
				`			genericNestedWithValidationsIR = append(genericNestedWithValidationsIR, genericNestedWithValidationsIIR`,
				`	o.GenericNestedWithValidations = genericNestedWithValidationsIR`,
			},
		},

		// load expectations for parameters in operation get_nested_required_parameters.go
		"getNestedRequired": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedRequiredParams() GetNestedRequiredParams {`,
				`	return GetNestedRequiredParams{`,
				`type GetNestedRequiredParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ObjectNestedRequired [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`func (o *GetNestedRequiredParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("objectNestedRequired", "body", "", err)`,
				`		} else {`,
				`			o.ObjectNestedRequired = body`,
				`			if err := o.validateObjectNestedRequiredBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedRequiredParams) validateObjectNestedRequiredBody(formats strfmt.Registry) error {`,
				`	objectNestedRequiredIC := o.ObjectNestedRequired`,
				`	var objectNestedRequiredIR [][][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`	for i, objectNestedRequiredIV := range objectNestedRequiredIC {`,
				`		objectNestedRequiredIIC := objectNestedRequiredIV`,
				`		if len(objectNestedRequiredIIC) > 0 {`,
				`			var objectNestedRequiredIIR [][][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`			for ii, objectNestedRequiredIIV := range objectNestedRequiredIIC {`,
				`				objectNestedRequiredIIIC := objectNestedRequiredIIV`,
				`				if len(objectNestedRequiredIIIC) > 0 {`,
				`					var objectNestedRequiredIIIR [][]*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`					for iii, objectNestedRequiredIIIV := range objectNestedRequiredIIIC {`,
				`						objectNestedRequiredIIIIC := objectNestedRequiredIIIV`,
				`						if len(objectNestedRequiredIIIIC) > 0 {`,
				`							var objectNestedRequiredIIIIR []*models.GetNestedRequiredParamsBodyItemsItemsItemsItems`,
				`							for iiii, objectNestedRequiredIIIIV := range objectNestedRequiredIIIIC {`,
				`								if objectNestedRequiredIIIIV == nil {`,
				`									if err := objectNestedRequiredIIII.Validate(formats); err != nil {`,
				`										if ve, ok := err.(*errors.Validation); ok {`,
				`											return ve.ValidateName(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "objectNestedRequired", i), ii), iii), iiii)`,
				`								objectNestedRequiredIIIIR = append(objectNestedRequiredIIIIR, objectNestedRequiredIIII`,
				`							objectNestedRequiredIIIR = append(objectNestedRequiredIIIR, objectNestedRequiredIIIIR`,
				`					objectNestedRequiredIIR = append(objectNestedRequiredIIR, objectNestedRequiredIIIR`,
				`			objectNestedRequiredIR = append(objectNestedRequiredIR, objectNestedRequiredIIR`,
				`	o.ObjectNestedRequired = objectNestedRequiredIR`,
			},
		},

		// load expectations for parameters in operation get_simple_array_with_slice_validation_parameters.go
		"getSimpleArrayWithSliceValidation": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetSimpleArrayWithSliceValidationParams() GetSimpleArrayWithSliceValidationParams {`,
				`	return GetSimpleArrayWithSliceValidationParams{`,
				`type GetSimpleArrayWithSliceValidationParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	SimpleArrayWithSliceValidation []int64`,
				`func (o *GetSimpleArrayWithSliceValidationParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("simpleArrayWithSliceValidation", "body", "", err)`,
				`		} else {`,
				`			o.SimpleArrayWithSliceValidation = body`,
				`			if err := o.validateSimpleArrayWithSliceValidationBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetSimpleArrayWithSliceValidationParams) validateSimpleArrayWithSliceValidationBody(formats strfmt.Registry) error {`,
				`	if err := validate.Enum("simpleArrayWithSliceValidation", "body", o.SimpleArrayWithSliceValidation, []interface{}{[]interface{}{1, 2, 3}, []interface{}{4, 5, 6}}); err != nil {`,
			},
		},

		// load expectations for parameters in operation get_simple_parameters.go
		"getSimple": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetSimpleParams() GetSimpleParams {`,
				`	return GetSimpleParams{`,
				`type GetSimpleParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	SimpleBody *models.GetSimpleParamsBody`,
				`func (o *GetSimpleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body models.GetSimpleParamsBody`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("simpleBody", "body", "", err)`,
				`		} else {`,
				`			if err := body.Validate(route.Formats); err != nil {`,
				`			if len(res) == 0 {`,
				`				o.SimpleBody = &body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_simple_array_with_validation_parameters.go
		"getSimpleArrayWithValidation": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetSimpleArrayWithValidationParams() GetSimpleArrayWithValidationParams {`,
				`	return GetSimpleArrayWithValidationParams{`,
				`type GetSimpleArrayWithValidationParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	SimpleArrayWithValidation []int64`,
				`func (o *GetSimpleArrayWithValidationParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("simpleArrayWithValidation", "body", "", err)`,
				`		} else {`,
				`			o.SimpleArrayWithValidation = body`,
				`			if err := o.validateSimpleArrayWithValidationBody(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetSimpleArrayWithValidationParams) validateSimpleArrayWithValidationBody(formats strfmt.Registry) error {`,
				`	simpleArrayWithValidationIC := o.SimpleArrayWithValidation`,
				`	var simpleArrayWithValidationIR []int64`,
				`	for i, simpleArrayWithValidationIV := range simpleArrayWithValidationIC {`,
				`		simpleArrayWithValidationI := simpleArrayWithValidationIV`,
				`		if err := validate.MaximumInt(fmt.Sprintf("%s.%v", "simpleArrayWithValidation", i), "body", int64(simpleArrayWithValidationI), 12, false); err != nil {`,
				`		simpleArrayWithValidationIR = append(simpleArrayWithValidationIR, simpleArrayWithValidationI`,
				`	o.SimpleArrayWithValidation = simpleArrayWithValidationIR`,
			},
		},

		// load expectations for parameters in operation get_nested_parameters.go
		"getNested": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedParams() GetNestedParams {`,
				`	return GetNestedParams{`,
				`type GetNestedParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	GenericNested [][][][]interface{`,
				`func (o *GetNestedParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][][]interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("genericNested", "body", "", err)`,
				`		} else {`,
				`			o.GenericNested = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_simple_array_parameters.go
		"getSimpleArray": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetSimpleArrayParams() GetSimpleArrayParams {`,
				`	return GetSimpleArrayParams{`,
				`type GetSimpleArrayParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	SimpleArray []int64`,
				`func (o *GetSimpleArrayParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("simpleArray", "body", "", err)`,
				`		} else {`,
				`			o.SimpleArray = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},
	}
	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1536", "fixture-1536-2.yaml"))
}

func TestGenParameter_Issue1536_Maps(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	fixtureConfig := map[string]map[string][]string{

		// load expectations for parameters in operation get_map_interface_parameters.go
		"getMapInterface": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapInterfaceParams() GetMapInterfaceParams {`,
				`	return GetMapInterfaceParams{`,
				`type GetMapInterfaceParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfInterface map[string]models.ModelInterface`,
				`func (o *GetMapInterfaceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]models.ModelInterface`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfInterface", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfInterface", "body", "", err)`,
				`		} else {`,
				`			o.MapOfInterface = body`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfInterface", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_array_of_interface_parameters.go
		"getArrayOfInterface": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetArrayOfInterfaceParams() GetArrayOfInterfaceParams {`,
				`	return GetArrayOfInterfaceParams{`,
				`type GetArrayOfInterfaceParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ArrayOfInterface []interface{`,
				`func (o *GetArrayOfInterfaceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []interface{`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("arrayOfInterface", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("arrayOfInterface", "body", "", err)`,
				`		} else {`,
				`			o.ArrayOfInterface = body`,
				`	} else {`,
				`		res = append(res, errors.Required("arrayOfInterface", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_map_array_with_max_parameters.go
		"getMapArrayWithMax": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapArrayWithMaxParams() GetMapArrayWithMaxParams {`,
				`	return GetMapArrayWithMaxParams{`,
				`type GetMapArrayWithMaxParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfArrayWithMax map[string]models.ModelArrayWithMax`,
				`func (o *GetMapArrayWithMaxParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]models.ModelArrayWithMax`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfArrayWithMax", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfArrayWithMax", "body", "", err)`,
				`		} else {`,
				`			for k := range body {`,
				`				if err := o.MapOfArrayWithMax[k].Validate(route.Formats); err != nil {`,
				`					break`,
				`			if len(res) == 0 {`,
				`				o.MapOfArrayWithMax = body`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfArrayWithMax", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_array_nested_simple_parameters.go
		"getArrayNestedSimple": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetArrayNestedSimpleParams() GetArrayNestedSimpleParams {`,
				`	return GetArrayNestedSimpleParams{`,
				`type GetArrayNestedSimpleParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ArrayOfarraySimple [][]string`,
				`func (o *GetArrayNestedSimpleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][]string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("arrayOfarraySimple", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("arrayOfarraySimple", "body", "", err)`,
				`		} else {`,
				`			o.ArrayOfarraySimple = body`,
				`			if err := o.validateArrayOfarraySimpleBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("arrayOfarraySimple", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetArrayNestedSimpleParams) validateArrayOfarraySimpleBody(formats strfmt.Registry) error {`,
				`	arrayOfarraySimpleIC := o.ArrayOfarraySimple`,
				`	var arrayOfarraySimpleIR [][]string`,
				`	for i, arrayOfarraySimpleIV := range arrayOfarraySimpleIC {`,
				`		arrayOfarraySimpleIIC := arrayOfarraySimpleIV`,
				`		if len(arrayOfarraySimpleIIC) > 0 {`,
				`			var arrayOfarraySimpleIIR []string`,
				`			for ii, arrayOfarraySimpleIIV := range arrayOfarraySimpleIIC {`,
				`				arrayOfarraySimpleII := arrayOfarraySimpleIIV`,
				`				if err := validate.MaxLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "arrayOfarraySimple", i), ii), "", arrayOfarraySimpleII, 100); err != nil {`,
				`				arrayOfarraySimpleIIR = append(arrayOfarraySimpleIIR, arrayOfarraySimpleII`,
				`			arrayOfarraySimpleIR = append(arrayOfarraySimpleIR, arrayOfarraySimpleIIR`,
				`	o.ArrayOfarraySimple = arrayOfarraySimpleIR`,
			},
		},

		// load expectations for parameters in operation get_map_of_format_parameters.go
		"getMapOfFormat": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapOfFormatParams() GetMapOfFormatParams {`,
				`	return GetMapOfFormatParams{`,
				`type GetMapOfFormatParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfFormat map[string]strfmt.UUID`,
				`func (o *GetMapOfFormatParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]strfmt.UUID`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfFormat", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfFormat", "body", "", err)`,
				`		} else {`,
				`			o.MapOfFormat = body`,
				`			if err := o.validateMapOfFormatBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfFormat", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapOfFormatParams) validateMapOfFormatBody(formats strfmt.Registry) error {`,
				`	mapOfFormatIC := o.MapOfFormat`,
				`	mapOfFormatIR := make(map[string]strfmt.UUID, len(mapOfFormatIC)`,
				`	for k, mapOfFormatIV := range mapOfFormatIC {`,
				`		mapOfFormatI := mapOfFormatIV`,
				`		if err := validate.FormatOf(fmt.Sprintf("%s.%v", "mapOfFormat", k), "body", "uuid", mapOfFormatI.String(), formats); err != nil {`,
				`		mapOfFormatIR[k] = mapOfFormatI`,
				`	o.MapOfFormat = mapOfFormatIR`,
			},
		},

		// load expectations for parameters in operation get_array_of_map_parameters.go
		"getArrayOfMap": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetArrayOfMapParams() GetArrayOfMapParams {`,
				`	return GetArrayOfMapParams{`,
				`type GetArrayOfMapParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ArrayOfMap []map[string][]int32`,
				`func (o *GetArrayOfMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []map[string][]int32`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("arrayOfMap", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("arrayOfMap", "body", "", err)`,
				`		} else {`,
				`			o.ArrayOfMap = body`,
				`			if err := o.validateArrayOfMapBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("arrayOfMap", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetArrayOfMapParams) validateArrayOfMapBody(formats strfmt.Registry) error {`,
				`	arrayOfMapSize := int64(len(o.ArrayOfMap)`,
				`	if err := validate.MaxItems("arrayOfMap", "body", arrayOfMapSize, 50); err != nil {`,
				`	arrayOfMapIC := o.ArrayOfMap`,
				`	var arrayOfMapIR []map[string][]int32`,
				`	for i, arrayOfMapIV := range arrayOfMapIC {`,
				`		arrayOfMapIIC := arrayOfMapIV`,
				`		arrayOfMapIIR := make(map[string][]int32, len(arrayOfMapIIC)`,
				`		for kk, arrayOfMapIIV := range arrayOfMapIIC {`,
				`			arrayOfMapIIIC := arrayOfMapIIV`,
				`			var arrayOfMapIIIR []int32`,
				`			for iii, arrayOfMapIIIV := range arrayOfMapIIIC {`,
				`				arrayOfMapIII := arrayOfMapIIIV`,
				`				if err := validate.MaximumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "arrayOfMap", i), kk), iii), "", int64(arrayOfMapIII), 100, false); err != nil {`,
				`				arrayOfMapIIIR = append(arrayOfMapIIIR, arrayOfMapIII`,
				`			arrayOfMapIIR[kk] = arrayOfMapIIIR`,
				`		arrayOfMapIR = append(arrayOfMapIR, arrayOfMapIIR`,
				`	o.ArrayOfMap = arrayOfMapIR`,
			},
		},

		// load expectations for parameters in operation get_map_anon_array_with_x_nullable_parameters.go
		"getMapAnonArrayWithXNullable": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapAnonArrayWithXNullableParams() GetMapAnonArrayWithXNullableParams {`,
				`	return GetMapAnonArrayWithXNullableParams{`,
				`type GetMapAnonArrayWithXNullableParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfAnonArrayWithXNullable map[string][]*int64`,
				`func (o *GetMapAnonArrayWithXNullableParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]*int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfAnonArrayWithXNullable", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfAnonArrayWithXNullable", "body", "", err)`,
				`		} else {`,
				`			o.MapOfAnonArrayWithXNullable = body`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfAnonArrayWithXNullable", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_array_nested_parameters.go
		"getArrayNested": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetArrayNestedParams() GetArrayNestedParams {`,
				`	return GetArrayNestedParams{`,
				`type GetArrayNestedParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ArrayOfarray [][]*models.ModelObject`,
				`func (o *GetArrayNestedParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][]*models.ModelObject`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("arrayOfarray", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("arrayOfarray", "body", "", err)`,
				`		} else {`,
				`			o.ArrayOfarray = body`,
				`			if err := o.validateArrayOfarrayBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("arrayOfarray", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetArrayNestedParams) validateArrayOfarrayBody(formats strfmt.Registry) error {`,
				`	arrayOfarrayIC := o.ArrayOfarray`,
				`	var arrayOfarrayIR [][]*models.ModelObject`,
				`	for i, arrayOfarrayIV := range arrayOfarrayIC {`,
				`		arrayOfarrayIIC := arrayOfarrayIV`,
				`		if len(arrayOfarrayIIC) > 0 {`,
				`			var arrayOfarrayIIR []*models.ModelObject`,
				`			for ii, arrayOfarrayIIV := range arrayOfarrayIIC {`,
				`				if arrayOfarrayIIV == nil {`,
				`					if err := arrayOfarrayII.Validate(formats); err != nil {`,
				`						if ve, ok := err.(*errors.Validation); ok {`,
				`							return ve.ValidateName(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "arrayOfarray", i), ii)`,
				`				arrayOfarrayIIR = append(arrayOfarrayIIR, arrayOfarrayII`,
				`			arrayOfarrayIR = append(arrayOfarrayIR, arrayOfarrayIIR`,
				`	o.ArrayOfarray = arrayOfarrayIR`,
			},
		},

		// load expectations for parameters in operation get_map_array_parameters.go
		"getMapArray": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapArrayParams() GetMapArrayParams {`,
				`	return GetMapArrayParams{`,
				`type GetMapArrayParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfArray models.GetMapArrayParamsBody`,
				`func (o *GetMapArrayParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body models.GetMapArrayParamsBody`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfArray", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfArray", "body", "", err)`,
				`		} else {`,
				`			if err := body.Validate(route.Formats); err != nil {`,
				`			if len(res) == 0 {`,
				`				o.MapOfArray = body`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfArray", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_map_anon_array_with_nullable_parameters.go
		"getMapAnonArrayWithNullable": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapAnonArrayWithNullableParams() GetMapAnonArrayWithNullableParams {`,
				`	return GetMapAnonArrayWithNullableParams{`,
				`type GetMapAnonArrayWithNullableParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfAnonArrayWithNullable map[string][]*int64`,
				`func (o *GetMapAnonArrayWithNullableParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]*int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfAnonArrayWithNullable", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfAnonArrayWithNullable", "body", "", err)`,
				`		} else {`,
				`			o.MapOfAnonArrayWithNullable = body`,
				`			if err := o.validateMapOfAnonArrayWithNullableBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfAnonArrayWithNullable", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapAnonArrayWithNullableParams) validateMapOfAnonArrayWithNullableBody(formats strfmt.Registry) error {`,
				`	mapOfAnonArrayWithNullableIC := o.MapOfAnonArrayWithNullable`,
				`	mapOfAnonArrayWithNullableIR := make(map[string][]*int64, len(mapOfAnonArrayWithNullableIC)`,
				`	for k, mapOfAnonArrayWithNullableIV := range mapOfAnonArrayWithNullableIC {`,
				`		mapOfAnonArrayWithNullableIIC := mapOfAnonArrayWithNullableIV`,
				`		var mapOfAnonArrayWithNullableIIR []*int64`,
				`		for ii, mapOfAnonArrayWithNullableIIV := range mapOfAnonArrayWithNullableIIC {`,
				`			mapOfAnonArrayWithNullableII := mapOfAnonArrayWithNullableIIV`,
				`			if err := validate.MinimumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "mapOfAnonArrayWithNullable", k), ii), "", int64(*mapOfAnonArrayWithNullableII), 0, false); err != nil {`,
				`			mapOfAnonArrayWithNullableIIR = append(mapOfAnonArrayWithNullableIIR, mapOfAnonArrayWithNullableII`,
				`		mapOfAnonArrayWithNullableIR[k] = mapOfAnonArrayWithNullableIIR`,
				`	o.MapOfAnonArrayWithNullable = mapOfAnonArrayWithNullableIR`,
			},
		},

		// load expectations for parameters in operation get_map_of_anon_array_parameters.go
		"getMapOfAnonArray": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapOfAnonArrayParams() GetMapOfAnonArrayParams {`,
				`	return GetMapOfAnonArrayParams{`,
				`type GetMapOfAnonArrayParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfAnonArray map[string][]int64`,
				`func (o *GetMapOfAnonArrayParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfAnonArray", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfAnonArray", "body", "", err)`,
				`		} else {`,
				`			o.MapOfAnonArray = body`,
				`			if err := o.validateMapOfAnonArrayBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfAnonArray", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapOfAnonArrayParams) validateMapOfAnonArrayBody(formats strfmt.Registry) error {`,
				`	mapOfAnonArrayIC := o.MapOfAnonArray`,
				`	mapOfAnonArrayIR := make(map[string][]int64, len(mapOfAnonArrayIC)`,
				`	for k, mapOfAnonArrayIV := range mapOfAnonArrayIC {`,
				`		mapOfAnonArrayIIC := mapOfAnonArrayIV`,
				`		var mapOfAnonArrayIIR []int64`,
				`		for ii, mapOfAnonArrayIIV := range mapOfAnonArrayIIC {`,
				`			mapOfAnonArrayII := mapOfAnonArrayIIV`,
				`			if err := validate.MaximumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "mapOfAnonArray", k), ii), "", int64(mapOfAnonArrayII), 100, false); err != nil {`,
				`			mapOfAnonArrayIIR = append(mapOfAnonArrayIIR, mapOfAnonArrayII`,
				`		mapOfAnonArrayIR[k] = mapOfAnonArrayIIR`,
				`	o.MapOfAnonArray = mapOfAnonArrayIR`,
			},
		},

		// load expectations for parameters in operation get_map_of_anon_map_parameters.go
		"getMapOfAnonMap": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapOfAnonMapParams() GetMapOfAnonMapParams {`,
				`	return GetMapOfAnonMapParams{`,
				`type GetMapOfAnonMapParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfAnonMap map[string]map[string][]int64`,
				`func (o *GetMapOfAnonMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string][]int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfAnonMap", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfAnonMap", "body", "", err)`,
				`		} else {`,
				`			o.MapOfAnonMap = body`,
				`			if err := o.validateMapOfAnonMapBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfAnonMap", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapOfAnonMapParams) validateMapOfAnonMapBody(formats strfmt.Registry) error {`,
				`	mapOfAnonMapIC := o.MapOfAnonMap`,
				`	mapOfAnonMapIR := make(map[string]map[string][]int64, len(mapOfAnonMapIC)`,
				`	for k, mapOfAnonMapIV := range mapOfAnonMapIC {`,
				`		mapOfAnonMapIIC := mapOfAnonMapIV`,
				`		mapOfAnonMapIIR := make(map[string][]int64, len(mapOfAnonMapIIC)`,
				`		for kk, mapOfAnonMapIIV := range mapOfAnonMapIIC {`,
				`			mapOfAnonMapIIIC := mapOfAnonMapIIV`,
				`			var mapOfAnonMapIIIR []int64`,
				`			for iii, mapOfAnonMapIIIV := range mapOfAnonMapIIIC {`,
				`				mapOfAnonMapIII := mapOfAnonMapIIIV`,
				`				if err := validate.MaximumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "mapOfAnonMap", k), kk), iii), "", int64(mapOfAnonMapIII), 100, false); err != nil {`,
				`				mapOfAnonMapIIIR = append(mapOfAnonMapIIIR, mapOfAnonMapIII`,
				`			mapOfAnonMapIIR[kk] = mapOfAnonMapIIIR`,
				`		mapOfAnonMapIR[k] = mapOfAnonMapIIR`,
				`	o.MapOfAnonMap = mapOfAnonMapIR`,
			},
		},

		// load expectations for parameters in operation get_map_anon_array_parameters.go
		"getMapAnonArray": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapAnonArrayParams() GetMapAnonArrayParams {`,
				`	return GetMapAnonArrayParams{`,
				`type GetMapAnonArrayParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfAnonArray map[string][]int64`,
				`func (o *GetMapAnonArrayParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfAnonArray", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfAnonArray", "body", "", err)`,
				`		} else {`,
				`			o.MapOfAnonArray = body`,
				`			if err := o.validateMapOfAnonArrayBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfAnonArray", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapAnonArrayParams) validateMapOfAnonArrayBody(formats strfmt.Registry) error {`,
				`	mapOfAnonArrayIC := o.MapOfAnonArray`,
				`	mapOfAnonArrayIR := make(map[string][]int64, len(mapOfAnonArrayIC)`,
				`	for k, mapOfAnonArrayIV := range mapOfAnonArrayIC {`,
				`		mapOfAnonArrayIIC := mapOfAnonArrayIV`,
				`		var mapOfAnonArrayIIR []int64`,
				`		for ii, mapOfAnonArrayIIV := range mapOfAnonArrayIIC {`,
				`			mapOfAnonArrayII := mapOfAnonArrayIIV`,
				`			if err := validate.MinimumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "mapOfAnonArray", k), ii), "", int64(mapOfAnonArrayII), 10, false); err != nil {`,
				`			mapOfAnonArrayIIR = append(mapOfAnonArrayIIR, mapOfAnonArrayII`,
				`		mapOfAnonArrayIR[k] = mapOfAnonArrayIIR`,
				`	o.MapOfAnonArray = mapOfAnonArrayIR`,
			},
		},

		// load expectations for parameters in operation get_map_of_primitive_parameters.go
		"getMapOfPrimitive": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapOfPrimitiveParams() GetMapOfPrimitiveParams {`,
				`	return GetMapOfPrimitiveParams{`,
				`type GetMapOfPrimitiveParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfPrimitive map[string]int64`,
				`func (o *GetMapOfPrimitiveParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfPrimitive", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfPrimitive", "body", "", err)`,
				`		} else {`,
				`			o.MapOfPrimitive = body`,
				`			if err := o.validateMapOfPrimitiveBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfPrimitive", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapOfPrimitiveParams) validateMapOfPrimitiveBody(formats strfmt.Registry) error {`,
				`	mapOfPrimitiveIC := o.MapOfPrimitive`,
				`	mapOfPrimitiveIR := make(map[string]int64, len(mapOfPrimitiveIC)`,
				`	for k, mapOfPrimitiveIV := range mapOfPrimitiveIC {`,
				`		mapOfPrimitiveI := mapOfPrimitiveIV`,
				`		if err := validate.MaximumInt(fmt.Sprintf("%s.%v", "mapOfPrimitive", k), "body", int64(mapOfPrimitiveI), 100, false); err != nil {`,
				`		mapOfPrimitiveIR[k] = mapOfPrimitiveI`,
				`	o.MapOfPrimitive = mapOfPrimitiveIR`,
			},
		},

		// load expectations for parameters in operation get_array_parameters.go
		"getArray": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetArrayParams() GetArrayParams {`,
				`	return GetArrayParams{`,
				`type GetArrayParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	ArrayOfObject []*models.ModelObject`,
				`func (o *GetArrayParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []*models.ModelObject`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("arrayOfObject", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("arrayOfObject", "body", "", err)`,
				`		} else {`,
				`			for i := range body {`,
				`				if body[i] == nil {`,
				`					if err := body[i].Validate(route.Formats); err != nil {`,
				`						break`,
				`			if len(res) == 0 {`,
				`				o.ArrayOfObject = body`,
				`	} else {`,
				`		res = append(res, errors.Required("arrayOfObject", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_map_object_parameters.go
		"getMapObject": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapObjectParams() GetMapObjectParams {`,
				`	return GetMapObjectParams{`,
				`type GetMapObjectParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfObject models.GetMapObjectParamsBody`,
				`func (o *GetMapObjectParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body models.GetMapObjectParamsBody`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfObject", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfObject", "body", "", err)`,
				`		} else {`,
				`			if err := body.Validate(route.Formats); err != nil {`,
				`			if len(res) == 0 {`,
				`				o.MapOfObject = body`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfObject", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_map_of_map_parameters.go
		"getMapOfMap": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetMapOfMapParams() GetMapOfMapParams {`,
				`	return GetMapOfMapParams{`,
				`type GetMapOfMapParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	MapOfMap map[string]map[string]models.ModelArrayWithMax`,
				`func (o *GetMapOfMapParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]models.ModelArrayWithMax`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("mapOfMap", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("mapOfMap", "body", "", err)`,
				`		} else {`,
				`			o.MapOfMap = body`,
				`			if err := o.validateMapOfMapBody(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("mapOfMap", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetMapOfMapParams) validateMapOfMapBody(formats strfmt.Registry) error {`,
				`	mapOfMapIC := o.MapOfMap`,
				`	mapOfMapIR := make(map[string]map[string]models.ModelArrayWithMax, len(mapOfMapIC)`,
				`	for k, mapOfMapIV := range mapOfMapIC {`,
				`		mapOfMapIIC := mapOfMapIV`,
				`		mapOfMapIIR := make(map[string]models.ModelArrayWithMax, len(mapOfMapIIC)`,
				`		for kk, mapOfMapIIV := range mapOfMapIIC {`,
				`			mapOfMapII := mapOfMapIIV`,
				`			if err := mapOfMapII.Validate(formats); err != nil {`,
				`				if ve, ok := err.(*errors.Validation); ok {`,
				`					return ve.ValidateName(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "mapOfMap", k), kk)`,
				`			mapOfMapIIR[kk] = mapOfMapII`,
				`		mapOfMapIR[k] = mapOfMapIIR`,
				`	o.MapOfMap = mapOfMapIR`,
			},
		},
	}
	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1536", "fixture-1536-3.yaml"))
}

func TestGenParameter_Issue1536_MoreMaps(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	// testing fixture-1536-4.yaml with flatten
	// param body with maps

	fixtureConfig := map[string]map[string][]string{

		// load expectations for parameters in operation get_nested_map04_parameters.go
		"getNestedMap04": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMap04Params() GetNestedMap04Params {`,
				`	return GetNestedMap04Params{`,
				`type GetNestedMap04Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMap04 map[string]map[string]map[string]*bool`,
				`func (o *GetNestedMap04Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]map[string]*bool`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMap04", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMap04", "body", "", err)`,
				`		} else {`,
				`			o.NestedMap04 = body`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMap04", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_slice_and_map01_parameters.go
		"getNestedSliceAndMap01": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedSliceAndMap01Params() GetNestedSliceAndMap01Params {`,
				`	return GetNestedSliceAndMap01Params{`,
				`type GetNestedSliceAndMap01Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Unique: true`,
				`	  In: body`,
				`	*/`,
				`	NestedSliceAndMap01 []map[string][]map[string]strfmt.Date`,
				`func (o *GetNestedSliceAndMap01Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []map[string][]map[string]strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedSliceAndMap01", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedSliceAndMap01", "body", "", err)`,
				`		} else {`,
				`			o.NestedSliceAndMap01 = body`,
				`			if err := o.validateNestedSliceAndMap01Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedSliceAndMap01", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedSliceAndMap01Params) validateNestedSliceAndMap01Body(formats strfmt.Registry) error {`,
				`	if err := validate.UniqueItems("nestedSliceAndMap01", "body", o.NestedSliceAndMap01); err != nil {`,
				`	nestedSliceAndMap01IC := o.NestedSliceAndMap01`,
				`	var nestedSliceAndMap01IR []map[string][]map[string]strfmt.Date`,
				`	for i, nestedSliceAndMap01IV := range nestedSliceAndMap01IC {`,
				`		nestedSliceAndMap01IIC := nestedSliceAndMap01IV`,
				`		nestedSliceAndMap01IIR := make(map[string][]map[string]strfmt.Date, len(nestedSliceAndMap01IIC)`,
				`		for kk, nestedSliceAndMap01IIV := range nestedSliceAndMap01IIC {`,
				`			nestedSliceAndMap01IIIC := nestedSliceAndMap01IIV`,
				`			if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedSliceAndMap01", i), kk), "", nestedSliceAndMap01IIIC); err != nil {`,
				`			var nestedSliceAndMap01IIIR []map[string]strfmt.Date`,
				`			for iii, nestedSliceAndMap01IIIV := range nestedSliceAndMap01IIIC {`,
				`				nestedSliceAndMap01IIIIC := nestedSliceAndMap01IIIV`,
				`				nestedSliceAndMap01IIIIR := make(map[string]strfmt.Date, len(nestedSliceAndMap01IIIIC)`,
				`				for kkkk, nestedSliceAndMap01IIIIV := range nestedSliceAndMap01IIIIC {`,
				`					nestedSliceAndMap01IIII := nestedSliceAndMap01IIIIV`,
				`					if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedSliceAndMap01", i), kk), iii), kkkk), "", "date", nestedSliceAndMap01IIII.String(), formats); err != nil {`,
				`					nestedSliceAndMap01IIIIR[kkkk] = nestedSliceAndMap01IIII`,
				`				nestedSliceAndMap01IIIR = append(nestedSliceAndMap01IIIR, nestedSliceAndMap01IIIIR`,
				`			nestedSliceAndMap01IIR[kk] = nestedSliceAndMap01IIIR`,
				`		nestedSliceAndMap01IR = append(nestedSliceAndMap01IR, nestedSliceAndMap01IIR`,
				`	o.NestedSliceAndMap01 = nestedSliceAndMap01IR`,
			},
		},

		// load expectations for parameters in operation get_nested_map_and_slice02_parameters.go
		"getNestedMapAndSlice02": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMapAndSlice02Params() GetNestedMapAndSlice02Params {`,
				`	return GetNestedMapAndSlice02Params{`,
				`type GetNestedMapAndSlice02Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMapAndSlice02 map[string][]map[string][]map[string]*int64`,
				`func (o *GetNestedMapAndSlice02Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]map[string][]map[string]*int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMapAndSlice02", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMapAndSlice02", "body", "", err)`,
				`		} else {`,
				`			o.NestedMapAndSlice02 = body`,
				`			if err := o.validateNestedMapAndSlice02Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMapAndSlice02", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedMapAndSlice02Params) validateNestedMapAndSlice02Body(formats strfmt.Registry) error {`,
				`	nestedMapAndSlice02IC := o.NestedMapAndSlice02`,
				`	nestedMapAndSlice02IR := make(map[string][]map[string][]map[string]*int64, len(nestedMapAndSlice02IC)`,
				`	for k, nestedMapAndSlice02IV := range nestedMapAndSlice02IC {`,
				`		nestedMapAndSlice02IIC := nestedMapAndSlice02IV`,
				`		if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "nestedMapAndSlice02", k), "body", nestedMapAndSlice02IIC); err != nil {`,
				`		var nestedMapAndSlice02IIR []map[string][]map[string]*int64`,
				`		for ii, nestedMapAndSlice02IIV := range nestedMapAndSlice02IIC {`,
				`			nestedMapAndSlice02IIIC := nestedMapAndSlice02IIV`,
				`			nestedMapAndSlice02IIIR := make(map[string][]map[string]*int64, len(nestedMapAndSlice02IIIC)`,
				`			for kkk, nestedMapAndSlice02IIIV := range nestedMapAndSlice02IIIC {`,
				`				nestedMapAndSlice02IIIIC := nestedMapAndSlice02IIIV`,
				`				if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMapAndSlice02", k), ii), kkk), "", nestedMapAndSlice02IIIIC); err != nil {`,
				`				var nestedMapAndSlice02IIIIR []map[string]*int64`,
				`				for iiii, nestedMapAndSlice02IIIIV := range nestedMapAndSlice02IIIIC {`,
				`					nestedMapAndSlice02IIIIIC := nestedMapAndSlice02IIIIV`,
				`					nestedMapAndSlice02IIIIIR := make(map[string]*int64, len(nestedMapAndSlice02IIIIIC)`,
				`					for kkkkk, nestedMapAndSlice02IIIIIV := range nestedMapAndSlice02IIIIIC {`,
				`						if nestedMapAndSlice02IIIIIV == nil {`,
				`						nestedMapAndSlice02IIIII := nestedMapAndSlice02IIIIIV`,
				`						if err := validate.MinimumInt(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMapAndSlice02", k), ii), kkk), iiii), kkkkk), "", int64(*nestedMapAndSlice02IIIII), 0, false); err != nil {`,
				`						nestedMapAndSlice02IIIIIR[kkkkk] = nestedMapAndSlice02IIIII`,
				`					nestedMapAndSlice02IIIIR = append(nestedMapAndSlice02IIIIR, nestedMapAndSlice02IIIIIR`,
				`				nestedMapAndSlice02IIIR[kkk] = nestedMapAndSlice02IIIIR`,
				`			nestedMapAndSlice02IIR = append(nestedMapAndSlice02IIR, nestedMapAndSlice02IIIR`,
				`		nestedMapAndSlice02IR[k] = nestedMapAndSlice02IIR`,
				`	o.NestedMapAndSlice02 = nestedMapAndSlice02IR`,
			},
		},

		// load expectations for parameters in operation get_nested_slice_and_map03_parameters.go
		"getNestedSliceAndMap03": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedSliceAndMap03Params() GetNestedSliceAndMap03Params {`,
				`	return GetNestedSliceAndMap03Params{`,
				`type GetNestedSliceAndMap03Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Unique: true`,
				`	  In: body`,
				`	*/`,
				`	NestedSliceAndMap03 []map[string][]map[string]string`,
				`func (o *GetNestedSliceAndMap03Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []map[string][]map[string]string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedSliceAndMap03", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedSliceAndMap03", "body", "", err)`,
				`		} else {`,
				`			o.NestedSliceAndMap03 = body`,
				`			if err := o.validateNestedSliceAndMap03Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedSliceAndMap03", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedSliceAndMap03Params) validateNestedSliceAndMap03Body(formats strfmt.Registry) error {`,
				`	if err := validate.UniqueItems("nestedSliceAndMap03", "body", o.NestedSliceAndMap03); err != nil {`,
				`	nestedSliceAndMap03IC := o.NestedSliceAndMap03`,
				`	var nestedSliceAndMap03IR []map[string][]map[string]string`,
				`	for i, nestedSliceAndMap03IV := range nestedSliceAndMap03IC {`,
				`		nestedSliceAndMap03IIC := nestedSliceAndMap03IV`,
				`		nestedSliceAndMap03IIR := make(map[string][]map[string]string, len(nestedSliceAndMap03IIC)`,
				`		for kk, nestedSliceAndMap03IIV := range nestedSliceAndMap03IIC {`,
				`			nestedSliceAndMap03IIIC := nestedSliceAndMap03IIV`,
				`			if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedSliceAndMap03", i), kk), "", nestedSliceAndMap03IIIC); err != nil {`,
				`			var nestedSliceAndMap03IIIR []map[string]string`,
				`			for _, nestedSliceAndMap03IIIV := range nestedSliceAndMap03IIIC {`,
				`				nestedSliceAndMap03IIIIC := nestedSliceAndMap03IIIV`,
				`				nestedSliceAndMap03IIIIR := make(map[string]string, len(nestedSliceAndMap03IIIIC)`,
				`				for kkkk, nestedSliceAndMap03IIIIV := range nestedSliceAndMap03IIIIC {`,
				`					nestedSliceAndMap03IIII := nestedSliceAndMap03IIIIV`,
				`					nestedSliceAndMap03IIIIR[kkkk] = nestedSliceAndMap03IIII`,
				`				nestedSliceAndMap03IIIR = append(nestedSliceAndMap03IIIR, nestedSliceAndMap03IIIIR`,
				`			nestedSliceAndMap03IIR[kk] = nestedSliceAndMap03IIIR`,
				`		nestedSliceAndMap03IR = append(nestedSliceAndMap03IR, nestedSliceAndMap03IIR`,
				`	o.NestedSliceAndMap03 = nestedSliceAndMap03IR`,
			},
		},

		// load expectations for parameters in operation get_nested_array03_parameters.go
		"getNestedArray03": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedArray03Params() GetNestedArray03Params {`,
				`	return GetNestedArray03Params{`,
				`type GetNestedArray03Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedArray03 [][][]string`,
				`func (o *GetNestedArray03Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][]string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedArray03", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedArray03", "body", "", err)`,
				`		} else {`,
				`			o.NestedArray03 = body`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedArray03", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_array04_parameters.go
		"getNestedArray04": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedArray04Params() GetNestedArray04Params {`,
				`	return GetNestedArray04Params{`,
				`type GetNestedArray04Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Unique: true`,
				`	  In: body`,
				`	*/`,
				`	NestedArray03 [][][]string`,
				`func (o *GetNestedArray04Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][]string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedArray03", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedArray03", "body", "", err)`,
				`		} else {`,
				`			o.NestedArray03 = body`,
				`			if err := o.validateNestedArray03Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedArray03", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedArray04Params) validateNestedArray03Body(formats strfmt.Registry) error {`,
				`	if err := validate.UniqueItems("nestedArray03", "body", o.NestedArray03); err != nil {`,
				`	nestedArray03IC := o.NestedArray03`,
				`	var nestedArray03IR [][][]string`,
				`	for i, nestedArray03IV := range nestedArray03IC {`,
				`		nestedArray03IIC := nestedArray03IV`,
				`		if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "nestedArray03", i), "body", nestedArray03IIC); err != nil {`,
				`		if len(nestedArray03IIC) > 0 {`,
				`			var nestedArray03IIR [][]string`,
				`			for ii, nestedArray03IIV := range nestedArray03IIC {`,
				`				nestedArray03IIIC := nestedArray03IIV`,
				`				if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray03", i), ii), "", nestedArray03IIIC); err != nil {`,
				`				if len(nestedArray03IIIC) > 0 {`,
				`					var nestedArray03IIIR []string`,
				`					for _, nestedArray03IIIV := range nestedArray03IIIC {`,
				`						nestedArray03III := nestedArray03IIIV`,
				`						nestedArray03IIIR = append(nestedArray03IIIR, nestedArray03III`,
				`					nestedArray03IIR = append(nestedArray03IIR, nestedArray03IIIR`,
				`			nestedArray03IR = append(nestedArray03IR, nestedArray03IIR`,
				`	o.NestedArray03 = nestedArray03IR`,
			},
		},

		// load expectations for parameters in operation get_nested_map_and_slice01_parameters.go
		"getNestedMapAndSlice01": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMapAndSlice01Params() GetNestedMapAndSlice01Params {`,
				`	return GetNestedMapAndSlice01Params{`,
				`type GetNestedMapAndSlice01Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMapAndSlice01 map[string][]map[string][]map[string]strfmt.Date`,
				`func (o *GetNestedMapAndSlice01Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]map[string][]map[string]strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMapAndSlice01", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMapAndSlice01", "body", "", err)`,
				`		} else {`,
				`			o.NestedMapAndSlice01 = body`,
				`			if err := o.validateNestedMapAndSlice01Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMapAndSlice01", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedMapAndSlice01Params) validateNestedMapAndSlice01Body(formats strfmt.Registry) error {`,
				`	nestedMapAndSlice01IC := o.NestedMapAndSlice01`,
				`	nestedMapAndSlice01IR := make(map[string][]map[string][]map[string]strfmt.Date, len(nestedMapAndSlice01IC)`,
				`	for k, nestedMapAndSlice01IV := range nestedMapAndSlice01IC {`,
				`		nestedMapAndSlice01IIC := nestedMapAndSlice01IV`,
				`		if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "nestedMapAndSlice01", k), "body", nestedMapAndSlice01IIC); err != nil {`,
				`		var nestedMapAndSlice01IIR []map[string][]map[string]strfmt.Date`,
				`		for ii, nestedMapAndSlice01IIV := range nestedMapAndSlice01IIC {`,
				`			nestedMapAndSlice01IIIC := nestedMapAndSlice01IIV`,
				`			nestedMapAndSlice01IIIR := make(map[string][]map[string]strfmt.Date, len(nestedMapAndSlice01IIIC)`,
				`			for kkk, nestedMapAndSlice01IIIV := range nestedMapAndSlice01IIIC {`,
				`				nestedMapAndSlice01IIIIC := nestedMapAndSlice01IIIV`,
				`				if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMapAndSlice01", k), ii), kkk), "", nestedMapAndSlice01IIIIC); err != nil {`,
				`				var nestedMapAndSlice01IIIIR []map[string]strfmt.Date`,
				`				for iiii, nestedMapAndSlice01IIIIV := range nestedMapAndSlice01IIIIC {`,
				`					nestedMapAndSlice01IIIIIC := nestedMapAndSlice01IIIIV`,
				`					nestedMapAndSlice01IIIIIR := make(map[string]strfmt.Date, len(nestedMapAndSlice01IIIIIC)`,
				`					for kkkkk, nestedMapAndSlice01IIIIIV := range nestedMapAndSlice01IIIIIC {`,
				`						nestedMapAndSlice01IIIII := nestedMapAndSlice01IIIIIV`,
				`						if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMapAndSlice01", k), ii), kkk), iiii), kkkkk), "", "date", nestedMapAndSlice01IIIII.String(), formats); err != nil {`,
				`						nestedMapAndSlice01IIIIIR[kkkkk] = nestedMapAndSlice01IIIII`,
				`					nestedMapAndSlice01IIIIR = append(nestedMapAndSlice01IIIIR, nestedMapAndSlice01IIIIIR`,
				`				nestedMapAndSlice01IIIR[kkk] = nestedMapAndSlice01IIIIR`,
				`			nestedMapAndSlice01IIR = append(nestedMapAndSlice01IIR, nestedMapAndSlice01IIIR`,
				`		nestedMapAndSlice01IR[k] = nestedMapAndSlice01IIR`,
				`	o.NestedMapAndSlice01 = nestedMapAndSlice01IR`,
			},
		},

		// load expectations for parameters in operation get_nested_slice_and_map03_ref_parameters.go
		"getNestedSliceAndMap03Ref": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedSliceAndMap03RefParams() GetNestedSliceAndMap03RefParams {`,
				`	return GetNestedSliceAndMap03RefParams{`,
				`type GetNestedSliceAndMap03RefParams struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  In: body`,
				`	*/`,
				`	NestedSliceAndMap03Ref models.NestedSliceAndMap03Ref`,
				`func (o *GetNestedSliceAndMap03RefParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body models.NestedSliceAndMap03Ref`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("nestedSliceAndMap03Ref", "body", "", err)`,
				`		} else {`,
				`			if err := body.Validate(route.Formats); err != nil {`,
				`			if len(res) == 0 {`,
				`				o.NestedSliceAndMap03Ref = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_map02_parameters.go
		"getNestedMap02": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMap02Params() GetNestedMap02Params {`,
				`	return GetNestedMap02Params{`,
				`type GetNestedMap02Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMap02 map[string]map[string]map[string]*string`,
				`func (o *GetNestedMap02Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]map[string]*string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMap02", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMap02", "body", "", err)`,
				`		} else {`,
				`			o.NestedMap02 = body`,
				`			if err := o.validateNestedMap02Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMap02", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedMap02Params) validateNestedMap02Body(formats strfmt.Registry) error {`,
				`	nestedMap02IC := o.NestedMap02`,
				`	nestedMap02IR := make(map[string]map[string]map[string]*string, len(nestedMap02IC)`,
				`	for k, nestedMap02IV := range nestedMap02IC {`,
				`		nestedMap02IIC := nestedMap02IV`,
				`		nestedMap02IIR := make(map[string]map[string]*string, len(nestedMap02IIC)`,
				`		for kk, nestedMap02IIV := range nestedMap02IIC {`,
				`			nestedMap02IIIC := nestedMap02IIV`,
				`			nestedMap02IIIR := make(map[string]*string, len(nestedMap02IIIC)`,
				`			for kkk, nestedMap02IIIV := range nestedMap02IIIC {`,
				`				if nestedMap02IIIV == nil {`,
				`				nestedMap02III := nestedMap02IIIV`,
				`				if err := validate.MinLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMap02", k), kk), kkk), "", (*nestedMap02III), 0); err != nil {`,
				`				nestedMap02IIIR[kkk] = nestedMap02III`,
				`			nestedMap02IIR[kk] = nestedMap02IIIR`,
				`		nestedMap02IR[k] = nestedMap02IIR`,
				`	o.NestedMap02 = nestedMap02IR`,
			},
		},

		// load expectations for parameters in operation get_nested_map_and_slice03_parameters.go
		"getNestedMapAndSlice03": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMapAndSlice03Params() GetNestedMapAndSlice03Params {`,
				`	return GetNestedMapAndSlice03Params{`,
				`type GetNestedMapAndSlice03Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMapAndSlice03 map[string][]map[string][]map[string]int64`,
				`func (o *GetNestedMapAndSlice03Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string][]map[string][]map[string]int64`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMapAndSlice03", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMapAndSlice03", "body", "", err)`,
				`		} else {`,
				`			o.NestedMapAndSlice03 = body`,
				`			if err := o.validateNestedMapAndSlice03Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMapAndSlice03", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedMapAndSlice03Params) validateNestedMapAndSlice03Body(formats strfmt.Registry) error {`,
				`	nestedMapAndSlice03IC := o.NestedMapAndSlice03`,
				`	nestedMapAndSlice03IR := make(map[string][]map[string][]map[string]int64, len(nestedMapAndSlice03IC)`,
				`	for k, nestedMapAndSlice03IV := range nestedMapAndSlice03IC {`,
				`		nestedMapAndSlice03IIC := nestedMapAndSlice03IV`,
				`		if err := validate.UniqueItems(fmt.Sprintf("%s.%v", "nestedMapAndSlice03", k), "body", nestedMapAndSlice03IIC); err != nil {`,
				`		var nestedMapAndSlice03IIR []map[string][]map[string]int64`,
				`		for ii, nestedMapAndSlice03IIV := range nestedMapAndSlice03IIC {`,
				`			nestedMapAndSlice03IIIC := nestedMapAndSlice03IIV`,
				`			nestedMapAndSlice03IIIR := make(map[string][]map[string]int64, len(nestedMapAndSlice03IIIC)`,
				`			for kkk, nestedMapAndSlice03IIIV := range nestedMapAndSlice03IIIC {`,
				`				nestedMapAndSlice03IIIIC := nestedMapAndSlice03IIIV`,
				`				if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMapAndSlice03", k), ii), kkk), "", nestedMapAndSlice03IIIIC); err != nil {`,
				`				var nestedMapAndSlice03IIIIR []map[string]int64`,
				`				for _, nestedMapAndSlice03IIIIV := range nestedMapAndSlice03IIIIC {`,
				`					nestedMapAndSlice03IIIIIC := nestedMapAndSlice03IIIIV`,
				`					nestedMapAndSlice03IIIIIR := make(map[string]int64, len(nestedMapAndSlice03IIIIIC)`,
				`					for kkkkk, nestedMapAndSlice03IIIIIV := range nestedMapAndSlice03IIIIIC {`,
				`						nestedMapAndSlice03IIIII := nestedMapAndSlice03IIIIIV`,
				`						nestedMapAndSlice03IIIIIR[kkkkk] = nestedMapAndSlice03IIIII`,
				`					nestedMapAndSlice03IIIIR = append(nestedMapAndSlice03IIIIR, nestedMapAndSlice03IIIIIR`,
				`				nestedMapAndSlice03IIIR[kkk] = nestedMapAndSlice03IIIIR`,
				`			nestedMapAndSlice03IIR = append(nestedMapAndSlice03IIR, nestedMapAndSlice03IIIR`,
				`		nestedMapAndSlice03IR[k] = nestedMapAndSlice03IIR`,
				`	o.NestedMapAndSlice03 = nestedMapAndSlice03IR`,
			},
		},

		// load expectations for parameters in operation get_nested_slice_and_map02_parameters.go
		"getNestedSliceAndMap02": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedSliceAndMap02Params() GetNestedSliceAndMap02Params {`,
				`	return GetNestedSliceAndMap02Params{`,
				`type GetNestedSliceAndMap02Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Unique: true`,
				`	  In: body`,
				`	*/`,
				`	NestedSliceAndMap02 []map[string][]map[string]*string`,
				`func (o *GetNestedSliceAndMap02Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body []map[string][]map[string]*string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedSliceAndMap02", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedSliceAndMap02", "body", "", err)`,
				`		} else {`,
				`			o.NestedSliceAndMap02 = body`,
				`			if err := o.validateNestedSliceAndMap02Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedSliceAndMap02", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedSliceAndMap02Params) validateNestedSliceAndMap02Body(formats strfmt.Registry) error {`,
				`	if err := validate.UniqueItems("nestedSliceAndMap02", "body", o.NestedSliceAndMap02); err != nil {`,
				`	nestedSliceAndMap02IC := o.NestedSliceAndMap02`,
				`	var nestedSliceAndMap02IR []map[string][]map[string]*string`,
				`	for i, nestedSliceAndMap02IV := range nestedSliceAndMap02IC {`,
				`		nestedSliceAndMap02IIC := nestedSliceAndMap02IV`,
				`		nestedSliceAndMap02IIR := make(map[string][]map[string]*string, len(nestedSliceAndMap02IIC)`,
				`		for kk, nestedSliceAndMap02IIV := range nestedSliceAndMap02IIC {`,
				`			nestedSliceAndMap02IIIC := nestedSliceAndMap02IIV`,
				`			if err := validate.UniqueItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedSliceAndMap02", i), kk), "", nestedSliceAndMap02IIIC); err != nil {`,
				`			var nestedSliceAndMap02IIIR []map[string]*string`,
				`			for iii, nestedSliceAndMap02IIIV := range nestedSliceAndMap02IIIC {`,
				`				nestedSliceAndMap02IIIIC := nestedSliceAndMap02IIIV`,
				`				nestedSliceAndMap02IIIIR := make(map[string]*string, len(nestedSliceAndMap02IIIIC)`,
				`				for kkkk, nestedSliceAndMap02IIIIV := range nestedSliceAndMap02IIIIC {`,
				`					if nestedSliceAndMap02IIIIV == nil {`,
				`					nestedSliceAndMap02IIII := nestedSliceAndMap02IIIIV`,
				`					if err := validate.MinLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedSliceAndMap02", i), kk), iii), kkkk), "", (*nestedSliceAndMap02IIII), 0); err != nil {`,
				`					nestedSliceAndMap02IIIIR[kkkk] = nestedSliceAndMap02IIII`,
				`				nestedSliceAndMap02IIIR = append(nestedSliceAndMap02IIIR, nestedSliceAndMap02IIIIR`,
				`			nestedSliceAndMap02IIR[kk] = nestedSliceAndMap02IIIR`,
				`		nestedSliceAndMap02IR = append(nestedSliceAndMap02IR, nestedSliceAndMap02IIR`,
				`	o.NestedSliceAndMap02 = nestedSliceAndMap02IR`,
			},
		},

		// load expectations for parameters in operation get_nested_map01_parameters.go
		"getNestedMap01": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMap01Params() GetNestedMap01Params {`,
				`	return GetNestedMap01Params{`,
				`type GetNestedMap01Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMap01 map[string]map[string]map[string]strfmt.Date`,
				`func (o *GetNestedMap01Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]map[string]strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMap01", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMap01", "body", "", err)`,
				`		} else {`,
				`			o.NestedMap01 = body`,
				`			if err := o.validateNestedMap01Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMap01", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedMap01Params) validateNestedMap01Body(formats strfmt.Registry) error {`,
				`	nestedMap01IC := o.NestedMap01`,
				`	nestedMap01IR := make(map[string]map[string]map[string]strfmt.Date, len(nestedMap01IC)`,
				`	for k, nestedMap01IV := range nestedMap01IC {`,
				`		nestedMap01IIC := nestedMap01IV`,
				`		nestedMap01IIR := make(map[string]map[string]strfmt.Date, len(nestedMap01IIC)`,
				`		for kk, nestedMap01IIV := range nestedMap01IIC {`,
				`			nestedMap01IIIC := nestedMap01IIV`,
				`			nestedMap01IIIR := make(map[string]strfmt.Date, len(nestedMap01IIIC)`,
				`			for kkk, nestedMap01IIIV := range nestedMap01IIIC {`,
				`				nestedMap01III := nestedMap01IIIV`,
				`				if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedMap01", k), kk), kkk), "", "date", nestedMap01III.String(), formats); err != nil {`,
				`				nestedMap01IIIR[kkk] = nestedMap01III`,
				`			nestedMap01IIR[kk] = nestedMap01IIIR`,
				`		nestedMap01IR[k] = nestedMap01IIR`,
				`	o.NestedMap01 = nestedMap01IR`,
			},
		},

		// load expectations for parameters in operation get_nested_map03_parameters.go
		"getNestedMap03": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedMap03Params() GetNestedMap03Params {`,
				`	return GetNestedMap03Params{`,
				`type GetNestedMap03Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  In: body`,
				`	*/`,
				`	NestedMap03 map[string]map[string]map[string]string`,
				`func (o *GetNestedMap03Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]map[string]string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedMap03", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedMap03", "body", "", err)`,
				`		} else {`,
				`			o.NestedMap03 = body`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedMap03", "body")`,
				`		return errors.CompositeValidationError(res...`,
			},
		},

		// load expectations for parameters in operation get_nested_array01_parameters.go
		"getNestedArray01": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedArray01Params() GetNestedArray01Params {`,
				`	return GetNestedArray01Params{`,
				`type GetNestedArray01Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Max Items: 10`,
				`	  In: body`,
				`	*/`,
				`	NestedArray01 [][][]strfmt.Date`,
				`func (o *GetNestedArray01Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][]strfmt.Date`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedArray01", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedArray01", "body", "", err)`,
				`		} else {`,
				`			o.NestedArray01 = body`,
				`			if err := o.validateNestedArray01Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedArray01", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedArray01Params) validateNestedArray01Body(formats strfmt.Registry) error {`,
				`	nestedArray01Size := int64(len(o.NestedArray01)`,
				`	if err := validate.MaxItems("nestedArray01", "body", nestedArray01Size, 10); err != nil {`,
				`	nestedArray01IC := o.NestedArray01`,
				`	var nestedArray01IR [][][]strfmt.Date`,
				`	for i, nestedArray01IV := range nestedArray01IC {`,
				`		nestedArray01IIC := nestedArray01IV`,
				`		nestedArray01ISize := int64(len(nestedArray01IIC)`,
				`		if err := validate.MaxItems(fmt.Sprintf("%s.%v", "nestedArray01", i), "body", nestedArray01ISize, 10); err != nil {`,
				`		if len(nestedArray01IIC) > 0 {`,
				`			var nestedArray01IIR [][]strfmt.Date`,
				`			for ii, nestedArray01IIV := range nestedArray01IIC {`,
				`				nestedArray01IIIC := nestedArray01IIV`,
				`				nestedArray01IiiSize := int64(len(nestedArray01IIIC)`,
				`				if err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray01", i), ii), "", nestedArray01IiiSize, 10); err != nil {`,
				`				if len(nestedArray01IIIC) > 0 {`,
				`					var nestedArray01IIIR []strfmt.Date`,
				`					for iii, nestedArray01IIIV := range nestedArray01IIIC {`,
				`						nestedArray01III := nestedArray01IIIV`,
				`						if err := validate.FormatOf(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray01", i), ii), iii), "", "date", nestedArray01III.String(), formats); err != nil {`,
				`						nestedArray01IIIR = append(nestedArray01IIIR, nestedArray01III`,
				`					nestedArray01IIR = append(nestedArray01IIR, nestedArray01IIIR`,
				`			nestedArray01IR = append(nestedArray01IR, nestedArray01IIR`,
				`	o.NestedArray01 = nestedArray01IR`,
			},
		},

		// load expectations for parameters in operation get_nested_array02_parameters.go
		"getNestedArray02": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedArray02Params() GetNestedArray02Params {`,
				`	return GetNestedArray02Params{`,
				`type GetNestedArray02Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  Required: true`,
				`	  Max Items: 10`,
				`	  In: body`,
				`	*/`,
				`	NestedArray01 [][][]*string`,
				`func (o *GetNestedArray02Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body [][][]*string`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			if err == io.EOF {`,
				`				res = append(res, errors.Required("nestedArray01", "body")`,
				`			} else {`,
				`				res = append(res, errors.NewParseError("nestedArray01", "body", "", err)`,
				`		} else {`,
				`			o.NestedArray01 = body`,
				`			if err := o.validateNestedArray01Body(route.Formats); err != nil {`,
				`	} else {`,
				`		res = append(res, errors.Required("nestedArray01", "body")`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedArray02Params) validateNestedArray01Body(formats strfmt.Registry) error {`,
				`	nestedArray01Size := int64(len(o.NestedArray01)`,
				`	if err := validate.MaxItems("nestedArray01", "body", nestedArray01Size, 10); err != nil {`,
				`	nestedArray01IC := o.NestedArray01`,
				`	var nestedArray01IR [][][]*string`,
				`	for i, nestedArray01IV := range nestedArray01IC {`,
				`		nestedArray01IIC := nestedArray01IV`,
				`		nestedArray01ISize := int64(len(nestedArray01IIC)`,
				`		if err := validate.MaxItems(fmt.Sprintf("%s.%v", "nestedArray01", i), "body", nestedArray01ISize, 10); err != nil {`,
				`		if len(nestedArray01IIC) > 0 {`,
				`			var nestedArray01IIR [][]*string`,
				`			for ii, nestedArray01IIV := range nestedArray01IIC {`,
				`				nestedArray01IIIC := nestedArray01IIV`,
				`				nestedArray01IiiSize := int64(len(nestedArray01IIIC)`,
				`				if err := validate.MaxItems(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray01", i), ii), "", nestedArray01IiiSize, 10); err != nil {`,
				`				if len(nestedArray01IIIC) > 0 {`,
				`					var nestedArray01IIIR []*string`,
				`					for iii, nestedArray01IIIV := range nestedArray01IIIC {`,
				`						if nestedArray01IIIV == nil {`,
				// do we need Required on nullable in items?
				// without Required
				`							continue`,
				// with Required
				//`							return errors.Required(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray01", i), ii), iii), ""`,
				`						nestedArray01III := nestedArray01IIIV`,
				`						if err := validate.MinLength(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedArray01", i), ii), iii), "", (*nestedArray01III), 0); err != nil {`,
				`						nestedArray01IIIR = append(nestedArray01IIIR, nestedArray01III`,
				`					nestedArray01IIR = append(nestedArray01IIR, nestedArray01IIIR`,
				`			nestedArray01IR = append(nestedArray01IR, nestedArray01IIR`,
				`	o.NestedArray01 = nestedArray01IR`,
			},
		},
		// load expectations for parameters in operation get_nested_ref_no_validation01_parameters.go
		"getNestedRefNoValidation01": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedRefNoValidation01Params() GetNestedRefNoValidation01Params {`,
				`	return GetNestedRefNoValidation01Params{`,
				`type GetNestedRefNoValidation01Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  In: body`,
				`	*/`,
				`	NestedRefNovalidation01 map[string]models.NestedRefNoValidation`,
				`func (o *GetNestedRefNoValidation01Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]models.NestedRefNoValidation`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("nestedRefNovalidation01", "body", "", err)`,
				`		} else {`,
				`			for k := range body {`,
				`				if err := o.NestedRefNovalidation01[k].Validate(route.Formats); err != nil {`,
				`					break`,
				`			if len(res) == 0 {`,
				`				o.NestedRefNovalidation01 = body`,
				`		return errors.CompositeValidationError(res...`,
			},
		},
		// load expectations for parameters in operation get_nested_ref_no_validation02_parameters.go
		"getNestedRefNoValidation02": map[string][]string{ // fixture index
			"serverParameter": []string{ // executed template
				// expected code lines
				`func NewGetNestedRefNoValidation02Params() GetNestedRefNoValidation02Params {`,
				`	return GetNestedRefNoValidation02Params{`,
				`type GetNestedRefNoValidation02Params struct {`,
				"	HTTPRequest *http.Request `json:\"-\"`",
				`	/*`,
				`	  In: body`,
				`	*/`,
				`	NestedRefNovalidation02 map[string]map[string]models.NestedRefNoValidation`,
				`func (o *GetNestedRefNoValidation02Params) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {`,
				`	o.HTTPRequest = r`,
				`	if runtime.HasBody(r) {`,
				`		defer r.Body.Close(`,
				`		var body map[string]map[string]models.NestedRefNoValidation`,
				`		if err := route.Consumer.Consume(r.Body, &body); err != nil {`,
				`			res = append(res, errors.NewParseError("nestedRefNovalidation02", "body", "", err)`,
				`		} else {`,
				`			o.NestedRefNovalidation02 = body`,
				`			if err := o.validateNestedRefNovalidation02Body(route.Formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedRefNoValidation02Params) validateNestedRefNovalidation02Body(formats strfmt.Registry) error {`,
				`	nestedRefNovalidation02IC := o.NestedRefNovalidation02`,
				`	nestedRefNovalidation02IR := make(map[string]map[string]models.NestedRefNoValidation, len(nestedRefNovalidation02IC)`,
				`	for k, nestedRefNovalidation02IV := range nestedRefNovalidation02IC {`,
				`		nestedRefNovalidation02IIC := nestedRefNovalidation02IV`,
				`		nestedRefNovalidation02IIR := make(map[string]models.NestedRefNoValidation, len(nestedRefNovalidation02IIC)`,
				`		for kk, nestedRefNovalidation02IIV := range nestedRefNovalidation02IIC {`,
				`			nestedRefNovalidation02II := nestedRefNovalidation02IIV`,
				`			if err := nestedRefNovalidation02II.Validate(formats); err != nil {`,
				`				if ve, ok := err.(*errors.Validation); ok {`,
				`					return ve.ValidateName(fmt.Sprintf("%s.%v", fmt.Sprintf("%s.%v", "nestedRefNovalidation02", k), kk)`,
				`			nestedRefNovalidation02IIR[kk] = nestedRefNovalidation02II`,
				`		nestedRefNovalidation02IR[k] = nestedRefNovalidation02IIR`,
				`	o.NestedRefNovalidation02 = nestedRefNovalidation02IR`,
			},
		},
	}

	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1536", "fixture-1536-4.yaml"))
}
