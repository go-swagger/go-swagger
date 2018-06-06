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
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

// Enable long running tests by using cmd line arg,
// Usage: go test ... -args [-enable-long|-enable-go-swagger]
//
// -enable-long:       enable spec_test.go:TestIssue18 and messages_test.go:Test_Quality*
// -enable-go-swagger: enable non-regression tests against go-swagger fixtures (validation status)
//                     in swagger_test.go:Test_GoSwagger  (running about 110 specs...)
//
// If none enabled, these tests are skipped
//
// NOTE: replacing with go test -short and testing.Short() means that
// by default, every test is launched. With -enable-long, we just get the
// opposite...
var enableLongTests bool
var enableGoSwaggerTests bool

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	flag.BoolVar(&enableLongTests, "enable-long", false, "enable long runnning tests")
	flag.BoolVar(&enableGoSwaggerTests, "enable-go-swagger", false, "enable go-swagger non-regression test")
}

func skipNotify(t *testing.T) {
	t.Log("To enable this long running test, use -args -enable-long in your go test command line")
}

func TestSpec_ExpandResponseLocalFile(t *testing.T) {
	fp := filepath.Join("fixtures", "local_expansion", "spec.yaml")
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		if assert.NotNil(t, doc) {
			validator := NewSpecValidator(doc.Schema(), strfmt.Default)
			res, _ := validator.Validate(doc)
			assert.True(t, res.IsValid())
			assert.Empty(t, res.Errors)
		}
	}
}

func TestSpec_ExpandResponseRecursive(t *testing.T) {
	fp := filepath.Join("fixtures", "recursive_expansion", "spec.yaml")
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		if assert.NotNil(t, doc) {
			validator := NewSpecValidator(doc.Schema(), strfmt.Default)
			res, _ := validator.Validate(doc)
			assert.True(t, res.IsValid())
			assert.Empty(t, res.Errors)
		}
	}
}

// Spec with no path
func TestSpec_Issue52(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "52", "swagger.json")
	jstext, _ := ioutil.ReadFile(fp)

	// as json schema
	var sch spec.Schema
	if assert.NoError(t, json.Unmarshal(jstext, &sch)) {
		validator := NewSchemaValidator(spec.MustLoadSwagger20Schema(), nil, "", strfmt.Default)
		res := validator.Validate(&sch)
		assert.False(t, res.IsValid())
		assert.EqualError(t, res.Errors[0], ".paths in body is required")
	}

	// as swagger spec: path is set to nil
	// Here, validation stops as paths is initialized to empty
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.False(t, res.IsValid())
		assert.EqualError(t, res.Errors[0], ".paths in body is required")
	}
	// Here, validation continues, with invalid path from early checks as null.
	// This provides an additional (hopefully more informative) message.
	doc, err = loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		validator.SetContinueOnErrors(true)
		res, _ := validator.Validate(doc)
		assert.False(t, res.IsValid())
		var verifiedErrors []string
		for _, e := range res.Errors {
			verifiedErrors = append(verifiedErrors, e.Error())
		}
		assert.Len(t, verifiedErrors, 2, "Unexpected number of error messages returned")
		assert.Contains(t, verifiedErrors, ".paths in body is required")
		assert.Contains(t, verifiedErrors, "spec has no valid path defined")
	}
}

func TestSpec_Issue53(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "53", "noswagger.json")
	jstext, _ := ioutil.ReadFile(fp)

	// as json schema
	var sch spec.Schema
	if assert.NoError(t, json.Unmarshal(jstext, &sch)) {
		validator := NewSchemaValidator(spec.MustLoadSwagger20Schema(), nil, "", strfmt.Default)
		res := validator.Validate(&sch)
		assert.False(t, res.IsValid())
		assert.EqualError(t, res.Errors[0], ".swagger in body is required")
	}

	// as swagger despec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		if assert.False(t, res.IsValid()) {
			assert.EqualError(t, res.Errors[0], ".swagger in body is required")
		}
	}
}

func TestSpec_Issue62(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "62", "swagger.json")

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.NotEmpty(t, res.Errors)
		assert.True(t, res.HasErrors())
	}
}

