package commands

import (
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
)

// Commands requires at least one arg
func TestCmd_Expand(t *testing.T) {
	v := &ExpandSpec{}
	testRequireParam(t, v)
}

func TestCmd_Expand_NoError(t *testing.T) {
	specDoc := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536.yaml")
	outDir, output := getOutput(t, specDoc, "flatten", "fixture-1536-flat-expand.json")
	defer os.RemoveAll(outDir)
	v := &ExpandSpec{
		Format:  "json",
		Compact: false,
		Output:  flags.Filename(output),
	}
	testProduceOutput(t, v, specDoc, output)
}

func TestCmd_Expand_Error(t *testing.T) {
	v := &ExpandSpec{}
	testValidRefs(t, v)
}
