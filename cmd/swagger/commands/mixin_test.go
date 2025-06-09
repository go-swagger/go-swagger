package commands

import (
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Commands requires at least one arg
func TestCmd_Mixin(t *testing.T) {
	var v MixinSpec
	result := v.Execute([]string{})
	require.Error(t, result)

	result = v.Execute([]string{"nowhere.json"})
	require.Error(t, result)

	result = v.Execute([]string{"nowhere.json", "notThere.json"})
	require.Error(t, result)

	specDoc1 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	result = v.Execute([]string{specDoc1, "notThere.json"})
	require.Error(t, result)
}

func TestCmd_Mixin_NoError(t *testing.T) {
	specDoc1 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	specDoc2 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536-2.yaml")
	outDir, err := os.MkdirTemp(filepath.Dir(specDoc1), "mixed")
	require.NoError(t, err)
	defer os.RemoveAll(outDir)
	v := MixinSpec{
		ExpectedCollisionCount: 3,
		Format:                 "yaml",
		Output:                 flags.Filename(filepath.Join(outDir, "fixture-1536-mixed.yaml")),
	}

	result := v.Execute([]string{specDoc1, specDoc2})
	require.NoError(t, result)
	_, exists := os.Stat(filepath.Join(outDir, "fixture-1536-mixed.yaml"))
	assert.False(t, os.IsNotExist(exists))
}

func TestCmd_Mixin_BothConflictsAndIgnoreConflictsSpecified(t *testing.T) {
	v := MixinSpec{
		ExpectedCollisionCount: 1,
		IgnoreConflicts:        true,
	}
	err := v.Execute([]string{"test.json", "test2.json"})
	require.Error(t, err)
	assert.Equal(t, ignoreConflictsAndCollisionsSpecified, err.Error())
}

func TestCmd_Mixin_IgnoreConflicts(t *testing.T) {
	specDoc1 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	specDoc2 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536-2.yaml")
	outDir, err := os.MkdirTemp(filepath.Dir(specDoc1), "mixed")
	require.NoError(t, err)
	defer os.RemoveAll(outDir)
	v := MixinSpec{
		IgnoreConflicts: true,
		Format:          "yaml",
		Output:          flags.Filename(filepath.Join(outDir, "fixture-1536-mixed.yaml")),
	}

	result := v.Execute([]string{specDoc1, specDoc2})
	require.NoError(t, result)
	_, exists := os.Stat(filepath.Join(outDir, "fixture-1536-mixed.yaml"))
	assert.False(t, os.IsNotExist(exists))
}
