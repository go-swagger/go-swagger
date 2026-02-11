// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/spec"
)

const (
	gcBadEnum = "bad_enum"
)

func getParameter(sctx *scanCtx, nm string) *entityDecl {
	for _, v := range sctx.app.Parameters {
		param := v
		if v.Ident.Name == nm {
			return param
		}
	}
	return nil
}

func TestScanFileParam(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	operations := make(map[string]*spec.Operation)
	for _, rn := range []string{"OrderBodyParams", "MultipleOrderParams", "ComplexerOneParams", "NoParams", "NoParamsAlias", "MyFileParams", "MyFuncFileParams", "EmbeddedFileParams"} {
		td := getParameter(sctx, rn)

		prs := &parameterBuilder{
			ctx:  sctx,
			decl: td,
		}
		require.NoError(t, prs.Build(operations))
	}
	assert.Len(t, operations, 10)

	of, ok := operations["myOperation"]
	assert.TrueT(t, ok)
	assert.Len(t, of.Parameters, 1)
	fileParam := of.Parameters[0]
	assert.EqualT(t, "MyFormFile desc.", fileParam.Description)
	assert.EqualT(t, "formData", fileParam.In)
	assert.EqualT(t, "file", fileParam.Type)
	assert.FalseT(t, fileParam.Required)

	emb, ok := operations["myOtherOperation"]
	assert.TrueT(t, ok)
	assert.Len(t, emb.Parameters, 2)
	fileParam = emb.Parameters[0]
	assert.EqualT(t, "MyFormFile desc.", fileParam.Description)
	assert.EqualT(t, "formData", fileParam.In)
	assert.EqualT(t, "file", fileParam.Type)
	assert.FalseT(t, fileParam.Required)
	extraParam := emb.Parameters[1]
	assert.EqualT(t, "ExtraParam desc.", extraParam.Description)
	assert.EqualT(t, "formData", extraParam.In)
	assert.EqualT(t, "integer", extraParam.Type)
	assert.TrueT(t, extraParam.Required)

	ffp, ok := operations["myFuncOperation"]
	assert.TrueT(t, ok)
	assert.Len(t, ffp.Parameters, 1)
	fileParam = ffp.Parameters[0]
	assert.EqualT(t, "MyFormFile desc.", fileParam.Description)
	assert.EqualT(t, "formData", fileParam.In)
	assert.EqualT(t, "file", fileParam.Type)
	assert.FalseT(t, fileParam.Required)
}

