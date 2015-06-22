package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// GenerateServerOperation generates a parameter model, parameter validator, http handler implementations for a given operation
// It also generates an operation handler interface that uses the parameter model for handling a valid request.
// Allows for specifying a list of tags to include only certain tags for the generation
func GenerateServerOperation(operationNames, tags []string, includeHandler, includeParameters bool, opts GenOpts) error {
	// Load the spec
	specPath, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	if len(operationNames) == 0 {
		operationNames = specDoc.OperationIDs()
	}

	for _, operationName := range operationNames {
		operation, ok := specDoc.OperationForName(operationName)
		if !ok {
			return fmt.Errorf("operation %q not found in %s", operationName, specPath)
		}

		generator := operationGenerator{
			Name:                 operationName,
			APIPackage:           opts.APIPackage,
			ModelsPackage:        opts.ModelPackage,
			ClientPackage:        opts.ClientPackage,
			ServerPackage:        opts.ServerPackage,
			Operation:            *operation,
			SecurityRequirements: specDoc.SecurityRequirementsFor(operation),
			Principal:            opts.Principal,
			Target:               filepath.Join(opts.Target, opts.APIPackage),
			Tags:                 tags,
			IncludeHandler:       includeHandler,
			IncludeParameters:    includeParameters,
			DumpData:             opts.DumpData,
			Doc:                  specDoc,
		}
		if err := generator.Generate(); err != nil {
			return err
		}
	}
	return nil
}

type operationGenerator struct {
	Name                 string
	Authorized           bool
	APIPackage           string
	ModelsPackage        string
	ServerPackage        string
	ClientPackage        string
	Operation            spec.Operation
	SecurityRequirements []spec.SecurityRequirement
	Principal            string
	Target               string
	Tags                 []string
	data                 interface{}
	pkg                  string
	cname                string
	IncludeHandler       bool
	IncludeParameters    bool
	DumpData             bool
	Doc                  *spec.Document
}

type codeGenOpBuilder struct {
	Name          string
	APIPackage    string
	ModelsPackage string
	Principal     string
	Target        string
	Operation     spec.Operation
	Doc           *spec.Document
	Authed        bool
}

func (o *operationGenerator) Generate() error {
	// Build a list of codegen operations based on the tags,
	// the tag decides the actual package for an operation
	// the user specified package serves as root for generating the directory structure
	var operations []genOperation
	authed := len(o.SecurityRequirements) > 0

	var bldr codeGenOpBuilder
	bldr.Name = o.Name
	bldr.ModelsPackage = o.ModelsPackage
	bldr.Principal = o.Principal
	bldr.Target = o.Target
	bldr.Operation = o.Operation
	bldr.Authed = authed
	bldr.Doc = o.Doc

	for _, tag := range o.Operation.Tags {
		if len(o.Tags) == 0 {
			bldr.APIPackage = tag
			op, err := makeCodegenOperation(bldr)
			if err != nil {
				return err
			}
			operations = append(operations, op)
			continue
		}
		for _, ft := range o.Tags {
			if ft == tag {
				bldr.APIPackage = tag
				op, err := makeCodegenOperation(bldr)
				if err != nil {
					return err
				}
				operations = append(operations, op)
				break
			}
		}

	}
	if len(operations) == 0 {
		bldr.APIPackage = o.APIPackage
		op, err := makeCodegenOperation(bldr)
		if err != nil {
			return err
		}
		operations = append(operations, op)
	}

	for _, op := range operations {
		if o.DumpData {
			bb, _ := json.MarshalIndent(swag.ToDynamicJSON(op), "", " ")
			fmt.Fprintln(os.Stdout, string(bb))
			continue
		}
		o.data = op
		o.pkg = op.Package
		o.cname = op.ClassName

		if o.IncludeHandler {
			if err := o.generateHandler(); err != nil {
				return fmt.Errorf("handler: %s", err)
			}
			log.Println("generated handler", op.Package+"."+op.ClassName)
		}

		if o.IncludeParameters && len(o.Operation.Parameters) > 0 {
			if err := o.generateParameterModel(); err != nil {
				return fmt.Errorf("parameters: %s", err)
			}
			log.Println("generated parameters", op.Package+"."+op.ClassName+"Parameters")
		}

		if len(o.Operation.Parameters) == 0 {
			log.Println("no parameters for operation", op.Package+"."+op.ClassName)
		}
	}

	return nil
}

