package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Items a limited subset of JSON-Schema's items object.
// It is used by parameter definitions that are not located in "body".
//
// For more information: http://goo.gl/8us55a#items-object-
type Items struct {
	Ref              string        `swagger:"-"`
	Type             string        `swagger:"type,omitempty"`
	Format           string        `swagger:"format,omitempty"`
	Items            *Items        `swagger:"items,omitempty"`
	CollectionFormat string        `swagger:"collectionFormat,omitempty"`
	Default          interface{}   `swagger:"default,omitempty"`
	Maximum          float64       `swagger:"maximum,omitempty"`
	ExclusiveMaximum bool          `swagger:"exclusiveMaximum,omitempty"`
	Minimum          float64       `swagger:"minimum,omitempty"`
	ExclusiveMinimum bool          `swagger:"exclusiveMinimum,omitempty"`
	MaxLength        int64         `swagger:"maxLength,omitempty"`
	MinLength        int64         `swagger:"minLength,omitempty"`
	Pattern          string        `swagger:"pattern,omitempty"`
	MaxItems         int64         `swagger:"maxItems,omitempty"`
	MinItems         int64         `swagger:"minItems,omitempty"`
	UniqueItems      bool          `swagger:"uniqueItems,omitempty"`
	MultipleOf       float64       `swagger:"multipleOf,omitempty"`
	Enum             []interface{} `swagger:"enum,omitempty"`
}

// MarshalMap converts this items object to a map
func (i Items) MarshalMap() map[string]interface{} {
	if i.Ref != "" {
		return map[string]interface{}{"$ref": i.Ref}
	}
	return reflection.MarshalMapRecursed(i)
}

// MarshalJSON converts this items object to JSON
func (i Items) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.MarshalMap())
}

// MarshalYAML converts this items object to YAML
func (i Items) MarshalYAML() (interface{}, error) {
	return i.MarshalMap(), nil
}
