// Package cmdtest provides test utilities
// to assert the output of CLI commands
package cmdtest

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertReadersContent compares the contents from io.Readers, optionally stripping blanks
func AssertReadersContent(t testing.TB, noBlanks bool, expected, actual io.Reader) bool {
	t.Helper()

	e, err := io.ReadAll(expected)
	require.NoError(t, err)

	a, err := io.ReadAll(actual)
	require.NoError(t, err)

	var wants, got strings.Builder
	_, _ = wants.Write(e)
	_, _ = got.Write(a)

	if noBlanks {
		r := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "")
		return assert.Equalf(t, r.Replace(wants.String()), r.Replace(got.String()), "expected:\n%s\ngot:\n%s", wants.String(), got.String())
	}
	return assert.Equal(t, wants.String(), got.String())
}

// CatchStdOut captures the standard output from a runnable function.
// You shouln't run this in parallel.
func CatchStdOut(t testing.TB, runnable func() error) ([]byte, error) {
	t.Helper()

	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = fakeStdout
	e := runnable()
	// need to close here, otherwise ReadAll never gets "EOF".
	require.NoError(t, fakeStdout.Close())
	newOutBytes, err := io.ReadAll(r)
	require.NoError(t, err)
	require.NoError(t, r.Close())
	return newOutBytes, e
}
