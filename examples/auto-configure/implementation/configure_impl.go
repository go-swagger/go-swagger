package implementation

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/examples/auto-configure/restapi/operations"
)

type ConfigureImpl struct {
	flags Flags
}

type Flags struct {
	Example1 string `long:"example1" description:"Sample for showing how to configure cmd-line flags"`
	Example2 string `long:"example2" description:"Further info at https://github.com/jessevdk/go-flags"`
}

func (i *ConfigureImpl) ConfigureFlags(api *operations.AToDoListApplicationAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Example Flags",
			LongDescription:  "",
			Options:          &i.flags,
		},
	}
}

func (i *ConfigureImpl) ConfigureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

func (i *ConfigureImpl) ConfigureServer(s *http.Server, scheme, addr string) {
	if i.flags.Example1 != "something" {
		log.Println("example1 argument is not something")
	}
}

func (i *ConfigureImpl) SetupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func (i *ConfigureImpl) SetupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Recieved request on path: %v", req.URL.String())
		handler.ServeHTTP(w, req)
	})
}

func (i *ConfigureImpl) CustomConfigure(api *operations.AToDoListApplicationAPI) {
	api.Logger = log.Printf
	api.ServerShutdown = func() {
		log.Printf("Running ServerShutdown function")
	}
}
