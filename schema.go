package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
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

// Float32Property creates a float32/float property
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
	Ref                  string                 `swagger:"-"`
	Description          string                 `swagger:"description,omitempty"`
	Maximum              *float64               `swagger:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `swagger:"exclusiveMaximum,omitempty"`
	Minimum              *float64               `swagger:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `swagger:"exclusiveMinimum,omitempty"`
	MaxLength            *int64                 `swagger:"maxLength,omitempty"`
	MinLength            *int64                 `swagger:"minLength,omitempty"`
	Pattern              string                 `swagger:"pattern,omitempty"`
	MaxItems             *int64                 `swagger:"maxItems,omitempty"`
	MinItems             *int64                 `swagger:"minItems,omitempty"`
	UniqueItems          bool                   `swagger:"uniqueItems,omitempty"`
	MultipleOf           *float64               `swagger:"multipleOf,omitempty"`
	Enum                 []interface{}          `swagger:"enum,omitempty"`
	Type                 *StringOrArray         `swagger:"type,omitempty,byValue"`
	Format               string                 `swagger:"format,omitempty"`
	Title                string                 `swagger:"title,omitempty"`
	Default              interface{}            `swagger:"default,omitempty"`
	MaxProperties        *int64                 `swagger:"maxProperties,omitempty"`
	MinProperties        *int64                 `swagger:"minProperties,omitempty"`
	Required             []string               `swagger:"required,omitempty"`
	Items                *SchemaOrArray         `swagger:"items,omitempty,byValue"`
	AllOf                []Schema               `swagger:"allOf,omitempty"`
	Properties           map[string]Schema      `swagger:"properties,omitempty"`
	Discriminator        string                 `swagger:"discriminator,omitempty"`
	ReadOnly             bool                   `swagger:"readOnly,omitempty"`
	XML                  *XMLObject             `swagger:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `swagger:"externalDocs,omitempty"`
	Example              interface{}            `swagger:"example,omitempty"`
	AdditionalProperties *Schema                `swagger:"additionalProperties,omitempty"`
	Extensions           map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
}

// UnmarshalMap hydrates this schema instance with the data from the map
func (s *Schema) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if ref, ok := dict["$ref"]; ok {
		s.Ref = ref.(string)
	}
	s.Extensions = readExtensions(dict)
	return reflection.UnmarshalMapRecursed(dict, s)
}

// UnmarshalJSON hydrates this schema instance with the data from JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return s.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this schema instance with the data from YAML
func (s *Schema) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return s.UnmarshalMap(value)
}

// MarshalMap converts this schema object to a map
func (s Schema) MarshalMap() map[string]interface{} {
	result := reflection.MarshalMapRecursed(s)
	if s.Ref != "" {
		result["$ref"] = s.Ref
	}
	addExtensions(result, s.Extensions)
	return result
}

// MarshalJSON converts this schema object to JSON
func (s Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.MarshalMap())
}

// MarshalYAML converts this schema object to YAML
func (s Schema) MarshalYAML() (interface{}, error) {
	return s.MarshalMap(), nil
}
