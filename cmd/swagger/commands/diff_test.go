package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"testing"

	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var assertThat = then.AssertThat
var equals = is.EqualTo

const (
	basePath = "../../../fixtures/diff"
)

type testCaseData struct {
	name          string
	oldSpec       string
	newSpec       string
	expectedLines string
}

// TestDiffForVariousCombinations - computes the diffs for a number
// of scenarios and compares the computed diff with expected diffs
func TestDiffForVariousCombinations(t *testing.T) {

	diffRootPath := basePath + "/"
	pattern := diffRootPath + "*.diff.txt"

	// To filter cases for debugging poke an individual case here eg "path", "enum" etc
	// see the test cases in fixtures/diff
	// Don't forget to remove it once you're done.
	// (There's a test at the end to check all cases were run)
	matches := []string{"path"}

	allTests, err := filepath.Glob(pattern)

	if err != nil || len(allTests) == 0 {
		t.Fatalf("Couldn't find files")
	}

	if len(matches) == 0 {
		matches = allTests
	}

	testCases := []testCaseData{}

	for _, eachFile := range matches {
		base := filepath.Base(eachFile)
		parts := strings.Split(base, ".diff.txt")
		namePart := parts[0]
		testCases = append(
			testCases, testCaseData{
				name:          namePart,
				oldSpec:       diffRootPath + namePart + ".v1.json",
				newSpec:       diffRootPath + namePart + ".v2.json",
				expectedLines: LinesInFile(diffRootPath + namePart + ".diff.txt"),
			})

	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			diffs, err := getDiffs(tc.oldSpec, tc.newSpec)

			assertThat(t, err, is.Nil())

			if err == nil {

				diffsStr := catchStdOut(t, func() {
					err = diffs.ReportAllDiffs(os.Stdout, false)
					assertThat(t, err, is.Not(is.Nil()))
				})
				assertThat(t, diffsStr, is.EqualToIgnoringWhitespace(tc.expectedLines))
			}
		})
	}
}

func TestDestParamRedirectsToFile(t *testing.T) {
	// create a random filename
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tmpfile.Close()
	os.Remove(tmpfile.Name()) // clean up

	diffRootPath := basePath + "/"
	// do a diff with that set to the output
	// should be not much sent to std out
	// output should be infile
	cmd := DiffCommand{
		Destination: tmpfile.Name(),
	}

	namePart := "enum"
	spec1 := diffRootPath + namePart + ".v1.json"
	spec2 := diffRootPath + namePart + ".v2.json"
	expectedDiffContent := LinesInFile(diffRootPath + namePart + ".diff.txt")

	diffsStr := catchStdOut(t, func() {
		err := cmd.Execute([]string{spec1, spec2})
		assertThat(t, err, is.Not(is.Nil()))
	})

	diffsInFile := LinesInFile(cmd.Destination)

	assertThat(t, diffsInFile, is.EqualToIgnoringWhitespace(expectedDiffContent))
	assertThat(t, diffsStr, is.EqualToIgnoringWhitespace(""))

}

func TestReadIgnores(t *testing.T) {

	diffRootPath := basePath + "/"
	ignorePath := diffRootPath + "ignoreFile.json"
	ignores, err := readIgnores(ignorePath)

	assertThat(t, err, is.Nil())
	assertThat(t, len(ignores), is.Not(equals(0)))

	isIn := diff.SpecDifference{DifferenceLocation: diff.DifferenceLocation{
		Method:   "get",
		Response: 0,
		URL:      "/a/",
		Node:     &diff.Node{Field: "Query", ChildNode: &diff.Node{Field: "personality"}},
	},
		Code:          diff.AddedEnumValue,
		Compatibility: diff.NonBreaking,
		DiffInfo:      "extrovert",
	}
	assertThat(t, ignores.Contains(isIn), equals(true))

}

func dieOn(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessIgnores(t *testing.T) {
	diffRootPath := basePath + "/"
	namePart := "enum"
	tc := testCaseData{
		name:          namePart,
		oldSpec:       diffRootPath + namePart + ".v1.json",
		newSpec:       diffRootPath + namePart + ".v2.json",
		expectedLines: LinesInFile(diffRootPath + "ignoreDiffs.json"),
	}

	cmd := DiffCommand{
		Format:      "json",
		IgnoreFile:  diffRootPath + "ignoreFile.json",
		Destination: "stdout",
	}

	diffsStr := catchStdOut(t, func() {
		err := cmd.Execute([]string{tc.oldSpec, tc.newSpec})
		assertThat(t, err, is.Nil())
	})
	assertThat(t, diffsStr, is.EqualToIgnoringWhitespace(tc.expectedLines))
}

func TestNoArgs(t *testing.T) {

	cmd := DiffCommand{
		Format:     "json",
		IgnoreFile: "",
	}

	err := cmd.Execute([]string{})
	assertThat(t, err, is.Not(is.Nil()))
}

func LinesInFile(fileName string) string {
	bytes, _ := ioutil.ReadFile(fileName)
	return string(bytes)
}

func catchStdOut(t *testing.T, runnable func()) string {

	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, err := os.Pipe()
	dieOn(err, t)
	os.Stdout = fakeStdout
	runnable()
	// need to close here, otherwise ReadAll never gets "EOF".
	dieOn(fakeStdout.Close(), t)
	newOutBytes, err := ioutil.ReadAll(r)
	dieOn(err, t)
	dieOn(r.Close(), t)
	return string(newOutBytes)
}
