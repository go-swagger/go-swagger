package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	props := make(map[string]genModelProperty)
	resolver := &typeResolver{
		ModelsPackage: "",
		Doc:           specDoc,
	}
	for pn, p := range schema.Properties {
		var required bool
		for _, v := range schema.Required {
			if v == pn {
				required = true
				break
			}
		}

		gmp, err := makeGenModelProperty(propGenBuildParams{
			Path:         pn,
			Name:         pn,
			Receiver:     receiver,
			IndexVar:     "i",
			ValueExpr:    receiver + "." + swag.ToGoName(pn),
			Schema:       p,
			Required:     required,
			TypeResolver: resolver,
		})
		if err != nil {
			return nil, err
		}
		props[pn] = gmp
	}

	var allOf []genModelProperty
	for _, p := range schema.AllOf {
		mod, err := makeGenModelProperty(propGenBuildParams{
			Name:         name,
			Path:         name + ".allOf",
			Receiver:     receiver,
			IndexVar:     "a",
			Schema:       p,
			TypeResolver: resolver,
		})
		if err != nil {
			return nil, err
		}
		allOf = append(allOf, mod)
	}

	var additionalProperties *genModelProperty
	var hasAdditionalProperties bool
	if schema.AdditionalProperties != nil {
		addp := schema.AdditionalProperties
		hasAdditionalProperties = addp.Allows || addp.Schema != nil
		if addp.Schema != nil {
			mod, err := makeGenModelProperty(propGenBuildParams{
				Name:               name,
				Path:               name + ".additionalProperties",
				Receiver:           receiver,
				IndexVar:           "p",
				Schema:             *addp.Schema,
				TypeResolver:       resolver,
				AdditionalProperty: true,
			})
			if err != nil {
				return nil, err
			}
			additionalProperties = &mod
		}
	}

	// TODO: add support for oneOf?
	// this would require a struct with unexported fields, custom json marshaller etc

	var properties []genModelProperty
	var hasValidations bool
	for _, v := range props {
		if v.HasValidations {
			hasValidations = v.HasValidations
		}
		properties = append(properties, v)
	}

	return &genModel{
		Package:        filepath.Base(pkg),
		Name:           name,
		ReceiverName:   receiver,
		Properties:     properties,
		Description:    schema.Description,
		Title:          schema.Title,
		HasValidations: hasValidations,
		AllOf:          allOf,
		HasAdditionalProperties: hasAdditionalProperties,
		AdditionalProperties:    additionalProperties,
	}, nil
}

type genModel struct {
	Package                 string
	ReceiverName            string
	Name                    string
	Path                    string
	Title                   string
	Description             string
	Properties              []genModelProperty
	Imports                 map[string]string
	DefaultImports          []string
	HasValidations          bool
	ExtraModels             []genModel
	Type                    *resolvedType
	IsAnonymous             bool // never actually set, because by definition this is named
	IsAdditionalProperties  bool // never actually set, keeps templates happy
	AllOf                   []genModelProperty
	AdditionalProperties    *genModelProperty
	HasAdditionalProperties bool
}

func modelDocString(className, desc string) string {
	return commentedLines(fmt.Sprintf("%s %s", className, desc))
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
	return pg
}