func (o *operationGenerator) generateHandler() error {
	buf := bytes.NewBuffer(nil)

	if err := operationTemplate.Execute(buf, o.data); err != nil {
		return err
	}
	log.Println("rendered handler template:", o.pkg+"."+o.cname)

	fp := filepath.Join(o.ServerPackage, o.Target)
	if len(o.Operation.Tags) > 0 {
		fp = filepath.Join(fp, o.pkg)
	}
	return writeToFile(fp, o.Name, buf.Bytes())
}

func (o *operationGenerator) generateParameterModel() error {
	buf := bytes.NewBuffer(nil)

	if err := parameterTemplate.Execute(buf, o.data); err != nil {
		return err
	}
	log.Println("rendered parameters template:", o.pkg+"."+o.cname+"Parameters")

	fp := filepath.Join(o.ServerPackage, o.Target)
	if len(o.Operation.Tags) > 0 {
		fp = filepath.Join(fp, o.pkg)
	}
	return writeToFile(fp, o.Name+"Parameters", buf.Bytes())
}

func makeCodegenOperation(b codeGenOpBuilder) (genOperation, error) {
	receiver := "o"
	resolver := typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}

	operation := b.Operation
	var params, qp, pp, hp, fp []genParameter
	var hasQueryParams bool
	for _, p := range operation.Parameters {
		cp, err := makeCodegenParameter(receiver, &resolver, p)
		if err != nil {
			return genOperation{}, err
		}
		if cp.IsQueryParam {
			hasQueryParams = true
			qp = append(qp, cp)
		}
		if cp.IsFormParam {
			fp = append(fp, cp)
		}
		if cp.IsPathParam {
			pp = append(pp, cp)
		}
		if cp.IsHeaderParam {
			hp = append(hp, cp)
		}
		params = append(params, cp)
	}

	var successModel string
	var returnsPrimitive, returnsFormatted, returnsContainer, returnsMap bool
	var responses map[int]genResponse
	var defaultResponse *genResponse
	if operation.Responses != nil {
		if r, ok := operation.Responses.StatusCodeResponses[200]; ok {
			tn, err := resolver.ResolveSchema(r.Schema, true)
			if err != nil {
				return genOperation{}, err
			}
			_, returnsPrimitive = primitives[tn.GoType]
			_, returnsFormatted = customFormatters[tn.GoType]
			returnsContainer = tn.IsArray
			returnsMap = tn.IsMap
			successModel = tn.GoType
		}
		for k, v := range operation.Responses.StatusCodeResponses {
			gr, err := makeGenResponse(b.APIPackage, receiver, swag.ToGoName(b.Name+" "+httpkit.Statuses[k]), k/100 == 2, &resolver, v)
			if err != nil {
				return genOperation{}, err
			}
			if responses == nil {
				responses = make(map[int]genResponse)
			}
			responses[k] = gr
		}
		if operation.Responses.Default != nil {
			gr, err := makeGenResponse(b.APIPackage, receiver, swag.ToGoName(b.Name+" Default Response"), false, &resolver, *operation.Responses.Default)
			if err != nil {
				return genOperation{}, err
			}
			defaultResponse = &gr
		}
	}

	prin := b.Principal
	if prin == "" {
		prin = "interface{}"
	}

	zero, ok := zeroes[successModel]
	if !ok {
		zero = "nil"
	}

	return genOperation{
		Package:              b.APIPackage,
		ClassName:            swag.ToGoName(b.Name),
		Name:                 swag.ToJSONName(b.Name),
		Description:          operation.Description,
		DocString:            operationDocString(swag.ToGoName(b.Name), operation),
		ReceiverName:         receiver,
		HumanClassName:       swag.ToHumanNameLower(swag.ToGoName(b.Name)),
		DefaultImports:       []string{filepath.Join(baseImport(b.Target), b.ModelsPackage)},
		Params:               params,
		Summary:              operation.Summary,
		QueryParams:          qp,
		PathParams:           pp,
		HeaderParams:         hp,
		FormParams:           fp,
		HasQueryParams:       hasQueryParams,
		SuccessModel:         successModel,
		SuccessZero:          zero,
		ReturnsPrimitive:     returnsPrimitive,
		ReturnsFormatted:     returnsFormatted,
		ReturnsContainer:     returnsContainer,
		ReturnsMap:           returnsMap,
		ReturnsComplexObject: !returnsPrimitive && !returnsFormatted && !returnsContainer && !returnsMap,
		Authorized:           b.Authed,
		Principal:            prin,
		Responses:            responses,
		DefaultResponse:      defaultResponse,
	}, nil
}

