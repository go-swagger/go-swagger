package validate

import (
	"math"
	"reflect"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/reflection"
	"github.com/casualjim/go-swagger/spec"
)

// same as ECMA Number.MAX_SAFE_INTEGER and Number.MIN_SAFE_INTEGER
const (
	MaxJSONFloat = float64(1<<53 - 1)  // 9007199254740991.0 	 	 2^53 - 1
	MinJSONFloat = -float64(1<<53 - 1) //-9007199254740991.0	-2^53 - 1
)

// allow for integers [-2^53, 2^53-1] inclusive
func isFloat64AnInteger(f float64) bool {
	if math.IsNaN(f) || math.IsInf(f, 0) || f < MinJSONFloat || f > MaxJSONFloat {
		return false
	}

	return f == float64(int64(f)) || f == float64(uint64(f))
}

type typeValidator struct {
	Type   spec.StringOrArray
	Format string
	In     string
	Path   string
}

var jsonTypeNames = map[string]struct{}{
	"array":   struct{}{},
	"boolean": struct{}{},
	"integer": struct{}{},
	"null":    struct{}{},
	"number":  struct{}{},
	"object":  struct{}{},
	"string":  struct{}{},
}

func (t *typeValidator) schemaInfoForType(data interface{}) (string, string) {
	switch data.(type) {
	case []byte:
		return "string", "byte"
	case swagger.Date, *swagger.Date:
		return "string", "date"
	case swagger.DateTime, *swagger.DateTime:
		return "string", "date-time"
	case swagger.File, *swagger.File:
		return "file", ""
	default:
		val := reflect.ValueOf(data)
		tpe := val.Type()
		switch tpe.Kind() {
		case reflect.Bool:
			return "boolean", ""
		case reflect.String:
			return "string", ""
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			return "integer", "int32"
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			return "integer", "int64"
		case reflect.Float32:
			return "number", "float32"
		case reflect.Float64:
			return "number", "float64"
		case reflect.Slice:
			return "array", ""
		case reflect.Map, reflect.Struct:
			return "object", ""
		case reflect.Interface:
			// What to do here?
			panic("dunno what to do here")
		case reflect.Ptr:
			return t.schemaInfoForType(reflect.Indirect(val).Interface())
		}
	}
	return "", ""
}

func (t *typeValidator) SetPath(path string) {
	t.Path = path
}

func (t *typeValidator) Applies(source interface{}, kind reflect.Kind) bool {
	return len(t.Type) > 0 && reflect.TypeOf(source) == specSchemaType
}

func (t *typeValidator) Validate(data interface{}) *Result {
	if data == nil || reflection.IsZero(reflect.ValueOf(data)) {
		if len(t.Type) > 0 && !t.Type.Contains("null") {
			return sErr(errors.InvalidType(t.Path, t.In, strings.Join(t.Type, ","), "null"))
		}
		return &Result{}
	}

	// check if the type matches, should be used in every validator chain as first item
	val := reflect.Indirect(reflect.ValueOf(data))

	schType, format := t.schemaInfoForType(data)
	isLowerInt := t.Format == "int64" && format == "int32"
	isLowerFloat := t.Format == "float64" && format == "float32"
	allowedStringFormat := schType == "string" && format == ""

	if !allowedStringFormat && t.Format != "" && !(format == t.Format || isLowerInt || isLowerFloat) {
		return sErr(errors.InvalidType(t.Path, t.In, t.Format, format))
	}
	if t.Type.Contains(schType) && allowedStringFormat {
		return &Result{}
	}

	isFloatInt := schType == "number" && isFloat64AnInteger(val.Float()) && t.Type.Contains("integer")
	isIntFloat := schType == "integer" && t.Type.Contains("number")
	if !(t.Type.Contains(schType) || isFloatInt || isIntFloat) {
		return sErr(errors.InvalidType(t.Path, t.In, strings.Join(t.Type, ","), schType))
	}
	return &Result{}
}
