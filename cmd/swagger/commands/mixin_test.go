// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Commands requires at least one arg.
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
	output := filepath.Join(t.TempDir(), "fixture-1536-mixed.yaml")

	v := MixinSpec{
		ExpectedCollisionCount: 3,
		Format:                 "yaml",
		Output:                 flags.Filename(output),
	}

	require.NoError(t, v.Execute([]string{specDoc1, specDoc2}))
	require.FileExists(t, output)
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
	output := filepath.Join(t.TempDir(), "fixture-1536-mixed.yaml")

	v := MixinSpec{
		IgnoreConflicts: true,
		Format:          "yaml",
		Output:          flags.Filename(output),
	}

	result := v.Execute([]string{specDoc1, specDoc2})
	require.NoError(t, result)
	_, exists := os.Stat(output)
	assert.False(t, os.IsNotExist(exists))
}
