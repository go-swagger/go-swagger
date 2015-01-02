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
// description      Defines schemaDocument, the main entry to every schemas.
//                  Contains the parsing logic and error checking.
//
// created          26-02-2013

package jsonschema

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/casualjim/go-swagger/jsonreference"
)

// Load loads a json schema document from a loader
func Load(loader Loader) (*Document, error) {
	var err error
	d := Document{}
	d.pool = newSchemaPool()
	d.referencePool = newSchemaReferencePool()

	d.documentReference, err = jsonreference.New(loader.URL())
	if err != nil {
		return nil, err
	}

	spd, err := d.pool.GetDocumentFromLoader(d.documentReference, loader)
	if err != nil {
		return nil, err
	}

	err = d.parse(spd.Document)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// New creates a new schema document for the given interface
func New(document interface{}) (*Document, error) {

	internalLog("New schema document :")

	var err error

	d := Document{}
	d.pool = newSchemaPool()
	d.referencePool = newSchemaReferencePool()

	switch document.(type) {

	// document is a reference, file or http scheme
	case string:

		internalLog(fmt.Sprintf(" From http or file (%s)", document.(string)))

		d.documentReference, err = jsonreference.New(document.(string))
		spd, err := d.pool.GetDocument(d.documentReference)
		if err != nil {
			return nil, err
		}

		err = d.parse(spd.Document)
		if err != nil {
			return nil, err
		}

	// document is json
	case map[string]interface{}:

		internalLog(" From map")

		d.documentReference, err = jsonreference.New("#")
		d.pool.SetStandaloneDocument(document)
		if err != nil {
			return nil, err
		}

		err = d.parse(document.(map[string]interface{}))
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("Invalid argument, must be a jsonReference string or Json as map[string]interface{}")
	}

	return &d, nil
}

// Document is a json schema document
type Document struct {
	documentReference jsonreference.Ref
	rootSchema        *jsonSchema
	pool              *schemaPool
	referencePool     *schemaReferencePool
}

func (d *Document) parse(document interface{}) error {
	d.rootSchema = &jsonSchema{property: rootSchemaProperty}
	return d.parseSchema(document, d.rootSchema)
}

// SetRootSchemaName sets the root schema name
func (d *Document) SetRootSchemaName(name string) {
	d.rootSchema.property = name
}

// Parses a schema
//
// Pretty long function ( sorry :) )... but pretty straight forward, repetitive and boring
// Not much magic involved here, most of the job is to validate the key names and their values,
// then the values are copied into schema struct
//
func (d *Document) parseSchema(documentNode interface{}, currentSchema *jsonSchema) error {

	if internalLogEnabled {
		documentJSON, err := toJSONString(documentNode)
		if err == nil && documentJSON != nil {
			internalLog(fmt.Sprintf("Parsing schema %s", *documentJSON))
		} else {
			internalLog(fmt.Sprintf("Parsing schema %v", documentNode))
		}
	}

	if !isKind(documentNode, reflect.Map) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, typeSchema, typeObject)
	}

	m := documentNode.(map[string]interface{})

	if currentSchema == d.rootSchema {
		currentSchema.ref = &d.documentReference
	}

	// $schema
	if existsMapKey(m, KeySchema) {
		if !isKind(m[KeySchema], reflect.String) {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeySchema, typeString)
		}
		schemaRef := m[KeySchema].(string)
		schemaReference, err := jsonreference.New(schemaRef)
		currentSchema.schema = &schemaReference
		if err != nil {
			return err
		}
	}

	// $ref
	if existsMapKey(m, KeyRef) && !isKind(m[KeyRef], reflect.String) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, KeyRef, typeString)
	}
	if k, ok := m[KeyRef].(string); ok {

		if sch, ok := d.referencePool.GetSchema(currentSchema.ref.String() + k); ok {

			currentSchema.refSchema = sch

		} else {

			var err error
			err = d.parseReference(documentNode, currentSchema, k)
			if err != nil {
				return err
			}

			return nil
		}
	}

	// definitions
	if existsMapKey(m, KeyDefinitions) {
		if isKind(m[KeyDefinitions], reflect.Map) {
			currentSchema.definitions = make(map[string]*jsonSchema)
			for dk, dv := range m[KeyDefinitions].(map[string]interface{}) {
				if isKind(dv, reflect.Map) {
					newSchema := &jsonSchema{property: KeyDefinitions, parent: currentSchema, ref: currentSchema.ref}
					currentSchema.definitions[dk] = newSchema
					err := d.parseSchema(dv, newSchema)
					if err != nil {
						return errors.New(err.Error())
					}
				} else {
					return fmt.Errorf(errMessageXMustBeOfTypeY, KeyDefinitions, typeArrayOfSchemas)
				}
			}
		} else {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeyDefinitions, typeArrayOfSchemas)
		}

	}

	// id
	if existsMapKey(m, KeyID) && !isKind(m[KeyID], reflect.String) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, KeyID, typeString)
	}
	if k, ok := m[KeyID].(string); ok {
		currentSchema.id = &k
	}

	// title
	if existsMapKey(m, KeyTitle) && !isKind(m[KeyTitle], reflect.String) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, KeyTitle, typeString)
	}
	if k, ok := m[KeyTitle].(string); ok {
		currentSchema.title = &k
	}

	// description
	if existsMapKey(m, KeyDescription) && !isKind(m[KeyDescription], reflect.String) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, KeyDescription, typeString)
	}
	if k, ok := m[KeyDescription].(string); ok {
		currentSchema.description = &k
	}

	// type
	if existsMapKey(m, KeyType) {
		if isKind(m[KeyType], reflect.String) {
			if k, ok := m[KeyType].(string); ok {
				err := currentSchema.types.Add(k)
				if err != nil {
					return err
				}
			}
		} else {
			if isKind(m[KeyType], reflect.Slice) {
				arrayOfTypes := m[KeyType].([]interface{})
				for _, typeInArray := range arrayOfTypes {
					if reflect.ValueOf(typeInArray).Kind() != reflect.String {
						return fmt.Errorf(errMessageXMustBeOfTypeY, KeyType, typeString+"/"+typeArrayOfStrings)
					}
					currentSchema.types.Add(typeInArray.(string))
				}

			} else {
				return fmt.Errorf(errMessageXMustBeOfTypeY, KeyType, typeString+"/"+typeArrayOfStrings)
			}
		}
	}

	// properties
	if existsMapKey(m, KeyProperties) {
		err := d.parseProperties(m[KeyProperties], currentSchema)
		if err != nil {
			return err
		}
	}

	// additionalProperties
	if existsMapKey(m, KeyAdditionalProperties) {
		if isKind(m[KeyAdditionalProperties], reflect.Bool) {
			currentSchema.additionalProperties = m[KeyAdditionalProperties].(bool)
		} else if isKind(m[KeyAdditionalProperties], reflect.Map) {
			newSchema := &jsonSchema{property: KeyAdditionalProperties, parent: currentSchema, ref: currentSchema.ref}
			currentSchema.additionalProperties = newSchema
			err := d.parseSchema(m[KeyAdditionalProperties], newSchema)
			if err != nil {
				return errors.New(err.Error())
			}
		} else {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeyAdditionalProperties, typeBoolean+"/"+typeSchema)
		}
	}

	// patternProperties
	if existsMapKey(m, KeyPatternProperties) {
		if isKind(m[KeyPatternProperties], reflect.Map) {
			patternPropertiesMap := m[KeyPatternProperties].(map[string]interface{})
			if len(patternPropertiesMap) > 0 {
				currentSchema.patternProperties = make(map[string]*jsonSchema)
				for k, v := range patternPropertiesMap {
					_, err := regexp.MatchString(k, "")
					if err != nil {
						return fmt.Errorf("Invalid regex pattern '%s'", k)
					}
					newSchema := &jsonSchema{property: k, parent: currentSchema, ref: currentSchema.ref}
					err = d.parseSchema(v, newSchema)
					if err != nil {
						return errors.New(err.Error())
					}
					currentSchema.patternProperties[k] = newSchema
				}
			}
		} else {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeyPatternProperties, typeSchema)
		}
	}

	// dependencies
	if existsMapKey(m, KeyDependencies) {
		err := d.parseDependencies(m[KeyDependencies], currentSchema)
		if err != nil {
			return err
		}
	}

	// items
	if existsMapKey(m, KeyItems) {
		if isKind(m[KeyItems], reflect.Slice) {
			for _, itemElement := range m[KeyItems].([]interface{}) {
				if isKind(itemElement, reflect.Map) {
					newSchema := &jsonSchema{parent: currentSchema, property: KeyItems}
					newSchema.ref = currentSchema.ref
					currentSchema.AddItemsChild(newSchema)
					err := d.parseSchema(itemElement, newSchema)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf(errMessageXMustBeOfTypeY, KeyItems, typeSchema+"/"+typeArrayOfSchemas)
				}
				currentSchema.itemsChildrenIsSingleSchema = false
			}
		} else if isKind(m[KeyItems], reflect.Map) {
			newSchema := &jsonSchema{parent: currentSchema, property: KeyItems}
			newSchema.ref = currentSchema.ref
			currentSchema.AddItemsChild(newSchema)
			err := d.parseSchema(m[KeyItems], newSchema)
			if err != nil {
				return err
			}
			currentSchema.itemsChildrenIsSingleSchema = true
		} else {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeyItems, typeSchema+"/"+typeArrayOfSchemas)
		}
	}

	// additionalItems
	if existsMapKey(m, KeyAdditionalItems) {
		if isKind(m[KeyAdditionalItems], reflect.Bool) {
			currentSchema.additionalItems = m[KeyAdditionalItems].(bool)
		} else if isKind(m[KeyAdditionalItems], reflect.Map) {
			newSchema := &jsonSchema{property: KeyAdditionalItems, parent: currentSchema, ref: currentSchema.ref}
			currentSchema.additionalItems = newSchema
			err := d.parseSchema(m[KeyAdditionalItems], newSchema)
			if err != nil {
				return errors.New(err.Error())
			}
		} else {
			return fmt.Errorf(errMessageXMustBeOfTypeY, KeyAdditionalItems, typeBoolean+"/"+typeSchema)
		}
	}

	// validation : number / integer

	if existsMapKey(m, KeyMultipleOf) {
		multipleOfValue := mustBeNumber(m[KeyMultipleOf])
		if multipleOfValue == nil {
			return errors.New("multipleOf must be a number")
		}
		if *multipleOfValue <= 0 {
			return errors.New("multipleOf must be strictly greater than 0")
		}
		currentSchema.multipleOf = multipleOfValue
	}

	if existsMapKey(m, KeyMinimum) {
		minimumValue := mustBeNumber(m[KeyMinimum])
		if minimumValue == nil {
			return errors.New("minimum must be a number")
		}
		currentSchema.minimum = minimumValue
	}

	if existsMapKey(m, KeyExclusiveMinimum) {
		if isKind(m[KeyExclusiveMinimum], reflect.Bool) {
			if currentSchema.minimum == nil {
				return errors.New("exclusiveMinimum cannot exist without minimum")
			}
			exclusiveMinimumValue := m[KeyExclusiveMinimum].(bool)
			currentSchema.exclusiveMinimum = exclusiveMinimumValue
		} else {
			return errors.New("exclusiveMinimum must be a boolean")
		}
	}

	if existsMapKey(m, KeyMaximum) {
		maximumValue := mustBeNumber(m[KeyMaximum])
		if maximumValue == nil {
			return errors.New("maximum must be a number")
		}
		currentSchema.maximum = maximumValue
	}

	if existsMapKey(m, KeyExclusiveMaximum) {
		if isKind(m[KeyExclusiveMaximum], reflect.Bool) {
			if currentSchema.maximum == nil {
				return errors.New("exclusiveMaximum cannot exist without maximum")
			}
			exclusiveMaximumValue := m[KeyExclusiveMaximum].(bool)
			currentSchema.exclusiveMaximum = exclusiveMaximumValue
		} else {
			return errors.New("exclusiveMaximum must be a boolean")
		}
	}

	if currentSchema.minimum != nil && currentSchema.maximum != nil {
		if *currentSchema.minimum > *currentSchema.maximum {
			return errors.New("minimum cannot be greater than maximum")
		}
	}

	// validation : string

	if existsMapKey(m, KeyMinLength) {
		minLengthIntegerValue := mustBeInteger(m[KeyMinLength])
		if minLengthIntegerValue == nil {
			return errors.New("minLength must be an integer")
		}
		if *minLengthIntegerValue < 0 {
			return errors.New("minLength must be greater than or equal to 0")
		}
		currentSchema.minLength = minLengthIntegerValue
	}

	if existsMapKey(m, KeyMaxLength) {
		maxLengthIntegerValue := mustBeInteger(m[KeyMaxLength])
		if maxLengthIntegerValue == nil {
			return errors.New("maxLength must be an integer")
		}
		if *maxLengthIntegerValue < 0 {
			return errors.New("maxLength must be greater than or equal to 0")
		}
		currentSchema.maxLength = maxLengthIntegerValue
	}

	if currentSchema.minLength != nil && currentSchema.maxLength != nil {
		if *currentSchema.minLength > *currentSchema.maxLength {
			return errors.New("minLength cannot be greater than maxLength")
		}
	}

	if existsMapKey(m, KeyPattern) {
		if isKind(m[KeyPattern], reflect.String) {
			regexpObject, err := regexp.Compile(m[KeyPattern].(string))
			if err != nil {
				return errors.New("pattern must be a valid regular expression")
			}
			currentSchema.pattern = regexpObject
		} else {
			return errors.New("pattern must be a string")
		}
	}

	// validation : object

	if existsMapKey(m, KeyMinProperties) {
		minPropertiesIntegerValue := mustBeInteger(m[KeyMinProperties])
		if minPropertiesIntegerValue == nil {
			return errors.New("minProperties must be an integer")
		}
		if *minPropertiesIntegerValue < 0 {
			return errors.New("minProperties must be greater than or equal to 0")
		}
		currentSchema.minProperties = minPropertiesIntegerValue
	}

	if existsMapKey(m, KeyMaxProperties) {
		maxPropertiesIntegerValue := mustBeInteger(m[KeyMaxProperties])
		if maxPropertiesIntegerValue == nil {
			return errors.New("maxProperties must be an integer")
		}
		if *maxPropertiesIntegerValue < 0 {
			return errors.New("maxProperties must be greater than or equal to 0")
		}
		currentSchema.maxProperties = maxPropertiesIntegerValue
	}

	if currentSchema.minProperties != nil && currentSchema.maxProperties != nil {
		if *currentSchema.minProperties > *currentSchema.maxProperties {
			return errors.New("minProperties cannot be greater than maxProperties")
		}
	}

	if existsMapKey(m, KeyRequired) {
		if isKind(m[KeyRequired], reflect.Slice) {
			requiredValues := m[KeyRequired].([]interface{})
			for _, requiredValue := range requiredValues {
				if isKind(requiredValue, reflect.String) {
					err := currentSchema.AddRequired(requiredValue.(string))
					if err != nil {
						return err
					}
				} else {
					return errors.New("required items must be string")
				}
			}
		} else {
			return errors.New("required must be an array")
		}
	}

	// validation : array

	if existsMapKey(m, KeyMinItems) {
		minItemsIntegerValue := mustBeInteger(m[KeyMinItems])
		if minItemsIntegerValue == nil {
			return errors.New("minItems must be an integer")
		}
		if *minItemsIntegerValue < 0 {
			return errors.New("minItems must be greater than or equal to 0")
		}
		currentSchema.minItems = minItemsIntegerValue
	}

	if existsMapKey(m, KeyMaxItems) {
		maxItemsIntegerValue := mustBeInteger(m[KeyMaxItems])
		if maxItemsIntegerValue == nil {
			return errors.New("maxItems must be an integer")
		}
		if *maxItemsIntegerValue < 0 {
			return errors.New("maxItems must be greater than or equal to 0")
		}
		currentSchema.maxItems = maxItemsIntegerValue
	}

	if existsMapKey(m, KeyUniqueItems) {
		if isKind(m[KeyUniqueItems], reflect.Bool) {
			currentSchema.uniqueItems = m[KeyUniqueItems].(bool)
		} else {
			return errors.New("uniqueItems must be an boolean")
		}
	}

	// validation : all

	if existsMapKey(m, KeyEnum) {
		if isKind(m[KeyEnum], reflect.Slice) {
			for _, v := range m[KeyEnum].([]interface{}) {
				err := currentSchema.AddEnum(v)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("enum must be an array")
		}
	}

	// validation : schema

	if existsMapKey(m, KeyOneOf) {
		if isKind(m[KeyOneOf], reflect.Slice) {
			for _, v := range m[KeyOneOf].([]interface{}) {
				newSchema := &jsonSchema{property: KeyOneOf, parent: currentSchema, ref: currentSchema.ref}
				currentSchema.AddOneOf(newSchema)
				err := d.parseSchema(v, newSchema)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("oneOf must be an array")
		}
	}

	if existsMapKey(m, KeyAnyOf) {
		if isKind(m[KeyAnyOf], reflect.Slice) {
			for _, v := range m[KeyAnyOf].([]interface{}) {
				newSchema := &jsonSchema{property: KeyAnyOf, parent: currentSchema, ref: currentSchema.ref}
				currentSchema.AddAnyOf(newSchema)
				err := d.parseSchema(v, newSchema)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("anyOf must be an array")
		}
	}

	if existsMapKey(m, KeyAllOf) {
		if isKind(m[KeyAllOf], reflect.Slice) {
			for _, v := range m[KeyAllOf].([]interface{}) {
				newSchema := &jsonSchema{property: KeyAllOf, parent: currentSchema, ref: currentSchema.ref}
				currentSchema.AddAllOf(newSchema)
				err := d.parseSchema(v, newSchema)
				if err != nil {
					return err
				}
			}
		} else {
			return errors.New("anyOf must be an array")
		}
	}

	if existsMapKey(m, KeyNot) {
		if isKind(m[KeyNot], reflect.Map) {
			newSchema := &jsonSchema{property: KeyNot, parent: currentSchema, ref: currentSchema.ref}
			currentSchema.SetNot(newSchema)
			err := d.parseSchema(m[KeyNot], newSchema)
			if err != nil {
				return err
			}
		} else {
			return errors.New("not must be an object")
		}
	}

	return nil
}

func (d *Document) parseReference(documentNode interface{}, currentSchema *jsonSchema, reference string) (e error) {

	var err error

	jsonReference, err := jsonreference.New(reference)
	if err != nil {
		return err
	}

	standaloneDocument := d.pool.GetStandaloneDocument()

	if jsonReference.HasFullURL || standaloneDocument != nil {
		currentSchema.ref = &jsonReference
	} else {
		inheritedReference, err := currentSchema.ref.Inherits(jsonReference)
		if err != nil {
			return err
		}
		currentSchema.ref = inheritedReference
	}

	jsonPointer := currentSchema.ref.GetPointer()

	var refdDocumentNode interface{}

	if standaloneDocument != nil {

		var err error
		refdDocumentNode, _, err = jsonPointer.Get(standaloneDocument)
		if err != nil {
			return err
		}

	} else {

		var err error
		dsp, err := d.pool.GetDocument(*currentSchema.ref)
		if err != nil {
			return err
		}

		refdDocumentNode, _, err = jsonPointer.Get(dsp.Document)
		if err != nil {
			return err
		}

	}

	if !isKind(refdDocumentNode, reflect.Map) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, typeSchema, typeObject)
	}

	// returns the loaded referenced schema for the caller to update its current schema
	newSchemaDocument := refdDocumentNode.(map[string]interface{})

	newSchema := &jsonSchema{property: KeyRef, parent: currentSchema, ref: currentSchema.ref}
	d.referencePool.AddSchema(currentSchema.ref.String()+reference, newSchema)

	err = d.parseSchema(newSchemaDocument, newSchema)
	if err != nil {
		return err
	}

	currentSchema.refSchema = newSchema

	return nil

}

