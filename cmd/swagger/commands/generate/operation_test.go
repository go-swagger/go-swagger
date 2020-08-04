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
			m := &generate.Operation{}
			if i == 0 {
				m.Shared.CopyrightFile = flags.Filename(filepath.Join(base, "LICENSE"))
			}
			_, _ = flags.ParseArgs(m, []string{"--name=listTasks"})
			m.Shared.Spec = flags.Filename(path)
			m.Shared.Target = flags.Filename(generated)
			m.Shared.StrictResponders = strict

			if err := m.Execute([]string{}); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateOperation_Check(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	m := &generate.Operation{}
	_, _ = flags.ParseArgs(m, []string{"--name=op1", "--name=op2"})
	m.Shared.DumpData = true
	m.Name = []string{"op1", "op2"}
	err := m.Execute([]string{})
	assert.Error(t, err)
}
