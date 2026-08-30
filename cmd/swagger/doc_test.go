// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestDoc(t *testing.T) {
	parser, err := register()
	require.NoError(t, err)

	d := &docCommand{
		Destination: t.TempDir(),
		parser:      parser,
	}

	require.NoError(t, d.gendoc())
}
