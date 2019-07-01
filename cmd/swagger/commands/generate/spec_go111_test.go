// +build go1.11

package generate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-swagger/go-swagger/codescan"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestGenerateJSONSpec(t *testing.T) {
	opts := codescan.Options{
		WorkDir:  basePath,
		Packages: []string{"./..."},
	}

	swspec, err := codescan.Run(&opts)
	assert.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(jsonResultFile)
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

	expected, err := ioutil.ReadFile(yamlResultFile)
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
