package generate_test

import (
	"io"
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
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

	testcases := []struct {
		name      string
		spec      string
		wantError bool
	}{
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

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", tc.spec)
			generated, err := os.MkdirTemp(filepath.Dir(path), "generated")
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

func TestGenerateCli_Check(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Cli{}
	_, _ = flags.Parse(m)
	err := m.Execute([]string{})
	assert.Error(t, err)
}

// This test runs cli generation on various swagger specs, for sanity check.
// Skipped in by default. Only run by developer locally.
func TestVariousCli(t *testing.T) {
	// comment out this skip to run test
	t.Skip()

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

	// change to true to run test case with runOnly set true
	runOnlyTest := false

	testcases := []struct {
		skip          bool
		name          string
		spec          string
		wantError     bool
		wantVetError  bool
		preserveFiles bool // force to preserve files
		runOnly       bool // run only this test, and skip all others
	}{
		{
			skip:         true, // do not run this since it is known to have bug
			name:         "crazy-alias",
			spec:         "fixtures/bugs/1260/fixture-realiased-types.yaml",
			wantError:    false, // generate files should success
			wantVetError: true,  // polymorphism is not supported. model import is not right. TODO: fix this.
		},
		{
			name:          "multi-auth",
			spec:          "examples/composed-auth/swagger.yml",
			preserveFiles: true,
		},
		// not working because of model generation order.
		// {
		// 	name:          "enum",
		// 	spec:          "fixtures/enhancements/1623/swagger.yml",
		// 	preserveFiles: true,
		// 	runOnly:       true,
		// },
	}

	for _, tc := range testcases {
		if runOnlyTest && !tc.runOnly {
			continue
		}
		t.Run(tc.name, func(tt *testing.T) {
			if tc.skip {
				tt.Skip()
			}
			path := filepath.Join(base, tc.spec)
			generated, err := os.MkdirTemp(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				// only clean up if success, and leave the files around for developer to inspect
				if !tt.Failed() {
					if !tc.preserveFiles {
						_ = os.RemoveAll(generated)
					}
				} else {
					// stop all tests, since it will generate too many files to inspect
					t.FailNow()
				}
			}()
			m := &generate.Cli{}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			err = m.Execute([]string{})
			if tc.wantError {
				assert.Error(tt, err)
			} else {
				require.NoError(tt, err)
				// always run go vet on generated files
				runVet := true
				if runVet {
					vet := exec.Command("go", "vet", generated+"/...")
					output, err := vet.CombinedOutput()
					if !tc.wantVetError {
						assert.NoError(tt, err, string(output))
					} else {
						assert.Error(t, err)
					}
				}
			}
		})
	}
}
