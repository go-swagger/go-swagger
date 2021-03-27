// +build ignore

package bug_1472

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/fixtures/bugs/1472/restapi"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1472/restapi/operations"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1472/restapi/operations/ops"
)

type User struct {
	ID string
}

func TestGenServer_1472_securityFromPrincipal_middleware(t *testing.T) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewMyAPI(swaggerSpec)
	api.OpBearerAuth = func(token string) (interface{}, error) {
		// Return the same user for testing
		return &User{ID: "someID"}, nil
	}
	api.OpsGetEndpointHandler = ops.GetEndpointHandlerFunc(
		func(params ops.GetEndpointParams, principal interface{}) middleware.Responder {
			return ops.NewGetEndpointOK()
		})

	apiHandler := api.Serve(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			assert.Nil(t, middleware.SecurityPrincipalFrom(r)) // Did not run through the handler yet!
			h.ServeHTTP(rw, r)
			assert.NotNil(t, middleware.SecurityPrincipalFrom(r)) // Should have the principal now!
			assert.Equal(t, middleware.SecurityPrincipalFrom(r).(*User).ID, "someID")
		})
	})

	ts := httptest.NewServer(apiHandler)
	defer ts.Close()

	request, err := http.NewRequest("GET", ts.URL+"/endpoint", nil)
	assert.Nil(t, err)
	request.Header.Add("Authorization", "Bearer sometoken")

	response, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
}
