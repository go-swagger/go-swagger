package validate

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/jsonschema"
	"github.com/casualjim/go-swagger/spec"
)

// Parameter creates a parameter validator
func Parameter(param *spec.Parameter) (*jsonschema.Document, error) {
	if param.In == "body" {
		return Schema(param.Schema)
	}
	b, err := param.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return loadSchema(b)
}

func loadSchema(schemaJSON []byte) (*jsonschema.Document, error) {
	var doc interface{}
	if err := json.Unmarshal(schemaJSON, &doc); err != nil {
		return nil, err
	}
	return jsonschema.New(doc)
}

// Schema creates a schema validator
func Schema(schema *spec.Schema) (*jsonschema.Document, error) {
	b, err := schema.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return loadSchema(b)
}

// type parameterContext struct {
// 	request *http.Request
// 	route   *router.MatchedRoute
// 	model   map[string]interface{}
// 	param   *swagger.Parameter
// 	errors  []errors.Error
// }

// // Parameter validates a request parameter
// func Parameter(request *http.Request, route *router.MatchedRoute, model map[string]interface{}, param *swagger.Parameter) []errors.Error {
// 	var data string
// 	switch param.In {
// 	case "query":

// 	default:
// 		errors.NotImplemented(param.In + " parsing is not implemented yet")
// 	}
// 	return nil
// }
