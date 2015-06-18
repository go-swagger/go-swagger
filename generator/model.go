package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// GenerateDefinition generates a model file for a schema defintion.
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
		ExtraSchemas:   pg.ExtraSchemas,
	}, nil
}

// GenDefinition contains all the properties to generate a
// defintion from a swagger spec
type GenDefinition struct {
	GenSchema
	Package          string
	Imports          map[string]string
	DefaultImports   []string
	ExtraSchemas     []GenSchema
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
	Index              int

	GenSchema    GenSchema
	Dependencies []string
	ExtraSchemas []GenSchema
}

func (sg *schemaGenContext) NewSliceBranch(schema *spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	indexVar := pg.IndexVar
	if pg.Path == "" {
		pg.Path = "strconv.Itoa(" + indexVar + ")"
	} else {
		pg.Path = pg.Path + "+ \".\" + strconv.Itoa(" + indexVar + ")"
	}
	pg.IndexVar = indexVar + "i"
	pg.ValueExpr = pg.ValueExpr + "[" + indexVar + "]"
	pg.Schema = *schema
	pg.Required = false
	return pg
}

func (sg *schemaGenContext) NewStructBranch(name string, schema spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	if sg.Path == "" {
		pg.Path = fmt.Sprintf("%q", name)
	} else {
		pg.Path = pg.Path + "+\".\"+" + fmt.Sprintf("%q", name)
	}
	pg.Name = name
	pg.ValueExpr = pg.ValueExpr + "." + swag.ToGoName(name)
	pg.Schema = schema
	for _, fn := range sg.Schema.Required {
		if name == fn {
			pg.Required = true
			break
		}
	}
	return pg
}

func (sg *schemaGenContext) shallowClone() *schemaGenContext {
	pg := new(schemaGenContext)
	*pg = *sg
	pg.GenSchema = GenSchema{}
	pg.Dependencies = nil
	pg.ExtraSchemas = nil
	pg.Named = false
	pg.Index = 0
	return pg
}

func (sg *schemaGenContext) NewCompositionBranch(schema spec.Schema, index int) *schemaGenContext {
	pg := sg.shallowClone()
	pg.Schema = schema
	pg.Name = "AO" + strconv.Itoa(index)
	if sg.Name != sg.TypeResolver.ModelName {
		pg.Name = sg.Name + pg.Name
	}
	pg.Index = index
	return pg
}

func (sg *schemaGenContext) NewAdditionalProperty(schema spec.Schema) *schemaGenContext {
	pg := sg.shallowClone()
	pg.Schema = schema
	pg.Name = "additionalProperties"
	pg.ValueExpr = pg.ValueExpr + ".AdditionalProperties"
	return pg
}

