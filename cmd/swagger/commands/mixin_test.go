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

var fixtureBase = filepath.FromSlash("../../../fixtures")

// Commands requires at least one arg
func TestCmd_Mixin(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	v := MixinSpec{}
	result := v.Execute([]string{})
	assert.Error(t, result)

	result = v.Execute([]string{"nowhere.json"})
	assert.Error(t, result)

	result = v.Execute([]string{"nowhere.json", "notThere.json"})
	assert.Error(t, result)

	specDoc1 := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536.yaml")
	result = v.Execute([]string{specDoc1, "notThere.json"})
	assert.Error(t, result)
}

func TestCmd_Mixin_NoError(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	specDoc1 := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536.yaml")
	specDoc2 := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536-2.yaml")
	outDir, err := ioutil.TempDir(filepath.Dir(specDoc1), "mixed")
	assert.NoError(t, err)
	defer os.RemoveAll(outDir)
	v := MixinSpec{
		ExpectedCollisionCount: 3,
		Format:                 "yaml",
		Output:                 flags.Filename(filepath.Join(outDir, "fixture-1536-mixed.yaml")),
	}

	result := v.Execute([]string{specDoc1, specDoc2})
	assert.NoError(t, result)
	_, exists := os.Stat(filepath.Join(outDir, "fixture-1536-mixed.yaml"))
	assert.True(t, !os.IsNotExist(exists))
}

func TestCmd_Mixin_BothConflictsAndIgnoreConflictsSpecified(t *testing.T) {
	v := MixinSpec{
		ExpectedCollisionCount: 1,
		IgnoreConflicts:        true,
	}
	err := v.Execute([]string{"test.json", "test2.json"})
	assert.Error(t, err)
	assert.Equal(t, ignoreConflictsAndCollisionsSpecified, err.Error())
}

func TestCmd_Mixin_IgnoreConflicts(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	specDoc1 := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536.yaml")
	specDoc2 := filepath.Join(fixtureBase, "bugs", "1536", "fixture-1536-2.yaml")
	outDir, err := ioutil.TempDir(filepath.Dir(specDoc1), "mixed")
	assert.NoError(t, err)
	defer os.RemoveAll(outDir)
	v := MixinSpec{
		IgnoreConflicts: true,
		Format:          "yaml",
		Output:          flags.Filename(filepath.Join(outDir, "fixture-1536-mixed.yaml")),
	}

	result := v.Execute([]string{specDoc1, specDoc2})
	assert.NoError(t, result)
	_, exists := os.Stat(filepath.Join(outDir, "fixture-1536-mixed.yaml"))
	assert.True(t, !os.IsNotExist(exists))
}
