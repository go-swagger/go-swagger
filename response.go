package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Response describes a single response from an API Operation.
//
// For more information: http://goo.gl/8us55a#responseObject
type Response struct {
	Description string            `swagger:"description,omitempty"`
	Ref         string            `swagger:"-"`
	Schema      *Schema           `swagger:"schema,omitempty"`
	Headers     map[string]Header `swagger:"headers,omitempty"`
	Examples    interface{}       `swagger:"examples,omitempty"`
}

// UnmarshalMap hydrates this response instance with the data from the map
func (r *Response) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if ref, ok := dict["$ref"]; ok {
		*r = Response{Ref: ref.(string)}
		return nil
	}
	return reflection.UnmarshalMapRecursed(dict, r)
}

// UnmarshalJSON hydrates this response instance with the data from JSON
func (r *Response) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return r.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this response instance with the data from YAML
func (r *Response) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return r.UnmarshalMap(value)
}

// MarshalMap converts this response object to a map
func (r Response) MarshalMap() map[string]interface{} {
	if r.Ref != "" {
		return map[string]interface{}{"$ref": r.Ref}
	}

	return reflection.MarshalMapRecursed(r)
}

// MarshalJSON converts this response object to JSON
func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.MarshalMap())
}

// MarshalYAML converts this response object to YAML
func (r Response) MarshalYAML() (interface{}, error) {
	return r.MarshalMap(), nil
}
