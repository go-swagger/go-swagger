// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// TestPrepare_EquivalentToLegacySequence asserts that the single Prepare() call
// produces the same finalized state as the historical
// EnsureDefaults -> CheckOpts -> loadTemplates sequence.
func TestPrepare_EquivalentToLegacySequence(t *testing.T) {
	defer discardOutput()()

	spec := filepath.Join("..", "testdata", "codegen", "simplesearch.yml")

	mk := func() *GenOpts {
		g := &GenOpts{}
		g.Spec = spec
		g.Target = "."
		g.APIPackage = defaultAPIPackage
		g.ModelPackage = defaultModelPackage
		g.ServerPackage = defaultServerPackage
		g.ClientPackage = defaultClientPackage
		g.IncludeModel = true
		g.IncludeHandler = true
		g.IncludeParameters = true
		g.IncludeResponses = true
		g.IncludeSupport = true
		return g
	}

	// historical sequence, as the generator entry points still drive it
	legacy := mk()
	require.NoError(t, ensureMachinery(legacy))
	require.NoError(t, validateOpts(legacy))
	require.NoError(t, legacy.loadTemplates())

	// new single, idempotent entry point
	prep := mk()
	require.NoError(t, prep.Prepare())

	// machinery is built
	require.NotNil(t, prep.LanguageOpts)
	require.NotNil(t, prep.funcMap)
	require.NotNil(t, prep.templates)
	require.NotNil(t, prep.FlattenOpts)

	// spec normalized identically (absolute path)
	assert.TrueT(t, filepath.IsAbs(prep.Spec))
	assert.EqualT(t, legacy.Spec, prep.Spec)

	// scalar defaults match
	assert.EqualT(t, legacy.DefaultScheme, prep.DefaultScheme)
	assert.EqualT(t, legacy.DefaultConsumes, prep.DefaultConsumes)
	assert.EqualT(t, legacy.DefaultProduces, prep.DefaultProduces)
	assert.EqualT(t, legacy.Principal, prep.Principal)
	assert.EqualT(t, legacy.IncludeValidator, prep.IncludeValidator)

	// render plan (sections) match
	assert.Len(t, prep.Sections.Application, len(legacy.Sections.Application))
	assert.Len(t, prep.Sections.Models, len(legacy.Sections.Models))
	assert.Len(t, prep.Sections.Operations, len(legacy.Sections.Operations))

	// func map populated equivalently
	assert.Len(t, prep.funcMap, len(legacy.funcMap))

	// Prepare is idempotent
	require.NoError(t, prep.Prepare())
}

// TestPrepare_ConfigLayoutOverridesDefaults asserts that a partial config
// `layout:` overrides only the sections it specifies, while the unspecified
// sections keep their defaults (rather than being wiped, as the historical
// wholesale-replace bug did).
func TestPrepare_ConfigLayoutOverridesDefaults(t *testing.T) {
	defer discardOutput()()

	const partialLayout = `
layout:
  models:
    - name: custom-model
      source: asset:model
      target: "{{ joinFilePath .Target (toPackagePath .ModelPackage) }}"
      file_name: "{{ (snakize (pascalize .Name)) }}.go"
`
	cfg := viper.New()
	cfg.SetConfigType("yaml")
	require.NoError(t, cfg.ReadConfig(strings.NewReader(partialLayout)))

	g := &GenOpts{}
	g.Spec = filepath.Join("..", "testdata", "codegen", "simplesearch.yml")
	g.Target = "."
	g.ServerPackage = defaultServerPackage
	g.IncludeHandler = true
	g.IncludeParameters = true
	g.IncludeResponses = true
	g.Viper = cfg

	require.NoError(t, g.Prepare())

	// the models section is overridden by the config layout
	require.Len(t, g.Sections.Models, 1)
	assert.EqualT(t, "custom-model", g.Sections.Models[0].Name)

	// the sections NOT mentioned in the config keep their defaults (not wiped)
	assert.NotEmpty(t, g.Sections.Operations)
	assert.NotEmpty(t, g.Sections.Application)
}

// TestPrepare_ValidationFailsBeforeMutation asserts that a validation failure is
// reported without leaving the options half-built.
func TestPrepare_ValidationFailsBeforeMutation(t *testing.T) {
	defer discardOutput()()

	g := &GenOpts{}
	// an absolute --server-package fails validation; use filepath.Abs so the path
	// is absolute on every platform (a leading separator is not absolute on Windows).
	abs, err := filepath.Abs(filepath.Join("absolute", "server", "pkg"))
	require.NoError(t, err)
	g.ServerPackage = abs

	require.Error(t, g.Prepare())

	// nothing finalized: machinery not built, options not marked prepared
	assert.FalseT(t, g.prepared)
	assert.FalseT(t, g.machineryBuilt)
}

