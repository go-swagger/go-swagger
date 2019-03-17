// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/go-swagger/go-swagger/examples/task-tracker-strict-buildapi/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/task-tracker-strict-buildapi/restapi/operations/tasks"
)

//go:generate swagger generate server --target ../../task-tracker-strict-buildapi --name IssueTracker --spec ../swagger.yml --strict

func configureFlags(api *operations.IssueTrackerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.IssueTrackerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "token" query is set
	api.APIKeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (api_key) token from query param [token] has not yet been implemented")
	}
	// Applies when the "X-Token" header is set
	api.TokenHeaderAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (token_header) X-Token from header param [X-Token] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.TasksAddCommentToTaskHandler = tasks.AddCommentToTaskHandlerFunc(func(params tasks.AddCommentToTaskParams, principal interface{}) tasks.AddCommentToTaskResponder {
		return tasks.AddCommentToTaskNotImplemented()
	})
	api.TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal interface{}) tasks.CreateTaskResponder {
		return tasks.CreateTaskNotImplemented()
	})
	api.TasksDeleteTaskHandler = tasks.DeleteTaskHandlerFunc(func(params tasks.DeleteTaskParams, principal interface{}) tasks.DeleteTaskResponder {
		return tasks.DeleteTaskNotImplemented()
	})
	api.TasksGetTaskCommentsHandler = tasks.GetTaskCommentsHandlerFunc(func(params tasks.GetTaskCommentsParams) tasks.GetTaskCommentsResponder {
		return tasks.GetTaskCommentsNotImplemented()
	})
	api.TasksGetTaskDetailsHandler = tasks.GetTaskDetailsHandlerFunc(func(params tasks.GetTaskDetailsParams) tasks.GetTaskDetailsResponder {
		return tasks.GetTaskDetailsNotImplemented()
	})
	api.TasksListTasksHandler = tasks.ListTasksHandlerFunc(func(params tasks.ListTasksParams) tasks.ListTasksResponder {
		return tasks.ListTasksNotImplemented()
	})
	api.TasksUpdateTaskHandler = tasks.UpdateTaskHandlerFunc(func(params tasks.UpdateTaskParams, principal interface{}) tasks.UpdateTaskResponder {
		return tasks.UpdateTaskNotImplemented()
	})
	api.TasksUploadTaskFileHandler = tasks.UploadTaskFileHandlerFunc(func(params tasks.UploadTaskFileParams, principal interface{}) tasks.UploadTaskFileResponder {
		return tasks.UploadTaskFileNotImplemented()
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
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
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
