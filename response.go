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

func (r Response) MarshalMap() map[string]interface{} {
	if r.Ref != "" {
		return map[string]interface{}{"$ref": r.Ref}
	}

	return reflection.MarshalMapRecursed(r)
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.MarshalMap())
}

func (r Response) MarshalYAML() (interface{}, error) {
	return r.MarshalMap(), nil
}
