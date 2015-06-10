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
	"github.com/casualjim/go-swagger/swag"
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
	for pn, p := range schema.Properties {
		var required bool
		for _, v := range schema.Required {
			if v == pn {
				required = true
				break
			}
		}

		gmp, err := makeGenModelProperty2(propGenBuildParams{
			Path:      "\"" + pn + "\"",
			ParamName: swag.ToJSONName(pn),
			Accessor:  swag.ToGoName(pn),
			Receiver:  receiver,
			IndexVar:  "i",
			ValueExpr: receiver + "." + swag.ToGoName(pn),
			Schema:    p,
			Required:  required,
		})
		if err != nil {
			return nil, err
		}
		props[swag.ToJSONName(pn)] = gmp
	}
	for _, p := range schema.AllOf {
		if p.Ref.GetURL() != nil {
			tn := filepath.Base(p.Ref.GetURL().Fragment)
			p = specDoc.Spec().Definitions[tn]
		}
		mod, err := makeCodegenModel(name, pkg, p, specDoc)
		if err != nil {
			return nil, err
		}
		if mod != nil {
			for _, prop := range mod.Properties {
				props[prop.ParamName] = prop
			}
		}
	}

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
		ClassName:      swag.ToGoName(name),
		Name:           swag.ToJSONName(name),
		ReceiverName:   receiver,
		Properties:     properties,
		Description:    schema.Description,
		DocString:      modelDocString(swag.ToGoName(name), schema.Description),
		HumanClassName: swag.ToHumanNameLower(swag.ToGoName(name)),
		HasValidations: hasValidations,
	}, nil
}

type genModel struct {
	Package        string             //`json:"package,omitempty"`
	ReceiverName   string             //`json:"receiverName,omitempty"`
	ClassName      string             //`json:"classname,omitempty"`
	Name           string             //`json:"name,omitempty"`
	Description    string             //`json:"description,omitempty"`
	Properties     []genModelProperty //`json:"properties,omitempty"`
	DocString      string             //`json:"docString,omitempty"`
	HumanClassName string             //`json:"humanClassname,omitempty"`
	Imports        map[string]string  //`json:"imports,omitempty"`
	DefaultImports []string           //`json:"defaultImports,omitempty"`
	HasValidations bool               //`json:"hasValidatins,omitempty"`
}

func modelDocString(className, desc string) string {
	return commentedLines(fmt.Sprintf("%s %s", className, desc))
}

type propGenBuildParams struct {
	Path         string
	ParamName    string
	Accessor     string
	Receiver     string
	IndexVar     string
	ValueExpr    string
	Schema       spec.Schema
	Required     bool
	TypeResolver *typeResolver
}

func (pg propGenBuildParams) NewSliceBranch(schema *spec.Schema) propGenBuildParams {
	indexVar := pg.IndexVar
	pg.Path = "fmt.Sprintf(\"%s.%v\", " + pg.Path + ", " + indexVar + ")"
	pg.IndexVar = indexVar + "i"
	pg.ValueExpr = pg.ValueExpr + "[" + indexVar + "]"
	pg.Schema = *schema
	pg.Required = false
	return pg
}

func makeGenModelProperty2(params propGenBuildParams) (genModelProperty, error) {
	// log.Printf("property: (path %s) (param %s) (accessor %s) (receiver %s) (indexVar %s) (expr %s) required %t", path, paramName, accessor, receiver, indexVar, valueExpression, required)
	ex := ""
	if params.Schema.Example != nil {
		ex = fmt.Sprintf("%#v", params.Schema.Example)
	}
	validations, err := modelValidations2(params)
	if err != nil {
		return genModelProperty{}, err
	}

	ctx := makeGenValidations(validations)

	singleSchemaSlice := params.Schema.Items != nil && params.Schema.Items.Schema != nil
	var items []genModelProperty
	if singleSchemaSlice {
		ctx.HasSliceValidations = true

		elProp, err := makeGenModelProperty2(params.NewSliceBranch(params.Schema.Items.Schema))
		if err != nil {
			return genModelProperty{}, err
		}
		items = []genModelProperty{
			elProp,
			//makeGenModelProperty("fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", paramName, accessor, receiver, indexVar+"i", valueExpression+"["+indexVar+"]", *params.Schema.Items.Schema, false),
		}
	} else if params.Schema.Items != nil {
		for _, s := range params.Schema.Items.Schemas {
			elProp, err := makeGenModelProperty2(params.NewSliceBranch(&s))
			if err != nil {
				return genModelProperty{}, err
			}
			items = append(items, elProp)
			//items = append(items, makeGenModelProperty("fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", paramName, accessor, receiver, indexVar+"i", valueExpression+"["+indexVar+"]", s, false))
		}
	}

	allowsAdditionalItems :=
		params.Schema.AdditionalItems != nil &&
			(params.Schema.AdditionalItems.Allows || params.Schema.AdditionalItems.Schema != nil)
	hasAdditionalItems := allowsAdditionalItems && !singleSchemaSlice
	var additionalItems *genModelProperty
	if params.Schema.AdditionalItems != nil && params.Schema.AdditionalItems.Schema != nil {
		it, err := makeGenModelProperty2(params.NewSliceBranch(params.Schema.AdditionalItems.Schema))
		if err != nil {
			return genModelProperty{}, err
		}
		additionalItems = &it
	}

	ctx.HasSliceValidations = len(items) > 0 || hasAdditionalItems
	ctx.HasValidations = ctx.HasValidations || ctx.HasSliceValidations

	xmlName := params.ParamName
	if params.Schema.XML != nil {
		if params.Schema.XML.Name != "" {
			xmlName = params.Schema.XML.Name
			if params.Schema.XML.Attribute {
				xmlName += ",attr"
			}
		}
	}

	return genModelProperty{
		sharedParam:     ctx,
		DataType:        ctx.Type,
		Example:         ex,
		DocString:       propertyDocString(params.Accessor, params.Schema.Description, ex),
		Description:     params.Schema.Description,
		ReceiverName:    params.Receiver,
		IsComplexObject: !ctx.IsPrimitive && !ctx.IsCustomFormatter && !ctx.IsContainer,

		HasAdditionalItems:    hasAdditionalItems,
		AllowsAdditionalItems: allowsAdditionalItems,
		AdditionalItems:       additionalItems,

		Items:             items,
		ItemsLen:          len(items),
		SingleSchemaSlice: singleSchemaSlice,

		XMLName: xmlName,
	}, nil
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
	XMLName               string             //`json:"xmlName,omitempty"`
}

func modelValidations2(params propGenBuildParams) (commonValidations, error) {
	tpe, err := params.TypeResolver.ResolveSchema(&params.Schema)
	if err != nil {
		return commonValidations{}, err
	}

	_, isPrimitive := primitives[tpe.GoType]
	_, isCustomFormatter := customFormatters[tpe.GoType]
	model := params.Schema

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      params.Accessor,
			ParamName:         params.ParamName,
			ValueExpression:   params.ValueExpr,
			IndexVar:          params.IndexVar,
			Path:              params.Path,
			IsContainer:       tpe.IsArray,
			IsPrimitive:       isPrimitive,
			IsCustomFormatter: isCustomFormatter,
			IsMap:             tpe.IsMap,
		},
		Required:         params.Required,
		Type:             tpe.GoType,
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
	}, nil
}

func propertyDocString(propertyName, description, example string) string {
	ex := ""
	if strings.TrimSpace(example) != "" {
		ex = " eg.\n\n    " + example
	}
	return commentedLines(fmt.Sprintf("%s %s%s", propertyName, description, ex))
}
