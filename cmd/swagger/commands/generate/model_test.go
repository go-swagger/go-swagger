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

	base := t.TempDir()

	for i, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

			m := &generate.Model{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Models.ExistingModels = "nonExisting"
			}
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)

			t.Run("go mod", gomodinit(generated))

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
