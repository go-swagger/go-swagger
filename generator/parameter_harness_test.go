// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"fmt"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func (ctx *paramTestContext) assertParameter(t *testing.T) (result bool) {
	t.Helper()

	defer func() {
		result = !t.Failed()
	}()

	_, _, op, err := ctx.B.Analyzed.OperationForName(ctx.OpID)

	require.True(t, err)
	require.NotNil(t, op)

	resolver := &typeResolver{ModelsPackage: ctx.B.ModelsPackage, Doc: ctx.B.Doc}
	resolver.KnownDefs = make(map[string]struct{})

	for k := range ctx.B.Doc.Spec().Definitions {
		resolver.KnownDefs[k] = struct{}{}
	}
	for _, param := range op.Parameters {
		if param.Name != ctx.Name {
			continue
		}

		gp, err := ctx.B.MakeParameter("a", resolver, param, nil)
		require.NoError(t, err)

		assert.True(t, ctx.assertGenParam(t, param, gp))
	}

	return !t.Failed()
}

func (ctx *paramTestContext) assertGenParam(t *testing.T, param spec.Parameter, gp GenParameter) bool {
	t.Helper()

	if !ctx.assertGenParamDefinitions(t, param, gp) {
		return false
	}
	if !ctx.assertGenParamValidations(t, param, gp) {
		return false
	}
	if !ctx.assertGenParamInternals(t, param, gp) {
		return false
	}

	return true
}

func (ctx *paramTestContext) assertGenParamDefinitions(t *testing.T, param spec.Parameter, gp GenParameter) bool {
	t.Helper()

	if !assert.Equal(t, param.In, gp.Location) {
		return false
	}
	if !assert.Equal(t, param.Name, gp.Name) {
		return false
	}
	if !assert.Equal(t, fmt.Sprintf("%q", param.Name), gp.Path) {
		return false
	}
	if !assert.Equal(t, param.Description, gp.Description) {
		return false
	}
	if !assert.Equal(t, param.CollectionFormat, gp.CollectionFormat) {
		return false
	}

	if param.In == body {
		return assertBodyParam(t, param, gp)
	}

	if ctx.Items != nil {
		return ctx.Items.Assert(t, param.Items, gp.Child)
	}

	return true
}

func (ctx *paramTestContext) assertGenParamValidations(t *testing.T, param spec.Parameter, gp GenParameter) bool {
	t.Helper()

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
	if !assert.Equal(t, param.Enum, gp.Enum) {
		return false
	}
	if !assert.Equal(t, param.Type, gp.SwaggerType) {
		return false
	}
	if !assert.Equal(t, param.Format, gp.SwaggerFormat) {
		return false
	}

	return true
}

func (ctx *paramTestContext) assertGenParamInternals(t *testing.T, param spec.Parameter, gp GenParameter) bool {
	t.Helper()

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

	if _, ok := primitives[gp.GoType]; ok {
		if !assert.True(t, gp.IsPrimitive) {
			return false
		}
	} else {
		if !assert.False(t, gp.IsPrimitive) {
			return false
		}
	}

	return true
}

func (ctx *paramItemsTestContext) Assert(t *testing.T, pItems *spec.Items, gpItems *GenItems) bool {
	t.Helper()

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
	if !assert.Equal(t, pItems.Enum, gpItems.Enum) {
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

func assertBodyParam(t *testing.T, param spec.Parameter, gp GenParameter) bool {
	t.Helper()

	if !assert.Equal(t, body, param.In) || !assert.Equal(t, body, gp.Location) {
		return false
	}
	if !assert.NotNil(t, gp.Schema) {
		return false
	}

	return true
}
