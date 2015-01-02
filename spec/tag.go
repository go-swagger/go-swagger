package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/util"
)

type tagProps struct {
	Description  string                 `json:"description,omitempty"`
	Name         string                 `json:"name,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// NewTag creates a new tag
func NewTag(name, description string, externalDocs *ExternalDocumentation) Tag {
	return Tag{tagProps: tagProps{description, name, externalDocs}}
}

// Tag allows adding meta data to a single tag that is used by the [Operation Object](http://goo.gl/8us55a#operationObject).
// It is not mandatory to have a Tag Object per tag used there.
//
// For more information: http://goo.gl/8us55a#tagObject
type Tag struct {
	vendorExtensible
	tagProps
}

// MarshalJSON marshal this to JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(t.tagProps)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(t.vendorExtensible)
	if err != nil {
		return nil, err
	}
	return util.ConcatJSON(b1, b2), nil
}

// UnmarshalJSON marshal this from JSON
func (t *Tag) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &t.tagProps); err != nil {
		return err
	}
	return json.Unmarshal(data, &t.vendorExtensible)
}
