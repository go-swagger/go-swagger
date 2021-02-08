// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-swagger/go-swagger/examples/stream-server/biz"

	"github.com/go-swagger/go-swagger/examples/stream-server/restapi/operations"
)

//go:generate swagger generate server --target .. --name Countdown --spec ../swagger.yml

func configureFlags(api *operations.CountdownAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CountdownAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	myCounter := &biz.MyCounter{}
	api.ElapseHandler = operations.ElapseHandlerFunc(func(params operations.ElapseParams) middleware.Responder {
		if params.Length == 11 {
			return operations.NewElapseForbidden()
		}
		return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
			f, _ := rw.(http.Flusher)
			rw.WriteHeader(200)
			_ = myCounter.Down(params.Length, &flushWriter{f: f, w: rw})
		})

	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// Via https://play.golang.org/p/PpbPyXbtEs
type flushWriter struct {
	f http.Flusher
	w io.Writer
}

// Via https://play.golang.org/p/PpbPyXbtEs
func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
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
