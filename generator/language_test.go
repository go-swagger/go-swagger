package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolang_MangleFileName(t *testing.T) {
	o := &LanguageOpts{}
	o.Init()
	res := o.MangleFileName("aFileEndingInOsNameWindows")
	assert.True(t, strings.HasSuffix(res, "_windows"))

	o = GoLangOpts()
	res = o.MangleFileName("aFileEndingInOsNameWindows")
	assert.True(t, strings.HasSuffix(res, "_windows_swagger"))
	res = o.MangleFileName("aFileEndingInOsNameWindowsAmd64")
	assert.True(t, strings.HasSuffix(res, "_windows_amd64_swagger"))
	res = o.MangleFileName("aFileEndingInTest")
	assert.True(t, strings.HasSuffix(res, "_test_swagger"))
}

func TestGolang_ManglePackage(t *testing.T) {
	o := GoLangOpts()

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
		assert.Equal(t, v.expectedPath, res)
		res = o.ManglePackageName(v.tested, "default")
		assert.Equal(t, v.expectedName, res)
	}
}

// Go literal initializer func
func TestGolang_SliceInitializer(t *testing.T) {
	o := GoLangOpts()
	goSliceInitializer := o.ArrayInitializerFunc

	a0 := []interface{}{"a", "b"}
	res, err := goSliceInitializer(a0)
	assert.NoError(t, err)
	assert.Equal(t, `{"a","b",}`, res)

	a1 := []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}}
	res, err = goSliceInitializer(a1)
	assert.NoError(t, err)
	assert.Equal(t, `{{"a","b",},{"c","d",},}`, res)

	a2 := map[string]interface{}{"a": "y", "b": "z"}
	res, err = goSliceInitializer(a2)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"y","b":"z",}`, res)

	_, err = goSliceInitializer(struct {
		A string `json:"a"`
		B func() string
	}{A: "good", B: func() string { return "" }})
	require.Error(t, err)

	a3 := []interface{}{}
	res, err = goSliceInitializer(a3)
	assert.NoError(t, err)
	assert.Equal(t, `{}`, res)
}

func TestGolangInit(t *testing.T) {
	opts := &LanguageOpts{}
	assert.Equal(t, "", opts.baseImport("x"))
	res, err := opts.FormatContent("x", []byte("y"))
	require.NoError(t, err)
	assert.Equal(t, []byte("y"), res)
	opts = GoLangOpts()
	o := opts
	o.Init()
	assert.Equal(t, opts, o)
}
