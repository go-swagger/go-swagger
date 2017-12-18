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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
)

/*
Rewrite specification document first:

* anonymous objects
* tuples
* extensible objects (properties + additionalProperties)
* AllOfs when they match the rewrite criteria (not a nullable allOf)

Find string enums and generate specialized idiomatic enum with them

Every action that happens tracks the path which is a linked list of refs


*/

// GenerateDefinition generates a model file for a schema definition.
func GenerateDefinition(modelNames []string, opts *GenOpts) error {
	if opts == nil {
		return errors.New("gen opts are required")
	}

	if opts.TemplateDir != "" {
		if err := templates.LoadDir(opts.TemplateDir); err != nil {
			return err
		}
	}

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
			return fmt.Errorf("model %q not found in definitions given by %q", modelName, specPath)
		}

		// generate files
		generator := definitionGenerator{
			Name:    modelName,
			Model:   model,
			SpecDoc: specDoc,
			Target:  filepath.Join(opts.Target, opts.ModelPackage),
			opts:    opts,
		}

		if err := generator.Generate(); err != nil {
			return err
		}
	}

	return nil
}

type definitionGenerator struct {
	Name    string
	Model   spec.Schema
	SpecDoc *loads.Document
	Target  string
	opts    *GenOpts
}

func (m *definitionGenerator) Generate() error {

	mod, err := makeGenDefinition(m.Name, m.Target, m.Model, m.SpecDoc, m.opts)
	if err != nil {
		return fmt.Errorf("could not generate definitions for model %s on target %s: %v", m.Name, m.Target, err)
	}

	if m.opts.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(mod), "", " ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	if m.opts.IncludeModel {
		log.Println("including additional model")
		if err := m.generateModel(mod); err != nil {
			return fmt.Errorf("could not generate model: %v", err)
		}
	}
	log.Println("generated model", m.Name)

	return nil
}

func (m *definitionGenerator) generateModel(g *GenDefinition) error {
	if Debug {
		log.Printf("rendering definitions for %+v", *g)
	}
	return m.opts.renderDefinition(g)
}

func makeGenDefinition(name, pkg string, schema spec.Schema, specDoc *loads.Document, opts *GenOpts) (*GenDefinition, error) {
	return makeGenDefinitionHierarchy(name, pkg, "", schema, specDoc, opts)
}

func makeGenDefinitionHierarchy(name, pkg, container string, schema spec.Schema, specDoc *loads.Document, opts *GenOpts) (*GenDefinition, error) {

	_, ok := schema.Extensions["x-go-type"]
	if ok {
		return nil, nil
	}

	receiver := "m"
	resolver := newTypeResolver("", specDoc)
	resolver.ModelName = name
	analyzed := analysis.New(specDoc.Spec())

	di := discriminatorInfo(analyzed)

	pg := schemaGenContext{
		Path:             "",
		Name:             name,
		Receiver:         receiver,
		IndexVar:         "i",
		ValueExpr:        receiver,
		Schema:           schema,
		Required:         false,
		TypeResolver:     resolver,
		Named:            true,
		ExtraSchemas:     make(map[string]GenSchema),
		Discrimination:   di,
		Container:        container,
		IncludeValidator: opts.IncludeValidator,
		IncludeModel:     opts.IncludeModel,
	}
	if err := pg.makeGenSchema(); err != nil {
		return nil, fmt.Errorf("could not generate schema for %s: %v", name, err)
	}
	dsi, ok := di.Discriminators["#/definitions/"+name]
	if ok {
		// when these 2 are true then the schema will render as an interface
		pg.GenSchema.IsBaseType = true
		pg.GenSchema.IsExported = true
		pg.GenSchema.DiscriminatorField = dsi.FieldName

		if pg.GenSchema.Discriminates == nil {
			pg.GenSchema.Discriminates = make(map[string]string)
		}
		pg.GenSchema.Discriminates[name] = dsi.GoType
		pg.GenSchema.DiscriminatorValue = name

		for _, v := range dsi.Children {
			pg.GenSchema.Discriminates[v.FieldValue] = v.GoType
		}

		for j := range pg.GenSchema.Properties {
			pg.GenSchema.Properties[j].ValueExpression += "()"
		}
	}

	dse, ok := di.Discriminated["#/definitions/"+name]
	if ok {
		pg.GenSchema.DiscriminatorField = dse.FieldName
		pg.GenSchema.DiscriminatorValue = dse.FieldValue
		pg.GenSchema.IsSubType = true
		knownProperties := make(map[string]struct{})

		// find the referenced definitions
		// check if it has a discriminator defined
		// when it has a discriminator get the schema and run makeGenSchema for it.
		// replace the ref with this new genschema
		swsp := specDoc.Spec()
		for i, ss := range schema.AllOf {
			ref := ss.Ref
			for ref.String() != "" {
				var rsch *spec.Schema
				var err error
				rsch, err = spec.ResolveRef(swsp, &ref)
				if err != nil {
					return nil, err
				}
				ref = rsch.Ref
				if rsch != nil && rsch.Ref.String() != "" {
					ref = rsch.Ref
					continue
				}
				ref = spec.Ref{}
				if rsch != nil && rsch.Discriminator != "" {
					gs, err := makeGenDefinitionHierarchy(strings.TrimPrefix(ss.Ref.String(), "#/definitions/"), pkg, pg.GenSchema.Name, *rsch, specDoc, opts)
					if err != nil {
						return nil, err
					}
					gs.GenSchema.IsBaseType = true
					gs.GenSchema.IsExported = true
					pg.GenSchema.AllOf[i] = gs.GenSchema
					schPtr := &(pg.GenSchema.AllOf[i])
					if schPtr.AdditionalItems != nil {
						schPtr.AdditionalItems.IsBaseType = true
					}
					if schPtr.AdditionalProperties != nil {
						schPtr.AdditionalProperties.IsBaseType = true
					}
					for j := range schPtr.Properties {
						schPtr.Properties[j].IsBaseType = true
						knownProperties[schPtr.Properties[j].Name] = struct{}{}
					}
				}
			}
		}

		// dedupe the fields
		alreadySeen := make(map[string]struct{})
		for i, ss := range pg.GenSchema.AllOf {
			var remainingProperties GenSchemaList
			for _, p := range ss.Properties {
				if _, ok := knownProperties[p.Name]; !ok || ss.IsBaseType {
					if _, seen := alreadySeen[p.Name]; !seen {
						remainingProperties = append(remainingProperties, p)
						alreadySeen[p.Name] = struct{}{}
					}
				}
			}
			pg.GenSchema.AllOf[i].Properties = remainingProperties
		}

	}

	defaultImports := []string{
		"github.com/go-openapi/errors",
		"github.com/go-openapi/runtime",
		"github.com/go-openapi/swag",
		"github.com/go-openapi/validate",
	}
	var extras []GenSchema
	var extraKeys []string
	for k := range pg.ExtraSchemas {
		extraKeys = append(extraKeys, k)
	}
	sort.Strings(extraKeys)
	for _, k := range extraKeys {
		extras = append(extras, pg.ExtraSchemas[k])
	}

	return &GenDefinition{
		GenCommon: GenCommon{
			Copyright:        opts.Copyright,
			TargetImportPath: filepath.ToSlash(opts.LanguageOpts.baseImport(opts.Target)),
		},
		Package:        opts.LanguageOpts.MangleName(filepath.Base(pkg), "definitions"),
		GenSchema:      pg.GenSchema,
		DependsOn:      pg.Dependencies,
		DefaultImports: defaultImports,
		ExtraSchemas:   extras,
		Imports:        findImports(&pg.GenSchema),
	}, nil
}

