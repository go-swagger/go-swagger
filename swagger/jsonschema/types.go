// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           sigu-399
// author-github    https://github.com/sigu-399
// author-mail      sigu.399@gmail.com
//
// repository-name  jsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description      Contains consts, types and (most) error messages.
//
// created          28-02-2013

package jsonschema

const (
	// KeySchema the schema property name
	KeySchema = "$schema"
	// KeyID the id property name
	KeyID = "$id"
	// KeyRef the ref property name
	KeyRef = "$ref"
	// KeyTitle the title property name
	KeyTitle = "title"
	// KeyDescription the description property name
	KeyDescription = "description"
	// KeyType the type property name
	KeyType = "type"
	// KeyItems the items property name
	KeyItems = "items"
	// KeyAdditionalItems the additional items property name
	KeyAdditionalItems = "additionalItems"
	// KeyProperties the properties property name
	KeyProperties = "properties"
	// KeyPatternProperties the pattern properties property name
	KeyPatternProperties = "patternProperties"
	// KeyAdditionalProperties the additional properties property name
	KeyAdditionalProperties = "additionalProperties"
	// KeyDefinitions the definitions property name
	KeyDefinitions = "definitions"
	// KeyMultipleOf the multipleOf property name
	KeyMultipleOf = "multipleOf"
	// KeyMinimum the minimum property name
	KeyMinimum = "minimum"
	// KeyMaximum the maximum property name
	KeyMaximum = "maximum"
	// KeyExclusiveMinimum the exclusiveMinimum property name
	KeyExclusiveMinimum = "exclusiveMinimum"
	// KeyExclusiveMaximum the exclusiveMaximum property name
	KeyExclusiveMaximum = "exclusiveMaximum"
	// KeyMinLength the minLength property name
	KeyMinLength = "minLength"
	// KeyMaxLength the maxLength property name
	KeyMaxLength = "maxLength"
	// KeyPattern the pattern property name
	KeyPattern = "pattern"
	// KeyMinProperties the minProperties property name
	KeyMinProperties = "minProperties"
	// KeyMaxProperties the maxProperties property name
	KeyMaxProperties = "maxProperties"
	// KeyDependencies the dependencies property name
	KeyDependencies = "dependencies"
	// KeyRequired the required property name
	KeyRequired = "required"
	// KeyMinItems the minItems property name
	KeyMinItems = "minItems"
	// KeyMaxItems the maxItems property name
	KeyMaxItems = "maxItems"
	// KeyUniqueItems the uniqueItems property name
	KeyUniqueItems = "uniqueItems"
	// KeyEnum the enum property name
	KeyEnum = "enum"
	// KeyOneOf the oneOf property name
	KeyOneOf = "oneOf"
	// KeyAnyOf the anyOf property name
	KeyAnyOf = "anyOf"
	// KeyAllOf the allOf property name
	KeyAllOf = "allOf"
	// KeyNot the not property name
	KeyNot = "not"
)

const (
	typeString         = "string"
	typeBoolean        = "boolean"
	typeArrayOfStrings = "array of strings"
	typeArrayOfSchemas = "array of schemas"
	typeObject         = "object"
	typeSchema         = "schema"
)

const (
	stringSchemaOrArrayOfStrings = "schema or array of strings"
	stringProperties             = "properties"
	stringDependency             = "dependency"
)

const (
	contextRoot        = "(root)"
	rootSchemaProperty = "(root)"
)

const (
	errMessageXMustBeOfTypeY                   = `%s must be of type %s`
	errMessageXIsMissingAndRequired            = `%s is missing and required`
	errMessageMustBeOfTypeX                    = `must be of type %s`
	errMessageArrayItemsMustBeUnique           = `array items must be unique`
	errMessageDoesNotMatchPattern              = `does not match pattern '%s'`
	errMessageMustMatchOneEnumValue            = `must match one of the enum values [%s]`
	errMessageStringLengthMustBeGreaterOrEqual = `string length must be greater than or equal to %d`
	errMessageStringLengthMustBeLessOrEqual    = `string length must be less than or equal to %d`
	errMessageNumberMustBeLessOrEqual          = `must be less than or equal to %s`
	errMessageNumberMustBeLess                 = `must be less than %s`
	errMessageNumberMustBeGreaterOrEqual       = `must be greater than or equal to %s`
	errMessageNumberMustBeGreater              = `must be greater than %s`
	errMessageMustValidateAllOf                = `must validate all the schemas (allOf)`
	errMessageMustValidateOneOf                = `must validate one and only one schema (oneOf)`
	errMessageMustValidateAnyOf                = `must validate at least one schema (anyOf)`
	errMessageMustValidateNot                  = `must not validate the schema (not)`
	errMessageArrayMinItems                    = `array must have at least %d items`
	errMessageArrayMaxItems                    = `array must have at the most %d items`
	errMessageArrayMinProperties               = `must have at least %d properties`
	errMessageArrayMaxProperties               = `must have at the most %d properties`
	errMessageHasDependencyOn                  = `has a dependency on %s`
	errMessageMultipleOf                       = `must be a multiple of %s`
	errMessageArrayNoAdditionalItems           = `no additional item allowed on array`
	errMessageAdditionalPropertyNotAllowed     = `additional property "%s" is not allowed`
	errMessageInvalidPatternProperty           = `property "%s" does not match pattern %s`
	errMessageInternal                         = `internal error %s`
)

const (
	// JSONTypeArray name for JSON/Schema array type
	JSONTypeArray = `array`
	// JSONTypeBoolean name for JSON/Schema boolean type
	JSONTypeBoolean = `boolean`
	// JSONTypeInteger name for JSON/Schema integer type
	JSONTypeInteger = `integer`
	// JSONTypeNumber name for JSON/Schema number type
	JSONTypeNumber = `number`
	// JSONTypeNull name for JSON/Schema null type
	JSONTypeNull = `null`
	// JSONTypeObject name for JSON/Schema object type
	JSONTypeObject = `object`
	// JSONTypeString name for JSON/Schema string type
	JSONTypeString = `string`
)

// JSONTypes the json types
var JSONTypes []string

// SchemaTypes the schema types
var SchemaTypes []string

func init() {
	JSONTypes = []string{
		JSONTypeArray,
		JSONTypeBoolean,
		JSONTypeInteger,
		JSONTypeNumber,
		JSONTypeNull,
		JSONTypeObject,
		JSONTypeString}

	SchemaTypes = []string{
		JSONTypeArray,
		JSONTypeBoolean,
		JSONTypeInteger,
		JSONTypeNumber,
		JSONTypeObject,
		JSONTypeString}
}
