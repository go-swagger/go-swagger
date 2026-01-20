// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/3idey/codescan/codescan"
	"github.com/go-openapi/spec"
)

func TestSpecEmbeddedDescriptionAndTags_Issue3125(t *testing.T) {
	// test the full repro cas provided by the OP, issue #3125
	workDir := prepare3125(t)

	t.Run("should NOT render siblings with $ref", testEmbeddedDescriptionAndTagsFull(workDir, false, "expected_swagger_noallow.yaml"))

	t.Run("should render siblings with $ref", testEmbeddedDescriptionAndTagsFull(workDir, true, "expected_swagger_allow.yaml"))
}

func prepare3125(t *testing.T) string {
	t.Helper()

	workDir := filepath.Join("..", "..", "..", "..", "fixtures", "bugs", "3125", "full")
	source := os.DirFS(workDir)
	target := filepath.Join(t.TempDir(), "swagger")

	require.NoError(t, os.CopyFS(target, source))
	t.Run("go mod", gomodinit(target))
	t.Run("go mod download", gomoddownload(target))
	t.Run("go mod tidy", gomodtidy(target))

	return target
}

func testEmbeddedDescriptionAndTagsFull(srcPath string, allowDescWithRef bool, expectedYAML string) func(*testing.T) {
	return func(t *testing.T) {
		opts := &codescan.Options{
			Packages: []string{
				"./...",
			},
			WorkDir:     srcPath,
			ScanModels:  true,
			DescWithRef: allowDescWithRef,
			BuildTags:   "testintegration",
		}

		swspec, err := codescan.Run(opts)
		require.NoError(t, err)

		data, err := marshalToYAMLFormat(swspec)
		require.NoError(t, err)

		yamlResultTagInRef := filepath.Join(srcPath, expectedYAML)
		expected, err := os.ReadFile(yamlResultTagInRef)
		require.NoError(t, err)

		/* NOTE: uncomment those lines if you need to get a glimpse at the output to debug.
		require.NoError(t,
			os.WriteFile(fmt.Sprintf("expected_desc_ref_%t.yaml", allowDescWithRef), expected, 0o600),
		)
		require.NoError(t,
			os.WriteFile(fmt.Sprintf("generated_desc_ref_%t.yaml", allowDescWithRef), data, 0o600),
		)
		*/

		verifyYAMLData(t, data, expected)
	}
}

func verifyYAMLData(t *testing.T, data, expectedYAML []byte) {
	t.Helper()

	var got, expected any

	require.NoError(t, yaml.Unmarshal(data, &got))
	require.NoError(t, yaml.Unmarshal(expectedYAML, &expected))
	assert.Equal(t, expected, got)
}

func marshalToYAMLFormat(swspec *spec.Swagger) ([]byte, error) {
	b, err := json.Marshal(swspec)
	if err != nil {
		return nil, err
	}

	var jsonObj any
	if err := yaml.Unmarshal(b, &jsonObj); err != nil {
		return nil, err
	}

	return yaml.Marshal(jsonObj)
}