func findImports(sch *GenSchema) map[string]string {
	imp := map[string]string{}
	t := sch.resolvedType
	if t.Pkg != "" && t.PkgAlias != "" {
		imp[t.PkgAlias] = t.Pkg
	}
	if sch.Items != nil {
		sub := findImports(sch.Items)
		for k, v := range sub {
			imp[k] = v
		}
	}
	if sch.AdditionalItems != nil {
		sub := findImports(sch.AdditionalItems)
		for k, v := range sub {
			imp[k] = v
		}
	}
	if sch.Object != nil {
		sub := findImports(sch.Object)
		for k, v := range sub {
			imp[k] = v
		}
	}
	if sch.Properties != nil {
		for _, p := range sch.Properties {
			sub := findImports(&p)
			for k, v := range sub {
				imp[k] = v
			}
		}
	}
	if sch.AdditionalProperties != nil {
		sub := findImports(sch.AdditionalProperties)
		for k, v := range sub {
			imp[k] = v
		}
	}
	if sch.AllOf != nil {
		for _, p := range sch.AllOf {
			sub := findImports(&p)
			for k, v := range sub {
				imp[k] = v
			}
		}
	}
	return imp
}

type schemaGenContext struct {
	Required           bool
	AdditionalProperty bool
	Untyped            bool
	Named              bool
	RefHandled         bool
	IsVirtual          bool
	IsTuple            bool
	IncludeValidator   bool
	IncludeModel       bool
	Index              int

	Path         string
	Name         string
	ParamName    string
	Accessor     string
	Receiver     string
	IndexVar     string
	KeyVar       string
	ValueExpr    string
	Container    string
	Schema       spec.Schema
	TypeResolver *typeResolver

	GenSchema      GenSchema
	Dependencies   []string
	ExtraSchemas   map[string]GenSchema
	Discriminator  *discor
	Discriminated  *discee
	Discrimination *discInfo
}

func (sg *schemaGenContext) NewSliceBranch(schema *spec.Schema) *schemaGenContext {
	if Debug {
		log.Printf("new slice branch %s (model: %s)", sg.Name, sg.TypeResolver.ModelName)
	}
	pg := sg.shallowClone()
	indexVar := pg.IndexVar
	if pg.Path == "" {
		pg.Path = "strconv.Itoa(" + indexVar + ")"
	} else {
		pg.Path = pg.Path + "+ \".\" + strconv.Itoa(" + indexVar + ")"
	}
	// check who is parent, if it's a base type then rewrite the value expression
	var rewriteValueExpr bool
	if sg.Discrimination != nil && sg.Discrimination.Discriminators != nil {
		_, rewriteValueExpr = sg.Discrimination.Discriminators["#/definitions/"+sg.TypeResolver.ModelName]
		if (pg.IndexVar == "i" && rewriteValueExpr) || sg.GenSchema.ElemType.IsBaseType {
			pg.ValueExpr = sg.Receiver + "." + swag.ToJSONName(sg.GenSchema.Name) + "Field"
		}
	}
	sg.GenSchema.IsBaseType = sg.GenSchema.ElemType.HasDiscriminator
	pg.IndexVar = indexVar + "i"
	pg.ValueExpr = pg.ValueExpr + "[" + indexVar + "]"
	pg.Schema = *schema
	pg.Required = false
	if sg.IsVirtual {
		resolver := newTypeResolver(sg.TypeResolver.ModelsPackage, sg.TypeResolver.Doc)
		resolver.ModelName = sg.TypeResolver.ModelName
		pg.TypeResolver = resolver
	}

	// when this is an anonymous complex object, this needs to become a ref
	return pg
}

func (sg *schemaGenContext) NewAdditionalItems(schema *spec.Schema) *schemaGenContext {
	if Debug {
		log.Printf("new additional items\n")
	}

	pg := sg.shallowClone()
	indexVar := pg.IndexVar
	pg.Name = sg.Name + " items"
	itemsLen := 0
	if sg.Schema.Items != nil {
		itemsLen = sg.Schema.Items.Len()
	}
	var mod string
	if itemsLen > 0 {
		mod = "+" + strconv.Itoa(itemsLen)
	}
	if pg.Path == "" {
		pg.Path = "strconv.Itoa(" + indexVar + mod + ")"
	} else {
		pg.Path = pg.Path + "+ \".\" + strconv.Itoa(" + indexVar + mod + ")"
	}
	pg.IndexVar = indexVar
	pg.ValueExpr = sg.ValueExpr + "." + sg.GoName() + "Items[" + indexVar + "]"
	pg.Schema = spec.Schema{}
	if schema != nil {
		pg.Schema = *schema
	}
	pg.Required = false
	return pg
}

