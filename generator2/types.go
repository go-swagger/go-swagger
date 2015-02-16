package generator2

import (
	"path/filepath"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

func typeForSchemaOrArray(schemas *spec.SchemaOrArray) string {
	if schemas == nil || len(schemas.Schemas) > 0 {
		return "interface{}"
	}
	return typeForSchema(schemas.Schema)
}

var reservedGoWords = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

var defaultGoImports = []string{
	"bool", "int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"float32", "float64", "interface{}", "string",
	"byte", "rune",
}

var goTypeMapping = map[string]string{
	"array":    "[]",
	"map":      "map",
	"List":     "[]",
	"boolean":  "bool",
	"int":      "int32",
	"float":    "float32",
	"number":   "inf.Dec",
	"DateTime": "swagger.DateTime",
	"long":     "int64",
	"short":    "int16",
	"char":     "rune",
	"double":   "float64",
	"object":   "interface{}",
	"integer":  "int32",
}

var goImports = map[string]string{
	"inf.Dec":   "speter.net/go/exp/math/dec/inf",
	"big.Int":   "math/big",
	"swagger.*": "github.com/casualjim/go-swagger",
}

var stringConverters = map[string]string{
	"int8":    "util.ConvertInt8",
	"int16":   "util.ConvertInt16",
	"int32":   "util.ConvertInt32",
	"int64":   "util.ConvertInt64",
	"uint8":   "util.ConvertUint8",
	"uint16":  "util.ConvertUint16",
	"uint32":  "util.ConvertUint32",
	"uint64":  "util.ConvertUint64",
	"bool":    "util.ConvertBool",
	"float32": "util.ConvertFloat32",
	"float64": "util.ConvertFloat64",
}

var typeMapping = map[string]string{
	"byte":       "[]byte{}",
	"date":       "swagger.Date",
	"datetime":   "swagger.DateTime",
	"uri":        "swagger.URI",
	"email":      "swagger.Email",
	"hostname":   "swagger.Hostname",
	"ipv4":       "swagger.IPv4",
	"ipv6":       "swagger.IPv6",
	"uuid":       "swagger.UUID",
	"uuid3":      "swagger.UUID3",
	"uuid4":      "swagger.UUID4",
	"uuid5":      "swagger.UUID5",
	"isbn":       "swagger.ISBN",
	"isbn10":     "swagger.ISBN10",
	"isbn13":     "swagger.ISBN13",
	"creditcard": "swagger.CreditCard",
	"ssn":        "swagger.SSN",
	"hexcolor":   "swagger.HexColor",
	"rgbcolor":   "swagger.RGBColor",
	"char":       "rune",
	"int":        "int64",
	"int8":       "int8",
	"int16":      "int16",
	"int32":      "int32",
	"int64":      "int64",
	"uint":       "uint64",
	"uint8":      "uint8",
	"uint16":     "uint16",
	"uint32":     "uint32",
	"uint64":     "uint64",
	"float":      "float32",
	"double":     "float64",
	"number":     "float64",
	"integer":    "int64",
	"boolean":    "bool",
	"file":       "swagger.File",
}

var swaggerTypeMapping map[string]string

func init() {
	swaggerTypeMapping = make(map[string]string)
	for k, v := range typeMapping {
		swaggerTypeMapping[v] = k
	}
}

func typeForParameter(param spec.Parameter) string {
	return resolveSimpleType(param.Type, param.Format, param.Items)
}

func resolveSimpleType(tn, fmt string, items *spec.Items) string {
	if fmt != "" {
		if tpe, ok := typeMapping[strings.Replace(fmt, "-", "", -1)]; ok {
			return tpe
		}
	}
	if tpe, ok := typeMapping[tn]; ok {
		return tpe
	}

	if tn == "array" {
		if items == nil {
			return "[]interface{}"
		}
		return "[]" + resolveSimpleType(items.Type, items.Format, items.Items)
	}
	return tn
}

func typeForSchema(schema *spec.Schema) string {
	if schema == nil {
		return "interface{}"
	}
	if schema.Ref.GetURL() != nil {
		return util.ToGoName(filepath.Base(schema.Ref.GetURL().Fragment))
	}
	if schema.Format != "" {
		if tpe, ok := typeMapping[strings.Replace(schema.Format, "-", "", -1)]; ok {
			return tpe
		}
	}
	if schema.Type.Contains("array") {
		return "[]" + typeForSchemaOrArray(schema.Items)
	}
	if schema.Type.Contains("file") {
		return typeMapping["file"]
	}
	if schema.Type.Contains("number") {
		return typeMapping["number"]
	}
	if schema.Type.Contains("integer") {
		return typeMapping["integer"]
	}
	if schema.Type.Contains("boolean") {
		return typeMapping["boolean"]
	}
	if schema.Type.Contains("string") {
		return "string"
	}
	if schema.Type.Contains("object") || schema.Type.Contains("") || len(schema.Type) == 0 {
		return "map[string]interface{}"
	}
	return "interface{}"
}

var primitives = map[string]struct{}{
	"bool":       struct{}{},
	"uint":       struct{}{},
	"uint8":      struct{}{},
	"uint16":     struct{}{},
	"uint32":     struct{}{},
	"uint64":     struct{}{},
	"int":        struct{}{},
	"int8":       struct{}{},
	"int16":      struct{}{},
	"int32":      struct{}{},
	"int64":      struct{}{},
	"float32":    struct{}{},
	"float64":    struct{}{},
	"string":     struct{}{},
	"complex64":  struct{}{},
	"complex128": struct{}{},
	"byte":       struct{}{},
	"[]byte":     struct{}{},
	"rune":       struct{}{},
}

var customFormatters = map[string]struct{}{
	"swagger.DateTime":   struct{}{},
	"swagger.Time":       struct{}{},
	"swagger.URI":        struct{}{},
	"swagger.Email":      struct{}{},
	"swagger.Hostname":   struct{}{},
	"swagger.IPv4":       struct{}{},
	"swagger.UUID":       struct{}{},
	"swagger.UUID3":      struct{}{},
	"swagger.UUID4":      struct{}{},
	"swagger.UUID5":      struct{}{},
	"swagger.ISBN":       struct{}{},
	"swagger.ISBN10":     struct{}{},
	"swagger.ISBN13":     struct{}{},
	"swagger.CreditCard": struct{}{},
	"swagger.SSN":        struct{}{},
	"swagger.HexColor":   struct{}{},
	"swagger.RGBColor":   struct{}{},
}
