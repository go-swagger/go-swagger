package main

import (
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/examples/todo-list/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/todo-list/restapi/operations/todos"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.SimpleToDoListAPIAPI) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.KeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (key) x-petstore-token from header has not yet been implemented")
	}

	api.AddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation addOne has not yet been implemented")
	})
	api.DestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation destroyOne has not yet been implemented")
	})
	api.FindHandler = todos.FindHandlerFunc(func(principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation find has not yet been implemented")
	})
	api.UpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation updateOne has not yet been implemented")
	})

}
