// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestDeepCloneSpec_SimpleSpec(t *testing.T) {
	// Load a simple fixture spec
	specDoc, err := loads.Spec("../fixtures/codegen/simplesearch.yml")
	require.NoError(t, err)

	original := specDoc.Spec()
	require.NotNil(t, original)

	cloned, err := deepCloneSpec(original)
	require.NoError(t, err)
	require.NotNil(t, cloned)

	// Verify basic properties are preserved
	assert.Equal(t, original.Swagger, cloned.Swagger)
	assert.Equal(t, original.Info.Title, cloned.Info.Title)
	assert.Equal(t, original.Info.Version, cloned.Info.Version)

	// Verify definitions are preserved
	assert.Len(t, cloned.Definitions, len(original.Definitions))

	// Verify deep independence: modifying clone should not affect original
	cloned.Info.Title = "Modified API"
	assert.Equal(t, "Modified API", cloned.Info.Title)
	assert.Equal(t, "Simple Search API", original.Info.Title)
}

func TestDeepCloneSpec_NilInput(t *testing.T) {
	cloned, err := deepCloneSpec(nil)
	require.NoError(t, err)
	// deepCloneSpec returns an empty Swagger spec for nil input (from JSON unmarshaling nil)
	require.NotNil(t, cloned)
}

func TestDeepCloneSpec_ComplexSpec(t *testing.T) {
	// Load a more complex fixture spec
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	original := specDoc.Spec()
	require.NotNil(t, original)

	cloned, err := deepCloneSpec(original)
	require.NoError(t, err)
	require.NotNil(t, cloned)

	// Verify definitions are preserved
	assert.Len(t, cloned.Definitions, len(original.Definitions))

	// Verify each definition is cloned
	for name := range original.Definitions {
		assert.Contains(t, cloned.Definitions, name)
		// Verify the cloned definition values match
		originalDef := original.Definitions[name]
		clonedDef := cloned.Definitions[name]
		assert.Equal(t, originalDef.Type, clonedDef.Type)
	}
}

func TestDeepCloneSpec_PathsAndOperations(t *testing.T) {
	// Load a spec with paths
	specDoc, err := loads.Spec("../fixtures/codegen/simplesearch.yml")
	require.NoError(t, err)

	original := specDoc.Spec()
	require.NotNil(t, original)

	cloned, err := deepCloneSpec(original)
	require.NoError(t, err)
	require.NotNil(t, cloned)

	// Verify paths are preserved
	assert.NotNil(t, cloned.Paths)
	if original.Paths != nil {
		assert.Len(t, cloned.Paths.Paths, len(original.Paths.Paths))

		// Verify each path is cloned
		for path := range original.Paths.Paths {
			assert.Contains(t, cloned.Paths.Paths, path)
		}
	}
}

func TestAnalyzedSpecCache_ThreadSafety(t *testing.T) {
	// Test that read/write operations on the cache are thread-safe
	opts := testGenOpts()
	done := make(chan bool)

	// Set thread 1
	go func() {
		defer func() { done <- true }()
		opts.setCachedRawSpec(nil)
	}()

	// Set thread 2
	go func() {
		defer func() { done <- true }()
		opts.setCachedRawSpec(nil)
	}()

	// Get thread 1
	go func() {
		defer func() { done <- true }()
		_ = opts.getAnalyzedSpec()
	}()

	// Get thread 2
	go func() {
		defer func() { done <- true }()
		_ = opts.getAnalyzedSpec()
	}()

	<-done
	<-done
	<-done
	<-done
	// If we got here without deadlock, thread safety works
}

func TestAnalyzedSpecCache_BasicOperations(t *testing.T) {
	opts := testGenOpts()

	// Initially, cache should be nil
	assert.Nil(t, opts.getAnalyzedSpec())

	// Set the cache
	opts.setCachedRawSpec(nil)

	// Retrieve the cache
	assert.Nil(t, opts.getAnalyzedSpec())

	// Set again
	opts.setCachedRawSpec(nil)

	// Retrieve again
	assert.Nil(t, opts.getAnalyzedSpec())
}

