package diff

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/go-openapi/loads"
)

const (
	basePath = "../../../../fixtures/diff"
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
	matches := []string{}

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
					err = diffs.ReportAllDiffs(false)
					if diffs.BreakingChangeCount() > 0 {
						assertThat(t, err, is.Not(is.Nil()))
					}
				})
				assertThat(t, diffsStr, is.EqualToIgnoringWhitespace(tc.expectedLines))
			}
		})
	}

	assertThat(t, len(matches), is.EqualTo(len(allTests)).Reason("All test cases were not run. Remove filter."))
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

func dieOn(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
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
