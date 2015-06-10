package spec

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-swagger/go-swagger/jsonpointer"
	"github.com/go-swagger/go-swagger/swag"
)

// BooleanProperty creates a boolean property
func BooleanProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"boolean"}}}
}

// BoolProperty creates a boolean property
func BoolProperty() *Schema { return BooleanProperty() }

// StringProperty creates a string property
func StringProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"string"}}}
}

// CharProperty creates a string property
func CharProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"string"}}}
}

// Float64Property creates a float64/double property
func Float64Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"number"}, Format: "double"}}
}

// Float32Property creates a float32/float property
func Float32Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"number"}, Format: "float"}}
}

// Int8Property creates an int8 property
func Int8Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"integer"}, Format: "int8"}}
}

// Int16Property creates an int16 property
func Int16Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"integer"}, Format: "int16"}}
}

// Int32Property creates an int32 property
func Int32Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"integer"}, Format: "int32"}}
}

// Int64Property creates an int64 property
func Int64Property() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"integer"}, Format: "int64"}}
}

// StrFmtProperty creates a property for the named string format
func StrFmtProperty(format string) *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"string"}, Format: format}}
}

// DateProperty creates a date property
func DateProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"string"}, Format: "date"}}
}

// DateTimeProperty creates a date time property
func DateTimeProperty() *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"string"}, Format: "date-time"}}
}

// MapProperty creates a map property
func MapProperty(property *Schema) *Schema {
	return &Schema{schemaProps: schemaProps{Type: []string{"object"}, AdditionalProperties: &SchemaOrBool{Allows: true, Schema: property}}}
}

// RefProperty creates a ref property
func RefProperty(name string) *Schema {
	return &Schema{schemaProps: schemaProps{Ref: MustCreateRef(name)}}
}

// ArrayProperty creates an array property
func ArrayProperty(items *Schema) *Schema {
	if items == nil {
		return &Schema{schemaProps: schemaProps{Type: []string{"array"}}}
	}
	return &Schema{schemaProps: schemaProps{Items: &SchemaOrArray{Schema: items}, Type: []string{"array"}}}
}

// SchemaURL represents a schema url
type SchemaURL string

// MarshalJSON marshal this to JSON
func (r SchemaURL) MarshalJSON() ([]byte, error) {
	if r == "" {
		return []byte("{}"), nil
	}
	v := map[string]interface{}{"$schema": string(r)}
	return json.Marshal(v)
}

// UnmarshalJSON unmarshal this from JSON
func (r *SchemaURL) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	if vv, ok := v["$schema"]; ok {
		if str, ok := vv.(string); ok {
			u, err := url.Parse(str)
			if err != nil {
				return err
			}

			*r = SchemaURL(u.String())
		}
	}
	return nil
}

// type extraSchemaProps map[string]interface{}

// // JSONSchema represents a structure that is a json schema draft 04
// type JSONSchema struct {
// 	schemaProps
// 	extraSchemaProps
// }

// // MarshalJSON marshal this to JSON
// func (s JSONSchema) MarshalJSON() ([]byte, error) {
// 	b1, err := json.Marshal(s.schemaProps)
// 	if err != nil {
// 		return nil, err
// 	}
// 	b2, err := s.Ref.MarshalJSON()
// 	if err != nil {
// 		return nil, err
// 	}
// 	b3, err := s.Schema.MarshalJSON()
// 	if err != nil {
// 		return nil, err
// 	}
// 	b4, err := json.Marshal(s.extraSchemaProps)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return util.ConcatJSON(b1, b2, b3, b4), nil
// }

// // UnmarshalJSON marshal this from JSON
// func (s *JSONSchema) UnmarshalJSON(data []byte) error {
// 	var sch JSONSchema
// 	if err := json.Unmarshal(data, &sch.schemaProps); err != nil {
// 		return err
// 	}
// 	if err := json.Unmarshal(data, &sch.Ref); err != nil {
// 		return err
// 	}
// 	if err := json.Unmarshal(data, &sch.Schema); err != nil {
// 		return err
// 	}
// 	if err := json.Unmarshal(data, &sch.extraSchemaProps); err != nil {
// 		return err
// 	}
// 	*s = sch
// 	return nil
// }

