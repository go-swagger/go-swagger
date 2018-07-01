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
	specs := []string{
		"tasklist.basic.yml",
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
			m := &generate.Client{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
			}
			m.Spec = flags.Filename(path)
			m.Target = flags.Filename(generated)

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateClient_Check(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Client{}
	_, _ = flags.Parse(m)
	m.CopyrightFile = "nullePart"
	err := m.Execute([]string{})
	assert.Error(t, err)
}
