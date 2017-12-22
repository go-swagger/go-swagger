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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
)

var (
	continueOnErrors = false
)

// SetContinueOnErrors ...
// For extended error reporting, it's better to pass all validations.
// For faster validation, it's better to give up early.
// SetContinueOnError(true) will set the validator to continue to the end of its checks.
func SetContinueOnErrors(c bool) {
	continueOnErrors = c
}

// Spec validates an OpenAPI 2.0 specification document.
// It validates the spec json against the json schema for swagger
// and then validates a number of extra rules that can't be expressed in json schema.
//
// Returns an error flattening in a single standard error, all validation messages.
//
// Reported as errors:
// 	- definition can't declare a property that's already defined by one of its ancestors
// 	- definition's ancestor can't be a descendant of the same model
// 	- path uniqueness: each api path should be non-verbatim (account for path param names) unique per method
// 	- each security reference should contain only unique scopes
// 	- each security scope in a security definition should be unique
//  - parameters in path must be unique
// 	- each path parameter must correspond to a parameter placeholder and vice versa
// 	- each referencable definition must have references
// 	- each definition property listed in the required array must be defined in the properties of the model
// 	- each parameter should have a unique `name` and `type` combination
// 	- each operation should have only 1 parameter of type body
// 	- each reference must point to a valid object
// 	- every default value that is specified must validate against the schema for that property
// 	- items property is required for all schemas/definitions of type `array`
//  - path parameters must be declared a required
//  - headers must not contain $ref
//  - schema and property examples provided must validate against their respective object's schema
//
// Reported as warnings:
//  - path parameters should not contain any of [{,},\w]
//  - empty path
//  - TODO: warnings id or $id
//  - TODO: $ref should not have siblings
//
// NOTE:
// - SecurityScopes are maps: no need to check uniqueness
//
func Spec(doc *loads.Document, formats strfmt.Registry) error {
	errs, _ /*warns*/ := NewSpecValidator(doc.Schema(), formats).Validate(doc)
	if errs.HasErrors() {
		return errors.CompositeValidationError(errs.Errors...)
	}
	return nil
}

// AgainstSchema validates the specified data with the provided schema, when no schema
// is provided it uses the json schema as default
func AgainstSchema(schema *spec.Schema, data interface{}, formats strfmt.Registry) error {
	res := NewSchemaValidator(schema, nil, "", formats).Validate(data)
	if res.HasErrors() {
		return errors.CompositeValidationError(res.Errors...)
	}
	return nil
}

// SpecValidator validates a swagger spec
type SpecValidator struct {
	schema       *spec.Schema // swagger 2.0 schema
	spec         *loads.Document
	analyzer     *analysis.Spec
	expanded     *loads.Document
	KnownFormats strfmt.Registry
	brokenRefs   bool // Safety indicator to scan analyzer resources
	brokenParam  bool // Safety indicator to scan analyzer resources
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
	var sd *loads.Document

	switch v := data.(type) {
	case *loads.Document:
		sd = v
	}
	if sd == nil {
		// TODO: should use a constant (from errors package?)
		errs = sErr(errors.New(500, "spec validator can only validate spec.Document objects"))
		return
	}
	s.spec = sd
	s.analyzer = analysis.New(sd.Spec())

	errs = new(Result)
	warnings = new(Result)

	schv := NewSchemaValidator(s.schema, nil, "", s.KnownFormats)
	var obj interface{}

	// Raw spec unmarshalling errors
	if err := json.Unmarshal(sd.Raw(), &obj); err != nil {
		errs.AddErrors(err)
		return
	}

	errs.Merge(schv.Validate(obj)) // error -
	// There may be a point in continuing to try and determine more accurate errors
	if !continueOnErrors && errs.HasErrors() {
		return // no point in continuing
	}

	errs.Merge(s.validateReferencesValid()) // error -
	// There may be a point in continuing to try and determine more accurate errors
	if !continueOnErrors && errs.HasErrors() {
		return // no point in continuing
	}

	errs.Merge(s.validateDuplicateOperationIDs())
	errs.Merge(s.validateDuplicatePropertyNames())         // error -
	errs.Merge(s.validateParameters())                     // error -
	errs.Merge(s.validateItems())                          // error -
	errs.Merge(s.validateRequiredDefinitions())            // error -
	errs.Merge(s.validateDefaultValueValidAgainstSchema()) // error -
	errs.Merge(s.validateExamplesValidAgainstSchema())     // error -
	errs.Merge(s.validateNonEmptyPathParamNames())

	warnings.MergeAsErrors(s.validateRefNoSibling()) // warning only
	warnings.MergeAsErrors(s.validateReferenced())   // warning only

	// errs holds all errors and warnings,
	// warnings only all warnings
	errs.MergeAsWarnings(warnings)
	warnings.AddErrors(errs.Warnings...)
	return
}

