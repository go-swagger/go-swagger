package generate

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-swagger/go-swagger/codescan"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

const (
	basePath       = "../../../../fixtures/goparsing/spec"
	jsonResultFile = basePath + "/api_spec_go111.json"
	yamlResultFile = basePath + "/api_spec_go111.yml"
)

func TestSpecFileExecute(t *testing.T) {
	files := []string{"", "spec.json", "spec.yml", "spec.yaml"}
	for _, outputFile := range files {
		spec := &SpecFile{
			WorkDir: basePath,
			Output:  flags.Filename(outputFile),
		}

		err := spec.Execute(nil)
		assert.NoError(t, err)
		if outputFile != "" {
			_ = os.Remove(outputFile)
		}
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

	var got map[string]interface{}
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Len(t, got["definitions"], 2)
	require.Contains(t, got["definitions"], "Item")
	itemDefinition := got["definitions"].(map[string]interface{})["Item"].(map[string]interface{})
	require.Contains(t, itemDefinition["properties"], "Value1")
	value1Property := itemDefinition["properties"].(map[string]interface{})["Value1"].(map[string]interface{})
	require.Contains(t, value1Property, "x-nullable")
	assert.Equal(t, true, value1Property["x-nullable"])
}

func TestGenerateJSONSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:  basePath,
		Packages: []string{"./..."},
	}

	swspec, err := codescan.Run(&opts)
	assert.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	assert.NoError(t, err)

	expected, err := os.ReadFile(jsonResultFile)
	assert.NoError(t, err)

	verifyJSONData(t, data, expected)
}

func TestGenerateYAMLSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:  basePath,
		Packages: []string{"./..."},
	}

	swspec, err := codescan.Run(&opts)
	assert.NoError(t, err)

	data, err := marshalToYAMLFormat(swspec)
	assert.NoError(t, err)

	expected, err := os.ReadFile(yamlResultFile)
	assert.NoError(t, err)

	verifyYAMLData(t, data, expected)
}

func verifyJSONData(t *testing.T, data, expectedJSON []byte) {
	var got interface{}
	var expected interface{}

	err := json.Unmarshal(data, &got)
	assert.NoError(t, err)

	err = json.Unmarshal(expectedJSON, &expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func verifyYAMLData(t *testing.T, data, expectedYAML []byte) {
	var got interface{}
	var expected interface{}

	err := yaml.Unmarshal(data, &got)
	assert.NoError(t, err)

	err = yaml.Unmarshal(expectedYAML, &expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
