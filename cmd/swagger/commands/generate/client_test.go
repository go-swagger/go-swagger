package generate_test

import (
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", tt.spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)
			m := &generate.Client{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.Template = tt.template

			if tt.prepare != nil {
				tt.prepare(m)
			}

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
	require.Error(t, m.Execute([]string{}))
}

func testBase() string {
	return filepath.FromSlash("../../../../")
}

func testTempDir(t testing.TB, path string) (generated string, cleanup func()) {
	t.Helper()
	var err error

	generated, err = os.MkdirTemp(filepath.Dir(path), "generated")
	require.NoErrorf(t, err, "TempDir()=%s", generated)

	return generated, func() {
		_ = os.RemoveAll(generated)
	}
}