func TestSpec_Issue63(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "63", "swagger.json")

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.True(t, res.IsValid())
	}
}

func TestSpec_Issue61_MultipleRefs(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "61", "multiple-refs.json")

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors)
		assert.True(t, res.IsValid())
	}
}

func TestSpec_Issue61_ResolvedRef(t *testing.T) {
	fp := filepath.Join("fixtures", "bugs", "61", "unresolved-ref-for-name.json")

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors)
		assert.True(t, res.IsValid())
	}
}

// No error with this one
func TestSpec_Issue123(t *testing.T) {
	path := "swagger.yml"
	fp := filepath.Join("fixtures", "bugs", "123", path)

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.True(t, res.IsValid())

		var verifiedErrors []string
		for _, e := range res.Errors {
			verifiedErrors = append(verifiedErrors, e.Error())
		}
		switch {
		case strings.Contains(path, "swagger.yml"):
			assert.Empty(t, verifiedErrors)
		default:
			t.Logf("Returned error messages: %v", verifiedErrors)
			t.Fatal("fixture not tested. Please add assertions for messages")
		}

		if DebugTest && t.Failed() {
			if len(verifiedErrors) > 0 {
				t.Logf("DEVMODE: Returned error messages validating %s ", path)
				for _, v := range verifiedErrors {
					t.Logf("%s", v)
				}
			}
		}
	}
}

func TestSpec_Issue6(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("fixtures", "bugs", "6", "*.json"))
	for _, path := range files {
		t.Logf("Tested spec=%s", path)
		doc, err := loads.Spec(path)
		if assert.NoError(t, err) {
			validator := NewSpecValidator(doc.Schema(), strfmt.Default)
			res, _ := validator.Validate(doc)
			assert.False(t, res.IsValid())

			var verifiedErrors []string
			for _, e := range res.Errors {
				verifiedErrors = append(verifiedErrors, e.Error())
			}
			switch {
			case strings.Contains(path, "empty-responses.json"):
				assert.Contains(t, verifiedErrors, "\"paths./foo.get.responses\" must not validate the schema (not)")
				assert.Contains(t, verifiedErrors, "paths./foo.get.responses in body should have at least 1 properties")
			case strings.Contains(path, "no-responses.json"):
				assert.Contains(t, verifiedErrors, "paths./foo.get.responses in body is required")
			default:
				t.Logf("Returned error messages: %v", verifiedErrors)
				t.Fatal("fixture not tested. Please add assertions for messages")
			}
			if DebugTest && t.Failed() {
				if len(verifiedErrors) > 0 {
					t.Logf("DEVMODE:Returned error messages validating %s ", path)
					for _, v := range verifiedErrors {
						t.Logf("%s", v)
					}
				}
			}
		}
	}
}

