package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Operation describes a single API operation on a path.
//
// For more information: http://goo.gl/8us55a#operationObject
type Operation struct {
	Description  string                 `swagger:"description,omitempty"`
	Extensions   map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Consumes     []string               `swagger:"consumes,omitempty"`
	Produces     []string               `swagger:"produces,omitempty"`
	Schemes      []string               `swagger:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Tags         []string               `swagger:"tags,omitempty"`
	Summary      string                 `swagger:"summary,omitempty"`
	ExternalDocs *ExternalDocumentation `swagger:"externalDocs,omitempty"`
	ID           string                 `swagger:"operationId,omitempty"`
	Deprecated   bool                   `swagger:"deprecated,omitempty"`
	Security     []map[string][]string  `swagger:"security,omitempty"`
	Parameters   []Parameter            `swagger:"parameters,omitempty"`
	Responses    Responses              `swagger:"responses,omitempty"`
}

// MarshalMap converts this operation to a map
func (o Operation) MarshalMap() map[string]interface{} {
	res := reflection.MarshalMapRecursed(o)
	addExtensions(res, o.Extensions)
	return res
}

// MarshalJSON converts this operation to JSON
func (o Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.MarshalMap())
}

// MarshalYAML converts this operation to YAML
func (o Operation) MarshalYAML() (interface{}, error) {
	return o.MarshalMap(), nil
}

// UnmarshalMap hydrates this operation instance with the data from the map
func (o *Operation) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, o); err != nil {
		return err
	}
	o.Extensions = readExtensions(dict)
	return nil
}

// UnmarshalJSON hydrates this operation instance with the data from JSON
func (o *Operation) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return o.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this operation instance with the data from YAML
func (o *Operation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return o.UnmarshalMap(value)
}
