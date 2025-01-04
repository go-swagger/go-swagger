package codescan

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestRouteExpression(t *testing.T) {
	t.Parallel()

	patternDelete := `swagger:route\s+DELETE\s+/orders/\{id\}\s+deleteOrder`
	patternGet := `swagger:route\s+GET\s+/v1\.2/something\s+deleteOrder`

	require.Regexp(t, patternDelete, "swagger:route DELETE /orders/{id} deleteOrder", "delete route pattern should match")
	require.Regexp(t, patternGet, "swagger:route GET /v1.2/something deleteOrder", "get route pattern should match")
}

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string, extensions spec.Extensions) {
	require.NotNil(t, op)
	require.Equal(t, summary, op.Summary)
	require.Equal(t, description, op.Description)
	require.Equal(t, id, op.ID)

	expectedConsumes := []string{"application/json", "application/x-protobuf"}
	expectedProduces := []string{"application/json", "application/x-protobuf"}
	expectedSchemes := []string{"http", "https", "ws", "wss"}

	require.ElementsMatch(t, tags, op.Tags)
	require.ElementsMatch(t, expectedConsumes, op.Consumes)
	require.ElementsMatch(t, expectedProduces, op.Produces)
	require.ElementsMatch(t, expectedSchemes, op.Schemes)

	require.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	require.True(t, ok)
	require.Empty(t, akv, "api_key security scope should be defined but empty")

	vv, ok := op.Security[1]["oauth"]
	require.True(t, ok)
	require.ElementsMatch(t, scopes, vv)

	require.NotNil(t, op.Responses.Default)
	require.Equal(t, "#/responses/genericError", op.Responses.Default.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	require.True(t, ok)
	require.Equal(t, "#/responses/someResponse", rsp.Ref.String())

	rsp, ok = op.Responses.StatusCodeResponses[422]
	require.True(t, ok)
	require.Equal(t, "#/responses/validationError", rsp.Ref.String())

	if extensions != nil {
		expectedJSON, err := json.Marshal(extensions)
		require.NoError(t, err)
		actualJSON, err := json.Marshal(op.Extensions)
		require.NoError(t, err)
		require.JSONEq(t, string(expectedJSON), string(actualJSON))
	} else {
		require.Empty(t, op.Extensions)
	}
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	require.NotNil(t, op)
	require.Equal(t, summary, op.Summary)
	require.Equal(t, description, op.Description)
	require.Equal(t, id, op.ID)

	expectedConsumes := []string{"application/json", "application/x-protobuf"}
	expectedProduces := []string{"application/json", "application/x-protobuf"}
	expectedSchemes := []string{"http", "https", "ws", "wss"}

	require.ElementsMatch(t, tags, op.Tags)
	require.ElementsMatch(t, expectedConsumes, op.Consumes)
	require.ElementsMatch(t, expectedProduces, op.Produces)
	require.ElementsMatch(t, expectedSchemes, op.Schemes)

	require.Len(t, op.Security, 2)
	akv, ok := op.Security[0]["api_key"]
	require.True(t, ok)
	require.Empty(t, akv, "api_key security scope should be defined but empty")

	vv, ok := op.Security[1]["oauth"]
	require.True(t, ok)
	require.ElementsMatch(t, scopes, vv)

	require.NotNil(t, op.Responses.Default)
	require.Empty(t, op.Responses.Default.Ref.String())
	require.Equal(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	require.True(t, ok)
	require.Empty(t, rsp.Ref.String())
	require.Equal(t, "#/definitions/someResponse", rsp.Schema.Ref.String())

	rsp, ok = op.Responses.StatusCodeResponses[422]
	require.True(t, ok)
	require.Empty(t, rsp.Ref.String())
	require.Equal(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}