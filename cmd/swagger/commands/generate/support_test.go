package generate_test

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
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

	for i, spec := range specs {
		t.Run(spec, func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)

			m := &generate.Support{}
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			require.NoError(t, m.Execute([]string{}))
		})
	}
}
