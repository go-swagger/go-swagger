package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

func BooleanProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "boolean"}}
}

func StringProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}}
}
func Float64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "double"}
}

func Float32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "float"}
}

func Int32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int32"}
}

func Int64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int64"}
}

func DateProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date"}
}
func DateTimeProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date-time"}
}
func MapProperty(property *Schema) *Schema {
	return &Schema{Type: &StringOrArray{Single: "object"}, AdditionalProperties: property}
}
func RefProperty(name string) *Schema {
	return &Schema{Ref: name}
}

func ArrayProperty(items *Schema) *Schema {
	return &Schema{Items: &SchemaOrArray{Single: items}, Type: &StringOrArray{Single: "array"}}
}

type Schema struct {
	Ref                  string                 `structs:"-"`
	Description          string                 `structs:"description,omitempty"`
	Maximum              float64                `structs:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `structs:"exclusiveMaximum,omitempty"`
	Minimum              float64                `structs:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `structs:"exclusiveMinimum,omitempty"`
	MaxLength            int64                  `structs:"maxLength,omitempty"`
	MinLength            int64                  `structs:"minLength,omitempty"`
	Pattern              string                 `structs:"pattern,omitempty"`
	MaxItems             int64                  `structs:"maxItems,omitempty"`
	MinItems             int64                  `structs:"minItems,omitempty"`
	UniqueItems          bool                   `structs:"uniqueItems,omitempty"`
	MultipleOf           float64                `structs:"multipleOf,omitempty"`
	Enum                 []interface{}          `structs:"enum,omitempty"`
	Type                 *StringOrArray         `structs:"-"`
	Format               string                 `structs:"format,omitempty"`
	Title                string                 `structs:"title,omitempty"`
	Default              interface{}            `structs:"default,omitempty"`
	MaxProperties        int64                  `structs:"maxProperties,omitempty"`
	MinProperties        int64                  `structs:"minProperties,omitempty"`
	Required             []string               `structs:"required,omitempty"`
	Items                *SchemaOrArray         `structs:"-"`
	AllOf                []Schema               `structs:"-"`
	Properties           map[string]Schema      `structs:"-"`
	Discriminator        string                 `structs:"discriminator,omitempty"`
	ReadOnly             bool                   `structs:"readOnly,omitempty"`
	XML                  *XMLObject             `structs:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `structs:"externalDocs,omitempty"`
	Example              interface{}            `structs:"example,omitempty"`
	AdditionalProperties *Schema                `structs:"-"`
}

func (s Schema) Map() map[string]interface{} {
	if s.Ref != "" {
		return map[string]interface{}{"$ref": s.Ref}
	}
	res := structs.Map(s)

	if len(s.AllOf) > 0 {
		var ser []map[string]interface{}
		for _, sch := range s.AllOf {
			ser = append(ser, sch.Map())
		}
		res["allOf"] = ser
	}

	if len(s.Properties) > 0 {
		ser := make(map[string]interface{})
		for k, v := range s.Properties {
			ser[k] = v.Map()
		}
		res["properties"] = ser
	}
	if s.AdditionalProperties != nil {
		res["additionalProperties"] = s.AdditionalProperties.Map()
	}

	if s.Type != nil {
		var value interface{} = s.Type.Multi
		if s.Type.Single != "" && len(s.Type.Multi) == 0 {
			value = s.Type.Single
		}
		res["type"] = value
	}

	if s.Items != nil {
		var value interface{} = s.Items.Multi
		if len(s.Items.Multi) == 0 && s.Items.Single != nil {
			value = s.Items.Single
		}
		res["items"] = value
	}

	return res
}

func (s Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s Schema) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}

