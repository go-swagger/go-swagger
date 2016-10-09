package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/middleware/untyped"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// our spec doesn't have application/json in the consumes or produces
	// so we need to clear those settings out
	api := untyped.NewAPI(specDoc).WithoutJSONDefaults()

	// register serializers
	mediaType := "application/io.goswagger.examples.todo-list.v1+json"
	api.DefaultConsumes = mediaType
	api.DefaultProduces = mediaType
	api.RegisterConsumer(mediaType, runtime.JSONConsumer())
	api.RegisterProducer(mediaType, runtime.JSONProducer())

	// register the operation handlers
	api.RegisterOperation("GET", "/", notImplemented)
	api.RegisterOperation("POST", "/", notImplemented)
	api.RegisterOperation("PUT", "/{id}", notImplemented)
	api.RegisterOperation("DELETE", "/{id}", notImplemented)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}

	// construct the application context for this server
	// use the loaded spec document and the api descriptor with the default router
	app := middleware.NewContext(specDoc, api, nil)

	log.Println("serving", specDoc.Spec().Info.Title, "at http://localhost:8000")
	// serve the api
	if err := http.ListenAndServe(":8000", app.APIHandler(nil)); err != nil {
		log.Fatalln(err)
	}
}

var notImplemented = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	return middleware.NotImplemented("not implemented"), nil
})
