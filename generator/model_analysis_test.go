// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/testify/v2/require"
)

// TestModelAnalysis_SharedIsEquivalentToPerModelAnalysis checks that the shared analysis yields the
// same definitions as analyzing the spec afresh for every model.
func TestModelAnalysis_SharedIsEquivalentToPerModelAnalysis(t *testing.T) {
	defer discardOutput()()

	const specPath = "../testdata/codegen/todolist.discriminators.yml"

	opts := opts()

	sharedDoc, err := loads.Spec(specPath)
	require.NoError(t, err)
	sharedAnalysis := newModelAnalysis(sharedDoc, analysis.New(sharedDoc.Spec()), opts)

	names := make([]string, 0, len(sharedDoc.Spec().Definitions))
	for name := range sharedDoc.Spec().Definitions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			perModelDoc, err := loads.Spec(specPath)
			require.NoError(t, err)

			perModel, err := makeGenDefinitionWithAnalysis(name, "models", perModelDoc.Spec().Definitions[name], perModelDoc, opts, nil)
			require.NoError(t, err)

			shared, err := makeGenDefinitionWithAnalysis(name, "models", sharedDoc.Spec().Definitions[name], sharedDoc, opts, sharedAnalysis)
			require.NoError(t, err)

			require.Equal(t, perModel, shared)
		})
	}
}

// TestModelAnalysis_IsSnapshotNotAnalysis checks that model generation is not affected by spec mutations:
// the shared discriminator info is a snapshot taken before any model is built.
func TestModelAnalysis_IsSnapshotNotAnalysis(t *testing.T) {
	defer discardOutput()()

	raw, err := json.Marshal(nestedAllOfSpec())
	require.NoError(t, err)

	mutatedDoc, err := loads.Analyzed(raw, "2.0")
	require.NoError(t, err)
	pristineDoc, err := loads.Analyzed(raw, "2.0")
	require.NoError(t, err)

	opts := opts()
	shared := newModelAnalysis(mutatedDoc, analysis.New(mutatedDoc.Spec()), opts)
	pristine := newModelAnalysis(pristineDoc, analysis.New(pristineDoc.Spec()), opts)

	// generating Outer first adds definitions to the spec
	_, err = makeGenDefinitionWithAnalysis("Outer", "models", mutatedDoc.Spec().Definitions["Outer"], mutatedDoc, opts, shared)
	require.NoError(t, err)

	fromMutatedSpec := discriminatorInfo(analysis.New(mutatedDoc.Spec()), opts)
	require.NotEqual(t, pristine.discriminators, fromMutatedSpec,
		"expected a fresh analysis of the mutated spec to see synthetic discriminator children")

	want, err := makeGenDefinitionWithAnalysis("Base", "models", pristineDoc.Spec().Definitions["Base"], pristineDoc, opts, pristine)
	require.NoError(t, err)
	got, err := makeGenDefinitionWithAnalysis("Base", "models", mutatedDoc.Spec().Definitions["Base"], mutatedDoc, opts, shared)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

// nestedAllOfSpec builds a spec where generating Outer registers a synthetic definition
// referencing the polymorphic Base.
func nestedAllOfSpec() *spec.Swagger {
	property := func(name string) map[string]spec.Schema {
		return map[string]spec.Schema{name: {SchemaProps: spec.SchemaProps{Type: []string{"string"}}}}
	}

	base := spec.Schema{SchemaProps: spec.SchemaProps{
		Type:       []string{"object"},
		Properties: property("kind"),
	}}
	base.Discriminator = "kind"

	outer := spec.Schema{SchemaProps: spec.SchemaProps{AllOf: []spec.Schema{{SchemaProps: spec.SchemaProps{
		AllOf: []spec.Schema{
			{SchemaProps: spec.SchemaProps{Ref: spec.MustCreateRef("#/definitions/Base")}},
			{SchemaProps: spec.SchemaProps{Type: []string{"object"}, Properties: property("x")}},
		},
	}}}}}

	return &spec.Swagger{SwaggerProps: spec.SwaggerProps{
		Swagger:     "2.0",
		Info:        &spec.Info{InfoProps: spec.InfoProps{Title: "disc", Version: "1.0.0"}},
		Paths:       &spec.Paths{Paths: map[string]spec.PathItem{}},
		Definitions: spec.Definitions{"Base": base, "Outer": outer},
	}}
}
