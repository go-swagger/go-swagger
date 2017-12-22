// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

const (
	missingArgMsg  = "The validate command requires the swagger document url to be specified"
	validSpecMsg   = "\nThe swagger spec at %q is valid against swagger specification %s\n"
	invalidSpecMsg = "\nThe swagger spec at %q is invalid against swagger specification %s. See errors below:\n"
	warningSpecMsg = "\nThe swagger spec at %q showed up some valid but possiby unwanted constructs. See warnings below:\n"
)

// ValidateSpec is a command that validates a swagger document
// against the swagger json schema
type ValidateSpec struct {
	// SchemaURL string `long:"schema" description:"The schema url to use" default:"http://swagger.io/v2/schema.json"`
}

// Execute validates the spec
func (c *ValidateSpec) Execute(args []string) error {
	// TODO: make optional
	showWarnings := true

	if len(args) == 0 {
		return errors.New(missingArgMsg)
	}

	swaggerDoc := args[0]

	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}

	// Attempts to report about all errors
	// TODO: as arg
	validate.SetContinueOnErrors(true)

	v := validate.NewSpecValidator(specDoc.Schema(), strfmt.Default)
	result, _ := v.Validate(specDoc) // returns fully detailed result with errors and warnings
	//result := validate.Spec(specDoc, strfmt.Default)		// returns single error

	if result.IsValid() {
		log.Printf(validSpecMsg, swaggerDoc, specDoc.Version())
	}
	if result.HasWarnings() {
		log.Printf(warningSpecMsg, swaggerDoc)
		if showWarnings {
			for _, desc := range result.Warnings {
				log.Printf("- %s\n", desc.Error())
			}
		}
	}
	if result.HasErrors() {
		str := fmt.Sprintf(invalidSpecMsg, swaggerDoc, specDoc.Version())
		for _, desc := range result.Errors {
			str += fmt.Sprintf("- %s\n", desc.Error())
		}
		return errors.New(str)
	}

	return nil
}
