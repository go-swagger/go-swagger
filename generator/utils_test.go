package generator

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/generator/internal/gentest"
)

var (
	discardOutput = gentest.DiscardOutput
	captureOutput = gentest.CaptureOutput
	goExecInDir   = gentest.GoExecInDir
)

// testing utilities for codegen assertions

func reqm(str string) *regexp.Regexp {
	return regexp.MustCompile(regexp.QuoteMeta(str))
}

func reqOri(str string) *regexp.Regexp {
	return regexp.MustCompile(str)
}

func assertInCode(t testing.TB, expr, code string) bool {
	return assert.Regexp(t, reqm(expr), code)
}

func assertRegexpInCode(t testing.TB, expr, code string) bool {
	return assert.Regexp(t, reqOri(expr), code)
}

func assertNotInCode(t testing.TB, expr, code string) bool {
	return assert.NotRegexp(t, reqm(expr), code)
}

func assertRegexpNotInCode(t testing.TB, expr, code string) bool {
	return assert.NotRegexp(t, reqOri(expr), code)
}

// Unused
// func assertRegexpNotInCode(t testing.TB, expr, code string) bool {
// 	return assert.NotRegexp(t, reqOri(expr), code)
// }

func requireValidation(t testing.TB, pth, expr string, gm GenSchema) {
	if !assertValidation(t, pth, expr, gm) {
		t.FailNow()
	}
}

func assertValidation(t testing.TB, pth, expr string, gm GenSchema) bool {
	if !assert.True(t, gm.HasValidations, "expected the schema to have validations") {
		return false
	}
	if !assert.Equal(t, pth, gm.Path, "paths don't match") {
		return false
	}
	if !assert.Equal(t, expr, gm.ValueExpression, "expressions don't match") {
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

func testCwd(t testing.TB) string {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	return cwd
}
