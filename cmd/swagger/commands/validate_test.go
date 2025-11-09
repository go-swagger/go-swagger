// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Commands requires at least one arg.
func TestCmd_Validate_MissingArgs(t *testing.T) {
	var v ValidateSpec
	require.Error(t, v.Execute([]string{}))
	require.Error(t, v.Execute([]string{"nowhere.json"}))
}

// Test proper validation: items in object error.
func TestCmd_Validate_Issue1238(t *testing.T) {
	var v ValidateSpec
	specDoc := filepath.Join(fixtureBase(), "bugs", "1238", "swagger.yaml")
	result := v.Execute([]string{specDoc})
	require.Error(t, result)
	/*
		The swagger spec at "../../../fixtures/bugs/1238/swagger.yaml" is invalid against swagger specification 2.0. see errors :
			- definitions.RRSets in body must be of type array
	*/
	assert.Contains(t, result.Error(), "is invalid against swagger specification 2.0")
	assert.Contains(t, result.Error(), "definitions.RRSets in body must be of type array")
}

// Test proper validation: missing items in array error.
func TestCmd_Validate_Issue1171(t *testing.T) {
	var v ValidateSpec
	specDoc := filepath.Join(fixtureBase(), "bugs", "1171", "swagger.yaml")
	require.Error(t, v.Execute([]string{specDoc}))
}

// Test proper validation: reference to inner property in schema.
func TestCmd_Validate_Issue342_ForbiddenProperty(t *testing.T) {
	var v ValidateSpec
	specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342.yaml")
	require.Error(t, v.Execute([]string{specDoc}))
}

// fixture 342-2 (a variant of invalid specification) (cannot unmarshal)
// Test proper validation: reference to shared top level parameter, but with incorrect
// yaml syntax: use map key instead of array item.
//
// NOTE: this error message is not clear enough. The role of this test
// is to determine that the validation does not panic and correctly states the spec is invalid.
// Open a dedicated issue on message relevance. This test shall be updated with the finalized message.
func TestCmd_Validate_Issue342_CannotUnmarshal(t *testing.T) {
	v := ValidateSpec{}
	specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342-2.yaml")
	require.NotPanics(t, func() {
		_ = v.Execute([]string{specDoc})
	})

	result := v.Execute([]string{specDoc})
	require.Error(t, result, "This spec should not pass validation")
	// assert.Contains(t, result.Error(), "is invalid against swagger specification 2.0")
	assert.Contains(t, result.Error(), "json: cannot unmarshal object into Go struct field")
	assert.Contains(t, result.Error(), "of type []spec.Parameter")
}

// This one is a correct version of issue#342 and it validates.
func TestCmd_Validate_Issue342_Correct(t *testing.T) {
	var v ValidateSpec
	specDoc := filepath.Join(fixtureBase(), "bugs", "342", "fixture-342-3.yaml")
	require.NoError(t, v.Execute([]string{specDoc}))
}
