package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/internal/cmdtest"
)

type testCaseData struct {
	name            string
	oldSpec         string
	newSpec         string
	expectedError   bool
	expectedWarning bool
	expectedLines   io.ReadCloser
	expectedFile    string
}

// TestDiffForVariousCombinations - computes the diffs for a number
// of scenarios and compares the computed diff with expected diffs.
func TestDiffForVariousCombinations(t *testing.T) {
	pattern := fixtureDiffPath("*.diff.txt")

	// To filter cases for debugging poke an individual case here eg "path", "enum" etc
	// see the test cases in fixtures/diff
	// Don't forget to remove it once you're done.
	// (There's a test at the end to check all cases were run)
	allTests, err := filepath.Glob(pattern)
	require.NoErrorf(t, err, "could not find test files")
	require.NotEmptyf(t, allTests, "could not find test files")

	testCases := makeTestCases(t, allTests)

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd := DiffCommand{}
			cmd.Args.OldSpec = tc.oldSpec
			cmd.Args.NewSpec = tc.newSpec
			diffs, err := cmd.getDiffs()

			if tc.expectedError {
				// edge cases with error
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			out, err, warn := diffs.ReportAllDiffs(false)
			require.NoError(t, err)

			// breaking changes reported with a warning
			if tc.expectedWarning {
				require.Error(t, warn)
			} else {
				require.NoError(t, warn)
			}

			if !cmdtest.AssertReadersContent(t, true, tc.expectedLines, out) {
				t.Logf("unexpected content for fixture %q[%d] (file: %s)", tc.name, i, tc.expectedFile)
			}
		})
	}
}

func TestDiffReadIgnores(t *testing.T) {
	cmd := DiffCommand{
		IgnoreFile: fixtureDiffPath("ignoreFile.json"),
	}

	ignores, err := cmd.readIgnores()
	require.NoError(t, err)
	require.NotEmpty(t, ignores)

	isIn := diff.SpecDifference{
		DifferenceLocation: diff.DifferenceLocation{
			Method:   "get",
			Response: 200,
			URL:      "/b/",
			Node:     &diff.Node{Field: "Body", TypeName: "A1", IsArray: true, ChildNode: &diff.Node{Field: "personality", TypeName: "string"}},
		},
		Code:          diff.DeletedEnumValue,
		Compatibility: diff.NonBreaking,
		DiffInfo:      "crazy",
	}
	assert.Contains(t, ignores, isIn)

	// edge case
	cmd = DiffCommand{
		IgnoreFile: "/someplace/wrong",
	}
	_, err = cmd.readIgnores()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/someplace/wrong")
}

func TestDiffProcessIgnores(t *testing.T) {
	const namePart = "enum"
	tc := testCaseData{
		// name:          namePart,
		oldSpec:       fixtureDiffPath(namePart, ".v1.json"),
		newSpec:       fixtureDiffPath(namePart, ".v2.json"),
		expectedLines: linesInFile(t, fixtureDiffPath("ignoreDiffs.json")),
	}

	reportFile, err := os.CreateTemp(t.TempDir(), "report.txt")
	require.NoError(t, err)
	require.NoError(t, reportFile.Close())

	defer func() {
		_ = os.Remove(reportFile.Name())
	}()

	cmd := DiffCommand{
		Format:      "json",
		IgnoreFile:  fixtureDiffPath("ignoreFile.json"),
		Destination: reportFile.Name(),
	}
	cmd.Args.OldSpec = tc.oldSpec
	cmd.Args.NewSpec = tc.newSpec

	require.NoError(t,
		cmd.Execute([]string{tc.oldSpec, tc.newSpec}),
	)

	output, err := os.Open(cmd.Destination)
	require.NoError(t, err)
	defer func() {
		_ = output.Close()
	}()

	cmdtest.AssertReadersContent(t, true, tc.expectedLines, output)
}

func TestDiffNoArgs(t *testing.T) {
	cmd := DiffCommand{
		Format:     "json",
		IgnoreFile: "",
	}
	require.Error(t, cmd.Execute(nil))

	cmd.Args.NewSpec = "x"
	require.Error(t, cmd.Execute(nil))
}

