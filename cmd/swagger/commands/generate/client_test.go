package generate_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
)

func TestGenerateClient(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

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
				c.Shared.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
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

			path := filepath.Join(base, "fixtures/codegen", tt.spec)
			generated, err := ioutil.TempDir(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				_ = os.RemoveAll(generated)
			}()
			m := &generate.Client{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.Template = tt.template

			if tt.prepare != nil {
				tt.prepare(m)
			}

			err = m.Execute([]string{})
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateClient_Check(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Client{}
	_, _ = flags.Parse(m)
	m.Shared.CopyrightFile = "nullePart"
	err := m.Execute([]string{})
	assert.Error(t, err)
}
