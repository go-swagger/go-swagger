package reflection

import (
	"reflect"
	"strings"
)

const (
	TagName = "swagger"
)

type MapMarshaller interface {
	MarshalMap() (map[string]interface{}, error)
}

type parsedTag struct {
	Name       string
	OmitEmpty  bool
	Inline     bool
	ShouldSkip bool
	ByValue    bool
}

func parseTag(tag string, name string) *parsedTag {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return &parsedTag{
			Name: name,
		}
	}
	nm := parts[0]
	if nm == "" {
		nm = name
	}
	shouldSkip := nm == "-"
	var omitEmpty, inline, byValue bool
	for _, p := range parts[1:] {
		if p == "omitempty" {
			omitEmpty = true
		}
		if p == "inline" {
			inline = true
		}
		if p == "byValue" {
			byValue = true
		}
	}
	return &parsedTag{
		Name:       nm,
		OmitEmpty:  omitEmpty,
		Inline:     inline,
		ShouldSkip: shouldSkip,
		ByValue:    byValue,
	}
}

func isStruct(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return false
	}

	return v.Kind() == reflect.Struct
}

func MarshalMap(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return nil, nil
	}

	var result map[string]interface{}

	val := reflect.ValueOf(data)
	tpe := val.Type()
	if tpe.Kind() == reflect.Ptr {
		val = val.Elem()
		tpe = tpe.Elem()
	}

	for i := 0; i < tpe.NumField(); i++ {
		fld := val.Field(i)
		fldTpe := tpe.Field(i)
		if fldTpe.PkgPath != "" {
			continue
		}

		tag := parseTag(fldTpe.Tag.Get(TagName), fldTpe.Name)
		if tag != nil && !tag.ShouldSkip {
			if result == nil {
				result = map[string]interface{}{}
			}

			if tag.OmitEmpty {
				zero := reflect.Zero(fld.Type()).Interface()
				if reflect.DeepEqual(fld.Interface(), zero) {
					continue
				}
			}

			var value interface{}
			if isStruct(fld) && !tag.ByValue {
				v, err := MarshalMap(fld.Interface())
				if err != nil {
					return nil, err
				}
				value = v
			} else {
				value = fld.Interface()
			}

			if !tag.Inline {
				result[tag.Name] = value
			}

		}
	}
	return result, nil
}

type MapUnmarshaller interface {
	UnmarshalMap(map[string]interface{}) error
}

func UnmarshalMap(data map[string]interface{}, obj interface{}) error {
	return nil
}
