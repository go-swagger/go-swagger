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
	"github.com/kr/pretty"
)

const (
	iface   = "interface{}"
	array   = "array"
	file    = "file"
	number  = "number"
	integer = "integer"
	boolean = "boolean"
	str     = "string"
	object  = "object"
	binary  = "binary"
	sHTTP   = "http"
	body    = "body"
)

// Extensions supported by go-swagger
const (
	xClass       = "x-class"         // class name used by discriminator
	xGoCustomTag = "x-go-custom-tag" // additional tag for serializers on struct fields
	xGoName      = "x-go-name"       // name of the generated go variable
	xGoType      = "x-go-type"       // reuse existing type (do not generate)
	xIsNullable  = "x-isnullable"
	xNullable    = "x-nullable" // turns the schema into a pointer
	xOmitEmpty   = "x-omitempty"
	xSchemes     = "x-schemes" // additional schemes supported for operations (server generation)
)

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

	if tn == file {
		result.IsPrimitive = true
		result.GoType = typeMapping[binary]
		result.IsStream = true
		return
	}

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

func newTypeResolver(pkg string, doc *loads.Document) *typeResolver {
	resolver := typeResolver{ModelsPackage: pkg, Doc: doc}
	resolver.KnownDefs = make(map[string]struct{}, 64)
	for k, sch := range doc.Spec().Definitions {
		tpe, _, _ := knownDefGoType(k, sch, nil)
		resolver.KnownDefs[tpe] = struct{}{}
	}
	return &resolver
}

func debugLog(format string, args ...interface{}) {
	if Debug {
		_, file, pos, _ := runtime.Caller(2)
		log.Printf("%s:%d: "+format, append([]interface{}{filepath.Base(file), pos}, args...)...)
	}
}

// knownDefGoType returns go type, package and package alias for definition
func knownDefGoType(def string, schema spec.Schema, clear func(string) string) (string, string, string) {
	debugLog("known def type: %q", def)
	ext := schema.Extensions
	if nm, ok := ext.GetString(xGoName); ok {
		if clear == nil {
			debugLog("known def type %s no clear: %q", xGoName, nm)
			return nm, "", ""
		}
		debugLog("known def type %s clear: %q -> %q", xGoName, nm, clear(nm))
		return clear(nm), "", ""
	}
	v, ok := ext[xGoType]
	if !ok {
		if clear == nil {
			debugLog("known def type no clear: %q", def)
			return def, "", ""
		}
		debugLog("known def type clear: %q -> %q", def, clear(def))
		return clear(def), "", ""
	}
	xt := v.(map[string]interface{})
	t := xt["type"].(string)
	imp := xt["import"].(map[string]interface{})
	pkg := imp["package"].(string)
	al, ok := imp["alias"]
	var alias string
	if ok {
		alias = al.(string)
	} else {
		alias = filepath.Base(pkg)
	}
	debugLog("known def type %s no clear: %q", xGoType, alias+"."+t, pkg, alias)
	return alias + "." + t, pkg, alias
}

type typeResolver struct {
	Doc           *loads.Document
	ModelsPackage string
	ModelName     string
	KnownDefs     map[string]struct{}
}

func (t *typeResolver) NewWithModelName(name string) *typeResolver {
	return &typeResolver{
		Doc:           t.Doc,
		ModelsPackage: t.ModelsPackage,
		ModelName:     name,
		KnownDefs:     t.KnownDefs,
	}
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
		var ref *spec.Schema
		var er error

		ref, er = spec.ResolveRef(t.Doc.Spec(), &schema.Ref)
		if er != nil {
			if Debug {
				log.Print("error resolving", er)
			}
			err = er
			return
		}
		res, er := t.ResolveSchema(ref, false, isRequired)
		if er != nil {
			err = er
			return
		}
		result = res

		tn := filepath.Base(schema.Ref.GetURL().Fragment)
		tpe, pkg, alias := knownDefGoType(tn, *ref, t.goTypeName)
		if Debug {
			log.Printf("type name %s, package %s, alias %s", tpe, pkg, alias)
		}
		if tpe != "" {
			result.GoType = tpe
			result.Pkg = pkg
			result.PkgAlias = alias
		}
		result.HasDiscriminator = res.HasDiscriminator
		result.IsBaseType = result.HasDiscriminator
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
			// propagate extensions in resolvedType
			result.Extensions = schema.Extensions

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
	check := func(extension string) (bool, bool) {
		v, found := schema.Extensions[extension]
		nullable, cast := v.(bool)
		return nullable, found && cast
	}

	if nullable, ok := check(xIsNullable); ok {
		return nullable
	}
	if nullable, ok := check(xNullable); ok {
		return nullable
	}
	return len(schema.Properties) > 0
}