func makeGenModelProperty(params propGenBuildParams) (genModelProperty, error) {
	// log.Printf("property: (path %s) (param %s) (accessor %s) (receiver %s) (indexVar %s) (expr %s) required %t", path, paramName, accessor, receiver, indexVar, valueExpression, required)
	ex := ""
	if params.Schema.Example != nil {
		ex = fmt.Sprintf("%#v", params.Schema.Example)
	}

	ctx := modelValidations(params)
	tpe, err := params.TypeResolver.ResolveSchema(&params.Schema, true)
	if err != nil {
		return genModelProperty{}, err
	}
	var properties []genModelProperty
	for k, v := range params.Schema.Properties {
		emprop, err := makeGenModelProperty(params.NewStructBranch(k, v))
		if err != nil {
			return genModelProperty{}, err
		}
		properties = append(properties, emprop)
	}

	var allOf []genModelProperty
	for _, sch := range params.Schema.AllOf {
		comprop, err := makeGenModelProperty(params.NewCompositionBranch(sch))
		if err != nil {
			return genModelProperty{}, err
		}
		allOf = append(allOf, comprop)
	}

	var additionalProperties *genModelProperty
	var hasAdditionalProperties bool
	if params.Schema.AdditionalProperties != nil {
		addp := params.Schema.AdditionalProperties
		hasAdditionalProperties = addp.Allows || addp.Schema != nil
		if addp.Schema != nil {
			comprop, err := makeGenModelProperty(params.NewAdditionalProperty(*addp.Schema))
			if err != nil {
				return genModelProperty{}, err
			}
			additionalProperties = &comprop
		}
	}

	singleSchemaSlice := params.Schema.Items != nil && params.Schema.Items.Schema != nil
	var items []genModelProperty
	if singleSchemaSlice {
		ctx.HasSliceValidations = true

		elProp, err := makeGenModelProperty(params.NewSliceBranch(params.Schema.Items.Schema))
		if err != nil {
			return genModelProperty{}, err
		}
		items = []genModelProperty{elProp}
	} else if params.Schema.Items != nil {
		for _, s := range params.Schema.Items.Schemas {
			elProp, err := makeGenModelProperty(params.NewSliceBranch(&s))
			if err != nil {
				return genModelProperty{}, err
			}
			items = append(items, elProp)
		}
	}

	allowsAdditionalItems :=
		params.Schema.AdditionalItems != nil &&
			(params.Schema.AdditionalItems.Allows || params.Schema.AdditionalItems.Schema != nil)
	hasAdditionalItems := allowsAdditionalItems && !singleSchemaSlice
	var additionalItems *genModelProperty
	if params.Schema.AdditionalItems != nil && params.Schema.AdditionalItems.Schema != nil {
		it, err := makeGenModelProperty(params.NewSliceBranch(params.Schema.AdditionalItems.Schema))
		if err != nil {
			return genModelProperty{}, err
		}
		additionalItems = &it
	}

	ctx.HasSliceValidations = len(items) > 0 || hasAdditionalItems
	ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

	var xmlName string
	if params.Schema.XML != nil {
		xmlName = params.ParamName
		if params.Schema.XML.Name != "" {
			xmlName = params.Schema.XML.Name
			if params.Schema.XML.Attribute {
				xmlName += ",attr"
			}
		}
	}

	return genModelProperty{
		resolvedType:      tpe,
		sharedValidations: ctx,
		Example:           ex,
		Path:              params.Path,
		Name:              params.Name,
		Title:             params.Schema.Title,
		Description:       params.Schema.Description,
		ReceiverName:      params.Receiver,

		Properties: properties,
		AllOf:      allOf,
		HasAdditionalProperties: hasAdditionalProperties,
		AdditionalProperties:    additionalProperties,
		IsAdditionalProperties:  params.AdditionalProperty,

		HasAdditionalItems:    hasAdditionalItems,
		AllowsAdditionalItems: allowsAdditionalItems,
		AdditionalItems:       additionalItems,

		Items:             items,
		ItemsLen:          len(items),
		SingleSchemaSlice: singleSchemaSlice,

		XMLName: xmlName,
	}, nil
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

// the adapter
func modelValidations(params propGenBuildParams) sharedValidations {

	model := params.Schema

	var enum string
	if len(params.Schema.Enum) > 0 {
		enum = fmt.Sprintf("%#v", model.Enum)
	}

	return sharedValidations{
		Required:         params.Required,
		Maximum:          model.Maximum,
		ExclusiveMaximum: model.ExclusiveMaximum,
		Minimum:          model.Minimum,
		ExclusiveMinimum: model.ExclusiveMinimum,
		MaxLength:        model.MaxLength,
		MinLength:        model.MinLength,
		Pattern:          model.Pattern,
		MaxItems:         model.MaxItems,
		MinItems:         model.MinItems,
		UniqueItems:      model.UniqueItems,
		MultipleOf:       model.MultipleOf,
		Enum:             enum,
	}
}

func propertyDocString(propertyName, description, example string) string {
	ex := ""
	if strings.TrimSpace(example) != "" {
		ex = " eg.\n\n    " + example
	}
	return commentedLines(fmt.Sprintf("%s %s%s", propertyName, description, ex))
}
