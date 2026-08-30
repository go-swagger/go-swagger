// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestSwagger(t *testing.T) {
	stdout := os.Stdout
	t.Cleanup(func() {
		os.Stdout = stdout
	})
	os.Stdout = discard(t)
	parser, err := register()
	require.NoError(t, err)

	err = run(parser, []string{"--help"})
	require.Error(t, err)
}

func discard(t *testing.T) *os.File {
	discard, err := os.Open(os.DevNull)
	require.NoError(t, err)

	return discard
}
