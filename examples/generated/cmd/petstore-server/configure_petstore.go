package main

import (
	"io"
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/pet"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/store"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/user"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.PetstoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.XMLConsumer = httpkit.ConsumerFunc(func(r io.Reader, target interface{}) error {
		return errors.NotImplemented("xml consumer has not yet been implemented")
	})

	api.JSONProducer = httpkit.JSONProducer()

	api.XMLProducer = httpkit.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("xml producer has not yet been implemented")
	})

	api.APIKeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (api_key) api_key from header has not yet been implemented")
	}

	api.AddPetHandler = pet.AddPetHandlerFunc(func(params pet.AddPetParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation addPet has not yet been implemented")
	})
	api.CreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		return middleware.NotImplemented("operation createUser has not yet been implemented")
	})
	api.CreateUsersWithArrayInputHandler = user.CreateUsersWithArrayInputHandlerFunc(func(params user.CreateUsersWithArrayInputParams) middleware.Responder {
		return middleware.NotImplemented("operation createUsersWithArrayInput has not yet been implemented")
	})
	api.CreateUsersWithListInputHandler = user.CreateUsersWithListInputHandlerFunc(func(params user.CreateUsersWithListInputParams) middleware.Responder {
		return middleware.NotImplemented("operation createUsersWithListInput has not yet been implemented")
	})
	api.DeleteOrderHandler = store.DeleteOrderHandlerFunc(func(params store.DeleteOrderParams) middleware.Responder {
		return middleware.NotImplemented("operation deleteOrder has not yet been implemented")
	})
	api.DeletePetHandler = pet.DeletePetHandlerFunc(func(params pet.DeletePetParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation deletePet has not yet been implemented")
	})
	api.DeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
		return middleware.NotImplemented("operation deleteUser has not yet been implemented")
	})
	api.FindPetsByStatusHandler = pet.FindPetsByStatusHandlerFunc(func(params pet.FindPetsByStatusParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation findPetsByStatus has not yet been implemented")
	})
	api.FindPetsByTagsHandler = pet.FindPetsByTagsHandlerFunc(func(params pet.FindPetsByTagsParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation findPetsByTags has not yet been implemented")
	})
	api.GetOrderByIDHandler = store.GetOrderByIDHandlerFunc(func(params store.GetOrderByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation getOrderById has not yet been implemented")
	})
	api.GetPetByIDHandler = pet.GetPetByIDHandlerFunc(func(params pet.GetPetByIDParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation getPetById has not yet been implemented")
	})
	api.GetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation getUserByName has not yet been implemented")
	})
	api.LoginUserHandler = user.LoginUserHandlerFunc(func(params user.LoginUserParams) middleware.Responder {
		return middleware.NotImplemented("operation loginUser has not yet been implemented")
	})
	api.LogoutUserHandler = user.LogoutUserHandlerFunc(func() middleware.Responder {
		return middleware.NotImplemented("operation logoutUser has not yet been implemented")
	})
	api.PlaceOrderHandler = store.PlaceOrderHandlerFunc(func(params store.PlaceOrderParams) middleware.Responder {
		return middleware.NotImplemented("operation placeOrder has not yet been implemented")
	})
	api.UpdatePetHandler = pet.UpdatePetHandlerFunc(func(params pet.UpdatePetParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation updatePet has not yet been implemented")
	})
	api.UpdatePetWithFormHandler = pet.UpdatePetWithFormHandlerFunc(func(params pet.UpdatePetWithFormParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation updatePetWithForm has not yet been implemented")
	})
	api.UpdateUserHandler = user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) middleware.Responder {
		return middleware.NotImplemented("operation updateUser has not yet been implemented")
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