type schemaProps struct {
	ID                   string            `json:"id,omitempty"`
	Ref                  Ref               `json:"-,omitempty"`
	Schema               SchemaURL         `json:"-,omitempty"`
	Description          string            `json:"description,omitempty"`
	Type                 StringOrArray     `json:"type,omitempty"`
	Format               string            `json:"format,omitempty"`
	Title                string            `json:"title,omitempty"`
	Default              interface{}       `json:"default,omitempty"`
	Maximum              *float64          `json:"maximum,omitempty"`
	ExclusiveMaximum     bool              `json:"exclusiveMaximum,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty"`
	ExclusiveMinimum     bool              `json:"exclusiveMinimum,omitempty"`
	MaxLength            *int64            `json:"maxLength,omitempty"`
	MinLength            *int64            `json:"minLength,omitempty"`
	Pattern              string            `json:"pattern,omitempty"`
	MaxItems             *int64            `json:"maxItems,omitempty"`
	MinItems             *int64            `json:"minItems,omitempty"`
	UniqueItems          bool              `json:"uniqueItems,omitempty"`
	MultipleOf           *float64          `json:"multipleOf,omitempty"`
	Enum                 []interface{}     `json:"enum,omitempty"`
	MaxProperties        *int64            `json:"maxProperties,omitempty"`
	MinProperties        *int64            `json:"minProperties,omitempty"`
	Required             []string          `json:"required,omitempty"`
	Items                *SchemaOrArray    `json:"items,omitempty"`
	AllOf                []Schema          `json:"allOf,omitempty"`
	OneOf                []Schema          `json:"oneOf,omitempty"`
	AnyOf                []Schema          `json:"anyOf,omitempty"`
	Not                  *Schema           `json:"not,omitempty"`
	Properties           map[string]Schema `json:"properties,omitempty"`
	AdditionalProperties *SchemaOrBool     `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]Schema `json:"patternProperties,omitempty"`
	Dependencies         Dependencies      `json:"dependencies,omitempty"`
	AdditionalItems      *SchemaOrBool     `json:"additionalItems,omitempty"`
	Definitions          Definitions       `json:"definitions,omitempty"`
}

type swaggerSchemaProps struct {
	Discriminator string                 `json:"discriminator,omitempty"`
	ReadOnly      bool                   `json:"readOnly,omitempty"`
	XML           *XMLObject             `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Example       interface{}            `json:"example,omitempty"`
}

// Schema the schema object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is based on the [JSON Schema Specification Draft 4](http://json-schema.org/)
// and uses a predefined subset of it.
// On top of this subset, there are extensions provided by this specification to allow for more complete documentation.
//
// For more information: http://goo.gl/8us55a#schemaObject
type Schema struct {
	vendorExtensible
	schemaProps
	swaggerSchemaProps
	ExtraProps map[string]interface{} `json:"-"`
}

// JSONLookup implements an interface to customize json pointer lookup
func (s Schema) JSONLookup(token string) (interface{}, error) {
	if ex, ok := s.Extensions[token]; ok {
		return &ex, nil
	}

	if ex, ok := s.ExtraProps[token]; ok {
		return &ex, nil
	}

	r, _, err := jsonpointer.GetForToken(s.schemaProps, token)
	if r != nil || err != nil {
		return r, err
	}
	r, _, err = jsonpointer.GetForToken(s.swaggerSchemaProps, token)
	return r, err
}

// WithProperties sets the properties for this schema
func (s *Schema) WithProperties(schemas map[string]Schema) *Schema {
	s.Properties = schemas
	return s
}

// SetProperty sets a property on this schema
func (s *Schema) SetProperty(name string, schema Schema) *Schema {
	if s.Properties == nil {
		s.Properties = make(map[string]Schema)
	}
	s.Properties[name] = schema
	return s
}

// WithAllOf sets the all of property
func (s *Schema) WithAllOf(schemas ...Schema) *Schema {
	s.AllOf = schemas
	return s
}

// WithMaxProperties sets the max number of properties an object can have
func (s *Schema) WithMaxProperties(max int64) *Schema {
	s.MaxProperties = &max
	return s
}

// WithMinProperties sets the min number of properties an object must have
func (s *Schema) WithMinProperties(min int64) *Schema {
	s.MinProperties = &min
	return s
}

