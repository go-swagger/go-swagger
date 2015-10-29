package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
)

// SpecValidator validates a swagger spec
type SpecValidator struct {
	schema       *spec.Schema // swagger 2.0 schema
	spec         *spec.Document
	expanded     *spec.Document
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
	errs.Merge(schv.Validate(sd.Spec())) // error -
	if errs.HasErrors() {
		return // no point in continuing
	}

	errs.Merge(s.validateReferencesValid()) // error -
	if errs.HasErrors() {
		return // no point in continuing
	}

	errs.Merge(s.validateDuplicatePropertyNames())         // error -
	errs.Merge(s.validateParameters())                     // error -
	errs.Merge(s.validateItems())                          // error -
	errs.Merge(s.validateRequiredDefinitions())            // error -
	errs.Merge(s.validateDefaultValueValidAgainstSchema()) // error

	warnings.Merge(s.validateUniqueSecurityScopes())            // warning
	warnings.Merge(s.validateUniqueScopesSecurityDefinitions()) // warning
	warnings.Merge(s.validateReferenced())                      // warning

	return
}

type dupProp struct {
	Name       string
	Definition string
}

func (s *SpecValidator) validateDuplicatePropertyNames() *Result {
	// definition can't declare a property that's already defined by one of its ancestors
	res := new(Result)
	for k, sch := range s.spec.Spec().Definitions {
		if len(sch.AllOf) == 0 {
			continue
		}

		knownanc := map[string]struct{}{
			"#/definitions/" + k: struct{}{},
		}

		ancs := s.validateCircularAncestry(k, sch, knownanc)
		if len(ancs) > 0 {
			res.AddErrors(errors.New(422, "definition %q has circular ancestry: %v", k, ancs))
			return res
		}

		knowns := make(map[string]struct{})
		dups := s.validateSchemaPropertyNames(k, sch, knowns)
		if len(dups) > 0 {
			var pns []string
			for _, v := range dups {
				pns = append(pns, v.Definition+"."+v.Name)
			}
			res.AddErrors(errors.New(422, "definition %q contains duplicate properties: %v", k, pns))
		}

	}
	return res
}

func (s *SpecValidator) validateSchemaPropertyNames(nm string, sch spec.Schema, knowns map[string]struct{}) []dupProp {
	var dups []dupProp

	schn := nm
	schc := &sch
	if sch.Ref.GetURL() != nil {
		// gather property names
		reso, err := spec.ResolveRef(s.spec.Spec(), &sch.Ref)
		if err != nil {
			panic(err)
		}
		schc = reso
		schn = sch.Ref.String()
	}

	if len(schc.AllOf) > 0 {
		for _, chld := range schc.AllOf {
			dups = append(dups, s.validateSchemaPropertyNames(schn, chld, knowns)...)
		}
		return dups
	}

	for k := range schc.Properties {
		_, ok := knowns[k]
		if ok {
			dups = append(dups, dupProp{Name: k, Definition: schn})
		} else {
			knowns[k] = struct{}{}
		}
	}

	return dups
}

