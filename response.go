package swagger

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// Response describes a single response from an API Operation.
//
// For more information: http://goo.gl/8us55a#responseObject
type Response struct {
	Description string            `structs:"description,omitempty"`
	Ref         string            `structs:"-"`
	Schema      *Schema           `structs:"-"`
	Headers     map[string]Header `structs:"-"`
	Examples    interface{}       `structs:"examples,omitempty"`
}

func (r Response) Map() map[string]interface{} {
	if r.Ref != "" {
		return map[string]interface{}{"$ref": r.Ref}
	}

	res := structs.Map(r)

	if r.Schema != nil {
		res["schema"] = r.Schema.Map()
	}

	if len(r.Headers) > 0 {
		headers := make(map[string]map[string]interface{}, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Map()
		}
		res["headers"] = headers
	}

	return res
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Map())
}

func (r Response) MarshalYAML() (interface{}, error) {
	return r.Map(), nil
}
