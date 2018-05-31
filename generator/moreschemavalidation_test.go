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

package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// modelExpectations is a test structure to capture expected codegen lines of code
type modelExpectations struct {
	GeneratedFile    string
	ExpectedLines    []string
	NotExpectedLines []string
	ExpectedLogs     []string
	NotExpectedLogs  []string
	ExpectFailure    bool
}

// modelTestRun is a test structure to configure generations options to test a spec
type modelTestRun struct {
	FixtureOpts *GenOpts
	Definitions map[string]*modelExpectations
}

// AddExpectations adds expected / not expected sets of lines of code to the current run
func (r *modelTestRun) AddExpectations(file string, expectedCode, notExpectedCode, expectedLogs, notExpectedLogs []string) {
	k := strings.ToLower(swag.ToJSONName(strings.TrimSuffix(file, ".go")))
	if def, ok := r.Definitions[k]; ok {
		def.ExpectedLines = append(def.ExpectedLines, expectedCode...)
		def.NotExpectedLines = append(def.NotExpectedLines, notExpectedCode...)
		def.ExpectedLogs = append(def.ExpectedLogs, expectedLogs...)
		def.NotExpectedLogs = append(def.NotExpectedLogs, notExpectedLogs...)
		return
	}
	r.Definitions[k] = &modelExpectations{
		GeneratedFile:    file,
		ExpectedLines:    expectedCode,
		NotExpectedLines: notExpectedCode,
		ExpectedLogs:     expectedLogs,
		NotExpectedLogs:  notExpectedLogs,
	}
}

// ExpectedFor returns the map of model expectations from the run for a given model definition
func (r *modelTestRun) ExpectedFor(definition string) *modelExpectations {
	if def, ok := r.Definitions[strings.ToLower(definition)]; ok {
		return def
	}
	return nil
}

// modelFixture is a test structure to launch configurable test runs on a given spec
type modelFixture struct {
	SpecFile    string
	Description string
	Runs        []*modelTestRun
}

// Add adds a new run to the provided model fixture
func (f *modelFixture) AddRun(skipFlattenSpec bool) *modelTestRun {
	opts := &GenOpts{}
	opts.IncludeValidator = true
	opts.IncludeModel = true
	// sets gen options (e.g. flatten vs expand) - flatten is the default setting
	opts.ValidateSpec = false
	opts.FlattenSpec = !skipFlattenSpec
	opts.Spec = f.SpecFile
	if err := opts.EnsureDefaults(); err != nil {
		panic(err)
	}
	defs := make(map[string]*modelExpectations, 150)
	run := &modelTestRun{
		FixtureOpts: opts,
		Definitions: defs,
	}
	f.Runs = append(f.Runs, run)
	return run
}

// ExpectedBy returns the expectations from another run of the current fixture, recalled by its index in the list of planned runs
func (f *modelFixture) ExpectedFor(index int, definition string) *modelExpectations {
	if index > len(f.Runs)-1 {
		return nil
	}
	if def, ok := f.Runs[index].Definitions[strings.ToLower(definition)]; ok {
		return def
	}
	return nil
}

// newModelFixture is a test utility to build a new test plan for a spec file.
// The returned structure may be then used to add runs and expectations to each run.
func newModelFixture(specFile string, description string) *modelFixture {
	// lookup if already here
	for _, fix := range testedModels {
		if fix.SpecFile == specFile {
			return fix
		}
	}
	runs := make([]*modelTestRun, 0, 2)
	fix := &modelFixture{
		SpecFile:    specFile,
		Description: description,
		Runs:        runs,
	}
	testedModels = append(testedModels, fix)
	return fix
}

// all tested specs: init at the end of this source file
// you may append to those with different initXXX() funcs below.
var (
	testedModels []*modelFixture

	// convenient vars for (not) matching some lines
	noLines     []string
	todo        []string
	validatable []string
	warning     []string
)

func init() {
	testedModels = make([]*modelFixture, 0, 20)
	noLines = []string{}
	todo = []string{`TODO`}
	validatable = append(todo, `Validate(`)
	warning = []string{`warning`}
}

