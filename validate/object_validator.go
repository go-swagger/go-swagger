package validate

import (
	"reflect"
	"regexp"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
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
}

func (o *objectValidator) SetPath(path string) {
	o.Path = path
}

func (o *objectValidator) Applies(source interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(source) == specSchemaType && kind == reflect.Map
}

func (o *objectValidator) Validate(data interface{}) *Result {
	val := data.(map[string]interface{})
	numKeys := int64(len(val))

	if o.MinProperties != nil && numKeys < *o.MinProperties {
		return sErr(errors.New(422, "must have at least %d properties", *o.MinProperties))
	}
	if o.MaxProperties != nil && numKeys > *o.MaxProperties {
		return sErr(errors.New(422, "must have at most %d properties", *o.MaxProperties))
	}

	res := new(Result)
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
			if !(regularProperty || matched) {
				res.AddErrors(errors.New(422, "%s.%s in %s is a forbidden property", o.Path, k, o.In))
			}

		}
	} else {
		for key, value := range val {
			_, regularProperty := o.Properties[key]
			matched, succeededOnce, _ := o.validatePatternProperty(key, value, res)
			if !(regularProperty || matched || succeededOnce) {
				if o.AdditionalProperties != nil && o.AdditionalProperties.Schema != nil {
					res.Merge(newSchemaValidator(o.AdditionalProperties.Schema, o.Root, o.Path+"."+key).Validate(value))
				} else if regularProperty && !(matched || succeededOnce) {
					res.AddErrors(errors.New(422, "%s.%s in %s failed all pattern properties", o.Path, key, o.In))
				}
			}
		}
	}

	for pName, pSchema := range o.Properties {
		if v, ok := val[pName]; ok {
			res.Merge(newSchemaValidator(&pSchema, o.Root, o.Path+"."+pName).Validate(v))
		}
	}

	// Pattern Properties
	res.Inc()
	return res
}

func (o *objectValidator) validatePatternProperty(key string, value interface{}, result *Result) (bool, bool, []string) {
	matched := false
	succeededOnce := false
	var patterns []string

	for k, schema := range o.PatternProperties {
		patterns = append(patterns, k)
		if match, _ := regexp.MatchString(k, key); match {
			matched = true
			validator := newSchemaValidator(&schema, o.Root, o.Path+"."+key)

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
