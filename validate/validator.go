package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
)

type valueValidator interface {
	SetPath(path string)
	Applies(interface{}, reflect.Kind) bool
	Validate(interface{}) *result
}

type result struct {
	Errors     []errors.Error
	MatchCount int
}

func (r *result) Merge(other *result) *result {
	if other == nil {
		return r
	}
	r.AddErrors(other.Errors...)
	r.MatchCount += other.MatchCount
	return r
}

func (r *result) AddErrors(errors ...errors.Error) {
	r.Errors = append(r.Errors, errors...)
}

func (r *result) IsValid() bool {
	return len(r.Errors) == 0
}

func (r *result) HasErrors() bool {
	return !r.IsValid()
}

func (r *result) Inc() {
	r.MatchCount++
}

type itemsValidator struct {
	items      *spec.Items
	root       interface{}
	path       string
	in         string
	validators []valueValidator
}

func newItemsValidator(path, in string, items *spec.Items, root interface{}) *itemsValidator {
	iv := &itemsValidator{path: path, in: in, items: items, root: root}
	iv.validators = []valueValidator{
		iv.stringValidator(),
		iv.numberValidator(),
		iv.sliceValidator(),
		iv.commonValidator(),
	}
	return iv
}

func (i *itemsValidator) Validate(index int, data interface{}) *result {
	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()
	mainResult := &result{}
	path := fmt.Sprintf("%s.%d", i.path, index)

	for _, validator := range i.validators {
		validator.SetPath(path)
		if validator.Applies(i.root, kind) {
			result := validator.Validate(data)
			mainResult.Merge(result)
			mainResult.Inc()
			if result != nil && result.HasErrors() {
				return mainResult
			}
		}
	}
	return mainResult
}

func (i *itemsValidator) commonValidator() valueValidator {
	return &basicCommonValidator{
		In:      i.in,
		Default: i.items.Default,
		Enum:    i.items.Enum,
	}
}

func (i *itemsValidator) sliceValidator() valueValidator {
	return &basicSliceValidator{
		In:          i.in,
		Default:     i.items.Default,
		MaxItems:    i.items.MaxItems,
		MinItems:    i.items.MinItems,
		UniqueItems: i.items.UniqueItems,
		Source:      i.root,
		Items:       i.items.Items,
	}
}

func (i *itemsValidator) numberValidator() valueValidator {
	return &numberValidator{
		In:               i.in,
		Default:          i.items.Default,
		MultipleOf:       i.items.MultipleOf,
		Maximum:          i.items.Maximum,
		ExclusiveMaximum: i.items.ExclusiveMaximum,
		Minimum:          i.items.Minimum,
		ExclusiveMinimum: i.items.ExclusiveMinimum,
	}
}

func (i *itemsValidator) stringValidator() valueValidator {
	return &stringValidator{
		In:        i.in,
		Default:   i.items.Default,
		MaxLength: i.items.MaxLength,
		MinLength: i.items.MinLength,
		Pattern:   i.items.Pattern,
	}
}

// a param has very limited subset of validations to apply
type paramValidator struct {
	param      *spec.Parameter
	validators []valueValidator
}

func newParamValidator(param *spec.Parameter) *paramValidator {
	p := &paramValidator{param: param}
	p.validators = []valueValidator{
		p.stringValidator(),
		p.numberValidator(),
		p.sliceValidator(),
		p.commonValidator(),
	}
	return p
}

func (p *paramValidator) Validate(data interface{}) *result {
	result := &result{}
	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()

	for _, validator := range p.validators {
		if validator.Applies(p.param, kind) {
			if err := validator.Validate(data); err != nil {
				result.Merge(err)
				if err.HasErrors() {
					return result
				}
			}
		}
	}
	return nil
}

func (p *paramValidator) commonValidator() valueValidator {
	return &basicCommonValidator{
		Path:    p.param.Name,
		In:      p.param.In,
		Default: p.param.Default,
		Enum:    p.param.Enum,
	}
}

type basicCommonValidator struct {
	Path    string
	In      string
	Default interface{}
	Enum    []interface{}
}

func (b *basicCommonValidator) SetPath(path string) {
	b.Path = path
}

func (b *basicCommonValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Schema:
		return true
	}
	return false
}

func (b *basicCommonValidator) Validate(data interface{}) (res *result) {
	if len(b.Enum) > 0 {

		for _, enumValue := range b.Enum {
			if data != nil && reflect.DeepEqual(enumValue, data) {
				return nil
			}
		}
		return &result{Errors: []errors.Error{errors.EnumFail(b.Path, b.In, data, b.Enum)}}
	}
	return nil
}