func (t *typeResolver) IsEmptyOmitted(schema *spec.Schema) bool {
	v, found := schema.Extensions[xOmitEmpty]
	omitted, cast := v.(bool)
	return found && cast && omitted
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
	result.IsEmptyOmitted = t.IsEmptyOmitted(schema)
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
	result.Extensions = schema.Extensions

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
		log.Printf("%s:%d: resolving object %s (anon: %t, req: %t)\n", filepath.Base(file), pos, t.ModelName, isAnonymous, false) //, bbb)
	}

	result.IsAnonymous = isAnonymous

	result.IsBaseType = schema.Discriminator != ""
	if !isAnonymous {
		result.SwaggerType = object
		tpe, pkg, alias := knownDefGoType(t.ModelName, *schema, t.goTypeName)
		result.GoType = tpe
		result.Pkg = pkg
		result.PkgAlias = alias
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
		sch := schema.AdditionalProperties.Schema
		et, er := t.ResolveSchema(schema.AdditionalProperties.Schema, sch.Ref.String() == "", false)
		if er != nil {
			err = er
			return
		}
		result.IsMap = !result.IsComplexObject
		result.SwaggerType = object

		et.IsNullable = t.isNullable(schema.AdditionalProperties.Schema)
		result.GoType = "map[string]" + et.GoType
		if et.IsNullable {
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

	isMin := schema.Minimum != nil && (*schema.Minimum != 0 || schema.ExclusiveMinimum)
	bcMin := schema.Minimum != nil && *schema.Minimum == 0 && !schema.ExclusiveMinimum
	isMax := schema.Minimum == nil && (schema.Maximum != nil && (*schema.Maximum != 0 || schema.ExclusiveMaximum))
	bcMax := schema.Maximum != nil && *schema.Maximum == 0 && !schema.ExclusiveMaximum
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
	logDebug("resolving schema (anon: %t, req: %t) %s\n", isAnonymous, isRequired, t.ModelName /*bbb*/)
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
			logDebug("not anonymous ref")
		}
		logDebug("returning after ref")
		return
	}

	if t.firstType(schema) == file {
		result.SwaggerType = file
		result.IsPrimitive = true
		result.IsNullable = false
		result.GoType = typeMapping[binary]
		result.IsStream = true
		return
	}

	returns, result, err = t.resolveFormat(schema, isAnonymous, isRequired)
	if returns {
		logDebug("returning after resolve format: %s", pretty.Sprint(result))
		return
	}

	result.IsNullable = t.isNullable(schema) || isRequired
	tpe := t.firstType(schema)
	switch tpe {
	case array:
		return t.resolveArray(schema, isAnonymous, false)

	case file, number, integer, boolean:
		result.Extensions = schema.Extensions
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
		result.Extensions = schema.Extensions
		return

	case object:
		rt, err2 := t.resolveObject(schema, isAnonymous)
		if err2 != nil {
			return resolvedType{}, err2
		}
		rt.HasDiscriminator = schema.Discriminator != ""
		return rt, nil

	case "null":
		result.GoType = iface
		result.SwaggerType = object
		result.IsNullable = false
		result.IsInterface = true
		return

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
	IsEmptyOmitted    bool

	// A tuple gets rendered as an anonymous struct with P{index} as property name
	IsTuple            bool
	HasAdditionalItems bool
	IsComplexObject    bool
	IsBaseType         bool

	GoType        string
	Pkg           string
	PkgAlias      string
	AliasedType   string
	SwaggerType   string
	SwaggerFormat string
	Extensions    spec.Extensions

	ElemType *resolvedType
}

func (rt *resolvedType) Zero() string {
	// if type is aliased, provide zero from the aliased type
	if rt.IsAliased {
		if zr, ok := zeroes[rt.AliasedType]; ok {
			return rt.GoType + "(" + zr + ")"
		}
	}
	// zero function provided as native or by strfmt function
	if zr, ok := zeroes[rt.GoType]; ok {
		return zr
	}
	// map and slice initializer
	if rt.IsMap || rt.IsArray {
		return "make(" + rt.GoType + ", 0, 50)"
	}
	// object initializer
	if rt.IsTuple || rt.IsComplexObject {
		if rt.IsNullable {
			return "new(" + rt.GoType + ")"
		}
		return rt.GoType + "{}"
	}
	// interface initializer
	if rt.IsInterface {
		return "nil"
	}

	return ""
}
