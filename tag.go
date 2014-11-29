package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Tag allows adding meta data to a single tag that is used by the [Operation Object](http://goo.gl/8us55a#operationObject).
// It is not mandatory to have a Tag Object per tag used there.
//
// For more information: http://goo.gl/8us55a#tagObject
type Tag struct {
	Description  string                 `swagger:"description,omitempty"`
	Extensions   map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Name         string                 `swagger:"name"`
	ExternalDocs *ExternalDocumentation `swagger:"externalDocs,omitempty"`
}

// MarshalMap converts this tag object into a map
func (t Tag) MarshalMap() map[string]interface{} {
	res := reflection.MarshalMapRecursed(t)
	addExtensions(res, t.Extensions)
	return res
}

// MarshalJSON converts this tag object into JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.MarshalMap())
}

// MarshalYAML converts this tag object into YAML
func (t Tag) MarshalYAML() (interface{}, error) {
	return t.MarshalMap(), nil
}

// UnmarshalMap hydrates this tag instance with the data from the map
func (t *Tag) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, t); err != nil {
		return err
	}
	t.Extensions = readExtensions(dict)
	return nil
}

// UnmarshalJSON hydrates this tag instance with the data from JSON
func (t *Tag) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return t.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this tag instance with the data from YAML
func (t *Tag) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return t.UnmarshalMap(value)
}
