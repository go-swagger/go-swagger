package generator

import "github.com/go-swagger/go-swagger/spec"

// NewAppData creates a new application struct.
func NewAppData(name string, modelNames, operationIDs []string, opts GenOpts) error {
	// NOTE: perhaps prior to generating the actual models first
	// build their data representation
	// generate a structure that contains all the models
	// that have default values and enums
	// for each of those defined json structures we want to create a defined golang structure.
	// So we will just render all of them with fmt.Sprintf("%#v") and print that as a json document
	// to stdout, which will then be parsed and used for completing the data creation for a template

	return nil
}

// The AppData struct contains all the information converted in a format
// suitable for code generations from a swagger spec document.
// This is complemented with functions in the templates that provide
// a dynamic way to access information from the spec and transform it into
// something a template knows how to deal with
type AppData struct {
	Package             string
	ReceiverName        string
	Name                string
	Principal           string
	DefaultConsumes     string
	DefaultProduces     string
	Info                *spec.Info
	ExternalDocs        *spec.ExternalDocumentation
	Imports             map[string]string
	DefaultImports      []string
	Consumes            []SerializerGroup
	Produces            []SerializerGroup
	SecurityDefinitions []SecuritySchemeData
	Models              []GenDefinition
	OperationGroups     map[string]genOperationGroup
	SwaggerJSON         string
}

// SerializerGroup represents a group of serializer for use in a template.
type SerializerGroup struct {
	Name           string
	MediaType      string
	Implementation string
	AllSerializers []SerializerData
}

// SerializerData represents a serializer for a media type
type SerializerData struct {
	Name           string
	MediaType      string
	Implementation string
}

// SecuritySchemeData represents the data for a security scheme
type SecuritySchemeData struct {
	Name         string
	IsBasicAuth  bool
	IsAPIKeyAuth bool
	Source       string
	Principal    string
}

// ModelData the data for a named schema (from definitions property)
type ModelData struct {
	Name           string
	Title          string
	Description    string
	Properties     []ModelPropertyData
	Imports        map[string]string
	DefaultImports []string
	HasValidations bool
}

// ModelPropertyData represents a schema property.
// This can contain its own set of properties when it's an anonymous schema
type ModelPropertyData struct {
	genValidations
	resolvedType
	validationData
	Name         string
	Path         string
	ReceiverName string
	Title        string
	Description  string
	Properties   map[string]ModelPropertyData
	XMLName      string
	Enums        []ModelPropertyData
}

type swaggerValidation struct {
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool
	MaxLength        *int64
	MinLength        *int64
	Pattern          string
	MaxItems         *int64
	MinItems         *int64
	UniqueItems      bool
	MultipleOf       *float64
	Enum             []interface{}
}

type validationData struct {
	MaxLength           int64
	MinLength           int64
	Pattern             string
	MultipleOf          float64
	Minimum             float64
	Maximum             float64
	ExclusiveMinimum    bool
	ExclusiveMaximum    bool
	HasValidations      bool
	MinItems            int64
	MaxItems            int64
	UniqueItems         bool
	HasSliceValidations bool
}

type modelDataBuilder struct {
	Spec     *spec.Document
	Schema   *spec.Schema
	Resolver *typeResolver
	Path     string
	Name     string
	Receiver string
	Required bool
	Result   *ModelData
}

func (mdb modelDataBuilder) Build() error {
	return nil
}

type modelPropertyDataBuilder struct {
	Spec     *spec.Document
	Schema   *spec.Schema
	Resolver *typeResolver
	Path     string
	Name     string
	Receiver string
	Required bool
	Result   *ModelPropertyData
}

func (mdb modelPropertyDataBuilder) Build() error {
	mdb.Result = new(ModelPropertyData)
	mdb.Result.Name = mdb.Name
	mdb.Result.Required = mdb.Required
	mdb.Result.ReceiverName = mdb.Receiver
	mdb.Result.Path = mdb.Path

	tpe, err := mdb.Resolver.ResolveSchema(mdb.Schema, true)
	if err != nil {
		return err
	}
	mdb.Result.resolvedType = tpe

	if err := mdb.buildValidations(); err != nil {
		return err
	}

	if err := mdb.addDefaultValue(); err != nil {
		return err
	}

	return nil
}

func (mdb modelPropertyDataBuilder) buildValidations() error {
	adapter := swaggerValidation{
		Maximum:          mdb.Schema.Maximum,
		ExclusiveMaximum: mdb.Schema.ExclusiveMaximum,
		Minimum:          mdb.Schema.Minimum,
		ExclusiveMinimum: mdb.Schema.ExclusiveMinimum,
		MaxLength:        mdb.Schema.MaxLength,
		MinLength:        mdb.Schema.MinLength,
		Pattern:          mdb.Schema.Pattern,
		MaxItems:         mdb.Schema.MaxItems,
		MinItems:         mdb.Schema.MinItems,
		UniqueItems:      mdb.Schema.UniqueItems,
		MultipleOf:       mdb.Schema.MultipleOf,
		Enum:             mdb.Schema.Enum,
	}
	d, err := makeValidationData(adapter, mdb.Required)
	if err != nil {
		return err
	}
	mdb.Result.validationData = d
	return nil
}

func (mdb modelPropertyDataBuilder) addDefaultValue() error {
	// a default value is gotten by generating a random, predictable name
	// as a package property, and making it marshal the json
	// this uses swag.MustMarshalJSON(bytes, target)
	return nil
}

func makeValidationData(validation swaggerValidation, required bool) (validationData, error) {
	hasValidations := required
	var maxLength int64
	if validation.MaxLength != nil {
		hasValidations = true
		maxLength = *validation.MaxLength
	}

	var minLength int64
	if validation.MinLength != nil {
		hasValidations = true
		minLength = *validation.MinLength
	}

	var minimum float64
	if validation.Minimum != nil {
		hasValidations = true
		minimum = *validation.Minimum
	}

	var maximum float64
	if validation.Maximum != nil {
		hasValidations = true
		maximum = *validation.Maximum
	}

	var multipleOf float64
	if validation.MultipleOf != nil {
		hasValidations = true
		multipleOf = *validation.MultipleOf
	}

	hasSliceValidations := validation.UniqueItems
	var maxItems int64
	if validation.MaxItems != nil {
		hasSliceValidations = true
		maxItems = *validation.MaxItems
	}

	var minItems int64
	if validation.MinItems != nil {
		hasSliceValidations = true
		minItems = *validation.MinItems
	}

	return validationData{
		MaxLength:        maxLength,
		MinLength:        minLength,
		Pattern:          validation.Pattern,
		MultipleOf:       multipleOf,
		Minimum:          minimum,
		Maximum:          maximum,
		ExclusiveMinimum: validation.ExclusiveMinimum,
		ExclusiveMaximum: validation.ExclusiveMaximum,
		//Enum:                enum,
		HasValidations:      hasValidations,
		MinItems:            minItems,
		MaxItems:            maxItems,
		UniqueItems:         validation.UniqueItems,
		HasSliceValidations: hasSliceValidations,
	}, nil
}
