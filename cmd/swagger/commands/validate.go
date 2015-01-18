package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
	"github.com/casualjim/go-swagger/validate"
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
	b, err := util.JSONDoc(swaggerDoc)
	if err != nil {
		return err
	}
	var toValidate interface{}
	if err := json.Unmarshal(b, &toValidate); err != nil {
		return err
	}
	result := validate.WithSchema(spec.MustLoadSwagger20Schema(), toValidate)

	if result.IsValid() {
		fmt.Printf("The swagger spec at %q is valid against swagger specification %s\n", swaggerDoc, "2.0")
	} else {
		str := fmt.Sprintf("The swagger spec at %q is invalid against swagger specification %s. see errors :\n", swaggerDoc, "2.0")
		for _, desc := range result.Errors {
			str += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(str)
	}
	return nil
}
