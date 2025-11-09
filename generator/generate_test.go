// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"
)

const testServerPkg = "nrcodegen-server"

func TestGenerateAndTest(t *testing.T) {
	// Full codegen test, with advanced assertions, which may be:
	// * simple assertions about generated folders or code
	// * inclusion of a test program to exercise the generated program
	// * inclusion of test dependencies to exercise the build
	//
	// All generated code is built as go modules in t.TempDir().
	t.Parallel()
	defer discardOutput()()

	root := t.TempDir()
	t.Run("server build", testHarnessBuildServers(root, generateServerFixtures())) // debugging mode: add selected fixture keys to limit the scope
	t.Run("client build", testHarnessBuildClients(root, generateClientFixtures())) // debugging mode: add selected fixtures keys to limit the scope
}

func generateServerFixtures() map[string]generateFixture {
	return map[string]generateFixture{
		"issue 1943":                       fixtureServer1943(),
		"packages_mangling":                fixtureServerPackageMangling(),
		"packages_flattening":              fixtureServerPackageFlattening(),
		"main_package":                     fixtureServerMainPackage(),
		"external_model":                   fixtureServerExternalModel(),
		"external_models_hints":            fixtureServerExternalModelsHints(),
		"conflict_name_api_issue_2405_1":   fixtureServerNameConflict2405_1(),
		"conflict_name_api_issue_2405_2":   fixtureServerNameConflict2405_2(),
		"conflict_name_api_issue_2405_3":   fixtureServerNameConflict2405_3(),
		"ext_types_issue_2385":             fixtureServerExternalTypes2385(),
		"ext_types_full_example":           fixtureServerExternalTypesFull(),
		"conflict_name_server_issue_2730":  fixtureServerNameConflictServer2730(),
		"tag_package_name_issue_2866":      fixtureServerTagPackageName2866(),
		"tag_package_name_regression_3143": fixtureServerTagPackageName3143(),
	}
}

func generateClientFixtures() map[string]generateFixture {
	return map[string]generateFixture{
		"issue1083":                       fixtureClientRoundTrip1083(),
		"conflict_name_client_issue_2730": fixtureClientNameConflict2730(),
		"type conversions":                fixtureClientTypeConversions(),
	}
}

func fixtureServer1943() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/1943/fixture-1943.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				input, err := os.ReadFile("../fixtures/bugs/1943/datarace_test.go")
				require.NoError(t, err)

				// rewrite imports for the relocated test program
				rebasedContent := bytes.ReplaceAll(
					input,
					[]byte("github.com/go-swagger/go-swagger/fixtures/bugs/1943"),
					[]byte(gentest.SanitizeGoModPath(opts.Target)),
				)

				rebasedContent = removeBuildTags(rebasedContent)
				require.NoError(t, os.WriteFile(filepath.Join(opts.Target, "datarace_test.go"), rebasedContent, 0o600))
				opts.ExcludeSpec = false
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				if runtime.GOOS == winOS {
					// don't run race tests on windows (why so?)
					t.Skipf("warn: race test skipped on os %s", runtime.GOOS)

					return
				}

				testPrg := "datarace_test.go"
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("go get", gentest.GoExecInDir(target, "get", "./..."))
				t.Run("running data race test on generated server",
					gentest.GoExecInDir(target, "test", "-v", "-race", testPrg),
				)
			}
		},
	}
}

