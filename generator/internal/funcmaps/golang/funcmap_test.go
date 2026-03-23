// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/swag"
)

func TestMain(m *testing.M) {
	swag.GoNamePrefixFunc = PrefixForName //nolint:staticcheck // tracked for migration to mangling.WithGoNamePrefixFunc
	os.Exit(m.Run())
}

func TestFuncMap(t *testing.T) {
	fm := FuncMap()

	// Verify all expected keys are present
	for _, key := range []string{
		"pascalize", "camelize", "humanize", "dasherize",
		"pluralizeFirstWord", "json", "prettyjson",
		"hasInsecure", "hasSecure",
		"dropPackage", "containsPkgStr", "contains",
		"padSurround", "joinFilePath", "joinPath",
		"comment", "blockcomment", "inspect",
		"cleanPath", "mediaTypeName", "mediaGoName",
		"dict", "isInteger", "hasPrefix", "stringContains",
		"trimSpace", "mdBlock", "httpStatus",
		"cleanupEnumVariant", "gt0",
		"escapeBackticks",
		"flagNameVar", "flagValueVar", "flagDefaultVar", "flagModelVar", "flagDescriptionVar",
		"printGoLiteral",
	} {
		assert.NotNil(t, fm[key], "expected funcmap key %q", key)
	}
}

func TestDropPackage(t *testing.T) {
	assert.EqualT(t, "trail", DropPackage("base.trail"))
	assert.EqualT(t, "trail", DropPackage("base.another.trail"))
	assert.EqualT(t, "trail", DropPackage("trail"))
}

func TestPascalize(t *testing.T) {
	assert.EqualT(t, "Plus1", Pascalize("+1"))
	assert.EqualT(t, "Plus", Pascalize("+"))
	assert.EqualT(t, "Minus1", Pascalize("-1"))
	assert.EqualT(t, "Minus", Pascalize("-"))
	assert.EqualT(t, "Nr8", Pascalize("8"))
	assert.EqualT(t, "Asterisk", Pascalize("*"))
	assert.EqualT(t, "ForwardSlash", Pascalize("/"))
	assert.EqualT(t, "EqualSign", Pascalize("="))

	assert.EqualT(t, "Hello", Pascalize("+hello"))

	// other values from swag rules
	assert.EqualT(t, "At8", Pascalize("@8"))
	assert.EqualT(t, "AtHello", Pascalize("@hello"))
	assert.EqualT(t, "Bang8", Pascalize("!8"))
	assert.EqualT(t, "At", Pascalize("@"))

	// # values
	assert.EqualT(t, "Hello", Pascalize("#hello"))
	assert.EqualT(t, "BangHello", Pascalize("#!hello"))
	assert.EqualT(t, "HashTag8", Pascalize("#8"))
	assert.EqualT(t, "HashTag", Pascalize("#"))

	// single '_'
	assert.EqualT(t, "Nr", Pascalize("_"))
	assert.EqualT(t, "Hello", Pascalize("_hello"))

	// remove spaces
	assert.EqualT(t, "HelloWorld", Pascalize("# hello world"))
	assert.EqualT(t, "HashTag8HelloWorld", Pascalize("# 8 hello world"))

	assert.EqualT(t, "Empty", Pascalize(""))
}

func TestAsJSON(t *testing.T) {
	for _, jsonFunc := range []func(any) (string, error){
		AsJSON,
		AsPrettyJSON,
	} {
		res, err := jsonFunc(struct {
			A string `json:"a"`
			B int
		}{A: "good", B: 3})
		require.NoError(t, err)
		assert.JSONEqT(t, `{"a":"good","B":3}`, res)

		_, err = jsonFunc(struct {
			A string `json:"a"`
			B func() string
		}{A: "good", B: func() string { return "" }})
		require.Error(t, err)
	}
}

func TestDict(t *testing.T) {
	d, err := dict("a", "b", "c", "d")
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"a": "b", "c": "d"}, d)

	// odd number of arguments
	_, err = dict("a", "b", "c")
	require.Error(t, err)

	// none-string key
	_, err = dict("a", "b", 3, "d")
	require.Error(t, err)
}

