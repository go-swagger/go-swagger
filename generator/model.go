package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

var (
	modelTemplate          *template.Template
	modelValidatorTemplate *template.Template
)

func init() {
	bv, _ := Asset("templates/modelvalidator.gotmpl")
	modelValidatorTemplate = template.Must(template.New("modelvalidator").Parse(string(bv)))

	bm, _ := Asset("templates/model.gotmpl")
	modelTemplate = template.Must(template.New("model").Parse(string(bm)))

}

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
	mod := makeCodegenModel(m.Name, m.Target, m.Model, m.SpecDoc)
	if m.DumpData {
		bb, _ := json.MarshalIndent(util.ToDynamicJSON(mod), "", " ")
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

func makeCodegenModel(name, pkg string, schema spec.Schema, specDoc *spec.Document) *genModel {
	receiver := "m"
	props := make(map[string]genModelProperty)
	for pn, p := range schema.Properties {
		var required bool
		for _, v := range schema.Required {
			if v == pn {
				required = true
				break
			}
		}
		props[util.ToJSONName(pn)] = makeGenModelProperty(
			"\""+pn+"\"",
			util.ToJSONName(pn),
			util.ToGoName(pn),
			receiver,
			"i",
			receiver+"."+util.ToGoName(pn),
			p,
			required)
	}
	for _, p := range schema.AllOf {
		if p.Ref.GetURL() != nil {
			tn := filepath.Base(p.Ref.GetURL().Fragment)
			p = specDoc.Spec().Definitions[tn]
		}
		mod := makeCodegenModel(name, pkg, p, specDoc)
		if mod != nil {
			for _, prop := range mod.Properties {
				props[prop.ParamName] = prop
			}
		}
	}

	var properties []genModelProperty
	for _, v := range props {
		properties = append(properties, v)
	}

	return &genModel{
		Package:        filepath.Base(pkg),
		ClassName:      util.ToGoName(name),
		Name:           util.ToJSONName(name),
		ReceiverName:   receiver,
		Properties:     properties,
		Description:    schema.Description,
		DocString:      modelDocString(util.ToGoName(name), schema.Description),
		HumanClassName: util.ToHumanNameLower(util.ToGoName(name)),
	}
}

type genModel struct {
	Package        string             //`json:"package,omitempty"`
	ReceiverName   string             ////`json:"receiverName,omitempty"`
	ClassName      string             //`json:"classname,omitempty"`
	Name           string             //`json:"name,omitempty"`
	Description    string             //`json:"description,omitempty"`
	Properties     []genModelProperty //`json:"properties,omitempty"`
	DocString      string             //`json:"docString,omitempty"`
	HumanClassName string             //`json:"humanClassname,omitempty"`
	Imports        []string           //`json:"imports,omitempty"`
}

func modelDocString(className, desc string) string {
	return commentedLines(fmt.Sprintf("%s %s", className, desc))
}

func makeGenModelProperty(path, paramName, accessor, receiver, indexVar, valueExpression string, schema spec.Schema, required bool) genModelProperty {
	// log.Printf("property: (path %s) (param %s) (accessor %s) (receiver %s) (indexVar %s) (expr %s) required %t", path, paramName, accessor, receiver, indexVar, valueExpression, required)
	ex := ""
	if schema.Example != nil {
		ex = fmt.Sprintf("%#v", schema.Example)
	}

	ctx := makeGenValidations(modelValidations(path, paramName, accessor, indexVar, valueExpression, "", required, schema))

	singleSchemaSlice := schema.Items != nil && schema.Items.Schema != nil
	var items []genModelProperty
	if singleSchemaSlice {
		ctx.HasSliceValidations = true
		items = []genModelProperty{
			makeGenModelProperty("fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", paramName, accessor, receiver, indexVar+"i", valueExpression+"["+indexVar+"]", *schema.Items.Schema, false),
		}
	} else if schema.Items != nil {
		for _, s := range schema.Items.Schemas {
			items = append(items, makeGenModelProperty("fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", paramName, accessor, receiver, indexVar+"i", valueExpression+"["+indexVar+"]", s, false))
		}
	}

	allowsAdditionalItems :=
		schema.AdditionalItems != nil &&
			(schema.AdditionalItems.Allows || schema.AdditionalItems.Schema != nil)
	hasAdditionalItems := allowsAdditionalItems && !singleSchemaSlice
	var additionalItems *genModelProperty
	if schema.AdditionalItems != nil && schema.AdditionalItems.Schema != nil {
		it := makeGenModelProperty("fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", paramName, accessor, receiver, indexVar+"i", valueExpression+"["+indexVar+"]", *schema.AdditionalItems.Schema, false)
		additionalItems = &it
	}

	ctx.HasSliceValidations = len(items) > 0 || hasAdditionalItems
	ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

	return genModelProperty{
		sharedParam:     ctx,
		DataType:        ctx.Type,
		Example:         ex,
		DocString:       propertyDocString(accessor, schema.Description, ex),
		Description:     schema.Description,
		ReceiverName:    receiver,
		IsComplexObject: !ctx.IsPrimitive && !ctx.IsCustomFormatter && !ctx.IsContainer,

		HasAdditionalItems:    hasAdditionalItems,
		AllowsAdditionalItems: allowsAdditionalItems,
		AdditionalItems:       additionalItems,

		Items:             items,
		ItemsLen:          len(items),
		SingleSchemaSlice: singleSchemaSlice,
	}
}

// TODO:
// untyped data requires a cast somehow to the inner type
//
// wants an IsNested or IsAnonymous flag for schemas with properties
// wants an IsMap or IsDynamic flag for schemas with additional properties set
//

type genModelProperty struct {
	sharedParam
	Example               string             //`json:"example,omitempty"`
	Description           string             //`json:"description,omitempty"`
	DataType              string             //`json:"dataType,omitempty"`
	DocString             string             //`json:"docString,omitempty"`
	Location              string             //`json:"location,omitempty"`
	ReceiverName          string             //`json:"receiverName,omitempty"`
	IsComplexObject       bool               //`json:"isComplex,omitempty"` // not slice, custom formatter or primitive
	SingleSchemaSlice     bool               //`json:"singleSchemaSlice,omitempty"`
	Items                 []genModelProperty //`json:"items,omitempty"`
	ItemsLen              int                //`json:"itemsLength,omitempty"`
	AllowsAdditionalItems bool               //`json:"allowsAdditionalItems,omitempty"`
	HasAdditionalItems    bool               //`json:"hasAdditionalItems,omitempty"`
	AdditionalItems       *genModelProperty  //`json:"additionalItems,omitempty"`
	Object                *genModelProperty  //`json:"object,omitempty"`
	Imports               []string           //`json:"imports,omitempty"`
}

func modelValidations(path, paramName, accessor, indexVar, valueExpression, pkg string, required bool, model spec.Schema) commonValidations {
	tpe := typeForSchema(&model, pkg)

	_, isPrimitive := primitives[tpe]
	_, isCustomFormatter := customFormatters[tpe]

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      accessor,
			ParamName:         paramName,
			ValueExpression:   valueExpression,
			IndexVar:          indexVar,
			Path:              path,
			IsContainer:       model.Items != nil || model.Type.Contains("array"),
			IsPrimitive:       isPrimitive,
			IsCustomFormatter: isCustomFormatter,
		},
		Required:         required,
		Type:             tpe,
		Format:           model.Format,
		Default:          model.Default,
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
		Enum:             model.Enum,
	}
}

func propertyDocString(propertyName, description, example string) string {
	ex := ""
	if strings.TrimSpace(example) != "" {
		ex = " eg.\n\n    " + example
	}
	return commentedLines(fmt.Sprintf("%s %s%s", propertyName, description, ex))
}
