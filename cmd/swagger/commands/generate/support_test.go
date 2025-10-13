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

func TestGenerateSupport(t *testing.T) {
	testGenerateSupport(t, false)
}

func TestGenerateSupportWithStrictResponders(t *testing.T) {
	testGenerateSupport(t, true)
}

func testGenerateSupport(t *testing.T, strict bool) {
	specs := []string{
		"tasklist.basic.yml",
	}

	base := t.TempDir()

	for i, spec := range specs {
		t.Run(spec, func(t *testing.T) {
			pth := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated := filepath.Join(base, "codegen-"+strconv.Itoa(i))
			require.NoError(t, os.MkdirAll(generated, fs.ModePerm))

			m := &generate.Support{}
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			_, _ = flags.Parse(m)
			m.Shared.Spec = flags.Filename(pth)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			t.Run("go mod", gomodinit(generated))

			require.NoError(t, m.Execute([]string{}))
		})
	}
}
