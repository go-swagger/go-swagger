package load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/casualjim/go-swagger/swagger/spec"
	"gopkg.in/yaml.v2"
)

// YAMLSpec loads a swagger spec document
func YAMLSpec(path string) (*spec.Document, error) {
	yamlDoc, err := YAMLDoc(path)
	if err != nil {
		return nil, err
	}

	data, err := yamlToJSON(yamlDoc)
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

func yamlToJSON(data interface{}) (json.RawMessage, error) {
	jm, err := transformData(data)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(jm)
	return json.RawMessage(b), err
}

func transformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return nil, fmt.Errorf("types don't match: expect map key string or int get: %T", k)
			}
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = transformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	default:
		return in, nil
	}
	return in, nil
}