// Typed sets the type of this schema for a single value item
func (s *Schema) Typed(tpe, format string) *Schema {
	s.Type = []string{tpe}
	s.Format = format
	return s
}

// AddType adds a type with potential format to the types for this schema
func (s *Schema) AddType(tpe, format string) *Schema {
	s.Type = append(s.Type, tpe)
	if format != "" {
		s.Format = format
	}
	return s
}

// CollectionOf a fluent builder method for an array parameter
func (s *Schema) CollectionOf(items Schema) *Schema {
	s.Type = []string{"array"}
	s.Items = &SchemaOrArray{Schema: &items}
	return s
}

// WithDefault sets the default value on this parameter
func (s *Schema) WithDefault(defaultValue interface{}) *Schema {
	s.Default = defaultValue
	return s
}

// WithRequired flags this parameter as required
func (s *Schema) WithRequired(items ...string) *Schema {
	s.Required = items
	return s
}

// WithMaxLength sets a max length value
func (s *Schema) WithMaxLength(max int64) *Schema {
	s.MaxLength = &max
	return s
}

// WithMinLength sets a min length value
func (s *Schema) WithMinLength(min int64) *Schema {
	s.MinLength = &min
	return s
}

// WithPattern sets a pattern value
func (s *Schema) WithPattern(pattern string) *Schema {
	s.Pattern = pattern
	return s
}

// WithMultipleOf sets a multiple of value
func (s *Schema) WithMultipleOf(number float64) *Schema {
	s.MultipleOf = &number
	return s
}

// WithMaximum sets a maximum number value
func (s *Schema) WithMaximum(max float64, exclusive bool) *Schema {
	s.Maximum = &max
	s.ExclusiveMaximum = exclusive
	return s
}

// WithMinimum sets a minimum number value
func (s *Schema) WithMinimum(min float64, exclusive bool) *Schema {
	s.Minimum = &min
	s.ExclusiveMinimum = exclusive
	return s
}

// WithEnum sets a the enum values (replace)
func (s *Schema) WithEnum(values ...interface{}) *Schema {
	s.Enum = append([]interface{}{}, values...)
	return s
}

// WithMaxItems sets the max items
func (s *Schema) WithMaxItems(size int64) *Schema {
	s.MaxItems = &size
	return s
}

// WithMinItems sets the min items
func (s *Schema) WithMinItems(size int64) *Schema {
	s.MinItems = &size
	return s
}

// UniqueValues dictates that this array can only have unique items
func (s *Schema) UniqueValues() *Schema {
	s.UniqueItems = true
	return s
}

// AllowDuplicates this array can have duplicates
func (s *Schema) AllowDuplicates() *Schema {
	s.UniqueItems = false
	return s
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
	b3, err := s.Ref.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("ref prop %v", err)
	}
	b4, err := s.Schema.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("schema prop %v", err)
	}
	b5, err := json.Marshal(s.swaggerSchemaProps)
	if err != nil {
		return nil, fmt.Errorf("common validations %v", err)
	}
	var b6 []byte
	if s.ExtraProps != nil {
		jj, err := json.Marshal(s.ExtraProps)
		if err != nil {
			return nil, fmt.Errorf("extra props %v", err)
		}
		b6 = jj
	}
	return swag.ConcatJSON(b1, b2, b3, b4, b5, b6), nil
}

// UnmarshalJSON marshal this from JSON
func (s *Schema) UnmarshalJSON(data []byte) error {
	var sch Schema
	if err := json.Unmarshal(data, &sch.schemaProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &sch.Ref); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &sch.Schema); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &sch.swaggerSchemaProps); err != nil {
		return err
	}

	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	delete(d, "$ref")
	delete(d, "$schema")
	for _, pn := range swag.DefaultJSONNameProvider.GetJSONNames(s) {
		delete(d, pn)
	}

	for k, vv := range d {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-") {
			if sch.Extensions == nil {
				sch.Extensions = map[string]interface{}{}
			}
			sch.Extensions[k] = vv
			continue
		}
		if sch.ExtraProps == nil {
			sch.ExtraProps = map[string]interface{}{}
		}
		sch.ExtraProps[k] = vv
	}

	*s = sch

	return nil
}