func (sg *schemaGenContext) NewTupleElement(schema *spec.Schema, index int) *schemaGenContext {
	if Debug {
		log.Printf("New tuple element\n")
	}

	pg := sg.shallowClone()
	if pg.Path == "" {
		pg.Path = "\"" + strconv.Itoa(index) + "\""
	} else {
		pg.Path = pg.Path + "+ \".\"+\"" + strconv.Itoa(index) + "\""
	}
	pg.ValueExpr = pg.ValueExpr + ".P" + strconv.Itoa(index)
	pg.Required = true
	pg.IsTuple = true
	pg.Schema = *schema
	return pg
}

func (sg *schemaGenContext) NewStructBranch(name string, schema spec.Schema) *schemaGenContext {
	if Debug {
		log.Printf("new struct branch %s (parent %s)", sg.Name, sg.Container)
	}
	pg := sg.shallowClone()
	if sg.Path == "" {
		pg.Path = fmt.Sprintf("%q", name)
	} else {
		pg.Path = pg.Path + "+\".\"+" + fmt.Sprintf("%q", name)
	}
	pg.Name = name
	pg.ValueExpr = pg.ValueExpr + "." + pascalize(goName(&schema, name))
	pg.Schema = schema
	for _, fn := range sg.Schema.Required {
		if name == fn {
			pg.Required = true
			break
		}
	}
	if Debug {
		log.Printf("made new struct branch %s (parent %s)", pg.Name, pg.Container)
	}
	return pg
}

func (sg *schemaGenContext) shallowClone() *schemaGenContext {
	if Debug {
		log.Printf("cloning context %s\n", sg.Name)
	}
	pg := new(schemaGenContext)
	*pg = *sg
	if pg.Container == "" {
		pg.Container = sg.Name
	}
	pg.GenSchema = GenSchema{}
	pg.Dependencies = nil
	pg.Named = false
	pg.Index = 0
	pg.IsTuple = false
	pg.IncludeValidator = sg.IncludeValidator
	pg.IncludeModel = sg.IncludeModel
	return pg
}

func (sg *schemaGenContext) NewCompositionBranch(schema spec.Schema, index int) *schemaGenContext {
	if Debug {
		log.Printf("new composition branch %s (parent: %s, index: %d)", sg.Name, sg.Container, index)
	}
	pg := sg.shallowClone()
	pg.Schema = schema
	pg.Name = "AO" + strconv.Itoa(index)
	if sg.Name != sg.TypeResolver.ModelName {
		pg.Name = sg.Name + pg.Name
	}
	pg.Index = index
	if Debug {
		log.Printf("made new composition branch %s (parent: %s)", pg.Name, pg.Container)
	}
	return pg
}

func (sg *schemaGenContext) NewAdditionalProperty(schema spec.Schema) *schemaGenContext {
	if Debug {
		log.Printf("new additional property %s (expr: %s)", sg.Name, sg.ValueExpr)
	}
	pg := sg.shallowClone()
	pg.Schema = schema
	if pg.KeyVar == "" {
		pg.ValueExpr = sg.ValueExpr
	}
	pg.KeyVar += "k"
	pg.ValueExpr += "[" + pg.KeyVar + "]"
	pg.Path = pg.KeyVar
	pg.GenSchema.Suffix = "Value"
	if sg.Path != "" {
		pg.Path = sg.Path + "+\".\"+" + pg.KeyVar
	}
	return pg
}

func hasValidations(model *spec.Schema, isRequired bool) (needsValidation bool, hasValidation bool) {
	hasNumberValidation := model.Maximum != nil || model.Minimum != nil || model.MultipleOf != nil
	hasStringValidation := model.MaxLength != nil || model.MinLength != nil || model.Pattern != ""
	hasSliceValidations := model.MaxItems != nil || model.MinItems != nil || model.UniqueItems
	simpleObject := len(model.Properties) > 0 && model.Discriminator == ""

	needsValidation = hasNumberValidation || hasStringValidation || hasSliceValidations || len(model.Enum) > 0
	hasValidation = isRequired || needsValidation || simpleObject
	return
}

// handleFormatConflicts handles all conflicting model properties when a format is set
func handleFormatConflicts(model *spec.Schema) {
	switch model.Format {
	case "date", "datetime", "uuid", "bsonobjectid", "base64", "duration":
		model.MinLength = nil
		model.MaxLength = nil
		model.Pattern = ""
		// more cases should be inserted here if they arise
	}
}

func (sg *schemaGenContext) schemaValidations() sharedValidations {
	model := sg.Schema
	// resolve any conflicting properties if the model has a format
	handleFormatConflicts(&model)

	isRequired := sg.Required
	if model.Default != nil || model.ReadOnly {
		isRequired = false
	}
	hasSliceValidations := model.MaxItems != nil || model.MinItems != nil || model.UniqueItems
	needsValidation, hasValidation := hasValidations(&model, isRequired)

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
		Enum:                model.Enum,
		HasValidations:      hasValidation,
		HasSliceValidations: hasSliceValidations,
		NeedsValidation:     needsValidation,
		NeedsRequired:       isRequired,
	}
}
func (sg *schemaGenContext) MergeResult(other *schemaGenContext, liftsRequired bool) {
	if other.GenSchema.AdditionalProperties != nil && other.GenSchema.AdditionalProperties.HasValidations {
		sg.GenSchema.HasValidations = true
	}
	if other.GenSchema.AdditionalProperties != nil && other.GenSchema.AdditionalProperties.NeedsValidation {
		sg.GenSchema.NeedsValidation = true
	}
	if liftsRequired && other.GenSchema.AdditionalProperties != nil && other.GenSchema.AdditionalProperties.NeedsRequired {
		sg.GenSchema.NeedsRequired = true
	}
	if liftsRequired && other.GenSchema.AdditionalProperties != nil && other.GenSchema.AdditionalProperties.Required {
		sg.GenSchema.Required = true
	}
	if other.GenSchema.HasValidations {
		sg.GenSchema.HasValidations = other.GenSchema.HasValidations
	}
	if other.GenSchema.NeedsValidation {
		sg.GenSchema.NeedsValidation = other.GenSchema.NeedsValidation
	}
	if liftsRequired && other.GenSchema.NeedsRequired {
		sg.GenSchema.NeedsRequired = other.GenSchema.NeedsRequired
	}
	if liftsRequired && other.GenSchema.Required {
		sg.GenSchema.Required = other.GenSchema.Required
	}
	if other.GenSchema.HasBaseType {
		sg.GenSchema.HasBaseType = other.GenSchema.HasBaseType
	}
	sg.Dependencies = append(sg.Dependencies, other.Dependencies...)
	for k, v := range other.ExtraSchemas {
		sg.ExtraSchemas[k] = v
	}
}

