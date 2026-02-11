// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/spec"
)

func getResponse(sctx *scanCtx, nm string) *entityDecl {
	for _, v := range sctx.app.Responses {
		if v.Ident.Name == nm {
			return v
		}
	}
	return nil
}

func TestParseResponses(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	responses := make(map[string]spec.Response)
	for _, rn := range []string{"ComplexerOne", "SimpleOnes", "SimpleOnesFunc", "ComplexerPointerOne", "SomeResponse", "ValidationError", "Resp", "FileResponse", "GenericError", "ValidationError"} {
		td := getResponse(sctx, rn)
		prs := &responseBuilder{
			ctx:  sctx,
			decl: td,
		}
		require.NoError(t, prs.Build(responses))
	}

	require.Len(t, responses, 9)
	cr, ok := responses["complexerOne"]
	assert.TrueT(t, ok)
	assert.Len(t, cr.Headers, 7)
	for k, header := range cr.Headers {
		switch k {
		case "id":
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int64", header.Format)
		case "name":
			assert.EqualT(t, "string", header.Type)
			assert.Empty(t, header.Format)
		case "age":
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int32", header.Format)
		case "notes":
			assert.EqualT(t, "string", header.Type)
			assert.Empty(t, header.Format)
		case "extra":
			assert.EqualT(t, "string", header.Type)
			assert.Empty(t, header.Format)
		case "createdAt":
			assert.EqualT(t, "string", header.Type)
			assert.EqualT(t, "date-time", header.Format)
		case "NoTagName":
			assert.EqualT(t, "string", header.Type)
			assert.Empty(t, header.Format)
		default:
			assert.Fail(t, "unknown header: "+k)
		}
	}

	cpr, ok := responses["complexerPointerOne"]
	assert.TrueT(t, ok)
	assert.Len(t, cpr.Headers, 4)
	for k, header := range cpr.Headers {
		switch k {
		case "id":
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int64", header.Format)
		case "name":
			assert.EqualT(t, "string", header.Type)
			assert.Empty(t, header.Format)
		case "age":
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int32", header.Format)
		case "extra":
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int64", header.Format)
		default:
			assert.Fail(t, "unknown header: "+k)
		}
	}

	sos, ok := responses["simpleOnes"]
	assert.TrueT(t, ok)
	assert.Len(t, sos.Headers, 1)

	sosf, ok := responses["simpleOnesFunc"]
	assert.TrueT(t, ok)
	assert.Len(t, sosf.Headers, 1)

	res, ok := responses["someResponse"]
	assert.TrueT(t, ok)
	assert.Len(t, res.Headers, 7)

	for k, header := range res.Headers {
		switch k {
		case "id":
			assert.EqualT(t, "ID of this some response instance.\nids in this application start at 11 and are smaller than 1000", header.Description)
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int64", header.Format)
			// assert.Equal(t, "ID", header.Extensions["x-go-name"])
			require.NotNil(t, header.Maximum)
			assert.InDeltaT(t, 1000.00, *header.Maximum, epsilon)
			assert.TrueT(t, header.ExclusiveMaximum)
			require.NotNil(t, header.Minimum)
			assert.InDeltaT(t, 10.00, *header.Minimum, epsilon)
			assert.TrueT(t, header.ExclusiveMinimum)
			assert.EqualValues(t, 11, header.Default, "ID default value is incorrect")

		case "score":
			assert.EqualT(t, "The Score of this model", header.Description)
			assert.EqualT(t, "integer", header.Type)
			assert.EqualT(t, "int32", header.Format)
			// assert.Equal(t, "Score", header.Extensions["x-go-name"])
			require.NotNil(t, header.Maximum)
			assert.InDeltaT(t, 45.00, *header.Maximum, epsilon)
			assert.FalseT(t, header.ExclusiveMaximum)
			require.NotNil(t, header.Minimum)
			assert.InDeltaT(t, 3.00, *header.Minimum, epsilon)
			assert.FalseT(t, header.ExclusiveMinimum)
			assert.EqualValues(t, 27, header.Example)

		case "x-hdr-name":
			assert.EqualT(t, "Name of this some response instance", header.Description)
			assert.EqualT(t, "string", header.Type)
			// assert.Equal(t, "Name", header.Extensions["x-go-name"])
			require.NotNil(t, header.MinLength)
			assert.InDelta(t, 4.00, *header.MinLength, epsilon)
			require.NotNil(t, header.MaxLength)
			assert.InDelta(t, 50.00, *header.MaxLength, epsilon)
			assert.EqualT(t, "[A-Za-z0-9-.]*", header.Pattern)

		case "active":
			assert.EqualT(t, "Active state of the record", header.Description)
			assert.EqualT(t, "boolean", header.Type)
			active, ok2 := header.Default.(bool)
			assert.TrueT(t, ok2)
			assert.TrueT(t, active)

		case "created":
			assert.EqualT(t, "Created holds the time when this entry was created", header.Description)
			assert.EqualT(t, "string", header.Type)
			assert.EqualT(t, "date-time", header.Format)
			// assert.Equal(t, "Created", header.Extensions["x-go-name"])

		case "foo_slice":
			assert.EqualT(t, "a FooSlice has foos which are strings", header.Description)
			// assert.Equal(t, "FooSlice", header.Extensions["x-go-name"])
			assert.EqualT(t, "array", header.Type)
			assert.TrueT(t, header.UniqueItems)
			assert.EqualT(t, "pipe", header.CollectionFormat)
			assert.NotNil(t, header.Items, "foo_slice should have had an items property")
			require.NotNil(t, header.MinItems)
			assert.EqualT(t, int64(3), *header.MinItems, "'foo_slice' should have had 3 min items")
			require.NotNil(t, header.MaxItems)
			assert.EqualT(t, int64(10), *header.MaxItems, "'foo_slice' should have had 10 max items")

			itprop := header.Items
			require.NotNil(t, itprop.MinLength)
			assert.EqualT(t, int64(3), *itprop.MinLength, "'foo_slice.items.minLength' should have been 3")
			require.NotNil(t, itprop.MaxLength)
			assert.EqualT(t, int64(10), *itprop.MaxLength, "'foo_slice.items.maxLength' should have been 10")
			assert.EqualT(t, "\\w+", itprop.Pattern, "'foo_slice.items.pattern' should have \\w+")
			assert.Equal(t, "foo", itprop.Example)

		case "bar_slice":
			assert.EqualT(t, "a BarSlice has bars which are strings", header.Description)
			assert.EqualT(t, "array", header.Type)
			assert.TrueT(t, header.UniqueItems)
			assert.EqualT(t, "pipe", header.CollectionFormat)
			require.NotNil(t, header.Items, "bar_slice should have had an items property")
			require.NotNil(t, header.MinItems)
			assert.EqualT(t, int64(3), *header.MinItems, "'bar_slice' should have had 3 min items")
			require.NotNil(t, header.MaxItems)
			assert.EqualT(t, int64(10), *header.MaxItems, "'bar_slice' should have had 10 max items")

			itprop := header.Items
			require.NotNil(t, itprop)
			require.NotNil(t, itprop.MinItems)
			assert.EqualT(t, int64(4), *itprop.MinItems, "'bar_slice.items.minItems' should have been 4")
			require.NotNil(t, itprop.MaxItems)
			assert.EqualT(t, int64(9), *itprop.MaxItems, "'bar_slice.items.maxItems' should have been 9")

			itprop2 := itprop.Items
			require.NotNil(t, itprop2)
			require.NotNil(t, itprop.MinItems)
			assert.EqualT(t, int64(5), *itprop2.MinItems, "'bar_slice.items.items.minItems' should have been 5")
			require.NotNil(t, itprop.MaxItems)
			assert.EqualT(t, int64(8), *itprop2.MaxItems, "'bar_slice.items.items.maxItems' should have been 8")

			itprop3 := itprop2.Items
			require.NotNil(t, itprop3)
			require.NotNil(t, itprop3.MinLength)
			assert.EqualT(t, int64(3), *itprop3.MinLength, "'bar_slice.items.items.items.minLength' should have been 3")
			require.NotNil(t, itprop3.MaxLength)
			assert.EqualT(t, int64(10), *itprop3.MaxLength, "'bar_slice.items.items.items.maxLength' should have been 10")
			assert.EqualT(t, "\\w+", itprop3.Pattern, "'bar_slice.items.items.items.pattern' should have \\w+")

		default:
			assert.Fail(t, "unknown property: "+k)
		}
	}

	assert.NotNil(t, res.Schema)
	aprop := res.Schema
	assert.EqualT(t, "array", aprop.Type[0])
	require.NotNil(t, aprop.Items)
	require.NotNil(t, aprop.Items.Schema)

	itprop := aprop.Items.Schema
	assert.Len(t, itprop.Properties, 4)
	assert.Len(t, itprop.Required, 3)
	assertProperty(t, itprop, "integer", "id", "int32", "ID")

	iprop, ok := itprop.Properties["id"]
	assert.TrueT(t, ok)
	assert.EqualT(t, "ID of this some response instance.\nids in this application start at 11 and are smaller than 1000", iprop.Description)
	require.NotNil(t, iprop.Maximum)
	assert.InDeltaT(t, 1000.00, *iprop.Maximum, epsilon)
	assert.TrueT(t, iprop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	require.NotNil(t, iprop.Minimum)
	assert.InDeltaT(t, 10.00, *iprop.Minimum, epsilon)
	assert.TrueT(t, iprop.ExclusiveMinimum, "'id' should have had an exclusive minimum")

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

	res, ok = responses["resp"]
	assert.TrueT(t, ok)
	assert.NotNil(t, res.Schema)
	assert.EqualT(t, "#/definitions/user", res.Schema.Ref.String())
}

func TestParseResponses_TransparentAliases(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages:           []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/transparentalias"},
		TransparentAliases: true,
		ScanModels:         true,
	})
	require.NoError(t, err)

	td := getResponse(sctx, "TransparentAliasResponse")
	require.NotNil(t, td)

	// Build the response map using the transparent alias fixtures.
	responses := make(map[string]spec.Response)
	prs := &responseBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(responses))

	resp, ok := responses["transparentAliasResponse"]
	require.TrueT(t, ok)
	require.NotNil(t, resp.Schema)
	assert.TrueT(t, resp.Schema.Type.Contains("object"))
	assert.Empty(t, resp.Schema.Ref.String())

	payload, ok := resp.Schema.Properties["payload"]
	require.TrueT(t, ok)
	// The response payload alias should expand inline and retain field metadata.
	assert.TrueT(t, payload.Type.Contains("object"))
	assert.Empty(t, payload.Ref.String())
	assert.Equal(t, "Payload", payload.Extensions["x-go-name"])
}

