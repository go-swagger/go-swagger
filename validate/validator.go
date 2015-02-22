package validate

import (
	"fmt"
	"reflect"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

type entityValidator interface {
	Validate(interface{}) *Result
}

type valueValidator interface {
	SetPath(path string)
	Applies(interface{}, reflect.Kind) bool
	Validate(interface{}) *Result
}

// Result represents a validation result
type Result struct {
	Errors     []errors.Error
	MatchCount int
}

// Merge merges this result with the other one, preserving match counts etc
func (r *Result) Merge(other *Result) *Result {
	if other == nil {
		return r
	}
	r.AddErrors(other.Errors...)
	r.MatchCount += other.MatchCount
	return r
}

// AddErrors adds errors to this validation result
func (r *Result) AddErrors(errors ...errors.Error) {
	r.Errors = append(r.Errors, errors...)
}

// IsValid returns true when this result is valid
func (r *Result) IsValid() bool {
	return len(r.Errors) == 0
}

// HasErrors returns true when this result is invalid
func (r *Result) HasErrors() bool {
	return !r.IsValid()
}

// Inc increments the match count
func (r *Result) Inc() {
	r.MatchCount++
}

type itemsValidator struct {
	items        *spec.Items
	root         interface{}
	path         string
	in           string
	validators   []valueValidator
	KnownFormats strfmt.Registry
}

func newItemsValidator(path, in string, items *spec.Items, root interface{}, formats strfmt.Registry) *itemsValidator {
	iv := &itemsValidator{path: path, in: in, items: items, root: root, KnownFormats: formats}
	iv.validators = []valueValidator{
		iv.stringValidator(),
		iv.formatValidator(),
		iv.numberValidator(),
		iv.sliceValidator(),
		iv.commonValidator(),
	}
	return iv
}

func (i *itemsValidator) Validate(index int, data interface{}) *Result {
	tpe := reflect.TypeOf(data)
	kind := tpe.Kind()
	mainResult := &Result{}
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
		In:           i.in,
		Default:      i.items.Default,
		MaxItems:     i.items.MaxItems,
		MinItems:     i.items.MinItems,
		UniqueItems:  i.items.UniqueItems,
		Source:       i.root,
		Items:        i.items.Items,
		KnownFormats: i.KnownFormats,
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

func (i *itemsValidator) formatValidator() valueValidator {
	return &formatValidator{
		In:           i.in,
		Default:      i.items.Default,
		Format:       i.items.Format,
		KnownFormats: i.KnownFormats,
	}
}

// a param has very limited subset of validations to apply
type paramValidator struct {
	param        *spec.Parameter
	validators   []valueValidator
	KnownFormats strfmt.Registry
}

func newParamValidator(param *spec.Parameter, formats strfmt.Registry) *paramValidator {
	p := &paramValidator{param: param, KnownFormats: formats}
	p.validators = []valueValidator{
		p.stringValidator(),
		p.formatValidator(),
		p.numberValidator(),
		p.sliceValidator(),
		p.commonValidator(),
	}
	return p
}

func (p *paramValidator) Validate(data interface{}) *Result {
	result := &Result{}
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

func (b *basicCommonValidator) Validate(data interface{}) (res *Result) {
	if len(b.Enum) > 0 {
		for _, enumValue := range b.Enum {
			if data != nil && reflect.DeepEqual(enumValue, data) {
				return nil
			}
		}
		return &Result{Errors: []errors.Error{errors.EnumFail(b.Path, b.In, data, b.Enum)}}
	}
	return nil
}

func (p *paramValidator) sliceValidator() valueValidator {
	return &basicSliceValidator{
		Path:         p.param.Name,
		In:           p.param.In,
		Default:      p.param.Default,
		MaxItems:     p.param.MaxItems,
		MinItems:     p.param.MinItems,
		UniqueItems:  p.param.UniqueItems,
		Items:        p.param.Items,
		Source:       p.param,
		KnownFormats: p.KnownFormats,
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

func (p *paramValidator) formatValidator() valueValidator {
	return &formatValidator{
		Path:         p.param.Name,
		In:           p.param.In,
		Default:      p.param.Default,
		Format:       p.param.Format,
		KnownFormats: p.KnownFormats,
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
	KnownFormats   strfmt.Registry
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

func sErr(err errors.Error) *Result {
	return &Result{Errors: []errors.Error{err}}
}

func (s *basicSliceValidator) Validate(data interface{}) *Result {
	val := reflect.ValueOf(data)

	size := int64(val.Len())
	if s.MinItems != nil {
		if err := MinItems(s.Path, s.In, size, *s.MinItems); err != nil {
			return sErr(err)
		}
	}

	if s.MaxItems != nil {
		if err := MaxItems(s.Path, s.In, size, *s.MaxItems); err != nil {
			return sErr(err)
		}
	}

	if s.UniqueItems {
		if err := UniqueItems(s.Path, s.In, data); err != nil {
			return sErr(err)
		}
	}

	if s.itemsValidator == nil && s.Items != nil {
		s.itemsValidator = newItemsValidator(s.Path, s.In, s.Items, s.Source, s.KnownFormats)
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
	case *spec.Parameter, *spec.Schema, *spec.Items:
		isInt := kind >= reflect.Int && kind <= reflect.Uint64
		isFloat := kind == reflect.Float32 || kind == reflect.Float64
		r := isInt || isFloat
		// fmt.Printf("schema props validator for %q applies %t for %T (kind: %v)\n", n.Path, r, source, kind)
		return r
	}
	// fmt.Printf("schema props validator for %q applies %t for %T (kind: %v)\n", n.Path, false, source, kind)
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

func (n *numberValidator) Validate(val interface{}) *Result {
	data := n.convertToFloat(val)

	if n.MultipleOf != nil {
		if err := MultipleOf(n.Path, n.In, data, *n.MultipleOf); err != nil {
			return sErr(err)
		}
	}
	if n.Maximum != nil {
		if err := Maximum(n.Path, n.In, data, *n.Maximum, n.ExclusiveMaximum); err != nil {
			return sErr(err)
		}
	}
	if n.Minimum != nil {
		if err := Minimum(n.Path, n.In, data, *n.Minimum, n.ExclusiveMinimum); err != nil {
			return sErr(err)
		}
	}
	result := &Result{}
	result.Inc()
	return result
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
	case *spec.Parameter, *spec.Schema, *spec.Items:
		r := kind == reflect.String
		// fmt.Printf("string validator for %q applies %t for %T (kind: %v)\n", s.Path, r, source, kind)
		return r
	}
	// fmt.Printf("string validator for %q applies %t for %T (kind: %v)\n", s.Path, false, source, kind)
	return false
}

func (s *stringValidator) Validate(val interface{}) *Result {
	data := val.(string)

	if s.Required && (s.Default == nil || s.Default == "") {
		if err := RequiredString(s.Path, s.In, data); err != nil {
			return sErr(err)
		}
	}

	if s.MaxLength != nil {
		if err := MaxLength(s.Path, s.In, data, *s.MaxLength); err != nil {
			return sErr(err)
		}
	}

	if s.MinLength != nil {
		if err := MinLength(s.Path, s.In, data, *s.MinLength); err != nil {
			return sErr(err)
		}
	}

	if s.Pattern != "" {
		if err := Pattern(s.Path, s.In, data, s.Pattern); err != nil {
			return sErr(err)
		}
	}
	return nil
}