// initModelFixtures loads all tests to be performed
func initModelFixtures() {
	initFixtureSimpleAllOf()
	initFixtureComplexAllOf()
	initFixtureIsNullable()
	initFixtureItching()
	initFixtureAdditionalProps()
	initFixtureTuple()
	initFixture1479Part()
	initFixture1198()
	initFixture1042()
	initFixture1042V2()
	initFixture979()
	initFixture842()
	initFixture607()
	initFixture1336()
	initFixtureErrors()
	initFixture844Variations()
	initFixtureMoreAddProps()
	// a more stringent verification of this known fixture
	initTodolistSchemavalidation()
	initFixture1537()
	initFixture1537v2()

	// more maps and nullability checks
	initFixture15365()
	initFixtureNestedMaps()
}

/* Template initTxxx() to prepare and load a fixture:

func initTxxx() {
	// testing xxx.yaml with expand (--skip-flatten)
	f := newModelFixture("xxx.yaml", "A test blg")

	// makes a run with skipFlattenSpec=false
	thisRun := f.AddRun(false)

	// loads expectations for model abc
	thisRun.AddExpectations("abc.go", []string{
		`line {`,
		`	more codegen  		`,
		`}`,
	},
		// not expected
		noLines,
		// output in Log
		noLines,
		noLines)

	// loads expectations for model abcDef
	thisRun.AddExpectations("abc_def.go", []string{}, []string{}, noLines, noLines)
}

*/

func TestModelGenerateDefinition(t *testing.T) {
	// exercise the top level model generation func
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	fixtureSpec := "../fixtures/bugs/1487/fixture-is-nullable.yaml"
	assert := assert.New(t)
	gendir, erd := ioutil.TempDir(".", "model-test")
	defer func() {
		os.RemoveAll(gendir)
	}()
	if assert.NoError(erd) {
		opts := &GenOpts{}
		opts.IncludeValidator = true
		opts.IncludeModel = true
		// sets gen options (e.g. flatten vs expand) - flatten is the default setting
		opts.ValidateSpec = false
		opts.FlattenSpec = true
		opts.Spec = fixtureSpec
		opts.ModelPackage = "models"
		opts.Target = gendir
		if err := opts.EnsureDefaults(); err != nil {
			panic(err)
		}

		err := GenerateDefinition([]string{"thingWithNullableDates"}, opts)
		assert.NoErrorf(err, "Expected GenerateDefinition() to run without error")

		err = GenerateDefinition(nil, opts)
		assert.NoErrorf(err, "Expected GenerateDefinition() to run without error")

		opts.TemplateDir = gendir
		err = GenerateDefinition([]string{"thingWithNullableDates"}, opts)
		assert.NoErrorf(err, "Expected GenerateDefinition() to run without error")

		err = GenerateDefinition([]string{"thingWithNullableDates"}, nil)
		assert.Errorf(err, "Expected GenerateDefinition() return an error when no option is passed")

		opts.TemplateDir = "templates"
		err = GenerateDefinition([]string{"thingWithNullableDates"}, opts)
		assert.Errorf(err, "Expected GenerateDefinition() to croak about protected templates")

		opts.TemplateDir = ""
		err = GenerateDefinition([]string{"myAbsentDefinition"}, opts)
		assert.Errorf(err, "Expected GenerateDefinition() to return an error when the model is not in spec")

		opts.Spec = "pathToNowhere"
		err = GenerateDefinition([]string{"thingWithNullableDates"}, opts)
		assert.Errorf(err, "Expected GenerateDefinition() to return an error when the spec is not reachable")
	}
}

