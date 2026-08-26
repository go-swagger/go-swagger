// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate_test

import (
	"encoding/json"
	"os"
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

// TestMarkdownFlags asserts the flag set of the markdown command: it exposes the options that
// bear on the generated document, and none of those that only shape generated go code.
func TestMarkdownFlags(t *testing.T) {
	t.Run("should expose the options that alter the generated markdown", func(t *testing.T) {
		for _, arg := range []string{
			"--spec=swagger.yaml", "--skip-validation", "--with-expand", "--with-flatten=full",
			"--restricted", "--rooted=.",
			"--output=doc.md", "--target=.", "--ensure-target",
			"--template-dir=.", "--config-file=cfg.yaml", "--allow-template-override",
			"--additional-initialism=ABC", "--dump-data",
			"--model=Pet", "--keep-spec-order", "--all-definitions",
			"--operation=addPet", "--tags=pet", "--skip-tag-packages",
		} {
			_, err := flags.ParseArgs(&generate.Markdown{}, []string{arg})
			require.NoErrorf(t, err, "expected the markdown command to accept %q", arg)
		}
	})

	t.Run("should not expose the options that only shape generated go code", func(t *testing.T) {
		for _, arg := range []string{
			"--copyright-file=LICENSE", "--template=stratoscale", "--strict-responders",
			"--return-errors", "--with-custom-formatter",
			"--model-package=models", "--existing-models=github.com/foo/bar",
			"--strict-additional-properties", "--struct-tags=yaml", "--rooted-error-path",
			"--with-stringer", "--generate-getters", "--no-default-omit-empty",
			"--api-package=operations", "--with-enum-ci",
		} {
			_, err := flags.ParseArgs(&generate.Markdown{}, []string{arg})
			require.Errorf(t, err, "expected the markdown command to reject %q", arg)
		}
	})

	t.Run("should keep the full flag set on the code generation commands", func(t *testing.T) {
		for _, arg := range []string{
			"--copyright-file=LICENSE", "--strict-responders", "--return-errors",
			"--with-custom-formatter", "--struct-tags=yaml", "--with-stringer", "--generate-getters", "--with-enum-ci",
			"--model-package=models", "-m=models", "--api-package=operations", "-a=operations",
		} {
			_, err := flags.ParseArgs(&generate.Server{}, []string{arg})
			require.NoErrorf(t, err, "expected the server command to accept %q", arg)
		}
	})
}

func TestMarkdownDumpData(t *testing.T) {
	m := &generate.Markdown{}
	_, _ = flags.ParseArgs(m, []string{"--skip-validation"})
	m.Shared.Spec = flags.Filename(filepath.Join(testBase(), "testdata", "enhancements", "184", "fixture-184.yaml"))
	m.Shared.Target = flags.Filename(t.TempDir())
	m.Shared.DumpData = true
	m.Output = flags.Filename("markdown.md")

	// the dump is larger than a pipe buffer, so it is captured through a file rather than
	// with cmdtest.CatchStdOut.
	dump := filepath.Join(t.TempDir(), "dump.json")
	dumped := captureStdOutToFile(t, dump, func() error { return m.Execute([]string{}) })

	var data map[string]any
	require.NoError(t, json.Unmarshal(dumped, &data))
	require.Contains(t, data, "Models")

	_, err := os.Stat(filepath.Join(string(m.Shared.Target), "markdown.md"))
	require.ErrorIs(t, err, os.ErrNotExist)
}

// captureStdOutToFile redirects the standard output to a file while runnable executes,
// and returns what was written.
func captureStdOutToFile(t *testing.T, name string, runnable func() error) []byte {
	t.Helper()

	fakeStdout, err := os.Create(name)
	require.NoError(t, err)

	realStdout := os.Stdout
	os.Stdout = fakeStdout
	runErr := runnable()
	os.Stdout = realStdout

	require.NoError(t, fakeStdout.Close())
	require.NoError(t, runErr)

	captured, err := os.ReadFile(name)
	require.NoError(t, err)

	return captured
}
