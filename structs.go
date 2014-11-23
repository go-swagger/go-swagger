package swagger

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

// License information for the exposed API.
type License struct {
	Name string `structs:"name"`
	URL  string `structs:"url"`
}

// Contact information for the exposed API.
type ContactInfo struct {
	Name  string `structs:"name"`
	URL   string `structs:"url"`
	Email string `structs:"email"`
}

func addExtensions(res map[string]interface{}, extensions map[string]interface{}) {
	for k, v := range extensions {
		key := k
		if key != "" {
			if !strings.HasPrefix(key, "x-") {
				key = "x-" + key
			}
			if !structs.IsZero(v) {
				res[key] = structs.Map(v)
			}
		}
	}
}

// The object provides metadata about the API.
// The metadata can be used by the clients if needed, and can be presented in the Swagger-UI for convenience.
type Info struct {
	Extensions     map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Description    string                 `structs:"description,omitempty"`
	Title          string                 `structs:"title,omitempty"`
	TermsOfService string                 `structs:"termsOfService,omitempty"`
	Contact        *ContactInfo           `structs:"contact,omitempty"`
	License        *License               `structs:"license,omitempty"`
	Version        string                 `structs:"version,omitempty"`
}

func (i Info) Map() map[string]interface{} {
	res := structs.Map(i)
	addExtensions(res, i.Extensions)
	return res
}

func (i Info) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Map())
}

func (i Info) MarshalYAML() (interface{}, error) {
	return i.Map(), nil
}

// This is the root document object for the API specification.
// It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
type Swagger struct {
	Consumes            []string               `structs:"consumes,omitempty"`
	Produces            []string               `structs:"produces,omitempty"`
	Schemes             []string               `structs:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Swagger             string                 `structs:"swagger"`
	Info                Info                   `structs:"-"`
	Host                string                 `structs:"host,omitempty"`
	BasePath            string                 `structs:"basePath,omitempty"` // must start with a leading "/"
	Paths               Paths                  `structs:"-"`                  // required
	Definitions         Definitions            `structs:"-"`
	Parameters          []Parameter            `structs:"-"`
	Responses           ResponsesMap           `structs:"-"`
	SecurityDefinitions SecurityDefinitions    `structs:"-"`
	Security            SecurityRequirements   `structs:"security,omitempty"`
	Tags                []Tag                  `structs:"-"`
	ExternalDocs        *ExternalDocumentation `structs:"externalDocs,omitempty"`
}

func (s Swagger) Map() map[string]interface{} {
	res := structs.Map(s)
	res["info"] = s.Info.Map()
	res["paths"] = s.Paths.Map()
	res["responses"] = s.Responses.Map()
	res["definitions"] = s.Definitions.Map()
	res["securityDefinitions"] = s.SecurityDefinitions.Map()

	var params []map[string]interface{}
	for _, param := range s.Parameters {
		params = append(params, param.Map())
	}
	res["parameters"] = params

	var tags []map[string]interface{}
	for _, t := range s.Tags {
		tags = append(tags, t.Map())
	}
	res["tags"] = tags
	return res
}

func (s Swagger) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s Swagger) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}

type SecurityDefinitions map[string]*SecurityScheme

func (s SecurityDefinitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range s {
		res[k] = v.Map()
	}
	return res
}

type Definitions map[string]Schema

func (d Definitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range d {
		res[k] = v.Map()
	}
	return res
}

type ResponsesMap map[string]Response

func (r ResponsesMap) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range r {
		res[k] = v.Map()
	}
	return res
}

type Tag struct {
	Description  string                 `structs:"description,omitempty"`
	Extensions   map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Name         string                 `structs:"name"`
	ExternalDocs *ExternalDocumentation `structs:"externalDocs,omitempty"`
}

func (t Tag) Map() map[string]interface{} {
	res := structs.Map(t)
	addExtensions(res, t.Extensions)
	return res
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Map())
}

func (t Tag) MarshalYAML() (interface{}, error) {
	return t.Map(), nil
}

// Holds the relative paths to the individual endpoints.
// The path is appended to the [`basePath`](#swaggerBasePath) in order
// to construct the full URL.
// The Paths may be empty, due to [ACL constraints](#securityFiltering).
type Paths struct {
	Extensions map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Paths      map[string]PathItem    `structs:"-"` // custom serializer to flatten this, each entry must start with "/"
}

func (p Paths) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range p.Paths {
		key := k
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}
		res[key] = v.Map()
	}
	addExtensions(res, p.Extensions)
	return res
}

func (p Paths) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Map())
}

func (p Paths) MarshalYAML() (interface{}, error) {
	return p.Map(), nil
}

// Describes the operations available on a single path.
// A Path Item may be empty, due to [ACL constraints](#securityFiltering).
// The path itself is still exposed to the documentation viewer but they will
// not know which operations and parameters are available.
type PathItem struct {
	Ref        string                 `structs:"-"`
	Extensions map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Get        *Operation             `structs:"get,omitempty"`
	Put        *Operation             `structs:"put,omitempty"`
	Post       *Operation             `structs:"post,omitempty"`
	Delete     *Operation             `structs:"delete,omitempty"`
	Options    *Operation             `structs:"options,omitempty"`
	Head       *Operation             `structs:"head,omitempty"`
	Patch      *Operation             `structs:"patch,omitempty"`
	Parameters []Parameter            `structs:"-"`
}

func (p PathItem) Map() map[string]interface{} {
	if p.Ref != "" {
		return map[string]interface{}{"$ref": p.Ref}
	}

	res := make(map[string]interface{})
	addOp := func(key string, op *Operation) {
		if op != nil {
			res[key] = op.Map()
		}
	}
	addOp("get", p.Get)
	addOp("put", p.Put)
	addOp("post", p.Post)
	addOp("delete", p.Delete)
	addOp("head", p.Head)
	addOp("options", p.Options)
	addOp("patch", p.Patch)

	var params []map[string]interface{}
	for _, param := range p.Parameters {
		params = append(params, param.Map())
	}
	res["parameters"] = params

	addExtensions(res, p.Extensions)

	return res
}

func (p PathItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Map())
}
func (p PathItem) MarshalYAML() (interface{}, error) {
	return p.Map(), nil
}

// Describes a single API operation on a path.
type Operation struct {
	Description  string                 `structs:"description,omitempty"`
	Extensions   map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Consumes     []string               `structs:"consumes,omitempty"`
	Produces     []string               `structs:"produces,omitempty"`
	Schemes      []string               `structs:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Tags         []string               `structs:"tags,omitempty"`
	Summary      string                 `structs:"summary,omitempty"`
	ExternalDocs *ExternalDocumentation `structs:"externalDocs,omitempty"`
	ID           string                 `structs:"operationId"`
	Deprecated   bool                   `structs:"deprecated,omitempty"`
	Security     []SecurityRequirement  `structs:"security,omitempty"`
	Parameters   []Parameter            `structs:"-"`
	Responses    Responses              `structs:"-"`
}

