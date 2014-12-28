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
// description      Extends JsonSchemaDocument and jsonSchema, implements the validation phase.
//
// created          28-02-2013

package jsonschema

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (v *JsonSchemaDocument) Validate(document interface{}) *ValidationResult {
	result := &ValidationResult{}
	context := consJsonContext(CONTEXT_ROOT, nil)
	v.rootSchema.validateRecursive(v.rootSchema, document, result, context)
	return result
}

func (v *jsonSchema) Validate(document interface{}, context *jsonContext) *ValidationResult {
	result := &ValidationResult{}
	v.validateRecursive(v, document, result, context)
	return result
}

func convertDocumentNode(val interface{}) interface{} {
	if lval, ok := val.([]interface{}); ok {
		res := []interface{}{}
		for _, v := range lval {
			res = append(res, convertDocumentNode(v))
		}
		return res
	}
	if mval, ok := val.(map[interface{}]interface{}); ok {
		res := map[string]interface{}{}
		for k, v := range mval {
			res[k.(string)] = convertDocumentNode(v)
		}
		return res
	}
	return val
}

// Walker function to validate the json recursively against the schema
func (v *jsonSchema) validateRecursive(currentSchema *jsonSchema, currentNode interface{}, result *ValidationResult, context *jsonContext) {

	// Handle referenced schemas, returns directly when a $ref is found
	if currentSchema.refSchema != nil {
		v.validateRecursive(currentSchema.refSchema, currentNode, result, context)
		return
	}

	// Check for null value
	if currentNode == nil {

		if currentSchema.types.IsTyped() && !currentSchema.types.Contains(TYPE_NULL) {
			result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
			return
		}

		currentSchema.validateSchema(currentSchema, currentNode, result, context)
		v.validateCommon(currentSchema, currentNode, result, context)

	} else { // Not a null value

		rValue := reflect.ValueOf(currentNode)
		rKind := rValue.Kind()

		switch rKind {

		// Slice => JSON array

		case reflect.Slice:

			if currentSchema.types.IsTyped() && !currentSchema.types.Contains(TYPE_ARRAY) {
				result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
				return
			}

			castCurrentNode := currentNode.([]interface{})

			currentSchema.validateSchema(currentSchema, castCurrentNode, result, context)

			v.validateArray(currentSchema, castCurrentNode, result, context)
			v.validateCommon(currentSchema, castCurrentNode, result, context)

		// Map => JSON object

		case reflect.Map:
			if currentSchema.types.IsTyped() && !currentSchema.types.Contains(TYPE_OBJECT) {
				result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
				return
			}

			castCurrentNode, ok := currentNode.(map[string]interface{})
			if !ok {
				castCurrentNode = convertDocumentNode(currentNode).(map[string]interface{})
			}

			currentSchema.validateSchema(currentSchema, castCurrentNode, result, context)

			v.validateObject(currentSchema, castCurrentNode, result, context)
			v.validateCommon(currentSchema, castCurrentNode, result, context)

			for _, pSchema := range currentSchema.propertiesChildren {
				nextNode, ok := castCurrentNode[pSchema.property]
				if ok {
					subContext := consJsonContext(pSchema.property, context)
					v.validateRecursive(pSchema, nextNode, result, subContext)
				}
			}

		// Simple JSON values : string, number, boolean

		case reflect.Bool:

			if currentSchema.types.IsTyped() && !currentSchema.types.Contains(TYPE_BOOLEAN) {
				result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
				return
			}

			value := currentNode.(bool)

			currentSchema.validateSchema(currentSchema, value, result, context)
			v.validateNumber(currentSchema, value, result, context)
			v.validateCommon(currentSchema, value, result, context)
			v.validateString(currentSchema, value, result, context)

		case reflect.String:

			if currentSchema.types.IsTyped() && !currentSchema.types.Contains(TYPE_STRING) {
				result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
				return
			}

			value := currentNode.(string)

			currentSchema.validateSchema(currentSchema, value, result, context)
			v.validateNumber(currentSchema, value, result, context)
			v.validateCommon(currentSchema, value, result, context)
			v.validateString(currentSchema, value, result, context)

		case reflect.Float64:

			value := currentNode.(float64)

			// Note: JSON only understand one kind of numeric ( can be float or int )
			// JSON schema make a distinction between fload and int
			// An integer can be a number, but a number ( with decimals ) cannot be an integer
			isInteger := isFloat64AnInteger(value)
			validType := currentSchema.types.Contains(TYPE_NUMBER) || (isInteger && currentSchema.types.Contains(TYPE_INTEGER))

			if currentSchema.types.IsTyped() && !validType {
				result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_MUST_BE_OF_TYPE_X, currentSchema.types.String()))
				return
			}

			currentSchema.validateSchema(currentSchema, value, result, context)
			v.validateNumber(currentSchema, value, result, context)
			v.validateCommon(currentSchema, value, result, context)
			v.validateString(currentSchema, value, result, context)
		}
	}

	result.incrementScore()
}

