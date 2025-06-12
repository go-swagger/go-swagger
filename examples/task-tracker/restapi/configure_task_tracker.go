// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/examples/task-tracker/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/task-tracker/restapi/operations/tasks"
)

//go:generate swagger generate server --target ../../task-tracker --name TaskTracker --spec ../swagger.yml --principal any

func configureFlags(api *operations.TaskTrackerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.TaskTrackerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...any)
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "token" query is set
	if api.APIKeyAuth == nil {
		api.APIKeyAuth = func(token string) (any, error) {
			_ = token

			return nil, errors.NotImplemented("api key auth (api_key) token from query param [token] has not yet been implemented")
		}
	}
	// Applies when the "X-Token" header is set
	if api.TokenHeaderAuth == nil {
		api.TokenHeaderAuth = func(token string) (any, error) {
			_ = token

			return nil, errors.NotImplemented("api key auth (token_header) X-Token from header param [X-Token] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// tasks.UploadTaskFileMaxParseMemory = 32 << 20

	if api.TasksAddCommentToTaskHandler == nil {
		api.TasksAddCommentToTaskHandler = tasks.AddCommentToTaskHandlerFunc(func(params tasks.AddCommentToTaskParams, principal any) middleware.Responder {
			_ = params
			_ = principal

			return middleware.NotImplemented("operation tasks.AddCommentToTask has not yet been implemented")
		})
	}
	if api.TasksCreateTaskHandler == nil {
		api.TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal any) middleware.Responder {
			_ = params
			_ = principal

			return middleware.NotImplemented("operation tasks.CreateTask has not yet been implemented")
		})
	}
	if api.TasksDeleteTaskHandler == nil {
		api.TasksDeleteTaskHandler = tasks.DeleteTaskHandlerFunc(func(params tasks.DeleteTaskParams, principal any) middleware.Responder {
			_ = params
			_ = principal

			return middleware.NotImplemented("operation tasks.DeleteTask has not yet been implemented")
		})
	}
	if api.TasksGetTaskCommentsHandler == nil {
		api.TasksGetTaskCommentsHandler = tasks.GetTaskCommentsHandlerFunc(func(params tasks.GetTaskCommentsParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation tasks.GetTaskComments has not yet been implemented")
		})
	}
	if api.TasksGetTaskDetailsHandler == nil {
		api.TasksGetTaskDetailsHandler = tasks.GetTaskDetailsHandlerFunc(func(params tasks.GetTaskDetailsParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation tasks.GetTaskDetails has not yet been implemented")
		})
	}
	if api.TasksListTasksHandler == nil {
		api.TasksListTasksHandler = tasks.ListTasksHandlerFunc(func(params tasks.ListTasksParams) middleware.Responder {
			_ = params

			return middleware.NotImplemented("operation tasks.ListTasks has not yet been implemented")
		})
	}
	if api.TasksUpdateTaskHandler == nil {
		api.TasksUpdateTaskHandler = tasks.UpdateTaskHandlerFunc(func(params tasks.UpdateTaskParams, principal any) middleware.Responder {
			_ = params
			_ = principal

			return middleware.NotImplemented("operation tasks.UpdateTask has not yet been implemented")
		})
	}
	if api.TasksUploadTaskFileHandler == nil {
		api.TasksUploadTaskFileHandler = tasks.UploadTaskFileHandlerFunc(func(params tasks.UploadTaskFileParams, principal any) middleware.Responder {
			_ = params
			_ = principal

			return middleware.NotImplemented("operation tasks.UploadTaskFile has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	_ = tlsConfig
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(server *http.Server, scheme, addr string) {
	_ = server
	_ = scheme
	_ = addr
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