func TestDeepCloneSpec_IndependentMutations(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/simplesearch.yml")
	require.NoError(t, err)

	// Clone the spec multiple times
	cloned1, err1 := deepCloneSpec(specDoc.Spec())
	cloned2, err2 := deepCloneSpec(specDoc.Spec())
	require.NoError(t, err1)
	require.NoError(t, err2)

	// Modify first clone
	cloned1.Info.Title = "Clone 1"
	cloned1.Info.Version = "2.0.0"

	// Modify second clone
	cloned2.Info.Title = "Clone 2"
	cloned2.Info.Version = "3.0.0"

	// Verify each clone has its own modifications
	assert.Equal(t, "Clone 1", cloned1.Info.Title)
	assert.Equal(t, "2.0.0", cloned1.Info.Version)

	assert.Equal(t, "Clone 2", cloned2.Info.Title)
	assert.Equal(t, "3.0.0", cloned2.Info.Version)

	// Verify original is unchanged
	assert.Equal(t, "Simple Search API", specDoc.Spec().Info.Title)
	assert.Equal(t, "1.0.0", specDoc.Spec().Info.Version)
}

// TestAnalyzedSpecCache_WithEmptyCache verifies fallback behavior when cache is empty.
func TestAnalyzedSpecCache_WithEmptyCache(t *testing.T) {
	opts := testGenOpts()

	// Initially, getAnalyzedSpec should return nil
	analyzed := opts.getAnalyzedSpec()
	assert.Nil(t, analyzed)
}

// TestDeepCloneSpec_PreservesJSONRoundtrip verifies that deep cloning
// preserves the spec through JSON round-trips.
func TestDeepCloneSpec_PreservesJSONRoundtrip(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.models.yml")
	require.NoError(t, err)

	original := specDoc.Spec()
	require.NotNil(t, original)

	// Deep clone
	cloned, err := deepCloneSpec(original)
	require.NoError(t, err)
	require.NotNil(t, cloned)

	// Deep clone again (second roundtrip)
	cloned2, err := deepCloneSpec(cloned)
	require.NoError(t, err)
	require.NotNil(t, cloned2)

	// Verifying the structure remains consistent through multiple clones
	assert.Equal(t, original.Swagger, cloned2.Swagger)
	assert.Equal(t, original.Info.Title, cloned2.Info.Title)
	assert.Equal(t, original.Info.Version, cloned2.Info.Version)
	assert.Len(t, cloned2.Definitions, len(original.Definitions))
}

// TestAnalyzedSpecCache_FreshAnalysis verifies that getAnalyzedSpec()
// returns a freshly analyzed spec on each call, preventing internal state
// mutations from affecting subsequent retrievals.
func TestAnalyzedSpecCache_FreshAnalysis(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/simplesearch.yml")
	require.NoError(t, err)

	opts := testGenOpts()
	opts.setCachedRawSpec(specDoc.Spec())

	// Get analyzed spec twice
	analyzed1 := opts.getAnalyzedSpec()
	analyzed2 := opts.getAnalyzedSpec()

	require.NotNil(t, analyzed1)
	require.NotNil(t, analyzed2)

	// Verify they are different objects (fresh analysis each time)
	assert.NotSame(t, analyzed1, analyzed2,
		"getAnalyzedSpec() should return a fresh analyzed spec on each call")

	// Verify both have the same number of definition references
	defs1 := analyzed1.AllDefinitions()
	defs2 := analyzed2.AllDefinitions()
	assert.Equal(t, len(defs1), len(defs2),
		"Both analyzed specs should have the same number of definitions")

	// Verify both have the same operation IDs
	ops1 := analyzed1.OperationIDs()
	ops2 := analyzed2.OperationIDs()
	assert.ElementsMatch(t, ops1, ops2,
		"Both analyzed specs should have the same operation IDs")
}