func TestParamsParser(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	operations := make(map[string]*spec.Operation)
	for _, rn := range []string{"OrderBodyParams", "MultipleOrderParams", "ComplexerOneParams", "NoParams", "NoParamsAlias", "MyFileParams", "MyFuncFileParams", "EmbeddedFileParams"} {
		td := getParameter(sctx, rn)

		prs := &parameterBuilder{
			ctx:  sctx,
			decl: td,
		}
		require.NoError(t, prs.Build(operations))
	}

	assert.Len(t, operations, 10)
	cr, okParam := operations["yetAnotherOperation"]
	assert.TrueT(t, okParam)
	assert.Len(t, cr.Parameters, 8)
	for _, param := range cr.Parameters {
		switch param.Name {
		case "id":
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int64", param.Format)
		case "name":
			assert.EqualT(t, "string", param.Type)
			assert.Empty(t, param.Format)
		case "age":
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int32", param.Format)
		case "notes":
			assert.EqualT(t, "string", param.Type)
			assert.Empty(t, param.Format)
		case "extra":
			assert.EqualT(t, "string", param.Type)
			assert.Empty(t, param.Format)
		case "createdAt":
			assert.EqualT(t, "string", param.Type)
			assert.EqualT(t, "date-time", param.Format)
		case "informity":
			assert.EqualT(t, "string", param.Type)
			assert.EqualT(t, "formData", param.In)
		case "NoTagName":
			assert.EqualT(t, "string", param.Type)
			assert.Empty(t, param.Format)
		default:
			assert.Fail(t, "unknown property: "+param.Name)
		}
	}

	ob, okParam := operations["updateOrder"]
	assert.TrueT(t, okParam)
	assert.Len(t, ob.Parameters, 1)
	bodyParam := ob.Parameters[0]
	assert.EqualT(t, "The order to submit.", bodyParam.Description)
	assert.EqualT(t, "body", bodyParam.In)
	assert.EqualT(t, "#/definitions/order", bodyParam.Schema.Ref.String())
	assert.TrueT(t, bodyParam.Required)

	mop, okParam := operations["getOrders"]
	assert.TrueT(t, okParam)
	require.Len(t, mop.Parameters, 2)
	ordersParam := mop.Parameters[0]
	assert.EqualT(t, "The orders", ordersParam.Description)
	assert.TrueT(t, ordersParam.Required)
	assert.EqualT(t, "array", ordersParam.Type)
	otherParam := mop.Parameters[1]
	assert.EqualT(t, "And another thing", otherParam.Description)

	op, okParam := operations["someOperation"]
	assert.TrueT(t, okParam)
	assert.Len(t, op.Parameters, 12)

	for _, param := range op.Parameters {
		switch param.Name {
		case "id":
			assert.EqualT(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", param.Description)
			assert.EqualT(t, "path", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int64", param.Format)
			assert.TrueT(t, param.Required)
			assert.Equal(t, "ID", param.Extensions["x-go-name"])
			require.NotNil(t, param.Maximum)
			assert.InDeltaT(t, 1000.00, *param.Maximum, epsilon)
			assert.TrueT(t, param.ExclusiveMaximum)
			require.NotNil(t, param.Minimum)
			assert.InDeltaT(t, 10.00, *param.Minimum, epsilon)
			assert.TrueT(t, param.ExclusiveMinimum)
			assert.Equal(t, 1, param.Default, "%s default value is incorrect", param.Name)

		case "score":
			assert.EqualT(t, "The Score of this model", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int32", param.Format)
			assert.TrueT(t, param.Required)
			assert.Equal(t, "Score", param.Extensions["x-go-name"])
			require.NotNil(t, param.Maximum)
			assert.InDeltaT(t, 45.00, *param.Maximum, epsilon)
			assert.FalseT(t, param.ExclusiveMaximum)
			require.NotNil(t, param.Minimum)
			assert.InDeltaT(t, 3.00, *param.Minimum, epsilon)
			assert.FalseT(t, param.ExclusiveMinimum)
			assert.EqualValues(t, 2, param.Default, "%s default value is incorrect", param.Name)
			assert.EqualValues(t, 27, param.Example)

		case "x-hdr-name":
			assert.EqualT(t, "Name of this no model instance", param.Description)
			assert.EqualT(t, "header", param.In)
			assert.EqualT(t, "string", param.Type)
			assert.TrueT(t, param.Required)
			assert.Equal(t, "Name", param.Extensions["x-go-name"])
			require.NotNil(t, param.MinLength)
			assert.EqualT(t, int64(4), *param.MinLength)
			require.NotNil(t, param.MaxLength)
			assert.EqualT(t, int64(50), *param.MaxLength)
			assert.EqualT(t, "[A-Za-z0-9-.]*", param.Pattern)

		case "created":
			assert.EqualT(t, "Created holds the time when this entry was created", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "string", param.Type)
			assert.EqualT(t, "date-time", param.Format)
			assert.FalseT(t, param.Required)
			assert.Equal(t, "Created", param.Extensions["x-go-name"])

		case "category_old":
			assert.EqualT(t, "The Category of this model (old enum format)", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "string", param.Type)
			assert.TrueT(t, param.Required)
			assert.Equal(t, "CategoryOld", param.Extensions["x-go-name"])
			assert.Equal(t, []any{"foo", "bar", "none"}, param.Enum, "%s enum values are incorrect", param.Name)
			assert.Equal(t, "bar", param.Default, "%s default value is incorrect", param.Name)
		case "category":
			assert.EqualT(t, "The Category of this model", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "string", param.Type)
			assert.TrueT(t, param.Required)
			assert.Equal(t, "Category", param.Extensions["x-go-name"])
			assert.Equal(t, []any{"foo", "bar", "none"}, param.Enum, "%s enum values are incorrect", param.Name)
			assert.Equal(t, "bar", param.Default, "%s default value is incorrect", param.Name)
		case "type_old":
			assert.EqualT(t, "Type of this model (old enum format)", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.Equal(t, []any{1, 3, 5}, param.Enum, "%s enum values are incorrect", param.Name)
		case "type":
			assert.EqualT(t, "Type of this model", param.Description)
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.Equal(t, []any{1, 3, 5}, param.Enum, "%s enum values are incorrect", param.Name)
		case gcBadEnum:
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.Equal(t, []any{1, "rsq", "qaz"}, param.Enum, "%s enum values are incorrect", param.Name)
		case "foo_slice":
			assert.EqualT(t, "a FooSlice has foos which are strings", param.Description)
			assert.Equal(t, "FooSlice", param.Extensions["x-go-name"])
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "array", param.Type)
			assert.FalseT(t, param.Required)
			assert.TrueT(t, param.UniqueItems)
			assert.EqualT(t, "pipe", param.CollectionFormat)
			require.NotNil(t, param.MinItems)
			assert.EqualT(t, int64(3), *param.MinItems, "'foo_slice' should have had 3 min items")
			require.NotNil(t, param.MaxItems)
			assert.EqualT(t, int64(10), *param.MaxItems, "'foo_slice' should have had 10 max items")

			itprop := param.Items
			require.NotNil(t, itprop.MinLength)
			assert.EqualT(t, int64(3), *itprop.MinLength, "'foo_slice.items.minLength' should have been 3")
			require.NotNil(t, itprop.MaxLength)
			assert.EqualT(t, int64(10), *itprop.MaxLength, "'foo_slice.items.maxLength' should have been 10")
			assert.EqualT(t, "\\w+", itprop.Pattern, "'foo_slice.items.pattern' should have \\w+")
			assert.EqualValues(t, "bar", itprop.Default, "'foo_slice.items.default' should have bar default value")

		case "items":
			assert.Equal(t, "Items", param.Extensions["x-go-name"])
			assert.EqualT(t, "body", param.In)
			assert.NotNil(t, param.Schema)
			aprop := param.Schema
			assert.EqualT(t, "array", aprop.Type[0])
			require.NotNil(t, aprop.Items)
			require.NotNil(t, aprop.Items.Schema)

			itprop := aprop.Items.Schema
			assert.Len(t, itprop.Properties, 4)
			assert.Len(t, itprop.Required, 3)
			assertProperty(t, itprop, "integer", "id", "int32", "ID")
			iprop, ok := itprop.Properties["id"]
			assert.TrueT(t, ok)
			assert.EqualT(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", iprop.Description)
			require.NotNil(t, iprop.Maximum)
			assert.InDeltaT(t, 1000.00, *iprop.Maximum, epsilon)
			assert.TrueT(t, iprop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
			require.NotNil(t, iprop.Minimum)
			assert.InDeltaT(t, 10.00, *iprop.Minimum, epsilon)
			assert.TrueT(t, iprop.ExclusiveMinimum, "'id' should have had an exclusive minimum")
			assert.Equal(t, 3, iprop.Default, "Items.ID default value is incorrect")

			assertRef(t, itprop, "pet", "Pet", "#/definitions/pet")
			_, ok = itprop.Properties["pet"]
			assert.TrueT(t, ok)
			// if itprop.Ref.String() == "" {
			// 	assert.Equal(t, "The Pet to add to this NoModel items bucket.\nPets can appear more than once in the bucket", iprop.Description)
			// }

			assertProperty(t, itprop, "integer", "quantity", "int16", "Quantity")
			iprop, ok = itprop.Properties["quantity"]
			assert.TrueT(t, ok)
			assert.EqualT(t, "The amount of pets to add to this bucket.", iprop.Description)
			require.NotNil(t, iprop.Minimum)
			assert.InDeltaT(t, 1.00, *iprop.Minimum, epsilon)
			require.NotNil(t, iprop.Maximum)
			assert.InDeltaT(t, 10.00, *iprop.Maximum, epsilon)

			assertProperty(t, itprop, "string", "notes", "", "Notes")
			iprop, ok = itprop.Properties["notes"]
			assert.TrueT(t, ok)
			assert.EqualT(t, "Notes to add to this item.\nThis can be used to add special instructions.", iprop.Description)

		case "bar_slice":
			assert.EqualT(t, "a BarSlice has bars which are strings", param.Description)
			assert.Equal(t, "BarSlice", param.Extensions["x-go-name"])
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "array", param.Type)
			assert.FalseT(t, param.Required)
			assert.TrueT(t, param.UniqueItems)
			assert.EqualT(t, "pipe", param.CollectionFormat)
			require.NotNil(t, param.Items, "bar_slice should have had an items property")
			require.NotNil(t, param.MinItems)
			assert.EqualT(t, int64(3), *param.MinItems, "'bar_slice' should have had 3 min items")
			require.NotNil(t, param.MaxItems)
			assert.EqualT(t, int64(10), *param.MaxItems, "'bar_slice' should have had 10 max items")

			itprop := param.Items
			require.NotNil(t, itprop)
			require.NotNil(t, itprop.MinItems)
			assert.EqualT(t, int64(4), *itprop.MinItems, "'bar_slice.items.minItems' should have been 4")
			require.NotNil(t, itprop.MaxItems)
			assert.EqualT(t, int64(9), *itprop.MaxItems, "'bar_slice.items.maxItems' should have been 9")

			itprop2 := itprop.Items
			require.NotNil(t, itprop2)
			require.NotNil(t, itprop2.MinItems)
			assert.EqualT(t, int64(5), *itprop2.MinItems, "'bar_slice.items.items.minItems' should have been 5")
			require.NotNil(t, itprop2.MaxItems)
			assert.EqualT(t, int64(8), *itprop2.MaxItems, "'bar_slice.items.items.maxItems' should have been 8")

			itprop3 := itprop2.Items
			require.NotNil(t, itprop3)
			require.NotNil(t, itprop3.MinLength)
			assert.EqualT(t, int64(3), *itprop3.MinLength, "'bar_slice.items.items.items.minLength' should have been 3")
			require.NotNil(t, itprop3.MaxLength)
			assert.EqualT(t, int64(10), *itprop3.MaxLength, "'bar_slice.items.items.items.maxLength' should have been 10")
			assert.EqualT(t, "\\w+", itprop3.Pattern, "'bar_slice.items.items.items.pattern' should have \\w+")

		default:
			assert.Fail(t, "unknown property: "+param.Name)
		}
	}

	// assert that the order of the parameters is maintained
	order, ok := operations["anotherOperation"]
	assert.TrueT(t, ok)
	assert.Len(t, order.Parameters, 12)

	for index, param := range order.Parameters {
		switch param.Name {
		case "id":
			assert.EqualT(t, 0, index, "%s index incorrect", param.Name)
		case "score":
			assert.EqualT(t, 1, index, "%s index incorrect", param.Name)
		case "x-hdr-name":
			assert.EqualT(t, 2, index, "%s index incorrect", param.Name)
		case "created":
			assert.EqualT(t, 3, index, "%s index incorrect", param.Name)
		case "category_old":
			assert.EqualT(t, 4, index, "%s index incorrect", param.Name)
		case "category":
			assert.EqualT(t, 5, index, "%s index incorrect", param.Name)
		case "type_old":
			assert.EqualT(t, 6, index, "%s index incorrect", param.Name)
		case "type":
			assert.EqualT(t, 7, index, "%s index incorrect", param.Name)
		case gcBadEnum:
			assert.EqualT(t, 8, index, "%s index incorrect", param.Name)
		case "foo_slice":
			assert.EqualT(t, 9, index, "%s index incorrect", param.Name)
		case "bar_slice":
			assert.EqualT(t, 10, index, "%s index incorrect", param.Name)
		case "items":
			assert.EqualT(t, 11, index, "%s index incorrect", param.Name)
		default:
			assert.Fail(t, "unknown property: "+param.Name)
		}
	}

	// check that aliases work correctly
	aliasOp, ok := operations["someAliasOperation"]
	assert.TrueT(t, ok)
	assert.Len(t, aliasOp.Parameters, 4)
	for _, param := range aliasOp.Parameters {
		switch param.Name {
		case "intAlias":
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int64", param.Format)
			assert.TrueT(t, param.Required)
			require.NotNil(t, param.Maximum)
			assert.InDeltaT(t, 10.00, *param.Maximum, epsilon)
			require.NotNil(t, param.Minimum)
			assert.InDeltaT(t, 1.00, *param.Minimum, epsilon)
		case "stringAlias":
			assert.EqualT(t, "query", param.In)
			assert.EqualT(t, "string", param.Type)
		case "intAliasPath":
			assert.EqualT(t, "path", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int64", param.Format)
		case "intAliasForm":
			assert.EqualT(t, "formData", param.In)
			assert.EqualT(t, "integer", param.Type)
			assert.EqualT(t, "int64", param.Format)
		default:
			assert.Fail(t, "unknown property: "+param.Name)
		}
	}
}

func TestParamsParser_TransparentAliases(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages:           []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/transparentalias"},
		TransparentAliases: true,
		ScanModels:         true,
	})
	require.NoError(t, err)

	td := getParameter(sctx, "TransparentAliasParams")
	require.NotNil(t, td)

	// Build the operation map from the transparent alias fixtures.
	operations := make(map[string]*spec.Operation)
	prs := &parameterBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(operations))

	op, ok := operations["transparentAlias"]
	require.TrueT(t, ok)
	require.Len(t, op.Parameters, 2)

	var bodyParam, queryParam *spec.Parameter
	for i := range op.Parameters {
		p := &op.Parameters[i]
		switch p.In {
		case "body":
			bodyParam = p
		case "query":
			queryParam = p
		}
	}

	require.NotNil(t, bodyParam)
	require.NotNil(t, queryParam)
	require.NotNil(t, bodyParam.Schema)

	// Body aliases should expand inline instead of producing a $ref definition.
	assert.EqualT(t, "aliasBody", bodyParam.Name)
	assert.TrueT(t, bodyParam.Schema.Type.Contains("object"))
	assert.Empty(t, bodyParam.Schema.Ref.String())
	assert.SliceContainsT(t, bodyParam.Schema.Required, "id")
	idSchema, ok := bodyParam.Schema.Properties["id"]
	require.TrueT(t, ok)
	assert.TrueT(t, idSchema.Type.Contains("integer"))
	assert.Equal(t, "ID", idSchema.Extensions["x-go-name"])
	// Query aliases should behave like their underlying scalar type.
	assert.EqualT(t, "aliasQuery", queryParam.Name)
	assert.EqualT(t, "string", queryParam.Type)
	assert.Empty(t, queryParam.Ref.String())
}

func TestParameterParser_Issue2007(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	operations := make(map[string]*spec.Operation)
	td := getParameter(sctx, "SetConfiguration")
	prs := &parameterBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(operations))

	op := operations["getConfiguration"]
	require.NotNil(t, op)
	require.Len(t, op.Parameters, 1)
	sch := op.Parameters[0].Schema
	require.NotNil(t, sch)

	require.TrueT(t, sch.Type.Contains("object"))
	require.NotNil(t, sch.AdditionalProperties)
	require.NotNil(t, sch.AdditionalProperties.Schema)
	require.TrueT(t, sch.AdditionalProperties.Schema.Type.Contains("string"))
}

func TestParameterParser_Issue2011(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	operations := make(map[string]*spec.Operation)
	td := getParameter(sctx, "NumPlates")
	prs := &parameterBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(operations))

	op := operations["putNumPlate"]
	require.NotNil(t, op)
	require.Len(t, op.Parameters, 1)
	sch := op.Parameters[0].Schema
	require.NotNil(t, sch)
}
