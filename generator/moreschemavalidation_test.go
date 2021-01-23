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
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func (m modelExpectations) ExpectLogs() bool {
	// does this test case assert output?
	return len(m.ExpectedLogs) > 0 || len(m.NotExpectedLogs) > 0
}

func (m modelExpectations) AssertModelLogs(t testing.TB, msg, definitionName, fixtureSpec string) {
	// assert logged output
	for line, logLine := range m.ExpectedLogs {
		if !assertInCode(t, strings.TrimSpace(logLine), msg) {
			t.Logf("log expected did not match for definition %q in fixture %s at (fixture) log line %d", definitionName, fixtureSpec, line)
		}
	}

	for line, logLine := range m.NotExpectedLogs {
		if !assertNotInCode(t, strings.TrimSpace(logLine), msg) {
			t.Logf("log unexpectedly matched for definition %q in fixture %s at (fixture) log line %d", definitionName, fixtureSpec, line)
		}
	}

	if t.Failed() {
		t.FailNow()
	}
}

func (m modelExpectations) AssertModelCodegen(t testing.TB, msg, definitionName, fixtureSpec string) {
	// assert generated code
	for line, codeLine := range m.ExpectedLines {
		if !assertInCode(t, strings.TrimSpace(codeLine), msg) {
			t.Logf("code expected did not match for definition %q in fixture %s at (fixture) line %d", definitionName, fixtureSpec, line)
		}
	}

	for line, codeLine := range m.NotExpectedLines {
		if !assertNotInCode(t, strings.TrimSpace(codeLine), msg) {
			t.Logf("code unexpectedly matched for definition %q in fixture %s at (fixture) line %d", definitionName, fixtureSpec, line)
		}
	}

	if t.Failed() {
		t.FailNow()
	}
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

func (r *modelTestRun) WithMinimalFlatten(minimal bool) *modelTestRun {
	r.FixtureOpts.FlattenOpts.Minimal = minimal
	return r
}

// modelFixture is a test structure to launch configurable test runs on a given spec
type modelFixture struct {
	SpecFile    string
	Description string
	Runs        []*modelTestRun
}

// Add adds a new run to the provided model fixture
func (f *modelFixture) AddRun(expandSpec bool) *modelTestRun {
	opts := &GenOpts{}
	opts.IncludeValidator = true
	opts.IncludeModel = true
	opts.ValidateSpec = false
	opts.Spec = f.SpecFile
	if err := opts.EnsureDefaults(); err != nil {
		panic(err)
	}

	// sets gen options (e.g. flatten vs expand) - full flatten is the default setting for this test (NOT the default CLI option!)
	opts.FlattenOpts.Expand = expandSpec
	opts.FlattenOpts.Minimal = false

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
	modelTestMutex = &sync.Mutex{} // mutex to protect log capture
	testedModels   []*modelFixture

	// convenient vars for (not) matching some lines
	noLines     []string
	todo        []string
	validatable []string
	warning     []string
)

func initSchemaValidationTest() {
	testedModels = make([]*modelFixture, 0, 50)
	noLines = []string{}
	todo = []string{`TODO`}
	validatable = append([]string{`Validate(`}, todo...)
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
	initFixtureDeepMaps()

	// format "byte" validation
	initFixture1548()

	// more tuples
	initFixtureSimpleTuple()

	// allOf with properties
	initFixture1617()

	// type realiasing
	initFixtureRealiasedTypes()

	// required base type
	initFixture1993()

	// allOf marshallers
	initFixture2071()

	// x-omitempty
	initFixture2116()

	// additionalProperties in base type (pending fix, non regression assertion only atm)
	initFixture2220()

	// allOf can be forced to non-nullable
	initFixture2364()

	// ReadOnly ContextValidate
	initFixture936ReadOnly()

	// required interface{}
	initFixture2081()

	// required map
	initFixture2300()

	// required $ref primitive
	initFixture2381()

	// required aliased primitive
	initFixture2400()

	// numerical validations
	initFixture2448()

	initFixtureGuardFormats()

	// min / maxProperties
	initFixture2444()

	// map of nullable array
	initFixture2494()
}

/* Template initTxxx() to prepare and load a fixture:

func initTxxx() {
	// testing xxx.yaml with expand (--with-expand)
	f := newModelFixture("xxx.yaml", "A test blg")

	// makes a run with expandSpec=false (full flattening)
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
		_ = os.RemoveAll(gendir)
	}()
	require.NoError(t, erd)

	opts := &GenOpts{}
	opts.IncludeValidator = true
	opts.IncludeModel = true
	opts.ValidateSpec = false
	opts.Spec = fixtureSpec
	opts.ModelPackage = "models"
	opts.Target = gendir
	if err := opts.EnsureDefaults(); err != nil {
		panic(err)
	}
	// sets gen options (e.g. flatten vs expand) - flatten is the default setting
	opts.FlattenOpts.Minimal = false

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

func TestMoreModelValidations(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	initModelFixtures()

	t.Logf("INFO: model specs tested: %d", len(testedModels))
	for _, toPin := range testedModels {
		fixture := toPin
		if fixture.SpecFile == "" {
			continue
		}
		fixtureSpec := fixture.SpecFile
		runTitle := strings.Join([]string{"codegen", strings.TrimSuffix(path.Base(fixtureSpec), path.Ext(fixtureSpec))}, "-")

		t.Run(runTitle, func(t *testing.T) {
			t.Parallel()
			log.SetOutput(ioutil.Discard)

			for _, fixtureRun := range fixture.Runs {
				opts := fixtureRun.FixtureOpts
				opts.Spec = fixtureSpec
				// this is the expanded or flattened spec
				newSpecDoc, err := opts.validateAndFlattenSpec()
				require.NoErrorf(t, err, "could not expand/flatten fixture %s: %v", fixtureSpec, err)

				definitions := newSpecDoc.Spec().Definitions
				for k, fixtureExpectations := range fixtureRun.Definitions {
					// pick definition to test
					definitionName, schema := findTestDefinition(k, definitions)
					require.NotNilf(t, schema, "expected to find definition %q in model fixture %s", k, fixtureSpec)

					checkDefinitionCodegen(t, definitionName, fixtureSpec, schema, newSpecDoc, opts, fixtureExpectations)
				}
			}
		})
	}
}

func findTestDefinition(k string, definitions spec.Definitions) (string, *spec.Schema) {
	var (
		schema         *spec.Schema
		definitionName string
	)

	for def, s := range definitions {
		// please do not inject fixtures with case conflicts on defs...
		// this one is just easier to retrieve model back from file names when capturing
		// the generated code.
		mangled := swag.ToJSONName(def)
		if strings.EqualFold(mangled, k) {
			definition := s
			schema = &definition
			definitionName = def
			break
		}
	}
	return definitionName, schema
}

func checkDefinitionCodegen(t testing.TB, definitionName, fixtureSpec string, schema *spec.Schema, specDoc *loads.Document, opts *GenOpts, fixtureExpectations *modelExpectations) {
	// prepare assertions on log output (e.g. generation warnings)
	var logCapture bytes.Buffer
	var msg string

	if fixtureExpectations.ExpectLogs() {
		// lock when capturing shared log resource (hopefully not for all testcases)
		modelTestMutex.Lock()
		log.SetOutput(&logCapture)

		defer func() {
			log.SetOutput(ioutil.Discard)
			modelTestMutex.Unlock()
		}()
	}

	// generate the schema for this definition
	genModel, err := makeGenDefinition(definitionName, "models", *schema, specDoc, opts)
	if fixtureExpectations.ExpectLogs() {
		msg = logCapture.String()
	}

	if fixtureExpectations.ExpectFailure {
		// expected an error here, and it has not happened
		require.Errorf(t, err, "Expected an error during generation of definition %q from spec fixture %s", definitionName, fixtureSpec)
	}

	// expected smooth generation
	require.NoErrorf(t, err, "could not generate model definition %q from spec fixture %s: %v", definitionName, fixtureSpec, err)

	if fixtureExpectations.ExpectLogs() {
		fixtureExpectations.AssertModelLogs(t, msg, definitionName, fixtureSpec)
	}

	// execute the model template with this schema
	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoErrorf(t, err, "could not render model template for definition %q in spec fixture %s: %v", definitionName, fixtureSpec, err)

	outputName := fixtureExpectations.GeneratedFile
	if outputName == "" {
		outputName = swag.ToFileName(definitionName) + ".go"
	}

	// run goimport, gofmt on the generated code
	formatted, err := opts.LanguageOpts.FormatContent(outputName, buf.Bytes())
	require.NoErrorf(t, err, "could not render model template for definition %q in spec fixture %s: %v", definitionName, fixtureSpec, err)

	// assert generated code (see fixture file)
	fixtureExpectations.AssertModelCodegen(t, string(formatted), definitionName, fixtureSpec)
}