func operationDocString(name string, operation spec.Operation) string {
	hdr := fmt.Sprintf("%s %s", name, operation.Description)
	ed := operation.ExternalDocs
	var txtFoot string
	if ed != nil {
		if ed.Description != "" && ed.URL != "" {
			txtFoot = fmt.Sprintf("\n%s\nSee: %s", ed.Description, ed.URL)
		}
		if ed.URL != "" {
			txtFoot = "\nSee: " + ed.URL
		}
	}
	return commentedLines(strings.Join([]string{hdr, txtFoot}, "\n"))
}

type genOperationGroup struct {
	Name       string
	Operations []genOperation

	DocString      string
	Imports        map[string]string
	DefaultImports []string
}

type genOperation struct {
	Package        string
	ReceiverName   string
	ClassName      string
	Name           string
	HumanClassName string

	Summary      string
	Description  string
	DocString    string
	ExternalDocs string

	Imports        map[string]string
	DefaultImports []string

	Authorized bool
	Principal  string

	SuccessModel         string
	SuccessZero          string
	ReturnsPrimitive     bool
	ReturnsFormatted     bool
	ReturnsContainer     bool
	ReturnsComplexObject bool
	ReturnsMap           bool
	Responses            map[int]genResponse
	DefaultResponse      *genResponse

	Params         []genParameter
	QueryParams    []genParameter
	PathParams     []genParameter
	HeaderParams   []genParameter
	FormParams     []genParameter
	HasQueryParams bool
	HasFormParams  bool
	HasFileParams  bool
}

func makeCodegenParameter(receiver string, resolver *typeResolver, param spec.Parameter) (genParameter, error) {
	var ctx sharedParam
	var child *genParameterItem

	if param.In == "body" {
		sc := schemaGenContext{
			Path:         param.Name,
			Name:         param.Name,
			Receiver:     receiver,
			IndexVar:     "i",
			Schema:       *param.Schema,
			Required:     param.Required,
			TypeResolver: resolver,
			Named:        false,
		}
		if err := sc.makeGenSchema(); err != nil {
			return genParameter{}, err
		}

	} else {
		ctx = makeGenValidations(paramValidations(receiver, param))
		thisItem := genParameterItem{}
		thisItem.sharedParam = ctx
		thisItem.ValueExpression = ctx.IndexVar + "c"
		thisItem.CollectionFormat = param.CollectionFormat
		thisItem.Converter = stringConverters[ctx.Type]
		thisItem.Location = param.In

		if param.Items != nil {
			it := makeCodegenParamItem(
				"fmt.Sprintf(\"%s.%v\", "+ctx.Path+", "+ctx.IndexVar+")",
				ctx.ParamName,
				ctx.PropertyName,
				ctx.IndexVar+"i",
				ctx.IndexVar+"c["+ctx.IndexVar+"]",
				thisItem,
				*param.Items,
			)
			child = &it
		}

	}

	return genParameter{
		sharedParam:      ctx,
		Description:      param.Description,
		ReceiverName:     receiver,
		IsQueryParam:     param.In == "query",
		IsBodyParam:      param.In == "body",
		IsHeaderParam:    param.In == "header",
		IsPathParam:      param.In == "path",
		IsFormParam:      param.In == "formData",
		IsFileParam:      param.Type == "file",
		CollectionFormat: param.CollectionFormat,
		Child:            child,
		Location:         param.In,
		Converter:        stringConverters[ctx.Type],
		Formatter:        stringFormatters[ctx.Type],
	}, nil
}