func TestParseResponses_Issue2007(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	responses := make(map[string]spec.Response)
	td := getResponse(sctx, "GetConfiguration")
	prs := &responseBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(responses))

	resp := responses["GetConfiguration"]
	require.Empty(t, resp.Headers)
	require.NotNil(t, resp.Schema)

	require.TrueT(t, resp.Schema.Type.Contains("object"))
	require.NotNil(t, resp.Schema.AdditionalProperties)
	require.NotNil(t, resp.Schema.AdditionalProperties.Schema)
	require.TrueT(t, resp.Schema.AdditionalProperties.Schema.Type.Contains("string"))
}

func TestParseResponses_Issue2011(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	responses := make(map[string]spec.Response)
	td := getResponse(sctx, "NumPlatesResp")
	prs := &responseBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(responses))

	resp := responses["NumPlatesResp"]
	require.Empty(t, resp.Headers)
	require.NotNil(t, resp.Schema)
}

func TestParseResponses_Issue2145(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages: []string{"github.com/go-swagger/go-swagger/fixtures/goparsing/product/..."},
	})
	require.NoError(t, err)
	responses := make(map[string]spec.Response)
	td := getResponse(sctx, "GetProductsResponse")
	prs := &responseBuilder{
		ctx:  sctx,
		decl: td,
	}
	require.NoError(t, prs.Build(responses))
	resp := responses["GetProductsResponse"]
	require.Empty(t, resp.Headers)
	require.NotNil(t, resp.Schema)

	assert.NotEmpty(t, prs.postDecls) // should have Product
}
