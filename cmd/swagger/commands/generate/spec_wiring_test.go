// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/codescan"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/jessevdk/go-flags"
)

// TestSpecFileToOptions pins the CLI-flag → codescan.Options wiring.
//
// It guards against the wiring failure modes: a flag mapped to the wrong Options field, a dropped
// mapping, and — most importantly — the polarity-sensitive cases (EnableAllOfCompounding is
// negated into SkipAllOfCompounding; the DescWithRef flag feeds EmitRefSiblings, not DescWithRef).
func TestSpecFileToOptions(t *testing.T) {
	t.Run("zero value maps to expected defaults", func(t *testing.T) {
		opts := (&SpecFile{}).toOptions(nil, nil, false)

		// polarity: no --enable-allof-compounding ⇒ compounding is skipped.
		assert.True(t, opts.SkipAllOfCompounding)
		// the legacy DescWithRef option is never wired.
		assert.False(t, opts.DescWithRef) //nolint:staticcheck // we precisely assert that there is no wiring of this deprecated field
		assert.False(t, opts.EmitRefSiblings)
	})

	tests := []struct {
		name   string
		set    func(*SpecFile)
		assert func(*testing.T, codescan.Options)
	}{
		{
			name: "DefaultAllOfForEmbeds",
			set:  func(s *SpecFile) { s.DefaultAllOfForEmbeds = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.DefaultAllOfForEmbeds)
			},
		},
		{
			name: "NameFromTags",
			set:  func(s *SpecFile) { s.NameFromTags = []string{"form", "json"} },
			assert: func(t *testing.T, o codescan.Options) {
				assert.Equal(t, []string{"form", "json"}, o.NameFromTags)
			},
		},
		{
			name: "SkipJSONifyInterfaceMethods",
			set:  func(s *SpecFile) { s.SkipJSONifyInterfaceMethods = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.SkipJSONifyInterfaceMethods)
			},
		},
		{
			name: "NameConcatBudget",
			set:  func(s *SpecFile) { s.NameConcatBudget = 0.42 },
			assert: func(t *testing.T, o codescan.Options) {
				assert.Equal(t, 0.42, o.NameConcatBudget)
			},
		},
		{
			name: "AfterDeclComments",
			set:  func(s *SpecFile) { s.AfterDeclComments = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.AfterDeclComments)
			},
		},
		{
			name: "CleanGoDoc",
			set:  func(s *SpecFile) { s.CleanGoDoc = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.CleanGoDoc)
			},
		},
		{
			name: "PruneUnusedModels",
			set:  func(s *SpecFile) { s.PruneUnusedModels = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.PruneUnusedModels)
			},
		},
		{
			name: "EnableAllOfCompounding is negated into SkipAllOfCompounding",
			set:  func(s *SpecFile) { s.EnableAllOfCompounding = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.False(t, o.SkipAllOfCompounding)
			},
		},
		{
			name: "DescWithRef flag feeds EmitRefSiblings, not DescWithRef",
			set:  func(s *SpecFile) { s.DescWithRef = true },
			assert: func(t *testing.T, o codescan.Options) {
				assert.True(t, o.EmitRefSiblings)
				assert.False(t, o.DescWithRef) //nolint:staticcheck // we precisely assert that there is no wiring of this deprecated field
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &SpecFile{}
			tc.set(s)
			tc.assert(t, s.toOptions(nil, nil, false))
		})
	}
}

// generateDefinitions runs `generate spec` end-to-end through Execute and returns the emitted
// definitions map. It exercises the full flag → Options → codescan.Run path.
func generateDefinitions(t *testing.T, s *SpecFile) map[string]any {
	t.Helper()

	out := filepath.Join(t.TempDir(), "spec.json")
	s.Output = flags.Filename(out)
	s.Quiet = true
	require.NoError(t, s.Execute(nil))

	data, err := os.ReadFile(out)
	require.NoError(t, err)

	var doc struct {
		Definitions map[string]any `json:"definitions"`
	}
	require.NoError(t, json.Unmarshal(data, &doc))

	return doc.Definitions
}

// TestSpecFileExecuteRespectsPruneUnusedModels is a behavioral smoke test: the --prune flag must
// actually reach codescan and drop the unreferenced model. Asserting on-vs-off guards against a
// vacuous pass (the fixture must genuinely carry an orphan to prune).
func TestSpecFileExecuteRespectsPruneUnusedModels(t *testing.T) {
	const workDir = "../../../../fixtures/enhancements/prune-unused"

	off := generateDefinitions(t, &SpecFile{WorkDir: workDir, ScanModels: true})
	require.Contains(t, off, "Used")
	require.Contains(t, off, "Orphan", "control: the orphan must be present without --prune, else the test is vacuous")

	on := generateDefinitions(t, &SpecFile{WorkDir: workDir, ScanModels: true, PruneUnusedModels: true})
	assert.Contains(t, on, "Used")
	assert.NotContains(t, on, "Orphan", "--prune must drop the unreferenced model")
}

// TestSpecFileExecuteRespectsDefaultAllOfForEmbeds is a behavioral smoke test: the
// --default-allof-embeds flag must flip a plain embed from inlined properties to an allOf
// composition. Asserting on-vs-off guards against a vacuous pass.
func TestSpecFileExecuteRespectsDefaultAllOfForEmbeds(t *testing.T) {
	const workDir = "../../../../fixtures/enhancements/default-allof-embeds"

	off := generateDefinitions(t, &SpecFile{WorkDir: workDir, ScanModels: true})
	derivedOff, ok := off["Derived"].(map[string]any)
	require.TrueT(t, ok)
	require.NotContains(t, derivedOff, "allOf", "control: default emits inlined properties, not allOf")
	require.Contains(t, derivedOff, "properties")

	on := generateDefinitions(t, &SpecFile{WorkDir: workDir, ScanModels: true, DefaultAllOfForEmbeds: true})
	derivedOn, ok := on["Derived"].(map[string]any)
	require.TrueT(t, ok)
	assert.Contains(t, derivedOn, "allOf", "--default-allof-embeds must compose the embed via allOf")
}
