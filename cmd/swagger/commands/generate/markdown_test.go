// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate_test

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"

	"github.com/go-openapi/testify/v2/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestMarkdown(t *testing.T) {
	generated := t.TempDir()

	m := &generate.Markdown{}
	_, _ = flags.ParseArgs(m, []string{"--skip-validation"})
	m.Shared.Spec = flags.Filename(filepath.Join(testBase(), "testdata", "enhancements", "184", "fixture-184.yaml"))
	m.Shared.Target = flags.Filename(generated)
	m.Output = flags.Filename("markdown.md")
	require.NoError(t, m.Execute([]string{}))
}

func TestMarkdownSpecArgument(t *testing.T) {
	spec := filepath.Join(testBase(), "testdata", "enhancements", "184", "fixture-184.yaml")

	newCmd := func(t *testing.T) *generate.Markdown {
		t.Helper()

		m := &generate.Markdown{}
		_, _ = flags.ParseArgs(m, []string{"--skip-validation"})
		m.Shared.Target = flags.Filename(t.TempDir())
		m.Output = flags.Filename("markdown.md")

		return m
	}

	t.Run("should generate with the spec passed as an argument", func(t *testing.T) {
		m := newCmd(t)
		require.NoError(t, m.Execute([]string{spec}))
		require.FileExists(t, filepath.Join(string(m.Shared.Target), "markdown.md"))
	})

	t.Run("should not accept the spec twice", func(t *testing.T) {
		m := newCmd(t)
		m.Shared.Spec = flags.Filename(spec)
		require.Error(t, m.Execute([]string{spec}))
	})

	t.Run("should not accept extra arguments", func(t *testing.T) {
		m := newCmd(t)
		require.Error(t, m.Execute([]string{spec, spec}))
	})
}
