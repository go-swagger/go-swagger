// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

const nonExistingSpec = "nowhere.json"

// Commands requires at least one arg.
func TestCmd_Mixin(t *testing.T) {
	const (
		yamlFormat           = "yaml"
		otherNonExistingSpec = "notThere.json"
	)

	t.Run("should require an argument", func(t *testing.T) {
		var v MixinSpec
		result := v.Execute([]string{})
		require.Error(t, result)
	})

	t.Run("first spec file must exists", func(t *testing.T) {
		var v MixinSpec
		result := v.Execute([]string{nonExistingSpec})
		require.Error(t, result)
	})

	t.Run("both spec files must exists", func(t *testing.T) {
		var v MixinSpec
		result := v.Execute([]string{nonExistingSpec, otherNonExistingSpec})
		require.Error(t, result)
	})

	specDoc1 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	t.Run("second spec file must exists", func(t *testing.T) {
		var v MixinSpec
		result := v.Execute([]string{specDoc1, otherNonExistingSpec})
		require.Error(t, result)
	})

	specDoc2 := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536-2.yaml")
	t.Run("should merge specs", func(t *testing.T) {
		output := filepath.Join(t.TempDir(), "fixture-1536-mixed.yaml")

		v := MixinSpec{
			ExpectedCollisionCount: 3,
			Format:                 yamlFormat,
			Output:                 flags.Filename(output),
		}

		require.NoError(t, v.Execute([]string{specDoc1, specDoc2}))
		require.FileExists(t, output)
	})

	t.Run("should error on inconsistent flags - ignore conflicts and count collisions are incompatible", func(t *testing.T) {
		v := MixinSpec{
			ExpectedCollisionCount: 1,
			IgnoreConflicts:        true,
		}

		err := v.Execute([]string{"test.json", "test2.json"})
		require.Error(t, err)
		assert.EqualT(t, ignoreConflictsAndCollisionsSpecified, err.Error())
	})

	t.Run("should ignore conflicts when specified", func(t *testing.T) {
		output := filepath.Join(t.TempDir(), "fixture-1536-mixed.yaml")

		v := MixinSpec{
			IgnoreConflicts: true,
			Format:          yamlFormat,
			Output:          flags.Filename(output),
		}

		result := v.Execute([]string{specDoc1, specDoc2})
		require.NoError(t, result)
		assert.FileExists(t, output)
	})
}