func (s *SpecValidator) validateNonEmptyPathParamNames() *Result {
	res := new(Result)
	if s.spec.Spec().Paths == nil {
		// There is no Paths object: error
		res.AddErrors(errors.New(errors.CompositeErrorCode, "spec has no valid path defined"))
	} else {
		if s.spec.Spec().Paths.Paths == nil {
			// Paths may be empty: warning
			res.AddWarnings(errors.New(errors.CompositeErrorCode, "spec has no valid path defined"))
		} else {
			for k := range s.spec.Spec().Paths.Paths {
				if strings.Contains(k, "{}") {
					res.AddErrors(errors.New(errors.CompositeErrorCode, "%q contains an empty path parameter", k))
				}
			}

		}

	}
	return res
}

// TODO: there is a catch here. Duplicate operationId are not strictly forbidden, but
// not supported by go-swagger. Shouldn't it be a warning?
func (s *SpecValidator) validateDuplicateOperationIDs() *Result {
	res := new(Result)
	known := make(map[string]int)
	for _, v := range s.analyzer.OperationIDs() {
		if v != "" {
			known[v]++
		}
	}
	for k, v := range known {
		if v > 1 {
			res.AddErrors(errors.New(errors.CompositeErrorCode, "%q is defined %d times", k, v))
		}
	}
	return res
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

		ancs, rec := s.validateCircularAncestry(k, sch, knownanc)
		if rec != nil && (rec.HasErrors() || !rec.HasWarnings()) {
			res.Merge(rec)
		}
		if len(ancs) > 0 {
			res.AddErrors(errors.New(errors.CompositeErrorCode, "definition %q has circular ancestry: %v", k, ancs))
			return res
		}

		knowns := make(map[string]struct{})
		dups, rep := s.validateSchemaPropertyNames(k, sch, knowns)
		if rep != nil && (rep.HasErrors() || rep.HasWarnings()) {
			res.Merge(rep)
		}
		if len(dups) > 0 {
			var pns []string
			for _, v := range dups {
				pns = append(pns, v.Definition+"."+v.Name)
			}
			res.AddErrors(errors.New(errors.CompositeErrorCode, "definition %q contains duplicate properties: %v", k, pns))
		}

	}
	return res
}

func (s *SpecValidator) resolveRef(ref *spec.Ref) (*spec.Schema, error) {
	if s.spec.SpecFilePath() != "" {
		return spec.ResolveRefWithBase(s.spec.Spec(), ref, &spec.ExpandOptions{RelativeBase: s.spec.SpecFilePath()})
	}
	return spec.ResolveRef(s.spec.Spec(), ref)
}

func (s *SpecValidator) validateSchemaPropertyNames(nm string, sch spec.Schema, knowns map[string]struct{}) ([]dupProp, *Result) {
	var dups []dupProp

	schn := nm
	schc := &sch
	res := new(Result)

	for schc.Ref.String() != "" {
		// gather property names
		reso, err := s.resolveRef(&schc.Ref)
		if err != nil {
			s.addPointerError(res, err, schc.Ref.String(), nm)
			return dups, res
		}
		schc = reso
		schn = sch.Ref.String()
	}

	if len(schc.AllOf) > 0 {
		for _, chld := range schc.AllOf {
			dup, rep := s.validateSchemaPropertyNames(schn, chld, knowns)
			if rep != nil && (rep.HasErrors() || rep.HasWarnings()) {
				res.Merge(rep)
			}
			dups = append(dups, dup...)
		}
		return dups, res
	}

	for k := range schc.Properties {
		_, ok := knowns[k]
		if ok {
			dups = append(dups, dupProp{Name: k, Definition: schn})
		} else {
			knowns[k] = struct{}{}
		}
	}

	return dups, res
}