func (o Operation) Map() map[string]interface{} {
	res := structs.Map(o)
	res["responses"] = o.Responses.Map()
	var params []map[string]interface{}
	for _, param := range o.Parameters {
		params = append(params, param.Map())
	}
	res["parameters"] = params
	addExtensions(res, o.Extensions)
	return res
}

func (o Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Map())
}

func (o Operation) MarshalYAML() (interface{}, error) {
	return o.Map(), nil
}

type Responses struct {
	Extensions          map[string]interface{}
	Default             *Response
	StatusCodeResponses map[int]Response
}

func (r Responses) Map() map[string]interface{} {
	res := make(map[string]interface{})
	if r.Default != nil {
		res["default"] = r.Default.Map()
	}
	for k, v := range r.StatusCodeResponses {
		res[strconv.Itoa(k)] = v.Map()
	}
	addExtensions(res, r.Extensions)
	return res
}

func (r Responses) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Map())
}
func (r Responses) MarshalYAML() (interface{}, error) {
	return r.Map(), nil
}

type Response struct {
	Description string            `structs:"description,omitempty"`
	Ref         string            `structs:"-"`
	Schema      *Schema           `structs:"-"`
	Headers     map[string]Header `structs:"-"`
	Examples    interface{}       `structs:"examples,omitempty"`
}

func (r Response) Map() map[string]interface{} {
	if r.Ref != "" {
		return map[string]interface{}{"$ref": r.Ref}
	}

	res := structs.Map(r)

	if r.Schema != nil {
		res["schema"] = r.Schema.Map()
	}

	if len(r.Headers) > 0 {
		headers := make(map[string]map[string]interface{}, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Map()
		}
		res["headers"] = headers
	}

	return res
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Map())
}

func (r Response) MarshalYAML() (interface{}, error) {
	return r.Map(), nil
}

