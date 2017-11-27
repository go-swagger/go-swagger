package generate_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
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
	for _, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", spec)
			generated, err := ioutil.TempDir(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer os.RemoveAll(generated)
			m := &generate.Model{}
			flags.Parse(m)
			m.Spec = flags.Filename(path)
			m.Target = flags.Filename(generated)
			m.NoValidator = true

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}
