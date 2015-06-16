package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// GenerateModel generates a model file for a schema defintion
func GenerateModel(modelNames []string, includeModel, includeValidator bool, opts GenOpts) error {
	// Load the spec
	specPath, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	if len(modelNames) == 0 {
		for k := range specDoc.Spec().Definitions {
			modelNames = append(modelNames, k)
		}
	}

	for _, modelName := range modelNames {
		// lookup schema
		model, ok := specDoc.Spec().Definitions[modelName]
		if !ok {
			return fmt.Errorf("model %q not found in definitions in %s", modelName, specPath)
		}

		// generate files
		generator := modelGenerator{
			Name:             modelName,
			Model:            model,
			SpecDoc:          specDoc,
			Target:           filepath.Join(opts.Target, opts.ModelPackage),
			IncludeModel:     includeModel,
			IncludeValidator: includeValidator,
			DumpData:         opts.DumpData,
		}

		if err := generator.Generate(); err != nil {
			return err
		}
	}

	return nil
}

type modelGenerator struct {
	Name             string
	Model            spec.Schema
	SpecDoc          *spec.Document
	Target           string
	IncludeModel     bool
	IncludeValidator bool
	Data             interface{}
	DumpData         bool
}

func (m *modelGenerator) Generate() error {
	mod, err := makeCodegenModel(m.Name, m.Target, m.Model, m.SpecDoc)
	if err != nil {
		return err
	}
	if m.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(mod), "", " ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	m.Data = mod

	if m.IncludeModel {
		if err := m.generateModel(); err != nil {
			return fmt.Errorf("model: %s", err)
		}
	}
	log.Println("generated model", m.Name)

	if m.IncludeValidator {
		if err := m.generateValidator(); err != nil {
			return fmt.Errorf("validator: %s", err)
		}
	}
	log.Println("generated validator", m.Name)
	return nil
}

func (m *modelGenerator) generateValidator() error {
	buf := bytes.NewBuffer(nil)
	if err := modelValidatorTemplate.Execute(buf, m.Data); err != nil {
		return err
	}
	log.Println("rendered validator template:", m.Name)
	return writeToFile(m.Target, m.Name+"Validator", buf.Bytes())
}

func (m *modelGenerator) generateModel() error {
	buf := bytes.NewBuffer(nil)

	if err := modelTemplate.Execute(buf, m.Data); err != nil {
		return err
	}
	log.Println("rendered model template:", m.Name)

	return writeToFile(m.Target, m.Name, buf.Bytes())
}

func makeCodegenModel(name, pkg string, schema spec.Schema, specDoc *spec.Document) (*genModel, error) {
	receiver := "m"
	resolver := &typeResolver{
		ModelsPackage: "",
		Doc:           specDoc,
	}
	pg := propGenBuildParams{
		Path:         "",
		Name:         name,
		Receiver:     receiver,
		IndexVar:     "i",
		ValueExpr:    receiver,
		Schema:       schema,
		Required:     false,
		TypeResolver: resolver,
	}
	mp, dependsOn, err := pg.makeGenModelProperty()
	if err != nil {
		return nil, err
	}

	return &genModel{
		Package:          filepath.Base(pkg),
		genModelProperty: mp,
		DependsOn:        dependsOn,
	}, nil
}

type genModel struct {
	genModelProperty
	Package        string
	Imports        map[string]string
	DefaultImports []string
	ExtraModels    []genModel
	DependsOn      []string
}

type propGenBuildParams struct {
	Path               string
	Name               string
	ParamName          string
	Accessor           string
	Receiver           string
	IndexVar           string
	ValueExpr          string
	Schema             spec.Schema
	Required           bool
	AdditionalProperty bool
	TypeResolver       *typeResolver
	Named              bool
}

func (pg propGenBuildParams) NewSliceBranch(schema *spec.Schema) propGenBuildParams {
	indexVar := pg.IndexVar
	pg.Path = pg.Path + "." + indexVar
	pg.IndexVar = indexVar + "i"
	pg.ValueExpr = pg.ValueExpr + "[" + indexVar + "]"
	pg.Schema = *schema
	pg.Required = false
	return pg
}

