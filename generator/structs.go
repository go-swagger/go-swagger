package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
)

// GenCommon contains common properties needed across
// definitions, app and operations
// TargetImportPath may be used by templates to import other (possibly
// generated) packages in the generation path (e.g. relative to GOPATH).
// TargetImportPath is NOT used by standard templates.
type GenCommon struct {
	Copyright        string
	TargetImportPath string
}

// GenDefinition contains all the properties to generate a
// definition from a swagger spec
type GenDefinition struct {
	GenCommon
	GenSchema
	Package        string
	Imports        map[string]string
	DefaultImports []string
	ExtraSchemas   []GenSchema
	DependsOn      []string
}

// GenSchemaList is a list of schemas for generation.
//
// It can be sorted by name to get a stable struct layout for
// version control and such
type GenSchemaList []GenSchema

func (g GenSchemaList) Len() int           { return len(g) }
func (g GenSchemaList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenSchemaList) Less(i, j int) bool { return g[i].Name < g[j].Name }

// GenSchema contains all the information needed to generate the code
// for a schema
type GenSchema struct {
	resolvedType
	sharedValidations
	Example                 string
	OriginalName            string
	Name                    string
	Suffix                  string
	Path                    string
	ValueExpression         string
	IndexVar                string
	KeyVar                  string
	Title                   string
	Description             string
	Location                string
	ReceiverName            string
	Items                   *GenSchema
	AllowsAdditionalItems   bool
	HasAdditionalItems      bool
	AdditionalItems         *GenSchema
	Object                  *GenSchema
	XMLName                 string
	CustomTag               string
	Properties              GenSchemaList
	AllOf                   []GenSchema
	HasAdditionalProperties bool
	IsAdditionalProperties  bool
	AdditionalProperties    *GenSchema
	ReadOnly                bool
	IsVirtual               bool
	IsBaseType              bool
	HasBaseType             bool
	IsSubType               bool
	IsExported              bool
	DiscriminatorField      string
	DiscriminatorValue      string
	Discriminates           map[string]string
	Parents                 []string
	IncludeValidator        bool
	IncludeModel            bool
	Default                 interface{}
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
	Enum                []interface{}
	ItemsEnum           []interface{}
	HasValidations      bool
	MinItems            *int64
	MaxItems            *int64
	UniqueItems         bool
	HasSliceValidations bool
	NeedsSize           bool
	NeedsValidation     bool
	NeedsRequired       bool
}

// GenResponse represents a response object for code generation
type GenResponse struct {
	Package       string
	ModelsPackage string
	ReceiverName  string
	Name          string
	Description   string

	IsSuccess bool

	Code               int
	Method             string
	Path               string
	Headers            GenHeaders
	Schema             *GenSchema
	AllowsForStreaming bool

	Imports        map[string]string
	DefaultImports []string

	Extensions map[string]interface{}
}

// GenHeader represents a header on a response for code generation
type GenHeader struct {
	resolvedType
	sharedValidations

	Package      string
	ReceiverName string
	IndexVar     string

	ID              string
	Name            string
	Path            string
	ValueExpression string

	Title       string
	Description string
	Default     interface{}
	HasDefault  bool

	CollectionFormat string

	Child  *GenItems
	Parent *GenItems

	Converter string
	Formatter string

	ZeroValue string
}

// GenHeaders is a sorted collection of headers for codegen
type GenHeaders []GenHeader

func (g GenHeaders) Len() int           { return len(g) }
func (g GenHeaders) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenHeaders) Less(i, j int) bool { return g[i].Name < g[j].Name }

// GenParameter is used to represent
// a parameter or a header for code generation.
type GenParameter struct {
	resolvedType
	sharedValidations

	ID              string
	Name            string
	ModelsPackage   string
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

	Default         interface{}
	HasDefault      bool
	ZeroValue       string
	AllowEmptyValue bool

	Extensions map[string]interface{}
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

// GenParameters represents a sorted parameter collection
type GenParameters []GenParameter

func (g GenParameters) Len() int           { return len(g) }
func (g GenParameters) Less(i, j int) bool { return g[i].Name < g[j].Name }
func (g GenParameters) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

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
	IndexVar string
}

// GenOperationGroup represents a named (tagged) group of operations
type GenOperationGroup struct {
	GenCommon
	Name       string
	Operations GenOperations

	Summary        string
	Description    string
	Imports        map[string]string
	DefaultImports []string
	RootPackage    string
	WithContext    bool
}

// GenOperationGroups is a sorted collection of operation groups
type GenOperationGroups []GenOperationGroup

func (g GenOperationGroups) Len() int           { return len(g) }
func (g GenOperationGroups) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenOperationGroups) Less(i, j int) bool { return g[i].Name < g[j].Name }

// GenStatusCodeResponses a container for status code responses
type GenStatusCodeResponses []GenResponse

func (g GenStatusCodeResponses) Len() int           { return len(g) }
func (g GenStatusCodeResponses) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenStatusCodeResponses) Less(i, j int) bool { return g[i].Code < g[j].Code }

