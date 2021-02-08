package generate_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
	flags "github.com/jessevdk/go-flags"
)

func TestMarkdown(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	base := filepath.FromSlash("../../../../")

	generated, err := ioutil.TempDir(".", "test-markdown")
	require.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(generated)
	}()

	m := &generate.Markdown{}
	_, _ = flags.ParseArgs(m, []string{"--skip-validation"})
	m.Shared.Spec = flags.Filename(filepath.Join(base, "fixtures", "enhancements", "184", "fixture-184.yaml"))
	m.Output = flags.Filename(filepath.Join(generated, "markdown.md"))
	require.NoError(t, m.Execute([]string{}))
}
