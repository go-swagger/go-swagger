package generate_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestGenerateServer(t *testing.T) {
	testGenerateServer(t, false)
}

func TestGenerateServerStrict(t *testing.T) {
	testGenerateServer(t, true)
}

func TestGenerateServer_Checks(t *testing.T) {
	t.Run("invalid provided copyright file should error", func(t *testing.T) {
		m := &generate.Server{}
		_, _ = flags.Parse(m)
		m.Shared.CopyrightFile = "nowhere"
		require.Error(t, m.Execute([]string{}))
	})
}

func TestRegressionIssue2601(t *testing.T) {
	specs := []string{
		"impl.yml",
	}

	base := t.TempDir()

	for i, spec := range specs {
		t.Run("should generate server from spec "+spec, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

			m := &generate.Server{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)

			// Error was coming from these two being set together
			m.Shared.StrictResponders = true
			m.ImplementationPackage = "github.com/go-swagger/go-swagger/fixtures/codegen/impl"

			// Load new copy of template
			m.Shared.AllowTemplateOverride = true
			m.Shared.TemplateDir = flags.Filename(filepath.Join(testBase(), "generator/templates"))

			t.Run("go mod", gomodinit(generated))

			require.NoError(t, m.Execute([]string{}))
		})
	}
}

func testGenerateServer(t *testing.T, strict bool) {
	specs := []string{
		"billforward.discriminators.yml",
		"todolist.simplequery.yml",
		"todolist.simplequery.yml",
	}

	for i, spec := range specs {
		t.Run("should generate server from spec "+spec, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated := t.TempDir()

			m := &generate.Server{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			switch i {
			case 1:
				m.FlagStrategy = "pflag"
			case 2:
				m.FlagStrategy = "flag"
			}
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			t.Run("go mod", gomodinit(generated))

			require.NoError(t, m.Execute([]string{}))
		})
	}
}
