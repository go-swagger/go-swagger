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

package validate

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// defaultValidator validates default values in a spec.
// According to Swagger spec, default values MUST validate their schema.
type defaultValidator struct {
	SpecValidator *SpecValidator
}

// Validate validates the default values declared in the swagger spec
func (d *defaultValidator) Validate() (errs *Result) {
	errs = new(Result)
	if d == nil || d.SpecValidator == nil {
		return errs
	}
	errs.Merge(d.validateDefaultValueValidAgainstSchema()) // error -
	return errs
}

func (d *defaultValidator) validateDefaultValueValidAgainstSchema() *Result {
	// every default value that is specified must validate against the schema for that property
	// headers, items, parameters, schema

	res := new(Result)
	s := d.SpecValidator

	for method, pathItem := range s.analyzer.Operations() {
		if pathItem != nil { // Safeguard
			for path, op := range pathItem {
				// parameters
				for _, param := range paramHelp.safeExpandedParamsFor(path, method, op.ID, res, s) {
					if param.Default != nil && param.Required {
						res.AddWarnings(requiredHasDefaultMsg(param.Name, param.In))
					}

					// Check simple parameters first
					// default values provided must validate against their inline definition (no explicit schema)
					if param.Default != nil && param.Schema == nil {
						// check param default value is valid
						red := NewParamValidator(&param, s.KnownFormats).Validate(param.Default)
						if red.HasErrorsOrWarnings() {
							res.AddErrors(defaultValueDoesNotValidateMsg(param.Name, param.In))
							res.Merge(red)
						}
					}

					// Recursively follows Items and Schemas
					if param.Items != nil {
						red := d.validateDefaultValueItemsAgainstSchema(param.Name, param.In, &param, param.Items)
						if red.HasErrorsOrWarnings() {
							res.AddErrors(defaultValueItemsDoesNotValidateMsg(param.Name, param.In))
							res.Merge(red)
						}
					}

					if param.Schema != nil {
						// Validate default value against schema
						red := d.validateDefaultValueSchemaAgainstSchema(param.Name, param.In, param.Schema)
						if red.HasErrorsOrWarnings() {
							res.AddErrors(defaultValueDoesNotValidateMsg(param.Name, param.In))
							res.Merge(red)
						}
					}
				}

				if op.Responses != nil {
					if op.Responses.Default != nil {
						// Same constraint on default Response
						res.Merge(d.validateDefaultInResponse(op.Responses.Default, "default", path, 0, op.ID))
					}
					// Same constraint on regular Responses
					if op.Responses.StatusCodeResponses != nil { // Safeguard
						for code, r := range op.Responses.StatusCodeResponses {
							res.Merge(d.validateDefaultInResponse(&r, "response", path, code, op.ID))
						}
					}
				} else {
					// Empty op.ID means there is no meaningful operation: no need to report a specific message
					if op.ID != "" {
						res.AddErrors(noValidResponseMsg(op.ID))
					}
				}
			}
		}
	}
	if s.spec.Spec().Definitions != nil { // Safeguard
		for nm, sch := range s.spec.Spec().Definitions {
			res.Merge(d.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("definitions.%s", nm), "body", &sch))
		}
	}
	return res
}

func (d *defaultValidator) validateDefaultInResponse(resp *spec.Response, responseType, path string, responseCode int, operationID string) *Result {
	s := d.SpecValidator

	response, res := responseHelp.expandResponseRef(resp, path, s)
	if !res.IsValid() {
		return res
	}

	responseName, responseCodeAsStr := responseHelp.responseMsgVariants(responseType, responseCode)

	if response.Headers != nil { // Safeguard
		for nm, h := range response.Headers {
			if h.Default != nil {
				red := NewHeaderValidator(nm, &h, s.KnownFormats).Validate(h.Default)
				if red.HasErrorsOrWarnings() {
					res.AddErrors(defaultValueHeaderDoesNotValidateMsg(operationID, nm, responseName))
					res.Merge(red)
				}
			}

			// Headers have inline definition, like params
			if h.Items != nil {
				red := d.validateDefaultValueItemsAgainstSchema(nm, "header", &h, h.Items)
				if red.HasErrorsOrWarnings() {
					res.AddErrors(defaultValueHeaderItemsDoesNotValidateMsg(operationID, nm, responseName))
					res.Merge(red)
				}
			}

			if _, err := compileRegexp(h.Pattern); err != nil {
				res.AddErrors(invalidPatternInHeaderMsg(operationID, nm, responseName, h.Pattern, err))
			}

			// Headers don't have schema
		}
	}
	if response.Schema != nil {
		red := d.validateDefaultValueSchemaAgainstSchema(responseCodeAsStr, "response", response.Schema)
		if red.HasErrorsOrWarnings() {
			// Additional message to make sure the context of the error is not lost
			res.AddErrors(defaultValueInDoesNotValidateMsg(operationID, responseName))
			res.Merge(red)
		}
	}
	return res
}

func (d *defaultValidator) validateDefaultValueSchemaAgainstSchema(path, in string, schema *spec.Schema) *Result {
	res := new(Result)
	s := d.SpecValidator
	if schema != nil { // Safeguard
		if schema.Default != nil {
			res.Merge(NewSchemaValidator(schema, s.spec.Spec(), path+".default", s.KnownFormats).Validate(schema.Default))
		}
		if schema.Items != nil {
			if schema.Items.Schema != nil {
				res.Merge(d.validateDefaultValueSchemaAgainstSchema(path+".items.default", in, schema.Items.Schema))
			}
			// Multiple schemas in items
			if schema.Items.Schemas != nil { // Safeguard
				for i, sch := range schema.Items.Schemas {
					res.Merge(d.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.items[%d].default", path, i), in, &sch))
				}
			}
		}
		if _, err := compileRegexp(schema.Pattern); err != nil {
			res.AddErrors(invalidPatternInMsg(path, in, schema.Pattern))
		}
		if schema.AdditionalItems != nil && schema.AdditionalItems.Schema != nil {
			// NOTE: we keep validating values, even though additionalItems is not supported by Swagger 2.0 (and 3.0 as well)
			res.Merge(d.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.additionalItems", path), in, schema.AdditionalItems.Schema))
		}
		for propName, prop := range schema.Properties {
			res.Merge(d.validateDefaultValueSchemaAgainstSchema(path+"."+propName, in, &prop))
		}
		for propName, prop := range schema.PatternProperties {
			res.Merge(d.validateDefaultValueSchemaAgainstSchema(path+"."+propName, in, &prop))
		}
		if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
			res.Merge(d.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.additionalProperties", path), in, schema.AdditionalProperties.Schema))
		}
		if schema.AllOf != nil {
			for i, aoSch := range schema.AllOf {
				res.Merge(d.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.allOf[%d]", path, i), in, &aoSch))
			}
		}
	}
	return res
}

func (d *defaultValidator) validateDefaultValueItemsAgainstSchema(path, in string, root interface{}, items *spec.Items) *Result {
	res := new(Result)
	s := d.SpecValidator
	if items != nil {
		if items.Default != nil {
			res.Merge(newItemsValidator(path, in, items, root, s.KnownFormats).Validate(0, items.Default))
		}
		if items.Items != nil {
			res.Merge(d.validateDefaultValueItemsAgainstSchema(path+"[0].default", in, root, items.Items))
		}
		if _, err := compileRegexp(items.Pattern); err != nil {
			res.AddErrors(invalidPatternInMsg(path, in, items.Pattern))
		}
	}
	return res
}