func (s *SpecValidator) validateCircularAncestry(nm string, sch spec.Schema, knowns map[string]struct{}) ([]string, *Result) {
	res := new(Result)

	if sch.Ref.String() == "" && len(sch.AllOf) == 0 {
		return nil, res
	}
	var ancs []string

	schn := nm
	schc := &sch

	for schc.Ref.String() != "" {
		reso, err := s.resolveRef(&schc.Ref)
		if err != nil {
			s.addPointerError(res, err, schc.Ref.String(), nm)
			return ancs, res
		}
		schc = reso
		schn = sch.Ref.String()
	}

	if schn != nm && schn != "" {
		if _, ok := knowns[schn]; ok {
			ancs = append(ancs, schn)
		}
		knowns[schn] = struct{}{}

		if len(ancs) > 0 {
			return ancs, res
		}
	}

	if len(schc.AllOf) > 0 {
		for _, chld := range schc.AllOf {
			if chld.Ref.String() != "" || len(chld.AllOf) > 0 {
				anc, rec := s.validateCircularAncestry(schn, chld, knowns)
				if rec != nil && (rec.HasErrors() || !rec.HasWarnings()) {
					res.Merge(rec)
				}
				ancs = append(ancs, anc...)
				if len(ancs) > 0 {
					return ancs, res
				}
			}
		}
	}
	return ancs, res
}

func (s *SpecValidator) validateItems() *Result {
	// validate parameter, items, schema and response objects for presence of item if type is array
	res := new(Result)

	// TODO: implement support for lookups of refs
	for method, pi := range s.analyzer.Operations() {
		for path, op := range pi {
			if !s.brokenRefs && !s.brokenParam { // Safeguard (since ParamsFor() panics on brokenRefs||brokenParam)
				for _, param := range s.analyzer.ParamsFor(method, path) {
					if param.TypeName() == "array" && param.ItemsTypeName() == "" {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "param %q for %q is a collection without an element type (array requires item definition)", param.Name, op.ID))
						continue
					}
					if param.In != "body" {
						if param.Items != nil {
							items := param.Items
							for items.TypeName() == "array" {
								if items.ItemsTypeName() == "" {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "param %q for %q is a collection without an element type (array requires item definition)", param.Name, op.ID))
									break
								}
								items = items.Items
							}
						}
					} else {
						// In: body
						if param.Schema != nil {
							if err := s.validateSchemaItems(*param.Schema, fmt.Sprintf("body param %q", param.Name), op.ID); err != nil {
								res.AddErrors(err)
							}
						}
					}
				}
			}

			var responses []spec.Response
			if op.Responses != nil {
				if op.Responses.Default != nil {
					responses = append(responses, *op.Responses.Default)
				}
				if op.Responses.StatusCodeResponses != nil {
					for _, v := range op.Responses.StatusCodeResponses {
						responses = append(responses, v)
					}
				}
			}

			for _, resp := range responses {
				// Response headers with array
				for hn, hv := range resp.Headers {
					if hv.TypeName() == "array" && hv.ItemsTypeName() == "" {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "header %q for %q is a collection without an element type (array requires items definition)", hn, op.ID))
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

// Verifies constraints on array type
func (s *SpecValidator) validateSchemaItems(schema spec.Schema, prefix, opID string) error {
	if !schema.Type.Contains("array") {
		return nil
	}

	if schema.Items == nil || schema.Items.Len() == 0 {
		return errors.New(errors.CompositeErrorCode, "%s for %q is a collection without an element type (array requires items definition)", prefix, opID)
	}

	if schema.Items.Schema != nil {
		schema = *schema.Items.Schema
		if _, err := compileRegexp(schema.Pattern); err != nil {
			return errors.New(errors.CompositeErrorCode, "%s for %q has invalid items pattern: %q", prefix, opID, schema.Pattern)
		}

		return s.validateSchemaItems(schema, prefix, opID)
	}

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
			res.Errors = append(res.Errors, errors.New(errors.CompositeErrorCode, "path param %q has no parameter definition", l))
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
			res.AddErrors(errors.New(errors.CompositeErrorCode, "path param %q is not present in path %q", p, path))
		}
	}

	return res
}

func (s *SpecValidator) validateReferenced() *Result {
	var res Result
	res.Merge(s.validateReferencedParameters())
	res.Merge(s.validateReferencedResponses())
	res.Merge(s.validateReferencedDefinitions())
	return &res
}

func (s *SpecValidator) validateReferencedParameters() *Result {
	// Each referenceable definition must have references.
	params := s.spec.Spec().Parameters
	if len(params) == 0 {
		return nil
	}

	expected := make(map[string]struct{})
	for k := range params {
		expected["#/parameters/"+jsonpointer.Escape(k)] = struct{}{}
	}
	for _, k := range s.analyzer.AllParameterReferences() {
		if _, ok := expected[k]; ok {
			delete(expected, k)
		}
	}

	if len(expected) == 0 {
		return nil
	}
	var result Result
	for k := range expected {
		result.AddErrors(errors.New(errors.CompositeErrorCode, "parameter %q is not used anywhere", k))
	}
	return &result
}

