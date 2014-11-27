package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// Allows adding meta data to a single tag that is used by the [Operation Object](http://goo.gl/8us55a#operationObject).
// It is not mandatory to have a Tag Object per tag used there.
//
// For more information: http://goo.gl/8us55a#tagObject
type Tag struct {
	Description  string                 `structs:"description,omitempty"`
	Extensions   map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Name         string                 `structs:"name"`
	ExternalDocs *ExternalDocumentation `structs:"externalDocs,omitempty"`
}

func (t Tag) Map() map[string]interface{} {
	res := structs.Map(t)
	addExtensions(res, t.Extensions)
	return res
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Map())
}

func (t Tag) MarshalYAML() (interface{}, error) {
	return t.Map(), nil
}
