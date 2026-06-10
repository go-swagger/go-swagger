// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"testing"
	"text/template"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-openapi/swag/conv"
	"github.com/go-openapi/swag/mangling"
)

func TestFuncMap(t *testing.T) { //nolint:maintidx // false positive
	t.Parallel()
	fm := testMap()

	const helloTitle = "Hello"

	t.Run("All expected keys should be present", func(t *testing.T) {
		t.Parallel()

		for _, key := range []string{
			"pascalize", "camelize", "humanize", "dasherize",
			"pluralizeFirstWord", "json", "prettyjson",
			"hasInsecure", "hasSecure",
			"dropPackage", "containsPkgStr", "contains",
			"padSurround", "joinFilePath", "joinPath",
			"lineComment", "linePadComment", "blockComment", "inspect",
			"cleanPath", "mediaTypeName", "mediaGoName",
			"dict", "isInteger", "hasPrefix", "stringContains",
			"trimSpace", "mdBlock", "httpStatus",
			"cleanupEnumVariant", "gt0",
			"escapeBackticks",
			"flagNameVar", "flagValueVar", "flagDefaultVar", "flagModelVar", "flagDescriptionVar",
			"printGoLiteral",
		} {
			assert.MapContainsTf(t, fm, key, "expected funcmap key %q", key)
		}
	})

	t.Run("dropPackage should only keep the base name", func(t *testing.T) {
		t.Parallel()

		dropPackage, ok := fm["dropPackage"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, dropPackage)

		assert.EqualT(t, "trail", dropPackage("base.trail"))
		assert.EqualT(t, "trail", dropPackage("base.another.trail"))
		assert.EqualT(t, "trail", dropPackage("trail"))
	})

	t.Run("pascalize should use custom prefix function", func(t *testing.T) {
		t.Parallel()

		pascalize, ok := fm["pascalize"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, pascalize)

		for _, tc := range []struct {
			Input    string
			Expected string
		}{
			{Expected: "Plus1", Input: "+1"},
			{Expected: "Plus", Input: "+"},
			{Expected: "Minus1", Input: "-1"},
			{Expected: "Minus", Input: "-"},
			{Expected: "Nr8", Input: "8"},
			{Expected: "Asterisk", Input: "*"},
			{Expected: "ForwardSlash", Input: "/"},
			{Expected: "EqualSign", Input: "="},
			{Expected: helloTitle, Input: "+hello"},
			// other values from swag rules
			{Expected: "At8", Input: "@8"},
			{Expected: "Bang8", Input: "!8"},
			{Expected: "At", Input: "@"},
			// # values
			{Expected: helloTitle, Input: "#hello"},
			{Expected: "BangHello", Input: "#!hello"},
			{Expected: "HashTag8", Input: "#8"},
			{Expected: "HashTag", Input: "#"},
			// single '_'
			{Expected: "Nr", Input: "_"},
			{Expected: helloTitle, Input: "_hello"},
			// remove spaces
			{Expected: "HelloWorld", Input: "# hello world"},
			{Expected: "HashTag8HelloWorld", Input: "# 8 hello world"},
			{Expected: "Empty", Input: ""},
		} {
			result := pascalize(tc.Input)
			assert.EqualTf(t, tc.Expected, result, "given %q, expected pascalize to yield %q, but got %q", tc.Input, tc.Expected, result)
		}
	})

	t.Run("asJSON should jsonify anything that is serializable to JSON", func(t *testing.T) {
		t.Parallel()

		asJSON, ok := fm["json"].(func(any) (string, error))
		require.TrueT(t, ok)
		require.NotNil(t, pascalize)

		asPrettyJSON, ok := fm["prettyjson"].(func(any) (string, error))
		require.TrueT(t, ok)
		require.NotNil(t, pascalize)

		for _, jsonFunc := range []func(any) (string, error){
			asJSON,
			asPrettyJSON,
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
	})

	t.Run("dict should render values as a map", func(t *testing.T) {
		t.Parallel()

		dict, ok := fm["dict"].(func(...any) (map[string]any, error))
		require.TrueT(t, ok)
		require.NotNil(t, dict)

		d, err := dict("a", "b", "c", "d")
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"a": "b", "c": "d"}, d)

		// odd number of arguments
		_, err = dict("a", "b", "c")
		require.Error(t, err)

		// none-string key
		_, err = dict("a", "b", 3, "d")
		require.Error(t, err)
	})

	t.Run("isInteger should detect integer values", func(t *testing.T) {
		t.Parallel()

		isInteger, ok := fm["isInteger"].(func(any) bool)
		require.TrueT(t, ok)
		require.NotNil(t, isInteger)

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
			conv.Pointer(int(4)),
			conv.Pointer(int32(4)),
			conv.Pointer(int64(4)),
			conv.Pointer(uint(4)),
			conv.Pointer(uint32(4)),
			conv.Pointer(uint64(4)),
			float32(12),
			float64(12),
			conv.Pointer(float32(12)),
			conv.Pointer(float64(12)),
			"12",
			conv.Pointer("12"),
		} {
			val := anInteger
			require.Truef(t, isInteger(val), "expected %#v to be detected an integer value", val)
		}

		for _, notAnInteger := range []any{
			float32(12.5),
			float64(12.5),
			conv.Pointer(float32(12.5)),
			conv.Pointer(float64(12.5)),
			[]string{"a"},
			struct{}{},
			nil,
			map[string]int{"a": 1},
			"abc",
			"2.34",
			conv.Pointer("2.34"),
			nilString,
			nilInt,
			nilFloat,
		} {
			val := notAnInteger
			require.Falsef(t, isInteger(val), "did not expect %#v to be detected an integer value", val)
		}
	})

	t.Run("gt0 should work with any *int64", func(t *testing.T) {
		t.Parallel()

		gt0, ok := fm["gt0"].(func(*int64) bool)
		require.TrueT(t, ok)
		require.NotNil(t, gt0)

		require.TrueT(t, gt0(conv.Pointer(int64(1))))
		require.FalseT(t, gt0(conv.Pointer(int64(0))))
		require.FalseT(t, gt0(nil))
	})

	t.Run("mediaMime should return mime type with parameters stripped", func(t *testing.T) {
		t.Parallel()

		mediaMime, ok := fm["mediaTypeName"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, mediaMime)

		assert.EqualT(t, "application/json", mediaMime("application/json"))
		assert.EqualT(t, "application/json", mediaMime("application/json;param=1;param=2"))
	})

	t.Run("mediaGoMime should return buid a go name from any mime", func(t *testing.T) {
		t.Parallel()

		mediaGoName, ok := fm["mediaGoName"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, mediaGoName)

		assert.EqualT(t, "StarStar", mediaGoName("*/*"))
	})

	t.Run("mediaGoMime should return buid a go name from any mime", func(t *testing.T) {
		t.Parallel()

		containsPkgStr, ok := fm["containsPkgStr"].(func(string) bool)
		require.TrueT(t, ok)
		require.NotNil(t, containsPkgStr)

		assert.TrueT(t, containsPkgStr("models.MyType"))
		assert.FalseT(t, containsPkgStr("MyType"))
		assert.FalseT(t, containsPkgStr(""))
	})

	t.Run("httpStatus should return the string of well-known codes", func(t *testing.T) {
		t.Parallel()

		httpStatus, ok := fm["httpStatus"].(func(int) string)
		require.TrueT(t, ok)
		require.NotNil(t, httpStatus)

		assert.EqualT(t, "OK", httpStatus(200))
		assert.EqualT(t, "Not Found", httpStatus(404))
		assert.EqualT(t, "Status 999", httpStatus(999))
	})

	t.Run("markdownBlock should trim space and handle new line", func(t *testing.T) {
		t.Parallel()

		markdownBlock, ok := fm["mdBlock"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, markdownBlock)

		assert.EqualT(t, "line1</br>line2", markdownBlock("line1\nline2"))
		assert.EqualT(t, "line1</br>line2", markdownBlock("line1\r\nline2"))
		assert.EqualT(t, "trimmed", markdownBlock("  trimmed  "))
	})

	t.Run("pluralizeFirstWord should plurarize a word using inflect", func(t *testing.T) {
		t.Parallel()

		pluralizeFirstWord, ok := fm["pluralizeFirstWord"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, pluralizeFirstWord)

		assert.EqualT(t, "ponies of the round table", pluralizeFirstWord("pony of the round table"))
		assert.EqualT(t, "dwarves", pluralizeFirstWord("dwarf"))
		assert.EqualT(t, "", pluralizeFirstWord(""))
	})

	t.Run("padSurround should plurarize a word using inflect", func(t *testing.T) {
		t.Parallel()

		padSurround, ok := fm["padSurround"].(func(string, string, int, int) string)
		require.TrueT(t, ok)
		require.NotNil(t, padSurround)

		assert.EqualT(t, "-,-,-,padme,-,-,-,-,-,-,-,-", padSurround("padme", "-", 3, 12))
		assert.EqualT(t, "padme,-,-,-,-,-,-,-,-,-,-,-", padSurround("padme", "-", 0, 12))
		assert.EqualT(t, "only", padSurround("only", "-", 0, 1))
	})

	t.Run("cleanupEnumVariant should transliterate special characters", func(t *testing.T) {
		t.Parallel()

		cleanupEnumVariant, ok := fm["cleanupEnumVariant"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, cleanupEnumVariant)

		assert.EqualT(t, "2-Dot-4Ghz", cleanupEnumVariant("2.4Ghz"))
		assert.EqualT(t, "-Plus-1", cleanupEnumVariant("+1"))
		assert.EqualT(t, "a-Dash-b-Hashtag-c", cleanupEnumVariant("a-b#c"))
		assert.EqualT(t, "plain", cleanupEnumVariant("plain"))
		assert.EqualT(t, "-Equal--Equal-", cleanupEnumVariant("=="))
		assert.EqualT(t, "-Equal--Tilde-", cleanupEnumVariant("=~"))
		assert.EqualT(t, "-GreaterThan--Equal-", cleanupEnumVariant(">="))
		assert.EqualT(t, "-LessThan--Equal-", cleanupEnumVariant("<="))
		assert.EqualT(t, "-Bang--Equal-", cleanupEnumVariant("!="))
		assert.EqualT(t, "-Bang--Tilde-", cleanupEnumVariant("!~"))
	})

	t.Run("hasInsecure should detect the http scheme as insecure", func(t *testing.T) {
		t.Parallel()

		hasInsecure, ok := fm["hasInsecure"].(func([]string) bool)
		require.TrueT(t, ok)
		require.NotNil(t, hasInsecure)

		assert.TrueT(t, hasInsecure([]string{"http"}))
		assert.TrueT(t, hasInsecure([]string{"ws"}))
		assert.FalseT(t, hasInsecure([]string{"https"}))
		assert.FalseT(t, hasInsecure([]string{"wss"}))
	})

	t.Run("hasSecure should detect the https scheme as secure", func(t *testing.T) {
		t.Parallel()

		hasSecure, ok := fm["hasSecure"].(func([]string) bool)
		require.TrueT(t, ok)
		require.NotNil(t, hasSecure)

		assert.TrueT(t, hasSecure([]string{"https"}))
		assert.TrueT(t, hasSecure([]string{"wss"}))
		assert.FalseT(t, hasSecure([]string{"http"}))
		assert.FalseT(t, hasSecure([]string{"ws"}))
	})

	t.Run("escapeBackicks should escape backticks in strings", func(t *testing.T) {
		t.Parallel()

		escapeBackticks, ok := fm["escapeBackticks"].(func(string) string)
		require.TrueT(t, ok)
		require.NotNil(t, escapeBackticks)

		assert.EqualT(t, "no ticks", escapeBackticks("no ticks"))
		assert.EqualT(t, "has`+\"`\"+`tick", escapeBackticks("has`tick"))
	})
}

