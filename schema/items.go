package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

type Items struct {
	Ref              string        `structs:"-"`
	Type             string        `structs:"type,omitempty"`
	Format           string        `structs:"format,omitempty"`
	Items            *Items        `structs:"-"`
	CollectionFormat string        `structs:"collectionFormat,omitempty"`
	Default          interface{}   `structs:"default,omitempty"`
	Maximum          float64       `structs:"maximum,omitempty"`
	ExclusiveMaximum bool          `structs:"exclusiveMaximum,omitempty"`
	Minimum          float64       `structs:"minimum,omitempty"`
	ExclusiveMinimum bool          `structs:"exclusiveMinimum,omitempty"`
	MaxLength        int64         `structs:"maxLength,omitempty"`
	MinLength        int64         `structs:"minLength,omitempty"`
	Pattern          string        `structs:"pattern,omitempty"`
	MaxItems         int64         `structs:"maxItems,omitempty"`
	MinItems         int64         `structs:"minItems,omitempty"`
	UniqueItems      bool          `structs:"uniqueItems,omitempty"`
	MultipleOf       float64       `structs:"multipleOf,omitempty"`
	Enum             []interface{} `structs:"enum,omitempty"`
}

func (i Items) Map() map[string]interface{} {
	if i.Ref != "" {
		return map[string]interface{}{"$ref": i.Ref}
	}
	res := structs.Map(i)
	if i.Items != nil {
		res["items"] = i.Items.Map()
	}
	return res
}

func (i Items) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Map())
}

func (i Items) MarshalYAML() (interface{}, error) {
	return i.Map(), nil
}