// Different kinds of validation there, schema / common / array / object / string...
func (v *jsonSchema) validateSchema(currentSchema *jsonSchema, currentNode interface{}, result *ValidationResult, context *jsonContext) {

	if len(currentSchema.anyOf) > 0 {

		validatedAnyOf := false
		var bestValidationResult *ValidationResult

		for _, anyOfSchema := range currentSchema.anyOf {
			if !validatedAnyOf {
				validationResult := anyOfSchema.Validate(currentNode, context)
				validatedAnyOf = validationResult.Valid()

				if !validatedAnyOf && (bestValidationResult == nil || validationResult.score > bestValidationResult.score) {
					bestValidationResult = validationResult
				}
			}
		}
		if !validatedAnyOf {

			result.addError(context, currentNode, ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ANYOF)

			if bestValidationResult != nil {
				// add error messages of closest matching schema as
				// that's probably the one the user was trying to match
				result.mergeErrors(bestValidationResult)
			}
		}
	}

	if len(currentSchema.oneOf) > 0 {

		nbValidated := 0
		var bestValidationResult *ValidationResult

		for _, oneOfSchema := range currentSchema.oneOf {
			validationResult := oneOfSchema.Validate(currentNode, context)
			if validationResult.Valid() {
				nbValidated++
			} else if nbValidated == 0 && (bestValidationResult == nil || validationResult.score > bestValidationResult.score) {
				bestValidationResult = validationResult
			}
		}

		if nbValidated != 1 {

			result.addError(context, currentNode, ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ONEOF)

			if nbValidated == 0 {
				// add error messages of closest matching schema as
				// that's probably the one the user was trying to match
				result.mergeErrors(bestValidationResult)
			}
		}

	}

	if len(currentSchema.allOf) > 0 {
		nbValidated := 0

		for _, allOfSchema := range currentSchema.allOf {
			validationResult := allOfSchema.Validate(currentNode, context)
			if validationResult.Valid() {
				nbValidated++
			}
			result.mergeErrors(validationResult)
		}

		if nbValidated != len(currentSchema.allOf) {
			result.addError(context, currentNode, ERROR_MESSAGE_NUMBER_MUST_VALIDATE_ALLOF)
		}
	}

	if currentSchema.not != nil {
		validationResult := currentSchema.not.Validate(currentNode, context)
		if validationResult.Valid() {
			result.addError(context, currentNode, ERROR_MESSAGE_NUMBER_MUST_VALIDATE_NOT)
		}
	}

	if currentSchema.dependencies != nil && len(currentSchema.dependencies) > 0 {
		if isKind(currentNode, reflect.Map) {
			for elementKey := range currentNode.(map[string]interface{}) {
				if dependency, ok := currentSchema.dependencies[elementKey]; ok {
					switch dependency := dependency.(type) {

					case []string:
						for _, dependOnKey := range dependency {
							if _, dependencyResolved := currentNode.(map[string]interface{})[dependOnKey]; !dependencyResolved {
								result.addError(context, currentNode, fmt.Sprintf(ERROR_MESSAGE_HAS_DEPENDENCY_ON, dependOnKey))
							}
						}

					case *jsonSchema:
						dependency.validateRecursive(dependency, currentNode, result, context)

					}
				}
			}
		}
	}

	result.incrementScore()
}

