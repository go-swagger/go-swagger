// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package diff

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestSpecChangeCode(t *testing.T) {
	log.SetOutput(io.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := NoChangeDetected
	assert.EqualT(t, toLongStringSpecChangeCode[NoChangeDetected], c.Description())
	assert.EqualT(t, "UNDEFINED", SpecChangeCode(9999999999999).Description())

	res, err := json.Marshal(c)
	require.NoError(t, err)
	assert.JSONEqT(t, `"NoChangeDetected"`, string(res))

	var d SpecChangeCode
	in := []byte(`"NoChangeDetected"`)
	err = json.Unmarshal(in, &d)
	require.NoError(t, err)
	assert.EqualT(t, NoChangeDetected, d)

	in = []byte(`"dummy"`) // invalid enum
	err = json.Unmarshal(in, &d)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "unknown enum value")

	in = []byte(`{"dummy"`) // invalid json
	err = json.Unmarshal(in, &d)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "JSON")
}

func TestCompatibiliyCode(t *testing.T) {
	log.SetOutput(io.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	c := Breaking
	assert.EqualT(t, toStringCompatibility[Breaking], c.String())

	res, err := json.Marshal(c)
	require.NoError(t, err)
	assert.JSONEqT(t, `"Breaking"`, string(res))

	var d Compatibility
	in := []byte(`"Breaking"`)
	err = json.Unmarshal(in, &d)
	require.NoError(t, err)
	assert.EqualT(t, Breaking, d)

	in = []byte(`"dummy"`) // invalid enum
	err = json.Unmarshal(in, &d)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "unknown enum value")

	in = []byte(`{"dummy"`) // invalid json
	err = json.Unmarshal(in, &d)
	require.Error(t, err)
	assert.StringContainsT(t, err.Error(), "JSON")
}
