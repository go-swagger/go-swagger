package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Header describes a header for a response of the API
//
// For more information: http://goo.gl/8us55a#headerObject
type Header struct {
	Description      string        `swagger:"description,omitempty"`
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
	Type             string        `swagger:"type,omitempty"`
	Format           string        `swagger:"format,omitempty"`
	Default          interface{}   `swagger:"default,omitempty"`
	Items            *Items        `swagger:"-"`
}

// UnmarshalJSON hydrates this header from JSON
func (h *Header) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, h)
}

// UnmarshalYAML hydrates this header from YAML
func (h *Header) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, h)
}

// MarshalJSON converts this header object to JSON
func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(h))
}

// MarshalYAML converts this header object to YAML
func (h Header) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(h), nil
}
