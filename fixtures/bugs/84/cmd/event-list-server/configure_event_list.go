package main

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations"
	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations/events"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.EventListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.DeleteEventByIDHandler = events.DeleteEventByIDHandlerFunc(func(params events.DeleteEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation deleteEventById has not yet been implemented")
	})
	api.GetEventByIDHandler = events.GetEventByIDHandlerFunc(func(params events.GetEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation getEventById has not yet been implemented")
	})
	api.GetEventsHandler = events.GetEventsHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation getEvents has not yet been implemented")
	})
	api.PostEventHandler = events.PostEventHandlerFunc(func(params events.PostEventParams) middleware.Responder {
		return middleware.NotImplemented("operation postEvent has not yet been implemented")
	})
	api.PutEventByIDHandler = events.PutEventByIDHandlerFunc(func(params events.PutEventByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation putEventById has not yet been implemented")
	})

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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
