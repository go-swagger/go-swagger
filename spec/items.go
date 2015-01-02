package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/util"
)

type refable struct {
	Ref Ref
}

func (j refable) MarshalJSON() ([]byte, error) {
	return j.Ref.MarshalJSON()
}

func (j *refable) UnmarshalJSON(d []byte) error {
	return j.Ref.UnmarshalJSON(d)
}

type simpleSchema struct {
	Type             string      `json:"type,omitempty"`
	Format           string      `json:"format,omitempty"`
	Items            *Items      `json:"items,omitempty"`
	CollectionFormat string      `json:"collectionFormat,omitempty"`
	Default          interface{} `json:"default,omitempty"`
}

type commonValidations struct {
	Maximum          *float64      `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64      `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64        `json:"maxLength,omitempty"`
	MinLength        *int64        `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         *int64        `json:"maxItems,omitempty"`
	MinItems         *int64        `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       *float64      `json:"multipleOf,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
}

// Items a limited subset of JSON-Schema's items object.
// It is used by parameter definitions that are not located in "body".
//
// For more information: http://goo.gl/8us55a#items-object-
type Items struct {
	refable
	commonValidations
	simpleSchema
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (i *Items) UnmarshalJSON(data []byte) error {
	var validations commonValidations
	if err := json.Unmarshal(data, &validations); err != nil {
		return err
	}
	var ref refable
	if err := json.Unmarshal(data, &ref); err != nil {
		return err
	}
	var simpleSchema simpleSchema
	if err := json.Unmarshal(data, &simpleSchema); err != nil {
		return err
	}
	i.refable = ref
	i.commonValidations = validations
	i.simpleSchema = simpleSchema
	return nil
}

// MarshalJSON converts this items object to JSON
func (i Items) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(i.commonValidations)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(i.simpleSchema)
	if err != nil {
		return nil, err
	}
	b3, err := json.Marshal(i.refable)
	if err != nil {
		return nil, err
	}
	return util.ConcatJSON(b3, b1, b2), nil
}
