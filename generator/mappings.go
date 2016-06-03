package generator

var zeroes map[string]string
var stringConverters map[string]string
var stringFormatters map[string]string
var typeMapping map[string]string
var swaggerTypeName map[string]string
var primitives map[string]struct{}
var customFormatters map[string]struct{}
var inEasyJSONMap map[string]string
var outEasyJSONMap map[string]string

func init() {

	zeroes = map[string]string{
		"string":            "\"\"",
		"int8":              "0",
		"int":               "0",
		"int16":             "0",
		"int32":             "0",
		"int64":             "0",
		"uint":              "0",
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

	stringConverters = map[string]string{
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

	stringFormatters = map[string]string{
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

	typeMapping = map[string]string{
		"byte":       "strfmt.Base64",
		"date":       "strfmt.Date",
		"datetime":   "strfmt.DateTime",
		"date-time":  "strfmt.DateTime",
		"uri":        "strfmt.URI",
		"email":      "strfmt.Email",
		"hostname":   "strfmt.Hostname",
		"ipv4":       "strfmt.IPv4",
		"ipv6":       "strfmt.IPv6",
		"mac":        "strfmt.MAC",
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
		"binary":     "io.ReadCloser",
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
		"file":       "runtime.File",
	}

	primitives = map[string]struct{}{
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

	customFormatters = map[string]struct{}{
		"strfmt.DateTime":   struct{}{},
		"strfmt.Date":       struct{}{},
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
		"strfmt.Duration":   struct{}{},
		"io.ReadCloser":     struct{}{},
		"io.Writer":         struct{}{},
	}

	inEasyJSONMap = map[string]string{
		"uint8":   "in.Uint8()",
		"uint16":  "in.Uint16()",
		"uint32":  "in.Uint32()",
		"uint":    "in.Uint()",
		"uint64":  "in.Uint64()",
		"int8":    "in.Int8()",
		"int16":   "in.Int16()",
		"int32":   "in.Int32()",
		"int":     "in.Int()",
		"int64":   "in.Int64()",
		"float32": "in.Float32()",
		"float64": "in.Float64()",
		"string":  "in.String()",
		"bool":    "in.Bool()",
	}

	outEasyJSONMap = map[string]string{
		"uint8":   "out.Uint8",
		"uint16":  "out.Uint16",
		"uint32":  "out.Uint32",
		"uint":    "out.Uint",
		"uint64":  "out.Uint64",
		"int8":    "out.Int8",
		"int16":   "out.Int16",
		"int32":   "out.Int32",
		"int":     "out.Int",
		"int64":   "out.Int64",
		"float32": "out.Float32",
		"float64": "out.Float64",
		"string":  "out.String",
		"bool":    "out.Bool",
	}

	swaggerTypeName = make(map[string]string)
	for k, v := range typeMapping {
		swaggerTypeName[v] = k
	}
}