type genParameter struct {
	sharedParam
	ReceiverName     string
	Title            string
	Description      string
	IsQueryParam     bool
	IsFormParam      bool
	IsPathParam      bool
	IsHeaderParam    bool
	IsBodyParam      bool
	IsFileParam      bool
	CollectionFormat string
	Child            *genParameterItem
	BodyParam        *genParameter
	Converter        string
	Formatter        string
	Parent           *genParameterItem
	Location         string
}

func makeCodegenParamItem(path, paramName, accessor, indexVar, valueExpression string, parent genParameterItem, items spec.Items) genParameterItem {
	ctx := makeGenValidations(paramItemValidations(path, paramName, accessor, indexVar, valueExpression, items))

	res := genParameterItem{}
	res.sharedParam = ctx
	res.CollectionFormat = items.CollectionFormat
	res.Parent = &parent
	res.Converter = stringConverters[ctx.Type]
	res.Formatter = stringFormatters[ctx.Type]
	res.Location = parent.Location
	res.ValueExpression = "value"

	var child *genParameterItem
	if items.Items != nil {
		it := makeCodegenParamItem(
			"fmt.Sprintf(\"%s.%v\", "+ctx.Path+", "+ctx.IndexVar+")",
			ctx.ParamName,
			ctx.PropertyName,
			ctx.IndexVar+"i",
			ctx.IndexVar+"c["+ctx.IndexVar+"]",
			res,
			*items.Items,
		)
		child = &it
	}
	res.Child = child

	return res
}

type genParameterItem struct {
	sharedParam
	CollectionFormat string
	Child            *genParameterItem
	Parent           *genParameterItem
	Converter        string
	Formatter        string
	Location         string
}

type sharedParam struct {
	genValidations
	propertyDescriptor
}

// GenParameter is used to represent
// a parameter or a header for code generation.
type GenParameter struct {
	sharedValidations
	Name              string
	Path              string
	ValueExpression   string
	IndexVar          string
	ReceiverName      string
	Title             string
	Description       string
	Converter         string
	Formatter         string
	Type              string
	Format            string
	Items             *spec.Items
	IsPrimitive       bool
	IsCustomFormatter bool
	IsArray           bool
	Default           interface{}
	Enum              []interface{}
}

func paramItemValidations(path, paramName, accessor, indexVar, valueExpression string, items spec.Items) commonValidations {
	tpe := resolveSimpleType(items.Type, items.Format, items.Items)
	_, isPrimitive := primitives[tpe]
	_, isCustomFormatter := customFormatters[tpe]

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      accessor,
			ParamName:         paramName,
			ValueExpression:   valueExpression,
			IndexVar:          indexVar,
			Path:              path,
			IsContainer:       items.Items != nil || tpe == "array",
			IsPrimitive:       isPrimitive,
			IsCustomFormatter: isCustomFormatter,
			IsMap:             strings.HasPrefix(tpe, "map"),
		},

		Type:             tpe,
		Format:           items.Format,
		Items:            items.Items,
		Default:          items.Default,
		Maximum:          items.Maximum,
		ExclusiveMaximum: items.ExclusiveMaximum,
		Minimum:          items.Minimum,
		ExclusiveMinimum: items.ExclusiveMinimum,
		MaxLength:        items.MaxLength,
		MinLength:        items.MinLength,
		Pattern:          items.Pattern,
		MaxItems:         items.MaxItems,
		MinItems:         items.MinItems,
		UniqueItems:      items.UniqueItems,
		MultipleOf:       items.MultipleOf,
		Enum:             items.Enum,
	}
}

