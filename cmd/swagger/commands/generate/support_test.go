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

func TestGenerateSupport(t *testing.T) {
	testGenerateSupport(t, false)
}

func TestGenerateSupportStrict(t *testing.T) {
	testGenerateSupport(t, true)
}

func testGenerateSupport(t *testing.T, strict bool) {
	specs := []string{
		"tasklist.basic.yml",
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
			m := &generate.Support{}
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
			}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}
