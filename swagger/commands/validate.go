package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateSpec is a command that validates a swagger document
// against the swagger json schema
type ValidateSpec struct {
	SchemaURL string `long:"schema" description:"The schema url to use" default:"http://swagger.io/v2/schema.json"`
}

// Execute validates the spec
func (c *ValidateSpec) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("The validate command requires the swagger document url to be specified")
	}

	swaggerDoc := args[0]
	schemaDocument, err := gojsonschema.NewJsonSchemaDocument(c.SchemaURL)
	if err != nil {
		return err
	}

	var jsonDocument interface{}
	if strings.HasPrefix(swaggerDoc, "http") {
		// Loads the JSON to validate from a http location
		jsonDocument, err = gojsonschema.GetHttpJson(swaggerDoc)
	} else {
		// Loads the JSON to validate from a local file
		jsonDocument, err = gojsonschema.GetFileJson(swaggerDoc)
	}
	if err != nil {
		return err
	}

	// Try to validate the Json against the schema
	result := schemaDocument.Validate(jsonDocument)

	// Deal with result
	if result.Valid() {
		fmt.Printf("The swagger spec at %q is valid against schema %q\n", swaggerDoc, c.SchemaURL)
	} else {
		fmt.Printf("The swagger spec at %q is valid against schema %q. see errors :\n", swaggerDoc, c.SchemaURL)
		// Loop through errors
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	return nil
}