func TestDiffCannotReport(t *testing.T) {
	cmd := DiffCommand{
		OnlyBreakingChanges: true,
		Format:              "txt",
		IgnoreFile:          "",
		Destination:         "/someplace/wrong",
	}
	const namePart = "enum"
	cmd.Args.OldSpec = fixtureDiffPath(namePart, ".v1.json")
	cmd.Args.NewSpec = fixtureDiffPath(namePart, ".v2.json")
	err := cmd.Execute(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/someplace/wrong")
}

func TestDiffOnlyBreaking(t *testing.T) {
	reportDir := t.TempDir()
	txtReport := filepath.Join(reportDir, "report.txt")

	cmd := DiffCommand{
		OnlyBreakingChanges: true,
		Format:              "txt",
		IgnoreFile:          "",
		Destination:         txtReport,
	}

	t.Run("diff should return an error when breaking changes are detected", func(t *testing.T) {
		const namePart = "enum"
		cmd.Args.OldSpec = fixtureDiffPath(namePart, ".v1.json")
		cmd.Args.NewSpec = fixtureDiffPath(namePart, ".v2.json")
		err := cmd.Execute(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compatibility test FAILED")
	})

	t.Run("diff should correctly identify breaking changes", func(t *testing.T) {
		actual, err := os.Open(txtReport)
		require.NoError(t, err)
		defer func() {
			_ = actual.Close()
		}()

		expected, err := os.Open(fixtureDiffPath("enum", ".diff.breaking.txt"))
		require.NoError(t, err)
		defer func() {
			_ = expected.Close()
		}()

		cmdtest.AssertReadersContent(t, true, expected, actual)

		t.Run("same expectations should hold when output is stdout", func(t *testing.T) {
			// assert stdout just the same (we do it just once, so there is no race condition on os.Stdout)
			cmd.Destination = "stdout"
			output, err := cmdtest.CatchStdOut(t, func() error { return cmd.Execute(nil) })
			require.Error(t, err)
			assert.Contains(t, err.Error(), "compatibility test FAILED")

			_, _ = expected.Seek(0, io.SeekStart)
			result := bytes.NewBuffer(output)
			cmdtest.AssertReadersContent(t, true, expected, result)
		})
	})
}

func fixturePart(file string) string {
	base := filepath.Base(file)
	parts := strings.Split(base, ".diff.txt")
	return parts[0]
}

func hasFixtureBreaking(part string) bool {
	// these fixtures expect some breaking changes
	switch part {
	case "enum", "kitchensink", "param", "path", "response", "refprop", "reqparam":
		return true
	default:
		return false
	}
}

func makeTestCases(t testing.TB, matches []string) []testCaseData {
	t.Helper()

	testCases := make([]testCaseData, 0, len(matches)+2)
	for _, eachFile := range matches {
		namePart := fixturePart(eachFile)

		testCases = append(
			testCases, testCaseData{
				name:            namePart,
				oldSpec:         fixtureDiffPath(namePart, ".v1.json"),
				newSpec:         fixtureDiffPath(namePart, ".v2.json"),
				expectedLines:   linesInFile(t, fixtureDiffPath(namePart, ".diff.txt")),
				expectedFile:    fixtureDiffPath(namePart, ".diff.txt"), // only for debugging failed tests
				expectedWarning: hasFixtureBreaking(namePart),
			})
	}

	// edge cases with errors
	testCases = append(testCases, testCaseData{
		name:          "failure to load old spec",
		oldSpec:       "nowhere.json",
		newSpec:       fixtureDiffPath("enum", ".v2.json"),
		expectedError: true,
	},
		testCaseData{
			name:          "failure to load new spec",
			oldSpec:       fixtureDiffPath("enum", ".v1.json"),
			newSpec:       "nowhere.json",
			expectedError: true,
		},
	)
	return testCases
}

func fixtureDiffPath(file string, parts ...string) string {
	return filepath.Join(fixtureBase(), "diff", strings.Join(append([]string{file}, parts...), ""))
}

func linesInFile(t testing.TB, fileName string) io.ReadCloser {
	t.Helper()

	file, err := os.Open(fileName)
	require.NoError(t, err)

	return file
}