func TestIsInteger(t *testing.T) {
	var (
		nilString *string
		nilInt    *int
		nilFloat  *float32
	)

	for _, anInteger := range []any{
		int8(4),
		int16(4),
		int32(4),
		int64(4),
		int(4),
		swag.Int(4),    //nolint:staticcheck // have to migrate to the new swag API
		swag.Int32(4),  //nolint:staticcheck // have to migrate to the new swag API
		swag.Int64(4),  //nolint:staticcheck // have to migrate to the new swag API
		swag.Uint(4),   //nolint:staticcheck // have to migrate to the new swag API
		swag.Uint32(4), //nolint:staticcheck // have to migrate to the new swag API
		swag.Uint64(4), //nolint:staticcheck // have to migrate to the new swag API
		float32(12),
		float64(12),
		swag.Float32(12), //nolint:staticcheck // have to migrate to the new swag API
		swag.Float64(12), //nolint:staticcheck // have to migrate to the new swag API
		"12",
		swag.String("12"), //nolint:staticcheck // have to migrate to the new swag API
	} {
		val := anInteger
		require.Truef(t, isInteger(val), "expected %#v to be detected an integer value", val)
	}

	for _, notAnInteger := range []any{
		float32(12.5),
		float64(12.5),
		swag.Float32(12.5), //nolint:staticcheck // have to migrate to the new swag API
		swag.Float64(12.5), //nolint:staticcheck // have to migrate to the new swag API
		[]string{"a"},
		struct{}{},
		nil,
		map[string]int{"a": 1},
		"abc",
		"2.34",
		swag.String("2.34"), //nolint:staticcheck // have to migrate to the new swag API
		nilString,
		nilInt,
		nilFloat,
	} {
		val := notAnInteger
		require.Falsef(t, isInteger(val), "did not expect %#v to be detected an integer value", val)
	}
}

func TestGt0(t *testing.T) {
	require.TrueT(t, gt0(swag.Int64(1)))  //nolint:staticcheck // have to migrate to the new swag API
	require.FalseT(t, gt0(swag.Int64(0))) //nolint:staticcheck // have to migrate to the new swag API
	require.FalseT(t, gt0(nil))
}

func TestMediaMime(t *testing.T) {
	assert.EqualT(t, "application/json", MediaMime("application/json"))
	assert.EqualT(t, "application/json", MediaMime("application/json;param=1;param=2"))
}

func TestMediaGoName(t *testing.T) {
	assert.EqualT(t, "StarStar", MediaGoName("*/*"))
}

func TestContainsPkgStr(t *testing.T) {
	assert.TrueT(t, ContainsPkgStr("models.MyType"))
	assert.FalseT(t, ContainsPkgStr("MyType"))
	assert.FalseT(t, ContainsPkgStr(""))
}

func TestPadComment(t *testing.T) {
	assert.EqualT(t, "line1\n// line2\n// line3", padComment("line1\nline2\nline3"))
	assert.EqualT(t, "line1\n//\tline2", padComment("line1\nline2", "\t"))
	assert.EqualT(t, "single", padComment("single"))
}

func TestBlockComment(t *testing.T) {
	assert.EqualT(t, "before [*]/ after", blockComment("before */ after"))
	assert.EqualT(t, "no end marker", blockComment("no end marker"))
}

func TestHttpStatus(t *testing.T) {
	assert.EqualT(t, "OK", httpStatus(200))
	assert.EqualT(t, "Not Found", httpStatus(404))
	assert.EqualT(t, "Status 999", httpStatus(999))
}

func TestMarkdownBlock(t *testing.T) {
	assert.EqualT(t, "line1</br>line2", markdownBlock("line1\nline2"))
	assert.EqualT(t, "line1</br>line2", markdownBlock("line1\r\nline2"))
	assert.EqualT(t, "trimmed", markdownBlock("  trimmed  "))
}

