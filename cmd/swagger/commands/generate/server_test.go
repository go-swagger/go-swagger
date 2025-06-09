package generate_test

import (
	"fmt"
	"path/filepath"
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

	for i, spec := range specs {
		t.Run(fmt.Sprintf("should generate server from spec %s", spec), func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)

			m := &generate.Server{}
			_, _ = flags.Parse(m)
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)

			// Error was coming from these two being set together
			m.Shared.StrictResponders = true
			m.ImplementationPackage = "github.com/go-swagger/go-swagger/fixtures/codegen/impl"

			// Load new copy of template
			m.Shared.AllowTemplateOverride = true
			m.Shared.TemplateDir = flags.Filename(filepath.Join(testBase(), "generator/templates"))

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
		t.Run(fmt.Sprintf("should generate server from spec %s", spec), func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)

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
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			require.NoError(t, m.Execute([]string{}))
		})
	}
}