// check if invalid patterns are indeed invalidated
func TestSpec_Issue18(t *testing.T) {
	if !enableLongTests {
		skipNotify(t)
		t.SkipNow()
	}
	files, _ := filepath.Glob(filepath.Join("fixtures", "bugs", "18", "*.json"))
	for _, path := range files {
		t.Logf("Tested spec=%s", path)
		doc, err := loads.Spec(path)
		if assert.NoError(t, err) {
			validator := NewSpecValidator(doc.Schema(), strfmt.Default)
			validator.SetContinueOnErrors(true)
			res, _ := validator.Validate(doc)
			assert.False(t, res.IsValid())

			var verifiedErrors []string
			for _, e := range res.Errors {
				verifiedErrors = append(verifiedErrors, e.Error())
			}
			switch {
			case strings.Contains(path, "headerItems.json"):
				assert.Contains(t, verifiedErrors, "X-Foo in header has invalid pattern: \")<-- bad pattern\"")
			case strings.Contains(path, "headers.json"):
				assert.Contains(t, verifiedErrors, "in operation \"\", header X-Foo for default response has invalid pattern \")<-- bad pattern\": error parsing regexp: unexpected ): `)<-- bad pattern`")
				//  in operation \"\", header X-Foo for default response has invalid pattern \")<-- bad pattern\": error parsing regexp: unexpected ): `)<-- bad pattern`
				assert.Contains(t, verifiedErrors, "in operation \"\", header X-Foo for response 402 has invalid pattern \")<-- bad pattern\": error parsing regexp: unexpected ): `)<-- bad pattern`")
				//  in operation "", header X-Foo for response 402 has invalid pattern ")<-- bad pattern": error parsing regexp: unexpected ): `)<-- bad pattern`

			case strings.Contains(path, "paramItems.json"):
				assert.Contains(t, verifiedErrors, "body param \"user\" for \"\" has invalid items pattern: \")<-- bad pattern\"")
				// Updated message: from "user.items in body has invalid pattern: \")<-- bad pattern\"" to:
				assert.Contains(t, verifiedErrors, "default value for user in body does not validate its schema")
				assert.Contains(t, verifiedErrors, "user.items.default in body has invalid pattern: \")<-- bad pattern\"")
			case strings.Contains(path, "parameters.json"):
				assert.Contains(t, verifiedErrors, "operation \"\" has invalid pattern in param \"userId\": \")<-- bad pattern\"")
			case strings.Contains(path, "schema.json"):
				// TODO: strange that the text does not say response "200"...
				assert.Contains(t, verifiedErrors, "200 in response has invalid pattern: \")<-- bad pattern\"")
			default:
				t.Logf("Returned error messages: %v", verifiedErrors)
				t.Fatal("fixture not tested. Please add assertions for messages")
			}

			if DebugTest && t.Failed() {
				if len(verifiedErrors) > 0 {
					t.Logf("DEVMODE: Returned error messages validating %s ", path)
					for _, v := range verifiedErrors {
						t.Logf("%s", v)
					}
				}
			}
		}
	}
}

// check if a fragment path parameter is recognized, without error
func TestSpec_Issue39(t *testing.T) {
	path := "swagger.yml"
	fp := filepath.Join("fixtures", "bugs", "39", path)

	// as swagger spec
	doc, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.True(t, res.IsValid())

		var verifiedErrors []string
		for _, e := range res.Errors {
			verifiedErrors = append(verifiedErrors, e.Error())
		}
		switch {
		case strings.Contains(path, "swagger.yml"):
			assert.Empty(t, verifiedErrors)
		default:
			t.Logf("Returned error messages: %v", verifiedErrors)
			t.Fatal("fixture not tested. Please add assertions for messages")
		}
		if DebugTest && t.Failed() {
			if len(verifiedErrors) > 0 {
				t.Logf("DEVMODE: Returned error messages validating %s ", path)
				for _, v := range verifiedErrors {
					t.Logf("%s", v)
				}
			}
		}
	}
}

func TestSpec_ValidateDuplicatePropertyNames(t *testing.T) {
	// simple allOf
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "duplicateprops.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateDuplicatePropertyNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)

	}

	// nested allOf
	doc, err = loads.Spec(filepath.Join("fixtures", "validation", "nestedduplicateprops.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateDuplicatePropertyNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)

	}
}

func TestSpec_ValidateNonEmptyPathParameterNames(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "empty-path-param-name.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateNonEmptyPathParamNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)

	}
}

func TestSpec_ValidateCircularAncestry(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "direct-circular-ancestor.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateDuplicatePropertyNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "validation", "indirect-circular-ancestor.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateDuplicatePropertyNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "validation", "recursive-circular-ancestor.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		res := validator.validateDuplicatePropertyNames()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)
	}

}

func TestSpec_ValidateReferenced(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "valid-referenced.yml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		validator.analyzer = analysis.New(doc.Spec())
		res := validator.validateReferenced()
		assert.Empty(t, res.Errors)
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "validation", "invalid-referenced.yml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		validator.analyzer = analysis.New(doc.Spec())
		res := validator.validateReferenced()
		assert.Empty(t, res.Errors)
		assert.NotEmpty(t, res.Warnings)
		assert.Len(t, res.Warnings, 3)
	}
}

