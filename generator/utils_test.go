// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"os"
	"regexp"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"
)

var (
	discardOutput = gentest.DiscardOutput
	captureOutput = gentest.CaptureOutput
)

// testing utilities for codegen assertions

func reqm(str string) *regexp.Regexp {
	return regexp.MustCompile(regexp.QuoteMeta(str))
}

func reqOri(str string) *regexp.Regexp {
	return regexp.MustCompile(str)
}

func assertInCode(tb testing.TB, expr, code string) bool {
	tb.Helper()
	return assert.RegexpT(tb, reqm(expr), code)
}

func assertRegexpInCode(tb testing.TB, expr, code string) bool {
	tb.Helper()
	return assert.RegexpT(tb, reqOri(expr), code)
}

func assertNotInCode(tb testing.TB, expr, code string) bool {
	tb.Helper()
	return assert.NotRegexpT(tb, reqm(expr), code)
}

func assertRegexpNotInCode(tb testing.TB, expr, code string) bool {
	tb.Helper()
	return assert.NotRegexpT(tb, reqOri(expr), code)
}

func requireValidation(tb testing.TB, pth, expr string, gm GenSchema) {
	if !assertValidation(tb, pth, expr, gm) {
		tb.FailNow()
	}
}

func assertValidation(tb testing.TB, pth, expr string, gm GenSchema) bool {
	tb.Helper()
	if !assert.TrueT(tb, gm.HasValidations, "expected the schema to have validations") {
		return false
	}
	if !assert.EqualT(tb, pth, gm.Path, "paths don't match") {
		return false
	}
	if !assert.EqualT(tb, expr, gm.ValueExpression, "expressions don't match") {
		return false
	}
	return true
}

func funcBody(code string, signature string) string {
	submatches := regexp.MustCompile(
		"\\nfunc \\([a-zA-Z_][a-zA-Z0-9_]* " + regexp.QuoteMeta(signature) + " {\\n" + // function signature
			"((([^}\\n][^\\n]*)?\\n)*)}\\n", // function body
	).FindStringSubmatch(code)

	if submatches == nil {
		return ""
	}
	return submatches[1]
}

// testing utilities for codegen build

func testCwd(tb testing.TB) string {
	tb.Helper()
	cwd, err := os.Getwd()
	require.NoError(tb, err)
	return cwd
}