type Header struct {
	Description      string        `structs:"description,omitempty"`
	Maximum          float64       `structs:"maximum,omitempty"`
	ExclusiveMaximum bool          `structs:"exclusiveMaximum,omitempty"`
	Minimum          float64       `structs:"minimum,omitempty"`
	ExclusiveMinimum bool          `structs:"exclusiveMinimum,omitempty"`
	MaxLength        int64         `structs:"maxLength,omitempty"`
	MinLength        int64         `structs:"minLength,omitempty"`
	Pattern          string        `structs:"pattern,omitempty"`
	MaxItems         int64         `structs:"maxItems,omitempty"`
	MinItems         int64         `structs:"minItems,omitempty"`
	UniqueItems      bool          `structs:"uniqueItems,omitempty"`
	MultipleOf       float64       `structs:"multipleOf,omitempty"`
	Enum             []interface{} `structs:"enum,omitempty"`
	Type             string        `structs:"type,omitempty"`
	Format           string        `structs:"format,omitempty"`
	Default          interface{}   `structs:"default,omitempty"`
	Items            *Items        `structs:"-"`
}

func (h Header) Map() map[string]interface{} {
	res := structs.Map(h)
	if h.Items != nil {
		res["items"] = h.Items.Map()
	}
	return res
}

func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Map())
}

func (h Header) MarshalYAML() (interface{}, error) {
	return h.Map(), nil
}

type ExternalDocumentation struct {
	Description string `structs:"description,omitempty"`
	URL         string `structs:"url"`
}

type SecurityRequirement map[string][]string
type SecurityRequirements []SecurityRequirement

func BasicAuth() *SecurityScheme {
	return &SecurityScheme{Type: "basic"}
}

func ApiKeyAuth(fieldName, valueSource string) *SecurityScheme {
	return &SecurityScheme{Type: "apiKey", Name: fieldName, In: valueSource}
}

func OAuth2Implicit(authorizationURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:             "oauth2",
		Flow:             "implicit",
		AuthorizationURL: authorizationURL,
	}
}

func OAuth2Password(tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:     "oauth2",
		Flow:     "password",
		TokenURL: tokenURL,
	}
}

func OAuth2Application(tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:     "oauth2",
		Flow:     "application",
		TokenURL: tokenURL,
	}
}

func OAuth2AccessToken(authorizationURL, tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:             "oauth2",
		Flow:             "accessCode",
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
	}
}

type SecurityScheme struct {
	Description      string                 `structs:"description,omitempty"`
	Extensions       map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Type             string                 `structs:"type"`
	Name             string                 `structs:"name,omitempty"`             // api key
	In               string                 `structs:"in,omitempty"`               // api key
	Flow             string                 `structs:"flow,omitempty"`             // oauth2
	AuthorizationURL string                 `structs:"authorizationUrl,omitempty"` // oauth2
	TokenURL         string                 `structs:"tokenUrl,omitempty"`         // oauth2
	Scopes           map[string]string      `structs:"scopes,omitempty"`           // oauth2
}

func (s *SecurityScheme) AddScope(scope, description string) {
	if s.Scopes == nil {
		s.Scopes = make(map[string]string)
	}
	s.Scopes[scope] = description
}

func (s SecurityScheme) Map() map[string]interface{} {
	res := structs.Map(s)
	addExtensions(res, s.Extensions)
	return res
}

func (s SecurityScheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s SecurityScheme) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}

func QueryParam() *Parameter {
	return &Parameter{In: "query"}
}

func HeaderParam() *Parameter {
	return &Parameter{In: "header", Required: true}
}

func PathParam() *Parameter {
	return &Parameter{In: "path", Required: true}
}

func BodyParam() *Parameter {
	return &Parameter{In: "body"}
}

