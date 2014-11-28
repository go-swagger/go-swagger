package reflection

import (
	"log"
	"reflect"
	"strings"
)

const (
	// TagName the name of the tag used for reflection
	TagName = "swagger"
)

// MapMarshaller is an interface for things that have a need to customize this marshalling process
type MapMarshaller interface {
	MarshalMap() map[string]interface{}
}

var (
	mapMarshallerType = reflect.TypeOf(new(MapMarshaller)).Elem()
	stringType        = reflect.TypeOf("")
)

type parsedTag struct {
	Name       string
	OmitEmpty  bool
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
	var omitEmpty, byValue bool
	for _, p := range parts[1:] {
		if p == "omitempty" {
			omitEmpty = true
		}
		if p == "byValue" {
			byValue = true
		}
	}
	return &parsedTag{
		Name:       nm,
		OmitEmpty:  omitEmpty,
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

// MarshalMapRecursed this method marshals an interface to a map but skips the check for a custom interface
// on the provided data
func MarshalMapRecursed(data interface{}) map[string]interface{} {
	return marshalMap(data, true)
}

// MarshalMap this method marshals an interface to a map
// when the data provided implements MapMarshaller it will use that marshaller to get to the map
func MarshalMap(data interface{}) map[string]interface{} {
	return marshalMap(data, false)
}
func marshalMap(data interface{}, skipInterface bool) map[string]interface{} {
	if data == nil {
		return nil
	}

	val := reflect.ValueOf(data)
	tpe := val.Type()
	if tpe.Kind() == reflect.Ptr {
		val = val.Elem()
		tpe = tpe.Elem()
	}

	if !skipInterface && tpe.Implements(mapMarshallerType) {
		return val.Interface().(MapMarshaller).MarshalMap()
	}

	result := map[string]interface{}{}
	if tpe.Kind() == reflect.Map {
		keys := val.MapKeys()
		if len(keys) == 0 {
			return nil
		}
		for _, key := range keys {
			if key.Type() != stringType {
				log.Println("Only maps with string keys are allowed")
				return nil
			}
			value := val.MapIndex(key)
			var mapValue interface{}
			if isStruct(value) {
				mapValue = MarshalMap(value.Interface())
			} else {
				mapValue = value.Interface()
			}
			result[key.Interface().(string)] = mapValue
		}
		return result
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
			if fld.Kind() == reflect.Slice {
				var content []interface{}
				for j := 0; j < fld.Len(); j++ {
					el := fld.Index(j)
					if el.Kind() == reflect.Ptr {
						el = el.Elem()
					}
					if isStruct(el) && !tag.ByValue {
						content = append(content, MarshalMap(el.Interface()))
					} else {
						content = append(content, el.Interface())
					}
				}
				result[tag.Name] = content
				continue
			}

			if isStruct(fld) && !tag.ByValue {
				v := MarshalMap(fld.Interface())
				value = v
			} else {
				value = fld.Interface()
			}
			result[tag.Name] = value
		}
	}
	return result
}
