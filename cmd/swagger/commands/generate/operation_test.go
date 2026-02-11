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

	"github.com/go-openapi/testify/v2/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestGenerateOperation(t *testing.T) {
	t.Run("should generate operation with simple Responder", testGenerateOperation(false))
}

func TestGenerateOperationStrict(t *testing.T) {
	t.Run("should generate operation with strict Responder", testGenerateOperation(true))
}

func testGenerateOperation(strict bool) func(*testing.T) {
	return func(t *testing.T) {
		specs := []string{
			"tasklist.basic.yml",
		}

		base := t.TempDir()

		for i, spec := range specs {
			t.Run(spec, func(t *testing.T) {
				pth := filepath.Join(testBase(), "fixtures/codegen", spec)
				generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
				require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

				m := &generate.Operation{}
				if i == 0 {
					m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
				}
				_, _ = flags.ParseArgs(m, []string{"--name=listTasks"})
				m.Shared.Spec = flags.Filename(pth)
				m.Shared.Target = flags.Filename(generated)
				m.Shared.StrictResponders = strict

				t.Run("go mod", gomodinit(generated))

				require.NoError(t, m.Execute([]string{}))
			})
		}
	}
}

func TestGenerateOperation_Check(t *testing.T) {
	m := &generate.Operation{}
	_, _ = flags.ParseArgs(m, []string{"--name=op1", "--name=op2"})
	m.Shared.DumpData = true
	m.Name = []string{"op1", "op2"}

	require.Error(t, m.Execute([]string{}))
}
