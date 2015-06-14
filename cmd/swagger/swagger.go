package main

import (
	"log"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands"
	"github.com/jessevdk/go-flags"
)

var opts struct{}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.ShortDescription = "helps you keep your API well described"
	parser.LongDescription = `
Swagger tries to support you as best as possible when building API's.

It aims to represent the contract of your API with a language agnostic description of your application in json or yaml.
`
	parser.AddCommand("validate", "validate the swagger document", "validate the provided swagger document against a swagger spec", &commands.ValidateSpec{})

	genpar, err := parser.AddCommand("generate", "genererate go code", "generate go code for the swagger spec file", &commands.Generate{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, cmd := range genpar.Commands() {
		switch cmd.Name {
		case "spec":
			cmd.ShortDescription = "generate a swagger spec document from a go application"
			cmd.LongDescription = cmd.ShortDescription
		case "client":
			cmd.ShortDescription = "generate all the files for a client library"
			cmd.LongDescription = cmd.ShortDescription
		case "server":
			cmd.ShortDescription = "generate all the files for a server application"
			cmd.LongDescription = cmd.ShortDescription
		case "model":
			cmd.ShortDescription = "generate one or more models from the swagger spec"
			cmd.LongDescription = cmd.ShortDescription
		case "support":
			cmd.ShortDescription = "generate supporting files like the main function and the api builder"
			cmd.LongDescription = cmd.ShortDescription
		case "operation":
			cmd.ShortDescription = "generate one or more server operations from the swagger spec"
			cmd.LongDescription = cmd.ShortDescription
		}
	}

	parser.Parse()
}
