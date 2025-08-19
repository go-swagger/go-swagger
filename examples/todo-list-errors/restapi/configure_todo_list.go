// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"

	oapierrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/examples/todo-list-errors/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/todo-list-errors/restapi/operations/todos"
)

//go:generate swagger generate server --target ../../todo-list-errors --name TodoList --spec ../swagger.yml --principal any --return-errors

func configureFlags(api *operations.TodoListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func catcher(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, errAlreadyExists) {
		slog.Info("we catch custom error! congratulations!")
	}

	oapierrors.ServeError(w, r, err)
}

//nolint:gochecknoglobals
var errAlreadyExists = errors.New("already exists")

func configureAPI(api *operations.TodoListAPI) http.Handler {
	// configure the api here
	api.ServeError = catcher

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...any)
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// todos.FindMaxParseMemory = 32 << 20

	api.TodosAddOneHandler = todos.AddOneHandlerFunc(
		func(params todos.AddOneParams) (middleware.Responder, error) {
			return nil, errAlreadyExists
		})

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