func paramValidations(receiver string, param spec.Parameter) commonValidations {
	accessor := swag.ToGoName(param.Name)
	paramName := swag.ToJSONName(param.Name)

	tpe := typeForParameter(param)
	_, isPrimitive := primitives[tpe]
	_, isCustomFormatter := customFormatters[tpe]

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      accessor,
			ParamName:         paramName,
			ValueExpression:   fmt.Sprintf("%s.%s", receiver, accessor),
			IndexVar:          "i",
			Path:              "\"" + paramName + "\"",
			IsContainer:       param.Items != nil || tpe == "array",
			IsPrimitive:       isPrimitive,
			IsCustomFormatter: isCustomFormatter,
			IsMap:             strings.HasPrefix(tpe, "map"),
		},
		Required:         param.Required,
		Type:             tpe,
		Format:           param.Format,
		Items:            param.Items,
		Default:          param.Default,
		Maximum:          param.Maximum,
		ExclusiveMaximum: param.ExclusiveMaximum,
		Minimum:          param.Minimum,
		ExclusiveMinimum: param.ExclusiveMinimum,
		MaxLength:        param.MaxLength,
		MinLength:        param.MinLength,
		Pattern:          param.Pattern,
		MaxItems:         param.MaxItems,
		MinItems:         param.MinItems,
		UniqueItems:      param.UniqueItems,
		MultipleOf:       param.MultipleOf,
		Enum:             param.Enum,
	}
}

func makeGenValidations(s commonValidations) sharedParam {
	hasValidations := s.Required

	var defVal string
	if s.Default != nil {
		hasValidations = false
		defVal = fmt.Sprintf("%#v", s.Default)
	}

	var format string
	if s.IsCustomFormatter {
		hasValidations = true
		format = s.Format
	}

	var maxLength int64
	if s.MaxLength != nil {
		hasValidations = true
		maxLength = *s.MaxLength
	}

	var minLength int64
	if s.MinLength != nil {
		hasValidations = true
		minLength = *s.MinLength
	}

	var minimum float64
	if s.Minimum != nil {
		hasValidations = true
		minimum = *s.Minimum
	}

	var maximum float64
	if s.Maximum != nil {
		hasValidations = true
		maximum = *s.Maximum
	}

	var multipleOf float64
	if s.MultipleOf != nil {
		hasValidations = true
		multipleOf = *s.MultipleOf
	}

	var needsSize bool
	hasSliceValidations := s.UniqueItems
	var maxItems int64
	if s.MaxItems != nil {
		hasSliceValidations = true
		needsSize = true
		maxItems = *s.MaxItems
	}

	var minItems int64
	if s.MinItems != nil {
		hasSliceValidations = true
		needsSize = true
		minItems = *s.MinItems
	}

	var enum string
	if len(s.Enum) > 0 {
		hasValidations = true
		enum = fmt.Sprintf("%#v", s.Enum)
	}

	return sharedParam{
		propertyDescriptor: s.propertyDescriptor,
		genValidations: genValidations{
			Type:                s.Type,
			Required:            s.Required,
			DefaultValue:        defVal,
			MaxLength:           maxLength,
			MinLength:           minLength,
			Pattern:             s.Pattern,
			MultipleOf:          multipleOf,
			Minimum:             minimum,
			Maximum:             maximum,
			ExclusiveMinimum:    s.ExclusiveMinimum,
			ExclusiveMaximum:    s.ExclusiveMaximum,
			Enum:                enum,
			HasValidations:      hasValidations,
			Format:              format,
			MinItems:            minItems,
			MaxItems:            maxItems,
			UniqueItems:         s.UniqueItems,
			HasSliceValidations: hasSliceValidations,
			NeedsSize:           needsSize,
		},
	}
}

