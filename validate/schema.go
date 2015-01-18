package validate

import (
	"reflect"

	"github.com/casualjim/go-swagger/spec"
)

var specSchemaType = reflect.TypeOf(&spec.Schema{})

// WithSchema validates the specified data with the provided schema, when no schema
// is provided it uses the json schema as default
func WithSchema(schema *spec.Schema, data interface{}) *Result {
	return newSchemaValidator(schema, nil, "(root)").Validate(data)
}

// like param validator but for a full json schema
type schemaValidator struct {
	Path       string
	in         string
	Schema     *spec.Schema
	validators []valueValidator
	Root       interface{}
}

func newSchemaValidator(schema *spec.Schema, rootSchema interface{}, root string) *schemaValidator {
	if rootSchema == nil {
		rootSchema = schema
	}
	if schema.ID != "" || schema.Ref.String() != "" || schema.Ref.IsRoot() {
		err := spec.ExpandSchema(schema, rootSchema, nil)
		if err != nil {
			panic(err)
		}
	}
	s := schemaValidator{Path: root, in: "body", Schema: schema, Root: rootSchema}

	s.validators = []valueValidator{
		s.typeValidator(),
		s.schemaValidator(),
		s.stringValidator(),
		s.numberValidator(),
		s.sliceValidator(),
		s.commonValidator(),
		s.objectValidator(),
	}
	return &s
}

func (s *schemaValidator) Validate(data interface{}) *Result {
	if data == nil {
		v := s.validators[0].Validate(data)
		v.Merge(s.validators[5].Validate(data))
		return v
	}
	result := &Result{}

	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()

	for _, v := range s.validators {
		if !v.Applies(s.Schema, kind) {
			continue
		}

		err := v.Validate(data)
		result.Merge(err)
		result.Inc()
	}
	result.Inc()
	return result
}

func (s *schemaValidator) typeValidator() valueValidator {
	return &typeValidator{Type: s.Schema.Type, Format: s.Schema.Format, In: s.in, Path: s.Path}
}

func (s *schemaValidator) commonValidator() valueValidator {
	return &basicCommonValidator{
		Path:    s.Path,
		In:      s.in,
		Default: s.Schema.Default,
		Enum:    s.Schema.Enum,
	}
}

func (s *schemaValidator) sliceValidator() valueValidator {
	return &schemaSliceValidator{
		Path:            s.Path,
		In:              s.in,
		MaxItems:        s.Schema.MaxItems,
		MinItems:        s.Schema.MinItems,
		UniqueItems:     s.Schema.UniqueItems,
		AdditionalItems: s.Schema.AdditionalItems,
		Items:           s.Schema.Items,
		Root:            s.Root,
	}
}

func (s *schemaValidator) numberValidator() valueValidator {
	return &numberValidator{
		Path:             s.Path,
		In:               s.in,
		Default:          s.Schema.Default,
		MultipleOf:       s.Schema.MultipleOf,
		Maximum:          s.Schema.Maximum,
		ExclusiveMaximum: s.Schema.ExclusiveMaximum,
		Minimum:          s.Schema.Minimum,
		ExclusiveMinimum: s.Schema.ExclusiveMinimum,
	}
}

func (s *schemaValidator) stringValidator() valueValidator {
	return &stringValidator{
		Path:      s.Path,
		In:        s.in,
		Default:   s.Schema.Default,
		MaxLength: s.Schema.MaxLength,
		MinLength: s.Schema.MinLength,
		Pattern:   s.Schema.Pattern,
	}
}

func (s *schemaValidator) schemaValidator() valueValidator {
	sch := s.Schema
	return newSchemaPropsValidator(s.Path, s.in, sch.AllOf, sch.OneOf, sch.AnyOf, sch.Not, sch.Dependencies, s.Root)
}

func (s *schemaValidator) objectValidator() valueValidator {
	return &objectValidator{
		Path:                 s.Path,
		In:                   s.in,
		MaxProperties:        s.Schema.MaxProperties,
		MinProperties:        s.Schema.MinProperties,
		Required:             s.Schema.Required,
		Properties:           s.Schema.Properties,
		AdditionalProperties: s.Schema.AdditionalProperties,
		PatternProperties:    s.Schema.PatternProperties,
		Root:                 s.Root,
	}
}
