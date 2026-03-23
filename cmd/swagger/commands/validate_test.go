// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestCmd_Validate(t *testing.T) {
	t.Run("should require an argument", func(t *testing.T) {
		// Command requires at least one arg.
		var v ValidateSpec
		require.Error(t, v.Execute([]string{}))
	})

	t.Run("spec file must exists", func(t *testing.T) {
		var v ValidateSpec
		require.Error(t, v.Execute([]string{nonExistingSpec}))
	})

	t.Run("should detect invalid spec - items in object", func(t *testing.T) {
		// issue #1238
		var v ValidateSpec
		specDoc := filepath.Join(fixtureBase(), "bugs", "1238", "swagger.yaml")
		result := v.Execute([]string{specDoc})
		require.Error(t, result)
		assert.StringContainsT(t, result.Error(), "is invalid against swagger specification 2.0")
		assert.StringContainsT(t, result.Error(), "definitions.RRSets in body must be of type array")
	})

	t.Run("should detect invalid spec - missing items in array", func(t *testing.T) {
		// issue #1171
		var v ValidateSpec
		specDoc := filepath.Join(fixtureBase(), "bugs", "1171", "swagger.yaml")
		require.Error(t, v.Execute([]string{specDoc}))
	})

	t.Run("should detect invalid spec - reference to inner (forbidden) property", func(t *testing.T) {
		// issue #342
		var v ValidateSpec
		specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342.yaml")
		require.Error(t, v.Execute([]string{specDoc}))
	})

	t.Run("should detect invalid spec - cannot unmarshal", func(t *testing.T) {
		// fixture 342-2 (a variant of invalid specification) (cannot unmarshal)
		// Test proper validation: reference to shared top level parameter, but with incorrect
		// yaml syntax: use map key instead of array item.
		//
		// NOTE: this error message is not clear enough. The role of this test
		// is to determine that the validation does not panic and correctly states the spec is invalid.
		// Open a dedicated issue on message relevance. This test shall be updated with the finalized message.
		v := ValidateSpec{}
		specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342-2.yaml")
		require.NotPanics(t, func() {
			_ = v.Execute([]string{specDoc})
		})

		result := v.Execute([]string{specDoc})
		require.Error(t, result, "This spec should not pass validation")
		assert.ErrorContains(t, result, "json: cannot unmarshal object into Go struct field")
		assert.ErrorContains(t, result, "of type []spec.Parameter")
	})

	t.Run("should detect valid spec", func(t *testing.T) {
		// This one is a correct version of issue#342 and it validates.
		var v ValidateSpec
		specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342-3.yaml")
		require.NoError(t, v.Execute([]string{specDoc}))
	})
}
