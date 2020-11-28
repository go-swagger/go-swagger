package generator

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndTest(t *testing.T) {
	defer discardOutput()()

	cwd := testCwd(t)
	const root = "generated"
	defer func() {
		_ = os.RemoveAll(filepath.Join(cwd, root))
	}()

	t.Run("server build", func(t *testing.T) {
		for name, cas := range generateFixtures(t) {
			thisCas := cas
			thisName := name

			t.Run(thisName, func(t *testing.T) {
				t.Parallel()

				log.SetOutput(ioutil.Discard)
				defer thisCas.warnFailed(t)

				// default opts
				opts := testGenOpts()

				// create directory layout, defer clean
				defer thisCas.prepareTarget(t, thisName, "server_test", root, opts)()

				// preparation before generation
				if thisCas.prepare != nil {
					thisCas.prepare(t, opts)
				}

				t.Logf("generating test server at: %s, from %s", opts.Target, opts.Spec)

				err := GenerateServer("", nil, nil, opts)
				if thisCas.wantError {
					require.Errorf(t, err, "expected an error for server build fixture: %s", opts.Spec)
				} else {
					require.NoError(t, err, "unexpected error for server build fixture: %s", opts.Spec)
				}

				// verify
				if thisCas.verify != nil {
					thisCas.verify(t, opts.Target)
				}

				// fixture-specific clean
				if thisCas.clean != nil {
					thisCas.clean()
				}
			})
		}
	})
}

type generateFixture struct {
	name      string
	spec      string
	target    string
	wantError bool
	prepare   func(testing.TB, *GenOpts)
	verify    func(testing.TB, string)
	clean     func()
}

func (f generateFixture) base(t testing.TB, root string) (string, func()) {
	// base generation target
	cwd := testCwd(t)

	base := filepath.Join(cwd, root)
	require.NoErrorf(t, os.MkdirAll(base, 0700), "error in test creating target dir")

	generated, err := ioutil.TempDir(base, "generated")
	require.NoErrorf(t, err, "error in test creating temp dir")

	return generated, func() {
		_ = os.RemoveAll(generated)
	}
}

func (f generateFixture) prepareTarget(t testing.TB, name, base, root string, opts *GenOpts) func() {
	if name == "" {
		name = f.name
	}

	spec := filepath.FromSlash(f.spec)
	opts.Spec = spec

	generated, clean := f.base(t, root)

	if f.target == "" {
		opts.Target = filepath.Join(generated, opts.LanguageOpts.ManglePackageName(name, base))
	} else {
		opts.Target = filepath.Join(generated, filepath.Base(f.target))
	}

	require.NoErrorf(t, os.MkdirAll(opts.Target, 0700), "error in test creating target dir")

	return clean
}

func (f generateFixture) warnFailed(t testing.TB) func() {
	return func() {
		if t.Failed() {
			t.Log("ERROR: generation failed")
		}
	}
}

