// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package language

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestGolang_MangleFileName(t *testing.T) {
	o := &Options{}
	o.Init()
	res := o.MangleFileName("aFileEndingInOsNameWindows")
	assert.TrueT(t, strings.HasSuffix(res, "_windows"))

	o = GolangOpts()
	res = o.MangleFileName("aFileEndingInOsNameWindows")
	assert.TrueT(t, strings.HasSuffix(res, "_windows_swagger"))
	res = o.MangleFileName("aFileEndingInOsNameWindowsAmd64")
	assert.TrueT(t, strings.HasSuffix(res, "_windows_amd64_swagger"))
	res = o.MangleFileName("aFileEndingInTest")
	assert.TrueT(t, strings.HasSuffix(res, "_test_swagger"))
}

func TestGolang_ManglePackage(t *testing.T) {
	o := GolangOpts()

	for _, v := range []struct {
		tested       string
		expectedPath string
		expectedName string
	}{
		{tested: "", expectedPath: "default", expectedName: "default"},
		{tested: "select", expectedPath: "select_default", expectedName: "select_default"},
		{tested: "x", expectedPath: "x", expectedName: "x"},
		{tested: "a/b/c-d/e_f/g", expectedPath: "a/b/c-d/e_f/g", expectedName: "g"},
		{tested: "a/b/c-d/e_f/g-h", expectedPath: "a/b/c-d/e_f/g_h", expectedName: "g_h"},
		{tested: "a/b/c-d/e_f/2A", expectedPath: "a/b/c-d/e_f/nr2_a", expectedName: "nr2_a"},
		{tested: "a/b/c-d/e_f/#", expectedPath: "a/b/c-d/e_f/hash_tag", expectedName: "hash_tag"},
		{tested: "#help", expectedPath: "hash_tag_help", expectedName: "hash_tag_help"},
		{tested: "vendor", expectedPath: "vendor_swagger", expectedName: "vendor_swagger"},
		{tested: "internal", expectedPath: "internal_swagger", expectedName: "internal_swagger"},
	} {
		res := o.ManglePackagePath(v.tested, "default")
		assert.EqualT(t, v.expectedPath, res)
		res = o.ManglePackageName(v.tested, "default")
		assert.EqualT(t, v.expectedName, res)
	}
}

// Go literal initializer func.
func TestGolang_SliceInitializer(t *testing.T) {
	o := GolangOpts()
	goSliceInitializer := o.ArrayInitializerFunc

	a0 := []any{"a", "b"}
	res, err := goSliceInitializer(a0)
	require.NoError(t, err)
	assert.EqualT(t, `{"a","b",}`, res)

	a1 := []any{[]any{"a", "b"}, []any{"c", "d"}}
	res, err = goSliceInitializer(a1)
	require.NoError(t, err)
	assert.EqualT(t, `{{"a","b",},{"c","d",},}`, res)

	a2 := map[string]any{"a": "y", "b": "z"}
	res, err = goSliceInitializer(a2)
	require.NoError(t, err)
	assert.EqualT(t, `{"a":"y","b":"z",}`, res)

	_, err = goSliceInitializer(struct {
		A string `json:"a"`
		B func() string
	}{A: "good", B: func() string { return "" }})
	require.Error(t, err)

	a3 := []any{}
	res, err = goSliceInitializer(a3)
	require.NoError(t, err)
	assert.EqualT(t, `{}`, res)
}

func TestGolang_Imports(t *testing.T) {
	o := GolangOpts()

	// empty map: returns ""
	assert.Empty(t, o.Imports(map[string]string{}))

	// unaliased import (name matches last path component)
	res := o.Imports(map[string]string{"fmt": "fmt"})
	assert.StringContainsT(t, res, `"fmt"`)

	// aliased import (name differs from last path component)
	res = o.Imports(map[string]string{"myalias": "github.com/example/pkg"})
	assert.StringContainsT(t, res, `myalias "github.com/example/pkg"`)
}

func TestDefaultGoFormatFunc(t *testing.T) {
	o := GolangOpts()

	src := []byte("package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hello\") }\n")
	res, err := o.FormatContent("test.go", src)
	require.NoError(t, err)
	assert.StringContainsT(t, string(res), "package main")
	assert.StringContainsT(t, string(res), `"fmt"`)
}

func TestRelPathToRelGoPath(t *testing.T) {
	assert.EqualT(t, "", relPathToRelGoPath("/base", "."))
	assert.EqualT(t, "/sub/pkg", relPathToRelGoPath("/base", "/base/sub/pkg"))
	assert.EqualT(t, "/pkg", relPathToRelGoPath("/base", "/base/pkg"))
}

func TestCheckPrefixAndFetchRelativePath(t *testing.T) {
	ok, rel := CheckPrefixAndFetchRelativePath("/home/user/go/src/mypackage", "/home/user/go/src")
	assert.TrueT(t, ok)
	assert.EqualT(t, "mypackage", rel)

	ok, _ = CheckPrefixAndFetchRelativePath("/other/path", "/home/user/go/src")
	assert.FalseT(t, ok)
}
