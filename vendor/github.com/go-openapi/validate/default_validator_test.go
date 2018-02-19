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
	"path/filepath"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestDefault_ValidatePetStore(t *testing.T) {
	doc, _ := loads.Analyzed(PetStoreJSONMessage, "")
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
	validator.spec = doc
	validator.analyzer = analysis.New(doc.Spec())
	myDefaultValidator := &defaultValidator{SpecValidator: validator}
	res := myDefaultValidator.Validate()
	assert.Empty(t, res.Errors)
}

func TestDefault_ValidateDefaults(t *testing.T) {
	tests := []string{
		"parameter",
		"parameter-required",
		"parameter-ref",
		"parameter-items",
		"header",
		"header-items",
		"schema",
		"schema-ref",
		"schema-additionalProperties",
		"schema-patternProperties",
		"schema-items",
		"schema-allOf",
		"parameter-schema",
		"default-response",
		"header-response",
		"header-items-default-response",
		"header-items-response",
		"header-pattern",
		"header-badpattern",
		"schema-items-allOf",
		"response-ref",
	}

	for _, tt := range tests {
		path := filepath.Join("fixtures", "validation", "default", "valid-default-value-"+tt+".json")
		if DebugTest {
			t.Logf("Testing valid default values for: %s", path)
		}
		doc, err := loads.Spec(path)
		if assert.NoError(t, err) {
			validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
			validator.spec = doc
			validator.analyzer = analysis.New(doc.Spec())
			myDefaultValidator := &defaultValidator{SpecValidator: validator}
			res := myDefaultValidator.Validate()
			assert.Empty(t, res.Errors, tt+" should not have errors")
			// Special case: warning only
			if tt == "parameter-required" {
				warns := []string{}
				for _, w := range res.Warnings {
					warns = append(warns, w.Error())
				}
				assert.Contains(t, warns, "limit in query has a default value and is required as parameter")
			}
		}

		path = filepath.Join("fixtures", "validation", "default", "invalid-default-value-"+tt+".json")
		if DebugTest {
			t.Logf("Testing invalid default values for: %s", path)
		}
		doc, err = loads.Spec(path)
		if assert.NoError(t, err) {
			validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), strfmt.Default)
			validator.spec = doc
			validator.analyzer = analysis.New(doc.Spec())
			myDefaultValidator := &defaultValidator{SpecValidator: validator}
			res := myDefaultValidator.Validate()
			assert.NotEmpty(t, res.Errors, tt+" should have errors")
			// Update: now we have an additional message to explain it's all about a default value
			// Example:
			// - default value for limit in query does not validate its Schema
			// - limit in query must be of type integer: "string"]
			assert.True(t, len(res.Errors) >= 1, tt+" should have at least 1 error")
		}
	}
}

func TestDefault_EdgeCase(t *testing.T) {
	// Testing guards
	var myDefaultvalidator *defaultValidator
	res := myDefaultvalidator.Validate()
	assert.True(t, res.IsValid())

	myDefaultvalidator = &defaultValidator{}
	res = myDefaultvalidator.Validate()
	assert.True(t, res.IsValid())
}
