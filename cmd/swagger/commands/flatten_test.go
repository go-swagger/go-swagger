package commands

import (
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

type executable interface {
	Execute([]string) error
}

// Commands requires at least one arg
func TestCmd_Flatten(t *testing.T) {
	v := &FlattenSpec{}
	testRequireParam(t, v)
}

func TestCmd_Flatten_Default(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	outDir, output := getOutput(t, specDoc, "flatten", "fixture-1536-flat-minimal.json")
	defer os.RemoveAll(outDir)
	v := &FlattenSpec{
		Format:  "json",
		Compact: true,
		Output:  flags.Filename(output),
	}
	testProduceOutput(t, v, specDoc, output)
}

func TestCmd_Flatten_Error(t *testing.T) {
	v := &FlattenSpec{}
	testValidRefs(t, v)
}

func TestCmd_Flatten_Issue2919(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "2919", "edge-api", "client.yml")
	outDir, output := getOutput(t, specDoc, "flatten", "fixture-2919-flat-minimal.yml")
	defer os.RemoveAll(outDir)

	v := &FlattenSpec{
		Format:  "yaml",
		Compact: true,
		Output:  flags.Filename(output),
	}
	testProduceOutput(t, v, specDoc, output)
}

func TestCmd_FlattenKeepNames_Issue2334(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "2334", "swagger.yaml")
	outDir, output := getOutput(t, specDoc, "flatten", "fixture-2334-flat-keep-names.yaml")
	defer os.RemoveAll(outDir)

	v := &FlattenSpec{
		Format:  "yaml",
		Compact: true,
		Output:  flags.Filename(output),
		FlattenCmdOptions: generate.FlattenCmdOptions{
			WithFlatten: []string{"keep-names"},
		},
	}
	testProduceOutput(t, v, specDoc, output)
	buf, err := os.ReadFile(output)
	require.NoError(t, err)
	spec := string(buf)

	require.Contains(t, spec, "$ref: '#/definitions/Bar'")
	require.Contains(t, spec, "Bar:")
	require.Contains(t, spec, "Baz:")
}

func testValidRefs(t *testing.T, v executable) {
	specDoc := filepath.Join(fixtureBase(), "expansion", "invalid-refs.json")
	result := v.Execute([]string{specDoc})
	require.Error(t, result)
}

func testRequireParam(t *testing.T, v executable) {
	result := v.Execute([]string{})
	require.Error(t, result)

	result = v.Execute([]string{"nowhere.json"})
	require.Error(t, result)
}

func getOutput(t *testing.T, specDoc, _, filename string) (string, string) {
	outDir, err := os.MkdirTemp(filepath.Dir(specDoc), "flatten")
	require.NoError(t, err)
	return outDir, filepath.Join(outDir, filename)
}

func testProduceOutput(t *testing.T, v executable, specDoc, output string) {
	require.NoError(t, v.Execute([]string{specDoc}))
	_, exists := os.Stat(output)
	assert.False(t, os.IsNotExist(exists))
}