func (sg *schemaGenContext) buildProperties() error {

	if Debug {
		log.Printf("building properties %s (parent: %s)", sg.Name, sg.Container)
	}

	for k, v := range sg.Schema.Properties {
		if Debug {
			bbb, _ := json.MarshalIndent(sg.Schema, "", "  ")
			log.Printf("building property %s[%q] (tup: %t) (BaseType: %t) %s\n", sg.Name, k, sg.IsTuple, sg.GenSchema.IsBaseType, bbb)
		}

		// check if this requires de-anonymizing, if so lift this as a new struct and extra schema
		tpe, err := sg.TypeResolver.ResolveSchema(&v, true, sg.IsTuple || containsString(sg.Schema.Required, k))
		if sg.Schema.Discriminator == k {
			tpe.IsNullable = false
		}
		if err != nil {
			return err
		}

		vv := v
		var hasValidation bool
		var needsValidation bool
		if tpe.IsComplexObject && tpe.IsAnonymous && len(v.Properties) > 0 {
			pg := sg.makeNewStruct(sg.Name+swag.ToGoName(k), v)
			pg.IsTuple = sg.IsTuple
			if sg.Path != "" {
				pg.Path = sg.Path + "+ \".\"+" + fmt.Sprintf("%q", k)
			} else {
				pg.Path = fmt.Sprintf("%q", k)
			}
			if err := pg.makeGenSchema(); err != nil {
				return err
			}
			if v.Discriminator != "" {
				pg.GenSchema.IsBaseType = true
				pg.GenSchema.IsExported = true
				pg.GenSchema.HasBaseType = true
			}

			vv = *spec.RefProperty("#/definitions/" + pg.Name)
			hasValidation = pg.GenSchema.HasValidations
			needsValidation = pg.GenSchema.NeedsValidation
			sg.MergeResult(pg, false)
			sg.ExtraSchemas[pg.Name] = pg.GenSchema
		}

		emprop := sg.NewStructBranch(k, vv)
		emprop.IsTuple = sg.IsTuple
		if err := emprop.makeGenSchema(); err != nil {
			return err
		}
		if hasValidation || emprop.GenSchema.HasValidations {
			emprop.GenSchema.HasValidations = true
			sg.GenSchema.HasValidations = true
		}
		if needsValidation || emprop.GenSchema.NeedsValidation {
			emprop.GenSchema.NeedsValidation = true
			sg.GenSchema.NeedsValidation = true
		}
		if emprop.Schema.Ref.String() != "" {
			ref := emprop.Schema.Ref
			var sch *spec.Schema
			for ref.String() != "" {
				var rsch *spec.Schema
				var err error
				specDoc := sg.TypeResolver.Doc
				rsch, err = spec.ResolveRef(specDoc.Spec(), &ref)
				if err != nil {
					return err
				}
				ref = rsch.Ref
				if rsch != nil && rsch.Ref.String() != "" {
					ref = rsch.Ref
					continue
				}
				ref = spec.Ref{}
				sch = rsch
			}
			if emprop.Discrimination != nil {
				if _, ok := emprop.Discrimination.Discriminators[emprop.Schema.Ref.String()]; ok {
					emprop.GenSchema.IsBaseType = true
					emprop.GenSchema.IsNullable = false
					emprop.GenSchema.HasBaseType = true
				}
				if _, ok := emprop.Discrimination.Discriminated[emprop.Schema.Ref.String()]; ok {
					emprop.GenSchema.IsSubType = true
				}
			}
			var nm = filepath.Base(emprop.Schema.Ref.GetURL().Fragment)
			var tn string
			if gn, ok := emprop.Schema.Extensions["x-go-name"]; ok {
				tn = gn.(string)
			} else {
				tn = swag.ToGoName(nm)
			}

			tr := newTypeResolver(sg.TypeResolver.ModelsPackage, sg.TypeResolver.Doc)
			tr.ModelName = tn
			ttpe, err := tr.ResolveSchema(sch, false, true)
			if err != nil {
				return err
			}
			if ttpe.IsAliased {
				emprop.GenSchema.IsAliased = true
			}
			nv, hv := hasValidations(sch, false)
			if hv {
				emprop.GenSchema.HasValidations = true
			}
			if nv {
				emprop.GenSchema.NeedsValidation = true
			}
		}
		if sg.Schema.Discriminator == k {
			emprop.GenSchema.IsNullable = false
		}
		if emprop.GenSchema.IsBaseType {
			sg.GenSchema.HasBaseType = true
		}
		sg.MergeResult(emprop, false)

		if customTag, found := emprop.Schema.Extensions["x-go-custom-tag"]; found {
			emprop.GenSchema.CustomTag = customTag.(string)
		}
		if emprop.GenSchema.HasDiscriminator {
			emprop.GenSchema.ValueExpression += "()"
		}
		sg.GenSchema.Properties = append(sg.GenSchema.Properties, emprop.GenSchema)
	}
	sort.Sort(sg.GenSchema.Properties)
	return nil
}

func (sg *schemaGenContext) buildAllOf() error {
	if len(sg.Schema.AllOf) > 0 {
		if sg.Container == "" {
			sg.Container = sg.Name
		}
	}
	if Debug {
		log.Printf("building all of for %d entries", len(sg.Schema.AllOf))
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}
	for i, sch := range sg.Schema.AllOf {
		if Debug {
			b, _ := json.MarshalIndent(sch, "", "  ")
			log.Println("trying", string(b))
		}
		comprop := sg.NewCompositionBranch(sch, i)
		if err := comprop.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(comprop, true)
		sg.GenSchema.AllOf = append(sg.GenSchema.AllOf, comprop.GenSchema)
	}
	if len(sg.Schema.AllOf) > 0 {
		sg.GenSchema.IsNullable = true
	}
	return nil
}

