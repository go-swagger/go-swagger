package generate

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/go-swagger/go-swagger/scan"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const (
	basePath       = "../../../../fixtures/goparsing/spec"
	jsonResultFile = "../../../../fixtures/goparsing/spec/api_spec.json"
	yamlResultFile = "../../../../fixtures/goparsing/spec/api_spec.yml"
)

func TestGenerateJSONSpec(t *testing.T) {
	opts := scan.Opts{
		BasePath: basePath,
	}

	swspec, err := scan.Application(opts)
	assert.NoError(t, err)

	data, err := marshalToJSONFormat(swspec, true)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(jsonResultFile)
	assert.NoError(t, err)

	varifyJSONData(t, data, expected)
}

func TestGenerateYAMLSpec(t *testing.T) {
	opts := scan.Opts{
		BasePath: basePath,
	}

	swspec, err := scan.Application(opts)
	assert.NoError(t, err)

	data, err := marshalToYAMLFormat(swspec)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(yamlResultFile)
	assert.NoError(t, err)

	varifyYAMLData(t, data, expected)
}

func varifyJSONData(t *testing.T, data, expectedJSON []byte) {
	var got interface{}
	var expected interface{}

	err := json.Unmarshal(data, &got)
	assert.NoError(t, err)

	err = json.Unmarshal(expectedJSON, &expected)
	assert.NoError(t, err)

	if !assert.ObjectsAreEqual(got, expected) {
		assert.Fail(t, "marshaled JSON data doesn't equal expected JSON data")
	}
}

func varifyYAMLData(t *testing.T, data, expectedYAML []byte) {
	var got interface{}
	var expected interface{}

	err := yaml.Unmarshal(data, &got)
	assert.NoError(t, err)

	err = yaml.Unmarshal(expectedYAML, &expected)
	assert.NoError(t, err)

	if !assert.ObjectsAreEqual(got, expected) {
		assert.Fail(t, "marshaled YAML data doesn't equal expected YAML data")
	}
}
