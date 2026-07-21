// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// TestSpec_RootedLoader verifies that the "-Rooted" codegen option produces a document loader that
// confines local reads to the root directory: in-root documents load, anything escaping the root
// (traversal or an absolute path outside it) is rejected.
func TestSpec_RootedLoader(t *testing.T) {
	defer discardOutput()()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "child.json"), []byte(`{"type":"string"}`), 0o600))
	parent := filepath.Dir(root)
	require.NoError(t, os.WriteFile(filepath.Join(parent, "secret.json"), []byte(`{"secret":true}`), 0o600))

	sa := newSpecAnalyzer(&GenOpts{Rooted: root})
	sa.setLoaderOptions()
	require.NotNil(t, sa.loader)

	t.Run("in-root document loads", func(t *testing.T) {
		b, err := sa.loader("child.json")
		require.NoError(t, err)
		assert.StringContainsT(t, string(b), "string")
	})

	t.Run("traversal escaping the root is rejected", func(t *testing.T) {
		_, err := sa.loader("../secret.json")
		require.Error(t, err)
	})

	t.Run("absolute path outside the root is rejected", func(t *testing.T) {
		_, err := sa.loader(filepath.Join(parent, "secret.json"))
		require.Error(t, err)
	})
}

// TestSpec_RestrictedLoader verifies that the "-Restricted" codegen option produces a document
// loader that refuses remote fetches to forbidden (loopback / private / link-local) addresses,
// including cloud-metadata endpoints.
func TestSpec_RestrictedLoader(t *testing.T) {
	defer discardOutput()()

	sa := newSpecAnalyzer(&GenOpts{Restricted: true})
	sa.setLoaderOptions()
	require.NotNil(t, sa.loader)

	for _, url := range []string{
		"http://127.0.0.1:65535/spec.json",        // loopback
		"http://169.254.169.254/latest/meta-data", // link-local (cloud metadata)
		"http://10.0.0.1/internal.json",           // private
	} {
		t.Run(url, func(t *testing.T) {
			_, err := sa.loader(url)
			require.Error(t, err)
		})
	}
}

// TestSpec_RootedModeBlocksEscapingRef verifies end to end that codegen preprocessing
// (validate + flatten) refuses to resolve a $ref escaping the root when run in rooted mode, while
// the same spec processes normally without confinement.
func TestSpec_RootedModeBlocksEscapingRef(t *testing.T) {
	defer discardOutput()()

	root := t.TempDir()
	parent := filepath.Dir(root)
	require.NoError(t, os.WriteFile(filepath.Join(parent, "secret.json"),
		[]byte(`{"definitions":{"Leak":{"type":"string"}}}`), 0o600))

	specPath := filepath.Join(root, "swagger.json")
	doc := `{"swagger":"2.0","info":{"title":"x","version":"1"},"paths":{},` +
		`"definitions":{"Ref":{"$ref":"../secret.json#/definitions/Leak"}}}`
	require.NoError(t, os.WriteFile(specPath, []byte(doc), 0o600))

	t.Run("unconfined resolves the escaping $ref", func(t *testing.T) {
		opts := testGenOpts()
		opts.Spec = specPath
		_, err := newSpecAnalyzer(opts).validateAndFlattenSpec()
		require.NoError(t, err)
	})

	t.Run("rooted mode blocks the escaping $ref", func(t *testing.T) {
		opts := testGenOpts()
		opts.Spec = specPath
		opts.Rooted = root
		_, err := newSpecAnalyzer(opts).validateAndFlattenSpec()
		require.Error(t, err)
	})
}
