// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type relativePathTest struct {
	childpath  string
	parentpath string
	ok         bool
	path       string
}

func prefixAndFetchRelativePathFixtures() []relativePathTest {
	return []relativePathTest{
		// Positive
		{"/", "/", true, "."},
		{"/User/Gopher", "/", true, "User/Gopher"},
		{"/User/Gopher/Go", "/User/Gopher/Go", true, "."},
		{"/User/../User/Gopher", "/", true, "User/Gopher"},
		// Negative cases
		{"/", "/var", false, ""},
		{"/User/Gopher", "/User/SomethingElse", false, ""},
		{"/var", "/etc", false, ""},
		{"/mnt/dev3", "/mnt/dev3/dir", false, ""},
	}
}

type baseImportTest struct {
	title        string
	path         []string
	gopath       string
	targetpath   string
	symlinksrc   string
	symlinkdest  string // symlink is the last dir in targetpath
	expectedpath string
}

func baseImportTestFixtures(tempdir string) []baseImportTest {
	suffixParts := []string{"github.com", "go-swagger"}
	goParts := []string{"root", "go"} // root/go
	srcParts := goParts
	srcParts = append(srcParts, "src") // e.g. root/go/src

	pathParts := srcParts
	pathParts = append(pathParts, suffixParts...) // e.g. root/go/src/github.com/go-swagger
	tmp := []string{tempdir}

	tmpGoParts := tmp
	tmpGoParts = append(tmpGoParts, goParts...) // e.g. /tmp/root/go

	tmpSrcParts := tmp
	tmpSrcParts = append(tmpSrcParts, srcParts...) // e.g. /tmp/root/go/src

	tmpPathParts := tmp
	tmpPathParts = append(tmpPathParts, pathParts...) // e.g. /tmp/root/go/src/github.com/go-swagger

	tmpSymLinkParts := []string{tempdir, "root", "symlink"}

	return []baseImportTest{
		{
			title:        "No sym link. Positive Test Case",
			path:         []string{filepath.Join(tmpPathParts...)},
			gopath:       filepath.Join(tmpGoParts...),
			targetpath:   filepath.Join(tmpPathParts...),
			symlinksrc:   "",
			symlinkdest:  "",
			expectedpath: filepath.Join(suffixParts...),
		},
		{
			title:        "Symlink points inside GOPATH",
			path:         []string{filepath.Join(tmpPathParts...)},
			gopath:       filepath.Join(tmpGoParts...),
			targetpath:   filepath.Join(tmpSymLinkParts...),
			symlinksrc:   filepath.Join(tmpSymLinkParts...),
			symlinkdest:  filepath.Join(tmpSrcParts...),
			expectedpath: ".",
		},
		{
			title:        "Symlink points inside GOPATH (2)",
			path:         []string{filepath.Join(tmpPathParts...)},
			gopath:       filepath.Join(tmpGoParts...),
			targetpath:   filepath.Join(tmpSymLinkParts...),
			symlinksrc:   filepath.Join(tmpSymLinkParts...),
			symlinkdest:  filepath.Join(tempdir, "root", "go", "src", "github.com"),
			expectedpath: "github.com",
		},
		{
			title: "Symlink point outside GOPATH : Targets Case 1: in baseImport implementation",
			path: []string{
				filepath.Join(tmpPathParts...),
				filepath.Join(tempdir, "root", "gopher", "go"),
			},
			gopath:       filepath.Join(tmpGoParts...),
			targetpath:   filepath.Join(tempdir, "root", "go", "src", "github.com", "gopher"),
			symlinksrc:   filepath.Join(tempdir, "root", "go", "src", "github.com", "gopher"),
			symlinkdest:  filepath.Join(tempdir, "root", "gopher", "go"),
			expectedpath: path.Join("github.com", "gopher"), // with a "/", on every platform
		},
	}
}

