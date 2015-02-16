package generator

import "github.com/casualjim/go-swagger/spec"

// Operation represents an operation for codegen
type Operation struct {
	HasParams            bool                `json:"hasParams,omitempty"`
	HasQueryParams       bool                `json:"hasQueryParams,omitempty"`
	ReturnsPrimitive     bool                `json:"returnTypeIsPrimitive,omitempty"`
	ReturnsSimple        bool                `json:"returnSimpleType,omitempty"`
	SubResourceOperation bool                `json:"subresourceOperation,omitempty"`
	Path                 string              `json:"path,omitempty"`
	OperationID          string              `json:"operationId,omitempty"`
	ReturnType           string              `json:"returnType,omitempty"`
	HTTPMethod           string              `json:"httpMethod,omitempty"`
	ReturnBaseType       string              `json:"returnBaseType,omitempty"`
	ReturnContainer      string              `json:"returnContainer,omitempty"`
	Summary              string              `json:"summary,omitempty"`
	Notes                string              `json:"notes,omitempty"`
	BaseName             string              `json:"baseName,omitempty"`
	DefaultResponse      string              `json:"defaultResponse,omitempty"`
	Consumes             []map[string]string `json:"consumes,omitempty"`
	Produces             []map[string]string `json:"produces,omitempty"`
	BodyParam            *Parameter          `json:"bodyParam,omitempty"`
	AllParams            []Parameter         `json:"allParams,omitempty"`
	BodyParams           []Parameter         `json:"bodyParams,omitempty"`
	PathParams           []Parameter         `json:"pathParams,omitempty"`
	QueryParams          []Parameter         `json:"queryParams,omitempty"`
	HeaderParams         []Parameter         `json:"headerParams,omitempty"`
	FormParams           []Parameter         `json:"formParams,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
	Responses            []string            `json:"responses,omitempty"`
	ResponseHeaders      []Property          `json:"responseHeaders,omitempty"`
	Imports              []string            `json:"imports,omitempty"`
	Examples             []map[string]string `json:"examples,omitempty"`
	ExternalDocs         string              `json:"externalDocs,omitempty"`
}

// Response represents a response for codegen
type Response struct {
	Code     string              `json:"code,omitempty"`
	Message  string              `json:"message,omitempty"`
	Examples []map[string]string `json:"examples,omitempty"`
	Schema   interface{}         `json:"schema,omitempty"`
}

// Parameter represents a parameter for codegen
// used to provide data to the mustache template
// the properties use the json flags because we convert it to a json map before
// handing it over to the mustache engine
type Parameter struct {
	HasMore          bool    `json:"hasMore,omitempty"`
	IsContainer      bool    `json:"isContainer,omitempty"`
	SecondaryParam   bool    `json:"secondaryParam,omitempty"`
	BaseName         string  `json:"baseName,omitempty"`
	ParamName        string  `json:"paramName,omitempty"`
	PropertyName     string  `json:"propertyName,omitempty"`
	VarName          string  `json:"varName,omitempty"`
	DataType         string  `json:"dataType,omitempty"`
	CollectionFormat string  `json:"collectionFormat,omitempty"`
	Description      string  `json:"description,omitempty"`
	BaseType         string  `json:"baseType,omitempty"`
	IsQueryParam     bool    `json:"isQueryParam,omitempty"`
	IsPathParam      bool    `json:"isPathParam,omitempty"`
	IsHeaderParam    bool    `json:"isHeaderParam,omitempty"`
	IsCookieParam    bool    `json:"isCookieParam,omitempty"`
	IsBodyParam      bool    `json:"isBodyParam,omitempty"`
	ReaderFunc       string  `json:"reader,omitempty"`
	Minimum          float64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool    `json:"exclusiveMinimum,omitempty"`
	Maximum          float64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool    `json:"exclusiveMaximum,omitempty"`
	MultipleOf       float64 `json:"multipleOf,omitempty"`
	DefaultValue     string  `json:"defaultValue,omitempty"` // contains fmt.Sprintf("%#v", schema.Default)
	/*
	 * Determines whether this parameter is mandatory. If the parameter is in "path",
	 * this property is required and its value MUST be true. Otherwise, the property
	 * MAY be included and its default value is false.
	 */
	Required      bool            `json:"required,omitempty"`
	Enum          string          `json:"enum,omitempty"` // contains fmt.Sprintf("%#v", schema.Enum)
	MinLength     int64           `json:"minLength,omitempty"`
	MaxLength     int64           `json:"maxLength,omitempty"`
	Pattern       string          `json:"pattern,omitempty"`
	ParamLocation string          `json:"paramLocation,omitempty"`
	ParamData     string          `json:"paramData,omitempty"`
	MinItems      int64           `json:"minItems,omitempty"`
	MaxItems      int64           `json:"maxItems,omitempty"`
	UniqueItems   bool            `json:"uniqueItems,omitempty"`
	IndexVar      string          `json:"indexVar,omitempty"`
	Child         []ParameterItem `json:"child,omitempty"`
}

type ParameterItem struct {
	Enum               string          `json:"enum,omitempty"` // contains fmt.Sprintf("%#v", schema.Enum)
	MinLength          int64           `json:"minLength,omitempty"`
	MaxLength          int64           `json:"maxLength,omitempty"`
	Pattern            string          `json:"pattern,omitempty"`
	ParamLocation      string          `json:"paramLocation,omitempty"`
	ParamData          string          `json:"paramData,omitempty"`
	MinItems           int64           `json:"minItems,omitempty"`
	MaxItems           int64           `json:"maxItems,omitempty"`
	UniqueItems        bool            `json:"uniqueItems,omitempty"`
	IndexVar           string          `json:"indexVar,omitempty"`
	IsContainer        bool            `json:"isContainer,omitempty"`
	ParentPropertyName string          `json:"parentPropertyName,omitempty"`
	Child              []ParameterItem `json:"child,omitempty"`
	PropertyName       string          `json:"propertyName,omitempty"`
	VarName            string          `json:"varName,omitempty"`
}

// Model represents a codegen model
type Model struct {
	Parent        string                      `json:"parent,omitempty"`
	Name          string                      `json:"name,omitempty"`
	ClassName     string                      `json:"classname,omitempty"`
	Description   string                      `json:"description,omitempty"`
	ClassVarName  string                      `json:"classVarName,omitempty"`
	ModelJSON     string                      `json:"modelJson,omitempty"`
	DefaultValue  string                      `json:"defaultValue,omitempty"`
	Vars          []Property                  `json:"vars,omitempty"`
	Imports       []string                    `json:"imports,omitempty"`
	HasVars       bool                        `json:"hasVars,omitempty"`
	EmptyVars     bool                        `json:"emptyVars,omitempty"`
	HasMoreModels bool                        `json:"hasMoreModels,omitempty"`
	ExternalDocs  *spec.ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Property represents a property of a codegen model
type Property struct {
	BaseName         string                 `json:"baseName,omitempty"`
	ComplexType      string                 `json:"complexType,omitempty"`
	Getter           string                 `json:"getter,omitempty"`
	Setter           string                 `json:"setter,omitempty"`
	Description      string                 `json:"description,omitempty"`
	DataType         string                 `json:"datatype,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Min              string                 `json:"min,omitempty"`
	Max              string                 `json:"max,omitempty"`
	DefaultValue     string                 `json:"defaultValue,omitempty"`
	BaseType         string                 `json:"baseType,omitempty"`
	ContainerType    string                 `json:"containerType,omitempty"`
	MaxLength        int64                  `json:"maxLength,omitempty"`
	MinLength        int64                  `json:"minLength,omitempty"`
	Pattern          string                 `json:"pattern,omitempty"`
	Example          string                 `json:"example,omitempty"`
	Minimum          float64                `json:"minimum,omitempty"`
	Maximum          float64                `json:"maximum,omitempty"`
	ExclusiveMinimum bool                   `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool                   `json:"exclusiveMaximum,omitempty"`
	HasMore          bool                   `json:"hasMore,omitempty"`
	Required         bool                   `json:"required,omitempty"`
	SecondaryParam   bool                   `json:"secondaryParam,omitempty"`
	IsPrimitiveType  bool                   `json:"isPrimitiveType,omitempty"`
	IsContainer      bool                   `json:"isContainer,omitempty"`
	IsNotContainer   bool                   `json:"isNotContainer,omitempty"`
	IsEnum           bool                   `json:"isEnum,omitempty"`
	Enum             []string               `json:"enum,omitempty"`
	AllowableValues  map[string]interface{} `json:"allowableValues,omitempty"`
}

// SupportingFile represents an extra file to be generated.
// this is helpful for generating build scripts and facades
type SupportingFile struct {
	TemplateFile    string `json:"templateFile,omitempty"`
	Folder          string `json:"folder,omitempty"`
	DestinationFile string `json:"destinationFilename,omitempty"`
}
