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
	m.Shared.Spec = flags.Filename(filepath.Join(testBase(), "fixtures", "enhancements", "184", "fixture-184.yaml"))
	m.Shared.Target = flags.Filename(generated)
	m.Output = flags.Filename("markdown.md")
	require.NoError(t, m.Execute([]string{}))
}
