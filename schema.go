package swagger

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// BooleanProperty creates a boolean property
func BooleanProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "boolean"}}
}

// BoolProperty creates a boolean property
func BoolProperty() *Schema { return BooleanProperty() }

// StringProperty creates a string property
func StringProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}}
}

// Float64Property creates a float64/double property
func Float64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "double"}
}

// Float32 creates a float32/float property
func Float32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "float"}
}

// Int32Property creates an int32 property
func Int32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int32"}
}

// Int64Property creates an int64 property
func Int64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int64"}
}

// DateProperty creates an date property
func DateProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date"}
}

// DateTimeProperty creates a date time property
func DateTimeProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date-time"}
}

// MapProperty creates a map property
func MapProperty(property *Schema) *Schema {
	return &Schema{Type: &StringOrArray{Single: "object"}, AdditionalProperties: property}
}

// RefProperty creates a ref property
func RefProperty(name string) *Schema {
	return &Schema{Ref: name}
}

// ArrayProperty creates an array property
func ArrayProperty(items *Schema) *Schema {
	return &Schema{Items: &SchemaOrArray{Single: items}, Type: &StringOrArray{Single: "array"}}
}

// Schema the schema object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is based on the [JSON Schema Specification Draft 4](http://json-schema.org/)
// and uses a predefined subset of it.
// On top of this subset, there are extensions provided by this specification to allow for more complete documentation.
//
// For more information: http://goo.gl/8us55a#schemaObject
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
