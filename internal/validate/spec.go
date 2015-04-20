package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

// SpecValidator validates a swagger spec
type SpecValidator struct {
	schema       *spec.Schema // swagger 2.0 schema
	spec         *spec.Document
	KnownFormats strfmt.Registry
}

// NewSpecValidator creates a new swagger spec validator instance
func NewSpecValidator(schema *spec.Schema, formats strfmt.Registry) *SpecValidator {
	return &SpecValidator{
		schema:       schema,
		KnownFormats: formats,
	}
}

// Validate validates the swagger spec
func (s *SpecValidator) Validate(data interface{}) (errs *Result, warnings *Result) {
	var sd *spec.Document

	switch v := data.(type) {
	case *spec.Document:
		sd = v
	}
	if sd == nil {
		errs = sErr(errors.New(500, "spec validator can only validate spec.Document objects"))
		return
	}
	s.spec = sd

	errs = new(Result)
	warnings = new(Result)

	schv := NewSchemaValidator(s.schema, nil, "", s.KnownFormats)
	errs.Merge(schv.Validate(sd.Spec()))                        // error -
	errs.Merge(s.validateItems())                               // error -
	warnings.Merge(s.validateUniqueSecurityScopes())            // warning
	warnings.Merge(s.validateUniqueScopesSecurityDefinitions()) // warning
	warnings.Merge(s.validateReferenced())                      // warning
	errs.Merge(s.validateRequiredDefinitions())                 // error
	errs.Merge(s.validateParameters())                          // error -
	errs.Merge(s.validateReferencesValid())                     // error
	errs.Merge(s.validateDefaultValueValidAgainstSchema())      // error

	return
}

func (s *SpecValidator) validateItems() *Result {
	// validate parameter, items, schema and response objects for presence of item if type is array
	res := new(Result)

	// TODO: implement support for lookups of refs
	for method, pi := range s.spec.Operations() {
		for path, op := range pi {
			for _, param := range s.spec.ParamsFor(method, path) {
				if param.TypeName() == "array" && param.ItemsTypeName() == "" {
					res.AddErrors(errors.New(422, "param %q for %q is a collection without an element type", param.Name, op.ID))
					continue
				}
				if param.In != "body" {
					if param.Items != nil {
						items := param.Items
						for items.TypeName() == "array" {
							if items.ItemsTypeName() == "" {
								res.AddErrors(errors.New(422, "param %q for %q is a collection without an element type", param.Name, op.ID))
								break
							}
							items = items.Items
						}
					}
				} else {
					if err := s.validateSchemaItems(*param.Schema, fmt.Sprintf("body param %q", param.Name), op.ID); err != nil {
						res.AddErrors(err)
					}
				}
			}

			var responses []spec.Response
			if op.Responses != nil {
				if op.Responses.Default != nil {
					responses = append(responses, *op.Responses.Default)
				}
				for _, v := range op.Responses.StatusCodeResponses {
					responses = append(responses, v)
				}
			}

			for _, resp := range responses {
				for hn, hv := range resp.Headers {
					if hv.TypeName() == "array" && hv.ItemsTypeName() == "" {
						res.AddErrors(errors.New(422, "header %q for %q is a collection without an element type", hn, op.ID))
					}
				}
				if resp.Schema != nil {
					if err := s.validateSchemaItems(*resp.Schema, "response body", op.ID); err != nil {
						res.AddErrors(err)
					}
				}
			}
		}
	}
	return res
}

func (s *SpecValidator) validateSchemaItems(schema spec.Schema, prefix, opID string) error {
	if !schema.Type.Contains("array") {
		return nil
	}

	if schema.Items == nil || schema.Items.Len() == 0 {
		return errors.New(422, "%s for %q is a collection without an element type", prefix, opID)
	}

	schemas := schema.Items.Schemas
	if schema.Items.Schema != nil {
		schemas = []spec.Schema{*schema.Items.Schema}
	}
	for _, sch := range schemas {
		if err := s.validateSchemaItems(sch, prefix, opID); err != nil {
			return err
		}
	}
	return nil
}

func (s *SpecValidator) validateUniqueSecurityScopes() *Result {
	// Each authorization/security reference should contain only unique scopes.
	// (Example: For an oauth2 authorization/security requirement, when listing the required scopes,
	// each scope should only be listed once.)
	return nil
}

func (s *SpecValidator) validateUniqueScopesSecurityDefinitions() *Result {
	// Each authorization/security scope in an authorization/security definition should be unique.
	return nil
}