// A unique parameter is defined by a combination of a [name](#parameterName) and [location](#parameterIn).
//
// There are five possible parameter types.
// * Path - Used together with [Path Templating](#pathTemplating), where the parameter value is actually part of the operation's URL. This does not include the host or base path of the API. For example, in `/items/{itemId}`, the path parameter is `itemId`.
// * Query - Parameters that are appended to the URL. For example, in `/items?id=###`, the query parameter is `id`.
// * Header - Custom headers that are expected as part of the request.
// * Body - The payload that's appended to the HTTP request. Since there can only be one payload, there can only be *one* body parameter. The name of the body parameter has no effect on the parameter itself and is used for documentation purposes only. Since Form parameters are also in the payload, body and form parameters cannot exist together for the same operation.
// * Form - Used to describe the payload of an HTTP request when either `application/x-www-form-urlencoded` or `multipart/form-data` are used as the content type of the request (in Swagger's definition, the [`consumes`](#operationConsumes) property of an operation). This is the only parameter type that can be used to send files, thus supporting the `file` type. Since form parameters are sent in the payload, they cannot be declared together with a body parameter for the same operation. Form parameters have a different format based on the content-type used (for further details, consult http://www.w3.org/TR/html401/interact/forms.html#h-17.13.4):
//   * `application/x-www-form-urlencoded` - Similar to the format of Query parameters but as a payload. For example, `foo=1&bar=swagger` - both `foo` and `bar` are form parameters. This is normally used for simple parameters that are being transferred.
//   * `multipart/form-data` - each parameter takes a section in the payload with an internal header. For example, for the header `Content-Disposition: form-data; name="submit-name"` the name of the parameter is `submit-name`. This type of form parameters is more commonly used for file transfers.
type Parameter struct {
	Description      string                 `structs:"description,omitempty"`
	Items            *Items                 `structs:"-"`
	Extensions       map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Ref              string                 `structs:"-"`
	Maximum          float64                `structs:"maximum,omitempty"`
	ExclusiveMaximum bool                   `structs:"exclusiveMaximum,omitempty"`
	Minimum          float64                `structs:"minimum,omitempty"`
	ExclusiveMinimum bool                   `structs:"exclusiveMinimum,omitempty"`
	MaxLength        int64                  `structs:"maxLength,omitempty"`
	MinLength        int64                  `structs:"minLength,omitempty"`
	Pattern          string                 `structs:"pattern,omitempty"`
	MaxItems         int64                  `structs:"maxItems,omitempty"`
	MinItems         int64                  `structs:"minItems,omitempty"`
	UniqueItems      bool                   `structs:"uniqueItems,omitempty"`
	MultipleOf       float64                `structs:"multipleOf,omitempty"`
	Enum             []interface{}          `structs:"enum,omitempty"`
	Type             string                 `structs:"type,omitempty"`
	Format           string                 `structs:"format,omitempty"`
	Name             string                 `structs:"name,omitempty"`
	In               string                 `structs:"in,omitempty"`
	Required         bool                   `structs:"required"`
	Schema           *Schema                `structs:"-"` // when in == "body"
	CollectionFormat string                 `structs:"collectionFormat,omitempty"`
	Default          interface{}            `structs:"default,omitempty"`
}

func (p Parameter) Map() map[string]interface{} {
	if p.Ref != "" {
		return map[string]interface{}{"$ref": p.Ref}
	}
	res := structs.Map(p)
	if p.Schema != nil {
		res["schema"] = p.Schema.Map()
	}
	if p.Items != nil {
		res["items"] = p.Items.Map()
	}
	addExtensions(res, p.Extensions)
	return res
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Map())
}

func (p Parameter) MarshalYAML() (interface{}, error) {
	return p.Map(), nil
}

type Items struct {
	Ref              string        `structs:"-"`
	Type             string        `structs:"type,omitempty"`
	Format           string        `structs:"format,omitempty"`
	Items            *Items        `structs:"-"`
	CollectionFormat string        `structs:"collectionFormat,omitempty"`
	Default          interface{}   `structs:"default,omitempty"`
	Maximum          float64       `structs:"maximum,omitempty"`
	ExclusiveMaximum bool          `structs:"exclusiveMaximum,omitempty"`
	Minimum          float64       `structs:"minimum,omitempty"`
	ExclusiveMinimum bool          `structs:"exclusiveMinimum,omitempty"`
	MaxLength        int64         `structs:"maxLength,omitempty"`
	MinLength        int64         `structs:"minLength,omitempty"`
	Pattern          string        `structs:"pattern,omitempty"`
	MaxItems         int64         `structs:"maxItems,omitempty"`
	MinItems         int64         `structs:"minItems,omitempty"`
	UniqueItems      bool          `structs:"uniqueItems,omitempty"`
	MultipleOf       float64       `structs:"multipleOf,omitempty"`
	Enum             []interface{} `structs:"enum,omitempty"`
}

func (i Items) Map() map[string]interface{} {
	if i.Ref != "" {
		return map[string]interface{}{"$ref": i.Ref}
	}
	res := structs.Map(i)
	if i.Items != nil {
		res["items"] = i.Items.Map()
	}
	return res
}

func (i Items) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Map())
}

func (i Items) MarshalYAML() (interface{}, error) {
	return i.Map(), nil
}

func BooleanProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "boolean"}}
}

func StringProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}}
}
func Float64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "double"}
}

func Float32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "float"}
}

