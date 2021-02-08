// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/pet"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/store"
	"github.com/go-swagger/go-swagger/examples/generated/restapi/operations/user"
)

//go:generate swagger generate server --target ../../generated --name Petstore --spec ../swagger-petstore.json --principal interface{}

func configureFlags(api *operations.PetstoreAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PetstoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.UrlformConsumer = runtime.DiscardConsumer
	api.XMLConsumer = runtime.XMLConsumer()

	api.JSONProducer = runtime.JSONProducer()
	api.XMLProducer = runtime.XMLProducer()

	// Applies when the "api_key" header is set
	if api.APIKeyAuth == nil {
		api.APIKeyAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (api_key) api_key from header param [api_key] has not yet been implemented")
		}
	}
	if api.PetstoreAuthAuth == nil {
		api.PetstoreAuthAuth = func(token string, scopes []string) (interface{}, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (petstore_auth) has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// pet.UpdatePetWithFormMaxParseMemory = 32 << 20

	if api.PetAddPetHandler == nil {
		api.PetAddPetHandler = pet.AddPetHandlerFunc(func(params pet.AddPetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.AddPet has not yet been implemented")
		})
	}
	if api.UserCreateUserHandler == nil {
		api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.CreateUser has not yet been implemented")
		})
	}
	if api.UserCreateUsersWithArrayInputHandler == nil {
		api.UserCreateUsersWithArrayInputHandler = user.CreateUsersWithArrayInputHandlerFunc(func(params user.CreateUsersWithArrayInputParams) middleware.Responder {
			return middleware.NotImplemented("operation user.CreateUsersWithArrayInput has not yet been implemented")
		})
	}
	if api.UserCreateUsersWithListInputHandler == nil {
		api.UserCreateUsersWithListInputHandler = user.CreateUsersWithListInputHandlerFunc(func(params user.CreateUsersWithListInputParams) middleware.Responder {
			return middleware.NotImplemented("operation user.CreateUsersWithListInput has not yet been implemented")
		})
	}
	if api.StoreDeleteOrderHandler == nil {
		api.StoreDeleteOrderHandler = store.DeleteOrderHandlerFunc(func(params store.DeleteOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation store.DeleteOrder has not yet been implemented")
		})
	}
	if api.PetDeletePetHandler == nil {
		api.PetDeletePetHandler = pet.DeletePetHandlerFunc(func(params pet.DeletePetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.DeletePet has not yet been implemented")
		})
	}
	if api.UserDeleteUserHandler == nil {
		api.UserDeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.DeleteUser has not yet been implemented")
		})
	}
	if api.PetFindPetsByStatusHandler == nil {
		api.PetFindPetsByStatusHandler = pet.FindPetsByStatusHandlerFunc(func(params pet.FindPetsByStatusParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.FindPetsByStatus has not yet been implemented")
		})
	}
	if api.PetFindPetsByTagsHandler == nil {
		api.PetFindPetsByTagsHandler = pet.FindPetsByTagsHandlerFunc(func(params pet.FindPetsByTagsParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.FindPetsByTags has not yet been implemented")
		})
	}
	if api.StoreGetOrderByIDHandler == nil {
		api.StoreGetOrderByIDHandler = store.GetOrderByIDHandlerFunc(func(params store.GetOrderByIDParams) middleware.Responder {
			return middleware.NotImplemented("operation store.GetOrderByID has not yet been implemented")
		})
	}
	if api.PetGetPetByIDHandler == nil {
		api.PetGetPetByIDHandler = pet.GetPetByIDHandlerFunc(func(params pet.GetPetByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.GetPetByID has not yet been implemented")
		})
	}
	if api.UserGetUserByNameHandler == nil {
		api.UserGetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
			return middleware.NotImplemented("operation user.GetUserByName has not yet been implemented")
		})
	}
	if api.UserLoginUserHandler == nil {
		api.UserLoginUserHandler = user.LoginUserHandlerFunc(func(params user.LoginUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.LoginUser has not yet been implemented")
		})
	}
	if api.UserLogoutUserHandler == nil {
		api.UserLogoutUserHandler = user.LogoutUserHandlerFunc(func(params user.LogoutUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.LogoutUser has not yet been implemented")
		})
	}
	if api.StorePlaceOrderHandler == nil {
		api.StorePlaceOrderHandler = store.PlaceOrderHandlerFunc(func(params store.PlaceOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation store.PlaceOrder has not yet been implemented")
		})
	}
	if api.PetUpdatePetHandler == nil {
		api.PetUpdatePetHandler = pet.UpdatePetHandlerFunc(func(params pet.UpdatePetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.UpdatePet has not yet been implemented")
		})
	}
	if api.PetUpdatePetWithFormHandler == nil {
		api.PetUpdatePetWithFormHandler = pet.UpdatePetWithFormHandlerFunc(func(params pet.UpdatePetWithFormParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation pet.UpdatePetWithForm has not yet been implemented")
		})
	}
	if api.UserUpdateUserHandler == nil {
		api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.UpdateUser has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

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

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
