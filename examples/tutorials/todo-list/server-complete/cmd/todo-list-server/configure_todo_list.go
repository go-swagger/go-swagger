package main

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete/models"
	"github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete/restapi/operations/todos"
)

// This file is safe to edit. Once it exists it will not be overwritten

var store = make(map[int64]models.Item)
var ids int64

func configureAPI(api *operations.TodoListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
		ids += 1
		item := *params.Body
		item.ID = ids
		store[item.ID] = item
		return todos.NewAddOneCreated().WithPayload(&item)
	})

	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
		delete(store, params.ID)
		return todos.NewDestroyOneNoContent()
	})

	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(func(params todos.FindTodosParams) middleware.Responder {
		items := make([]*models.Item, 0, params.Limit)
		for id := range store {
			if id > params.Since {
				item := store[id]
				items = append(items, &item)
				if len(items) == int(params.Limit) {
					break
				}
			}
		}

		return todos.NewFindTodosOK().WithPayload(items)
	})

	api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
		params.Body.ID = params.ID
		item := *params.Body
		store[params.ID] = item
		return todos.NewUpdateOneOK().WithPayload(&item)
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
