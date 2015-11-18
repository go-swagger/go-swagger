package main

import (
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations"
	"github.com/go-swagger/go-swagger/fixtures/bugs/84/restapi/operations/events"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.EventListAPI) {
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

}
