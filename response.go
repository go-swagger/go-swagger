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
