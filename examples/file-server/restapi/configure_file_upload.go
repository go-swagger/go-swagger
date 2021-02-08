// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/go-swagger/examples/file-server/restapi/operations"
	"github.com/go-swagger/go-swagger/examples/file-server/restapi/operations/uploads"
)

//go:generate swagger generate server --target ../../file-server --name FileUpload --spec ../swagger.yml --principal interface{}

func configureFlags(api *operations.FileUploadAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.FileUploadAPI) http.Handler {
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
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// uploads.UploadFileMaxParseMemory = 32 << 20

	uploadFolder, err := ioutil.TempDir(".", "upload")
	if err != nil {
		panic("could not create upload folder")
	}
	uploadCounter := 0

	api.UploadsUploadFileHandler = uploads.UploadFileHandlerFunc(func(params uploads.UploadFileParams) middleware.Responder {

		if params.File == nil {
			return middleware.Error(404, fmt.Errorf("no file provided"))
		}
		defer func() {
			_ = params.File.Close()
		}()

		if namedFile, ok := params.File.(*runtime.File); ok {
			log.Printf("received file name: %s", namedFile.Header.Filename)
			log.Printf("received file size: %d", namedFile.Header.Size)
		}

		// uploads file and save it locally
		filename := path.Join(uploadFolder, fmt.Sprintf("uploaded_file_%d.dat", uploadCounter))
		uploadCounter++
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return middleware.Error(500, fmt.Errorf("could not create file on server"))
		}

		n, err := io.Copy(f, params.File)
		if err != nil {
			return middleware.Error(500, fmt.Errorf("could not upload file on server"))
		}

		log.Printf("copied bytes %d", n)

		log.Printf("file uploaded copied as %s", filename)

		return uploads.NewUploadFileOK()
	})

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
