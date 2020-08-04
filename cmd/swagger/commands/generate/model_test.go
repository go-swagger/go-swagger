package generate_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
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
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")
	for i, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", spec)
			generated, err := ioutil.TempDir(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				_ = os.RemoveAll(generated)
			}()
			m := &generate.Model{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Models.ExistingModels = "nonExisting"
			}
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateModel_Check(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Model{}
	_, _ = flags.Parse(m)
	m.Shared.DumpData = true
	m.Name = []string{"model1", "model2"}
	err := m.Execute([]string{})
	assert.Error(t, err)
}
