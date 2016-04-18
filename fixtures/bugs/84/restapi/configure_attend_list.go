package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	httpkit "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations"
	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations/events"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.AttendListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AttendListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.EventsDeleteEventByIDHandler = events.DeleteEventByIDHandlerFunc(func(params events.DeleteEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation events.DeleteEventByID has not yet been implemented")
	})
	api.EventsGetEventByIDHandler = events.GetEventByIDHandlerFunc(func(params events.GetEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation events.GetEventByID has not yet been implemented")
	})
	api.EventsGetEventsHandler = events.GetEventsHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation events.GetEvents has not yet been implemented")
	})
	api.EventsPostEventHandler = events.PostEventHandlerFunc(func(params events.PostEventParams) middleware.Responder {
		return middleware.NotImplemented("operation events.PostEvent has not yet been implemented")
	})
	api.EventsPutEventByIDHandler = events.PutEventByIDHandlerFunc(func(params events.PutEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation events.PutEventByID has not yet been implemented")
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
