// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	"github.com/go-swagger/go-swagger/generator/internal/gentest"
)

const (
	defaultAPIPackage    = "operations"
	defaultClientPackage = "client"
	defaultModelPackage  = "models"
)

func TestGenerateAndBuild(t *testing.T) {
	// This test generates and actually compiles the output
	// of generated clients.
	//
	// We run this in parallel now. Therefore it is no more
	// possible to assert the output on stdout.
	//
	// NOTE: test cases are randomized (map)
	t.Parallel()
	defer gentest.DiscardOutput()()

	cases := map[string]struct {
		spec string
	}{
		"issue 844": {
			"../fixtures/bugs/844/swagger.json",
		},
		"issue 844 (with params)": {
			"../fixtures/bugs/844/swagger-bis.json",
		},
		"issue 1216": {
			"../fixtures/bugs/1216/swagger.yml",
		},
		"issue 2111": {
			"../fixtures/bugs/2111/fixture-2111.yaml",
		},
		"issue 2278": {
			"../fixtures/bugs/2278/fixture-2278.yaml",
		},
		"issue 2163": {
			"../fixtures/enhancements/2163/fixture-2163.yaml",
		},
		"issue 1771": {
			"../fixtures/enhancements/1771/fixture-1771.yaml",
		},
	}

	t.Run("build client", func(t *testing.T) {
		// This test builds a client as a go module outside of the go source tree.
		//
		// It needs a go mod initialized and a go mod tidy sync, which slows down the test
		// a bit but simulates more realistically a full-fledged build with modules.
		for name, toPin := range cases {
			cas := toPin

			t.Run(name, func(t *testing.T) {
				t.Parallel()

				specPath := filepath.Clean(filepath.FromSlash(cas.spec))
				generatedLocation := filepath.Join(t.TempDir(), filepath.Base(specPath), "generated")
				t.Run(fmt.Sprintf("building client in %q", generatedLocation), func(t *testing.T) {
					require.NoError(t, os.MkdirAll(generatedLocation, fs.ModePerm))

					module := gentest.SanitizeGoModPath(generatedLocation)
					t.Run(fmt.Sprintf("should initialize module %q", module),
						gentest.GoExecInDir(generatedLocation, "mod", "init", module),
					)

					t.Run("should build client", func(t *testing.T) {
						require.NoError(t, newTestClient(specPath, generatedLocation).Execute(nil))
					})

					t.Run("should go get imports", gentest.GoExecInDir(generatedLocation, "get", "./..."))
					t.Run("should go mod tidy", gentest.GoExecInDir(generatedLocation, "mod", "tidy"))
					t.Run("should build client", gentest.GoExecInDir(generatedLocation, "build", "./..."))
				})
			})
		}
	})
}

func newTestClient(input, output string) *generate.Client {
	c := &generate.Client{}

	c.DefaultScheme = "http"
	c.DefaultProduces = "application/json"
	c.Shared.Spec = flags.Filename(input)
	c.Shared.Target = flags.Filename(output)
	c.Operations.APIPackage = defaultAPIPackage
	c.Models.ModelPackage = defaultModelPackage
	c.ClientPackage = defaultClientPackage

	return c
}