func TestMoreModelValidations(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()
	continueOnErrors := false
	initModelFixtures()

	assert := assert.New(t)

	for _, fixture := range testedModels {
		if fixture.SpecFile == "" {
			continue
		}
		fixtureSpec := fixture.SpecFile
		specDoc, err := loads.Spec(fixtureSpec)
		if !assert.NoErrorf(err, "unexpected failure loading spec %s: %v", fixtureSpec, err) {
			t.FailNow()
			return
		}
		for _, fixtureRun := range fixture.Runs {
			// NOTE: could run in parallel, if no package level common structures, such as a cache etc..
			opts := fixtureRun.FixtureOpts
			t.Logf("codegen for  %s (%s) - run with FlattenSpec=%t", fixtureSpec, fixture.Description, opts.FlattenSpec)

			// this is the expanded or flattened spec
			newSpecDoc, er0 := validateAndFlattenSpec(opts, specDoc)
			if !assert.NoErrorf(er0, "could not expand/flatten fixture %s: %v", fixtureSpec, er0) {
				t.FailNow()
				return
			}
			definitions := newSpecDoc.Spec().Definitions
			for k, fixtureExpectations := range fixtureRun.Definitions {
				// pick definition to test
				var schema *spec.Schema
				var definitionName string
				for def, s := range definitions {
					// please do not inject in fixtures with case conflicts on defs...
					// this one is just easier to retrieve model back from file names when capturing
					// the generated code.
					if strings.EqualFold(def, k) {
						schema = &s
						definitionName = def
						break
					}
				}
				if !assert.NotNil(schema, "expected to find definition %q in model fixture %s", k, fixtureSpec) {
					t.FailNow()
					return
				}

				// prepare assertions on log output (e.g. generation warnings)
				var logCapture bytes.Buffer
				if len(fixtureExpectations.ExpectedLogs) > 0 || len(fixtureExpectations.NotExpectedLogs) > 0 {
					log.SetOutput(&logCapture)
				} else {
					log.SetOutput(ioutil.Discard)
				}

				// generate the schema for this definition
				genModel, er1 := makeGenDefinition(definitionName, "models", *schema, newSpecDoc, opts)

				if fixtureExpectations.ExpectFailure && !assert.Errorf(er1, "Expected an error during generation of definition %q from spec fixture %s", k, fixtureSpec) {
					// expected an error here, and it has not happened
					if continueOnErrors {
						t.Fail()
						continue
					} else {
						t.FailNow()
						return
					}
				}
				if !assert.NoErrorf(er1, "could not generate model definition %q from spec fixture %s: %v", k, fixtureSpec, er1) {
					// expected smooth generation
					if continueOnErrors {
						t.Fail()
						continue
					} else {
						t.FailNow()
						return
					}
				}
				if len(fixtureExpectations.ExpectedLogs) > 0 || len(fixtureExpectations.NotExpectedLogs) > 0 {
					// assert logged output
					res := logCapture.String()
					for line, logLine := range fixtureExpectations.ExpectedLogs {
						if !assertInCode(t, strings.TrimSpace(logLine), res) {
							t.Logf("Log expected did not match for definition %q in fixture %s at (fixture) log line %d", k, fixtureSpec, line)
						}
					}
					for line, logLine := range fixtureExpectations.NotExpectedLogs {
						if !assertNotInCode(t, strings.TrimSpace(logLine), res) {
							t.Logf("Log unexpectedly matched for definition %q in fixture %s at (fixture) log line %d", k, fixtureSpec, line)
						}
					}
					if t.Failed() && !continueOnErrors {
						t.FailNow()
						return
					}
					log.SetOutput(ioutil.Discard)
				}

				// execute the model template with this schema
				buf := bytes.NewBuffer(nil)
				er2 := templates.MustGet("model").Execute(buf, genModel)
				if !assert.NoErrorf(er2, "could not render model template for definition %q in spec fixture %s: %v", k, fixtureSpec, er2) {
					if continueOnErrors {
						t.Fail()
						continue
					} else {
						t.FailNow()
						return
					}
				}
				outputName := fixtureExpectations.GeneratedFile
				if outputName == "" {
					outputName = swag.ToFileName(k) + ".go"
				}

				// run goimport, gofmt on the generated code
				formatted, er3 := opts.LanguageOpts.FormatContent(outputName, buf.Bytes())
				if !assert.NoErrorf(er3, "could not render model template for definition %q in spec fixture %s: %v", k, fixtureSpec, er2) {
					t.Fail()
					if continueOnErrors {
						t.Fail()
						continue
					} else {
						t.FailNow()
						return
					}
				}

				// asserts generated code (see fixture file)
				res := string(formatted)
				for line, codeLine := range fixtureExpectations.ExpectedLines {
					if !assertInCode(t, strings.TrimSpace(codeLine), res) {
						t.Logf("Code expected did not match for definition %q in fixture %s at (fixture) line %d", k, fixtureSpec, line)
					}
				}
				for line, codeLine := range fixtureExpectations.NotExpectedLines {
					if !assertNotInCode(t, strings.TrimSpace(codeLine), res) {
						t.Logf("code unexpectedly matched for definition %q in fixture %s at (fixture) line %d", k, fixtureSpec, line)
					}
				}
				if t.Failed() && !continueOnErrors {
					t.FailNow()
					return
				}
			}
		}
	}
}
