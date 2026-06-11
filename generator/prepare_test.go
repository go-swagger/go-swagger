// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"path/filepath"
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

	spec := filepath.Join("..", "fixtures", "codegen", "simplesearch.yml")

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
	g.Spec = filepath.Join("..", "fixtures", "codegen", "simplesearch.yml")
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