func TestSpec_ValidateReferencesValid(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "valid-ref.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		validator.analyzer = analysis.New(doc.Spec())
		res := validator.validateReferencesValid()
		assert.Empty(t, res.Errors)
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "validation", "invalid-ref.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
		validator.spec = doc
		validator.analyzer = analysis.New(doc.Spec())
		res := validator.validateReferencesValid()
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)
	}
}

func TestSpec_ValidateRequiredDefinitions(t *testing.T) {
	doc, _ := loads.Analyzed(PetStoreJSONMessage, "")
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	res := validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	// properties
	sw := doc.Spec()
	def := sw.Definitions["Tag"]
	def.Required = append(def.Required, "type")
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)

	// pattern properties
	def.PatternProperties = make(map[string]spec.Schema)
	def.PatternProperties["ty.*"] = *spec.StringProperty()
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	def.PatternProperties = make(map[string]spec.Schema)
	def.PatternProperties["^ty.$"] = *spec.StringProperty()
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)

	// additional properties
	def.PatternProperties = nil
	def.AdditionalProperties = &spec.SchemaOrBool{Allows: true}
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	def.AdditionalProperties = &spec.SchemaOrBool{Allows: false}
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)
}

func TestSpec_ValidateParameters(t *testing.T) {
	doc, _ := loads.Analyzed(PetStoreJSONMessage, "")
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	res := validator.validateParameters()
	assert.Empty(t, res.Errors)

	sw := doc.Spec()
	sw.Paths.Paths["/pets"].Get.Parameters = append(sw.Paths.Paths["/pets"].Get.Parameters, *spec.QueryParam("limit").Typed("string", ""))
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)

	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	sw = doc.Spec()
	sw.Paths.Paths["/pets"].Post.Parameters = append(sw.Paths.Paths["/pets"].Post.Parameters, *spec.BodyParam("fake", spec.RefProperty("#/definitions/Pet")))
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 1)
	assert.Contains(t, res.Errors[0].Error(), "has more than 1 body param")

	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	sw = doc.Spec()
	pp := sw.Paths.Paths["/pets/{id}"]
	pp.Delete = nil
	var nameParams []spec.Parameter
	for _, p := range pp.Parameters {
		if p.Name == "id" {
			p.Name = "name"
			nameParams = append(nameParams, p)
		}
	}
	pp.Parameters = nameParams
	sw.Paths.Paths["/pets/{name}"] = pp

	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 1)
	assert.Contains(t, res.Errors[0].Error(), "overlaps with")

	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	sw = doc.Spec()
	pp = sw.Paths.Paths["/pets/{id}"]
	pp.Delete = nil
	pp.Get.Parameters = nameParams
	pp.Parameters = nil
	sw.Paths.Paths["/pets/{id}"] = pp

	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 2)
	assert.Contains(t, res.Errors[1].Error(), "is not present in path \"/pets/{id}\"")
	assert.Contains(t, res.Errors[0].Error(), "has no parameter definition")
}

func TestSpec_ValidateItems(t *testing.T) {
	doc, _ := loads.Analyzed(PetStoreJSONMessage, "")
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	res := validator.validateItems()
	assert.Empty(t, res.Errors)

	// in operation parameters
	sw := doc.Spec()
	sw.Paths.Paths["/pets"].Get.Parameters[0].Type = "array"
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items = spec.NewItems().Typed("string", "")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items = spec.NewItems().Typed("array", "")
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items.Items = spec.NewItems().Typed("string", "")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	// in global parameters
	sw.Parameters = make(map[string]spec.Parameter)
	sw.Parameters["other"] = *spec.SimpleArrayParam("other", "array", "csv")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	//pp := spec.SimpleArrayParam("other", "array", "")
	//pp.Items = nil
	//sw.Parameters["other"] = *pp
	//res = validator.validateItems()
	//assert.NotEmpty(t, res.Errors)

	// in shared path object parameters
	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	sw = doc.Spec()

	pa := sw.Paths.Paths["/pets"]
	pa.Parameters = []spec.Parameter{*spec.SimpleArrayParam("another", "array", "csv")}
	sw.Paths.Paths["/pets"] = pa
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	pa = sw.Paths.Paths["/pets"]
	pp := spec.SimpleArrayParam("other", "array", "")
	pp.Items = nil
	pa.Parameters = []spec.Parameter{*pp}
	sw.Paths.Paths["/pets"] = pa
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in body param schema
	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	pa.Post.Parameters[0].Schema = spec.ArrayProperty(nil)
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in response headers
	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	rp := pa.Post.Responses.StatusCodeResponses[200]
	var hdr spec.Header
	hdr.Type = "array"
	rp.Headers = make(map[string]spec.Header)
	rp.Headers["X-YADA"] = hdr
	pa.Post.Responses.StatusCodeResponses[200] = rp
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in response schema
	doc, _ = loads.Analyzed(PetStoreJSONMessage, "")
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	rp = pa.Post.Responses.StatusCodeResponses[200]
	rp.Schema = spec.ArrayProperty(nil)
	pa.Post.Responses.StatusCodeResponses[200] = rp
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)
}

