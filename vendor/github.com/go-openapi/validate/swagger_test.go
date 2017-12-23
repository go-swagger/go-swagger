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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

// Exercise validate will all tests cases from package go-swagger
// A copy of all fixtures available in in go-swagger/go-swagger
// is maintained in fixtures/go-swagger
func Test_GoSwaggerTestCases(t *testing.T) {
	if !enableLongTests {
		skipNotify(t)
		t.SkipNow()
	}
	// A list of test cases which fail on "swagger validate" at spec load time
	expectedLoadFailures := map[string]bool{
		"fixtures/go-swagger/bugs/342/fixture-342.yaml":   false,
		"fixtures/go-swagger/bugs/342/fixture-342-2.yaml": true,
	}

	// A list of test cases which fail on "swagger validate"
	expectedFailures := map[string]bool{
		"fixtures/go-swagger/bugs/1010/swagger.yml":                      true,
		"fixtures/go-swagger/bugs/103/swagger.json":                      true,
		"fixtures/go-swagger/bugs/106/swagger.json":                      true,
		"fixtures/go-swagger/bugs/1171/swagger.yaml":                     true,
		"fixtures/go-swagger/bugs/1238/swagger.yaml":                     true,
		"fixtures/go-swagger/bugs/1289/fixture-1289-2.yaml":              true,
		"fixtures/go-swagger/bugs/1289/fixture-1289.yaml":                true,
		"fixtures/go-swagger/bugs/193/spec2.json":                        true,
		"fixtures/go-swagger/bugs/195/swagger.json":                      true,
		"fixtures/go-swagger/bugs/248/swagger.json":                      true,
		"fixtures/go-swagger/bugs/249/swagger.json":                      true,
		"fixtures/go-swagger/bugs/342/fixture-342-2.yaml":                true,
		"fixtures/go-swagger/bugs/342/fixture-342.yaml":                  true,
		"fixtures/go-swagger/bugs/423/swagger.json":                      true,
		"fixtures/go-swagger/bugs/453/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/455/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/628/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/733/swagger.json":                      true,
		"fixtures/go-swagger/bugs/763/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/774/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/776/error.yaml":                        true,
		"fixtures/go-swagger/bugs/776/item.yaml":                         true,
		"fixtures/go-swagger/bugs/809/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/825/swagger.yml":                       true,
		"fixtures/go-swagger/bugs/890/path/health_check.yaml":            true,
		"fixtures/go-swagger/bugs/981/swagger.json":                      true,
		"fixtures/go-swagger/canary/docker/swagger.json":                 true,
		"fixtures/go-swagger/canary/ms-cog-sci/swagger.json":             true,
		"fixtures/go-swagger/codegen/azure-text-analyis.json":            true,
		"fixtures/go-swagger/codegen/issue72.json":                       true,
		"fixtures/go-swagger/codegen/simplesearch.yml":                   true,
		"fixtures/go-swagger/codegen/swagger-codegen-tests.json":         true,
		"fixtures/go-swagger/codegen/todolist.allparams.yml":             true,
		"fixtures/go-swagger/codegen/todolist.bodyparams.yml":            true,
		"fixtures/go-swagger/codegen/todolist.discriminators.yml":        true,
		"fixtures/go-swagger/codegen/todolist.enums.yml":                 true,
		"fixtures/go-swagger/codegen/todolist.models.yml":                true,
		"fixtures/go-swagger/codegen/todolist.responses.yml":             true,
		"fixtures/go-swagger/codegen/todolist.schemavalidation.yml":      true,
		"fixtures/go-swagger/codegen/todolist.simplepath.yml":            true,
		"fixtures/go-swagger/codegen/todolist.simple.yml":                true,
		"fixtures/go-swagger/codegen/todolist.url.basepath.yml":          true,
		"fixtures/go-swagger/codegen/todolist.url.simple.yml":            true,
		"fixtures/go-swagger/expansion/all-the-things.json":              true,
		"fixtures/go-swagger/expansion/circularRefs.json":                true,
		"fixtures/go-swagger/expansion/invalid-refs.json":                true,
		"fixtures/go-swagger/expansion/params.json":                      true,
		"fixtures/go-swagger/expansion/schemas1.json":                    true,
		"fixtures/go-swagger/expansion/schemas2.json":                    true,
		"fixtures/go-swagger/petstores/petstore-expanded.json":           true,
		"fixtures/go-swagger/petstores/petstore-simple.json":             true,
		"fixtures/go-swagger/petstores/petstore-with-external-docs.json": true,
		"fixtures/go-swagger/remotes/folder/folderInteger.json":          true,
		"fixtures/go-swagger/remotes/integer.json":                       true,
		"fixtures/go-swagger/remotes/subSchemas.json":                    true,
		"fixtures/go-swagger/specs/deeper/arrayProp.json":                true,
		"fixtures/go-swagger/specs/deeper/stringProp.json":               true,
		"fixtures/go-swagger/specs/refed.json":                           true,
		"fixtures/go-swagger/specs/resolution2.json":                     true,
		"fixtures/go-swagger/specs/resolution.json":                      true,
	}

	//t.SkipNow()
	state := continueOnErrors
	SetContinueOnErrors(true)
	defer func() {
		SetContinueOnErrors(state)
	}()
	if testGoSwaggerSpecs(t, "./fixtures/go-swagger", expectedFailures, expectedLoadFailures, true) != 0 {
		t.Fail()
	}
}

// A non regression test re "swagger validate" expectations
// Just validates all fixtures in ./fixtures/go-swagger (excluded codegen cases)
func testGoSwaggerSpecs(t *testing.T, path string, expectToFail, expectToFailOnLoad map[string]bool, haltOnErrors bool) (errs int) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			shouldNotLoad := false
			shouldFail := false
			if _, ok := expectToFailOnLoad[path]; ok {
				shouldNotLoad = expectToFailOnLoad[path]
			}
			if _, ok := expectToFail[path]; ok {
				shouldFail = expectToFail[path]
			}
			if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".json")) {
				// Checking invalid specs
				t.Logf("Testing messages for spec: %s", path)
				doc, err := loads.Spec(path)
				if shouldNotLoad {
					if !assert.Error(t, err, "Expected this spec not to load") {
						errs++
					}
				} else {
					if !assert.NoError(t, err, "Expected this spec to load without error") {
						errs++
					}
				}
				if errs > 0 {
					if haltOnErrors {
						return fmt.Errorf("Test halted: stop on error mode")
					}
					return nil
				}
				if shouldNotLoad {
					return nil
				}

				// Validate the spec document
				validator := NewSpecValidator(doc.Schema(), strfmt.Default)
				res, _ := validator.Validate(doc)
				if shouldFail {
					if !assert.False(t, res.IsValid(), "Expected this spec to be invalid") {
						errs++
					}
				} else {
					if !assert.True(t, res.IsValid(), "Expected this spec to be valid") {
						t.Logf("Errors reported by validation on %s", path)
						for _, e := range res.Errors {
							t.Log(e)
						}
						errs++
					}
				}
			}
			if haltOnErrors && errs > 0 {
				return fmt.Errorf("Test halted: stop on error mode")
			}
			return nil
		})
	if err != nil {
		t.Logf("%v", err)
		errs++
	}
	return
}
