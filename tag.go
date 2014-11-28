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