func (pg propGenBuildParams) NewStructBranch(name string, schema spec.Schema) propGenBuildParams {
	pg.Path = pg.Path + "." + name
	pg.Name = name
	pg.ValueExpr = pg.ValueExpr + "." + swag.ToGoName(name)
	pg.Schema = schema
	pg.Required = swag.ContainsStringsCI(schema.Required, name)
	return pg
}

func (pg propGenBuildParams) NewCompositionBranch(schema spec.Schema) propGenBuildParams {
	pg.Schema = schema
	return pg
}

func (pg propGenBuildParams) NewAdditionalProperty(schema spec.Schema) propGenBuildParams {
	pg.Schema = schema
	pg.AdditionalProperty = true
	pg.Name = "additionalProperties"
	return pg
}

func (pg propGenBuildParams) modelValidations() sharedValidations {

	model := pg.Schema

	isRequired := pg.Required
	if pg.Schema.Default != nil {
		isRequired = false
	}
	hasNumberValidation := model.Maximum != nil || model.Minimum != nil || model.MultipleOf != nil
	hasStringValidation := model.MaxLength != nil || model.MinLength != nil || model.Pattern != ""
	hasSliceValidations := model.MaxItems != nil || model.MinItems != nil || model.UniqueItems
	hasValidations := isRequired || hasNumberValidation || hasStringValidation || hasSliceValidations

	var enum string
	if len(pg.Schema.Enum) > 0 {
		hasValidations = true
		enum = fmt.Sprintf("%#v", model.Enum)
	}

	return sharedValidations{
		Required:            pg.Required,
		Maximum:             model.Maximum,
		ExclusiveMaximum:    model.ExclusiveMaximum,
		Minimum:             model.Minimum,
		ExclusiveMinimum:    model.ExclusiveMinimum,
		MaxLength:           model.MaxLength,
		MinLength:           model.MinLength,
		Pattern:             model.Pattern,
		MaxItems:            model.MaxItems,
		MinItems:            model.MinItems,
		UniqueItems:         model.UniqueItems,
		MultipleOf:          model.MultipleOf,
		Enum:                enum,
		HasValidations:      hasValidations,
		HasSliceValidations: hasSliceValidations,
	}
}

