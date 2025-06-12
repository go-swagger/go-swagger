package generate_test

import (
	"path/filepath"
	"testing"

	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
)

func TestMarkdown(t *testing.T) {
	path := filepath.Join(".", "test-markdown")
	generated, cleanup := testTempDir(t, path)
	t.Cleanup(cleanup)

	m := &generate.Markdown{}
	_, _ = flags.ParseArgs(m, []string{"--skip-validation"})
	m.Shared.Spec = flags.Filename(filepath.Join(testBase(), "fixtures", "enhancements", "184", "fixture-184.yaml"))
	m.Output = flags.Filename(filepath.Join(generated, "markdown.md"))
	require.NoError(t, m.Execute([]string{}))
}
