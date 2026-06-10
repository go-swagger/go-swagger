// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"strconv"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestIsGoIdentifier(t *testing.T) {
	valid := []string{"Foo", "foo", "_foo", "fooBar", "Föö", "ID", "f1"}
	for _, s := range valid {
		assert.Truef(t, isGoIdentifier(s), "expected %q to be a valid Go identifier", s)
	}

	invalid := []string{
		"",                      // empty
		"_",                     // blank identifier
		"1foo",                  // leading digit
		"foo bar",               // space
		"foo.bar",               // dotted
		"foo-bar",               // dash
		"func",                  // keyword
		"type",                  // keyword
		"foo\nbar",              // newline
		"foo`bar",               // backtick
		"}; func init(){}; var", // injection payload
	}
	for _, s := range invalid {
		assert.Falsef(t, isGoIdentifier(s), "expected %q to be rejected", s)
	}
}

func TestIsGoQualifiedType(t *testing.T) {
	valid := []string{"MyType", "pkg.MyType", "strfmt.DateTime"}
	for _, s := range valid {
		assert.Truef(t, isGoQualifiedType(s), "expected %q to be a valid qualified type", s)
	}

	invalid := []string{
		"",
		".",
		"pkg.",
		".Type",
		"pkg..Type",
		"[]byte",                        // composite is not a bare qualified ident
		"map[string]int",                // composite
		"Type } ; func init(){} ; type", // injection payload
		"Type\n}\nfunc init(){}",        // newline breakout
	}
	for _, s := range invalid {
		assert.Falsef(t, isGoQualifiedType(s), "expected %q to be rejected", s)
	}
}

func TestValidateGoExtensions(t *testing.T) {
	require.NoError(t, validateGoIdentifierExtension("MyName"))
	require.Error(t, validateGoIdentifierExtension("My Name"))
	require.Error(t, validateGoIdentifierExtension("`; func init(){}"))

	require.NoError(t, validateGoTypeExtension("pkg.MyType"))
	require.Error(t, validateGoTypeExtension("Type }; func init(){}"))
}

func TestSanitizeGoNameOverride(t *testing.T) {
	defer discardOutput()()

	assert.Equal(t, "", sanitizeGoNameOverride(""))
	assert.Equal(t, "MyName", sanitizeGoNameOverride("MyName"))
	assert.Equal(t, "", sanitizeGoNameOverride("My Name"))
	assert.Equal(t, "", sanitizeGoNameOverride("`; func init(){}"))
}

// TestPrintTagsCustomTagContained verifies that an x-go-custom-tag value cannot
// break out of the generated struct-tag literal, whichever literal form is used.
//
// The whole tag (payload included) must render as a single, valid Go string
// literal. A backtick is the only terminator of a raw `...` literal, and
// strconv.Quote escapes everything (including ") in the double-quoted fallback,
// so neither a backtick nor an inner double quote can escape.
func TestPrintTagsCustomTagContained(t *testing.T) {
	for _, payload := range []string{
		`db:"a,b"`, // benign custom tag with inner double quotes
		"`; }; func init(){ println(\"PWNED\") }; type _ struct { _ string `", // backtick breakout attempt
		`"; }; func init(){}; var _ = "`,                                      // double-quote breakout attempt (no backtick)
		"db:\"x\" `nested` \"y\"",                                             // mix of backticks and double quotes
	} {
		s := GenSchema{OriginalName: "field"}
		s.CustomTag = payload

		tag := s.PrintTags()

		// Must parse as exactly one Go string literal (raw or interpreted): the
		// payload cannot have terminated the literal early.
		unq, err := strconv.Unquote(tag)
		require.NoErrorf(t, err, "PrintTags must emit a single valid Go string literal for %q, got: %s", payload, tag)

		// The payload is retained verbatim inside that single literal.
		assert.Containsf(t, unq, payload, "custom tag payload must be contained inside the literal: %s", tag)

		// A backtick in the payload forces the double-quoted form (a raw literal
		// cannot contain a backtick).
		if strings.Contains(payload, "`") {
			assert.Falsef(t, strings.HasPrefix(tag, "`"),
				"a custom tag containing a backtick must not render as a raw-backtick literal: %s", tag)
		}
	}
}
