// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

// TestSanitizeGenHarness is the adversarial verification harness for
// untrusted-spec code injection (see .claude/plans/harden-generated-code-untrusted-spec.md).
//
// Each fixture stuffs a spec field with a Go breakout payload. Generation must
// either REJECT it (identifier overrides — x-go-name, x-go-type) or NEUTRALISE
// it (x-go-custom-tag), so the generated source never gains an injected
// top-level declaration.
//
// Identifier and tag vectors are addressed by Phase 1; comment vectors (summary /
// description) are addressed by the comment-helper unification (lineComment /
// linePadComment), which wraps every line so free text cannot break out. All are
// asserted here.
func TestSanitizeGenHarness(t *testing.T) {
	t.Parallel()
	defer discardOutput()()

	root := t.TempDir()

	t.Run("x-go-custom-tag is contained (model)", func(t *testing.T) {
		target := harnessTarget(t, "custom_tag", "model_test", root)
		opts := defaultServerOpts(t, "../fixtures/codegen/sanitize/x-go-custom-tag.yaml", target)

		require.NoError(t, GenerateModels([]string{"", ""}, opts))
		assertNoInjectedDecl(t, filepath.Join(opts.Target, defaultModelsTarget))
	})

	t.Run("x-go-type is contained (model)", func(t *testing.T) {
		target := harnessTarget(t, "x_go_type", "model_test", root)
		opts := defaultServerOpts(t, "../fixtures/codegen/sanitize/x-go-type.yaml", target)

		// The hostile extension is skipped (warning logged); generation succeeds
		// with the field falling back to a plain string.
		require.NoError(t, GenerateModels([]string{"", ""}, opts))
		assertNoInjectedDecl(t, filepath.Join(opts.Target, defaultModelsTarget))
	})

	t.Run("x-go-name on discriminator is contained (model)", func(t *testing.T) {
		target := harnessTarget(t, "go_name_disc", "model_test", root)
		opts := defaultServerOpts(t, "../fixtures/codegen/sanitize/x-go-name-discriminator.yaml", target)

		// The hostile overrides are skipped (warning logged); the types fall back
		// to their mangled schema names.
		require.NoError(t, GenerateModels([]string{"", ""}, opts))
		assertNoInjectedDecl(t, filepath.Join(opts.Target, defaultModelsTarget))
	})

	t.Run("x-go-name on parameter is rejected (server)", func(t *testing.T) {
		target := harnessTarget(t, "go_name_param", "server_test", root)
		opts := defaultServerOpts(t, "../fixtures/codegen/sanitize/x-go-name-param.yaml", target)

		require.Error(t, GenerateServer("", nil, nil, opts),
			"generation must reject a hostile x-go-name parameter override")
	})

	t.Run("comment injection is contained (server)", func(t *testing.T) {
		target := harnessTarget(t, "comments", "server_test", root)
		opts := defaultServerOpts(t, "../fixtures/codegen/sanitize/comments.yaml", target)

		// info.description and the operation summary carry Go breakout payloads
		// (block "*/", newline, //go: directive). Generation succeeds: the free
		// text is not rejected, but the comment helpers wrap every line so the
		// payload can never become a top-level declaration.
		require.NoError(t, GenerateServer("", nil, nil, opts))
		assertNoInjectedDecl(t, opts.Target)
	})

	t.Run("cli string-literal injection is contained (cli)", func(t *testing.T) {
		target := harnessTarget(t, "cli_injection", "cli_test", root)
		opts := NewGenOpts(ForCli(),
			WithSpec("../fixtures/codegen/sanitize/cli-injection.yaml"), WithTarget(target))

		// A security-definition name/description and an operation summary carry
		// payloads that try to break out of the Go string literals emitted by the
		// CLI generator (flag names, cobra Use/Long, APIKeyAuth args). printf %q and
		// escapeBackticks keep each one inside a single literal.
		require.NoError(t, GenerateClient("", nil, nil, opts))
		assertNoInjectedDecl(t, opts.Target)
	})
}

// assertNoInjectedDecl walks generated .go files under dir and asserts each is
// valid Go with no broken-out payload.
//
// Every sanitize fixture's payload either calls the println builtin — which
// generated go-swagger code never does — or declares an "Injected*" identifier.
// A contained payload survives only inside a comment or a string literal, neither
// of which yields a CallExpr or a bare Ident, so these AST nodes are a reliable,
// false-positive-free signal that a payload escaped its intended slot. (A blanket
// "no func init" rule cannot be used: generated servers legitimately emit one.)
func assertNoInjectedDecl(t *testing.T, dir string) {
	t.Helper()

	fset := token.NewFileSet()
	var files int

	require.NoError(t, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		files++

		f, perr := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		require.NoErrorf(t, perr, "generated file must be valid Go: %s", path)

		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if id, ok := x.Fun.(*ast.Ident); ok && id.Name == "println" {
					t.Errorf("unexpected injected println call in %s (code injection)", path)
				}
			case *ast.Ident:
				if strings.HasPrefix(x.Name, "Injected") {
					t.Errorf("unexpected injected identifier %q in %s (code injection)", x.Name, path)
				}
			}

			return true
		})

		return nil
	}))

	require.Positivef(t, files, "expected generated .go files under %s", dir)
}
