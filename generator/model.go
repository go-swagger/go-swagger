package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// GenerateDefinition generates a model file for a schema defintion.
//
//    defintion of primitive => type alias/name
//    defintion of array => type alias/name
//    definition of map => type alias/name
//    definition of object with properties => struct
//    definition of ref => type alias/name
//    object with only additional properties => map[string]T
//    object with additional properties and properties => custom serializer
//    schema with schema array in items => tuple (struct with properties, custom serializer)
//    schema with all of => struct
//      * all of schema with ref => embedded value
//      * all of schema with properties => properties are included in struct
//      * adding an all of schema with just "x-isnullable": true turns the schema into a pointer
func GenerateDefinition(modelNames []string, includeModel, includeValidator bool, opts GenOpts) error {
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

	// NOTE: Change this to work on a collection of models
	//       it should first build up a graph of the dependencies and
	//       when there is a circular dependency, log a message
	//       but still generate and hope for the best
	//       This is a nice to have, at worst the code will need an extra format run, if this is out of order
	//
	//       So this needs to become a 2-phased approach so that dependencies can be worked out
	for _, modelName := range modelNames {
		// lookup schema
		model, ok := specDoc.Spec().Definitions[modelName]
		if !ok {
			return fmt.Errorf("model %q not found in definitions in %s", modelName, specPath)
		}

		// generate files
		generator := definitionGenerator{
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

type definitionGenerator struct {
	Name             string
	Model            spec.Schema
	SpecDoc          *spec.Document
	Target           string
	IncludeModel     bool
	IncludeValidator bool
	Data             interface{}
	DumpData         bool
}

func (m *definitionGenerator) Generate() error {
	mod, err := makeGenDefinition(m.Name, m.Target, m.Model, m.SpecDoc)
	if err != nil {
		return err
	}
	if m.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(mod), "", " ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	mod.IncludeValidator = m.IncludeValidator
	m.Data = mod

	if m.IncludeModel {
		if err := m.generateModel(); err != nil {
			return fmt.Errorf("model: %s", err)
		}
	}
	log.Println("generated model", m.Name)

	return nil
}

func (m *definitionGenerator) generateModel() error {
	buf := bytes.NewBuffer(nil)

	if err := modelTemplate.Execute(buf, m.Data); err != nil {
		return err
	}
	log.Println("rendered model template:", m.Name)

	return writeToFile(m.Target, m.Name, buf.Bytes())
}

func makeGenDefinition(name, pkg string, schema spec.Schema, specDoc *spec.Document) (*GenDefinition, error) {
	receiver := "m"
	resolver := &typeResolver{
		ModelsPackage: "",
		ModelName:     name,
		Doc:           specDoc,
	}
	pg := schemaGenContext{
		Path:         "",
		Name:         name,
		Receiver:     receiver,
		IndexVar:     "i",
		ValueExpr:    receiver,
		Schema:       schema,
		Required:     false,
		TypeResolver: resolver,
		Named:        true,
	}
	if err := pg.makeGenSchema(); err != nil {
		return nil, err
	}

	var defaultImports []string
	if pg.GenSchema.HasValidations {
		defaultImports = []string{
			"github.com/go-swagger/go-swagger/errors",
			"github.com/go-swagger/go-swagger/strfmt",
			"github.com/go-swagger/go-swagger/httpkit/validate",
		}
	}

	return &GenDefinition{
		Package:        filepath.Base(pkg),
		GenSchema:      pg.GenSchema,
		DependsOn:      pg.Dependencies,
		DefaultImports: defaultImports,
	}, nil
}

// GenDefinition contains all the properties to generate a
// defintion from a swagger spec
type GenDefinition struct {
	GenSchema
	Package          string
	Imports          map[string]string
	DefaultImports   []string
	ExtraSchemas     []GenDefinition
	DependsOn        []string
	IncludeValidator bool
}

// GenSchemaList is a list of schemas for generation.
//
// It can be sorted by name to get a stable struct layout for
// version control and such
type GenSchemaList []GenSchema

func (g GenSchemaList) Len() int           { return len(g) }
func (g GenSchemaList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenSchemaList) Less(i, j int) bool { return g[i].Name < g[j].Name }

type schemaGenContext struct {
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

	GenSchema        GenSchema
	Dependencies     []string
	ExtraDefinitions []GenDefinition
}

func (sg *schemaGenContext) NewSliceBranch(schema *spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	indexVar := pg.IndexVar
	pg.Path = pg.Path + "." + indexVar
	pg.IndexVar = indexVar + "i"
	pg.ValueExpr = pg.ValueExpr + "[" + indexVar + "]"
	pg.Schema = *schema
	pg.Required = false
	return pg
}

func (sg *schemaGenContext) NewStructBranch(name string, schema spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	pg.Path = pg.Path + "." + name
	if sg.Path == "" {
		pg.Path = name
	}
	pg.Name = name
	pg.ValueExpr = pg.ValueExpr + "." + swag.ToGoName(name)
	pg.Schema = schema
	pg.Required = swag.ContainsStringsCI(schema.Required, name)
	return pg
}

func (sg *schemaGenContext) shallowClone() *schemaGenContext {
	pg := new(schemaGenContext)
	*pg = *sg
	pg.GenSchema = GenSchema{}
	pg.Dependencies = nil
	pg.ExtraDefinitions = nil
	pg.Named = false
	return pg
}

func (sg *schemaGenContext) NewCompositionBranch(schema spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	pg.Schema = schema
	return pg
}

func (sg *schemaGenContext) NewAdditionalProperty(schema spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	pg.Schema = schema
	pg.AdditionalProperty = true
	pg.Name = "additionalProperties"
	return pg
}

func (sg *schemaGenContext) schemaValidations() sharedValidations {
	model := sg.Schema

	isRequired := sg.Required
	if sg.Schema.Default != nil {
		isRequired = false
	}
	hasNumberValidation := model.Maximum != nil || model.Minimum != nil || model.MultipleOf != nil
	hasStringValidation := model.MaxLength != nil || model.MinLength != nil || model.Pattern != ""
	hasSliceValidations := model.MaxItems != nil || model.MinItems != nil || model.UniqueItems
	hasValidations := isRequired || hasNumberValidation || hasStringValidation || hasSliceValidations

	var enum string
	if len(sg.Schema.Enum) > 0 {
		hasValidations = true
		enum = fmt.Sprintf("%#v", model.Enum)
	}

	return sharedValidations{
		Required:            sg.Required,
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
func (sg *schemaGenContext) MergeResult(other *schemaGenContext) {
	if other.GenSchema.HasValidations {
		sg.GenSchema.HasValidations = other.GenSchema.HasValidations
	}
	sg.Dependencies = append(sg.Dependencies, other.Dependencies...)
}

func (sg *schemaGenContext) buildProperties() error {
	for k, v := range sg.Schema.Properties {
		emprop := sg.NewStructBranch(k, v)
		if err := emprop.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(emprop)
		sg.GenSchema.Properties = append(sg.GenSchema.Properties, emprop.GenSchema)
	}
	sort.Sort(sg.GenSchema.Properties)
	return nil
}

func (sg *schemaGenContext) buildAllOf() error {
	for _, sch := range sg.Schema.AllOf {
		comprop := sg.NewCompositionBranch(sch)
		if err := comprop.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(comprop)
		sg.GenSchema.AllOf = append(sg.GenSchema.AllOf, comprop.GenSchema)
	}
	return nil
}

func (sg *schemaGenContext) buildAdditionalProperties() error {
	if sg.Schema.AdditionalProperties != nil {
		sg.GenSchema.IsAdditionalProperties = sg.AdditionalProperty

		addp := sg.Schema.AdditionalProperties
		sg.GenSchema.HasAdditionalProperties = addp.Allows || addp.Schema != nil
		if addp.Schema != nil {
			comprop := sg.NewAdditionalProperty(*addp.Schema)
			if err := comprop.makeGenSchema(); err != nil {
				return err
			}
			sg.MergeResult(comprop)
			sg.GenSchema.AdditionalProperties = &comprop.GenSchema
		}
	}
	return nil
}

func (sg *schemaGenContext) buildItems() error {
	sg.GenSchema.SingleSchemaSlice = sg.Schema.Items != nil && sg.Schema.Items.Schema != nil
	if sg.GenSchema.SingleSchemaSlice {
		elProp := sg.NewSliceBranch(sg.Schema.Items.Schema)
		if err := elProp.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(elProp)
		sg.GenSchema.Items = []GenSchema{elProp.GenSchema}
	} else if sg.Schema.Items != nil {
		// This is a tuple, build a new model that represents this
		for _, s := range sg.Schema.Items.Schemas {
			elProp := sg.NewSliceBranch(&s)
			if err := elProp.makeGenSchema(); err != nil {
				return err
			}
			sg.MergeResult(elProp)
			sg.GenSchema.Items = append(sg.GenSchema.Items, elProp.GenSchema)
		}
	}
	return nil
}

func (sg *schemaGenContext) buildAdditionalItems() error {
	sg.GenSchema.AllowsAdditionalItems =
		sg.Schema.AdditionalItems != nil &&
			(sg.Schema.AdditionalItems.Allows || sg.Schema.AdditionalItems.Schema != nil)

	sg.GenSchema.HasAdditionalItems = sg.GenSchema.AllowsAdditionalItems && !sg.GenSchema.SingleSchemaSlice
	if sg.Schema.AdditionalItems != nil && sg.Schema.AdditionalItems.Schema != nil {
		it := sg.NewSliceBranch(sg.Schema.AdditionalItems.Schema)
		if err := it.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(it)
		sg.GenSchema.AdditionalItems = &it.GenSchema
	}
	return nil
}

func (sg *schemaGenContext) buildXMLName() error {
	if sg.Schema.XML != nil {
		sg.GenSchema.XMLName = sg.ParamName
		if sg.Schema.XML.Name != "" {
			sg.GenSchema.XMLName = sg.Schema.XML.Name
			if sg.Schema.XML.Attribute {
				sg.GenSchema.XMLName += ",attr"
			}
		}
	}
	return nil
}

func (sg *schemaGenContext) makeGenSchema() error {
	//log.Printf("property: (path %s) (named: %t) (name %s) (receiver %s) (indexVar %s) (expr %s) required %t", sg.Path, sg.Named, sg.Name, sg.Receiver, sg.IndexVar, sg.ValueExpr, sg.Required)
	ex := ""
	if sg.Schema.Example != nil {
		ex = fmt.Sprintf("%#v", sg.Schema.Example)
	}
	sg.GenSchema.Example = ex
	sg.GenSchema.Path = sg.Path
	sg.GenSchema.Name = sg.Name
	sg.GenSchema.Title = sg.Schema.Title
	sg.GenSchema.Description = sg.Schema.Description
	sg.GenSchema.ReceiverName = sg.Receiver

	// This if block ensures that a struct gets
	// rendered with the ref as embedded ref.
	if sg.Named && sg.Schema.Ref.GetURL() != nil {
		tpe := resolvedType{}
		tpe.GoType = sg.Name
		if sg.TypeResolver.ModelsPackage != "" {
			tpe.GoType = sg.TypeResolver.ModelsPackage + "." + sg.TypeResolver.ModelName
		}

		tpe.SwaggerType = "object"
		tpe.IsComplexObject = true
		tpe.IsMap = false
		tpe.IsAnonymous = false

		item := sg.NewCompositionBranch(sg.Schema)
		if err := item.makeGenSchema(); err != nil {
			return err
		}
		sg.GenSchema.resolvedType = tpe
		sg.GenSchema.AllOf = append(sg.GenSchema.AllOf, item.GenSchema)
		return nil
	}

	tpe, err := sg.TypeResolver.ResolveSchema(&sg.Schema, !sg.Named)
	if err != nil {
		return err
	}
	//log.Printf("%+v", tpe)

	if err := sg.buildProperties(); err != nil {
		return nil
	}

	if err := sg.buildAllOf(); err != nil {
		return err
	}

	if err := sg.buildAdditionalProperties(); err != nil {
		return err
	}

	if err := sg.buildItems(); err != nil {
		return err
	}

	if err := sg.buildAdditionalItems(); err != nil {
		return err
	}

	if err := sg.buildXMLName(); err != nil {
		return err
	}

	ctx := sg.schemaValidations()
	ctx.HasSliceValidations = len(sg.GenSchema.Items) > 0 || sg.GenSchema.HasAdditionalItems || sg.GenSchema.SingleSchemaSlice
	ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

	sg.GenSchema.resolvedType = tpe
	sg.GenSchema.sharedValidations = ctx
	sg.GenSchema.ReadOnly = sg.Schema.ReadOnly
	sg.GenSchema.ItemsLen = len(sg.GenSchema.Items)

	return nil
}

// NOTE:
// untyped data requires a cast somehow to the inner type
// I wonder if this is still a problem after adding support for tuples
// and anonymous structs. At that point there is very little that would
// end up being cast to interface, and if it does it truly is the best guess

// GenSchema contains all the information needed to generate the code
// for a schema
type GenSchema struct {
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
	Items                   []GenSchema
	ItemsLen                int
	AllowsAdditionalItems   bool
	HasAdditionalItems      bool
	AdditionalItems         *GenSchema
	Object                  *GenSchema
	XMLName                 string
	Properties              GenSchemaList
	AllOf                   []GenSchema
	HasAdditionalProperties bool
	IsAdditionalProperties  bool
	AdditionalProperties    *GenSchema
	ReadOnly                bool
}

type sharedValidations struct {
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
