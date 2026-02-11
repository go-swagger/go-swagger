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

func TestGenerateClient(t *testing.T) {
	tests := []struct {
		name      string
		spec      string
		template  string
		wantError bool
		prepare   func(c *generate.Client)
	}{
		{
			name:      "tasklist_basic",
			spec:      "tasklist.basic.yml",
			wantError: false,
		},
		{
			name:      "tasklist_simplequery",
			spec:      "todolist.simplequery.yml",
			wantError: false,
			prepare: func(c *generate.Client) {
				c.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			},
		},
		{
			name:      "generate_client_with_invalid_template",
			spec:      "todolist.simplequery.yml",
			template:  "NonExistingContributorTemplate",
			wantError: true,
		},
		{
			name:      "Existing_contributor",
			spec:      "todolist.simplequery.yml",
			template:  "stratoscale",
			wantError: false,
		},
	}

	// calling TempDir() within the loop exposes the test to unwanted waits or sometimes failures on windows (e.g. sub-test locking file).
	// Best to trigger the cleanup only once.
	base := t.TempDir()

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", tt.spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))
			m := &generate.Client{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.Template = tt.template

			if tt.prepare != nil {
				tt.prepare(m)
			}
			t.Run("go mod", gomodinit(generated))

			err := m.Execute([]string{})
			if tt.wantError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGenerateClient_Check(t *testing.T) {
	m := &generate.Client{}
	_, _ = flags.Parse(m)
	m.Shared.CopyrightFile = "nullePart"

	t.Run("should error when the copyright file does not exist", func(t *testing.T) {
		require.Error(t, m.Execute([]string{}))
	})
}
