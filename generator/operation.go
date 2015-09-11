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

func (o *operationGenerator) Generate() error {
	// Build a list of codegen operations based on the tags,
	// the tag decides the actual package for an operation
	// the user specified package serves as root for generating the directory structure
	var operations []GenOperation
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
			//op, err := makeCodegenOperation(bldr)
			op, err := bldr.MakeOperation()
			if err != nil {
				return err
			}
			operations = append(operations, op)
			continue
		}
		for _, ft := range o.Tags {
			if ft == tag {
				bldr.APIPackage = tag
				//op, err := makeCodegenOperation(bldr)
				op, err := bldr.MakeOperation()
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
		//op, err := makeCodegenOperation(bldr)
		op, err := bldr.MakeOperation()
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
		o.cname = swag.ToGoName(op.Name)

		if o.IncludeHandler {
			if err := o.generateHandler(); err != nil {
				return fmt.Errorf("handler: %s", err)
			}
			log.Println("generated handler", op.Package+"."+o.cname)
		}

		if o.IncludeParameters && len(o.Operation.Parameters) > 0 {
			if err := o.generateParameterModel(); err != nil {
				return fmt.Errorf("parameters: %s", err)
			}
			log.Println("generated parameters", op.Package+"."+o.cname+"Parameters")
		}

		if len(o.Operation.Parameters) == 0 {
			log.Println("no parameters for operation", op.Package+"."+o.cname)
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

type codeGenOpBuilder struct {
	Name          string
	APIPackage    string
	ModelsPackage string
	Principal     string
	Target        string
	Operation     spec.Operation
	Doc           *spec.Document
	Authed        bool
	ExtraSchemas  map[string]GenSchema
}

func (b codeGenOpBuilder) MakeOperation() (GenOperation, error) {
	resolver := typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	receiver := "o"

	operation := b.Operation
	var params, qp, pp, hp, fp []GenParameter
	var hasQueryParams bool
	for _, p := range operation.Parameters {
		cp, err := b.MakeParameter(receiver, &resolver, p)
		if err != nil {
			return GenOperation{}, err
		}
		if cp.IsQueryParam() {
			hasQueryParams = true
			qp = append(qp, cp)
		}
		if cp.IsFormParam() {
			fp = append(fp, cp)
		}
		if cp.IsPathParam() {
			pp = append(pp, cp)
		}
		if cp.IsHeaderParam() {
			hp = append(hp, cp)
		}
		params = append(params, cp)
	}

	var responses map[int]GenResponse
	var defaultResponse *GenResponse
	var successResponse *GenResponse
	if operation.Responses != nil {
		for k, v := range operation.Responses.StatusCodeResponses {
			isSuccess := k/100 == 2
			gr, err := b.MakeResponse(receiver, swag.ToJSONName(b.Name+" "+httpkit.Statuses[k]), isSuccess, &resolver, v)
			if err != nil {
				return GenOperation{}, err
			}
			if isSuccess {
				successResponse = &gr
			}
			if responses == nil {
				responses = make(map[int]GenResponse)
			}
			responses[k] = gr
		}

		if operation.Responses.Default != nil {
			gr, err := b.MakeResponse(receiver, b.Name+" default", false, &resolver, *operation.Responses.Default)
			if err != nil {
				return GenOperation{}, err
			}
			defaultResponse = &gr
		}
	}

	prin := b.Principal
	if prin == "" {
		prin = "interface{}"
	}

	fmt.Printf("the target directory for operation %s is %s\n", b.Name, b.Target)
	fmt.Printf("the base import path for operation %s is %s\n", b.Name, baseImport(b.Target))

	return GenOperation{
		Package:         b.APIPackage,
		Name:            b.Name,
		Description:     operation.Description,
		ReceiverName:    receiver,
		DefaultImports:  []string{filepath.ToSlash(filepath.Join(baseImport(b.Target), b.ModelsPackage))},
		Params:          params,
		Summary:         operation.Summary,
		QueryParams:     qp,
		PathParams:      pp,
		HeaderParams:    hp,
		FormParams:      fp,
		HasQueryParams:  hasQueryParams,
		Authorized:      b.Authed,
		Principal:       prin,
		Responses:       responses,
		DefaultResponse: defaultResponse,
		SuccessResponse: successResponse,
	}, nil
}

func (b codeGenOpBuilder) MakeResponse(receiver, name string, isSuccess bool, resolver *typeResolver, resp spec.Response) (GenResponse, error) {

	res := GenResponse{
		Package:        b.APIPackage,
		ReceiverName:   receiver,
		Name:           name,
		Description:    resp.Description,
		DefaultImports: nil,
		Imports:        nil,
		IsSuccess:      isSuccess,
	}

	for hName, header := range resp.Headers {
		res.Headers = append(res.Headers, b.MakeHeader(receiver, hName, header))
	}

	if resp.Schema != nil {
		sc := schemaGenContext{
			Path:         fmt.Sprintf("%q", name),
			Name:         name + "Body",
			Receiver:     receiver,
			ValueExpr:    receiver + "." + name,
			IndexVar:     "i",
			Schema:       *resp.Schema,
			Required:     true,
			TypeResolver: resolver,
			Named:        false,
			ExtraSchemas: make(map[string]GenSchema),
		}
		if err := sc.makeGenSchema(); err != nil {
			return GenResponse{}, err
		}

		for k, v := range sc.ExtraSchemas {
			b.ExtraSchemas[k] = v
		}

		schema := sc.GenSchema
		if schema.IsAnonymous {
			schema.Name = swag.ToGoName(sc.Name + " Body")
			nm := schema.Name
			if b.ExtraSchemas == nil {
				b.ExtraSchemas = make(map[string]GenSchema)
			}
			b.ExtraSchemas[schema.Name] = schema
			schema = GenSchema{}
			schema.IsAnonymous = false
			schema.GoType = nm
			schema.SwaggerType = nm
		}

		res.Schema = &schema
	}
	return res, nil
}

func (b codeGenOpBuilder) MakeHeader(receiver, name string, hdr spec.Header) GenHeader {
	hasNumberValidation := hdr.Maximum != nil || hdr.Minimum != nil || hdr.MultipleOf != nil
	hasStringValidation := hdr.MaxLength != nil || hdr.MinLength != nil || hdr.Pattern != ""
	hasSliceValidations := hdr.MaxItems != nil || hdr.MinItems != nil || hdr.UniqueItems
	hasValidations := hasNumberValidation || hasStringValidation || hasSliceValidations || len(hdr.Enum) > 0

	tpe := simpleResolvedType(hdr.Type, hdr.Format, hdr.Items)

	return GenHeader{
		sharedValidations: sharedValidations{
			Required:            true,
			Maximum:             hdr.Maximum,
			ExclusiveMaximum:    hdr.ExclusiveMaximum,
			Minimum:             hdr.Minimum,
			ExclusiveMinimum:    hdr.ExclusiveMinimum,
			MaxLength:           hdr.MaxLength,
			MinLength:           hdr.MinLength,
			Pattern:             hdr.Pattern,
			MaxItems:            hdr.MaxItems,
			MinItems:            hdr.MinItems,
			UniqueItems:         hdr.UniqueItems,
			MultipleOf:          hdr.MultipleOf,
			Enum:                hdr.Enum,
			HasValidations:      hasValidations,
			HasSliceValidations: hasSliceValidations,
		},
		resolvedType: tpe,
		Package:      b.APIPackage,
		ReceiverName: receiver,
		Name:         name,
		Path:         name,
		Description:  hdr.Description,
		Converter:    stringConverters[tpe.GoType],
		Formatter:    stringFormatters[tpe.GoType],
	}
}

func (b codeGenOpBuilder) MakeParameterItem(receiver, paramName, indexVar, path, valueExpression string, resolver *typeResolver, items, parent *spec.Items) (GenItems, error) {
	var res GenItems
	res.resolvedType = simpleResolvedType(items.Type, items.Format, items.Items)
	res.sharedValidations = sharedValidations{
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
	res.Name = paramName
	res.ValueExpression = valueExpression
	res.CollectionFormat = items.CollectionFormat
	res.Converter = stringConverters[res.GoType]
	res.Formatter = stringFormatters[res.GoType]

	if items.Items != nil {
		pi, err := b.MakeParameterItem(receiver, paramName+" "+indexVar, indexVar+"i", "fmt.Sprintf(\"%s.%v\", "+path+", "+indexVar+")", valueExpression+"["+indexVar+"]", resolver, items.Items, items)
		if err != nil {
			return GenItems{}, err
		}
		res.Child = &pi
		pi.Parent = &res
	}

	return res, nil
}

func (b codeGenOpBuilder) MakeParameter(receiver string, resolver *typeResolver, param spec.Parameter) (GenParameter, error) {
	var child *GenItems
	res := GenParameter{
		Name:             param.Name,
		Path:             fmt.Sprintf("%q", param.Name),
		ValueExpression:  fmt.Sprintf("%s.%s", receiver, swag.ToGoName(param.Name)),
		IndexVar:         "i",
		BodyParam:        nil,
		Default:          param.Default,
		Enum:             param.Enum,
		Description:      param.Description,
		ReceiverName:     receiver,
		CollectionFormat: param.CollectionFormat,
		Child:            child,
		Location:         param.In,
	}

	if param.In == "body" {
		sc := schemaGenContext{
			Path:         res.Path,
			Name:         res.Name,
			Receiver:     res.ReceiverName,
			ValueExpr:    res.ValueExpression,
			IndexVar:     res.IndexVar,
			Schema:       *param.Schema,
			Required:     param.Required,
			TypeResolver: resolver,
			Named:        false,
			ExtraSchemas: make(map[string]GenSchema),
		}
		if err := sc.makeGenSchema(); err != nil {
			return GenParameter{}, err
		}

		schema := sc.GenSchema
		if schema.IsAnonymous {
			schema.Name = swag.ToGoName(b.Operation.ID + " Body")
			nm := schema.Name
			schema.GoType = nm
			schema.IsAnonymous = false
			if b.ExtraSchemas == nil {
				b.ExtraSchemas = make(map[string]GenSchema)
			}
			b.ExtraSchemas[nm] = schema
			schema = GenSchema{}
			schema.IsAnonymous = false
			schema.GoType = nm
			schema.SwaggerType = nm
			schema.IsComplexObject = true
		}
		res.Schema = &schema
		res.resolvedType = schema.resolvedType
		res.sharedValidations = schema.sharedValidations

	} else {
		res.resolvedType = simpleResolvedType(param.Type, param.Format, param.Items)
		res.sharedValidations = sharedValidations{
			Required:         param.Required,
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

		if param.Items != nil {
			pi, err := b.MakeParameterItem(receiver, param.Name+" "+res.IndexVar, res.IndexVar+"i", "fmt.Sprintf(\"%s.%v\", "+res.Path+", "+res.IndexVar+")", res.ValueExpression+"["+res.IndexVar+"]", resolver, param.Items, nil)
			if err != nil {
				return GenParameter{}, err
			}
			res.Child = &pi
		}

	}

	hasNumberValidation := param.Maximum != nil || param.Minimum != nil || param.MultipleOf != nil
	hasStringValidation := param.MaxLength != nil || param.MinLength != nil || param.Pattern != ""
	hasSliceValidations := param.MaxItems != nil || param.MinItems != nil || param.UniqueItems
	hasValidations := hasNumberValidation || hasStringValidation || hasSliceValidations || len(param.Enum) > 0

	res.Converter = stringConverters[res.GoType]
	res.Formatter = stringFormatters[res.GoType]
	res.HasValidations = hasValidations
	res.HasSliceValidations = hasSliceValidations
	return res, nil
}

func modelValidations(gs GenSchema) commonValidations {

	return commonValidations{
		propertyDescriptor: propertyDescriptor{
			PropertyName:      swag.ToGoName(gs.Name),
			ParamName:         gs.Name,
			ValueExpression:   gs.ValueExpression,
			IndexVar:          gs.IndexVar,
			Path:              gs.Path,
			IsContainer:       gs.IsArray,
			IsPrimitive:       gs.IsPrimitive,
			IsCustomFormatter: gs.IsCustomFormatter,
			IsMap:             gs.IsMap,
		},
		sharedValidations: sharedValidations{
			Required:         gs.Required,
			Maximum:          gs.Maximum,
			ExclusiveMaximum: gs.ExclusiveMaximum,
			Minimum:          gs.Minimum,
			ExclusiveMinimum: gs.ExclusiveMinimum,
			MaxLength:        gs.MaxLength,
			MinLength:        gs.MinLength,
			Pattern:          gs.Pattern,
			MaxItems:         gs.MaxItems,
			MinItems:         gs.MinItems,
			UniqueItems:      gs.UniqueItems,
			MultipleOf:       gs.MultipleOf,
			Enum:             gs.Enum,
		},
		Type:   gs.GoType,
		Format: gs.SwaggerFormat,
		//Default:          model.Default,
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
		sharedValidations: sharedValidations{
			Required:         param.Required,
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
		},
		Type:    tpe,
		Format:  param.Format,
		Items:   param.Items,
		Default: param.Default,
	}
}

type sharedParam struct {
	genValidations
	propertyDescriptor
}

func paramItemValidations(path, paramName, accessor, indexVar, valueExpression string, items spec.Items) commonValidations {
	tpe := resolveSimpleType(items.Type, items.Format, items.Items)
	_, isPrimitive := primitives[tpe]
	_, isCustomFormatter := customFormatters[tpe]

	shv := sharedValidations{
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
		sharedValidations: shv,

		Type:    tpe,
		Format:  items.Format,
		Items:   items.Items,
		Default: items.Default,
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

// GenResponse represents a response object for code generation
type GenResponse struct {
	Package      string
	ReceiverName string
	Name         string
	Description  string

	IsSuccess bool

	Headers []GenHeader
	Schema  *GenSchema

	Imports        map[string]string
	DefaultImports []string
}

// GenHeader represents a header on a response for code generation
type GenHeader struct {
	resolvedType
	sharedValidations

	Package      string
	ReceiverName string

	Name string
	Path string

	Title       string
	Description string

	Converter string
	Formatter string
}

// GenParameter is used to represent
// a parameter or a header for code generation.
type GenParameter struct {
	resolvedType
	sharedValidations

	Name            string
	Path            string
	ValueExpression string
	IndexVar        string
	ReceiverName    string
	Location        string
	Title           string
	Description     string
	Converter       string
	Formatter       string

	Schema *GenSchema

	CollectionFormat string

	Child  *GenItems
	Parent *GenItems

	BodyParam *GenParameter

	Default interface{}
	Enum    []interface{}
}

// IsQueryParam returns true when this parameter is a query param
func (g *GenParameter) IsQueryParam() bool {
	return g.Location == "query"
}

// IsPathParam returns true when this parameter is a path param
func (g *GenParameter) IsPathParam() bool {
	return g.Location == "path"
}

// IsFormParam returns true when this parameter is a form param
func (g *GenParameter) IsFormParam() bool {
	return g.Location == "formData"
}

// IsHeaderParam returns true when this parameter is a header param
func (g *GenParameter) IsHeaderParam() bool {
	return g.Location == "header"
}

// IsBodyParam returns true when this parameter is a body param
func (g *GenParameter) IsBodyParam() bool {
	return g.Location == "body"
}

// IsFileParam returns true when this parameter is a file param
func (g *GenParameter) IsFileParam() bool {
	return g.SwaggerType == "file"
}

// GenItems represents the collection items for a collection parameter
type GenItems struct {
	sharedValidations
	resolvedType

	Name             string
	Path             string
	ValueExpression  string
	CollectionFormat string
	Child            *GenItems
	Parent           *GenItems
	Converter        string
	Formatter        string

	Location string
}

// GenOperationGroup represents a named (tagged) group of operations
type GenOperationGroup struct {
	Name       string
	Operations []GenOperation

	Summary        string
	Description    string
	Imports        map[string]string
	DefaultImports []string
}

// GenOperation represents an operation for code generation
type GenOperation struct {
	Package      string
	ReceiverName string
	Name         string
	Summary      string
	Description  string

	Imports        map[string]string
	DefaultImports []string

	Authorized bool
	Principal  string

	SuccessResponse *GenResponse
	Responses       map[int]GenResponse
	DefaultResponse *GenResponse

	Params         []GenParameter
	QueryParams    []GenParameter
	PathParams     []GenParameter
	HeaderParams   []GenParameter
	FormParams     []GenParameter
	HasQueryParams bool
	HasFormParams  bool
	HasFileParams  bool
}
