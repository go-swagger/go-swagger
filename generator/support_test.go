package generator

import (
	"testing"
	"os"
	goruntime "runtime"
	"path/filepath"
)

var checkprefixandfetchrelativepathtests = []struct {
	childpath  string
	parentpath string
	ok bool
	path string
}{
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

var checkbaseimporttest = []struct {
	path  []string
	gopath string
	targetpath string
	symlinksrc string
	symlinkdest string  // symlink is the last dir in targetpath
	expectedpath string
}{
	// No sym link. Positive Test Case
	{[]string {"/tmp/root/go/src/github.com/go-swagger",}, "/tmp/root/go/", "/tmp/root/go/src/github.com/go-swagger", "", "", "github.com/go-swagger"},
	// Symlink points inside GOPATH
	{[]string {"/tmp/root/go/src/github.com/go-swagger",}, "/tmp/root/go/", "/tmp/root/symlink", "/tmp/root/symlink", "/tmp/root/go/src/", "."},
	// Symlink points inside GOPATH
	{[]string {"/tmp/root/go/src/github.com/go-swagger",}, "/tmp/root/go/", "/tmp/root/symlink", "/tmp/root/symlink", "/tmp/root/go/src/github.com", "github.com"},
	// Symlink point outside GOPATH : Targets Case 1: in baseImport implementation.
	{[]string {"/tmp/root/go/src/github.com/go-swagger","/tmp/root/gopher/go/"}, "/tmp/root/go/", "/tmp/root/go/src/github.com/gopher", "/tmp/root/go/src/github.com/gopher", "/tmp/root/gopher/go", "github.com/gopher"},

}

func TestCheckPrefixFetchRelPath(t *testing.T) {

	for _,item := range checkprefixandfetchrelativepathtests {
		actualok, actualpath := checkPrefixAndFetchRelativePath(item.childpath, item.parentpath)

		if goruntime.GOOS == "windows" {
			item.path = filepath.Clean(item.path)
		}

		if actualok != item.ok {
			t.Errorf("checkPrefixAndFetchRelativePath(%s, %s): expected %v, actual %v", item.childpath, item.parentpath, item.ok, actualok)
		} else if actualpath != item.path {
			t.Errorf("checkPrefixAndFetchRelativePath(%s, %s): expected %s, actual %s", item.childpath, item.parentpath, item.path, actualpath)
		} else {
			continue
		}
	}

}

func TestBaseImport(t *testing.T) {

	// 1. Create a root folder /tmp/root
	// 2. Simulate scenario
	//	2.a No Symlink
	//	2.b Symlink from outside of GOPATH to inside
	//  2.c Symlink from inside of GOPATH to outside.
	// 3. Check results.

	oldgopath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldgopath)
	defer os.RemoveAll("/tmp/root")

	for _,item := range checkbaseimporttest {

		// Create Paths
		for _,paths := range item.path {
			os.MkdirAll(paths, 0777)
		}

		// Change GOPATH
		os.Setenv("GOPATH", item.gopath)

		// Create Symlink
		os.Symlink(item.symlinkdest, item.symlinksrc)

		// Test
		actualpath := baseImport(item.targetpath)

		if goruntime.GOOS == "windows" {
			actualpath = filepath.Clean(actualpath)
			item.expectedpath = filepath.Clean(item.expectedpath)
		}

		if actualpath != item.expectedpath {
			t.Errorf("baseImport(%s): expected %s, actual %s", item.targetpath, item.expectedpath, actualpath)
		}

		os.RemoveAll("/tmp/root")

	}

}
