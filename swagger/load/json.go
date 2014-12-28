package load

import (
	"encoding/json"
	"io/ioutil"

	"github.com/casualjim/go-swagger/swagger/spec"
)

// JSONSpec loads a spec from a json document
func JSONSpec(path string) (*spec.Document, error) {
	loader := loadStrategy(path, ioutil.ReadFile, loadHTTPBytes)

	data, err := loader(path)
	if err != nil {
		return nil, err
	}
	// convert to json
	return spec.New(json.RawMessage(data), "")
}
