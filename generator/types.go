package generator

import (
	"path/filepath"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

func typeForSchemaOrArray(schemas *spec.SchemaOrArray, modelsPkg string) string {
	if schemas == nil || len(schemas.Schemas) > 0 {
		return "interface{}"
	}
	return typeForSchema(schemas.Schema, modelsPkg)
}

var goImports = map[string]string{
	"inf.Dec":   "speter.net/go/exp/math/dec/inf",
	"big.Int":   "math/big",
	"swagger.*": "github.com/casualjim/go-swagger",
}

var zeroes = map[string]string{
	"string":            "\"\"",
	"int8":              "0",
	"int16":             "0",
	"int32":             "0",
	"int64":             "0",
	"uint8":             "0",
	"uint16":            "0",
	"uint32":            "0",
	"uint64":            "0",
	"bool":              "false",
	"float32":           "0",
	"float64":           "0",
	"strfmt.DateTime":   "strfmt.DateTime{}",
	"strfmt.Date":       "strfmt.Date{}",
	"strfmt.URI":        "strfmt.URI(\"\")",
	"strfmt.Email":      "strfmt.Email(\"\")",
	"strfmt.Hostname":   "strfmt.Hostname(\"\")",
	"strfmt.IPv4":       "strfmt.IPv4(\"\")",
	"strfmt.IPv6":       "strfmt.IPv6(\"\")",
	"strfmt.UUID":       "strfmt.UUID(\"\")",
	"strfmt.UUID3":      "strfmt.UUID3(\"\")",
	"strfmt.UUID4":      "strfmt.UUID4(\"\")",
	"strfmt.UUID5":      "strfmt.UUID5(\"\")",
	"strfmt.ISBN":       "strfmt.ISBN(\"\")",
	"strfmt.ISBN10":     "strfmt.ISBN10(\"\")",
	"strfmt.ISBN13":     "strfmt.ISBN13(\"\")",
	"strfmt.CreditCard": "strfmt.CreditCard(\"\")",
	"strfmt.SSN":        "strfmt.SSN(\"\")",
	"strfmt.HexColor":   "strfmt.HexColor(\"#000000\")",
	"strfmt.RGBColor":   "strfmt.RGBColor(\"rgb(0,0,0)\")",
	"strfmt.Base64":     "nil",
	"strfmt.Duration":   "0",
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

// typeMapping contais a mapping of format or type name to go type
var typeMapping = map[string]string{
	"byte":       "strfmt.Base64",
	"date":       "strfmt.Date",
	"datetime":   "strfmt.DateTime",
	"uri":        "strfmt.URI",
	"email":      "strfmt.Email",
	"hostname":   "strfmt.Hostname",
	"ipv4":       "strfmt.IPv4",
	"ipv6":       "strfmt.IPv6",
	"uuid":       "strfmt.UUID",
	"uuid3":      "strfmt.UUID3",
	"uuid4":      "strfmt.UUID4",
	"uuid5":      "strfmt.UUID5",
	"isbn":       "strfmt.ISBN",
	"isbn10":     "strfmt.ISBN10",
	"isbn13":     "strfmt.ISBN13",
	"creditcard": "strfmt.CreditCard",
	"ssn":        "strfmt.SSN",
	"hexcolor":   "strfmt.HexColor",
	"rgbcolor":   "strfmt.RGBColor",
	"duration":   "strfmt.Duration",
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

// swaggerTypeMapping contains a mapping from go type to swagger type or format
var swaggerTypeName map[string]string

func init() {
	swaggerTypeName = make(map[string]string)
	for k, v := range typeMapping {
		swaggerTypeName[v] = k
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

func typeForSchema(schema *spec.Schema, modelsPkg string) string {
	if schema == nil {
		return "interface{}"
	}
	if schema.Ref.GetURL() != nil {
		tn := util.ToGoName(filepath.Base(schema.Ref.GetURL().Fragment))
		if modelsPkg != "" {
			return modelsPkg + "." + tn
		}
		return tn
	}
	if schema.Format != "" {
		if tpe, ok := typeMapping[strings.Replace(schema.Format, "-", "", -1)]; ok {
			return tpe
		}
	}
	if schema.Type.Contains("array") {
		return "[]" + typeForSchemaOrArray(schema.Items, modelsPkg)
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
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		return "map[string]" + typeForSchema(schema.AdditionalProperties.Schema, modelsPkg)
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
	// "strfmt.DateTime":   struct{}{},
	// "strfmt.Date":       struct{}{},
	"strfmt.URI":        struct{}{},
	"strfmt.Email":      struct{}{},
	"strfmt.Hostname":   struct{}{},
	"strfmt.IPv4":       struct{}{},
	"strfmt.IPv6":       struct{}{},
	"strfmt.UUID":       struct{}{},
	"strfmt.UUID3":      struct{}{},
	"strfmt.UUID4":      struct{}{},
	"strfmt.UUID5":      struct{}{},
	"strfmt.ISBN":       struct{}{},
	"strfmt.ISBN10":     struct{}{},
	"strfmt.ISBN13":     struct{}{},
	"strfmt.CreditCard": struct{}{},
	"strfmt.SSN":        struct{}{},
	"strfmt.HexColor":   struct{}{},
	"strfmt.RGBColor":   struct{}{},
	"strfmt.Base64":     struct{}{},
	// "strfmt.Duration":   struct{}{},
}