type mapStack struct {
	Type     *spec.Schema
	Next     *mapStack
	Previous *mapStack
	ValueRef *schemaGenContext
	Context  *schemaGenContext
	NewObj   *schemaGenContext
}

func newMapStack(context *schemaGenContext) (first, last *mapStack, err error) {
	ms := &mapStack{
		Type:    &context.Schema,
		Context: context,
	}

	l := ms
	for l.HasMore() {
		tpe, err := l.Context.TypeResolver.ResolveSchema(l.Type.AdditionalProperties.Schema, true, true)
		if err != nil {
			return nil, nil, err
		}
		if !tpe.IsMap {
			if tpe.IsComplexObject && tpe.IsAnonymous {
				nw := l.Context.makeNewStruct(l.Context.Name+" Anon", *l.Type.AdditionalProperties.Schema)
				sch := spec.RefProperty("#/definitions/" + nw.Name)
				l.NewObj = nw
				l.Type.AdditionalProperties.Schema = sch
				l.ValueRef = l.Context.NewAdditionalProperty(*sch)
			}
			break
		}
		l.Next = &mapStack{
			Previous: l,
			Type:     l.Type.AdditionalProperties.Schema,
			Context:  l.Context.NewAdditionalProperty(*l.Type.AdditionalProperties.Schema),
		}
		l = l.Next
	}

	return ms, l, nil
}

func (mt *mapStack) Build() error {
	if mt.NewObj == nil && mt.ValueRef == nil && mt.Next == nil && mt.Previous == nil {
		csch := mt.Type.AdditionalProperties.Schema
		cp := mt.Context.NewAdditionalProperty(*csch)
		d := mt.Context.TypeResolver.Doc
		asch, err := analysis.Schema(analysis.SchemaOpts{
			Root:     d.Spec(),
			BasePath: d.SpecFilePath(),
			Schema:   csch,
		})
		if err != nil {
			return err
		}
		cp.Required = !asch.IsSimpleSchema && !asch.IsMap
		if err := cp.makeGenSchema(); err != nil {
			return err
		}
		mt.Context.MergeResult(cp, false)
		mt.Context.GenSchema.AdditionalProperties = &cp.GenSchema
		if Debug {
			log.Println("early mapstack exit, nullable:", cp.GenSchema.IsNullable)
		}
		return nil
	}
	cur := mt
	for cur != nil {
		if cur.NewObj != nil {
			if err := cur.NewObj.makeGenSchema(); err != nil {
				return err
			}
		}

		if cur.ValueRef != nil {
			if err := cur.ValueRef.makeGenSchema(); err != nil {
				return nil
			}
		}

		if cur.NewObj != nil {
			cur.Context.MergeResult(cur.NewObj, false)
			cur.Context.ExtraSchemas[cur.NewObj.Name] = cur.NewObj.GenSchema
		}

		if cur.ValueRef != nil {
			if err := cur.Context.makeGenSchema(); err != nil {
				return err
			}
			cur.ValueRef.GenSchema.HasValidations = cur.NewObj.GenSchema.HasValidations
			cur.ValueRef.GenSchema.NeedsValidation = cur.NewObj.GenSchema.NeedsValidation
			cur.Context.MergeResult(cur.ValueRef, false)
			cur.Context.GenSchema.AdditionalProperties = &cur.ValueRef.GenSchema
		}

		if cur.Previous != nil {
			if err := cur.Context.makeGenSchema(); err != nil {
				return err
			}
		}
		if cur.Next != nil {
			cur.Context.MergeResult(cur.Next.Context, false)
			cur.Context.GenSchema.AdditionalProperties = &cur.Next.Context.GenSchema
		}
		if cur.ValueRef != nil {
			cur.Context.MergeResult(cur.ValueRef, false)
			cur.Context.GenSchema.AdditionalProperties = &cur.ValueRef.GenSchema
		}
		cur = cur.Previous
	}

	return nil
}

func (mt *mapStack) HasMore() bool {
	return mt.Type.AdditionalProperties != nil && (mt.Type.AdditionalProperties.Schema != nil || mt.Type.AdditionalProperties.Allows)
}

func (mt *mapStack) Dict() map[string]interface{} {
	res := make(map[string]interface{})
	res["context"] = mt.Context.Schema
	if mt.Next != nil {
		res["next"] = mt.Next.Dict()
	}
	if mt.NewObj != nil {
		res["obj"] = mt.NewObj.Schema
	}
	if mt.ValueRef != nil {
		res["value"] = mt.ValueRef.Schema
	}
	return res
}

