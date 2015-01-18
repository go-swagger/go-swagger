package validate

import (
	"fmt"
	"reflect"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
)

type schemaSliceValidator struct {
	Path            string
	In              string
	MaxItems        *int64
	MinItems        *int64
	UniqueItems     bool
	AdditionalItems *spec.SchemaOrBool
	Items           *spec.SchemaOrArray
	Root            interface{}
}

func (s *schemaSliceValidator) SetPath(path string) {
	s.Path = path
}

func (s *schemaSliceValidator) Applies(source interface{}, kind reflect.Kind) bool {
	_, ok := source.(*spec.Schema)
	return ok && kind == reflect.Slice
}

func (s *schemaSliceValidator) Validate(data interface{}) *Result {
	result := new(Result)
	val := data.([]interface{})
	size := int64(len(val))

	if s.Items != nil && s.Items.Schema != nil {
		for i, value := range val {
			validator := newSchemaValidator(s.Items.Schema, s.Root, fmt.Sprintf("%s.%d", s.Path, i))
			result.Merge(validator.Validate(value))
		}
	}

	itemsSize := int64(0)
	if s.Items != nil && len(s.Items.Schemas) > 0 {
		itemsSize = int64(len(s.Items.Schemas))
		for i := int64(0); i < itemsSize; i++ {
			validator := newSchemaValidator(&s.Items.Schemas[i], s.Root, fmt.Sprintf("%s.%d", s.Path, i))
			result.Merge(validator.Validate(val[i]))
		}

	}
	if s.AdditionalItems != nil && itemsSize < size {
		if s.Items != nil && (s.Items.Schema != nil || len(s.Items.Schemas) > 0) && !s.AdditionalItems.Allows {
			result.AddErrors(errors.New(422, "array doesn't allow for additional items"))
		}
		if s.AdditionalItems.Schema != nil {
			for i := itemsSize; i < (size-itemsSize)+1; i++ {
				validator := newSchemaValidator(s.AdditionalItems.Schema, s.Root, fmt.Sprintf("%s.%d", s.Path, i))
				result.Merge(validator.Validate(val[i]))
			}
		}
	}

	if s.MinItems != nil && size < *s.MinItems {
		result.AddErrors(errors.TooFewItems(s.Path, s.In, *s.MinItems))
	}
	if s.MaxItems != nil && size > *s.MaxItems {
		result.AddErrors(errors.TooManyItems(s.Path, s.In, *s.MaxItems))
	}
	if s.UniqueItems && s.hasDuplicates(val, int(size)) {
		result.AddErrors(errors.DuplicateItems(s.Path, s.In))
	}
	result.Inc()
	return result
}

func (s *schemaSliceValidator) hasDuplicates(value []interface{}, size int) bool {
	var unique []interface{}
	for _, v := range value {
		for _, u := range unique {
			if reflect.DeepEqual(v, u) {
				return true
			}
		}
		unique = append(unique, v)
	}
	return false
}
