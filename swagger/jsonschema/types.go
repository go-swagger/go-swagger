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
	KEY_SCHEMA                = "$schema"
	KEY_ID                    = "$id"
	KEY_REF                   = "$ref"
	KEY_TITLE                 = "title"
	KEY_DESCRIPTION           = "description"
	KEY_TYPE                  = "type"
	KEY_ITEMS                 = "items"
	KEY_ADDITIONAL_ITEMS      = "additionalItems"
	KEY_PROPERTIES            = "properties"
	KEY_PATTERN_PROPERTIES    = "patternProperties"
	KEY_ADDITIONAL_PROPERTIES = "additionalProperties"
	KEY_DEFINITIONS           = "definitions"
	KEY_MULTIPLE_OF           = "multipleOf"
	KEY_MINIMUM               = "minimum"
	KEY_MAXIMUM               = "maximum"
	KEY_EXCLUSIVE_MINIMUM     = "exclusiveMinimum"
	KEY_EXCLUSIVE_MAXIMUM     = "exclusiveMaximum"
	KEY_MIN_LENGTH            = "minLength"
	KEY_MAX_LENGTH            = "maxLength"
	KEY_PATTERN               = "pattern"
	KEY_MIN_PROPERTIES        = "minProperties"
	KEY_MAX_PROPERTIES        = "maxProperties"
	KEY_DEPENDENCIES          = "dependencies"
	KEY_REQUIRED              = "required"
	KEY_MIN_ITEMS             = "minItems"
	KEY_MAX_ITEMS             = "maxItems"
	KEY_UNIQUE_ITEMS          = "uniqueItems"
	KEY_ENUM                  = "enum"
	KEY_ONE_OF                = "oneOf"
	KEY_ANY_OF                = "anyOf"
	KEY_ALL_OF                = "allOf"
	KEY_NOT                   = "not"

	STRING_STRING                     = "string"
	STRING_BOOLEAN                    = "boolean"
	STRING_ARRAY_OF_STRINGS           = "array of strings"
	STRING_ARRAY_OF_SCHEMAS           = "array of schemas"
	STRING_OBJECT                     = "object"
	STRING_SCHEMA                     = "schema"
	STRING_SCHEMA_OR_ARRAY_OF_STRINGS = "schema or array of strings"
	STRING_PROPERTIES                 = "properties"
	STRING_DEPENDENCY                 = "dependency"

	CONTEXT_ROOT         = "(root)"
	ROOT_SCHEMA_PROPERTY = "(root)"
)

const (
	ERROR_MESSAGE_X_MUST_BE_OF_TYPE_Y = `%s must be of type %s`

	ERROR_MESSAGE_X_IS_MISSING_AND_REQUIRED  = `%s is missing and required`
	ERROR_MESSAGE_MUST_BE_OF_TYPE_X          = `must be of type %s`
	ERROR_MESSAGE_ARRAY_ITEMS_MUST_BE_UNIQUE = `array items must be unique`
	ERROR_MESSAGE_DOES_NOT_MATCH_PATTERN     = `does not match pattern '%s'`
	ERROR_MESSAGE_MUST_MATCH_ONE_ENUM_VALUES = `must match one of the enum values [%s]`

	ERROR_MESSAGE_STRING_LENGTH_MUST_BE_GREATER_OR_EQUAL = `string length must be greater or equal to %d`
	ERROR_MESSAGE_STRING_LENGTH_MUST_BE_LOWER_OR_EQUAL   = `string length must be lower or equal to %d`

	ERROR_MESSAGE_NUMBER_MUST_BE_LOWER_OR_EQUAL   = `must be lower than or equal to %s`
	ERROR_MESSAGE_NUMBER_MUST_BE_LOWER            = `must be lower than %s`
	ERROR_MESSAGE_NUMBER_MUST_BE_GREATER_OR_EQUAL = `must be greater than or equal to %s`
	ERROR_MESSAGE_NUMBER_MUST_BE_GREATER          = `must be greater than %s`

	ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ALLOF = `must validate all the schemas (allOf)`
	ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ONEOF = `must validate one and only one schema (oneOf)`
	ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ANYOF = `must validate at least one schema (anyOf)`
	ERROR_MESSAGE_NUMBER_MUST_VALIDATE_NOT   = `must not validate the schema (not)`

	ERROR_MESSAGE_ARRAY_MIN_ITEMS = `array must have at least %d items`
	ERROR_MESSAGE_ARRAY_MAX_ITEMS = `array must have at the most %d items`

	ERROR_MESSAGE_ARRAY_MIN_PROPERTIES = `must have at least %d properties`
	ERROR_MESSAGE_ARRAY_MAX_PROPERTIES = `must have at the most %d properties`

	ERROR_MESSAGE_HAS_DEPENDENCY_ON = `has a dependency on %s`

	ERROR_MESSAGE_MULTIPLE_OF = `must be a multiple of %s`

	ERROR_MESSAGE_ARRAY_NO_ADDITIONAL_ITEM = `no additional item allowed on array`

	ERROR_MESSAGE_ADDITIONAL_PROPERTY_NOT_ALLOWED = `additional property "%s" is not allowed`
	ERROR_MESSAGE_INVALID_PATTERN_PROPERTY        = `property "%s" does not match pattern %s`

	ERROR_MESSAGE_INTERNAL = `internal error %s`
)

const (
	TYPE_ARRAY   = `array`
	TYPE_BOOLEAN = `boolean`
	TYPE_INTEGER = `integer`
	TYPE_NUMBER  = `number`
	TYPE_NULL    = `null`
	TYPE_OBJECT  = `object`
	TYPE_STRING  = `string`
)

var JSON_TYPES []string
var SCHEMA_TYPES []string

func init() {
	JSON_TYPES = []string{
		TYPE_ARRAY,
		TYPE_BOOLEAN,
		TYPE_INTEGER,
		TYPE_NUMBER,
		TYPE_NULL,
		TYPE_OBJECT,
		TYPE_STRING}

	SCHEMA_TYPES = []string{
		TYPE_ARRAY,
		TYPE_BOOLEAN,
		TYPE_INTEGER,
		TYPE_NUMBER,
		TYPE_OBJECT,
		TYPE_STRING}
}