func TestPrefixForName_Letter(t *testing.T) {
	// unicode.IsLetter branch: returns ""
	assert.EqualT(t, "", PrefixForName("hello"))
}

func TestReplaceSpecialChar(t *testing.T) {
	assert.EqualT(t, "-Plus-", replaceSpecialChar('+'))
	assert.EqualT(t, "-Dash-", replaceSpecialChar('-'))
	assert.EqualT(t, "-Hashtag-", replaceSpecialChar('#'))
	assert.EqualT(t, "-Dot-", replaceSpecialChar('.'))
	assert.EqualT(t, "x", replaceSpecialChar('x'))
}

func TestPluralizeFirstWord(t *testing.T) {
	assert.EqualT(t, "ponies of the round table", pluralizeFirstWord("pony of the round table"))
	assert.EqualT(t, "dwarves", pluralizeFirstWord("dwarf"))
	assert.EqualT(t, "", pluralizeFirstWord(""))
}

func TestPadSurround(t *testing.T) {
	assert.EqualT(t, "-,-,-,padme,-,-,-,-,-,-,-,-", padSurround("padme", "-", 3, 12))
	assert.EqualT(t, "padme,-,-,-,-,-,-,-,-,-,-,-", padSurround("padme", "-", 0, 12))
	assert.EqualT(t, "only", padSurround("only", "-", 0, 1))
}

func TestCleanupEnumVariant(t *testing.T) {
	assert.EqualT(t, "2-Dot-4Ghz", cleanupEnumVariant("2.4Ghz"))
	assert.EqualT(t, "-Plus-1", cleanupEnumVariant("+1"))
	assert.EqualT(t, "a-Dash-b-Hashtag-c", cleanupEnumVariant("a-b#c"))
	assert.EqualT(t, "plain", cleanupEnumVariant("plain"))
}

func TestFuncMap_HasInsecure(t *testing.T) {
	fm := FuncMap()
	fn, ok := fm["hasInsecure"].(func([]string) bool)
	require.TrueT(t, ok)

	assert.TrueT(t, fn([]string{"http"}))
	assert.TrueT(t, fn([]string{"ws"}))
	assert.FalseT(t, fn([]string{"https"}))
	assert.FalseT(t, fn([]string{"wss"}))
}

func TestFuncMap_HasSecure(t *testing.T) {
	fm := FuncMap()
	fn, ok := fm["hasSecure"].(func([]string) bool)
	require.TrueT(t, ok)

	assert.TrueT(t, fn([]string{"https"}))
	assert.TrueT(t, fn([]string{"wss"}))
	assert.FalseT(t, fn([]string{"http"}))
	assert.FalseT(t, fn([]string{"ws"}))
}

func TestFuncMap_EscapeBackticks(t *testing.T) {
	fm := FuncMap()
	fn, ok := fm["escapeBackticks"].(func(string) string)
	require.TrueT(t, ok)

	assert.EqualT(t, "no ticks", fn("no ticks"))
	assert.EqualT(t, "has`+\"`\"+`tick", fn("has`tick"))
}

func TestFuncMap_FlagVars(t *testing.T) {
	fm := FuncMap()

	for _, tc := range []struct {
		key      string
		expected string
	}{
		{"flagNameVar", "flagMyFieldName"},
		{"flagValueVar", "flagMyFieldValue"},
		{"flagDefaultVar", "flagMyFieldDefault"},
		{"flagModelVar", "flagMyFieldModel"},
		{"flagDescriptionVar", "flagMyFieldDescription"},
	} {
		fn, ok := fm[tc.key].(func(string) string)
		require.TrueT(t, ok)
		assert.EqualT(t, tc.expected, fn("myField"))
	}
}

func TestFuncMap_PrintGoLiteral(t *testing.T) {
	fm := FuncMap()
	fn, ok := fm["printGoLiteral"].(func(any) string)
	require.TrueT(t, ok)

	assert.EqualT(t, `"hello"`, fn("hello"))
	assert.EqualT(t, "42", fn(42))
}
