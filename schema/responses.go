package schema

import (
	"encoding/json"
	"strconv"
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
	Extensions          map[string]interface{}
	Default             *Response
	StatusCodeResponses map[int]Response
}

func (r Responses) Map() map[string]interface{} {
	res := make(map[string]interface{})
	if r.Default != nil {
		res["default"] = r.Default.Map()
	}
	for k, v := range r.StatusCodeResponses {
		res[strconv.Itoa(k)] = v.Map()
	}
	addExtensions(res, r.Extensions)
	return res
}

func (r Responses) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Map())
}
func (r Responses) MarshalYAML() (interface{}, error) {
	return r.Map(), nil
}