func TestFuncMap_FlagVars(t *testing.T) {
	fm := testMap()
	const (
		flagNameVar        = "flagNameVar"
		flagValueVar       = "flagValueVar"
		flagDefaultVar     = "flagDefaultVar"
		flagModelVar       = "flagModelVar"
		flagDescriptionVar = "flagDescriptionVar"
	)

	for _, tc := range []struct {
		key      string
		expected string
	}{
		{flagNameVar, "flagMyFieldName"},
		{flagValueVar, "flagMyFieldValue"},
		{flagDefaultVar, "flagMyFieldDefault"},
		{flagModelVar, "flagMyFieldModel"},
		{flagDescriptionVar, "flagMyFieldDescription"},
	} {
		fn, ok := fm[tc.key].(func(string) string)
		require.TrueT(t, ok)
		assert.EqualT(t, tc.expected, fn("myField"))
	}
}

func TestFuncMap_PrintGoLiteral(t *testing.T) {
	fm := testMap()

	fn, ok := fm["printGoLiteral"].(func(any) string)
	require.TrueT(t, ok)

	assert.EqualT(t, `"hello"`, fn("hello"))
	assert.EqualT(t, "42", fn(42))
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
	assert.EqualT(t, "-Equal-", replaceSpecialChar('='))
	assert.EqualT(t, "-Bang-", replaceSpecialChar('!'))
	assert.EqualT(t, "-Tilde-", replaceSpecialChar('~'))
	assert.EqualT(t, "-GreaterThan-", replaceSpecialChar('>'))
	assert.EqualT(t, "-LessThan-", replaceSpecialChar('<'))
	assert.EqualT(t, "-Star-", replaceSpecialChar('*'))
	assert.EqualT(t, "-Slash-", replaceSpecialChar('/'))
	assert.EqualT(t, "x", replaceSpecialChar('x'))
}

func testMap() template.FuncMap {
	m := mangling.NewNameMangler(mangling.WithGoNamePrefixFunc(PrefixForName))

	return FuncMap(m)
}