func (v *jsonSchema) validateCommon(currentSchema *jsonSchema, value interface{}, result *ValidationResult, context *jsonContext) {

	// enum:
	if len(currentSchema.enum) > 0 {
		has, err := currentSchema.ContainsEnum(value)
		if err != nil {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_INTERNAL, err))
		}
		if !has {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_MUST_MATCH_ONE_ENUM_VALUES, strings.Join(currentSchema.enum, ",")))
		}
	}

	result.incrementScore()
}

func (v *jsonSchema) validateArray(currentSchema *jsonSchema, value []interface{}, result *ValidationResult, context *jsonContext) {

	nbItems := len(value)

	// TODO explain
	if currentSchema.itemsChildrenIsSingleSchema {
		for i := range value {
			subContext := consJsonContext(strconv.Itoa(i), context)
			validationResult := currentSchema.itemsChildren[0].Validate(value[i], subContext)
			result.mergeErrors(validationResult)
		}
	} else {
		if currentSchema.itemsChildren != nil && len(currentSchema.itemsChildren) > 0 {

			nbItems := len(currentSchema.itemsChildren)
			nbValues := len(value)

			if nbItems == nbValues {
				for i := 0; i != nbItems; i++ {
					subContext := consJsonContext(strconv.Itoa(i), context)
					validationResult := currentSchema.itemsChildren[i].Validate(value[i], subContext)
					result.mergeErrors(validationResult)
				}
			} else if nbItems < nbValues {
				switch currentSchema.additionalItems.(type) {
				case bool:
					if !currentSchema.additionalItems.(bool) {
						result.addError(context, value, ERROR_MESSAGE_ARRAY_NO_ADDITIONAL_ITEM)
					}
				case *jsonSchema:
					additionalItemSchema := currentSchema.additionalItems.(*jsonSchema)
					for i := nbItems; i != nbValues; i++ {
						subContext := consJsonContext(strconv.Itoa(i), context)
						validationResult := additionalItemSchema.Validate(value[i], subContext)
						result.mergeErrors(validationResult)
					}
				}
			}
		}
	}

	// minItems & maxItems
	if currentSchema.minItems != nil {
		if nbItems < *currentSchema.minItems {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ARRAY_MIN_ITEMS, *currentSchema.minItems))
		}
	}
	if currentSchema.maxItems != nil {
		if nbItems > *currentSchema.maxItems {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ARRAY_MAX_ITEMS, *currentSchema.maxItems))
		}
	}

	// uniqueItems:
	if currentSchema.uniqueItems {
		var stringifiedItems []string
		for _, v := range value {
			vString, err := marshalToJsonString(v)
			if err != nil {
				result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_INTERNAL, err))
			}
			if isStringInSlice(stringifiedItems, *vString) {
				result.addError(context, value, ERROR_MESSAGE_ARRAY_ITEMS_MUST_BE_UNIQUE)
			}
			stringifiedItems = append(stringifiedItems, *vString)
		}
	}

	result.incrementScore()
}