func (p *paramValidator) sliceValidator() valueValidator {
	return &basicSliceValidator{
		Path:        p.param.Name,
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
		Path:             p.param.Name,
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
		Path:      p.param.Name,
		In:        p.param.In,
		Default:   p.param.Default,
		Required:  p.param.Required,
		MaxLength: p.param.MaxLength,
		MinLength: p.param.MinLength,
		Pattern:   p.param.Pattern,
	}
}

type basicSliceValidator struct {
	Path           string
	In             string
	Default        interface{}
	MaxItems       *int64
	MinItems       *int64
	UniqueItems    bool
	Items          *spec.Items
	Source         interface{}
	itemsValidator *itemsValidator
}

func (s *basicSliceValidator) SetPath(path string) {
	s.Path = path
}

func (s *basicSliceValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Items:
		return kind == reflect.Slice
	}
	return false
}

func sErr(err errors.Error) *result {
	return &result{Errors: []errors.Error{err}}
}

func (s *basicSliceValidator) Validate(data interface{}) *result {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil // no business to do for this thing
	}

	size := int64(val.Len())
	if s.MinItems != nil && size < *s.MinItems {
		return sErr(errors.TooFewItems(s.Path, s.In, *s.MinItems))
	}
	if s.MaxItems != nil && size > *s.MaxItems {
		return sErr(errors.TooManyItems(s.Path, s.In, *s.MaxItems))
	}
	if s.UniqueItems && s.hasDuplicates(val, int(size)) {
		return sErr(errors.DuplicateItems(s.Path, s.In))
	}
	if s.itemsValidator == nil && s.Items != nil {
		s.itemsValidator = newItemsValidator(s.Path, s.In, s.Items, s.Source)
	}
	if s.itemsValidator != nil {
		for i := 0; i < int(size); i++ {
			ele := val.Index(i)
			if err := s.itemsValidator.Validate(i, ele.Interface()); err != nil && err.HasErrors() {
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
	Path             string
	In               string
	Default          interface{}
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
}

func (n *numberValidator) SetPath(path string) {
	n.Path = path
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

func (n *numberValidator) Validate(val interface{}) *result {
	data := n.convertToFloat(val)
	if n.MultipleOf != nil && !isFloat64AnInteger(data / *n.MultipleOf) {
		return sErr(errors.NotMultipleOf(n.Path, n.In, *n.MultipleOf))
	}
	if n.Maximum != nil {
		max := *n.Maximum
		if (!n.ExclusiveMaximum && data > max) || (n.ExclusiveMaximum && data >= max) {
			return sErr(errors.ExceedsMaximum(n.Path, n.In, *n.Maximum, n.ExclusiveMaximum))
		}
	}
	if n.Minimum != nil {
		min := *n.Minimum
		if (!n.ExclusiveMinimum && data < min) || (n.ExclusiveMinimum && data <= min) {
			return sErr(errors.ExceedsMinimum(n.Path, n.In, *n.Minimum, n.ExclusiveMinimum))
		}
	}
	return &result{}
}

type stringValidator struct {
	Default   interface{}
	Required  bool
	MaxLength *int64
	MinLength *int64
	Pattern   string
	Path      string
	In        string
}

func (s *stringValidator) SetPath(path string) {
	s.Path = path
}

func (s *stringValidator) Applies(source interface{}, kind reflect.Kind) bool {
	switch source.(type) {
	case *spec.Parameter, *spec.Schema:
		return kind == reflect.String
	}
	return false
}

func (s *stringValidator) Validate(val interface{}) *result {
	data := val.(string)
	if s.Required && s.Default == nil && data == "" {
		return sErr(errors.Required(s.Path, s.In))
	}
	strLen := int64(utf8.RuneCount([]byte(data)))
	if s.MaxLength != nil && strLen > *s.MaxLength {
		return sErr(errors.TooLong(s.Path, s.In, *s.MaxLength))
	}

	if s.MinLength != nil && strLen < *s.MinLength {
		return sErr(errors.TooShort(s.Path, s.In, *s.MinLength))
	}

	if s.Pattern != "" {
		// TODO: translate back and forth from javascript syntax, peferrably cached?
		//       perhaps allow the option of using a different regex engine like pcre?
		re := regexp.MustCompile(s.Pattern)
		if !re.MatchString(data) {
			return sErr(errors.FailedPattern(s.Path, s.In, s.Pattern))
		}
	}
	return nil
}
