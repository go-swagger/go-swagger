package validate

import (
	"math"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
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

var formatCheckers = map[string]func(string) bool{
	"datetime":   IsDateTime,
	"date":       IsDate,
	"byte":       govalidator.IsBase64,
	"uri":        IsURI,
	"email":      govalidator.IsEmail,
	"hostname":   IsHostname,
	"ipv4":       govalidator.IsIPv4,
	"ipv6":       govalidator.IsIPv6,
	"uuid":       govalidator.IsUUID,
	"uuid3":      govalidator.IsUUIDv3,
	"uuid4":      govalidator.IsUUIDv4,
	"uuid5":      govalidator.IsUUIDv5,
	"isbn":       func(str string) bool { return govalidator.IsISBN10(str) || govalidator.IsISBN13(str) },
	"isbn10":     govalidator.IsISBN10,
	"isbn13":     govalidator.IsISBN13,
	"creditcard": govalidator.IsCreditCard,
	"ssn":        govalidator.IsSSN,
	"hexcolor":   govalidator.IsHexcolor,
	"rgbcolor":   govalidator.IsRGBcolor,
}

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
		return "string", "datetime"
	case swagger.File, *swagger.File:
		return "file", ""
	case swagger.URI, *swagger.URI:
		return "string", "uri"
	case swagger.Email, *swagger.Email:
		return "string", "email"
	case swagger.Hostname, *swagger.Hostname:
		return "string", "hostname"
	case swagger.IPv4, *swagger.IPv4:
		return "string", "ipv4"
	case swagger.IPv6, *swagger.IPv6:
		return "string", "ipv6"
	case swagger.UUID, *swagger.UUID:
		return "string", "uuid"
	case swagger.UUID3, *swagger.UUID3:
		return "string", "uuid3"
	case swagger.UUID4, *swagger.UUID4:
		return "string", "uuid4"
	case swagger.UUID5, *swagger.UUID5:
		return "string", "uuid5"
	case swagger.ISBN, *swagger.ISBN:
		return "string", "isbn"
	case swagger.ISBN10, *swagger.ISBN10:
		return "string", "isbn10"
	case swagger.ISBN13, *swagger.ISBN13:
		return "string", "isbn13"
	case swagger.CreditCard, *swagger.CreditCard:
		return "string", "creditcard"
	case swagger.SSN, *swagger.SSN:
		return "string", "ssn"
	case swagger.HexColor, *swagger.HexColor:
		return "string", "hexcolor"
	case swagger.RGBColor, *swagger.RGBColor:
		return "string", "rgbcolor"
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
	return (len(t.Type) > 0 || t.Format != "") && reflect.TypeOf(source) == specSchemaType
}

func (t *typeValidator) Validate(data interface{}) *Result {
	result := new(Result)
	result.Inc()
	if data == nil || reflection.IsZero(reflect.ValueOf(data)) {
		if len(t.Type) > 0 && !t.Type.Contains("null") {
			return sErr(errors.InvalidType(t.Path, t.In, strings.Join(t.Type, ","), "null"))
		}
		return result
	}

	// check if the type matches, should be used in every validator chain as first item
	val := reflect.Indirect(reflect.ValueOf(data))

	schType, format := t.schemaInfoForType(data)
	isLowerInt := t.Format == "int64" && format == "int32"
	isLowerFloat := t.Format == "float64" && format == "float32"

	if val.Kind() != reflect.String && t.Format != "" && !(format == t.Format || isLowerInt || isLowerFloat) {
		return sErr(errors.InvalidType(t.Path, t.In, t.Format, format))
	}
	if t.Format != "" && val.Kind() == reflect.String {
		return result
	}

	isFloatInt := schType == "number" && isFloat64AnInteger(val.Float()) && t.Type.Contains("integer")
	isIntFloat := schType == "integer" && t.Type.Contains("number")
	if !(t.Type.Contains(schType) || isFloatInt || isIntFloat) {
		return sErr(errors.InvalidType(t.Path, t.In, strings.Join(t.Type, ","), schType))
	}
	return result
}