// MarshalJSON marshals these responses to json
func (g GenStatusCodeResponses) MarshalJSON() ([]byte, error) {
	if g == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteRune('{')
	for i, v := range g {
		rb, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(fmt.Sprintf("%q:", strconv.Itoa(v.Code)))
		buf.Write(rb)
	}
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON unmarshals this GenStatusCodeResponses from json
func (g *GenStatusCodeResponses) UnmarshalJSON(data []byte) error {
	var dd map[string]GenResponse
	if err := json.Unmarshal(data, &dd); err != nil {
		return err
	}
	var gg GenStatusCodeResponses
	for _, v := range dd {
		gg = append(gg, v)
	}
	sort.Sort(gg)
	*g = gg
	return nil
}

// GenOperation represents an operation for code generation
type GenOperation struct {
	GenCommon
	Package      string
	ReceiverName string
	Name         string
	Summary      string
	Description  string
	Method       string
	Path         string
	BasePath     string
	Tags         []string
	RootPackage  string

	Imports        map[string]string
	DefaultImports []string
	ExtraSchemas   []GenSchema

	Authorized          bool
	Security            []analysis.SecurityRequirement
	SecurityDefinitions map[string]spec.SecurityScheme
	Principal           string

	SuccessResponse  *GenResponse
	SuccessResponses []GenResponse
	Responses        GenStatusCodeResponses
	DefaultResponse  *GenResponse

	Params               GenParameters
	QueryParams          GenParameters
	PathParams           GenParameters
	HeaderParams         GenParameters
	FormParams           GenParameters
	HasQueryParams       bool
	HasFormParams        bool
	HasFormValueParams   bool
	HasFileParams        bool
	HasStreamingResponse bool

	Schemes            []string
	ExtraSchemes       []string
	ProducesMediaTypes []string
	ConsumesMediaTypes []string
	WithContext        bool
	TimeoutName        string

	Extensions map[string]interface{}
}

// GenOperations represents a list of operations to generate
// this implements a sort by operation id
type GenOperations []GenOperation

func (g GenOperations) Len() int           { return len(g) }
func (g GenOperations) Less(i, j int) bool { return g[i].Name < g[j].Name }
func (g GenOperations) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

// GenApp represents all the meta data needed to generate an application
// from a swagger spec
type GenApp struct {
	GenCommon
	APIPackage          string
	Package             string
	ReceiverName        string
	Name                string
	Principal           string
	DefaultConsumes     string
	DefaultProduces     string
	Host                string
	BasePath            string
	Info                *spec.Info
	ExternalDocs        *spec.ExternalDocumentation
	Imports             map[string]string
	DefaultImports      []string
	Schemes             []string
	ExtraSchemes        []string
	Consumes            GenSerGroups
	Produces            GenSerGroups
	SecurityDefinitions []GenSecurityScheme
	Models              []GenDefinition
	Operations          GenOperations
	OperationGroups     GenOperationGroups
	SwaggerJSON         string
	// this is important for when the generated server adds routes
	// ideally this should be removed after we code-generate the router instead of relying on runtime
	// CAUTION: Could be problematic for big specs (might consume large amounts of memory)
	FlatSwaggerJSON string
	ExcludeSpec     bool
	WithContext     bool
	GenOpts         *GenOpts
}

// UseGoStructFlags returns true when no strategy is specified or it is set to "go-flags"
func (g *GenApp) UseGoStructFlags() bool {
	if g.GenOpts == nil {
		return true
	}
	return g.GenOpts.FlagStrategy == "" || g.GenOpts.FlagStrategy == "go-flags"
}

// UsePFlags returns true when the flag strategy is set to pflag
func (g *GenApp) UsePFlags() bool {
	return g.GenOpts != nil && strings.HasPrefix(g.GenOpts.FlagStrategy, "pflag")
}

// UseIntermediateMode for https://wiki.mozilla.org/Security/Server_Side_TLS#Intermediate_compatibility_.28default.29
func (g *GenApp) UseIntermediateMode() bool {
	return g.GenOpts != nil && g.GenOpts.CompatibilityMode == "intermediate"
}

// UseModernMode for https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
func (g *GenApp) UseModernMode() bool {
	return g.GenOpts == nil || g.GenOpts.CompatibilityMode == "" || g.GenOpts.CompatibilityMode == "modern"
}

// GenSerGroups sorted representation of serializer groups
type GenSerGroups []GenSerGroup

func (g GenSerGroups) Len() int           { return len(g) }
func (g GenSerGroups) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenSerGroups) Less(i, j int) bool { return g[i].MediaType < g[j].MediaType }

// GenSerGroup represents a group of serializers, most likely this is a media type to a list of
// prioritized serializers.
type GenSerGroup struct {
	ReceiverName   string
	AppName        string
	Name           string
	MediaType      string
	Implementation string
	AllSerializers GenSerializers
}

// GenSerializers sorted representation of serializers
type GenSerializers []GenSerializer

func (g GenSerializers) Len() int           { return len(g) }
func (g GenSerializers) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenSerializers) Less(i, j int) bool { return g[i].MediaType < g[j].MediaType }

// GenSerializer represents a single serializer for a particular media type
type GenSerializer struct {
	ReceiverName   string
	AppName        string
	Name           string
	MediaType      string
	Implementation string
}

// GenSecuritySchemes sorted representation of serializers
type GenSecuritySchemes []GenSecurityScheme

func (g GenSecuritySchemes) Len() int           { return len(g) }
func (g GenSecuritySchemes) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenSecuritySchemes) Less(i, j int) bool { return g[i].Name < g[j].Name }

// GenSecurityScheme represents a security scheme for code generation
type GenSecurityScheme struct {
	AppName      string
	ID           string
	Name         string
	ReceiverName string
	IsBasicAuth  bool
	IsAPIKeyAuth bool
	IsOAuth2     bool
	Scopes       []string
	Source       string
	Principal    string
}
