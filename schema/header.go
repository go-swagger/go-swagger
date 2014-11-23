package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

type Header struct {
	Description      string        `structs:"description,omitempty"`
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
	Type             string        `structs:"type,omitempty"`
	Format           string        `structs:"format,omitempty"`
	Default          interface{}   `structs:"default,omitempty"`
	Items            *Items        `structs:"-"`
}

func (h Header) Map() map[string]interface{} {
	res := structs.Map(h)
	if h.Items != nil {
		res["items"] = h.Items.Map()
	}
	return res
}

func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Map())
}

func (h Header) MarshalYAML() (interface{}, error) {
	return h.Map(), nil
}

type ExternalDocumentation struct {
	Description string `structs:"description,omitempty"`
	URL         string `structs:"url"`
}
