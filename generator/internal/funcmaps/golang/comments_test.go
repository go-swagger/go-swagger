// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package golang

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func ptrTo[T any](v T) *T { return &v }

func TestLineComment(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single", "a summary", "// a summary"},
		{"crlf normalized", "a\r\nb", "// a\n// b"},
		{"lone cr normalized", "a\rb", "// a\n// b"},
		{"blank inner line", "a\n\nb", "// a\n//\n// b"},
		// trailing blank lines are preserved (the formatter drops one that ends up
		// trailing the whole comment block); per-line trailing spaces are trimmed.
		{"single trailing newline preserved", "a\n", "// a\n//"},
		{"trailing blank lines preserved", "a\n\n", "// a\n//\n//"},
		{"all blank yields nothing", "\n\n", ""},
		{"trailing spaces trimmed", "a   \nb\t", "// a\n// b"},
		// directive safety: the guaranteed space after // means a leading
		// "go:" / "line" can never form a compiler directive.
		{"go directive defused", "go:embed secret", "// go:embed secret"},
		{"line directive defused", "line evil.go:1", "// line evil.go:1"},
		// newline breakout payload becomes inert comment lines.
		{"breakout becomes comment", "x\nfunc init(){}", "// x\n// func init(){}"},
	}
	for _, tc := range cases {
		assert.Equalf(t, tc.want, lineComment(tc.in), "lineComment(%q)", tc.in)
	}
}

// TestLineCommentVariadic covers the multi-argument form used to compose a line
// inline, e.g. {{ lineComment "MinProperties: " .MinProperties }}.
func TestLineCommentVariadic(t *testing.T) {
	cases := []struct {
		name string
		in   []any
		want string
	}{
		{"no args", nil, ""},
		{"string prefix with value", []any{"MinProperties: ", 5}, "// MinProperties: 5"},
		{"adjacent strings not spaced", []any{"a", "b"}, "// ab"},
		{"non-string operands spaced", []any{1, 2}, "// 1 2"},
		{"nil arg skipped", []any{"x", nil, "y"}, "// xy"},
		{"all nil yields nothing", []any{nil, nil}, ""},
		// pointer args are dereferenced (as text/template does for {{ .X }}), so
		// they print by value, not as an address.
		{"pointer deref to value", []any{"MinProperties: ", ptrTo(int64(20))}, "// MinProperties: 20"},
		{"nil pointer skipped", []any{"v: ", (*int64)(nil)}, "// v:"},
		{"empty strings yield nothing", []any{"", ""}, ""},
		// composed value carrying a newline stays fully commented.
		{"composed multiline", []any{"Pattern: ", "a\nfunc init(){}"}, "// Pattern: a\n// func init(){}"},
	}
	for _, tc := range cases {
		assert.Equalf(t, tc.want, lineComment(tc.in...), "lineComment(%v)", tc.in)
	}
}

// TestLinePadComment covers the indentation-preserving variant used for the
// swagger:meta doc block, where continuation lines must keep a fixed indent.
func TestLinePadComment(t *testing.T) {
	cases := []struct {
		name string
		pad  string
		in   []any
		want string
	}{
		{"empty", "  ", []any{""}, ""},
		{"single line two-space pad", "  ", []any{"Host: ", "example.com"}, "//  Host: example.com"},
		{"four-space pad", "    ", []any{"- ", "application/json"}, "//    - application/json"},
		// every continuation line keeps the pad (unlike lineComment's fixed "// ").
		{"multiline keeps pad", "  ", []any{"a\nb"}, "//  a\n//  b"},
		{"blank inner line collapses", "  ", []any{"a\n\nb"}, "//  a\n//\n//  b"},
		// directive safety: a pad lacking leading whitespace gets a space inserted.
		{"pad without space made safe", "x", []any{"go:embed evil"}, "// xgo:embed evil"},
		{"empty pad falls back to single space", "", []any{"text"}, "// text"},
		{"nil arg skipped", "  ", []any{"k: ", nil, "v"}, "//  k: v"},
	}
	for _, tc := range cases {
		assert.Equalf(t, tc.want, linePadComment(tc.pad, tc.in...), "linePadComment(%q, %v)", tc.pad, tc.in)
	}
}

func TestWrapBlockComment(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"single", "a summary", "/* a summary */"},
		{"multiline", "a\nb", "/*\na\nb\n*/"},
		{"crlf normalized", "a\r\nb", "/*\na\nb\n*/"},
		{"trailing trimmed", "a summary  ", "/* a summary */"},
		// terminator neutralized so the text cannot escape the comment.
		{"terminator neutralized", "safe */ func init(){}", "/* safe [*]/ func init(){} */"},
		{"terminator neutralized multiline", "a\n*/\nb", "/*\na\n[*]/\nb\n*/"},
	}
	for _, tc := range cases {
		assert.Equalf(t, tc.want, wrapBlockComment(tc.in), "wrapBlockComment(%q)", tc.in)
	}
}

// TestCommentHelpersRegistered ensures the helpers are wired into the funcmap
// under their template-facing names.
func TestCommentHelpersRegistered(t *testing.T) {
	fm := testMap()
	assert.Contains(t, fm, "lineComment")
	assert.Contains(t, fm, "linePadComment")
	assert.Contains(t, fm, "blockComment")
}