// TestEnsureTarget asserts the pre-flight check on the generation target: the target
// must be a directory this process may write to, and is only created when asked for.
func TestEnsureTarget(t *testing.T) {
	defer discardOutput()()

	t.Run("should accept an existing writable directory", func(t *testing.T) {
		g := &GenOpts{Target: t.TempDir()}

		require.NoError(t, g.ensureTarget())
		assert.TrueT(t, g.targetEnsured)
	})

	t.Run("should leave no probe file behind", func(t *testing.T) {
		target := t.TempDir()
		g := &GenOpts{Target: target}

		require.NoError(t, g.ensureTarget())

		entries, err := os.ReadDir(target)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("should default an empty target to the current directory", func(t *testing.T) {
		// dump-data mode, so the check stops at the normalization: no probe file is
		// written in the package directory.
		g := &GenOpts{Target: "", DumpData: true}

		require.NoError(t, g.ensureTarget())
		assert.EqualT(t, ".", g.Target)
	})

	t.Run("with a missing target", func(t *testing.T) {
		t.Run("should fail and point at the flag", func(t *testing.T) {
			target := filepath.Join(t.TempDir(), "missing")
			g := &GenOpts{Target: target}

			err := g.ensureTarget()
			require.ErrorContains(t, err, "--ensure-target")
			assertNoDir(t, target)
			assert.FalseT(t, g.targetEnsured)
		})

		t.Run("should create it, with its parents, when EnsureTarget is set", func(t *testing.T) {
			target := filepath.Join(t.TempDir(), "a", "b", "c")
			g := &GenOpts{Target: target, EnsureTarget: true}

			require.NoError(t, g.ensureTarget())
			assert.DirExists(t, target)
		})
	})

	t.Run("should fail when the target is not a directory", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "afile")
		require.NoError(t, os.WriteFile(target, []byte("x"), readAllFile))

		// EnsureTarget doesn't help here: the path exists, it is just not a directory.
		g := &GenOpts{Target: target, EnsureTarget: true}

		require.ErrorContains(t, g.ensureTarget(), "is not a directory")
	})

	t.Run("should fail when the target is not writable", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("file modes don't govern write access on windows")
		}
		if os.Geteuid() == 0 {
			t.Skip("root writes to a read-only directory")
		}

		target := filepath.Join(t.TempDir(), "readonly")
		const readOnlyDir = 0o500
		require.NoError(t, os.Mkdir(target, readOnlyDir))
		t.Cleanup(func() {
			_ = os.Chmod(target, readableDir) // so the temp dir may be cleaned up
		})

		g := &GenOpts{Target: target}

		require.ErrorContains(t, g.ensureTarget(), "is not writeable")
	})

	t.Run("should skip the check when dumping data", func(t *testing.T) {
		// nothing is written to the target in that mode
		target := filepath.Join(t.TempDir(), "missing")
		g := &GenOpts{Target: target, DumpData: true}

		require.NoError(t, g.ensureTarget())
		assertNoDir(t, target)
	})

	t.Run("should check the target only once", func(t *testing.T) {
		target := t.TempDir()
		g := &GenOpts{Target: target}
		require.NoError(t, g.ensureTarget())

		// the guard holds even though the target is now gone
		require.NoError(t, os.RemoveAll(target))
		require.NoError(t, g.ensureTarget())
	})
}

// TestPrepare_TargetCheckedAfterSpec asserts that Prepare resolves the spec before it
// touches the target, so a run that fails early leaves no target tree behind.
func TestPrepare_TargetCheckedAfterSpec(t *testing.T) {
	defer discardOutput()()

	target := filepath.Join(t.TempDir(), "target")
	g := &GenOpts{
		Spec:         filepath.Join("..", "testdata", "codegen", "nosuchspec.yml"),
		Target:       target,
		EnsureTarget: true,
	}

	require.Error(t, g.Prepare())
	assertNoDir(t, target)
}

// assertNoDir asserts that nothing exists at path.
func assertNoDir(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)
	assert.TrueTf(t, errors.Is(err, fs.ErrNotExist), "expected %q not to exist", path)
}
