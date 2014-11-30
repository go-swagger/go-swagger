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

// UnmarshalMap hydrates this items instance with the data from the map
func (i *Items) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if ref, ok := dict["$ref"]; ok {
		i.Ref = ref.(string)
	}
	return reflection.UnmarshalMapRecursed(dict, i)
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (i *Items) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return i.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this items instance with the data from YAML
func (i *Items) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.UnmarshalMap(value)
}

// MarshalMap converts this items object to a map
func (i Items) MarshalMap() map[string]interface{} {
	result := reflection.MarshalMapRecursed(i)
	if i.Ref != "" {
		result["$ref"] = i.Ref
	}
	return result
}

// MarshalJSON converts this items object to JSON
func (i Items) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.MarshalMap())
}

// MarshalYAML converts this items object to YAML
func (i Items) MarshalYAML() (interface{}, error) {
	return i.MarshalMap(), nil
}
