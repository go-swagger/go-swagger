// +build ignore

package main

import (
	"log"
	"testing"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1943/restapi"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1943/restapi/operations"
)

func Test_DataRace(t *testing.T) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPseudoServiceAPI(swaggerSpec)
	server := restapi.NewServer(api)

	server.ConfigureFlags()

	server.ConfigureAPI()

	go func() {
		time.Sleep(1 * time.Second)
		server.Shutdown()
	}()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