func (s *SpecValidator) validateReferencedResponses() *Result {
	// Each referenceable definition must have references.
	responses := s.spec.Spec().Responses
	if len(responses) == 0 {
		return nil
	}

	expected := make(map[string]struct{})
	for k := range responses {
		expected["#/responses/"+jsonpointer.Escape(k)] = struct{}{}
	}
	for _, k := range s.analyzer.AllResponseReferences() {
		if _, ok := expected[k]; ok {
			delete(expected, k)
		}
	}

	if len(expected) == 0 {
		return nil
	}
	var result Result
	for k := range expected {
		result.AddErrors(errors.New(errors.CompositeErrorCode, "response %q is not used anywhere", k))
	}
	return &result
}

func (s *SpecValidator) validateReferencedDefinitions() *Result {
	// Each referenceable definition must have references.
	defs := s.spec.Spec().Definitions
	if len(defs) == 0 {
		return nil
	}

	expected := make(map[string]struct{})
	for k := range defs {
		expected["#/definitions/"+jsonpointer.Escape(k)] = struct{}{}
	}
	for _, k := range s.analyzer.AllDefinitionReferences() {
		if _, ok := expected[k]; ok {
			delete(expected, k)
		}
	}

	if len(expected) == 0 {
		return nil
	}
	var result Result
	for k := range expected {
		result.AddErrors(errors.New(errors.CompositeErrorCode, "definition %q is not used anywhere", k))
	}
	return &result
}

func (s *SpecValidator) validateRequiredDefinitions() *Result {
	// Each definition property listed in the required array must be defined in the properties of the model
	res := new(Result)
	for d, v := range s.spec.Spec().Definitions {
		if v.Required != nil { // Safeguard
		REQUIRED:
			for _, pn := range v.Required {
				if _, ok := v.Properties[pn]; ok {
					continue
				}

				for pp := range v.PatternProperties {
					re, err := compileRegexp(pp)
					if err != nil {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "Pattern \"%q\" is invalid", pp))
						continue REQUIRED
					}
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

				res.AddErrors(errors.New(errors.CompositeErrorCode, "%q is present in required but not defined as property in definition %q", pn, d))
			}
		}
	}
	return res
}

