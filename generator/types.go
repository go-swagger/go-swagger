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

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
)

// var goImports = map[string]string{
// 	"inf.Dec":   "speter.net/go/exp/math/dec/inf",
// 	"big.Int":   "math/big",
// 	"swagger.*": "github.com/go-openapi/runtime",
// }

const (
	iface       = "interface{}"
	array       = "array"
	file        = "file"
	number      = "number"
	integer     = "integer"
	boolean     = "boolean"
	str         = "string"
	object      = "object"
	binary      = "binary"
	xNullable   = "x-nullable"
	xIsNullable = "x-isnullable"
	sHTTP       = "http"
)

var zeroes = map[string]string{
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

// typeMapping contains a mapping of format or type name to go type
var typeMapping = map[string]string{
	"byte":       "strfmt.Base64",
	"date":       "strfmt.Date",
	"datetime":   "strfmt.DateTime",
	"date-time":  "strfmt.DateTime",
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

// swaggerTypeMapping contains a mapping from go type to swagger type or format
var swaggerTypeName map[string]string

func init() {
	swaggerTypeName = make(map[string]string)
	for k, v := range typeMapping {
		swaggerTypeName[v] = k
	}
}

func simpleResolvedType(tn, fmt string, items *spec.Items) (result resolvedType) {
	result.SwaggerType = tn
	result.SwaggerFormat = fmt
	//_, result.IsPrimitive = primitives[tn]

	if fmt != "" {
		fmtn := strings.Replace(fmt, "-", "", -1)
		if tpe, ok := typeMapping[fmtn]; ok {
			result.GoType = tpe
			result.IsPrimitive = true
			_, result.IsCustomFormatter = customFormatters[tpe]
			result.IsStream = fmt == binary
			return
		}
	}

	if tpe, ok := typeMapping[tn]; ok {
		result.GoType = tpe
		_, result.IsPrimitive = primitives[tpe]
		result.IsPrimitive = ok
		return
	}

	if tn == array {
		result.IsArray = true
		result.IsPrimitive = false
		result.IsCustomFormatter = false
		result.IsNullable = false
		if items == nil {
			result.GoType = "[]" + iface
			return
		}
		res := simpleResolvedType(items.Type, items.Format, items.Items)
		result.GoType = "[]" + res.GoType
		return
	}
	result.GoType = tn
	_, result.IsPrimitive = primitives[tn]
	return
}

func typeForHeader(header spec.Header) resolvedType {
	return simpleResolvedType(header.Type, header.Format, header.Items)
}

//
// func typeForParameter(param spec.Parameter) string {
// 	return resolveSimpleType(param.Type, param.Format, param.Items)
// }

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

func newTypeResolver(pkg string, doc *loads.Document) *typeResolver {
	resolver := typeResolver{ModelsPackage: pkg, Doc: doc}
	resolver.KnownDefs = make(map[string]struct{})
	for k, sch := range doc.Pristine().Spec().Definitions {
		resolver.KnownDefs[k] = struct{}{}
		if nm, ok := sch.Extensions["x-go-name"]; ok {
			resolver.KnownDefs[nm.(string)] = struct{}{}
		}
	}
	return &resolver
}

type typeResolver struct {
	Doc           *loads.Document
	ModelsPackage string
	ModelName     string
	KnownDefs     map[string]struct{}
}

func (t *typeResolver) IsNullable(schema *spec.Schema) bool {
	nullable := t.isNullable(schema)
	return nullable || len(schema.AllOf) > 0
}

func (t *typeResolver) resolveSchemaRef(schema *spec.Schema, isRequired bool) (returns bool, result resolvedType, err error) {
	if schema.Ref.String() != "" {
		if Debug {
			_, file, pos, _ := runtime.Caller(1)
			log.Printf("%s:%d: resolving ref (anon: %t, req: %t) %s\n", filepath.Base(file), pos, false, isRequired, schema.Ref.String())
		}
		returns = true

		ref, er := spec.ResolveRef(t.Doc.Spec(), &schema.Ref)
		if er != nil {
			err = er
			return
		}
		var nm = filepath.Base(schema.Ref.GetURL().Fragment)
		var tn string
		if gn, ok := ref.Extensions["x-go-name"]; ok {
			tn = gn.(string)
			nm = tn
		} /*else {
			tn = swag.ToGoName(nm)
		}*/

		res, er := t.ResolveSchema(ref, false, isRequired)
		if er != nil {
			err = er
			return
		}
		result = res

		result.GoType = t.goTypeName(nm)
		result.HasDiscriminator = ref.Discriminator != ""
		result.IsNullable = t.IsNullable(ref)
		//result.IsAliased = true
		return

	}
	return
}

func (t *typeResolver) inferAliasing(result *resolvedType, schema *spec.Schema, isAnonymous bool, isRequired bool) {
	if !isAnonymous && t.ModelName != "" {
		result.AliasedType = result.GoType
		result.IsAliased = true
		result.GoType = t.goTypeName(t.ModelName)
	}
}

func (t *typeResolver) resolveFormat(schema *spec.Schema, isAnonymous bool, isRequired bool) (returns bool, result resolvedType, err error) {

	if schema.Format != "" {
		if Debug {
			_, file, pos, _ := runtime.Caller(1)
			log.Printf("%s:%d: resolving format (anon: %t, req: %t)\n", filepath.Base(file), pos, isAnonymous, isRequired) //, bbb)
		}
		schFmt := strings.Replace(schema.Format, "-", "", -1)
		if tpe, ok := typeMapping[schFmt]; ok {
			returns = true
			result.SwaggerType = str
			if len(schema.Type) > 0 {
				result.SwaggerType = schema.Type[0]
			}
			result.SwaggerFormat = schema.Format
			result.GoType = tpe
			t.inferAliasing(&result, schema, isAnonymous, isRequired)
			result.IsPrimitive = schFmt != binary
			result.IsStream = schFmt == binary
			_, result.IsCustomFormatter = customFormatters[tpe]

			switch result.SwaggerType {
			case str:
				result.IsNullable = nullableStrfmt(schema, isRequired)
			case number, integer:
				result.IsNullable = nullableNumber(schema, isRequired)
			default:
				result.IsNullable = t.IsNullable(schema)
			}
			return
		}
	}
	return
}

func (t *typeResolver) isNullable(schema *spec.Schema) bool {
	return t.checkIsNullable(xIsNullable, schema) || t.checkIsNullable(xNullable, schema)
}

func (t *typeResolver) checkIsNullable(extension string, schema *spec.Schema) bool {
	v, found := schema.Extensions[extension]
	nullable, cast := v.(bool)
	return (found && cast && nullable) || len(schema.Properties) > 0
}

func (t *typeResolver) firstType(schema *spec.Schema) string {
	if len(schema.Type) == 0 || schema.Type[0] == "" {
		return object
	}
	return schema.Type[0]
}

func (t *typeResolver) resolveArray(schema *spec.Schema, isAnonymous, isRequired bool) (result resolvedType, err error) {
	if Debug {
		_, file, pos, _ := runtime.Caller(1)
		log.Printf("%s:%d: resolving array (anon: %t, req: %t)\n", filepath.Base(file), pos, isAnonymous, isRequired) //, bbb)
	}

	result.IsArray = true
	result.IsNullable = false
	if schema.AdditionalItems != nil {
		result.HasAdditionalItems = (schema.AdditionalItems.Allows || schema.AdditionalItems.Schema != nil)
	}

	if schema.Items == nil {
		result.GoType = "[]" + iface
		result.SwaggerType = array
		result.SwaggerFormat = ""
		t.inferAliasing(&result, schema, isAnonymous, isRequired)

		return
	}

	if len(schema.Items.Schemas) > 0 {
		result.IsArray = false
		result.IsTuple = true
		result.SwaggerType = array
		result.SwaggerFormat = ""
		t.inferAliasing(&result, schema, isAnonymous, isRequired)

		return
	}

	rt, er := t.ResolveSchema(schema.Items.Schema, true, false)
	if er != nil {
		err = er
		return
	}

	rt.IsNullable = t.IsNullable(schema.Items.Schema) && !rt.HasDiscriminator
	result.GoType = "[]" + rt.GoType
	if rt.IsNullable && !strings.HasPrefix(rt.GoType, "*") {
		result.GoType = "[]*" + rt.GoType
	}

	result.ElemType = &rt
	result.SwaggerType = array
	result.SwaggerFormat = ""
	t.inferAliasing(&result, schema, isAnonymous, isRequired)

	return
}

func (t *typeResolver) goTypeName(nm string) string {
	if t.ModelsPackage == "" {
		return swag.ToGoName(nm)
	}
	if _, ok := t.KnownDefs[nm]; ok {
		return strings.Join([]string{t.ModelsPackage, swag.ToGoName(nm)}, ".")
	}
	return swag.ToGoName(nm)
}

func (t *typeResolver) resolveObject(schema *spec.Schema, isAnonymous bool) (result resolvedType, err error) {
	if Debug {
		_, file, pos, _ := runtime.Caller(1)
		log.Printf("%s:%d: resolving object (anon: %t, req: %t)\n", filepath.Base(file), pos, isAnonymous, false) //, bbb)
	}

	result.IsAnonymous = isAnonymous

	result.IsBaseType = schema.Discriminator != ""
	if !isAnonymous {
		result.SwaggerType = object
		result.GoType = t.goTypeName(t.ModelName)
	}
	if len(schema.AllOf) > 0 {
		result.GoType = t.goTypeName(t.ModelName)
		result.IsComplexObject = true
		var isNullable bool
		for _, p := range schema.AllOf {
			if t.IsNullable(&p) {
				isNullable = true
			}
		}
		result.IsNullable = isNullable
		result.SwaggerType = object
		return
	}

	// if this schema has properties, build a map of property name to
	// resolved type, this should also flag the object as anonymous,
	// when a ref is found, the anonymous flag will be reset
	if len(schema.Properties) > 0 {
		result.IsNullable = t.IsNullable(schema)
		result.IsComplexObject = true
		// no return here, still need to check for additional properties
	}

	// account for additional properties
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		et, er := t.ResolveSchema(schema.AdditionalProperties.Schema, true, false)
		if er != nil {
			err = er
			return
		}
		result.IsMap = !result.IsComplexObject
		result.SwaggerType = object
		et.IsNullable = t.IsNullable(schema.AdditionalProperties.Schema)
		result.GoType = "map[string]" + et.GoType
		if et.IsNullable { //&& et.IsComplexObject && !et.IsBaseType {
			result.GoType = "map[string]*" + et.GoType
		}
		t.inferAliasing(&result, schema, isAnonymous, false)
		result.ElemType = &et
		return
	}

	if len(schema.Properties) > 0 {
		return
	}
	result.GoType = iface
	result.IsMap = true
	result.IsMap = !result.IsComplexObject
	result.SwaggerType = object
	result.IsNullable = false
	result.IsInterface = len(schema.Properties) == 0
	return
}

func nullableBool(schema *spec.Schema, isRequired bool) bool {
	if nullable := nullableExtension(schema.Extensions); nullable != nil {
		return *nullable
	}
	required := isRequired && schema.Default == nil && !schema.ReadOnly
	optional := !isRequired && (schema.Default != nil || schema.ReadOnly)

	return required || optional
}

func nullableNumber(schema *spec.Schema, isRequired bool) bool {
	if nullable := nullableExtension(schema.Extensions); nullable != nil {
		return *nullable
	}
	hasDefault := schema.Default != nil && !swag.IsZero(schema.Default)

	isMin := schema.Minimum != nil && *schema.Minimum != 0
	bcMin := schema.Minimum != nil && *schema.Minimum == 0
	isMax := schema.Minimum == nil && (schema.Maximum != nil && *schema.Maximum != 0)
	bcMax := schema.Maximum != nil && *schema.Maximum == 0
	isMinMax := (schema.Minimum != nil && schema.Maximum != nil && *schema.Minimum < *schema.Maximum)
	bcMinMax := (schema.Minimum != nil && schema.Maximum != nil && (*schema.Minimum < 0 && 0 < *schema.Maximum))

	nullable := !schema.ReadOnly && (isRequired || (hasDefault && !(isMin || isMax || isMinMax)) || bcMin || bcMax || bcMinMax)
	return nullable
}

func nullableString(schema *spec.Schema, isRequired bool) bool {
	if nullable := nullableExtension(schema.Extensions); nullable != nil {
		return *nullable
	}
	hasDefault := schema.Default != nil && !swag.IsZero(schema.Default)

	isMin := schema.MinLength != nil && *schema.MinLength != 0
	bcMin := schema.MinLength != nil && *schema.MinLength == 0

	nullable := !schema.ReadOnly && (isRequired || (hasDefault && !isMin) || bcMin)
	return nullable
}

func nullableStrfmt(schema *spec.Schema, isRequired bool) bool {
	notBinary := schema.Format != binary
	if nullable := nullableExtension(schema.Extensions); nullable != nil && notBinary {
		return *nullable
	}
	hasDefault := schema.Default != nil && !swag.IsZero(schema.Default)

	nullable := !schema.ReadOnly && (isRequired || hasDefault)
	return notBinary && nullable
}

func nullableExtension(ext spec.Extensions) *bool {
	if ext == nil {
		return nil
	}

	if boolPtr := boolExtension(ext, xNullable); boolPtr != nil {
		return boolPtr
	}

	return boolExtension(ext, xIsNullable)
}

func boolExtension(ext spec.Extensions, key string) *bool {
	if v, ok := ext[key]; ok {
		if bb, ok := v.(bool); ok {
			return &bb
		}
	}
	return nil
}

func (t *typeResolver) ResolveSchema(schema *spec.Schema, isAnonymous, isRequired bool) (result resolvedType, err error) {
	if Debug {
		// bbb, _ := json.MarshalIndent(schema, "", "  ")
		_, file, pos, _ := runtime.Caller(1)
		log.Printf("%s:%d: resolving schema (anon: %t, req: %t) %s\n", filepath.Base(file), pos, isAnonymous, isRequired, t.ModelName /*bbb*/)
		// tt, _ := json.MarshalIndent(t, "", "  ")
		// log.Println("resolver", string(tt))
	}
	if schema == nil {
		result.IsInterface = true
		result.GoType = iface
		return
	}

	var returns bool
	returns, result, err = t.resolveSchemaRef(schema, isRequired)
	if returns {
		if !isAnonymous {
			result.IsMap = false
			result.IsComplexObject = true
		}
		return
	}

	returns, result, err = t.resolveFormat(schema, isAnonymous, isRequired)
	if returns {
		return
	}

	result.IsNullable = t.isNullable(schema) || isRequired
	tpe := t.firstType(schema)
	switch tpe {
	case array:
		return t.resolveArray(schema, isAnonymous, false)

	case file, number, integer, boolean:
		result.GoType = typeMapping[tpe]
		result.SwaggerType = tpe
		t.inferAliasing(&result, schema, isAnonymous, isRequired)

		switch tpe {
		case boolean:
			result.IsPrimitive = true
			result.IsCustomFormatter = false
			result.IsNullable = nullableBool(schema, isRequired)
		case number, integer:
			result.IsPrimitive = true
			result.IsCustomFormatter = false
			result.IsNullable = nullableNumber(schema, isRequired)
		case file:
		}
		return

	case str:
		result.GoType = str
		result.SwaggerType = str
		t.inferAliasing(&result, schema, isAnonymous, isRequired)

		result.IsPrimitive = true
		result.IsNullable = nullableString(schema, isRequired)
		return

	case object:
		rt, err2 := t.resolveObject(schema, isAnonymous)
		if err2 != nil {
			return resolvedType{}, err2
		}
		rt.HasDiscriminator = schema.Discriminator != ""
		return rt, nil

	default:
		err = fmt.Errorf("unresolvable: %v (format %q)", schema.Type, schema.Format)
		return
	}
}

// A resolvedType is a swagger type that has been resolved and analyzed for usage
// in a template
type resolvedType struct {
	IsAnonymous       bool
	IsArray           bool
	IsMap             bool
	IsInterface       bool
	IsPrimitive       bool
	IsCustomFormatter bool
	IsAliased         bool
	IsNullable        bool
	IsStream          bool
	HasDiscriminator  bool

	// A tuple gets rendered as an anonymous struct with P{index} as property name
	IsTuple            bool
	HasAdditionalItems bool
	IsComplexObject    bool
	IsBaseType         bool

	GoType        string
	AliasedType   string
	SwaggerType   string
	SwaggerFormat string

	ElemType *resolvedType
}

func (rt *resolvedType) Zero() string {
	if zr, ok := zeroes[rt.GoType]; ok {
		return zr
	}
	if rt.IsMap || rt.IsArray {
		return "make(" + rt.GoType + ")"
	}
	if rt.IsTuple || rt.IsComplexObject {
		if rt.IsNullable {
			return "new(" + rt.GoType + ")"
		}
		return rt.GoType + "{}"
	}
	if rt.IsInterface {
		return "nil"
	}

	return ""
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
}