func fixtureServerPackageMangling() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2111/fixture-2111.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				require.DirExists(t, filepath.Join(target, defaultServerTarget))
				assert.FileExists(t, filepath.Join(target, "cmd", "unsafe-tag-names-server", "main.go"))

				srvTarget := filepath.Join(target, defaultServerTarget)
				opsTarget := filepath.Join(srvTarget, defaultOperationsTarget)
				require.DirExists(t, srvTarget)
				require.DirExists(t, opsTarget)

				for _, fileOrDir := range []string{
					"abc_linux", "abc_test",
					apiPkg,
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
					if strings.HasSuffix(fileOrDir, ".go") {
						assert.FileExists(t, filepath.Join(opsTarget, fileOrDir))

						continue
					}
					assert.DirExists(t, filepath.Join(opsTarget, fileOrDir))
				}

				buf, err := os.ReadFile(filepath.Join(srvTarget, "configure_unsafe_tag_names.go"))
				require.NoError(t, err)

				code := string(buf)

				// assert imports, with deconfliction
				base := gentest.SanitizeGoModPath(target)
				baseImport := path.Join(base, `restapi/operations`)
				assertImports(t, baseImport, code)

				assertInCode(t, `api.APIGetConflictHandler = apiops.GetConflictHandlerFunc(`, code)
				assertInCode(t, `api.StrfmtGetAnotherConflictHandler = strfmtops.GetAnotherConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetNotagHandler = operations.GetNotagHandlerFunc(`, code)

				buf2, err := os.ReadFile(filepath.Join(opsTarget, "unsafe_tag_names_api.go"))
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
				assertInCode(t, `GetNotagHandler`, api)
				assertInCode(t, `AbcLinuxGetOtherReservedHandler abc_linux.GetOtherReservedHandler`, api)
				assertInCode(t, `PlusDonutsGetOtherUnsafeHandler plus_donuts.GetOtherUnsafeHandler`, api)
				assertInCode(t, `AbcTestGetReservedHandler abc_test.GetReservedHandler`, api)
				assertInCode(t, `GtlGetTestOverrideHandler gtl.GetTestOverrideHandler`, api)
				assertInCode(t, `HashTagDonutsGetUnsafeHandler hash_tag_donuts.GetUnsafeHandler`, api)
				assertInCode(t, `NrAtDonutsGetYetAnotherUnsafeHandler nr_at_donuts.GetYetAnotherUnsafeHandler`, api)
				assertInCode(t, `ForcedPostTestOverrideHandler forced.PostTestOverrideHandler`, api)
				assertInCode(t, `Nr12nastyPutTestOverrideHandler nr12nasty.PutTestOverrideHandler`, api)
				assertInCode(t, `Nr123abcTestIDHandler nr123abc.TestIDHandler`, api)
			}
		},
	}
}

func fixtureServerPackageFlattening() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2111/fixture-2111.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.SkipTagPackages = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				require.DirExists(t, filepath.Join(target, defaultServerTarget))

				srvTarget := filepath.Join(target, defaultServerTarget)
				opsTarget := filepath.Join(srvTarget, defaultOperationsTarget)
				require.DirExists(t, srvTarget)
				require.DirExists(t, opsTarget)

				for _, fileOrDir := range []string{
					"abc_linux", "abc_test",
					apiPkg,
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

				buf, err := os.ReadFile(filepath.Join(srvTarget, "configure_unsafe_tag_names.go"))
				require.NoError(t, err)
				code := string(buf)

				base := gentest.SanitizeGoModPath(target)
				baseImport := path.Join(base, `restapi/operations`)
				assertRegexpInCode(t, baseImport, code)

				assertInCode(t, `api.GetConflictHandler = operations.GetConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetAnotherConflictHandler = operations.GetAnotherConflictHandlerFunc(`, code)
				assertInCode(t, `api.GetNotagHandler = operations.GetNotagHandlerFunc(`, code)

				buf2, err := os.ReadFile(filepath.Join(opsTarget, "unsafe_tag_names_api.go"))
				require.NoError(t, err)
				api := string(buf2)

				assertInCode(t, `GetConflictHandler: GetConflictHandlerFunc(func(params GetConflictParams) middleware.Responder {`, api)
				assertInCode(t, `GetAnotherConflictHandler: GetAnotherConflictHandlerFunc(func(params GetAnotherConflictParams) middleware.Responder {`, api)
				assertInCode(t, `NotagHandler: GetNotagHandlerFunc(func(params GetNotagParams) middleware.Responder {`, api)

				assertInCode(t, `DeleteTestOverrideHandler`, api)
				assertInCode(t, `GetAnotherConflictHandler`, api)
				assertInCode(t, `GetConflictHandler`, api)
				assertInCode(t, `GetCustomHandler`, api)
				assertInCode(t, `GetMultipleHandler`, api)
				assertInCode(t, `GetNotagHandler`, api)
				assertInCode(t, `GetOtherReservedHandler`, api)
				assertInCode(t, `GetOtherUnsafeHandler`, api)
				assertInCode(t, `GetReservedHandler`, api)
				assertInCode(t, `GetTestOverrideHandler`, api)
				assertInCode(t, `GetUnsafeHandler`, api)
				assertInCode(t, `GetYetAnotherUnsafeHandler`, api)
				assertInCode(t, `PostTestOverrideHandler`, api)
				assertInCode(t, `PutTestOverrideHandler`, api)
				assertInCode(t, `TestIDHandler`, api)
			}
		},
	}
}

func fixtureServerMainPackage() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2111/fixture-2111.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.IncludeMain = true
				opts.MainPackage = "custom-api"
				opts.SkipTagPackages = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				assert.FileExists(t, filepath.Join(target, "cmd", "custom-api", "main.go"))
			}
		},
	}
}