func (s *SpecValidator) validateParameters() *Result {
	// - for each method, path is unique, regardless of path parameters
	//   e.g. GET:/petstore/{id}, GET:/petstore/{pet}, GET:/petstore are
	//   considered duplicate paths
	// - each parameter should have a unique `name` and `type` combination
	// - each operation should have only 1 parameter of type body
	// - there must be at most 1 parameter in body
	// - parameters with pattern property must specify valid patterns
	// - $ref in parameters must resolve
	// - path param must be required
	res := new(Result)
	rexGarbledPathSegment := mustCompileRegexp(`.*[{}\s]+.*`)
	for method, pi := range s.analyzer.Operations() {
		methodPaths := make(map[string]map[string]string)
		if pi != nil { // Safeguard
			for path, op := range pi {
				pathToAdd := stripParametersInPath(path)
				// Warn on garbled path afer param stripping
				if rexGarbledPathSegment.MatchString(pathToAdd) {
					res.AddWarnings(errors.New(errors.CompositeErrorCode, "path stripped from path parameters %s contains {,} or white space. This is probably no what you want.", pathToAdd))
				}
				// Check uniqueness of stripped paths
				if _, found := methodPaths[method][pathToAdd]; found {
					// Sort names for stable, testable output
					if strings.Compare(path, methodPaths[method][pathToAdd]) < 0 {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "path %s overlaps with %s", path, methodPaths[method][pathToAdd]))
					} else {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "path %s overlaps with %s", methodPaths[method][pathToAdd], path))
					}
				} else {
					if _, found := methodPaths[method]; !found {
						methodPaths[method] = map[string]string{}
					}
					methodPaths[method][pathToAdd] = path //Original non stripped path

				}

				ptypes := make(map[string]map[string]struct{})
				// Accurately report situations when more than 1 body param is declared (possibly unnamed)
				var bodyParams []string

				sw := s.spec.Spec()
				var paramNames []string

				// Check for duplicate parameters declaration in param section
				if op.Parameters != nil { // Safeguard
				PARAMETERS:
					for _, ppr := range op.Parameters {
						pr := ppr
						for pr.Ref.String() != "" {
							obj, _, err := pr.Ref.GetPointer().Get(sw)
							if err != nil {
								s.brokenRefs = true
								refPath := strings.Join([]string{"\"" + path + "\"", method}, ".")
								if ppr.Name != "" {
									refPath = strings.Join([]string{refPath, ppr.Name}, ".")
								}
								s.addPointerError(res, err, pr.Ref.String(), refPath)
								// Cannot continue on this one if ref not resolved
								if !continueOnErrors {
									// We have had enough. Please stop it
									break PARAMETERS
								} else {
									// We want to know more
									continue PARAMETERS
								}
							} else {
								if checkedObj, ok := s.checkedParamAssertion(obj, pr.Name, pr.In, op.ID, res); ok {
									pr = checkedObj
								} else {
									continue PARAMETERS
								}
							}
						}
						pnames, ok := ptypes[pr.In]
						if !ok {
							pnames = make(map[string]struct{})
							ptypes[pr.In] = pnames
						}

						_, ok = pnames[pr.Name]
						if ok {
							// TODO: test case to get there
							res.AddErrors(errors.New(errors.CompositeErrorCode, "duplicate parameter name %q for %q in operation %q", pr.Name, pr.In, op.ID))
						}
						pnames[pr.Name] = struct{}{}
					}
				}

				// NOTE: as for current, analyzer.ParamsFor() panics when at least
				// one $ref is broken (analyzer.paramsAsMap()).
				// Hence, we skip this if a $ref cannot resolved at this level.
				// Same with broken parameter type (e.g. Schema instead of Parameter)
				if s.brokenRefs || s.brokenParam {
					res.AddErrors(errors.New(errors.CompositeErrorCode, "some parameters definitions are broken in %q.%s. Cannot continue validating parameters for operation %s", path, method, op.ID))
				} else {
				PARAMETERS2:
					for _, ppr := range s.analyzer.ParamsFor(method, path) {
						pr := ppr
						for pr.Ref.String() != "" {
							obj, _, err := pr.Ref.GetPointer().Get(sw)
							if err != nil {
								// We should not be able to enter here, as (for now) ParamsFor() guarantees refs are resolved.
								refPath := strings.Join([]string{"\"" + path + "\"", method}, ".")
								if ppr.Name != "" {
									refPath = strings.Join([]string{refPath, ppr.Name}, ".")
								}
								s.addPointerError(res, err, pr.Ref.String(), refPath)
								// Cannot continue on this one if ref not resolved
								if !continueOnErrors {
									// We have had enough. Please stop it
									break PARAMETERS2
								} else {
									// We want to know more
									continue PARAMETERS2
								}
							} else {
								if checkedObj, ok := s.checkedParamAssertion(obj, pr.Name, pr.In, op.ID, res); ok {
									pr = checkedObj
								} else {
									continue PARAMETERS2
								}
							}
						}

						// Validate pattern for parameters with a pattern property
						if _, err := compileRegexp(pr.Pattern); err != nil {
							res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has invalid pattern in param %q: %q", op.ID, pr.Name, pr.Pattern))
						}

						// There must be at most one parameter in body: list them all
						if pr.In == "body" {
							bodyParams = append(bodyParams, fmt.Sprintf("%q", pr.Name))
						}

						if pr.In == "path" {
							paramNames = append(paramNames, pr.Name)
							// Path declared in path must have the required: true property
							if !pr.Required {
								res.AddErrors(errors.New(errors.CompositeErrorCode, "in operation %q,path param %q must be declared as required", op.ID, pr.Name))
							}
						}
					}
				}
				// There must be at most one body param
				if len(bodyParams) > 1 {
					sort.Strings(bodyParams)
					res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has more than 1 body param: %v", op.ID, bodyParams))
				}
				// Check uniqueness of parameters in path
				paramsInPath := extractPathParams(path)
				for i, p := range paramsInPath {
					for j, q := range paramsInPath {
						if p == q && i > j {
							res.AddErrors(errors.New(errors.CompositeErrorCode, "params in path %q must be unique: %q conflicts whith %q", path, p, q))
							break
						}

					}
				}

				// Warns about possible malformed params in path
				rexGarbledParam := mustCompileRegexp(`{.*[{}\s]+.*}`)
				for _, p := range paramsInPath {
					if rexGarbledParam.MatchString(p) {
						res.AddWarnings(errors.New(errors.CompositeErrorCode, "in path %q, param %q contains {,} or white space. Albeit not stricly illegal, this is probably no what you want", path, p))
					}
				}

				// Match params from path vs params from params section
				res.Merge(s.validatePathParamPresence(path, paramsInPath, paramNames))
			}
		}
	}
	return res
}

