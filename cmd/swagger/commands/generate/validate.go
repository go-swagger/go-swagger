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

package generate

import (
	"errors"
	"fmt"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	swaggererrors "github.com/go-openapi/errors"
)

// validateSpec is a shared function that validates a swagger document
// against the swagger json schema
func validateSpec(swaggerDoc string) error {
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		return err
	}

	result := validate.Spec(specDoc, strfmt.Default)
	if result == nil {
		return nil
	}

	str := fmt.Sprintf("The swagger spec at %q is invalid against swagger specification %s. see errors :\n", swaggerDoc, specDoc.Version())
	for _, desc := range result.(*swaggererrors.CompositeError).Errors {
		str += fmt.Sprintf("- %s\n", desc)
	}
	return errors.New(str)
}