func (sg *schemaGenContext) schemaValidations() sharedValidations {
	model := sg.Schema

	isRequired := sg.Required
	if sg.Schema.Default != nil || sg.Schema.ReadOnly {
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
	sg.ExtraSchemas = append(sg.ExtraSchemas, other.ExtraSchemas...)
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
	for i, sch := range sg.Schema.AllOf {
		comprop := sg.NewCompositionBranch(sch, i)
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
		addp := *sg.Schema.AdditionalProperties
		sg.GenSchema.HasAdditionalProperties = addp.Allows || addp.Schema != nil
		// flag swap
		if sg.GenSchema.IsComplexObject {
			sg.GenSchema.IsAdditionalProperties = sg.GenSchema.IsComplexObject
			sg.GenSchema.IsComplexObject = sg.GenSchema.IsMap
		}

		if addp.Schema != nil {
			if sg.GenSchema.IsMap || (sg.GenSchema.IsAdditionalProperties && sg.Named) {
				comprop := sg.NewAdditionalProperty(*addp.Schema)
				if err := comprop.makeGenSchema(); err != nil {
					return err
				}
				sg.MergeResult(comprop)
				sg.GenSchema.AdditionalProperties = &comprop.GenSchema
				return nil
			}
		}

		if sg.GenSchema.IsAdditionalProperties && !sg.Named {
			// for an anonoymous object, first build the new object
			// and then replace the current one with a $ref to the
			// new object
			var additionalProps schemaGenContext
			additionalProps = *sg
			additionalProps.Dependencies = nil
			additionalProps.ExtraSchemas = nil
			additionalProps.Named = true
			additionalProps.Name = swag.ToGoName(sg.GenSchema.Name + " AddedProps" + strconv.Itoa(sg.Index))
			ex := ""
			if additionalProps.Schema.Example != nil {
				ex = fmt.Sprintf("%#v", additionalProps.Schema.Example)
			}
			additionalProps.GenSchema.Example = ex
			additionalProps.GenSchema.Path = ""
			additionalProps.GenSchema.Name = swag.ToGoName(sg.GenSchema.Name)
			additionalProps.ExtraSchemas = nil
			additionalProps.Dependencies = nil
			if sg.TypeResolver.ModelName != "" {
				additionalProps.GenSchema.Name = swag.ToGoName(sg.TypeResolver.ModelName + " " + additionalProps.Name)
			}
			additionalProps.GenSchema.GoType = additionalProps.GenSchema.Name
			if sg.TypeResolver.ModelsPackage != "" {
				additionalProps.GenSchema.GoType = sg.TypeResolver.ModelsPackage + "." + additionalProps.GenSchema.Name
			}
			additionalProps.GenSchema.Title = additionalProps.GenSchema.Name + " a wrapper to serialize additional properties"
			additionalProps.GenSchema.Description = ""
			additionalProps.GenSchema.Properties = nil
			if err := (&additionalProps).buildProperties(); err != nil {
				return err
			}

			// rewrite to be a ref instead of a complex object
			sg.GenSchema.IsComplexObject = true
			sg.GenSchema.IsAnonymous = false
			sg.GenSchema.IsAdditionalProperties = false
			sg.GenSchema.HasAdditionalProperties = false
			sg.GenSchema.AdditionalProperties = nil
			sg.GenSchema.Properties = nil
			sg.GenSchema.GoType = additionalProps.GenSchema.GoType
			if addp.Schema != nil {
				comprop := additionalProps.NewAdditionalProperty(*addp.Schema)
				if err := comprop.makeGenSchema(); err != nil {
					return err
				}
				additionalProps.MergeResult(comprop)
				additionalProps.GenSchema.AdditionalProperties = &comprop.GenSchema
			}

			sg.ExtraSchemas = append(sg.ExtraSchemas, additionalProps.GenSchema)
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
		if sg.Named {
			sg.GenSchema.Name = sg.Name
			sg.GenSchema.GoType = swag.ToGoName(sg.Name)
			if sg.TypeResolver.ModelsPackage != "" {
				sg.GenSchema.GoType = sg.TypeResolver.ModelsPackage + "." + sg.GenSchema.GoType
			}
			for i, s := range sg.Schema.Items.Schemas {
				elProp := sg.NewSliceBranch(&s)
				if err := elProp.makeGenSchema(); err != nil {
					return err
				}
				sg.MergeResult(elProp)
				elProp.GenSchema.Name = "p" + strconv.Itoa(i)
				sg.GenSchema.Properties = append(sg.GenSchema.Properties, elProp.GenSchema)
			}
			return nil
		}

		// for an anonoymous object, first build the new object
		// and then replace the current one with a $ref to the
		// new tuple object
		var tup schemaGenContext
		tup = *sg
		tup.GenSchema.IsTuple = true
		tup.GenSchema.IsComplexObject = false
		tup.GenSchema.Name = swag.ToGoName(sg.GenSchema.Name + "Tuple" + strconv.Itoa(sg.Index))
		tup.Name = tup.GenSchema.Name
		if sg.TypeResolver.ModelName != "" {
			tup.GenSchema.Name = swag.ToGoName(sg.TypeResolver.ModelName + " " + tup.GenSchema.Name)
		}
		tup.GenSchema.GoType = tup.GenSchema.Name
		if sg.TypeResolver.ModelsPackage != "" {
			tup.GenSchema.GoType = sg.TypeResolver.ModelsPackage + "." + tup.GenSchema.Name
		}
		tup.GenSchema.Title = tup.GenSchema.Name + " a representation of an anonymous Tuple type"
		tup.GenSchema.Description = ""

		sg.GenSchema.IsComplexObject = true
		sg.GenSchema.IsTuple = false
		sg.GenSchema.GoType = tup.GenSchema.GoType

		for i, s := range sg.Schema.Items.Schemas {
			elProp := tup.NewSliceBranch(&s)
			if err := elProp.makeGenSchema(); err != nil {
				return err
			}
			tup.MergeResult(elProp)
			elProp.GenSchema.Name = "p" + strconv.Itoa(i)
			tup.GenSchema.Properties = append(tup.GenSchema.Properties, elProp.GenSchema)
		}
		sg.MergeResult(&tup)
		sg.ExtraSchemas = append(sg.ExtraSchemas, tup.GenSchema)
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

func (sg *schemaGenContext) shortCircuitNamedRef() (bool, error) {
	// This if block ensures that a struct gets
	// rendered with the ref as embedded ref.
	if sg.Named && sg.Schema.Ref.GetURL() != nil {
		nullableOverride := sg.GenSchema.IsNullable
		tpe := resolvedType{}
		tpe.GoType = sg.Name
		if sg.TypeResolver.ModelsPackage != "" {
			tpe.GoType = sg.TypeResolver.ModelsPackage + "." + sg.TypeResolver.ModelName
		}

		tpe.SwaggerType = "object"
		tpe.IsComplexObject = true
		tpe.IsMap = false
		tpe.IsAnonymous = false

		item := sg.NewCompositionBranch(sg.Schema, 0)
		if err := item.makeGenSchema(); err != nil {
			return true, err
		}
		sg.GenSchema.resolvedType = tpe
		sg.GenSchema.IsNullable = sg.GenSchema.IsNullable || nullableOverride
		sg.GenSchema.AllOf = append(sg.GenSchema.AllOf, item.GenSchema)
		return true, nil
	}
	return false, nil
}

func (sg *schemaGenContext) liftSpecialAllOf() error {
	// if there is only a $ref or a primitive and an x-isnullable schema then this is a nullable pointer
	if len(sg.Schema.AllOf) > 0 {
		var seenSchema int
		var seenNullable bool
		var schemaToLift spec.Schema

		for _, sch := range sg.Schema.AllOf {
			tpe, err := sg.TypeResolver.ResolveSchema(&sch, true)
			if err != nil {
				return err
			}
			if sg.TypeResolver.isNullable(&sch) {
				seenNullable = true
			}
			if len(sch.Type) > 0 || sch.Ref.GetURL() != nil {
				seenSchema++
				if (!tpe.IsAnonymous && tpe.IsComplexObject) || tpe.IsPrimitive {
					schemaToLift = sch
				}
			}
		}

		if seenSchema == 1 {
			sg.Schema = schemaToLift
			sg.GenSchema.IsNullable = seenNullable
		}
		return nil
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
	sg.GenSchema.Location = "body"
	sg.GenSchema.ValueExpression = sg.ValueExpr
	sg.GenSchema.Name = sg.Name
	sg.GenSchema.Title = sg.Schema.Title
	sg.GenSchema.Description = sg.Schema.Description
	sg.GenSchema.ReceiverName = sg.Receiver
	sg.GenSchema.sharedValidations = sg.schemaValidations()
	sg.GenSchema.ReadOnly = sg.Schema.ReadOnly

	returns, err := sg.shortCircuitNamedRef()
	if err != nil {
		return err
	}
	if returns {
		return nil
	}

	if err := sg.liftSpecialAllOf(); err != nil {
		return err
	}
	nullableOverride := sg.GenSchema.IsNullable

	tpe, err := sg.TypeResolver.ResolveSchema(&sg.Schema, !sg.Named)
	if err != nil {
		return err
	}
	tpe.IsNullable = tpe.IsNullable || nullableOverride
	sg.GenSchema.resolvedType = tpe

	if err := sg.buildProperties(); err != nil {
		return nil
	}

	if err := sg.buildAdditionalProperties(); err != nil {
		return err
	}

	if err := sg.buildAllOf(); err != nil {
		return err
	}

	if err := sg.buildXMLName(); err != nil {
		return err
	}

	if err := sg.buildAdditionalItems(); err != nil {
		return err
	}

	if err := sg.buildItems(); err != nil {
		return err
	}

	//ctx.HasSliceValidations = len(sg.GenSchema.Items) > 0 || sg.GenSchema.HasAdditionalItems || sg.GenSchema.SingleSchemaSlice
	//ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

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
	ValueExpression         string
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
