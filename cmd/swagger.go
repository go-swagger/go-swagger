package main

import (
	"github.com/casualjim/go-swagger/cmd/commands"
	"github.com/jessevdk/go-flags"
)

var opts struct{}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.AddCommand("validate", "validate the swagger document", "validate the provided swagger document against a swagger spec", &commands.ValidateSpec{})
	parser.Parse()
}
