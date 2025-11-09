// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

type cliTestCase struct {
	name         string
	spec         string
	skip         bool
	wantError    bool
	wantVetError bool
}

// Make sure generated code compiles.
func TestGenerateCLI(t *testing.T) {
	testcases := []cliTestCase{
		{
			name:      "tasklist_basic",
			spec:      "tasklist.basic.yml",
			wantError: false,
		},
		{
			name:      "tasklist_allparams",
			spec:      "todolist.allparams.yml",
			wantError: false,
		},
		{
			name:      "tasklist_arrayform",
			spec:      "todolist.arrayform.yml",
			wantError: false,
		},
		{
			name:      "tasklist_arrayquery",
			spec:      "todolist.arrayquery.yml",
			wantError: false,
		},
		{
			name:      "todolist_bodyparams",
			spec:      "todolist.bodyparams.yml",
			wantError: false,
		},
		{
			name:      "tasklist_simplequery",
			spec:      "todolist.simplequery.yml",
			wantError: false,
		},
		{
			name:      "todolist_responses",
			spec:      "todolist.responses.yml",
			wantError: false,
		},
		{
			name:      "todo_simple-fixed",
			spec:      "todolist.simple-fixed.yml",
			wantError: false,
		},
		{
			name:      "todo_simpleform",
			spec:      "todolist.simpleform.yml",
			wantError: false,
		},
		{
			name:      "todo_simpleheader",
			spec:      "todolist.simpleheader.yml",
			wantError: false,
		},
		{
			name:      "todo_simplepath",
			spec:      "todolist.simplepath.yml",
			wantError: false,
		},
	}

	base := t.TempDir()

	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", tc.spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

			m := &generate.Cli{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)

			t.Run("go mod", gomodinit(generated))

			err := m.Execute([]string{})
			if tc.wantError {
				require.Error(t, err)

				return
			}
			require.NoError(t, err)

			// run test with GOVET_TEST=1 go test -v -run CLI to run go vet on generated files.
			//
			// There are known issues there at this moment
			if enableGoVet {
				t.Run("go mod tidy", gomodtidy(generated))
				t.Run("go vet", govet(generated, tc.wantVetError))
			}
		})
	}
}

func TestGenerateCli_Check(t *testing.T) {
	m := &generate.Cli{}
	_, _ = flags.Parse(m)

	t.Run("generate CLI requires an argument", func(t *testing.T) {
		require.Error(t, m.Execute([]string{}))
	})
}

// This test runs cli generation on various swagger specs, for sanity check.
// Skipped in by default. Only run by developer locally.
func TestVariousCli(t *testing.T) {
	// comment out this skip to run test
	t.Skip()

	testcases := []cliTestCase{
		{
			skip:         true, // do not run this in CI since it is known to have bug
			name:         "crazy-alias",
			spec:         "fixtures/bugs/1260/fixture-realiased-types.yaml",
			wantError:    false, // generate files should success
			wantVetError: true,  // polymorphism is not supported. model import is not right. TODO: fix this.
		},
		{
			name: "multi-auth",
			spec: "examples/composed-auth/swagger.yml",
		},
		// not working because of model generation order.
		// {
		// 	name:          "enum",
		// 	spec:          "fixtures/enhancements/1623/swagger.yml",
		// },
	}

	base := t.TempDir()

	for i, tc := range testcases {
		t.Run(tc.name, func(tt *testing.T) {
			if tc.skip {
				tt.Skip()
			}

			pth := filepath.Join(testBase(), tc.spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

			m := &generate.Cli{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)

			t.Run("go mod", gomodinit(generated))

			err := m.Execute([]string{})
			if tc.wantError {
				require.Error(tt, err)

				return
			}
			require.NoError(tt, err)

			// run go vet on generated files
			if enableGoVet {
				t.Run("go mod tidy", gomodtidy(generated))
				t.Run("go vet", govet(generated, tc.wantVetError))
			}
		})
	}
}
