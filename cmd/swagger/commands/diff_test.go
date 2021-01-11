package commands

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/internal/cmdtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixturePath(file string, parts ...string) string {
	return filepath.Join("..", "..", "..", "fixtures", "diff", strings.Join(append([]string{file}, parts...), ""))
}

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
// of scenarios and compares the computed diff with expected diffs
func TestDiffForVariousCombinations(t *testing.T) {

	pattern := fixturePath("*.diff.txt")

	// To filter cases for debugging poke an individual case here eg "path", "enum" etc
	// see the test cases in fixtures/diff
	// Don't forget to remove it once you're done.
	// (There's a test at the end to check all cases were run)
	allTests, err := filepath.Glob(pattern)
	require.NoErrorf(t, err, "could not find test files")
	require.False(t, len(allTests) == 0, "could not find test files")

	testCases := makeTestCases(t, allTests)

	for i, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
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
				assert.Error(t, warn)
			} else {
				assert.NoError(t, warn)
			}

			if !cmdtest.AssertReadersContent(t, true, tc.expectedLines, out) {
				t.Logf("unexpected content for fixture %q[%d] (file: %s)", tc.name, i, tc.expectedFile)
			}
		})
	}
}

func TestDiffReadIgnores(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	cmd := DiffCommand{
		IgnoreFile: fixturePath("ignoreFile.json"),
	}

	ignores, err := cmd.readIgnores()
	require.NoError(t, err)
	require.True(t, len(ignores) > 0)

	isIn := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{
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
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	const namePart = "enum"
	tc := testCaseData{
		name:          namePart,
		oldSpec:       fixturePath(namePart, ".v1.json"),
		newSpec:       fixturePath(namePart, ".v2.json"),
		expectedLines: linesInFile(t, fixturePath("ignoreDiffs.json")),
	}

	reportFile, err := ioutil.TempFile("", "report.txt")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(reportFile.Name())
	}()

	cmd := DiffCommand{
		Format:      "json",
		IgnoreFile:  fixturePath("ignoreFile.json"),
		Destination: reportFile.Name(),
	}
	cmd.Args.OldSpec = tc.oldSpec
	cmd.Args.NewSpec = tc.newSpec

	err = cmd.Execute([]string{tc.oldSpec, tc.newSpec})
	require.NoError(t, err)

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
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	cmd := DiffCommand{
		OnlyBreakingChanges: true,
		Format:              "txt",
		IgnoreFile:          "",
		Destination:         "/someplace/wrong",
	}
	const namePart = "enum"
	cmd.Args.OldSpec = fixturePath(namePart, ".v1.json")
	cmd.Args.NewSpec = fixturePath(namePart, ".v2.json")
	err := cmd.Execute(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "/someplace/wrong")
}

func TestDiffOnlyBreaking(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	reportDir, err := ioutil.TempDir("", "diff-reports")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(reportDir)
	}()
	txtReport := filepath.Join(reportDir, "report.txt")

	cmd := DiffCommand{
		OnlyBreakingChanges: true,
		Format:              "txt",
		IgnoreFile:          "",
		Destination:         txtReport,
	}

	const namePart = "enum"
	cmd.Args.OldSpec = fixturePath(namePart, ".v1.json")
	cmd.Args.NewSpec = fixturePath(namePart, ".v2.json")
	err = cmd.Execute(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "compatibility test FAILED")

	actual, err := os.Open(txtReport)
	require.NoError(t, err)
	defer func() {
		_ = actual.Close()
	}()

	expected, err := os.Open(fixturePath("enum", ".diff.breaking.txt"))
	require.NoError(t, err)
	defer func() {
		_ = expected.Close()
	}()

	cmdtest.AssertReadersContent(t, true, expected, actual)

	// assert stdout just the same (we do it just once, so there is no race condition on os.Stdout)
	cmd.Destination = "stdout"
	output, err := cmdtest.CatchStdOut(t, func() error { return cmd.Execute(nil) })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "compatibility test FAILED")

	_, _ = expected.Seek(0, io.SeekStart)
	result := bytes.NewBuffer(output)
	cmdtest.AssertReadersContent(t, true, expected, result)
}

func fixturePart(file string) string {
	base := filepath.Base(file)
	parts := strings.Split(base, ".diff.txt")
	return parts[0]
}

func hasFixtureBreaking(part string) bool {
	// these fixtures expect some breaking changes
	switch part {
	case "enum", "kitchensink", "param", "path", "response", "refprop":
		return true
	default:
		return false
	}
}

func makeTestCases(t testing.TB, matches []string) []testCaseData {
	testCases := make([]testCaseData, 0, len(matches)+2)
	for _, eachFile := range matches {
		namePart := fixturePart(eachFile)

		testCases = append(
			testCases, testCaseData{
				name:            namePart,
				oldSpec:         fixturePath(namePart, ".v1.json"),
				newSpec:         fixturePath(namePart, ".v2.json"),
				expectedLines:   linesInFile(t, fixturePath(namePart, ".diff.txt")),
				expectedFile:    fixturePath(namePart, ".diff.txt"), // only for debugging failed tests
				expectedWarning: hasFixtureBreaking(namePart),
			})
	}

	// edge cases with errors
	testCases = append(testCases, testCaseData{
		name:          "failure to load old spec",
		oldSpec:       "nowhere.json",
		newSpec:       fixturePath("enum", ".v2.json"),
		expectedError: true,
	},
		testCaseData{
			name:          "failure to load new spec",
			oldSpec:       fixturePath("enum", ".v1.json"),
			newSpec:       "nowhere.json",
			expectedError: true,
		},
	)
	return testCases
}

func linesInFile(t testing.TB, fileName string) io.ReadCloser {
	file, err := os.Open(fileName)
	require.NoError(t, err)
	return file
}