// Reuse known validated cases through the higher level Spec() call
func TestSpec_ValidDoc(t *testing.T) {
	fp := filepath.Join("fixtures", "local_expansion", "spec.yaml")
	doc2, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		err := Spec(doc2, strfmt.Default)
		assert.NoError(t, err)
	}
}

// Check higher level behavior on invalid spec doc
func TestSpec_InvalidDoc(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "default", "invalid-default-value-parameter.json"))
	if assert.NoError(t, err) {
		err := Spec(doc, strfmt.Default)
		assert.Error(t, err)
	}
}

func TestSpec_Validate_InvalidInterface(t *testing.T) {
	fp := filepath.Join("fixtures", "local_expansion", "spec.yaml")
	doc2, err := loads.Spec(fp)
	if assert.NoError(t, err) {
		if assert.NotNil(t, doc2) {
			validator := NewSpecValidator(doc2.Schema(), strfmt.Default)
			bug := "bzzz"
			res, _ := validator.Validate(bug)
			assert.NotEmpty(t, res.Errors)
			assert.Contains(t, res.Errors[0].Error(), "can only validate spec.Document objects")
		}
	}
}

func TestSpec_ValidateBodyFormDataParams(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "validation", "invalid-formdata-body-params.json"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.NotEmpty(t, res.Errors)
		assert.Len(t, res.Errors, 1)
	}
}

func TestSpec_Issue73(t *testing.T) {
	doc, err := loads.Spec(filepath.Join("fixtures", "bugs", "73", "fixture-swagger.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, " in fixture-swagger.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "73", "fixture-swagger-2.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-swagger-2.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "73", "fixture-swagger-3.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-swagger-3.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "73", "fixture-swagger-good.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, " in fixture-swagger-good.yaml")
	}
}

func TestSpec_Issue1341(t *testing.T) {
	// testing recursive walk with defaults and examples
	doc, err := loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341-good.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, " in fixture-1341-good.yaml")
		assert.Len(t, res.Warnings, 1, " in fixture-1341-good.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-1341.yaml")
		assert.Empty(t, res.Warnings, "in fixture-1341.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341-2.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-1341-2.yaml")
		assert.Empty(t, res.Warnings, "in fixture-1341-2.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341-3.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-1341-3.yaml")
		assert.Len(t, res.Warnings, 4, "in fixture-1341-3.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341-4.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Empty(t, res.Errors, "in fixture-1341-4.yaml")
		assert.Empty(t, res.Warnings, "in fixture-1341-4.yaml")
	}

	doc, err = loads.Spec(filepath.Join("fixtures", "bugs", "1341", "fixture-1341-5.yaml"))
	if assert.NoError(t, err) {
		validator := NewSpecValidator(doc.Schema(), strfmt.Default)
		res, _ := validator.Validate(doc)
		assert.Len(t, res.Errors, 4, "in fixture-1341-5.yaml")
		assert.Empty(t, res.Warnings, "in fixture-1341-5.yaml")
	}
}