func (s *SpecValidator) validateCircularAncestry(nm string, sch spec.Schema, knowns map[string]struct{}) []string {
	var ancs []string

	schn := nm
	schc := &sch
	if sch.Ref.GetURL() != nil {
		reso, err := spec.ResolveRef(s.spec.Spec(), &sch.Ref)
		if err != nil {
			panic(err)
		}
		schc = reso
		schn = sch.Ref.String()
		knowns[schn] = struct{}{}
	}

	if _, ok := knowns[schn]; ok {
		ancs = append(ancs, schn)
	}
	if len(ancs) > 0 {
		return ancs
	}

	if len(schc.AllOf) > 0 {
		for _, chld := range schc.AllOf {
			ancs = append(ancs, s.validateCircularAncestry(schn, chld, knowns)...)
			if len(ancs) > 0 {
				return ancs
			}
		}
	}

	return ancs
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

func (s *SpecValidator) validatePathParamPresence(path string, fromPath, fromOperation []string) *Result {
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
			res.AddErrors(errors.New(422, "path param %q is not present in path %q", p, path))
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
			sw := s.spec.Spec()
			for _, ppr := range s.spec.ParamsFor(method, path) {
				pr := ppr
				if ppr.Ref.String() != "" {
					obj, _, err := ppr.Ref.GetPointer().Get(sw)
					if err != nil {
						res.AddErrors(err)
						break
					}
					pr = obj.(spec.Parameter)
				}

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
			res.Merge(s.validatePathParamPresence(path, fromPath, paramNames))
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
	res := new(Result)
	exp, err := s.spec.Expanded()
	if err != nil {
		res.AddErrors(err)
	}
	s.expanded = exp
	return res
}

func (s *SpecValidator) validateDefaultValueValidAgainstSchema() *Result {
	// every default value that is specified must validate against the schema for that property
	// headers, items, parameters, schema

	res := new(Result)

	for method, pathItem := range s.spec.Operations() {
		for path, op := range pathItem {
			// parameters
			for _, pr := range s.spec.ParamsFor(method, path) {
				// expand ref is necessary
				param := pr
				if pr.Ref.String() != "" {
					obj, _, err := pr.Ref.GetPointer().Get(s.spec.Spec())
					if err != nil {
						res.AddErrors(err)
						break
					}
					param = obj.(spec.Parameter)
				}
				// check simple paramters first
				if param.Default != nil && param.Schema == nil {
					//fmt.Println(param.Name, "in", param.In, "has a default without a schema")
					// check param valid
					res.Merge(NewParamValidator(&param, s.KnownFormats).Validate(param.Default))
				}

				if param.Items != nil {
					res.Merge(s.validateDefaultValueItemsAgainstSchema(param.Name, param.In, &param, param.Items))
				}

				if param.Schema != nil {
					res.Merge(s.validateDefaultValueSchemaAgainstSchema(param.Name, param.In, param.Schema))
				}

				//if param.Default != nil && param.Schema != nil {
				//res.Merge(NewSchemaValidator(param.Schema, nil, "", s.KnownFormats).Validate(param.Default))
				//}
			}

			if op.Responses.Default != nil {
				dr := op.Responses.Default
				for nm, h := range dr.Headers {
					if h.Default != nil {
						res.Merge(NewHeaderValidator(nm, &h, s.KnownFormats).Validate(h.Default))
					}
					if h.Items != nil {
						res.Merge(s.validateDefaultValueItemsAgainstSchema(nm, "header", &h, h.Items))
					}
				}
			}
			for _, r := range op.Responses.StatusCodeResponses {
				for nm, h := range r.Headers {
					if h.Default != nil {
						res.Merge(NewHeaderValidator(nm, &h, s.KnownFormats).Validate(h.Default))
					}
					if h.Items != nil {
						res.Merge(s.validateDefaultValueItemsAgainstSchema(nm, "header", &h, h.Items))
					}
				}
			}

		}
	}

	for nm, sch := range s.spec.Spec().Definitions {
		res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("definitions.%s", nm), "body", &sch))
	}

	return res
}

func (s *SpecValidator) validateDefaultValueSchemaAgainstSchema(path, in string, schema *spec.Schema) *Result {
	res := new(Result)
	if schema != nil {
		if schema.Default != nil {
			res.Merge(NewSchemaValidator(schema, s.spec.Spec(), path, s.KnownFormats).Validate(schema.Default))
		}
		if schema.Items != nil {
			if schema.Items.Schema != nil {
				res.Merge(s.validateDefaultValueSchemaAgainstSchema(path+".items", in, schema.Items.Schema))
			}
			for i, sch := range schema.Items.Schemas {
				res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.items[%d]", path, i), in, &sch))
			}
		}
		if schema.AdditionalItems != nil && schema.AdditionalItems.Schema != nil {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.additionalItems", path), in, schema.AdditionalItems.Schema))
		}
		for propName, prop := range schema.Properties {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(path+"."+propName, in, &prop))
		}
		for propName, prop := range schema.PatternProperties {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(path+"."+propName, in, &prop))
		}
		if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.additionalProperties", path), in, schema.AdditionalProperties.Schema))
		}
		for i, aoSch := range schema.AllOf {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.allOf[%d]", path, i), in, &aoSch))
		}

	}
	return res
}

func (s *SpecValidator) validateDefaultValueItemsAgainstSchema(path, in string, root interface{}, items *spec.Items) *Result {
	res := new(Result)
	if items != nil {
		if items.Default != nil {
			res.Merge(newItemsValidator(path, in, items, root, s.KnownFormats).Validate(0, items.Default))
		}
		if items.Items != nil {
			res.Merge(s.validateDefaultValueItemsAgainstSchema(path+"[0]", in, root, items.Items))
		}
	}
	return res
}

func (s *SpecValidator) isSwaggerType(tpe, format string, value interface{}) {
}