func Int32Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int32"}
}

func Int64Property() *Schema {
	return &Schema{Type: &StringOrArray{Single: "number"}, Format: "int64"}
}

func DateProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date"}
}
func DateTimeProperty() *Schema {
	return &Schema{Type: &StringOrArray{Single: "string"}, Format: "date-time"}
}
func MapProperty(property *Schema) *Schema {
	return &Schema{Type: &StringOrArray{Single: "object"}, AdditionalProperties: property}
}
func RefProperty(name string) *Schema {
	return &Schema{Ref: name}
}

func ArrayProperty(items *Schema) *Schema {
	return &Schema{Items: &SchemaOrArray{Single: items}, Type: &StringOrArray{Single: "array"}}
}

type Schema struct {
	Ref                  string                 `structs:"-"`
	Description          string                 `structs:"description,omitempty"`
	Maximum              float64                `structs:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `structs:"exclusiveMaximum,omitempty"`
	Minimum              float64                `structs:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `structs:"exclusiveMinimum,omitempty"`
	MaxLength            int64                  `structs:"maxLength,omitempty"`
	MinLength            int64                  `structs:"minLength,omitempty"`
	Pattern              string                 `structs:"pattern,omitempty"`
	MaxItems             int64                  `structs:"maxItems,omitempty"`
	MinItems             int64                  `structs:"minItems,omitempty"`
	UniqueItems          bool                   `structs:"uniqueItems,omitempty"`
	MultipleOf           float64                `structs:"multipleOf,omitempty"`
	Enum                 []interface{}          `structs:"enum,omitempty"`
	Type                 *StringOrArray         `structs:"-"`
	Format               string                 `structs:"format,omitempty"`
	Title                string                 `structs:"title,omitempty"`
	Default              interface{}            `structs:"default,omitempty"`
	MaxProperties        int64                  `structs:"maxProperties,omitempty"`
	MinProperties        int64                  `structs:"minProperties,omitempty"`
	Required             []string               `structs:"required,omitempty"`
	Items                *SchemaOrArray         `structs:"-"`
	AllOf                []Schema               `structs:"-"`
	Properties           map[string]Schema      `structs:"-"`
	Discriminator        string                 `structs:"discriminator,omitempty"`
	ReadOnly             bool                   `structs:"readOnly,omitempty"`
	XML                  *XMLObject             `structs:"xml,omitempty"`
	ExternalDocs         *ExternalDocumentation `structs:"externalDocs,omitempty"`
	Example              interface{}            `structs:"example,omitempty"`
	AdditionalProperties *Schema                `structs:"-"`
}

func (s Schema) Map() map[string]interface{} {
	if s.Ref != "" {
		return map[string]interface{}{"$ref": s.Ref}
	}
	res := structs.Map(s)

	if len(s.AllOf) > 0 {
		var ser []map[string]interface{}
		for _, sch := range s.AllOf {
			ser = append(ser, sch.Map())
		}
		res["allOf"] = ser
	}

	if len(s.Properties) > 0 {
		ser := make(map[string]interface{})
		for k, v := range s.Properties {
			ser[k] = v.Map()
		}
		res["properties"] = ser
	}
	if s.AdditionalProperties != nil {
		res["additionalProperties"] = s.AdditionalProperties.Map()
	}

	if s.Type != nil {
		var value interface{} = s.Type.Multi
		if s.Type.Single != "" && len(s.Type.Multi) == 0 {
			value = s.Type.Single
		}
		res["type"] = value
	}

	if s.Items != nil {
		var value interface{} = s.Items.Multi
		if len(s.Items.Multi) == 0 && s.Items.Single != nil {
			value = s.Items.Single
		}
		res["items"] = value
	}

	return res
}

func (s Schema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s Schema) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}

type XMLObject struct {
	Name      string `structs:"name,omitempty"`
	Namespace string `structs:"namespace,omitempty"`
	Prefix    string `structs:"prefix,omitempty"`
	Attribute bool   `structs:"attribute,omitempty"`
	Wrapped   bool   `structs:"wrapped,omitempty"`
}

type StringOrArray struct {
	Single string
	Multi  []string
}

func (s StringOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != "" {
		return s.Single, nil
	}
	return s.Multi, nil
}

func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != "" {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

type SchemaOrArray struct {
	Single *Schema
	Multi  []Schema
}

func (s SchemaOrArray) MarshalYAML() (interface{}, error) {
	if s.Single != nil {
		return s.Single, nil
	}
	return s.Multi, nil
}

func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != nil {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}
