package validate

import (
	"reflect"
	"regexp"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
	"github.com/casualjim/go-swagger/validate"
)

type objectValidator struct {
	Path                 string
	In                   string
	MaxProperties        *int64
	MinProperties        *int64
	Required             []string
	Properties           map[string]spec.Schema
	AdditionalProperties *spec.SchemaOrBool
	PatternProperties    map[string]spec.Schema
	Root                 interface{}
	KnownFormats         strfmt.Registry
}

func (o *objectValidator) SetPath(path string) {
	o.Path = path
}

func (o *objectValidator) Applies(source interface{}, kind reflect.Kind) bool {
	// TODO: this should also work for structs
	// there is a problem in the type validator where it will be unhappy about null values
	// so that requires more testing
	r := reflect.TypeOf(source) == specSchemaType && (kind == reflect.Map || kind == reflect.Struct)
	// fmt.Printf("object validator for %q applies %t for %T (kind: %v)\n", o.Path, r, source, kind)
	return r
}

func (o *objectValidator) Validate(data interface{}) *validate.Result {
	val := data.(map[string]interface{})
	numKeys := int64(len(val))

	if o.MinProperties != nil && numKeys < *o.MinProperties {
		return sErr(errors.New(422, "must have at least %d properties", *o.MinProperties))
	}
	if o.MaxProperties != nil && numKeys > *o.MaxProperties {
		return sErr(errors.New(422, "must have at most %d properties", *o.MaxProperties))
	}

	res := new(validate.Result)
	if len(o.Required) > 0 {
		for _, k := range o.Required {
			if _, ok := val[k]; !ok {
				res.AddErrors(errors.Required(o.Path+"."+k, o.In))
				continue
			}
			res.Inc()
		}
	}
	if o.AdditionalProperties != nil && !o.AdditionalProperties.Allows {
		for k := range val {
			_, regularProperty := o.Properties[k]
			matched := false

			for pk := range o.PatternProperties {
				if matches, _ := regexp.MatchString(pk, k); matches {
					matched = true
					break
				}
			}
			if !(regularProperty || k == "$schema" || k == "id" || matched) {
				res.AddErrors(errors.New(422, "%s.%s in %s is a forbidden property", o.Path, k, o.In))
			}

		}
	} else {
		for key, value := range val {
			_, regularProperty := o.Properties[key]
			matched, succeededOnce, _ := o.validatePatternProperty(key, value, res)
			if !(regularProperty || matched || succeededOnce) {
				if o.AdditionalProperties != nil && o.AdditionalProperties.Schema != nil {
					res.Merge(NewSchemaValidator(o.AdditionalProperties.Schema, o.Root, o.Path+"."+key, o.KnownFormats).Validate(value))
				} else if regularProperty && !(matched || succeededOnce) {
					res.AddErrors(errors.New(422, "%s.%s in %s failed all pattern properties", o.Path, key, o.In))
				}
			}
		}
	}

	for pName, pSchema := range o.Properties {
		rName := pName
		if o.Path != "" {
			rName = o.Path + "." + pName
		}
		if v, ok := val[pName]; ok {
			res.Merge(NewSchemaValidator(&pSchema, o.Root, rName, o.KnownFormats).Validate(v))
		}
	}

	// Pattern Properties
	res.Inc()
	return res
}

func (o *objectValidator) validatePatternProperty(key string, value interface{}, result *validate.Result) (bool, bool, []string) {
	matched := false
	succeededOnce := false
	var patterns []string

	for k, schema := range o.PatternProperties {
		patterns = append(patterns, k)
		if match, _ := regexp.MatchString(k, key); match {
			matched = true
			validator := NewSchemaValidator(&schema, o.Root, o.Path+"."+key, o.KnownFormats)

			res := validator.Validate(value)
			result.Merge(res)
			if res.IsValid() {
				succeededOnce = true
			}
		}
	}

	if succeededOnce {
		result.Inc()
	}

	return matched, succeededOnce, patterns
}
