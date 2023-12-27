//go:build ignore

package main

// this test is designed to be run dynamically by go-swagger test suite
// (generator.generate_test.go), after a test client has been generated.

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	rtclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/middleware/untyped"
	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1083/codegen/client"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1083/codegen/client/pet_operations"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1083/codegen/models"
	"github.com/stretchr/testify/require"
)

type operationHandler struct{}

func (h *operationHandler) Handle(in interface{}) (interface{}, error) {
	// this handler sends back the input path parameter in Name

	paramsMap, ok := in.(map[string]interface{})
	if !ok {
		return nil, errors.New(http.StatusInternalServerError, "unexpected params: wants a map")
	}
	val, ok := paramsMap["id"]
	if !ok {
		return nil, errors.New(http.StatusBadRequest, "unexpected params: {id} required")
	}
	id, ok := val.(string)
	if !ok {
		return nil, errors.New(http.StatusBadRequest, "unexpected params: {id} should be a string")
	}

	return models.Pet{
		ID:   swag.Int64(1),
		Name: swag.String(id),
		Tag:  "test",
	}, nil
}

func TestEscapedPathParam(t *testing.T) {
	// prepare an untyped test server
	buildTestServer := func(t testing.TB) (string, func()) {
		t.Helper()
		spec, err := loads.Spec("petstore.yaml")
		require.NoError(t, err)

		api := untyped.NewAPI(spec)
		api.RegisterConsumer("application/json", runtime.JSONConsumer())
		api.RegisterProducer("application/json", runtime.JSONProducer())
		api.RegisterOperation("GET", "/pets/{id}", &operationHandler{})
		handler := middleware.Serve(spec, api)

		server := httptest.NewServer(handler)
		u, err := url.Parse(server.URL)
		require.NoError(t, err)

		return u.Host, server.Close
	}

	buildTestClient := func(t testing.TB, host string) *client.Issue1083 {
		t.Helper()

		c := client.Default
		tr := rtclient.New(host, "/api", []string{"http"})
		tr.Debug = true
		c.SetTransport(tr)

		return c
	}

	submitAndAssert := func(c *client.Issue1083, pathParam string) func(*testing.T) {
		return func(t *testing.T) {
			params := pet_operations.NewGetPetsIDParams().WithID(pathParam)
			resp, err := c.PetOperations.GetPetsID(params)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int64(1), swag.Int64Value(resp.Payload.ID))
			require.Equal(t, pathParam, swag.StringValue(resp.Payload.Name))
		}
	}

	host, clean := buildTestServer(t)
	t.Cleanup(clean)

	c := buildTestClient(t, host)

	t.Run("should route with unescaped path param", submitAndAssert(c, "part"))
	t.Run("should route with escaped path param (1)", submitAndAssert(c, "part/ext"))
	t.Run("should route with escaped path param (2)", submitAndAssert(c, "part#ext"))
}
