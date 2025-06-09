package generate_test

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestGenerateOperation(t *testing.T) {
	testGenerateOperation(t, false)
}

func TestGenerateOperationStrict(t *testing.T) {
	testGenerateOperation(t, true)
}

func testGenerateOperation(t *testing.T, strict bool) {
	specs := []string{
		"tasklist.basic.yml",
	}

	for i, spec := range specs {
		t.Run(spec, func(t *testing.T) {
			path := filepath.Join(testBase(), "fixtures/codegen", spec)
			generated, cleanup := testTempDir(t, path)
			t.Cleanup(cleanup)

			m := &generate.Operation{}
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(testBase(), "LICENSE"))
			}
			_, _ = flags.ParseArgs(m, []string{"--name=listTasks"})
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			require.NoError(t, m.Execute([]string{}))
		})
	}
}

func TestGenerateOperation_Check(t *testing.T) {
	m := &generate.Operation{}
	_, _ = flags.ParseArgs(m, []string{"--name=op1", "--name=op2"})
	m.Shared.DumpData = true
	m.Name = []string{"op1", "op2"}

	require.Error(t, m.Execute([]string{}))
}
