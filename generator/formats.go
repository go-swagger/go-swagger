// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

// TODO: we may probably find a way to register most of this dynamically from strfmt

// map of function calls to be generated to get the zero value of a given type
var zeroes = map[string]string{
	"bool":    "false",
	"float32": "0",
	"float64": "0",
	"int":     "0",
	"int8":    "0",
	"int16":   "0",
	"int32":   "0",
	"int64":   "0",
	"string":  "\"\"",
	"uint":    "0",
	"uint8":   "0",
	"uint16":  "0",
	"uint32":  "0",
	"uint64":  "0",
	// Extended formats (23 formats corresponding to the Default registry
	// provided by go-openapi/strfmt)
	"strfmt.Base64":     "strfmt.Base64([]byte(nil))",
	"strfmt.CreditCard": "strfmt.CreditCard(\"\")",
	"strfmt.Date":       "strfmt.Date{}",
	"strfmt.DateTime":   "strfmt.DateTime{}",
	"strfmt.Duration":   "strfmt.Duration(0)",
	"strfmt.Email":      "strfmt.Email(\"\")",
	"strfmt.HexColor":   "strfmt.HexColor(\"#000000\")",
	"strfmt.Hostname":   "strfmt.Hostname(\"\")",
	"strfmt.IPv4":       "strfmt.IPv4(\"\")",
	"strfmt.IPv6":       "strfmt.IPv6(\"\")",
	"strfmt.ISBN":       "strfmt.ISBN(\"\")",
	"strfmt.ISBN10":     "strfmt.ISBN10(\"\")",
	"strfmt.ISBN13":     "strfmt.ISBN13(\"\")",
	"strfmt.MAC":        "strfmt.MAC(\"\")",
	"strfmt.ObjectId":   "strfmt.ObjectId(\"\")",
	"strfmt.Password":   "strfmt.Password(\"\")",
	"strfmt.RGBColor":   "strfmt.RGBColor(\"rgb(0,0,0)\")",
	"strfmt.SSN":        "strfmt.SSN(\"\")",
	"strfmt.URI":        "strfmt.URI(\"\")",
	"strfmt.UUID":       "strfmt.UUID(\"\")",
	"strfmt.UUID3":      "strfmt.UUID3(\"\")",
	"strfmt.UUID4":      "strfmt.UUID4(\"\")",
	"strfmt.UUID5":      "strfmt.UUID5(\"\")",
	//"file":       "runtime.File",
}

// conversion functions from string representation to a numerical or boolean
// primitive type
var stringConverters = map[string]string{
	"bool":    "swag.ConvertBool",
	"float32": "swag.ConvertFloat32",
	"float64": "swag.ConvertFloat64",
	"int8":    "swag.ConvertInt8",
	"int16":   "swag.ConvertInt16",
	"int32":   "swag.ConvertInt32",
	"int64":   "swag.ConvertInt64",
	"uint8":   "swag.ConvertUint8",
	"uint16":  "swag.ConvertUint16",
	"uint32":  "swag.ConvertUint32",
	"uint64":  "swag.ConvertUint64",
}

// formatting (string representation) functions from a native representation
// of a numerical or boolean primitive type
var stringFormatters = map[string]string{
	"bool":    "swag.FormatBool",
	"float32": "swag.FormatFloat32",
	"float64": "swag.FormatFloat64",
	"int8":    "swag.FormatInt8",
	"int16":   "swag.FormatInt16",
	"int32":   "swag.FormatInt32",
	"int64":   "swag.FormatInt64",
	"uint8":   "swag.FormatUint8",
	"uint16":  "swag.FormatUint16",
	"uint32":  "swag.FormatUint32",
	"uint64":  "swag.FormatUint64",
}

