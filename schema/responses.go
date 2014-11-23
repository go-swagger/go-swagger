package schema

import (
	"encoding/json"
	"strconv"
)

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
