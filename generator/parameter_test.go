package generator

import (
	"fmt"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/stretchr/testify/assert"
)

// tests the parameters for generation

var simplePathParams = []simpleParamContext{
	{"siBool", "simplePathParams", simpleResolvedType("boolean", "", nil), "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}},
	{"siString", "simplePathParams", simpleResolvedType("string", "", nil), "", "", codeGenOpBuilder{}},
	{"siInt", "simplePathParams", simpleResolvedType("integer", "", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siInt32", "simplePathParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siInt64", "simplePathParams", simpleResolvedType("integer", "int64", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siFloat", "simplePathParams", simpleResolvedType("number", "", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
	{"siFloat32", "simplePathParams", simpleResolvedType("number", "float", nil), "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}},
	{"siFloat64", "simplePathParams", simpleResolvedType("number", "double", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
}

func TestSimplePathParams(t *testing.T) {
	b, err := opBuilder("simplePathParams", "../fixtures/codegen/todolist.simplepath.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simplePathParams {
		v.B = b
		if !v.assertSimpleParameter(t) {
			t.FailNow()
		}
	}
}

var simpleHeaderParams = []simpleParamContext{
	{"id", "simpleHeaderParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siBool", "simpleHeaderParams", simpleResolvedType("boolean", "", nil), "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}},
	{"siString", "simpleHeaderParams", simpleResolvedType("string", "", nil), "", "", codeGenOpBuilder{}},
	{"siInt", "simpleHeaderParams", simpleResolvedType("integer", "", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siInt32", "simpleHeaderParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siInt64", "simpleHeaderParams", simpleResolvedType("integer", "int64", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siFloat", "simpleHeaderParams", simpleResolvedType("number", "", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
	{"siFloat32", "simpleHeaderParams", simpleResolvedType("number", "float", nil), "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}},
	{"siFloat64", "simpleHeaderParams", simpleResolvedType("number", "double", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
}

func TestSimpleHeaderParams(t *testing.T) {
	b, err := opBuilder("simpleHeaderParams", "../fixtures/codegen/todolist.simpleheader.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleHeaderParams {
		v.B = b
		if !v.assertSimpleParameter(t) {
			t.FailNow()
		}
	}
}

var simpleFormParams = []simpleParamContext{
	{"id", "simpleFormParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siBool", "simpleFormParams", simpleResolvedType("boolean", "", nil), "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}},
	{"siString", "simpleFormParams", simpleResolvedType("string", "", nil), "", "", codeGenOpBuilder{}},
	{"siInt", "simpleFormParams", simpleResolvedType("integer", "", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siInt32", "simpleFormParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siInt64", "simpleFormParams", simpleResolvedType("integer", "int64", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siFloat", "simpleFormParams", simpleResolvedType("number", "", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
	{"siFloat32", "simpleFormParams", simpleResolvedType("number", "float", nil), "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}},
	{"siFloat64", "simpleFormParams", simpleResolvedType("number", "double", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
}

func TestSimpleFormParams(t *testing.T) {
	b, err := opBuilder("simpleFormParams", "../fixtures/codegen/todolist.simpleform.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleFormParams {
		v.B = b
		if !v.assertSimpleParameter(t) {
			t.FailNow()
		}
	}
}

var simpleQueryParams = []simpleParamContext{
	{"id", "simpleQueryParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siBool", "simpleQueryParams", simpleResolvedType("boolean", "", nil), "swag.FormatBool", "swag.ConvertBool", codeGenOpBuilder{}},
	{"siString", "simpleQueryParams", simpleResolvedType("string", "", nil), "", "", codeGenOpBuilder{}},
	{"siInt", "simpleQueryParams", simpleResolvedType("integer", "", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siInt32", "simpleQueryParams", simpleResolvedType("integer", "int32", nil), "swag.FormatInt32", "swag.ConvertInt32", codeGenOpBuilder{}},
	{"siInt64", "simpleQueryParams", simpleResolvedType("integer", "int64", nil), "swag.FormatInt64", "swag.ConvertInt64", codeGenOpBuilder{}},
	{"siFloat", "simpleQueryParams", simpleResolvedType("number", "", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
	{"siFloat32", "simpleQueryParams", simpleResolvedType("number", "float", nil), "swag.FormatFloat32", "swag.ConvertFloat32", codeGenOpBuilder{}},
	{"siFloat64", "simpleQueryParams", simpleResolvedType("number", "double", nil), "swag.FormatFloat64", "swag.ConvertFloat64", codeGenOpBuilder{}},
}

func TestSimpleQueryParams(t *testing.T) {
	b, err := opBuilder("simpleQueryParams", "../fixtures/codegen/todolist.simplequery.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, v := range simpleQueryParams {
		v.B = b
		if !v.assertSimpleParameter(t) {
			t.FailNow()
		}
	}
}

type simpleParamContext struct {
	Name      string
	OpID      string
	Type      resolvedType
	Formatter string
	Converter string
	B         codeGenOpBuilder
}

func (ctx simpleParamContext) assertSimpleParameter(t testing.TB) bool {
	op, err := ctx.B.Doc.OperationForName(ctx.OpID)
	if assert.True(t, err) && assert.NotNil(t, op) {
		resolver := &typeResolver{ModelsPackage: ctx.B.ModelsPackage, Doc: ctx.B.Doc}
		for _, param := range op.Parameters {
			if param.Name == ctx.Name {
				gp, err := ctx.B.MakeParameter("a", resolver, param)
				if assert.NoError(t, err) {
					return assert.True(t, ctx.assertGenParam(t, param, gp))
				}
			}
		}
		return false
	}
	return false
}

func (ctx simpleParamContext) assertGenParam(t testing.TB, param spec.Parameter, gp GenParameter) bool {
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
	// verify rendered template
	if param.In == "body" {
		if !ctx.assertBodyParam(t, param, gp) {
			return false
		}
		return true
	}

	return ctx.assertParamItems(t, param, gp)
}

func (ctx simpleParamContext) assertBodyParam(t testing.TB, param spec.Parameter, gp GenParameter) bool {
	if !assert.Equal(t, "body", param.In) || !assert.Equal(t, "body", gp.Location) {
		return false
	}
	if !assert.NotNil(t, gp.Schema) {
		return false
	}
	return true
}

func (ctx simpleParamContext) assertParamItems(t testing.TB, param spec.Parameter, gp GenParameter) bool {
	if param.Items != nil {
		pItems, gpItems := param.Items, gp.Child
		// went with the verbose option here, easier to debug
		if !assert.Equal(t, param.CollectionFormat, gp.CollectionFormat) {
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

	}
	return true
}
