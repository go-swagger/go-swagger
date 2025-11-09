// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"
)

// Commands requires at least one arg.
func TestCmd_Expand(t *testing.T) {
	v := &ExpandSpec{}
	testRequireParam(t, v)
}

func TestCmd_Expand_NoError(t *testing.T) {
	specDoc := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	output := filepath.Join(t.TempDir(), "fixture-1536-flat-expand.json")
	v := &ExpandSpec{
		Format:  "json",
		Compact: false,
		Output:  flags.Filename(output),
	}
	testProduceOutput(t, v, specDoc, output)
}

func TestCmd_Expand_NoOutputFile(t *testing.T) {
	defaultWriter = io.Discard
	t.Cleanup(func() {
		defaultWriter = os.Stdout
	})
	specDoc := filepath.Join(fixtureBase(), "bugs", "1536", "fixture-1536.yaml")
	v := &ExpandSpec{
		Format:  "json",
		Compact: false,
		Output:  "",
	}
	require.NoError(t, v.Execute([]string{specDoc}))
}

func TestCmd_Expand_Error(t *testing.T) {
	v := new(ExpandSpec)
	testValidRefs(t, v)
}
