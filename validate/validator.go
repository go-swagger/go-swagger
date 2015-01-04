package validate

import (
	"fmt"
	"math"
	"reflect"
	"regexp"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
)

type schemaValidator struct {
	schema *spec.Schema
	parent interface{}
	path   string
	in     string
}

type itemsValidator struct {
	items  *spec.Items
	parent interface{}
	path   string
	in     string
}

func (i *itemsValidator) Validate(index int, data interface{}) *errors.Validation {
	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()

	path := fmt.Sprintf("%s.%d", i.path, index)
	validators := []valueValidator{
		i.stringValidator(path),
		i.numberValidator(path),
		i.sliceValidator(path),
		i.commonValidator(path),
	}

	for _, validator := range validators {
		if validator.Applies(i.parent, kind) {
			if err := validator.Validate(data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *itemsValidator) commonValidator(path string) valueValidator {
	return &basicCommonValidator{
		Name:    path,
		In:      i.in,
		Default: i.items.Default,
		Enum:    i.items.Enum,
	}
}

func (i *itemsValidator) sliceValidator(path string) valueValidator {
	return &basicSliceValidator{
		Name:        path,
		In:          i.in,
		Default:     i.items.Default,
		MaxItems:    i.items.MaxItems,
		MinItems:    i.items.MinItems,
		UniqueItems: i.items.UniqueItems,
		Source:      i.parent,
		Items:       i.items.Items,
	}
}

func (i *itemsValidator) numberValidator(path string) valueValidator {
	return &numberValidator{
		Name:             path,
		In:               i.in,
		Default:          i.items.Default,
		MultipleOf:       i.items.MultipleOf,
		Maximum:          i.items.Maximum,
		ExclusiveMaximum: i.items.ExclusiveMaximum,
		Minimum:          i.items.Minimum,
		ExclusiveMinimum: i.items.ExclusiveMinimum,
	}
}

func (i *itemsValidator) stringValidator(path string) valueValidator {
	return &stringValidator{
		Name:      path,
		In:        i.in,
		Default:   i.items.Default,
		MaxLength: i.items.MaxLength,
		MinLength: i.items.MinLength,
		Pattern:   i.items.Pattern,
	}
}

// a param has very limited subset of validations to apply
type paramValidator struct {
	param *spec.Parameter
	name  string
}

func (p *paramValidator) Validate(data interface{}) *errors.Validation {
	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()

	validators := []valueValidator{
		p.stringValidator(),
		p.numberValidator(),
		p.sliceValidator(),
		p.commonValidator(),
	}

	for _, validator := range validators {
		if validator.Applies(p.param, kind) {
			if err := validator.Validate(data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *paramValidator) commonValidator() valueValidator {
	return &basicCommonValidator{
		Name:    p.param.Name,
		In:      p.param.In,
		Default: p.param.Default,
		Enum:    p.param.Enum,
	}
}

type valueValidator interface {
	Applies(interface{}, reflect.Kind) bool
	Validate(interface{}) *errors.Validation
}

type basicCommonValidator struct {
	Name    string
	In      string
	Default interface{}
	Enum    []interface{}
}

func (b *basicCommonValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Schema:
		return true
	}
	return false
}

func (b *basicCommonValidator) Validate(data interface{}) *errors.Validation {
	if len(b.Enum) > 0 {
		for i := 0; i < len(b.Enum); i++ {
			if reflect.DeepEqual(b.Enum[i], data) {
				return nil
			}
		}
		return errors.EnumFail(b.Name, b.In, data, b.Enum)
	}
	return nil
}

func (p *paramValidator) sliceValidator() valueValidator {
	return &basicSliceValidator{
		Name:        p.param.Name,
		In:          p.param.In,
		Default:     p.param.Default,
		MaxItems:    p.param.MaxItems,
		MinItems:    p.param.MinItems,
		UniqueItems: p.param.UniqueItems,
		Items:       p.param.Items,
		Source:      p.param,
	}
}

func (p *paramValidator) numberValidator() valueValidator {
	return &numberValidator{
		Name:             p.param.Name,
		In:               p.param.In,
		Default:          p.param.Default,
		MultipleOf:       p.param.MultipleOf,
		Maximum:          p.param.Maximum,
		ExclusiveMaximum: p.param.ExclusiveMaximum,
		Minimum:          p.param.Minimum,
		ExclusiveMinimum: p.param.ExclusiveMinimum,
	}
}

func (p *paramValidator) stringValidator() valueValidator {
	return &stringValidator{
		Name:      p.param.Name,
		In:        p.param.In,
		Default:   p.param.Default,
		Required:  p.param.Required,
		MaxLength: p.param.MaxLength,
		MinLength: p.param.MinLength,
		Pattern:   p.param.Pattern,
	}
}

type basicSliceValidator struct {
	Name        string
	In          string
	Default     interface{}
	MaxItems    *int64
	MinItems    *int64
	UniqueItems bool
	Items       *spec.Items
	Source      interface{}
}

func (s *basicSliceValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Items, *spec.Schema:
		return kind == reflect.Slice
	}
	return false
}

func (s *basicSliceValidator) Validate(data interface{}) *errors.Validation {
	val := reflect.ValueOf(data) // YOLO: just going to assume this is an array
	if val.Kind() != reflect.Slice {
		return nil // no business to do for this thing
	}
	size := int64(val.Len())
	if s.MinItems != nil && size < *s.MinItems {
		return errors.TooFewItems(s.Name, s.In, *s.MinItems)
	}
	if s.MaxItems != nil && size > *s.MaxItems {
		return errors.TooManyItems(s.Name, s.In, *s.MaxItems)
	}
	if s.UniqueItems && s.hasDuplicates(val, int(size)) {
		return errors.DuplicateItems(s.Name, s.In)
	}
	if s.Items != nil {
		for i := 0; i < int(size); i++ {
			ele := val.Index(i)
			validator := itemsValidator{
				items:  s.Items,
				parent: s.Source,
				in:     s.In,
				path:   s.Name,
			}
			if err := validator.Validate(i, ele.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *basicSliceValidator) hasDuplicates(value reflect.Value, size int) bool {
	dict := make(map[interface{}]struct{})
	for i := 0; i < size; i++ {
		ele := value.Index(i)
		if _, ok := dict[ele.Interface()]; ok {
			return true
		}
		dict[ele.Interface()] = struct{}{}
	}
	return false
}

type numberValidator struct {
	Name             string
	In               string
	Default          interface{}
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
}

func (n *numberValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Schema:
		isInt := kind >= reflect.Int && kind <= reflect.Uint64
		isFloat := kind == reflect.Float32 || kind == reflect.Float64
		return isInt || isFloat
	}
	return false
}

func (n *numberValidator) convertToFloat(val interface{}) float64 {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}

func (n *numberValidator) Validate(val interface{}) *errors.Validation {
	data := n.convertToFloat(val)
	if n.MultipleOf != nil && math.Mod(data, *n.MultipleOf) != 0 {
		return errors.NotMultipleOf(n.Name, n.In, *n.MultipleOf)
	}
	if n.Maximum != nil {
		max := *n.Maximum
		if n.ExclusiveMaximum {
			max--
		}
		if max < data {
			return errors.ExceedsMaximum(n.Name, n.In, *n.Maximum, n.ExclusiveMaximum)
		}
	}
	if n.Minimum != nil {
		min := *n.Minimum
		if n.ExclusiveMinimum {
			min++
		}
		if min > data {
			return errors.ExceedsMinimum(n.Name, n.In, *n.Minimum, n.ExclusiveMinimum)
		}
	}
	return nil
}

type stringValidator struct {
	Default   interface{}
	Required  bool
	MaxLength *int64
	MinLength *int64
	Pattern   string
	Name      string
	In        string
}

func (s *stringValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Schema:
		return kind == reflect.String
	}
	return false
}

func (s *stringValidator) Validate(val interface{}) *errors.Validation {
	data := val.(string)
	if s.Required && s.Default == nil && data == "" {
		return errors.Required(s.Name, s.In)
	}

	if s.MaxLength != nil && int64(len(data)) > *s.MaxLength {
		return errors.TooLong(s.Name, s.In, *s.MaxLength)
	}

	if s.MinLength != nil && int64(len(data)) < *s.MinLength {
		return errors.TooShort(s.Name, s.In, *s.MinLength)
	}

	if s.Pattern != "" {
		// TODO: translate back and forth from javascript syntax, peferrably cached?
		//       perhaps allow the option of using a different regex engine like pcre?
		re := regexp.MustCompile(s.Pattern)
		if !re.MatchString(data) {
			return errors.FailedPattern(s.Name, s.In, s.Pattern)
		}
	}
	return nil
}