func makeGenResponse(pkg, receiver, name string, isSuccess bool, resolver *typeResolver, response spec.Response) (genResponse, error) {
	var headers []genHeader
	for hName, header := range response.Headers {
		headers = append(headers, makeGenHeader(receiver, hName, header))
	}
	tpe, err := resolver.ResolveSchema(response.Schema, true)
	if err != nil {
		return genResponse{}, err
	}

	var model string
	if response.Schema != nil {
		model = tpe.GoType
	}

	return genResponse{
		Package:         pkg,
		ReceiverName:    receiver,
		ClassName:       swag.ToGoName(name),
		Name:            swag.ToJSONName(name),
		HumanClassName:  swag.ToHumanNameLower(name),
		DefaultImports:  nil,
		Imports:         nil,
		Headers:         headers,
		IsSuccess:       isSuccess,
		Type:            model,
		IsPrimitive:     tpe.IsPrimitive,
		IsFormatted:     tpe.IsCustomFormatter,
		IsContainer:     tpe.IsArray,
		IsComplexObject: tpe.IsComplexObject,
	}, nil
}

type genResponse struct {
	Package        string
	ReceiverName   string
	ClassName      string
	Name           string
	HumanClassName string
	Description    string
	IsSuccess      bool

	Headers         []genHeader
	Zero            string
	IsPrimitive     bool
	IsFormatted     bool
	IsContainer     bool
	IsComplexObject bool
	IsMap           bool
	Type            string
	Format          string

	Imports        map[string]string
	DefaultImports []string
}

type genHeader struct {
	sharedParam

	ReceiverName string
	Title        string
	Description  string
	Converter    string
	Formatter    string
}

func makeGenHeader(receiver, name string, header spec.Header) (result genHeader) {
	ctx := makeGenValidations(headerValidations(receiver, name, header))
	result.sharedParam = ctx
	result.ReceiverName = receiver
	result.Description = header.Description
	result.Converter = stringConverters[ctx.Type]
	result.Formatter = stringFormatters[ctx.Type]
	return
}

func headerValidations(receiver, name string, header spec.Header) commonValidations {
	accessor := swag.ToGoName(name)
	paramName := name

	tpe := typeForHeader(header)
	_, isPrimitive := primitives[tpe]
	_, isCustomFormatter := customFormatters[tpe]

	//res := sharedValidations{
	//Maximum:          header.Maximum,
	//ExclusiveMaximum: header.ExclusiveMaximum,
	//Minimum:          header.Minimum,
	//ExclusiveMinimum: header.ExclusiveMinimum,
	//MaxLength:        header.MaxLength,
	//MinLength:        header.MinLength,
	//Pattern:          header.Pattern,
	//MaxItems:         header.MaxItems,
	//MinItems:         header.MinItems,
	//UniqueItems:      header.UniqueItems,
	//MultipleOf:       header.MultipleOf,
	//Enum:             header.Enum,
	//}

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      accessor,
			ParamName:         paramName,
			ValueExpression:   fmt.Sprintf("%s.%s", receiver, accessor),
			IndexVar:          "i",
			Path:              "\"" + paramName + "\"",
			IsContainer:       header.Items != nil || tpe == "array",
			IsPrimitive:       isPrimitive,
			IsCustomFormatter: isCustomFormatter,
			IsMap:             strings.HasPrefix(tpe, "map"),
		},
		Type:             tpe,
		Format:           header.Format,
		Items:            header.Items,
		Default:          header.Default,
		Maximum:          header.Maximum,
		ExclusiveMaximum: header.ExclusiveMaximum,
		Minimum:          header.Minimum,
		ExclusiveMinimum: header.ExclusiveMinimum,
		MaxLength:        header.MaxLength,
		MinLength:        header.MinLength,
		Pattern:          header.Pattern,
		MaxItems:         header.MaxItems,
		MinItems:         header.MinItems,
		UniqueItems:      header.UniqueItems,
		MultipleOf:       header.MultipleOf,
		Enum:             header.Enum,
	}
}
