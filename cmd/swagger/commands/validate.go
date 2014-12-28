package commands

import (
	"errors"
	"fmt"

	"github.com/casualjim/go-swagger/swagger/load"
)

// ValidateSpec is a command that validates a swagger document
// against the swagger json schema
type ValidateSpec struct {
	// SchemaURL string `long:"schema" description:"The schema url to use" default:"http://swagger.io/v2/schema.json"`
}

// Execute validates the spec
func (c *ValidateSpec) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("The validate command requires the swagger document url to be specified")
	}

	swaggerDoc := args[0]
	schemaDocument, err := load.JSONSpec(swaggerDoc)
	if err != nil {
		return err
	}

	result := schemaDocument.Validate()
	if result.Valid() {
		fmt.Printf("The swagger spec at %q is valid against swagger specification %s\n", swaggerDoc, schemaDocument.Version())
	} else {
		str := fmt.Sprintf("The swagger spec at %q is valid against swagger specification %s. see errors :\n", swaggerDoc, schemaDocument.Version())
		for _, desc := range result.Errors() {
			str += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(str)
	}
	return nil
}
