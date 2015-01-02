package spec

import (
	"encoding/json"
	"fmt"

	"github.com/casualjim/go-swagger/util"
)

// BooleanProperty creates a boolean property
func BooleanProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "boolean"}}}
}

// BoolProperty creates a boolean property
func BoolProperty() *Schema { return BooleanProperty() }

// StringProperty creates a string property
func StringProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "string"}}}
}

// Float64Property creates a float64/double property
func Float64Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "number"}, Format: "double"}}
}

// Float32Property creates a float32/float property
func Float32Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "number"}, Format: "float"}}
}

// Int32Property creates an int32 property
func Int32Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "number"}, Format: "int32"}}
}

// Int64Property creates an int64 property
func Int64Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "number"}, Format: "int64"}}
}

// DateProperty creates an date property
func DateProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "string"}, Format: "date"}}
}

// DateTimeProperty creates a date time property
func DateTimeProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "string"}, Format: "date-time"}}
}

// MapProperty creates a map property
func MapProperty(property *Schema) *Schema {
	return &Schema{schemaProps: schemaProps{Type: &StringOrArray{Single: "object"}, AdditionalProperties: property}}
}

// RefProperty creates a ref property
func RefProperty(name string) *Schema {
	return &Schema{refable: refable{Ref: Ref(name)}}
}

// ArrayProperty creates an array property
func ArrayProperty(items *Schema) *Schema {
	return &Schema{schemaProps: schemaProps{Items: &SchemaOrArray{Single: items}, Type: &StringOrArray{Single: "array"}}}
}

type schemaProps struct {
	Description          string                 `json:"description,omitempty"`
	Type                 *StringOrArray         `json:"type,omitempty,byValue"`
	Format               string                 `json:"format,omitempty"`
	Title                string                 `json:"title,omitempty"`
	Default              interface{}            `json:"default,omitempty"`
	MaxProperties        *int64                 `json:"maxProperties,omitempty"`
	MinProperties        *int64                 `json:"minProperties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Items                *SchemaOrArray         `json:"items,omitempty,byValue"`
	AllOf                []Schema               `json:"allOf,omitempty"`
	Properties           map[string]Schema      `json:"properties,omitempty"`
	Discriminator        string                 `json:"discriminator,omitempty"`
	ReadOnly             bool                   `json:"readOnly,omitempty"`
	XML                  *XMLObject             `json:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `json:"externalDocs,omitempty"`
	Example              interface{}            `json:"example,omitempty"`
	AdditionalProperties *Schema                `json:"additionalProperties,omitempty"`
}

// Schema the schema object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is based on the [JSON Schema Specification Draft 4](http://json-schema.org/)
// and uses a predefined subset of it.
// On top of this subset, there are extensions provided by this specification to allow for more complete documentation.
//
// For more information: http://goo.gl/8us55a#schemaObject
type Schema struct {
	refable
	vendorExtensible
	commonValidations
	schemaProps
}

// MarshalJSON marshal this to JSON
func (s Schema) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(s.schemaProps)
	if err != nil {
		return nil, fmt.Errorf("schema props %v", err)
	}
	b2, err := json.Marshal(s.vendorExtensible)
	if err != nil {
		return nil, fmt.Errorf("vendor props %v", err)
	}
	b3, err := json.Marshal(s.refable)
	if err != nil {
		return nil, fmt.Errorf("ref prop %v", err)
	}
	b4, err := json.Marshal(s.commonValidations)
	if err != nil {
		return nil, fmt.Errorf("common validations %v", err)
	}
	return util.ConcatJSON(b1, b2, b3, b4), nil
}

// UnmarshalJSON marshal this from JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.schemaProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.vendorExtensible); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.refable); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.commonValidations); err != nil {
		return err
	}
	return nil
}
