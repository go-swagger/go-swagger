package validate

import (
	"reflect"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

type formatValidator struct {
	Default      interface{}
	Format       string
	Path         string
	In           string
	KnownFormats strfmt.Registry
}

func (f *formatValidator) SetPath(path string) {
	f.Path = path
}

func (f *formatValidator) Applies(source interface{}, kind reflect.Kind) bool {
	doit := func() bool {
		if source == nil {
			return false
		}
		switch source.(type) {
		case *spec.Items:
			it := source.(*spec.Items)
			return kind == reflect.String && f.KnownFormats.ContainsName(it.Format)
		case *spec.Parameter:
			par := source.(*spec.Parameter)
			return kind == reflect.String && f.KnownFormats.ContainsName(par.Format)
		case *spec.Schema:
			sch := source.(*spec.Schema)
			return kind == reflect.String && f.KnownFormats.ContainsName(sch.Format)
		}
		return false
	}
	r := doit()
	// fmt.Printf("schema props validator for %q applies %t for %T (kind: %v)\n", f.Path, r, source, kind)
	return r
}

func (f *formatValidator) Validate(val interface{}) *Result {
	result := new(Result)

	if err := FormatOf(f.Path, f.In, f.Format, val.(string), f.KnownFormats); err != nil {
		result.AddErrors(err)
	}
	result.Inc()
	return result
}
