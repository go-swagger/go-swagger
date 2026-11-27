// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestSwagger(t *testing.T) {
	stdout := os.Stdout
	t.Cleanup(func() {
		os.Stdout = stdout
	})
	os.Stdout = discard(t)
	var doc docCommand
	parser, err := register(&doc)
	require.NoError(t, err)

	err = run(parser, []string{"--help"})
	require.Error(t, err)
}

func discard(t *testing.T) *os.File {
	discard, err := os.Open(os.DevNull)
	require.NoError(t, err)

	return discard
}

func TestVerify(t *testing.T) {
	invoke := func(t *testing.T, args ...string) error {
		t.Helper()
		parser, err := register(new(docCommand))
		require.NoError(t, err)
		return run(parser, args)
	}

	t.Run(
		"requires a spec",
		func(t *testing.T) {
			assert.ErrorContains(t, invoke(t, "verify"), "exactly one")
		},
	)

	t.Run("rejects a spec supplied twice", func(t *testing.T) {
		assert.ErrorContains(
			t,
			invoke(t, "verify", "--spec", "swagger.yaml", "swagger.yaml"),
			"specified twice",
		)
	})

	generateServer := func(t *testing.T) string {
		t.Helper()
		spec, err := filepath.Abs("../../testdata/codegen/tasklist.basic.yml")
		require.NoError(t, err)
		t.Chdir(t.TempDir())
		require.NoError(t, os.WriteFile("go.mod", []byte("module example.org/verify\n"), 0o600))
		require.NoError(t, invoke(t, "generate", "server", spec))
		return spec
	}

	t.Run("accepts current server files with either spec form", func(t *testing.T) {
		spec := generateServer(t)
		assert.NoError(t, invoke(t, "verify", spec))
		assert.NoError(t, invoke(t, "verify", "--spec", spec))
	})

	t.Run("reports missing and edited files without writing them", func(t *testing.T) {
		spec := generateServer(t)
		file := filepath.Join("models", "price.go")
		require.NoError(t, os.Remove(file))
		stale := filepath.Join("models", "currency.go")
		edit := []byte("// local edit\n")
		require.NoError(t, os.WriteFile(stale, edit, 0o600))
		err := invoke(t, "verify", spec)
		assert.ErrorContains(t, err, file)
		assert.ErrorContains(t, err, stale)
		data, err := os.ReadFile(stale)
		require.NoError(t, err)
		assert.Equal(t, edit, data)
		_, err = os.Stat(file)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}