func stripParametersInPath(path string) string {
	// Returns a path stripped from all path parameters, with multiple or trailing slashes removed
	// Stripping is performed on a slash-separated basis, e.g '/a{/b}' remains a{/b} and not /a
	// - Trailing "/" make a difference, e.g. /a/ !~ /a (ex: canary/bitbucket.org/swagger.json)
	// - presence or absence of a parameter makes a difference, e.g. /a/{log} !~ /a/ (ex: canary/kubernetes/swagger.json)

	// Regexp to extract parameters from path, with surrounding {}.
	// NOTE: important non-greedy modifier
	rexParsePathParam := mustCompileRegexp(`{[^{}]+?}`)
	strippedSegments := []string{}

	for _, segment := range strings.Split(path, "/") {
		strippedSegments = append(strippedSegments, rexParsePathParam.ReplaceAllString(segment, "X"))
	}
	return strings.Join(strippedSegments, "/")
}

func extractPathParams(path string) (params []string) {
	// Extracts all params from a path, with surrounding "{}"
	rexParsePathParam := mustCompileRegexp(`{[^{}]+?}`)

	for _, segment := range strings.Split(path, "/") {
		for _, v := range rexParsePathParam.FindAllStringSubmatch(segment, -1) {
			params = append(params, v...)
		}
	}
	return
}

func (s *SpecValidator) validateReferencesValid() *Result {
	// each reference must point to a valid object
	res := new(Result)
	for _, r := range s.analyzer.AllRefs() {
		// TODO: test case
		if !r.IsValidURI(s.spec.SpecFilePath()) {
			res.AddErrors(errors.New(404, "invalid ref %q", r.String()))
		}
	}
	if !res.HasErrors() {
		// NOTE: with default settings, loads.Document.Expanded()
		// stops on first error. Anyhow, the expand option to continue
		// on errors fails to report errors at all.
		exp, err := s.spec.Expanded()
		if err != nil {
			res.AddErrors(fmt.Errorf("some references could not be resolved in spec. First found: %v", err))
		}
		s.expanded = exp
	}
	return res
}

func (s *SpecValidator) validateResponseExample(path string, r *spec.Response) *Result {
	// values provided as example in responses must validate the schema they examplify
	res := new(Result)

	// Recursively follow possible $ref's
	if r.Ref.String() != "" {
		nr, _, err := r.Ref.GetPointer().Get(s.spec.Spec())
		if err != nil {
			s.addPointerError(res, err, r.Ref.String(), strings.Join([]string{"\"" + path + "\"", r.ResponseProps.Schema.ID}, "."))
			return res
		}
		// Here we may expect type assertion to be guaranteed (not like in the Parameter case)
		rr := nr.(spec.Response)
		return s.validateResponseExample(path, &rr)
	}

	// NOTE: "examples" in responses vs "example" is a misleading construct in swagger
	if r.Examples != nil {
		if r.Schema != nil {
			if example, ok := r.Examples["application/json"]; ok {
				res.Merge(NewSchemaValidator(r.Schema, s.spec.Spec(), path, s.KnownFormats).Validate(example))
			}

			// TODO: validate other media types too
		}
	}
	return res
}

func (s *SpecValidator) validateExamplesValidAgainstSchema() *Result {
	// validates all examples provided in a spec
	// - values provides as Examples in a response must validate the response's schema
	// - TODO: examples for params, etc..
	res := new(Result)

	for _ /*method*/, pathItem := range s.analyzer.Operations() {
		if pathItem != nil { // Safeguard
			for path, op := range pathItem {
				// Check Examples in Responses
				if op.Responses != nil {
					if op.Responses.Default != nil {
						dr := op.Responses.Default
						res.Merge(s.validateResponseExample(path, dr))
						// TODO: Check example in default response schema
					}
					if op.Responses.StatusCodeResponses != nil { // Safeguard
						for _ /*code*/, r := range op.Responses.StatusCodeResponses {
							res.Merge(s.validateResponseExample(path, &r))
							// TODO: Check example in response schema
						}
					}
				}
			}
		}
	}
	return res
}