func fixtureServerExternalModel() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/1897/fixture-1897.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				// generate a module for external models in {test dir}/external
				t.Run("should generate external model", generateExternalModel(
					opts,
					filepath.Join("..", "fixtures", "bugs", "1897", "model.yaml"),  // the spec for the external model
					"github.com/go-swagger/go-swagger/fixtures/bugs/1897/external", // the external package in imports
				))

				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			// verify that all dependencies are found and that a complete server can build
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", "repro1897-server")
				require.DirExists(t, location)

				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerExternalModelsHints() generateFixture {
	return generateFixture{
		spec: "../fixtures/enhancements/2224/fixture-2224.yaml",
		// in this test case, we have a mix of generated models and external models
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				// generate a module for external models in {test dir}/external
				t.Run("should generate external model", generateExternalModel(
					opts,
					filepath.Join("..", "fixtures", "enhancements", "2224", "fixture-2224-models.yaml"), // the spec for the external model
					"github.com/go-swagger/go-swagger/fixtures/enhancements/2224/external",              // the external package in imports
				))

				t.Run("external models should be available", func(t *testing.T) {
					require.DirExists(t, filepath.Join(opts.Target, "external"))

					for _, model := range []string{
						"access_point.go", "base.go",
						"hotspot.go", "hotspot_type.go",
						"incorrect.go", "json_message.go",
						"json_object.go", "json_object_with_alias.go",
						"object_with_embedded.go", "object_with_externals.go",
						"raw.go", "request.go",
						"request_pointer.go", "time_as_object.go", "time.go",
					} {
						require.FileExists(t, filepath.Join(opts.Target, "external", model))
					}
				})

				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				// generated models (not external)
				require.DirExists(t, filepath.Join(target, "models"))
				for _, model := range []string{"error.go", "external_with_embed.go"} {
					require.FileExists(t, filepath.Join(target, "models", model))
				}

				location := filepath.Join(target, "cmd", "external-types-with-hints-server")
				require.DirExists(t, location)
				// verify that all dependencies are found and that a complete server can build
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerNameConflict2405_1() generateFixture {
	return generateFixture{
		spec: "../examples/todo-list/swagger.yml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.ServerPackage = apiPkg
				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerNameConflict2405_2() generateFixture {
	return generateFixture{
		spec: "../examples/todo-list/swagger.yml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.ServerPackage = "loads"
				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerNameConflict2405_3() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2405/fixture-2405.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.ServerPackage = "server"
				opts.APIPackage = apiPkg
				opts.IncludeMain = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", "simple-to-do-list-api-server")
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerExternalTypes2385() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2385/fixture-2385.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				opts.MainPackage = testServerPkg
				opts.IncludeMain = true
				location := filepath.Join(opts.Target, "models")

				// add some custom model to the generated models
				addModelsToLocation(t, location, "my_type.go")
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", testServerPkg)
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))

				location = filepath.Join(target, "models")
				require.DirExists(t, location)
				t.Run("building generated models", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerExternalTypesFull() generateFixture {
	return generateFixture{
		spec: "../examples/external-types/example-external-types.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				opts.MainPackage = testServerPkg
				opts.IncludeMain = true
				opts.ValidateSpec = false // the spec contains AdditionalItems

				// add some custom model to the generated models
				location := filepath.Join(opts.Target, "models")
				addModelsToLocation(t, location, "my_type.go")
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", testServerPkg)
				require.DirExists(t, location)

				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
				location = filepath.Join(target, "models")
				t.Run("building generated models", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerNameConflictServer2730() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2730/2730.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.MainPackage = testServerPkg
				opts.IncludeMain = true
				opts.ValidateSpec = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", testServerPkg)
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerTagPackageName2866() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/2866/2866.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.MainPackage = testServerPkg
				opts.IncludeMain = true
				opts.ValidateSpec = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", testServerPkg)
				require.DirExists(t, location)
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "version1"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "version3"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "v2_validations"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "v3_validations"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "v3_actual"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "v3_planned"))

				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureServerTagPackageName3143() generateFixture {
	return generateFixture{
		spec: "../fixtures/bugs/3143/3143.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(_ *testing.T) {
				opts.MainPackage = testServerPkg
				opts.IncludeMain = true
				opts.ValidateSpec = true
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "cmd", testServerPkg)
				require.DirExists(t, location)
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "av2on"))
				require.DirExists(t, filepath.Join(target, "restapi", "operations", "trailingv2"))
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated server", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureClientRoundTrip1083() generateFixture {
	return generateFixture{
		// exercise generated client + untyped server
		spec: "../fixtures/bugs/1083/petstore.yaml",
		prepare: func(opts *GenOpts) func(*testing.T) {
			return func(t *testing.T) {
				targetImport := gentest.SanitizeGoModPath(opts.Target) // the generated module

				t.Run("should relocate test program", func(t *testing.T) {
					input, err := os.ReadFile(filepath.Join("..", "fixtures", "bugs", "1083", "pathparam_test.go"))
					require.NoError(t, err)

					// rewrite imports and relocates test program to the codegen target directory.
					//
					// Imports are rewritten such that there is no need for a replace directive in the generated go.mod
					rebasedContent := bytes.ReplaceAll(
						input,
						[]byte("github.com/go-swagger/go-swagger/fixtures/bugs/1083/codegen"),
						[]byte(targetImport),
					)
					rebasedContent = removeBuildTags(rebasedContent)
					require.NoError(t, os.WriteFile(filepath.Join(opts.Target, "pathparam_test.go"), rebasedContent, 0o600))
				})

				opts.ExcludeSpec = false

				t.Run("should copy spec for untyped usage", func(t *testing.T) {
					f, err := os.Open(filepath.Join("..", "fixtures", "bugs", "1083", "petstore.yaml"))
					require.NoError(t, err)
					defer func() {
						_ = f.Close()
					}()

					w, err := os.Create(filepath.Join(opts.Target, "petstore.yaml"))
					require.NoError(t, err)
					defer func() {
						_ = w.Close()
					}()
					_, err = io.Copy(w, f)
					require.NoError(t, err)
				})
			}
		},
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				testPrg := "pathparam_test.go"

				t.Run("go get", gentest.GoExecInDir(target, "get", "./..."))
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("running runtime request test on generated client",
					// This test runs a generated client against an untyped API server.
					// It verifies that path parameters are properly escaped and unescaped.
					// It exercises the full stack of runtime client and server.
					gentest.GoExecInDir(target, "test", "-v", testPrg),
				)
			}
		},
	}
}

func fixtureClientNameConflict2730() generateFixture {
	return generateFixture{
		spec:    "../fixtures/bugs/2730/2730.yaml",
		prepare: nil,
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "client")
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated client", gentest.GoBuild(location))
			}
		},
	}
}

