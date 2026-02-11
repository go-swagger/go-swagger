// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package codescan

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/testify/v2/assert"

	"github.com/go-openapi/spec"
)

func TestRouteExpression(t *testing.T) {
	assert.RegexpT(t, rxRoute, "swagger:route DELETE /orders/{id} deleteOrder")
	assert.RegexpT(t, rxRoute, "swagger:route GET /v1.2/something deleteOrder")
}

func TestRoutesParser(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	var ops spec.Paths
	for _, apiPath := range sctx.app.Routes {
		prs := &routesBuilder{
			ctx:        sctx,
			route:      apiPath,
			operations: make(map[string]*spec.Operation),
		}
		require.NoError(t, prs.Build(&ops))
	}

	assert.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	ext := make(spec.Extensions)
	ext.Add("x-some-flag", "true")
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
		ext,
	)
	assertOperation(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
		nil,
	)

	po, ok = ops.Paths["/orders"]
	ext = make(spec.Extensions)
	ext.Add("x-some-flag", "false")
	ext.Add("x-some-list", []string{"item1", "item2", "item3"})
	ext.Add("x-some-object", map[string]any{
		"key1": "value1",
		"key2": "value2",
		"subobject": map[string]any{
			"subkey1": "subvalue1",
			"subkey2": "subvalue2",
		},
		"key3": "value3",
	})
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
		ext,
	)
	assertOperation(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
		nil,
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperation(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
		nil,
	)

	assertOperation(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
		nil,
	)

	assertOperation(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
		nil,
	)

	// additional check description tag at Responses
	rsp, ok := po.Delete.Responses.StatusCodeResponses[202]
	assert.TrueT(t, ok)
	assert.EqualT(t, "Some description", rsp.Description)
	assert.Empty(t, rsp.Ref.String())
}

