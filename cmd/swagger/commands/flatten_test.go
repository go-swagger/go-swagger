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

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

type executable interface {
	Execute(args []string) error
}

// Commands requires at least one arg.
func TestCmd_Flatten(t *testing.T) {
	v := &FlattenSpec{}
	testRequireParam(t, v)
}

func TestCmd_Flatten_Default(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	output := filepath.Join(t.TempDir(), "fixture-1536-flat-minimal.json")
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
	output := filepath.Join(t.TempDir(), "fixture-2919-flat-minimal.yml")

	v := &FlattenSpec{
		Format:  "yaml",
		Compact: true,
		Output:  flags.Filename(output),
	}
	testProduceOutput(t, v, specDoc, output)
}

func TestCmd_FlattenKeepNames_Issue2334(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "2334", "swagger.yaml")
	output := filepath.Join(t.TempDir(), "fixture-2334-flat-keep-names.yaml")

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
	t.Helper()

	specDoc := filepath.Join(fixtureBase(), "expansion", "invalid-refs.json")
	result := v.Execute([]string{specDoc})
	require.Error(t, result)
}

func testRequireParam(t *testing.T, v executable) {
	t.Helper()

	result := v.Execute([]string{})
	require.Error(t, result)

	result = v.Execute([]string{"nowhere.json"})
	require.Error(t, result)
}

func testProduceOutput(t *testing.T, v executable, specDoc, output string) {
	t.Helper()

	require.NoError(t, v.Execute([]string{specDoc}))
	_, exists := os.Stat(output)
	assert.False(t, os.IsNotExist(exists))
}