func (s *SpecValidator) validateDefaultValueValidAgainstSchema() *Result {
	// every default value that is specified must validate against the schema for that property
	// headers, items, parameters, schema

	res := new(Result)

	for method, pathItem := range s.analyzer.Operations() {
		if pathItem != nil { // Safeguard
			for path, op := range pathItem {
				// parameters
				var hasForm, hasBody bool
				if !s.brokenRefs && !s.brokenParam { // Safeguard (since ParamsFor() panics on brokenRefs||brokenParam)
				PARAMETERS:
					for _, pr := range s.analyzer.ParamsFor(method, path) {
						// expand ref is necessary
						param := pr
						for param.Ref.String() != "" {
							obj, _, err := param.Ref.GetPointer().Get(s.spec.Spec())
							if err != nil {
								// We should not be able to enter here, as (for now) ParamsFor() guarantees refs are resolved.
								refPath := strings.Join([]string{"\"" + path + "\"", method}, ".")
								if param.Name != "" {
									refPath = strings.Join([]string{refPath, param.Name}, ".")
								}
								s.addPointerError(res, err, pr.Ref.String(), refPath)
								// Cannot continue on this one if ref not resolved
								if !continueOnErrors {
									// We have had enough. Please stop it
									break PARAMETERS
								} else {
									// We want to know more
									continue PARAMETERS
								}
							} else {
								if checkedObj, ok := s.checkedParamAssertion(obj, pr.Name, pr.In, op.ID, res); ok {
									pr = checkedObj
								} else {
									continue PARAMETERS
								}
							}
						}

						if param.In == "formData" {
							if hasBody && !hasForm {
								res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has both formData and body parameters. Only one such In: type may be used for a given operation", op.ID))
							}
							hasForm = true
						}
						if param.In == "body" {
							if hasForm && !hasBody {
								res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has both body and formData parameters. Only one such In: type may be used for a given operation", op.ID))
							}
							hasBody = true
						}
						// Check simple parameters first
						// default values provided must validate against their schema
						if param.Default != nil && param.Schema == nil {
							res.AddWarnings(errors.New(errors.CompositeErrorCode, "%s in %s has a default but no valid schema", param.Name, param.In))
							// check param valid
							red := NewParamValidator(&param, s.KnownFormats).Validate(param.Default)
							if red != nil {
								if !red.IsValid() {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in %s does not validate its schema", param.Name, param.In))
								}
								res.Merge(red)

							}
						}

						// Recursively follows Items and Schemas
						if param.Items != nil {
							red := s.validateDefaultValueItemsAgainstSchema(param.Name, param.In, &param, param.Items)
							if red != nil {
								if !red.IsValid() {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s.items in %s does not validate its schema", param.Name, param.In))
								}
								res.Merge(red)
							}
						}

						if param.Schema != nil {
							red := s.validateDefaultValueSchemaAgainstSchema(param.Name, param.In, param.Schema)
							if red != nil {
								if !red.IsValid() {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in %s does not validate its schema", param.Name, param.In))
								}
								res.Merge(red)
							}
						}
					}
				}

				// Same constraint on default Responses
				if op.Responses != nil {
					if op.Responses.Default != nil {
						dr := op.Responses.Default
						if dr.Headers != nil { // Safeguard
							for nm, h := range dr.Headers {
								if h.Default != nil {
									red := NewHeaderValidator(nm, &h, s.KnownFormats).Validate(h.Default)
									if red != nil {
										if !red.IsValid() {
											res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in header %s in default response does not validate its schema", nm))
										}
										res.Merge(red)
									}
								}
								if h.Items != nil {
									red := s.validateDefaultValueItemsAgainstSchema(nm, "header", &h, h.Items)
									if red != nil {
										if !red.IsValid() {
											res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in header.items %s in default response does not validate its schema", nm))
										}
										res.Merge(red)
									}
								}
								if _, err := compileRegexp(h.Pattern); err != nil {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has invalid pattern in default header %q: %q", op.ID, nm, h.Pattern))
								}
							}
						}
						if dr.Schema != nil {
							red := s.validateDefaultValueSchemaAgainstSchema("default", "response", dr.Schema)
							if red != nil {
								if !red.IsValid() {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "default value in default response does not validate its schema"))
								}
								res.Merge(red)
							}
						}
					}
					if op.Responses.StatusCodeResponses != nil { // Safeguard
						for code, r := range op.Responses.StatusCodeResponses {
							for nm, h := range r.Headers {
								if h.Default != nil {
									red := NewHeaderValidator(nm, &h, s.KnownFormats).Validate(h.Default)
									if red != nil {
										if !red.IsValid() {
											res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in header %s in response %d does not validate its schema", nm, code))
										}
										res.Merge(red)
									}
								}
								if h.Items != nil {
									red := s.validateDefaultValueItemsAgainstSchema(nm, "header", &h, h.Items)
									if red != nil {
										if !red.IsValid() {
											res.AddErrors(errors.New(errors.CompositeErrorCode, "default value for %s in header.items %s in response %d does not validate its schema", nm, code))
										}
										res.Merge(red)
									}
								}
								if _, err := compileRegexp(h.Pattern); err != nil {
									res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has invalid pattern in %v's header %q: %q", op.ID, code, nm, h.Pattern))
								}
							}
							if r.Schema != nil {
								red := s.validateDefaultValueSchemaAgainstSchema(strconv.Itoa(code), "response", r.Schema)
								if red != nil {
									if !red.IsValid() {
										res.AddErrors(errors.New(errors.CompositeErrorCode, "default value in response %d does not validate its schema", code))
									}
									res.Merge(red)
								}
							}
						}
					}
				} else {
					// Empty op.ID means there is no meaningful operation: no need to report a specific message
					if op.ID != "" {
						res.AddErrors(errors.New(errors.CompositeErrorCode, "operation %q has no valid response", op.ID))
					}
				}
			}
		}
	}
	if s.spec.Spec().Definitions != nil { // Safeguard
		for nm, sch := range s.spec.Spec().Definitions {
			res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("definitions.%s", nm), "body", &sch))
		}
	}
	return res
}

