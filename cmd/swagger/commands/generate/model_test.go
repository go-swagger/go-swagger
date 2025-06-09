package generate_test

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestGenerateModel(t *testing.T) {
	specs := []string{
		"billforward.discriminators.yml",
		"existing-model.yml",
		"instagram.yml",
		"shipyard.yml",
		"sodabooth.json",
		"tasklist.basic.yml",
		"todolist.simpleform.yml",
		"todolist.simpleheader.yml",
		"todolist.simplequery.yml",
	}

	for i, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)

			m := &generate.Model{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Models.ExistingModels = "nonExisting"
			}
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			require.NoError(t, m.Execute([]string{}))
		})
	}
}

func TestGenerateModel_Check(t *testing.T) {
	m := &generate.Model{}
	_, _ = flags.Parse(m)
	m.Shared.DumpData = true
	m.Name = []string{"model1", "model2"}
	require.Error(t, m.Execute([]string{}))
}
