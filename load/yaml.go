package load

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
	"gopkg.in/yaml.v2"
)

// YAMLSpec loads a swagger spec document
func YAMLSpec(path string) (*spec.Document, error) {
	yamlDoc, err := YAMLDoc(path)
	if err != nil {
		return nil, err
	}

	data, err := util.YAMLToJSON(yamlDoc)
	if err != nil {
		return nil, err
	}

	return spec.New(data, "")
}

// YAMLDoc loads a yaml document from either http or a file
func YAMLDoc(path string) (interface{}, error) {
	loadData := loadStrategy(path, ioutil.ReadFile, loadHTTPBytes)
	data, err := loadData(path)
	if err != nil {
		return nil, err
	}

	return bytesToYAMLDoc(data)
}

func loadStrategy(path string, local, remote func(string) ([]byte, error)) func(string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return remote
	}
	return local
}

func bytesToYAMLDoc(data []byte) (interface{}, error) {
	var document map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	return document, nil
}

func loadHTTPBytes(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not access document at %q [%s] ", path, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
