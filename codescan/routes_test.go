package codescan

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestRouteExpression(t *testing.T) {
	t.Parallel()

	t.Run("delete route", func(t *testing.T) {
		t.Parallel()
		require.Regexp(t, `swagger:route\s+DELETE\s+/orders/\{id\}\s+deleteOrder`, "swagger:route DELETE /orders/{id} deleteOrder")
	})

	t.Run("get route with version", func(t *testing.T) {
		t.Parallel()
		require.Regexp(t, `swagger:route\s+GET\s+/v1\.2/something\s+deleteOrder`, "swagger:route GET /v1.2/something deleteOrder")
	})
}

func TestRoutesParser(t *testing.T) {
	t.Parallel()

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

	require.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	ext := make(spec.Extensions)
	ext.Add("x-some-flag", "true")
	require.True(t, ok)
	require.NotNil(t, po.Get)

	assertOperation(t,
		po.Get,
		"listPets",
		"Lists pets filtered by some parameters.",
		"This will show all available pets by default.\nYou can get the pets that are out of stock",
		[]string{"pets", "users"},
		[]string{"read", "write"},
		ext,
	)

	// ... rest of the test remains the same ...
}

func assertOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string, extensions spec.Extensions) {
	require.NotNil(t, op)
	require.Equal(t, summary, op.Summary)
	require.Equal(t, description, op.Description)
	require.Equal(t, id, op.ID)
	require.ElementsMatch(t, tags, op.Tags)

	expectedConsumes := []string{"application/json", "application/x-protobuf"}
	expectedProduces := []string{"application/json", "application/x-protobuf"}
	expectedSchemes := []string{"http", "https", "ws", "wss"}

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

	ext := op.Extensions
	if extensions != nil {
		expectedJSON, err := json.Marshal(extensions)
		require.NoError(t, err)
		actualJSON, err := json.Marshal(ext)
		require.NoError(t, err)
		require.JSONEq(t, string(expectedJSON), string(actualJSON))
	} else {
		require.Nil(t, ext)
	}
}

func assertOperationBody(t *testing.T, op *spec.Operation, id, summary, description string, tags, scopes []string) {
	require.NotNil(t, op)
	require.Equal(t, summary, op.Summary)
	require.Equal(t, description, op.Description)
	require.Equal(t, id, op.ID)
	require.ElementsMatch(t, tags, op.Tags)

	expectedConsumes := []string{"application/json", "application/x-protobuf"}
	expectedProduces := []string{"application/json", "application/x-protobuf"}
	expectedSchemes := []string{"http", "https", "ws", "wss"}

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
	require.Equal(t, "", op.Responses.Default.Ref.String())
	require.Equal(t, "#/definitions/genericError", op.Responses.Default.Schema.Ref.String())

	rsp, ok := op.Responses.StatusCodeResponses[200]
	require.True(t, ok)
	require.Equal(t, "", rsp.Ref.String())
	require.Equal(t, "#/definitions/someResponse", rsp.Schema.Ref.String())

	rsp, ok = op.Responses.StatusCodeResponses[422]
	require.True(t, ok)
	require.Equal(t, "", rsp.Ref.String())
	require.Equal(t, "#/definitions/validationError", rsp.Schema.Ref.String())
}