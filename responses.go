package swagger

import (
	"encoding/json"
	"strconv"

	"github.com/casualjim/go-swagger/reflection"
)

// Responses is a container for the expected responses of an operation.
// The container maps a HTTP response code to the expected response.
// It is not expected from the documentation to necessarily cover all possible HTTP response codes,
// since they may not be known in advance. However, it is expected from the documentation to cover
// a successful operation response and any known errors.
//
// The `default` can be used a default response object for all HTTP codes that are not covered
// individually by the specification.
//
// The `Responses Object` MUST contain at least one response code, and it SHOULD be the response
// for a successful operation call.
//
// For more information: http://goo.gl/8us55a#responsesObject
type Responses struct {
	Extensions          map[string]interface{} `swagger:"-"`
	Default             *Response              `swagger:"-"`
	StatusCodeResponses map[int]Response       `swagger:"-"`
}

// MarshalMap converts this responses object to a map
func (r Responses) MarshalMap() map[string]interface{} {
	res := make(map[string]interface{})
	if r.Default != nil {
		res["default"] = reflection.MarshalMap(r.Default)
	}
	for k, v := range r.StatusCodeResponses {
		res[strconv.Itoa(k)] = reflection.MarshalMap(v)
	}
	addExtensions(res, r.Extensions)
	return res
}

// MarshalJSON converts this responses object to JSON
func (r Responses) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.MarshalMap())
}

// MarshalYAML converts this responses object to YAML
func (r Responses) MarshalYAML() (interface{}, error) {
	return r.MarshalMap(), nil
}
