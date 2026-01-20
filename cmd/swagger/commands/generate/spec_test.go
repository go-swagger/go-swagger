// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/3idey/codescan/codescan"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const (
	basePath       = "../../../../fixtures/goparsing/spec"
	jsonResultFile = basePath + "/api_spec_go111.json"
	yamlResultFile = basePath + "/api_spec_go111.yml"

	jsonResultFileRef = basePath + "/api_spec_go111_ref.json"
	yamlResultFileRef = basePath + "/api_spec_go111_ref.yml"

	jsonResultFileTransparent = basePath + "/api_spec_go111_transparent.json"
	yamlResultFileTransparent = basePath + "/api_spec_go111_transparent.yml"
)

var enableSpecOutput bool

func init() {
	flag.BoolVar(&enableSpecOutput, "enable-spec-output", false, "enable spec gen test to write output to a file")
}

func TestSpecFileExecute(t *testing.T) {
	files := []string{"", "spec.json", "spec.yml", "spec.yaml"}

	for _, outputFile := range files {
		name := outputFile
		if outputFile == "" {
			name = "to stdout"
		}

		t.Run("should produce spec file "+name, func(t *testing.T) {
			spec := &SpecFile{
				WorkDir: basePath,
				Output:  flags.Filename(outputFile),
			}
			if outputFile == "" {
				defaultWriter = io.Discard
			}
			t.Cleanup(func() {
				if outputFile != "" {
					_ = os.Remove(outputFile)
				} else {
					defaultWriter = os.Stdout
				}
			})

			require.NoError(t, spec.Execute(nil))
		})
	}
}

func TestSpecFileExecuteRespectsSetXNullableForPointersOption(t *testing.T) {
	outputFileName := "spec.json"
	spec := &SpecFile{
		WorkDir:                 "../../../../fixtures/enhancements/pointers-nullable-by-default",
		Output:                  flags.Filename(outputFileName),
		ScanModels:              true,
		SetXNullableForPointers: true,
	}

	defer func() { _ = os.Remove(outputFileName) }()

	err := spec.Execute(nil)
	require.NoError(t, err)

	data, err := os.ReadFile(outputFileName)
	require.NoError(t, err)

	var got map[string]any
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Len(t, got["definitions"], 2)
	require.Contains(t, got["definitions"], "Item")
	itemDefinition, ok := got["definitions"].(map[string]any)["Item"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, itemDefinition["properties"], "Value1")
	value1Property, ok := itemDefinition["properties"].(map[string]any)["Value1"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, value1Property, "x-nullable")
	assert.Equal(t, true, value1Property["x-nullable"])
}

func TestGenerateJSONSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:  basePath,
		Packages: []string{"./..."},
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	require.NoError(t, err)

	expected, err := os.ReadFile(jsonResultFile)
	require.NoError(t, err)

	verifyJSONData(t, data, expected)
}

func TestGenerateYAMLSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:  basePath,
		Packages: []string{"./..."},
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToYAMLFormat(swspec)
	require.NoError(t, err)

	expected, err := os.ReadFile(yamlResultFile)
	require.NoError(t, err)
	{
		var jsonObj any
		require.NoError(t, yaml.Unmarshal(expected, &jsonObj))

		rewritten, err := yaml.Marshal(jsonObj)
		require.NoError(t, err)
		expected = rewritten
	}

	if enableSpecOutput {
		require.NoError(t,
			os.WriteFile("expected.yaml", expected, 0o600),
		)
		require.NoError(t,
			os.WriteFile("generated.yaml", data, 0o600),
		)
	}

	verifyYAMLData(t, data, expected)
}

func TestGenerateJSONSpecWithSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:    basePath,
		Packages:   []string{"./..."},
		RefAliases: true,
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	require.NoError(t, err)

	expected, err := os.ReadFile(jsonResultFileRef)
	require.NoError(t, err)

	verifyJSONData(t, data, expected)
}

func TestGenerateYAMLSpecWithRefAliases(t *testing.T) {
	opts := codescan.Options{
		WorkDir:    basePath,
		Packages:   []string{"./..."},
		RefAliases: true,
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToYAMLFormat(swspec)
	require.NoError(t, err)

	expected, err := os.ReadFile(yamlResultFileRef)
	require.NoError(t, err)
	{
		var jsonObj any
		require.NoError(t, yaml.Unmarshal(expected, &jsonObj))

		rewritten, err := yaml.Marshal(jsonObj)
		require.NoError(t, err)
		expected = rewritten
	}

	if enableSpecOutput {
		require.NoError(t,
			os.WriteFile("expected_ref.yaml", expected, 0o600),
		)
		require.NoError(t,
			os.WriteFile("generated_ref.yaml", data, 0o600),
		)
	}

	verifyYAMLData(t, data, expected)
}

func TestGenerateJSONSpecWithTransparentAliases(t *testing.T) {
	opts := codescan.Options{
		WorkDir:            basePath,
		Packages:           []string{"./..."},
		TransparentAliases: true,
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	require.NoError(t, err)

	expected, err := os.ReadFile(jsonResultFileTransparent)
	require.NoError(t, err)

	verifyJSONData(t, data, expected)
}

func TestGenerateYAMLSpecWithTransparentAliases(t *testing.T) {
	opts := codescan.Options{
		WorkDir:            basePath,
		Packages:           []string{"./..."},
		TransparentAliases: true,
	}

	swspec, err := codescan.Run(&opts)
	require.NoError(t, err)

	data, err := marshalToYAMLFormat(swspec)
	require.NoError(t, err)

	expected, err := os.ReadFile(yamlResultFileTransparent)
	require.NoError(t, err)
	{
		var jsonObj any
		require.NoError(t, yaml.Unmarshal(expected, &jsonObj))

		rewritten, err := yaml.Marshal(jsonObj)
		require.NoError(t, err)
		expected = rewritten
	}

	if enableSpecOutput {
		require.NoError(t,
			os.WriteFile("expected_transparent.yaml", expected, 0o600),
		)
		require.NoError(t,
			os.WriteFile("generated_transparent.yaml", data, 0o600),
		)
	}

	verifyYAMLData(t, data, expected)
}

func verifyJSONData(t *testing.T, data, expectedJSON []byte) {
	t.Helper()

	var got, expected any

	require.NoError(t, json.Unmarshal(data, &got))
	require.NoError(t, json.Unmarshal(expectedJSON, &expected))
	assert.Equal(t, expected, got)
}

func verifyYAMLData(t *testing.T, data, expectedYAML []byte) {
	t.Helper()

	var got, expected any

	require.NoError(t, yaml.Unmarshal(data, &got))
	require.NoError(t, yaml.Unmarshal(expectedYAML, &expected))
	assert.Equal(t, expected, got)
}