func (pg propGenBuildParams) makeGenModelProperty() (genModelProperty, []string, error) {
	// log.Printf("property: (path %s) (param %s) (accessor %s) (receiver %s) (indexVar %s) (expr %s) required %t", path, paramName, accessor, receiver, indexVar, valueExpression, required)
	ex := ""
	if pg.Schema.Example != nil {
		ex = fmt.Sprintf("%#v", pg.Schema.Example)
	}

	ctx := pg.modelValidations()
	tpe, err := pg.TypeResolver.ResolveSchema(&pg.Schema, !pg.Named)
	if err != nil {
		return genModelProperty{}, nil, err
	}

	var discovered []string
	var properties []genModelProperty
	for k, v := range pg.Schema.Properties {
		emprop, disco, err := pg.NewStructBranch(k, v).makeGenModelProperty()
		if err != nil {
			return genModelProperty{}, nil, err
		}
		if emprop.HasValidations {
			ctx.HasValidations = emprop.HasValidations
		}
		properties = append(properties, emprop)
		discovered = append(discovered, disco...)
	}

	var allOf []genModelProperty
	for _, sch := range pg.Schema.AllOf {
		comprop, disco, err := pg.NewCompositionBranch(sch).makeGenModelProperty()
		if err != nil {
			return genModelProperty{}, nil, err
		}
		if comprop.HasValidations {
			ctx.HasValidations = comprop.HasValidations
		}
		allOf = append(allOf, comprop)
		discovered = append(discovered, disco...)
	}

	var additionalProperties *genModelProperty
	var hasAdditionalProperties bool
	if pg.Schema.AdditionalProperties != nil {
		addp := pg.Schema.AdditionalProperties
		hasAdditionalProperties = addp.Allows || addp.Schema != nil
		if addp.Schema != nil {
			comprop, disco, err := pg.NewAdditionalProperty(*addp.Schema).makeGenModelProperty()
			if err != nil {
				return genModelProperty{}, nil, err
			}
			if comprop.HasValidations {
				ctx.HasValidations = comprop.HasValidations
			}
			additionalProperties = &comprop
			discovered = append(discovered, disco...)
		}
	}

	singleSchemaSlice := pg.Schema.Items != nil && pg.Schema.Items.Schema != nil
	var items []genModelProperty
	if singleSchemaSlice {
		ctx.HasSliceValidations = true

		elProp, disco, err := pg.NewSliceBranch(pg.Schema.Items.Schema).makeGenModelProperty()
		if err != nil {
			return genModelProperty{}, nil, err
		}
		items = []genModelProperty{elProp}
		discovered = append(discovered, disco...)
	} else if pg.Schema.Items != nil {
		for _, s := range pg.Schema.Items.Schemas {
			elProp, disco, err := pg.NewSliceBranch(&s).makeGenModelProperty()
			if err != nil {
				return genModelProperty{}, nil, err
			}
			items = append(items, elProp)
			discovered = append(discovered, disco...)
		}
	}

	allowsAdditionalItems :=
		pg.Schema.AdditionalItems != nil &&
			(pg.Schema.AdditionalItems.Allows || pg.Schema.AdditionalItems.Schema != nil)
	hasAdditionalItems := allowsAdditionalItems && !singleSchemaSlice
	var additionalItems *genModelProperty
	if pg.Schema.AdditionalItems != nil && pg.Schema.AdditionalItems.Schema != nil {
		it, disco, err := pg.NewSliceBranch(pg.Schema.AdditionalItems.Schema).makeGenModelProperty()
		if err != nil {
			return genModelProperty{}, nil, err
		}
		additionalItems = &it
		discovered = append(discovered, disco...)
	}

	ctx.HasSliceValidations = len(items) > 0 || hasAdditionalItems
	ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

	var xmlName string
	if pg.Schema.XML != nil {
		xmlName = pg.ParamName
		if pg.Schema.XML.Name != "" {
			xmlName = pg.Schema.XML.Name
			if pg.Schema.XML.Attribute {
				xmlName += ",attr"
			}
		}
	}

	return genModelProperty{
		resolvedType:      tpe,
		sharedValidations: ctx,
		Example:           ex,
		Path:              pg.Path,
		Name:              pg.Name,
		Title:             pg.Schema.Title,
		Description:       pg.Schema.Description,
		ReceiverName:      pg.Receiver,
		ReadOnly:          pg.Schema.ReadOnly,

		Properties: properties,
		AllOf:      allOf,
		HasAdditionalProperties: hasAdditionalProperties,
		AdditionalProperties:    additionalProperties,
		IsAdditionalProperties:  pg.AdditionalProperty,

		HasAdditionalItems:    hasAdditionalItems,
		AllowsAdditionalItems: allowsAdditionalItems,
		AdditionalItems:       additionalItems,

		Items:             items,
		ItemsLen:          len(items),
		SingleSchemaSlice: singleSchemaSlice,

		XMLName: xmlName,
	}, discovered, nil
}

// NOTE:
// untyped data requires a cast somehow to the inner type
// I wonder if this is still a problem after adding support for tuples
// and anonymous structs. At that point there is very little that would
// end up being cast to interface, and if it does it truly is the best guess

type genModelProperty struct {
	resolvedType
	sharedValidations
	Example                 string
	Name                    string
	Path                    string
	Title                   string
	Description             string
	Location                string
	ReceiverName            string
	SingleSchemaSlice       bool
	Items                   []genModelProperty
	ItemsLen                int
	AllowsAdditionalItems   bool
	HasAdditionalItems      bool
	AdditionalItems         *genModelProperty
	Object                  *genModelProperty
	XMLName                 string
	Properties              []genModelProperty
	AllOf                   []genModelProperty
	HasAdditionalProperties bool
	IsAdditionalProperties  bool
	AdditionalProperties    *genModelProperty
	ReadOnly                bool
}

type sharedValidations struct {
	Type                resolvedType
	Required            bool
	MaxLength           *int64
	MinLength           *int64
	Pattern             string
	MultipleOf          *float64
	Minimum             *float64
	Maximum             *float64
	ExclusiveMinimum    bool
	ExclusiveMaximum    bool
	Enum                string
	HasValidations      bool
	MinItems            *int64
	MaxItems            *int64
	UniqueItems         bool
	HasSliceValidations bool
	NeedsSize           bool
}