func (sg *schemaGenContext) buildAdditionalProperties() error {
	if sg.Schema.AdditionalProperties == nil {
		return nil
	}
	addp := *sg.Schema.AdditionalProperties
	wantsAdditional := addp.Schema != nil || addp.Allows
	sg.GenSchema.HasAdditionalProperties = wantsAdditional
	if !wantsAdditional {
		return nil
	}
	// flag swap
	if sg.GenSchema.IsComplexObject {
		sg.GenSchema.IsAdditionalProperties = true
		sg.GenSchema.IsComplexObject = false
		sg.GenSchema.IsMap = false
	}

	if addp.Schema == nil {
		if addp.Allows { // map with interface{} value
			addp.Schema = &spec.Schema{}
			addp.Schema.Typed("object", "")
			sg.GenSchema.HasAdditionalProperties = true
			sg.GenSchema.IsComplexObject = false
			sg.GenSchema.IsMap = true

			sg.GenSchema.ValueExpression += "." + swag.ToGoName(sg.Name+" additionalProperties")
			cp := sg.NewAdditionalProperty(*addp.Schema)
			cp.Name += "AdditionalProperties"
			cp.Required = false
			if err := cp.makeGenSchema(); err != nil {
				return err
			}
			sg.MergeResult(cp, false)
			sg.GenSchema.AdditionalProperties = &cp.GenSchema
			if Debug {
				log.Println("added interface{} schema for additionalProperties[allows == true]", cp.GenSchema.IsInterface)
			}
		}
		return nil
	}

	if !sg.GenSchema.IsMap && (sg.GenSchema.IsAdditionalProperties && sg.Named) {
		sg.GenSchema.ValueExpression += "." + swag.ToGoName(sg.GenSchema.Name)
		comprop := sg.NewAdditionalProperty(*addp.Schema)
		d := sg.TypeResolver.Doc
		asch, err := analysis.Schema(analysis.SchemaOpts{
			Root:     d.Spec(),
			BasePath: d.SpecFilePath(),
			Schema:   addp.Schema,
		})
		if err != nil {
			return err
		}
		comprop.Required = !asch.IsSimpleSchema && !asch.IsMap
		if err := comprop.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(comprop, false)
		sg.GenSchema.AdditionalProperties = &comprop.GenSchema
		return nil
	}

	if sg.GenSchema.IsMap && wantsAdditional {
		// find out how deep this rabbit hole goes
		// descend, unwind and rewrite
		// This needs to be depth first, so it first goes as deep as it can and then
		// builds the result in reverse order.
		_, ls, err := newMapStack(sg)
		if err != nil {
			return err
		}
		if err := ls.Build(); err != nil {
			return err
		}

		return nil
	}

	if sg.GenSchema.IsAdditionalProperties && !sg.Named {
		// for an anonoymous object, first build the new object
		// and then replace the current one with a $ref to the
		// new object
		newObj := sg.makeNewStruct(sg.GenSchema.Name+" P"+strconv.Itoa(sg.Index), sg.Schema)
		if err := newObj.makeGenSchema(); err != nil {
			return err
		}

		sg.GenSchema = GenSchema{}
		sg.Schema = *spec.RefProperty("#/definitions/" + newObj.Name)
		if err := sg.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(newObj, false)
		sg.GenSchema.HasValidations = newObj.GenSchema.HasValidations
		sg.ExtraSchemas[newObj.Name] = newObj.GenSchema
		return nil
	}
	return nil
}

func (sg *schemaGenContext) makeNewStruct(name string, schema spec.Schema) *schemaGenContext {
	if Debug {
		log.Println("making new struct", name, sg.Container)
	}
	sp := sg.TypeResolver.Doc.Spec()
	name = swag.ToGoName(name)
	if sg.TypeResolver.ModelName != sg.Name {
		name = swag.ToGoName(sg.TypeResolver.ModelName + " " + name)
	}
	if sp.Definitions == nil {
		sp.Definitions = make(spec.Definitions)
	}
	sp.Definitions[name] = schema
	pg := schemaGenContext{
		Path:             "",
		Name:             name,
		Receiver:         sg.Receiver,
		IndexVar:         "i",
		ValueExpr:        sg.Receiver,
		Schema:           schema,
		Required:         false,
		Named:            true,
		ExtraSchemas:     make(map[string]GenSchema),
		Discrimination:   sg.Discrimination,
		Container:        sg.Container,
		IncludeValidator: sg.IncludeValidator,
		IncludeModel:     sg.IncludeModel,
	}
	if schema.Ref.String() == "" {
		resolver := newTypeResolver(sg.TypeResolver.ModelsPackage, sg.TypeResolver.Doc)
		resolver.ModelName = name //sg.TypeResolver.ModelName
		pg.TypeResolver = resolver
	}
	pg.GenSchema.IsVirtual = true

	sg.ExtraSchemas[name] = pg.GenSchema
	return &pg
}

func (sg *schemaGenContext) buildArray() error {
	tpe, err := sg.TypeResolver.ResolveSchema(sg.Schema.Items.Schema, true, false)
	if err != nil {
		return err
	}
	// check if the element is a complex object, if so generate a new type for it
	if tpe.IsComplexObject && tpe.IsAnonymous {
		pg := sg.makeNewStruct(sg.Name+" items"+strconv.Itoa(sg.Index), *sg.Schema.Items.Schema)
		if err := pg.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(pg, false)
		sg.ExtraSchemas[pg.Name] = pg.GenSchema
		sg.Schema.Items.Schema = spec.RefProperty("#/definitions/" + pg.Name)
		sg.IsVirtual = true
		return sg.makeGenSchema()
	}

	elProp := sg.NewSliceBranch(sg.Schema.Items.Schema)
	elProp.Required = true
	if err := elProp.makeGenSchema(); err != nil {
		return err
	}

	sg.MergeResult(elProp, false)
	sg.GenSchema.IsBaseType = elProp.GenSchema.IsBaseType
	sg.GenSchema.ItemsEnum = elProp.GenSchema.Enum
	elProp.GenSchema.Suffix = "Items"
	sg.GenSchema.GoType = "[]" + elProp.GenSchema.GoType

	elProp.GenSchema.IsNullable = tpe.IsNullable && !tpe.HasDiscriminator
	if elProp.GenSchema.IsNullable {
		sg.GenSchema.GoType = "[]*" + elProp.GenSchema.GoType
	}
	sg.GenSchema.IsArray = true
	sg.GenSchema.IsExported = !sg.GenSchema.ElemType.HasDiscriminator

	schemaCopy := elProp.GenSchema
	schemaCopy.Required = false
	hv, _ := hasValidations(sg.Schema.Items.Schema, false)
	schemaCopy.HasValidations = elProp.GenSchema.IsNullable || hv
	sg.GenSchema.Items = &schemaCopy
	if sg.Named {
		sg.GenSchema.AliasedType = sg.GenSchema.GoType
	}
	return nil
}

