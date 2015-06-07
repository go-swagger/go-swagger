package main

import (
	"log"

	"github.com/casualjim/go-swagger/fixtures/goparsing/petstore"
	"github.com/casualjim/go-swagger/fixtures/goparsing/petstore/rest"
)

var (
	// Version is a compile time constant, injected at build time
	Version string
)

// This is an application that doesn't actually do anything,
// it's used for testing the scanner
func main() {
	// this has no real purpose besides making the import present in this main package.
	// without this line the meta info for the swagger doc wouldn't be discovered
	petstore.APIVersion = Version

	// This servers na hypothetical API
	if err := rest.ServeAPI(); err != nil {
		log.Fatal(err)
	}
}
