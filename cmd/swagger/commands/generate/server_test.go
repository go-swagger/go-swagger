package generate_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	"github.com/jessevdk/go-flags"
)

func TestGenerateServer(t *testing.T) {
	testGenerateServer(t, false)
}

func TestGenerateServerStrict(t *testing.T) {
	testGenerateServer(t, true)
}

func testGenerateServer(t *testing.T, strict bool) {
	specs := []string{
		"billforward.discriminators.yml",
		"todolist.simplequery.yml",
		"todolist.simplequery.yml",
	}
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")
	for i, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", spec)
			generated, err := os.MkdirTemp(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				_ = os.RemoveAll(generated)
			}()
			m := &generate.Server{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
			}
			switch i {
			case 1:
				m.FlagStrategy = "pflag"
			case 2:
				m.FlagStrategy = "flag"
			}
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateServer_Checks(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Server{}
	_, _ = flags.Parse(m)
	m.Shared.CopyrightFile = "nowhere"
	err := m.Execute([]string{})
	assert.Error(t, err)
}

func TestRegressionIssue2601(t *testing.T) {
	specs := []string{
		"impl.yml",
	}
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")
	for i, spec := range specs {
		_ = t.Run(spec, func(t *testing.T) {
			path := filepath.Join(base, "fixtures/codegen", spec)
			generated, err := os.MkdirTemp(filepath.Dir(path), "generated")
			if err != nil {
				t.Fatalf("TempDir()=%s", generated)
			}
			defer func() {
				_ = os.RemoveAll(generated)
			}()
			m := &generate.Server{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
			}
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			// Error was coming from these two being set together
			m.Shared.StrictResponders = true
			m.ImplementationPackage = "github.com/go-swagger/go-swagger/fixtures/codegen/impl"

			// Load new copy of template
			m.Shared.AllowTemplateOverride = true
			m.Shared.TemplateDir = flags.Filename(filepath.Join(base, "generator/templates"))

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}