func (v *jsonSchema) validateObject(currentSchema *jsonSchema, value map[string]interface{}, result *ValidationResult, context *jsonContext) {

	// minProperties & maxProperties:
	if currentSchema.minProperties != nil {
		if len(value) < *currentSchema.minProperties {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ARRAY_MIN_PROPERTIES, *currentSchema.minProperties))
		}
	}
	if currentSchema.maxProperties != nil {
		if len(value) > *currentSchema.maxProperties {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ARRAY_MAX_PROPERTIES, *currentSchema.maxProperties))
		}
	}

	// required:
	for _, requiredProperty := range currentSchema.required {
		_, ok := value[requiredProperty]
		if ok {
			result.incrementScore()
		} else {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_X_IS_MISSING_AND_REQUIRED, fmt.Sprintf(`"%s" property`, requiredProperty)))
		}
	}

	// additionalProperty & patternProperty:
	if currentSchema.additionalProperties != nil {

		switch currentSchema.additionalProperties.(type) {
		case bool:

			if !currentSchema.additionalProperties.(bool) {

				for pk := range value {

					found := false
					for _, spValue := range currentSchema.propertiesChildren {
						if pk == spValue.property {
							found = true
						}
					}

					pp_has, pp_match := v.validatePatternProperty(currentSchema, pk, value[pk], result, context)

					if found {

						if pp_has && !pp_match {
							result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ADDITIONAL_PROPERTY_NOT_ALLOWED, pk))
						}

					} else {

						if !pp_has || !pp_match {
							result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_ADDITIONAL_PROPERTY_NOT_ALLOWED, pk))
						}

					}
				}
			}

		case *jsonSchema:

			additionalPropertiesSchema := currentSchema.additionalProperties.(*jsonSchema)
			for pk := range value {

				found := false
				for _, spValue := range currentSchema.propertiesChildren {
					if pk == spValue.property {
						found = true
					}
				}

				pp_has, pp_match := v.validatePatternProperty(currentSchema, pk, value[pk], result, context)

				if found {

					if pp_has && !pp_match {
						validationResult := additionalPropertiesSchema.Validate(value[pk], context)
						result.mergeErrors(validationResult)
					}

				} else {

					if !pp_has || !pp_match {
						validationResult := additionalPropertiesSchema.Validate(value[pk], context)
						result.mergeErrors(validationResult)
					}

				}

			}
		}
	} else {

		for pk := range value {

			pp_has, pp_match := v.validatePatternProperty(currentSchema, pk, value[pk], result, context)

			if pp_has && !pp_match {

				result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_INVALID_PATTERN_PROPERTY, pk, currentSchema.PatternPropertiesString()))
			}

		}
	}

	result.incrementScore()
}

func (v *jsonSchema) validatePatternProperty(currentSchema *jsonSchema, key string, value interface{}, result *ValidationResult, context *jsonContext) (has bool, matched bool) {

	has = false

	validatedkey := false

	for pk, pv := range currentSchema.patternProperties {
		if matches, _ := regexp.MatchString(pk, key); matches {
			has = true
			subContext := consJsonContext(key, context)
			validationResult := pv.Validate(value, subContext)
			result.mergeErrors(validationResult)
			if validationResult.Valid() {
				validatedkey = true
			}
		}
	}

	if !validatedkey {
		return has, false
	}

	result.incrementScore()

	return has, true
}

func (v *jsonSchema) validateString(currentSchema *jsonSchema, value interface{}, result *ValidationResult, context *jsonContext) {

	// Ignore non strings
	if !isKind(value, reflect.String) {
		return
	}

	stringValue := value.(string)

	// minLength & maxLength:
	if currentSchema.minLength != nil {
		if utf8.RuneCount([]byte(stringValue)) < *currentSchema.minLength {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_STRING_LENGTH_MUST_BE_GREATER_OR_EQUAL, *currentSchema.minLength))
		}
	}
	if currentSchema.maxLength != nil {
		if utf8.RuneCount([]byte(stringValue)) > *currentSchema.maxLength {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_STRING_LENGTH_MUST_BE_LOWER_OR_EQUAL, *currentSchema.maxLength))
		}
	}

	// pattern:
	if currentSchema.pattern != nil {
		if !currentSchema.pattern.MatchString(stringValue) {
			result.addError(context, value, fmt.Sprintf(ERROR_MESSAGE_DOES_NOT_MATCH_PATTERN, currentSchema.pattern))

		}
	}

	result.incrementScore()
}

