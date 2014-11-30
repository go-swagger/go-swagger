package swagger

import (
	"encoding/json"
	"regexp"
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

var onlyNumbers = regexp.MustCompile("\\d+")

// UnmarshalMap hydrates this responses instance with the data from the map
func (r *Responses) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if def, ok := dict["default"]; ok {
		value := &Response{}
		if err := value.UnmarshalMap(def.(map[string]interface{})); err != nil {
			return err
		}
		r.Default = value
		delete(dict, "default")
	}
	r.Extensions = readExtensions(dict)
	statusCodeResponses := make(map[int]Response)
	for k, v := range dict {
		if nk, err := strconv.Atoi(k); err == nil {
			resVal, ok := v.(map[string]interface{})
			if !ok {
				resVal = reflection.MarshalMap(resVal)
			}
			response := new(Response)
			if err := response.UnmarshalMap(resVal); err != nil {
				return err
			}
			statusCodeResponses[nk] = *response
		}
	}
	r.StatusCodeResponses = statusCodeResponses
	return nil
}

// UnmarshalJSON hydrates this responses instance with the data from JSON
func (r *Responses) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return r.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this responses instance with the data from YAML
func (r *Responses) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return r.UnmarshalMap(value)
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