func (sg *schemaGenContext) buildItems() error {
	presentsAsSingle := sg.Schema.Items != nil && sg.Schema.Items.Schema != nil
	if presentsAsSingle && sg.Schema.AdditionalItems != nil { // unsure if htis a valid of invalid schema
		return fmt.Errorf("single schema (%s) can't have additional items", sg.Name)
	}
	if presentsAsSingle {
		return sg.buildArray()
	}
	if sg.Schema.Items == nil {
		return nil
	}
	// This is a tuple, build a new model that represents this
	if sg.Named {
		sg.GenSchema.Name = sg.Name
		sg.GenSchema.GoType = sg.TypeResolver.goTypeName(sg.Name)
		for i, s := range sg.Schema.Items.Schemas {
			elProp := sg.NewTupleElement(&s, i)
			if err := elProp.makeGenSchema(); err != nil {
				return err
			}
			sg.MergeResult(elProp, false)
			elProp.GenSchema.Name = "p" + strconv.Itoa(i)
			sg.GenSchema.Properties = append(sg.GenSchema.Properties, elProp.GenSchema)
		}
		return nil
	}

	// for an anonoymous object, first build the new object
	// and then replace the current one with a $ref to the
	// new tuple object
	var sch spec.Schema
	sch.Typed("object", "")
	sch.Properties = make(map[string]spec.Schema)
	for i, v := range sg.Schema.Items.Schemas {
		sch.Required = append(sch.Required, "P"+strconv.Itoa(i))
		sch.Properties["P"+strconv.Itoa(i)] = v
	}
	sch.AdditionalItems = sg.Schema.AdditionalItems
	tup := sg.makeNewStruct(sg.GenSchema.Name+"Tuple"+strconv.Itoa(sg.Index), sch)
	tup.IsTuple = true
	if err := tup.makeGenSchema(); err != nil {
		return err
	}
	tup.GenSchema.IsTuple = true
	tup.GenSchema.IsComplexObject = false
	tup.GenSchema.Title = tup.GenSchema.Name + " a representation of an anonymous Tuple type"
	tup.GenSchema.Description = ""
	sg.ExtraSchemas[tup.Name] = tup.GenSchema

	sg.Schema = *spec.RefProperty("#/definitions/" + tup.Name)
	if err := sg.makeGenSchema(); err != nil {
		return err
	}
	sg.MergeResult(tup, false)
	return nil
}

func (sg *schemaGenContext) buildAdditionalItems() error {
	wantsAdditionalItems :=
		sg.Schema.AdditionalItems != nil &&
			(sg.Schema.AdditionalItems.Allows || sg.Schema.AdditionalItems.Schema != nil)

	sg.GenSchema.HasAdditionalItems = wantsAdditionalItems
	if wantsAdditionalItems {
		// check if the element is a complex object, if so generate a new type for it
		tpe, err := sg.TypeResolver.ResolveSchema(sg.Schema.AdditionalItems.Schema, true, true)
		if err != nil {
			return err
		}
		if tpe.IsComplexObject && tpe.IsAnonymous {
			pg := sg.makeNewStruct(sg.Name+" Items", *sg.Schema.AdditionalItems.Schema)
			if err := pg.makeGenSchema(); err != nil {
				return err
			}
			sg.Schema.AdditionalItems.Schema = spec.RefProperty("#/definitions/" + pg.Name)
			pg.GenSchema.HasValidations = true
			sg.MergeResult(pg, false)
			sg.ExtraSchemas[pg.Name] = pg.GenSchema
		}

		it := sg.NewAdditionalItems(sg.Schema.AdditionalItems.Schema)
		if tpe.IsInterface {
			it.Untyped = true
		}

		if err := it.makeGenSchema(); err != nil {
			return err
		}
		sg.MergeResult(it, true)
		sg.GenSchema.AdditionalItems = &it.GenSchema
	}
	return nil
}

func (sg *schemaGenContext) buildXMLName() error {
	if sg.Schema.XML == nil {
		return nil
	}
	sg.GenSchema.XMLName = sg.Name

	if sg.Schema.XML.Name != "" {
		sg.GenSchema.XMLName = sg.Schema.XML.Name
		if sg.Schema.XML.Attribute {
			sg.GenSchema.XMLName += ",attr"
		}
	}
	return nil
}

func (sg *schemaGenContext) shortCircuitNamedRef() (bool, error) {
	// This if block ensures that a struct gets
	// rendered with the ref as embedded ref.
	if sg.RefHandled || !sg.Named || sg.Schema.Ref.String() == "" {
		return false, nil
	}
	if Debug {
		log.Printf("short circuit named ref: %q", sg.Schema.Ref.String())
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}

	nullableOverride := sg.GenSchema.IsNullable
	tpe := resolvedType{}
	tpe.GoType = sg.TypeResolver.goTypeName(sg.Name)

	tpe.SwaggerType = "object"
	tpe.IsComplexObject = true
	tpe.IsMap = false
	tpe.IsAnonymous = false
	tpe.IsNullable = sg.TypeResolver.IsNullable(&sg.Schema)

	item := sg.NewCompositionBranch(sg.Schema, 0)
	if err := item.makeGenSchema(); err != nil {
		return true, err
	}
	sg.GenSchema.resolvedType = tpe
	sg.GenSchema.IsNullable = sg.GenSchema.IsNullable || nullableOverride
	sg.MergeResult(item, true)
	sg.GenSchema.AllOf = append(sg.GenSchema.AllOf, item.GenSchema)
	return true, nil
}

