package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/swag"
)

type responseProps struct {
	Description string            `json:"description,omitempty"`
	Schema      *Schema           `json:"schema,omitempty"`
	Headers     map[string]Header `json:"headers,omitempty"`
	Examples    interface{}       `json:"examples,omitempty"`
}

// Response describes a single response from an API Operation.
//
// For more information: http://goo.gl/8us55a#responseObject
type Response struct {
	refable
	responseProps
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (r *Response) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &r.responseProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &r.refable); err != nil {
		return err
	}
	return nil
}

// MarshalJSON converts this items object to JSON
func (r Response) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(r.responseProps)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(r.refable)
	if err != nil {
		return nil, err
	}
	return swag.ConcatJSON(b1, b2), nil
}