func (d *Document) parseProperties(documentNode interface{}, currentSchema *jsonSchema) error {

	if !isKind(documentNode, reflect.Map) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, stringProperties, typeObject)
	}

	m := documentNode.(map[string]interface{})
	for k := range m {
		schemaProperty := k
		newSchema := &jsonSchema{property: schemaProperty, parent: currentSchema, ref: currentSchema.ref}
		currentSchema.AddPropertiesChild(newSchema)
		err := d.parseSchema(m[k], newSchema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Document) parseDependencies(documentNode interface{}, currentSchema *jsonSchema) error {

	if !isKind(documentNode, reflect.Map) {
		return fmt.Errorf(errMessageXMustBeOfTypeY, KeyDependencies, typeObject)
	}

	m := documentNode.(map[string]interface{})
	currentSchema.dependencies = make(map[string]interface{})

	for k := range m {
		switch reflect.ValueOf(m[k]).Kind() {

		case reflect.Slice:
			values := m[k].([]interface{})
			var valuesToRegister []string

			for _, value := range values {
				if !isKind(value, reflect.String) {
					return fmt.Errorf(errMessageXMustBeOfTypeY, stringDependency, stringSchemaOrArrayOfStrings)
				}
				valuesToRegister = append(valuesToRegister, value.(string))
				currentSchema.dependencies[k] = valuesToRegister
			}

		case reflect.Map:
			depSchema := &jsonSchema{property: k, parent: currentSchema, ref: currentSchema.ref}
			err := d.parseSchema(m[k], depSchema)
			if err != nil {
				return err
			}
			currentSchema.dependencies[k] = depSchema

		default:
			return fmt.Errorf(errMessageXMustBeOfTypeY, stringDependency, stringSchemaOrArrayOfStrings)
		}

	}

	return nil
}