func (s *SpecValidator) validatePathParamPresence(fromPath, fromOperation []string) *Result {
	// Each defined operation path parameters must correspond to a named element in the API's path pattern.
	// (For example, you cannot have a path parameter named id for the following path /pets/{petId} but you must have a path parameter named petId.)
	res := new(Result)
	for _, l := range fromPath {
		var matched bool
		for _, r := range fromOperation {
			if l == "{"+r+"}" {
				matched = true
				break
			}
		}
		if !matched {
			res.Errors = append(res.Errors, errors.New(422, "path param %q has no parameter definition", l))
		}
	}

	for _, p := range fromOperation {
		var matched bool
		for _, r := range fromPath {
			if "{"+p+"}" == r {
				matched = true
				break
			}
		}
		if !matched {
			res.AddErrors(errors.New(422, "path param %q is not present in the path", p))
		}
	}

	return res
}

func (s *SpecValidator) validateReferenced() *Result {
	// Each referenceable definition must have references.
	return nil
}

func (s *SpecValidator) validateRequiredDefinitions() *Result {
	// Each definition property listed in the required array must be defined in the properties of the model
	res := new(Result)
	for d, v := range s.spec.Spec().Definitions {
	REQUIRED:
		for _, pn := range v.Required {
			if _, ok := v.Properties[pn]; ok {
				continue
			}

			for pp := range v.PatternProperties {
				re := regexp.MustCompile(pp)
				if re.MatchString(pn) {
					continue REQUIRED
				}
			}

			if v.AdditionalProperties != nil {
				if v.AdditionalProperties.Allows {
					continue
				}
				if v.AdditionalProperties.Schema != nil {
					continue
				}
			}

			res.AddErrors(errors.New(422, "%q is present in required but not defined as property in defintion %q", pn, d))
		}
	}
	return res
}

func (s *SpecValidator) validateParameters() *Result {
	// each parameter should have a unique `name` and `type` combination
	// each operation should have only 1 parameter of type body
	// each api path should be non-verbatim (account for path param names) unique per method
	res := new(Result)
	for method, pi := range s.spec.Operations() {
		knownPaths := make(map[string]string)
		for path, op := range pi {
			segments, params := parsePath(path)
			knowns := make([]string, 0, len(segments))
			for _, s := range segments {
				knowns = append(knowns, s)
			}
			var fromPath []string
			for _, i := range params {
				fromPath = append(fromPath, knowns[i])
				knowns[i] = "!"
			}
			knownPath := strings.Join(knowns, "/")
			if orig, ok := knownPaths[knownPath]; ok {
				res.AddErrors(errors.New(422, "path %s overlaps with %s", path, orig))
			} else {
				knownPaths[knownPath] = path
			}

			ptypes := make(map[string]map[string]struct{})
			var firstBodyParam string

			var paramNames []string
			for _, pr := range op.Parameters {
				pnames, ok := ptypes[pr.In]
				if !ok {
					pnames = make(map[string]struct{})
					ptypes[pr.In] = pnames
				}

				_, ok = pnames[pr.Name]
				if ok {
					res.AddErrors(errors.New(422, "duplicate parameter name %q for %q in operation %q", pr.Name, pr.In, op.ID))
				}
				pnames[pr.Name] = struct{}{}
			}
			for _, pr := range s.spec.ParamsFor(method, path) {
				if pr.In == "body" {
					if firstBodyParam != "" {
						res.AddErrors(errors.New(422, "operation %q has more than 1 body param (accepted: %q, dropped: %q)", op.ID, firstBodyParam, pr.Name))
					}
					firstBodyParam = pr.Name
				}

				if pr.In == "path" {
					paramNames = append(paramNames, pr.Name)
				}
			}
			res.Merge(s.validatePathParamPresence(fromPath, paramNames))
		}
	}
	return res
}

func parsePath(path string) (segments []string, params []int) {
	for i, p := range strings.Split(path, "/") {
		segments = append(segments, p)
		if len(p) > 0 && p[0] == '{' && p[len(p)-1] == '}' {
			params = append(params, i)
		}
	}
	return
}

func (s *SpecValidator) validateReferencesValid() *Result {
	// each reference must point to a valid object
	return nil
}

func (s *SpecValidator) validateDefaultValueValidAgainstSchema() *Result {
	// every default value that is specified must validate against the schema for that property
	return nil
}
