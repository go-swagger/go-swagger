// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestNewGenOpts_Presets(t *testing.T) {
	defer discardOutput()()

	t.Run("ForServer applies the standard server bundle", func(t *testing.T) {
		g := NewGenOpts(ForServer(),
			WithSpec(filepath.Join("..", "fixtures", "codegen", "simplesearch.yml")),
			WithTarget("."),
		)

		assert.FalseT(t, g.IsClient)
		assert.TrueT(t, g.IncludeModel)
		assert.TrueT(t, g.IncludeHandler)
		assert.TrueT(t, g.IncludeParameters)
		assert.TrueT(t, g.IncludeResponses)
		assert.TrueT(t, g.IncludeSupport)
		assert.EqualT(t, defaultServerTarget, g.ServerPackage)
		assert.EqualT(t, defaultModelsTarget, g.ModelPackage)

		// options-only construction builds no machinery
		assert.Nil(t, g.LanguageOpts)

		// and Prepare finalizes it
		require.NoError(t, g.Prepare())
		require.NotNil(t, g.LanguageOpts)
		assert.NotEmpty(t, g.Sections.Application)
	})

	t.Run("ForClient sets the client flag", func(t *testing.T) {
		g := NewGenOpts(ForClient())
		assert.TrueT(t, g.IsClient)
		assert.EqualT(t, defaultClientTarget, g.ClientPackage)
	})

	t.Run("With* setters populate the fields", func(t *testing.T) {
		g := NewGenOpts(WithSpec("s.yml"), WithTarget("/tmp/out"), WithTemplatePlugin("p.so"))
		assert.EqualT(t, "s.yml", g.Spec)
		assert.EqualT(t, "/tmp/out", g.Target)
		assert.EqualT(t, "p.so", g.TemplatePlugin)
	})
}
