package main

import (
	"log"

	"github.com/casualjim/go-swagger/fixtures/goparsing/petstore/rest"
)

func main() {
	// This is an application that doesn't actually do anything,
	// it's used for testing the scanner

	if err := rest.ServeAPI(); err != nil {
		log.Fatal(err)
	}
}
