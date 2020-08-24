package diff

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/internal/cmdtest"
	"github.com/stretchr/testify/require"
)

func fixturePath(file string, parts ...string) string {
	return filepath.Join("..", "..", "..", "..", "fixtures", "diff", strings.Join(append([]string{file}, parts...), ""))
}

type testCaseData struct {
	name          string
	oldSpec       string
	newSpec       string
	expectedLines io.Reader
	expectedFile  string
}

func fixturePart(file string) string {
	base := filepath.Base(file)
	parts := strings.Split(base, ".diff.txt")
	return parts[0]
}

// TestDiffForVariousCombinations - computes the diffs for a number
// of scenarios and compares the computed diff with expected diffs
func TestDiffForVariousCombinations(t *testing.T) {
	pattern := fixturePath("*.diff.txt")
	allTests, err := filepath.Glob(pattern)
	require.NoError(t, err)
	require.True(t, len(allTests) > 0)

	// To filter cases for debugging poke an individual case here eg "path", "enum" etc
	// see the test cases in fixtures/diff
	// Don't forget to remove it once you're done.
	// (There's a test at the end to check all cases were run)
	matches := allTests
	// matches := []string{"enum"}

	testCases := makeTestCases(t, matches)

	for i, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			diffs, err := getDiffs(tc.oldSpec, tc.newSpec)
			require.NoError(t, err)

			out, err, warn := diffs.ReportAllDiffs(false)
			require.NoError(t, err)

			if !cmdtest.AssertReadersContent(t, true, tc.expectedLines, out) {
				t.Logf("unexpected content for fixture %q[%d] (file: %s)", tc.name, i, tc.expectedFile)
			}

			if diffs.BreakingChangeCount() > 0 {
				require.Error(t, warn)
			}
		})
	}

	require.Equalf(t, len(allTests), len(matches), "All test cases were not run. Remove filter")
}

func getDiffs(oldSpecPath, newSpecPath string) (SpecDifferences, error) {
	swaggerDoc1 := oldSpecPath
	specDoc1, err := loads.Spec(swaggerDoc1)

	if err != nil {
		return nil, err
	}

	swaggerDoc2 := newSpecPath
	specDoc2, err := loads.Spec(swaggerDoc2)
	if err != nil {
		return nil, err
	}

	return Compare(specDoc1.Spec(), specDoc2.Spec())
}

func makeTestCases(t testing.TB, matches []string) []testCaseData {
	testCases := make([]testCaseData, 0, len(matches))
	for _, eachFile := range matches {
		namePart := fixturePart(eachFile)
		if _, err := os.Stat(fixturePath(namePart, ".v1.json")); err == nil {
			testCases = append(
				testCases, testCaseData{
					name:          namePart,
					oldSpec:       fixturePath(namePart, ".v1.json"),
					newSpec:       fixturePath(namePart, ".v2.json"),
					expectedLines: linesInFile(t, fixturePath(namePart, ".diff.txt")),
				})
		}
		if _, err := os.Stat(fixturePath(namePart, ".v1.yml")); err == nil {
			testCases = append(
				testCases, testCaseData{
					name:          namePart,
					oldSpec:       fixturePath(namePart, ".v1.yml"),
					newSpec:       fixturePath(namePart, ".v2.yml"),
					expectedLines: linesInFile(t, fixturePath(namePart, ".diff.txt")),
				})
		}
	}
	return testCases
}

func linesInFile(t testing.TB, fileName string) io.ReadCloser {
	file, err := os.Open(fileName)
	require.NoError(t, err)
	return file
}
