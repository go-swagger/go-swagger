
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/examples/todo-list/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/todo-list/restapi/operations/todos"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.ToDoListAPI) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.XPetstoreTokenAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth x-petstore-token from header has not yet been implemented")
	}

	api.AddOneHandler = todos.AddOneHandlerFunc(func(principal interface{}) middleware.Responder {
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