func TestRoutesParserBody(t *testing.T) {
	sctx, err := newScanCtx(&Options{
		Packages: []string{
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/operations",
			"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/operations_body",
		},
	})
	require.NoError(t, err)
	var ops spec.Paths
	for _, apiPath := range sctx.app.Routes {
		prs := &routesBuilder{
			ctx:        sctx,
			route:      apiPath,
			operations: make(map[string]*spec.Operation),
		}
		require.NoError(t, prs.Build(&ops))
	}

	assert.Len(t, ops.Paths, 4)

	po, ok := ops.Paths["/pets"]
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)
	assert.NotNil(t, po.Post)

	assertOperationBody(t,
		po.Post,
		"createPet",
		"Create a pet based on the parameters.",
		"",
		[]string{"pets", "users"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders"]
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"listOrders",
		"lists orders filtered by some parameters.",
		"",
		[]string{"orders"},
		[]string{"orders:read", "https://www.googleapis.com/auth/userinfo.email"},
	)
	assert.NotNil(t, po.Post)

	assertOperationBody(t,
		po.Post,
		"createOrder",
		"create an order based on the parameters.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	po, ok = ops.Paths["/orders/{id}"]
	assert.TrueT(t, ok)
	assert.NotNil(t, po.Get)
	assertOperationBody(t,
		po.Get,
		"orderDetails",
		"gets the details for an order.",
		"",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Put,
		"updateOrder",
		"Update the details for an order.",
		"When the order doesn't exist this will return an error.",
		[]string{"orders"},
		[]string{"read", "write"},
	)

	assertOperationBody(t,
		po.Delete,
		"deleteOrder",
		"delete a particular order.",
		"",
		nil,
		[]string{"read", "write"},
	)

	validateRoutesParameters(t, ops)
}

func validateRoutesParameters(t *testing.T, ops spec.Paths) {
	po := ops.Paths["/pets"]
	assert.Len(t, po.Post.Parameters, 2)

	// Testing standard param properties
	p := po.Post.Parameters[0]
	assert.EqualT(t, "request", p.Name)
	assert.EqualT(t, "body", p.In)
	assert.EqualT(t, "The request model.", p.Description)

	// Testing "required" and "allowEmpty"
	p = po.Post.Parameters[1]
	assert.EqualT(t, "id", p.Name)
	assert.EqualT(t, "The pet id", p.Description)
	assert.EqualT(t, "path", p.In)
	assert.TrueT(t, p.Required)
	assert.FalseT(t, p.AllowEmptyValue)

	po = ops.Paths["/orders"]
	assert.Len(t, po.Post.Parameters, 2)

	// Testing invalid value for "in"
	p = po.Post.Parameters[0]
	assert.EqualT(t, "id", p.Name)
	assert.EqualT(t, "The order id", p.Description)
	assert.Empty(t, p.In) // Invalid value should not be set
	assert.FalseT(t, p.Required)
	assert.TrueT(t, p.AllowEmptyValue)

	p = po.Post.Parameters[1]
	assert.EqualT(t, "request", p.Name)
	assert.EqualT(t, "body", p.In)
	assert.EqualT(t, "The request model.", p.Description)

	po = ops.Paths["/param-test"]
	assert.Len(t, po.Post.Parameters, 6)

	// Testing number param with "max" and "min" constraints
	p = po.Post.Parameters[0]
	assert.EqualT(t, "someNumber", p.Name)
	assert.EqualT(t, "some number", p.Description)
	assert.EqualT(t, "path", p.In)
	assert.TrueT(t, p.Required)
	assert.EqualT(t, "number", p.Type)
	minimum, maximum, def := float64(10), float64(20), float64(15)
	assert.Equal(t, &maximum, p.Maximum)
	assert.Equal(t, &minimum, p.Minimum)
	assert.InDelta(t, def, p.Default, epsilon)
	assert.Nil(t, p.Schema)

	// Testing array param provided as query string. Testing "minLength" and "maxLength" constraints for "array" types
	p = po.Post.Parameters[1]
	assert.EqualT(t, "someQuery", p.Name)
	assert.EqualT(t, "some query values", p.Description)
	assert.EqualT(t, "query", p.In)
	assert.FalseT(t, p.Required)
	assert.EqualT(t, "array", p.Type)
	minLen, maxLen := int64(5), int64(20)
	assert.Equal(t, &maxLen, p.MaxLength)
	assert.Equal(t, &minLen, p.MinLength)
	assert.Nil(t, p.Schema)

	// Testing boolean param with default value
	p = po.Post.Parameters[2]
	assert.EqualT(t, "someBoolean", p.Name)
	assert.EqualT(t, "some boolean", p.Description)
	assert.EqualT(t, "path", p.In)
	assert.FalseT(t, p.Required)
	assert.EqualT(t, "boolean", p.Type)
	someBoolean, ok := p.Default.(bool)
	assert.TrueT(t, ok)
	assert.TrueT(t, someBoolean)
	assert.Nil(t, p.Schema)

	// Testing that "min", "max", "minLength" and "maxLength" constraints will only be considered if the right type is provided
	p = po.Post.Parameters[3]
	assert.EqualT(t, "constraintsOnInvalidType", p.Name)
	assert.EqualT(t, "test constraints on invalid types", p.Description)
	assert.EqualT(t, "query", p.In)
	assert.EqualT(t, "boolean", p.Type)
	assert.Nil(t, p.Maximum)
	assert.Nil(t, p.Minimum)
	assert.Nil(t, p.MaxLength)
	assert.Nil(t, p.MinLength)
	assert.EqualT(t, "abcde", p.Format)
	constraintsOnInvalidType, ok2 := p.Default.(bool)
	assert.TrueT(t, ok2)
	assert.FalseT(t, constraintsOnInvalidType)
	assert.Nil(t, p.Schema)

	// Testing that when "type" is not provided, a schema will not be created
	p = po.Post.Parameters[4]
	assert.EqualT(t, "noType", p.Name)
	assert.EqualT(t, "test no type", p.Description)
	assert.Empty(t, p.Type)
	assert.Nil(t, p.Schema)

	// Testing a request body that takes a string value defined by a list of possible values in "enum"
	p = po.Post.Parameters[5]
	assert.EqualT(t, "request", p.Name)
	assert.EqualT(t, "The request model.", p.Description)
	assert.EqualT(t, "body", p.In)
	assert.EqualT(t, "string", p.Schema.Type[0])
	assert.Equal(t, "orange", p.Schema.Default)
	assert.Equal(t, []any{"apple", "orange", "pineapple", "peach", "plum"}, p.Schema.Enum)
	assert.Empty(t, p.Type)
}

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string, extensions spec.Extensions) {
	t.Helper()

	assert.NotNil(t, op)
	assert.EqualT(t, summary, op.Summary)
	assert.EqualT(t, description, op.Description)
	assert.EqualT(t, id, op.ID)
	assert.Equal(t, tags, op.Tags)
	assert.Equal(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.Equal(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.Equal(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.TrueT(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.TrueT(t, ok)
	assert.Equal(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.EqualT(t, "#/responses/genericError", op.Responses.Default.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.TrueT(t, ok)
	assert.EqualT(t, "#/responses/someResponse", rsp.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.TrueT(t, ok)
	assert.EqualT(t, "#/responses/validationError", rsp.Ref.String())

	ext := op.Extensions
	assert.Equal(t, extensions, ext)
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.EqualT(t, summary, op.Summary)
	assert.EqualT(t, description, op.Description)
	assert.EqualT(t, id, op.ID)
	assert.Equal(t, tags, op.Tags)
	assert.Equal(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.Equal(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.Equal(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	assert.TrueT(t, ok)
	// akv must be defined & not empty
	assert.NotNil(t, akv)
	assert.Empty(t, akv)

	vv, ok := op.Security[1]["oauth"]
	assert.TrueT(t, ok)
	assert.Equal(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Empty(t, op.Responses.Default.Ref.String())
	assert.EqualT(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.TrueT(t, ok)
	assert.Empty(t, rsp.Ref.String())
	assert.EqualT(t, "#/definitions/someResponse", rsp.Schema.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.TrueT(t, ok)
	assert.Empty(t, rsp.Ref.String())
	assert.EqualT(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}
