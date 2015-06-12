package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/swag"
)

var goImports = map[string]string{
	"inf.Dec":   "speter.net/go/exp/math/dec/inf",
	"big.Int":   "math/big",
	"swagger.*": "github.com/casualjim/go-swagger/httpkit",
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
	"strfmt.Password":   "strfmt.Password(\"\")",
	"strfmt.HexColor":   "strfmt.HexColor(\"#000000\")",
	"strfmt.RGBColor":   "strfmt.RGBColor(\"rgb(0,0,0)\")",
	"strfmt.Base64":     "nil",
	"strfmt.Duration":   "0",
}

var stringConverters = map[string]string{
	"int8":    "swag.ConvertInt8",
	"int16":   "swag.ConvertInt16",
	"int32":   "swag.ConvertInt32",
	"int64":   "swag.ConvertInt64",
	"uint8":   "swag.ConvertUint8",
	"uint16":  "swag.ConvertUint16",
	"uint32":  "swag.ConvertUint32",
	"uint64":  "swag.ConvertUint64",
	"bool":    "swag.ConvertBool",
	"float32": "swag.ConvertFloat32",
	"float64": "swag.ConvertFloat64",
}

var stringFormatters = map[string]string{
	"int8":    "swag.FormatInt8",
	"int16":   "swag.FormatInt16",
	"int32":   "swag.FormatInt32",
	"int64":   "swag.FormatInt64",
	"uint8":   "swag.FormatUint8",
	"uint16":  "swag.FormatUint16",
	"uint32":  "swag.FormatUint32",
	"uint64":  "swag.FormatUint64",
	"bool":    "swag.FormatBool",
	"float32": "swag.FormatFloat32",
	"float64": "swag.FormatFloat64",
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
	"password":   "strfmt.Password",
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
	"file":       "httpkit.File",
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
		// TODO: Items can't be nil per spec, this should return an error
		if items == nil {
			return "[]interface{}"
		}
		return "[]" + resolveSimpleType(items.Type, items.Format, items.Items)
	}
	return tn
}

type typeResolver struct {
	Doc           *spec.Document
	ModelsPackage string
}

func (t *typeResolver) ResolveSchema(schema *spec.Schema, isAnonymous bool) (result resolvedType, err error) {
	if schema == nil {
		result.IsInterface = true
		result.GoType = "interface{}"
		return
	}

	if schema.Ref.GetURL() != nil {
		ref, er := spec.ResolveRef(t.Doc.Spec(), &schema.Ref)
		if er != nil {
			err = er
			return
		}
		var tn string
		if gn, ok := ref.Extensions["x-go-name"]; ok {
			tn = gn.(string)
		} else {
			tn = swag.ToGoName(filepath.Base(schema.Ref.GetURL().Fragment))
		}

		if ref.Type.Contains("object") {
			result.GoType = t.ModelsPackage + "." + tn
			result.SwaggerType = "object"
			result.IsComplexObject = true
		} else {
			result, err = t.ResolveSchema(ref, true)
			if err != nil {
				return
			}
		}
		result.IsAnonymous = false

		return
	}
	if schema.Format != "" {
		schFmt := strings.Replace(schema.Format, "-", "", -1)
		if tpe, ok := typeMapping[schFmt]; ok {
			result.SwaggerType = "string"
			if len(schema.Type) > 0 {
				result.SwaggerType = schema.Type[0]
			}
			result.SwaggerFormat = schema.Format
			result.GoType = tpe
			return
		}
	}
	if schema.Type.Contains("array") {
		result.IsArray = true
		if schema.Items == nil {
			result.GoType = "[]interface{}"
			result.SwaggerType = "array"
			result.ElementType = &resolvedType{
				IsInterface: true,
				GoType:      "interface{}",
			}
			return
		}
		if len(schema.Items.Schemas) > 0 {
			result.IsArray = false
			result.IsTuple = true
			for _, esch := range schema.Items.Schemas {
				et, er := t.ResolveSchema(&esch, true)
				if er != nil {
					err = er
					return
				}
				result.TupleTypes = append(result.TupleTypes, &et)
			}
			return
		}
		rt, er := t.ResolveSchema(schema.Items.Schema, true)
		if er != nil {
			err = er
			return
		}
		result.GoType = "[]" + rt.GoType
		result.SwaggerType = "array"
		result.ElementType = &rt
		return
	}
	if schema.Type.Contains("file") {
		result.GoType = typeMapping["file"]
		result.SwaggerType = "file"
		return
	}
	if schema.Type.Contains("number") {
		result.GoType = typeMapping["number"]
		result.SwaggerType = "number"
		return
	}
	if schema.Type.Contains("integer") {
		result.GoType = typeMapping["integer"]
		result.SwaggerType = "integer"
		return
	}
	if schema.Type.Contains("boolean") {
		result.GoType = typeMapping["boolean"]
		result.SwaggerType = "boolean"
		return
	}
	if schema.Type.Contains("string") {
		result.GoType = "string"
		result.SwaggerType = "string"
		return
	}
	if schema.Type.Contains("object") || schema.Type.Contains("") || len(schema.Type) == 0 {
		// if this schema has properties, build a map of property name to
		// resolved type, this should also flag the object as anonymous,
		// when a ref is found, the anonymous flag will be reset
		if isAnonymous && len(schema.Properties) > 0 {
			result.PropertyTypes = make(map[string]*resolvedType)
			for key, prop := range schema.Properties {
				result.IsAnonymous = isAnonymous
				result.IsComplexObject = true
				prt, er := t.ResolveSchema(&prop, true)
				if er != nil {
					err = er
					return
				}
				result.PropertyTypes[key] = &prt
			}
			// no return here, still need to check for additional properties
		}

		// account for additional properties
		if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
			et, er := t.ResolveSchema(schema.AdditionalProperties.Schema, true)
			if er != nil {
				err = er
				return
			}
			result.GoType = "map[string]" + et.GoType
			result.ElementType = &et
			result.IsMap = true
			result.SwaggerType = "object"
			return
		}

		if len(schema.Properties) > 0 {
			return
		}
		result.GoType = "map[string]interface{}"
		result.ElementType = &resolvedType{
			IsInterface: true,
			GoType:      "interface{}",
		}

		result.IsMap = true
		result.SwaggerType = "object"
		return
	}
	err = fmt.Errorf("unresolvable: %v (format %q)", schema.Type, schema.Format)
	return
}

type resolvedType struct {
	IsAnonymous bool
	IsArray     bool
	IsMap       bool
	IsInterface bool
	// A tuple gets rendered as an anonymous struct with P[index] as property name
	IsTuple         bool
	IsComplexObject bool

	GoType        string
	SwaggerType   string
	SwaggerFormat string
	ElementType   *resolvedType
	TupleTypes    []*resolvedType
	PropertyTypes map[string]*resolvedType
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
	"strfmt.Password":   struct{}{},
	"strfmt.HexColor":   struct{}{},
	"strfmt.RGBColor":   struct{}{},
	"strfmt.Base64":     struct{}{},
	// "strfmt.Duration":   struct{}{},
}