func TestCheckPrefixFetchRelPath(t *testing.T) {
	for _, item := range prefixAndFetchRelativePathFixtures() {
		actualok, actualpath := checkPrefixAndFetchRelativePath(item.childpath, item.parentpath)

		item.path = filepath.FromSlash(item.path)

		assert.Equalf(t, item.ok, actualok, "checkPrefixAndFetchRelativePath(%s, %s): expected %v, actual %v", item.childpath, item.parentpath, item.ok, actualok)
		assert.Equal(t, item.path, actualpath, "checkPrefixAndFetchRelativePath(%s, %s): expected %s, actual %s", item.childpath, item.parentpath, item.path, actualpath)
	}
}

func TestBaseImport(t *testing.T) {
	// 1. Create a root folder /tmp/root
	// 2. Simulate scenario
	//	2.a No Symlink
	//	2.b Symlink from outside of GOPATH to inside
	//  2.c Symlink from inside of GOPATH to outside.
	// 3. Check results.

	// NOTE: on windows, this test requires that TempDir and the local target reside on the same drive.
	// This all stems from the use of RelPath to compute the import path, and this should not be really
	// needed (inherited from old GOPATH shinenigans).

	tempdir := t.TempDir()
	golang := GolangOpts()

	for _, item := range baseImportTestFixtures(tempdir) {
		t.Run(fmt.Sprintf("TestBaseImport(%q)", item.title), func(t *testing.T) {
			t.Run("should create paths", func(t *testing.T) {
				for _, paths := range item.path {
					require.NoError(t, os.MkdirAll(paths, 0o700))
				}
			})
			t.Cleanup(func() {
				_ = os.RemoveAll(filepath.Join(tempdir, "root"))
			})

			if item.symlinksrc == "" {
				return
			}

			t.Run("should create Symlink", func(t *testing.T) {
				_, err := os.Stat(item.symlinksrc)
				if os.IsNotExist(err) {
					// specifically for windows, we need to create the entry a symlink points to. linux doesn't care.
					require.NoError(t, os.MkdirAll(item.symlinkdest, fs.ModePerm))
				}

				require.NoErrorf(t, os.Symlink(item.symlinkdest, item.symlinksrc),
					"WARNING:TestBaseImport with symlink could not be carried on. Symlink creation failed for %s -> %s\n%s",
					item.symlinksrc, item.symlinkdest,
					"NOTE:TestBaseImport with symlink on Windows requires extended privileges (admin or a user with SeCreateSymbolicLinkPrivilege)",
				)
			})

			t.Run("baseImport should be "+item.expectedpath, func(t *testing.T) {
				t.Setenv("GOPATH", item.gopath)

				// Test (baseImport always with /)
				actualpath := golang.baseImport(item.targetpath)
				require.Equalf(t, item.expectedpath, actualpath, "baseImport(%s): expected %s, actual %s", item.targetpath, item.expectedpath, actualpath)
			})
		})
	}
}

func TestGenerateMarkdown(t *testing.T) {
	defer discardOutput()()

	t.Run("should generate doc for demo fixture", func(t *testing.T) {
		opts := testGenOpts()
		opts.Spec = "../fixtures/enhancements/184/fixture-184.yaml"
		output := filepath.Join(t.TempDir(), "markdown.md")

		require.NoError(t, GenerateMarkdown(output, nil, nil, opts))
		expectedCode := []string{
			"# Markdown generator demo",
		}

		code, err := os.ReadFile(output)
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", output, line, expectedCode[line])
			}
		}
	})

	t.Run("should handle new lines in descriptions", func(t *testing.T) {
		opts := testGenOpts()
		opts.Spec = "../fixtures/bugs/2700/2700.yaml"
		output := filepath.Join(t.TempDir(), "markdown.md")

		require.NoError(t, GenerateMarkdown(output, nil, nil, opts))
		expectedCode := []string{
			`| Filesystem type of the volume that you want to mount.</br>Tip: Ensure that the filesystem type is supported by the host operating system.</br>Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.</br>More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore</br></br>TODO: how do we prevent errors in the filesystem from compromising the machine |`,
		}

		code, err := os.ReadFile(output)
		require.NoError(t, err)

		for line, codeLine := range expectedCode {
			if !assertInCode(t, strings.TrimSpace(codeLine), string(code)) {
				t.Logf("Code expected did not match in codegenfile %s for expected line %d: %q", output, line, expectedCode[line])
			}
		}
	})
}