func (sg *schemaGenContext) liftSpecialAllOf() error {
	// if there is only a $ref or a primitive and an x-isnullable schema then this is a nullable pointer
	// so this should not compose several objects, just 1
	// if there is a ref with a discriminator then we look for x-class on the current definition to know
	// the value of the discriminator to instantiate the class
	if len(sg.Schema.AllOf) < 2 {
		return nil
	}
	var seenSchema int
	var seenNullable bool
	var schemaToLift spec.Schema

	for _, sch := range sg.Schema.AllOf {

		tpe, err := sg.TypeResolver.ResolveSchema(&sch, true, true)
		if err != nil {
			return err
		}
		if sg.TypeResolver.IsNullable(&sch) {
			seenNullable = true
		}
		if len(sch.Type) > 0 || len(sch.Properties) > 0 || sch.Ref.GetURL() != nil {
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

func (sg *schemaGenContext) buildAliased() error {
	if !sg.GenSchema.IsPrimitive && !sg.GenSchema.IsMap && !sg.GenSchema.IsArray && !sg.GenSchema.IsInterface {
		return nil
	}

	if sg.GenSchema.IsPrimitive {
		if sg.GenSchema.SwaggerType == "string" && sg.GenSchema.SwaggerFormat == "" {
			sg.GenSchema.IsAliased = sg.GenSchema.GoType != sg.GenSchema.SwaggerType
		}
		if sg.GenSchema.IsNullable && sg.Named {
			sg.GenSchema.IsNullable = false
		}
	}

	if sg.GenSchema.IsInterface {
		sg.GenSchema.IsAliased = sg.GenSchema.GoType != iface
	}
	if sg.GenSchema.IsMap {
		sg.GenSchema.IsAliased = !strings.HasPrefix(sg.GenSchema.GoType, "map[")
	}
	if sg.GenSchema.IsArray {
		sg.GenSchema.IsAliased = !strings.HasPrefix(sg.GenSchema.GoType, "[]")
	}
	return nil
}

func (sg *schemaGenContext) GoName() string {
	return goName(&sg.Schema, sg.Name)
}

func goName(sch *spec.Schema, orig string) string {
	name, _ := sch.Extensions.GetString("x-go-name")
	if name != "" {
		return name
	}
	return orig
}

func (sg *schemaGenContext) makeGenSchema() error {
	if Debug {
		log.Printf("making gen schema (anon: %t, req: %t, tuple: %t) %s\n", !sg.Named, sg.Required, sg.IsTuple, sg.Name)
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}

	ex := ""
	if sg.Schema.Example != nil {
		ex = fmt.Sprintf("%#v", sg.Schema.Example)
	}
	sg.GenSchema.IsExported = true
	sg.GenSchema.Example = ex
	sg.GenSchema.Path = sg.Path
	sg.GenSchema.IndexVar = sg.IndexVar
	sg.GenSchema.Location = body
	sg.GenSchema.ValueExpression = sg.ValueExpr
	sg.GenSchema.KeyVar = sg.KeyVar
	sg.GenSchema.OriginalName = sg.Name
	sg.GenSchema.Name = sg.GoName()
	sg.GenSchema.Title = sg.Schema.Title
	sg.GenSchema.Description = trimBOM(sg.Schema.Description)
	sg.GenSchema.ReceiverName = sg.Receiver
	sg.GenSchema.sharedValidations = sg.schemaValidations()
	sg.GenSchema.ReadOnly = sg.Schema.ReadOnly
	sg.GenSchema.IncludeValidator = sg.IncludeValidator
	sg.GenSchema.IncludeModel = sg.IncludeModel
	sg.GenSchema.Default = sg.Schema.Default

	var err error
	returns, err := sg.shortCircuitNamedRef()
	if err != nil {
		return err
	}
	if returns {
		return nil
	}
	if Debug {
		log.Println("after shortcuit named ref")
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}

	if e := sg.liftSpecialAllOf(); e != nil {
		return e
	}
	nullableOverride := sg.GenSchema.IsNullable
	if Debug {
		log.Println("after lifting special all of")
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}

	if sg.Container == "" {
		sg.Container = sg.GenSchema.Name
	}
	if e := sg.buildAllOf(); e != nil {
		return e
	}

	var tpe resolvedType
	if sg.Untyped {
		tpe, err = sg.TypeResolver.ResolveSchema(nil, !sg.Named, sg.IsTuple || sg.Required || sg.GenSchema.Required)
	} else {
		tpe, err = sg.TypeResolver.ResolveSchema(&sg.Schema, !sg.Named, sg.IsTuple || sg.Required || sg.GenSchema.Required)
	}
	if err != nil {
		return err
	}
	if Debug {
		log.Println("gschema rrequired", sg.GenSchema.Required, "nullable", sg.GenSchema.IsNullable)
	}
	tpe.IsNullable = tpe.IsNullable || nullableOverride
	sg.GenSchema.resolvedType = tpe
	sg.GenSchema.IsBaseType = tpe.IsBaseType
	sg.GenSchema.HasDiscriminator = tpe.HasDiscriminator
	if tpe.IsArray && tpe.ElemType.IsBaseType {
		sg.GenSchema.ValueExpression += "()"
	}

	if Debug {
		log.Println("gschema nullable", sg.GenSchema.IsNullable)
	}
	if e := sg.buildAdditionalProperties(); e != nil {
		return e
	}

	prev := sg.GenSchema
	if sg.Untyped {
		if Debug {
			log.Println("untyped resolve", sg.Named || sg.IsTuple || sg.Required || sg.GenSchema.Required)
			b, _ := json.MarshalIndent(sg.Schema, "", "  ")
			log.Println(string(b))
		}
		tpe, err = sg.TypeResolver.ResolveSchema(nil, !sg.Named, sg.Named || sg.IsTuple || sg.Required || sg.GenSchema.Required)
	} else {
		if Debug {
			log.Printf("typed resolve, isAnonymous(%t), n: %t, t: %t, sgr: %t, sr: %t, isRequired(%t), BaseType(%t)", !sg.Named, sg.Named, sg.IsTuple, sg.Required, sg.GenSchema.Required, sg.Named || sg.IsTuple || sg.Required || sg.GenSchema.Required, sg.GenSchema.IsBaseType)
			b, _ := json.MarshalIndent(sg.Schema, "", "  ")
			log.Println(string(b))
		}
		tpe, err = sg.TypeResolver.ResolveSchema(&sg.Schema, !sg.Named, sg.Named || sg.IsTuple || sg.Required || sg.GenSchema.Required)
	}
	if err != nil {
		return err
	}
	otn := tpe.IsNullable
	tpe.IsNullable = tpe.IsNullable || nullableOverride
	sg.GenSchema.resolvedType = tpe
	sg.GenSchema.IsComplexObject = prev.IsComplexObject
	sg.GenSchema.IsMap = prev.IsMap
	sg.GenSchema.IsAdditionalProperties = prev.IsAdditionalProperties
	sg.GenSchema.IsBaseType = sg.GenSchema.HasDiscriminator

	if Debug {
		log.Println("gschema nnullable", sg.GenSchema.IsNullable, otn, nullableOverride)
		b, _ := json.MarshalIndent(sg.Schema, "", "  ")
		log.Println(string(b))
	}
	if err := sg.buildProperties(); err != nil {
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

	if err := sg.buildAliased(); err != nil {
		return err
	}

	if Debug {
		log.Printf("finished gen schema for %q\n", sg.Name)
	}
	return nil
}