// typeMapping contains a mapping of format (or type name) to go type
var typeMapping = map[string]string{
	// Standard formats with native, straightforward, mapping
	"binary":  "io.ReadCloser",
	"boolean": "bool",
	"char":    "rune",
	"double":  "float64",
	"float":   "float32",
	"int":     "int64",
	"int8":    "int8",
	"int16":   "int16",
	"int32":   "int32",
	"int64":   "int64",
	"integer": "int64",
	"number":  "float64",
	"uint":    "uint64",
	"uint8":   "uint8",
	"uint16":  "uint16",
	"uint32":  "uint32",
	"uint64":  "uint64",
	// Extended format registry from go-openapi/strfmt.
	// Currently, 23 such formats are supported (default strftm registry),
	// plus the following aliases:
	//  - "datetime" alias for the more official "date-time"
	//  - "objectid" and "ObjectId" aliases for "bsonobjectid"
	"byte":         "strfmt.Base64",
	"creditcard":   "strfmt.CreditCard",
	"date":         "strfmt.Date",
	"date-time":    "strfmt.DateTime",
	"datetime":     "strfmt.DateTime",
	"duration":     "strfmt.Duration",
	"email":        "strfmt.Email",
	"hexcolor":     "strfmt.HexColor",
	"hostname":     "strfmt.Hostname",
	"ipv4":         "strfmt.IPv4",
	"ipv6":         "strfmt.IPv6",
	"isbn":         "strfmt.ISBN",
	"isbn10":       "strfmt.ISBN10",
	"isbn13":       "strfmt.ISBN13",
	"mac":          "strfmt.MAC",
	"bsonobjectid": "strfmt.ObjectId",
	"objectid":     "strfmt.ObjectId",
	"ObjectId":     "strfmt.ObjectId", // NOTE: does it work with uppercase?
	"password":     "strfmt.Password",
	"rgbcolor":     "strfmt.RGBColor",
	"ssn":          "strfmt.SSN",
	"uri":          "strfmt.URI",
	"uuid":         "strfmt.UUID",
	"uuid3":        "strfmt.UUID3",
	"uuid4":        "strfmt.UUID4",
	"uuid5":        "strfmt.UUID5",
	// For file producers
	"file": "runtime.File",
}

// go primitive types
var primitives = map[string]struct{}{
	"bool":       struct{}{},
	"byte":       struct{}{},
	"[]byte":     struct{}{},
	"complex64":  struct{}{},
	"complex128": struct{}{},
	"float32":    struct{}{},
	"float64":    struct{}{},
	"int":        struct{}{},
	"int8":       struct{}{},
	"int16":      struct{}{},
	"int32":      struct{}{},
	"int64":      struct{}{},
	"rune":       struct{}{},
	"string":     struct{}{},
	"uint":       struct{}{},
	"uint8":      struct{}{},
	"uint16":     struct{}{},
	"uint32":     struct{}{},
	"uint64":     struct{}{},
}

// Formats with a custom formatter.
// Currently, 23 such formats are supported
var customFormatters = map[string]struct{}{
	"strfmt.Base64":     struct{}{},
	"strfmt.CreditCard": struct{}{},
	"strfmt.Date":       struct{}{},
	"strfmt.DateTime":   struct{}{},
	"strfmt.Duration":   struct{}{},
	"strfmt.Email":      struct{}{},
	"strfmt.HexColor":   struct{}{},
	"strfmt.Hostname":   struct{}{},
	"strfmt.IPv4":       struct{}{},
	"strfmt.IPv6":       struct{}{},
	"strfmt.ISBN":       struct{}{},
	"strfmt.ISBN10":     struct{}{},
	"strfmt.ISBN13":     struct{}{},
	"strfmt.MAC":        struct{}{},
	"strfmt.ObjectId":   struct{}{},
	"strfmt.Password":   struct{}{},
	"strfmt.RGBColor":   struct{}{},
	"strfmt.SSN":        struct{}{},
	"strfmt.URI":        struct{}{},
	"strfmt.UUID":       struct{}{},
	"strfmt.UUID3":      struct{}{},
	"strfmt.UUID4":      struct{}{},
	"strfmt.UUID5":      struct{}{},
	// the following interfaces do not generate validations
	"io.ReadCloser": struct{}{}, // for "format": "binary" (server side)
	"io.Writer":     struct{}{}, // for "format": "binary" (client side)
	// NOTE: runtime.File is not a customFormatter
}