func (s *SpecValidator) validateDefaultValueSchemaAgainstSchema(path, in string, schema *spec.Schema) *Result {
	res := new(Result)
	if schema != nil { // Safeguard
		if schema.Default != nil {
			res.Merge(NewSchemaValidator(schema, s.spec.Spec(), path+".default", s.KnownFormats).Validate(schema.Default))
		}
		if schema.Items != nil {
			if schema.Items.Schema != nil {
				res.Merge(s.validateDefaultValueSchemaAgainstSchema(path+".items.default", in, schema.Items.Schema))
			}
			if schema.Items.Schemas != nil { // Safeguard
				for i, sch := range schema.Items.Schemas {
					res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.items[%d].default", path, i), in, &sch))
				}
			}
		}
		if _, err := compileRegexp(schema.Pattern); err != nil {
			res.AddErrors(errors.New(errors.CompositeErrorCode, "%s in %s has invalid pattern: %q", path, in, schema.Pattern))
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
		if schema.AllOf != nil {
			for i, aoSch := range schema.AllOf {
				res.Merge(s.validateDefaultValueSchemaAgainstSchema(fmt.Sprintf("%s.allOf[%d]", path, i), in, &aoSch))
			}
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
			res.Merge(s.validateDefaultValueItemsAgainstSchema(path+"[0].default", in, root, items.Items))
		}
		if _, err := compileRegexp(items.Pattern); err != nil {
			res.AddErrors(errors.New(errors.CompositeErrorCode, "%s in %s has invalid pattern: %q", path, in, items.Pattern))
		}
	}
	return res
}

// $ref may not have siblings
// Spec: $ref siblings are ignored. So this check produces a warning
// TODO: check that $refs are only found in schemas
func (s *SpecValidator) validateRefNoSibling() *Result {
	return nil
}

func (s *SpecValidator) addPointerError(res *Result, err error, ref string, fromPath string) *Result {
	// Provides more context on error messages
	// reported by the jsoinpointer package
	if err != nil {
		res.AddErrors(fmt.Errorf("could not resolve reference in %s to $ref %s: %v", fromPath, ref, err))

	}
	return res
}

func (s *SpecValidator) checkedParamAssertion(obj interface{}, path, in, operation string, res *Result) (spec.Parameter, bool) {
	// Secure parameter type assertion and try to explain failure
	if checkedObj, ok := obj.(spec.Parameter); ok {
		return checkedObj, true
	}
	// Try to explain why... best guess
	if _, ok := obj.(spec.Schema); ok {
		// Story of issue#342
		res.AddWarnings(errors.New(errors.CompositeErrorCode, "$ref property should have no sibling in %q.%s", operation, path))
		// Schema took over Parameter for an unexplained reason
		res.AddErrors(errors.New(errors.CompositeErrorCode, "invalid definition as Schema for parameter %s in %s in operation %q", path, in, operation))
	} else {
		// Another structure replaced spec.Parametern
		res.AddErrors(errors.New(errors.CompositeErrorCode, "invalid definition for parameter %s in %s in operation %q", path, in, operation))
	}
	s.brokenParam = true
	return spec.Parameter{}, false
}
