package main

import (
	"github.com/casualjim/go-swagger/cmd/swagger/commands"
	"github.com/jessevdk/go-flags"
)

var opts struct{}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.ShortDescription = "helps you keep your API well described"
	parser.LongDescription = `
Swagger tries to support you as best as possible when building API's
It aims to represent the contract of your API with a language agnostic description of your application in json or yaml.
`

	parser.AddCommand("validate", "validate the swagger document", "validate the provided swagger document against a swagger spec", &commands.ValidateSpec{})

	// parser.AddCommand("editor", "edit the swagger.json document", "serve the swagger editor with the specified spec file", commands.NewEditor())
	parser.AddCommand("ui", "api-docs for the swagger.json document", "serve the swagger ui application with the specified spec file", commands.NewUI())

	parser.AddCommand("generate", "genererate go code", "generate go code for the swagger spec file", &commands.Generate{})
	parser.Parse()
}
