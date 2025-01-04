package codescan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestRouteExpression(t *testing.T) {
	assert.Regexp(t, rxRoute, "swagger:route DELETE /orders/{id} deleteOrder")
	assert.Regexp(t, rxRoute, "swagger:route GET /v1.2/something deleteOrder")
}

// ... [previous test code remains the same until the assertOperation function] ...

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string, extensions spec.Extensions) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	assert.Empty(t, akv, "api_key security scope should be defined but empty")

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "#/responses/genericError", op.Responses.Default.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/someResponse", rsp.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "#/responses/validationError", rsp.Ref.String())

	ext := op.Extensions
	assert.Equal(t, extensions, ext)
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Consumes)
	assert.EqualValues(t, []string{"application/json", "application/x-protobuf"}, op.Produces)
	assert.EqualValues(t, []string{"http", "https", "ws", "wss"}, op.Schemes)
	assert.Len(t, op.Security, 2)
	
	akv, ok := op.Security[0]["api_key"]
	assert.True(t, ok)
	assert.Empty(t, akv, "api_key security scope should be defined but empty")

	vv, ok := op.Security[1]["oauth"]
	assert.True(t, ok)
	assert.EqualValues(t, scopes, vv)

	assert.NotNil(t, op.Responses.Default)
	assert.Equal(t, "", op.Responses.Default.Ref.String())
	assert.Equal(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/someResponse", rsp.Schema.Ref.String())
	rsp, ok = op.Responses.StatusCodeResponses[422]
	assert.True(t, ok)
	assert.Equal(t, "", rsp.Ref.String())
	assert.Equal(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}