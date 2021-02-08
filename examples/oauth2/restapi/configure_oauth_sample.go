// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/go-swagger/go-swagger/examples/oauth2/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/oauth2/restapi/operations/customers"

	models "github.com/go-swagger/go-swagger/examples/oauth2/models"
)

//go:generate swagger generate server --target .. --name oauthSample --spec ../swagger.yml --principal models.Principal

func configureFlags(api *operations.OauthSampleAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.OauthSampleAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.OauthSecurityAuth = func(token string, scopes []string) (*models.Principal, error) {
		ok, err := authenticated(token)
		if err != nil {
			return nil, errors.New(401, "error authenticate")
		}
		if !ok {
			return nil, errors.New(401, "invalid token")
		}
		prin := models.Principal(token)
		return &prin, nil

	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.GetAuthCallbackHandler = operations.GetAuthCallbackHandlerFunc(func(params operations.GetAuthCallbackParams) middleware.Responder {
		token, err := callback(params.HTTPRequest)
		if err != nil {
			return middleware.NotImplemented("operation .GetAuthCallback error")
		}
		log.Println("Token", token)
		return operations.NewGetAuthCallbackDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(token)})
	})
	api.GetLoginHandler = operations.GetLoginHandlerFunc(func(params operations.GetLoginParams) middleware.Responder {
		return login(params.HTTPRequest)
	})
	api.CustomersCreateHandler = customers.CreateHandlerFunc(func(params customers.CreateParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation customers.Create has not yet been implemented")
	})
	api.CustomersGetIDHandler = customers.GetIDHandlerFunc(func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
		log.Println("hit customer API")
		return middleware.NotImplemented("operation customers.GetID has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// This demonstrates how to enrich and pass custom context keys.
// In this case, we cache the current responseWriter in context.
type customContextKey int8

const (
	_ customContextKey = iota
	ctxResponseWriter
)

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	ourFunc := func(w http.ResponseWriter, r *http.Request) {
		rctx := context.WithValue(r.Context(), ctxResponseWriter, w)
		handler.ServeHTTP(w, r.WithContext(rctx))
	}
	return http.HandlerFunc(ourFunc)

}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
