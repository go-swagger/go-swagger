package main

import (
	"log"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
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

	log.Println("Would be serving:", specDoc.Spec().Info.Title)
}