func generateFixtures(t testing.TB) map[string]generateFixture {
	return map[string]generateFixture{
		"issue 1943": {
			spec:   "../fixtures/bugs/1943/fixture-1943.yaml",
			target: "../fixtures/bugs/1943",
			prepare: func(_ testing.TB, opts *GenOpts) {
				input, err := ioutil.ReadFile("../fixtures/bugs/1943/datarace_test.go")
				require.NoError(t, err)

				// rewrite imports for the relocated test program
				cwd := testCwd(t)
				rebased := bytes.ReplaceAll(
					input,
					[]byte("/fixtures/bugs/1943"),
					[]byte(filepath.ToSlash(strings.TrimPrefix(opts.Target, filepath.Dir(cwd)))),
				)

				require.NoError(t, ioutil.WriteFile(filepath.Join(opts.Target, "datarace_test.go"), rebased, 0600))
				opts.ExcludeSpec = false
			},
			verify: func(t testing.TB, target string) {
				if runtime.GOOS == "windows" {
					// don't run race tests on Appveyor CI
					t.Logf("warn: race test skipped on windows")
					return
				}

				const packages = "./..."
				testPrg := "datarace_test.go"

				goExecInDir(t, target, "get", packages)

				t.Log("running data race test on generated server")
				goExecInDir(t, target, "test", "-v", "-race", testPrg)
			},
		},
		"packages_mangling": {
			spec: "../fixtures/bugs/2111/fixture-2111.yaml",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				require.True(t, fileExists(target, defaultServerTarget))
				assert.True(t, fileExists(filepath.Join(target, "cmd", "unsafe-tag-names-server"), "main.go"))

				srvTarget := filepath.Join(target, defaultServerTarget)
				opsTarget := filepath.Join(srvTarget, defaultOperationsTarget)
				require.True(t, fileExists(opsTarget, ""))

				for _, fileOrDir := range []string{
					"abc_linux", "abc_test",
					"api",
					"custom",
					"hash_tag_donuts",
					"nr123abc", "nr_at_donuts", "plus_donuts",
					"strfmt",
					"forced",
					"gtl",
					"nr12nasty",
					"override",
					"get_notag.go",
					"operationsops",
				} {
					assert.True(t, fileExists(opsTarget, fileOrDir))
				}

				buf, err := ioutil.ReadFile(filepath.Join(srvTarget, "configure_unsafe_tag_names.go"))
				require.NoError(t, err)

				code := string(buf)

				// assert imports, with deconfliction
				cwd := testCwd(t)
				base := path.Join("github.com", "go-swagger", "go-swagger",
					filepath.ToSlash(strings.TrimPrefix(target, filepath.Dir(cwd))),
				)

				baseImport := path.Join(base, `restapi/operations`)
				assertImports(t, baseImport, code)

				assertInCode(t, `api.APIGetConflictHandler = apiops.GetConflictHandlerFunc(`, code)
				assertInCode(t, `api.StrfmtGetAnotherConflictHandler = strfmtops.GetAnotherConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetNotagHandler = operations.GetNotagHandlerFunc(`, code)

				buf2, err := ioutil.ReadFile(filepath.Join(opsTarget, "unsafe_tag_names_api.go"))
				require.NoError(t, err)

				api := string(buf2)
				assertImports(t, baseImport, api)

				assertInCode(t, `APIGetConflictHandler: apiops.GetConflictHandlerFunc(func(params apiops.GetConflictParams) middleware.Responder {`, api)
				assertInCode(t, `StrfmtGetAnotherConflictHandler: strfmtops.GetAnotherConflictHandlerFunc(func(params strfmtops.GetAnotherConflictParams) middleware.Responder {`, api)
				assertInCode(t, `GetNotagHandler: GetNotagHandlerFunc(func(params GetNotagParams) middleware.Responder {`, api)

				assertInCode(t, `OverrideDeleteTestOverrideHandler override.DeleteTestOverrideHandler`, api)
				assertInCode(t, `StrfmtGetAnotherConflictHandler strfmtops.GetAnotherConflictHandler`, api)
				assertInCode(t, `APIGetConflictHandler apiops.GetConflictHandler`, api)
				assertInCode(t, `CustomGetCustomHandler custom.GetCustomHandler`, api)
				assertInCode(t, `AbcLinuxGetMultipleHandler abc_linux.GetMultipleHandler`, api)
				assertInCode(t, `GetNotagHandler GetNotagHandler`, api)
				assertInCode(t, `AbcLinuxGetOtherReservedHandler abc_linux.GetOtherReservedHandler`, api)
				assertInCode(t, `PlusDonutsGetOtherUnsafeHandler plus_donuts.GetOtherUnsafeHandler`, api)
				assertInCode(t, `AbcTestGetReservedHandler abc_test.GetReservedHandler`, api)
				assertInCode(t, `GtlGetTestOverrideHandler gtl.GetTestOverrideHandler`, api)
				assertInCode(t, `HashTagDonutsGetUnsafeHandler hash_tag_donuts.GetUnsafeHandler`, api)
				assertInCode(t, `NrAtDonutsGetYetAnotherUnsafeHandler nr_at_donuts.GetYetAnotherUnsafeHandler`, api)
				assertInCode(t, `ForcedPostTestOverrideHandler forced.PostTestOverrideHandler`, api)
				assertInCode(t, `Nr12nastyPutTestOverrideHandler nr12nasty.PutTestOverrideHandler`, api)
				assertInCode(t, `Nr123abcTestIDHandler nr123abc.TestIDHandler`, api)
			},
		},
		"packages_flattening": {
			spec: "../fixtures/bugs/2111/fixture-2111.yaml",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.SkipTagPackages = true
			},
			verify: func(t testing.TB, target string) {
				require.True(t, fileExists(target, defaultServerTarget))

				srvTarget := filepath.Join(target, defaultServerTarget)
				opsTarget := filepath.Join(srvTarget, defaultOperationsTarget)
				require.True(t, fileExists(opsTarget, ""))

				for _, fileOrDir := range []string{
					"abc_linux", "abc_test",
					"api",
					"custom",
					"hash_tag_donuts",
					"nr123abc", "nr_at_donuts", "plus_donuts",
					"strfmt",
					"forced",
					"gtl",
					"nr12nasty",
					"override",
					"operationsops",
				} {
					assert.Falsef(t, fileExists(opsTarget, fileOrDir), "did not expect %s in %s", fileOrDir, opsTarget)
				}

				assert.Truef(t, fileExists(opsTarget, "get_notag.go"), "expected %s in %s", "get_notag.go", opsTarget)

				buf, err := ioutil.ReadFile(filepath.Join(srvTarget, "configure_unsafe_tag_names.go"))
				require.NoError(t, err)
				code := string(buf)

				cwd := testCwd(t)
				base := path.Join("github.com", "go-swagger", "go-swagger",
					filepath.ToSlash(strings.TrimPrefix(target, filepath.Dir(cwd))),
				)

				baseImport := path.Join(base, `restapi/operations`)
				assertRegexpInCode(t, baseImport, code)

				assertInCode(t, `api.GetConflictHandler = operations.GetConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetAnotherConflictHandler = operations.GetAnotherConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetNotagHandler = operations.GetNotagHandlerFunc(`, code)

				buf2, err := ioutil.ReadFile(filepath.Join(opsTarget, "unsafe_tag_names_api.go"))
				require.NoError(t, err)
				api := string(buf2)

				assertInCode(t, `GetConflictHandler: GetConflictHandlerFunc(func(params GetConflictParams) middleware.Responder {`, api)
				assertInCode(t, `GetAnotherConflictHandler: GetAnotherConflictHandlerFunc(func(params GetAnotherConflictParams) middleware.Responder {`, api)
				assertInCode(t, `NotagHandler: GetNotagHandlerFunc(func(params GetNotagParams) middleware.Responder {`, api)

				assertInCode(t, `DeleteTestOverrideHandler DeleteTestOverrideHandler`, api)
				assertInCode(t, `GetAnotherConflictHandler GetAnotherConflictHandler`, api)
				assertInCode(t, `GetConflictHandler GetConflictHandler`, api)
				assertInCode(t, `GetCustomHandler GetCustomHandler`, api)
				assertInCode(t, `GetMultipleHandler GetMultipleHandler`, api)
				assertInCode(t, `GetNotagHandler GetNotagHandler`, api)
				assertInCode(t, `GetOtherReservedHandler GetOtherReservedHandler`, api)
				assertInCode(t, `GetOtherUnsafeHandler GetOtherUnsafeHandler`, api)
				assertInCode(t, `GetReservedHandler GetReservedHandler`, api)
				assertInCode(t, `GetTestOverrideHandler GetTestOverrideHandler`, api)
				assertInCode(t, `GetUnsafeHandler GetUnsafeHandler`, api)
				assertInCode(t, `GetYetAnotherUnsafeHandler GetYetAnotherUnsafeHandler`, api)
				assertInCode(t, `PostTestOverrideHandler PostTestOverrideHandler`, api)
				assertInCode(t, `PutTestOverrideHandler PutTestOverrideHandler`, api)
				assertInCode(t, `TestIDHandler TestIDHandler`, api)
			},
		},
		"main_package": {
			spec: "../fixtures/bugs/2111/fixture-2111.yaml",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.IncludeMain = true
				opts.MainPackage = "custom-api"
				opts.SkipTagPackages = true
			},
			verify: func(t testing.TB, target string) {
				assert.True(t, fileExists(filepath.Join(target, "cmd", "custom-api"), "main.go"))
			},
		},
		"external_model": {
			spec: "../fixtures/bugs/1897/fixture-1897.yaml",
			prepare: func(t testing.TB, opts *GenOpts) {
				modelOpts := *opts
				modelOpts.AcceptDefinitionsOnly = true
				modelOpts.Spec = "../fixtures/bugs/1897/model.yaml"
				modelOpts.ModelPackage = "external"
				modelOpts.Target = filepath.Dir(modelOpts.Spec)

				require.NoError(t, GenerateModels(nil, &modelOpts))

				t.Logf("generated external model")
				require.True(t, fileExists(modelOpts.Target, filepath.Join("external")))
				require.True(t, fileExists(modelOpts.Target, filepath.Join("external", "error.go")))

				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				location := filepath.Join(target, "cmd", "repro1897-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")
			},
			clean: func() {
				// remove generated external models
				_ = os.RemoveAll(filepath.Join("..", "fixtures", "bugs", "1897", "external"))
			},
		},
		"external_models_hints": {
			spec:   "../fixtures/enhancements/2224/fixture-2224.yaml",
			target: "2224-hints",
			prepare: func(t testing.TB, opts *GenOpts) {
				modelOpts := *opts
				modelOpts.AcceptDefinitionsOnly = true
				modelOpts.Spec = "../fixtures/enhancements/2224/fixture-2224-models.yaml"
				modelOpts.ModelPackage = "external"
				modelOpts.Target = filepath.Dir(modelOpts.Spec)

				require.NoError(t, GenerateModels(nil, &modelOpts))

				t.Logf("generated external model")
				require.True(t, fileExists(modelOpts.Target, filepath.Join("external")))

				for _, model := range []string{
					"access_point.go", "base.go",
					"hotspot.go", "hotspot_type.go",
					"incorrect.go", "json_message.go",
					"json_object.go", "json_object_with_alias.go",
					"object_with_embedded.go", "object_with_externals.go",
					"raw.go", "request.go",
					"request_pointer.go", "time_as_object.go", "time.go",
				} {
					require.True(t, fileExists(modelOpts.Target, filepath.Join("external", model)))
				}

				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				// generated models (not external)
				require.True(t, fileExists(target, filepath.Join("models")))
				for _, model := range []string{"error.go", "external_with_embed.go"} {
					require.True(t, fileExists(target, filepath.Join("models", model)))
				}

				location := filepath.Join(target, "cmd", "external-types-with-hints-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")
			},
			clean: func() {
				// remove generated external models
				_ = os.RemoveAll(filepath.Join("..", "fixtures", "enhancements", "2224", "external"))
			},
		},
		"conflict_name_api_issue_2405_1": {
			spec:   "../examples/todo-list/swagger.yml",
			target: "2405-1",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.ServerPackage = "api"
				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")
			},
		},
		"conflict_name_api_issue_2405_2": {
			spec:   "../examples/todo-list/swagger.yml",
			target: "2405-2",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.ServerPackage = "loads"
				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")
			},
		},
		"conflict_name_api_issue_2405_3": {
			spec:   "../fixtures/bugs/2405/fixture-2405.yaml",
			target: "2405-3",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.ServerPackage = "server"
				opts.APIPackage = "api"
				opts.IncludeMain = true
			},
			verify: func(t testing.TB, target string) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")
			},
		},
		"ext_types_issue_2385": {
			spec:   "../fixtures/bugs/2385/fixture-2385.yaml",
			target: "2385",
			prepare: func(t testing.TB, opts *GenOpts) {
				opts.MainPackage = "nrcodegen-server"
				opts.IncludeMain = true
				location := filepath.Join(opts.Target, "models")

				// add some custom model to the generated models
				addModelsToLocation(t, location, "my_type.go")
			},
			verify: func(_ testing.TB, target string) {
				location := filepath.Join(target, "cmd", "nrcodegen-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")

				location = filepath.Join(target, "models")

				t.Log("building generated models")
				goExecInDir(t, location, "build")
			},
		},
		"ext_types_full_example": {
			spec:   "../examples/external-types/example-external-types.yaml",
			target: "external-full",
			prepare: func(_ testing.TB, opts *GenOpts) {
				opts.MainPackage = "nrcodegen-server"
				opts.IncludeMain = true
				opts.ValidateSpec = false // the spec contains AdditionalItems
				location := filepath.Join(opts.Target, "models")

				// add some custom model to the generated models
				addModelsToLocation(t, location, "my_type.go")
			},
			verify: func(t testing.TB, target string) {
				location := filepath.Join(target, "cmd", "nrcodegen-server")
				require.True(t, fileExists("", location))

				t.Log("building generated server")
				goExecInDir(t, location, "build")

				location = filepath.Join(target, "models")

				t.Log("building generated models")
				goExecInDir(t, location, "build")
			},
		},
	}
}

func addModelsToLocation(t testing.TB, location, file string) {
	// writes some external model to a file to supplement codegen
	// (test external types)
	require.NoError(t, os.MkdirAll(location, 0700))

	require.NoError(t, ioutil.WriteFile(filepath.Join(location, file), []byte(`
package models

import (
  "context"
  "io"
  "github.com/go-openapi/strfmt"
)

// MyType ...
type MyType string

// Validate MyType
func (MyType) Validate(strfmt.Registry) error { return nil }
func (MyType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyInteger ...
type MyInteger int

// Validate MyInteger
func (MyInteger) Validate(strfmt.Registry) error { return nil }
func (MyInteger) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyString ...
type MyString string

// Validate MyString
func (MyString) Validate(strfmt.Registry) error { return nil }
func (MyString) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyOtherType ...
type MyOtherType struct{}

// Validate MyOtherType
func (MyOtherType) Validate(strfmt.Registry) error { return nil }
func (MyOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyStreamer ...
type MyStreamer io.Reader
`),
		os.ModePerm))
}
