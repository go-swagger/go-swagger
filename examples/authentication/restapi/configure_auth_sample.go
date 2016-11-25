package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/examples/authentication/models"
	"github.com/go-swagger/go-swagger/examples/authentication/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/authentication/restapi/operations/customers"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name AuthSample --spec ../swagger.yml --principal models.Principal

func configureFlags(api *operations.AuthSampleAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AuthSampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	api.KeyAuth = func(token string) (*models.Principal, error) {
		if token == "abcdefuvwxyz" {
			prin := models.Principal(token)
			return &prin, nil
		}
		api.Logger("Access attempt with incorrect api key auth: %s", token)
		return nil, errors.New(401, "incorrect api key auth")
	}

	api.CustomersCreateHandler = customers.CreateHandlerFunc(func(params customers.CreateParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation customers.Create has not yet been implemented")
	})
	api.CustomersGetIDHandler = customers.GetIDHandlerFunc(func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation customers.GetID has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