func fixtureClientTypeConversions() generateFixture {
	return generateFixture{
		spec:    "../fixtures/codegen/conversions.yaml",
		prepare: nil,
		verify: func(target string) func(*testing.T) {
			return func(t *testing.T) {
				location := filepath.Join(target, "client")
				require.DirExists(t, location)
				t.Run("should tidy go mod", gentest.GoModTidy(target))
				t.Run("building generated client", gentest.GoBuild(location))
			}
		},
	}
}

func addModelsToLocation(t *testing.T, location, file string) {
	// writes some external model to a file to supplement codegen
	// (test external types)
	t.Helper()

	require.NoError(t, os.MkdirAll(location, 0o700))

	require.NoError(t, os.WriteFile(filepath.Join(location, file), []byte(`
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

// ContextValidate ...
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

// ContextValidate ...
func (MyOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyStreamer ...
type MyStreamer io.Reader
`),
		readableFile))
}

// generateExternalModel prepares some independently generated models, before we generate a server
// or client which imports them.
//
// Models are generated as their own go module, in folder "external" located in the temporary test directory.
// The module must be named exactly as it is imported (e.g. declared import location in the server spec) and
// a "replace" clause is added to the consuming module to locate the newly generated models.
func generateExternalModel(opts *GenOpts, modelSpecPath string, modelsPackage string) func(t *testing.T) {
	return func(t *testing.T) {
		// we first generate an external model from model.yaml, as its own module
		modelOpts := *opts
		modelOpts.AcceptDefinitionsOnly = true
		// the location of the spec for external models
		modelOpts.Spec = modelSpecPath
		modelOpts.ModelPackage = "external"
		targetPackageLocation := filepath.Join(modelOpts.Target, modelOpts.ModelPackage)
		require.NoError(t, os.MkdirAll(targetPackageLocation, readableDir))

		// generate module "external" with its fully qualified name referenced by imports
		t.Run("models mod init", gentest.GoModInit(targetPackageLocation, gentest.WithGoModuleName(modelsPackage)))

		// generate the external models package in package "external/models"
		require.NoError(t, GenerateModels(nil, &modelOpts))

		t.Run("should replace external package by test module",
			// in the module of the target server, replace the reference to the external module by its actual generated location
			gentest.GoExecInDir(
				opts.Target,
				"mod", "edit",
				"-replace", fmt.Sprintf("%s=%s", modelsPackage, targetPackageLocation),
			),
		)
	}
}
