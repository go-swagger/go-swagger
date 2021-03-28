package generate_test

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Make sure generated code compiles
func TestGenerateCLI(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

	testcases := []struct {
		name      string
		spec      string
		wantError bool
	}{
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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", tc.spec)
			generated, err := ioutil.TempDir(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				_ = os.RemoveAll(generated)
			}()
			m := &generate.Cli{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			err = m.Execute([]string{})
			if tc.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// change to true to run go vet on generated files
				runVet := false
				if runVet {
					vet := exec.Command("go", "vet", generated+"/...")
					output, err := vet.CombinedOutput()
					assert.NoError(t, err, string(output))
				}
			}
		})
	}
}

// This tests holds not yet working specs. It will be moved/combined with above test once fixed.
func TestGenerateCLIWorkInProgress(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

	testcases := []struct {
		name      string
		spec      string
		wantError bool
	}{
		// TODO: This one is not working for CLI
		// {
		// 	name:      "tasklist_basic",
		// 	spec:      "tasklist.basic.yml",
		// 	wantError: false,
		// },
		// TODO: Discriminator is not yet working for CLI.
		// {
		// 	name:      "todolist_discriminators",
		// 	spec:      "todolist.discriminators.yml",
		// 	wantError: false,
		// },
		// TODO: Not working yet. Need to find a way to detect pkg of model in schema
		// {
		// 	name:      "todolist_responses",
		// 	spec:      "todolist.responses.yml",
		// 	wantError: false,
		// },
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", tc.spec)
			generated, err := ioutil.TempDir(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			// defer func() {
			// 	_ = os.RemoveAll(generated)
			// }()
			m := &generate.Cli{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			err = m.Execute([]string{})
			if tc.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// run go vet on generated files
				vet := exec.Command("go", "vet", generated+"/...")
				output, err := vet.CombinedOutput()
				if err != nil {
					assert.NoError(t, err, string(output))
				}
			}
		})
	}
}
