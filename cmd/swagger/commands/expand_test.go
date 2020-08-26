package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
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

func TestCmd_Expand_NoOutputFile(t *testing.T) {
	specDoc := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536.yaml")
	v := &ExpandSpec{
		Format:  "json",
		Compact: false,
		Output:  "",
	}
	result := v.Execute([]string{specDoc})
	assert.Nil(t, result)
}

func TestCmd_Expand_DeterministicWithCyclicRef(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	specDoc := filepath.Join(fixtureBase, "bugs", "2393", "fixture-2393.yaml")
	set := map[string]struct{}{}

	// redirect stdout to a pipe which in turns redirect the input to a buffer
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	for i := 0; i < 100; i++ {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w

		v := &ExpandSpec{
			Format:  "json",
			Compact: false,
			Output:  "",
		}

		err = v.Execute([]string{specDoc})
		assert.NoError(t, err)

		w.Close()

		result, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

		set[string(result)] = struct{}{}
	}

	assert.True(t, len(set) == 1)
}

func TestCmd_Expand_Error(t *testing.T) {
	v := &ExpandSpec{}
	testValidRefs(t, v)
}