func (v *jsonSchema) validateNumber(currentSchema *jsonSchema, value interface{}, result *ValidationResult, context *jsonContext) {

	// Ignore non numbers
	if !isKind(value, reflect.Float64) {
		return
	}

	float64Value := value.(float64)

	// multipleOf:
	if currentSchema.multipleOf != nil {
		if !isFloat64AnInteger(float64Value / *currentSchema.multipleOf) {
			result.addError(context, validationErrorFormatNumber(float64Value), fmt.Sprintf(ERROR_MESSAGE_MULTIPLE_OF, validationErrorFormatNumber(*currentSchema.multipleOf)))
		}
	}

	//maximum & exclusiveMaximum:
	if currentSchema.maximum != nil {
		if currentSchema.exclusiveMaximum {
			if float64Value >= *currentSchema.maximum {
				result.addError(context, validationErrorFormatNumber(float64Value), fmt.Sprintf(ERROR_MESSAGE_NUMBER_MUST_BE_LOWER_OR_EQUAL, validationErrorFormatNumber(*currentSchema.maximum)))
			}
		} else {
			if float64Value > *currentSchema.maximum {
				result.addError(context, validationErrorFormatNumber(float64Value), fmt.Sprintf(ERROR_MESSAGE_NUMBER_MUST_BE_LOWER, validationErrorFormatNumber(*currentSchema.maximum)))
			}
		}
	}

	//minimum & exclusiveMinimum:
	if currentSchema.minimum != nil {
		if currentSchema.exclusiveMinimum {
			if float64Value <= *currentSchema.minimum {
				result.addError(context, validationErrorFormatNumber(float64Value), fmt.Sprintf(ERROR_MESSAGE_NUMBER_MUST_BE_GREATER_OR_EQUAL, validationErrorFormatNumber(*currentSchema.minimum)))
			}
		} else {
			if float64Value < *currentSchema.minimum {
				result.addError(context, validationErrorFormatNumber(float64Value), fmt.Sprintf(ERROR_MESSAGE_NUMBER_MUST_BE_GREATER, validationErrorFormatNumber(*currentSchema.minimum)))
			}
		}
	}

	result.incrementScore()
}

type ValidationError struct {
	Context     *jsonContext // Tree like notation of the part that failed the validation. ex (root).a.b ...
	Description string       // A human readable error message
	Value       interface{}  // Value given by the JSON file that is the source of the error
}

func (v ValidationError) String() string {

	// as a fallback, the value is displayed go style
	valueString := fmt.Sprintf("%v", v.Value)

	// marshall the go value value to json
	if v.Value == nil {
		valueString = TYPE_NULL
	} else {
		if vs, err := marshalToJsonString(v.Value); err == nil {
			if vs == nil {
				valueString = TYPE_NULL
			} else {
				valueString = *vs
			}
		}
	}

	return fmt.Sprintf("%s : %s, given %s", v.Context, v.Description, valueString)
}

type ValidationResult struct {
	errors []ValidationError
	// Scores how well the validation matched. Useful in generating
	// better error messages for anyOf and oneOf.
	score int
}

func (v *ValidationResult) Valid() bool {
	return len(v.errors) == 0
}

func (v *ValidationResult) Errors() []ValidationError {
	return v.errors
}

func (v *ValidationResult) addError(context *jsonContext, value interface{}, description string) {
	v.errors = append(v.errors, ValidationError{Context: context, Value: value, Description: description})
	v.score -= 2 // results in a net -1 when added to the +1 we get at the end of the validation function
}

// Used to copy errors from a sub-schema to the main one
func (v *ValidationResult) mergeErrors(otherResult *ValidationResult) {
	v.errors = append(v.errors, otherResult.Errors()...)
	v.score += otherResult.score
}

func (v *ValidationResult) incrementScore() {
	v.score++
}
