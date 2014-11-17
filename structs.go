package swagger

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

type License struct {
	Name string `structs:"name"`
	URL  string `structs:"url"`
}

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

type SecurityDefinitions map[string]SecurityScheme

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
	Schema      *Schema           `structs:"schema,omitempty"`
	Headers     map[string]Header `structs:"headers,omitempty"`
	Examples    interface{}       `structs:"examples,omitempty"`
}

func (r Response) Map() map[string]interface{} {
	if r.Ref != "" {
		return map[string]interface{}{"$ref": r.Ref}
	}

	return structs.Map(r)
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
	Items            *Items        `structs:"items,omitempty"`
}

type ExternalDocumentation struct {
	Description string `structs:"description,omitempty"`
	URL         string `structs:"url"`
}

type SecurityRequirement map[string][]string
type SecurityRequirements []SecurityRequirement

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

type Parameter struct {
	Description      string                 `structs:"description,omitempty"`
	Items            *Items                 `structs:"items,omitempty"`
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
	Required         bool                   `structs:"required,omitempty"`
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
	Type             string        `structs:"type,omitempty"`
	Format           string        `structs:"format,omitempty"`
	Items            *Items        `structs:"items,omitempty"`
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

type Schema struct {
	Ref              string                 `structs:"-"`
	Description      string                 `structs:"description,omitempty"`
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
	Type             *StringOrArray         `structs:"-"`
	Format           string                 `structs:"format,omitempty"`
	Title            string                 `structs:"title,omitempty"`
	Default          interface{}            `structs:"default,omitempty"`
	MaxProperties    int64                  `structs:"maxProperties,omitempty"`
	MinProperties    int64                  `structs:"minProperties,omitempty"`
	Required         bool                   `structs:"required,omitempty"`
	Items            *SchemaOrArray         `structs:"-"`
	AllOf            []Schema               `structs:"-"`
	Properties       map[string]Schema      `structs:"-"`
	Discriminator    string                 `structs:"discriminator,omitempty"`
	ReadOnly         bool                   `structs:"readOnly,omitempty"`
	XML              *XMLObject             `structs:"xml,omitempty"`
	ExternalDocs     *ExternalDocumentation `structs:"externalDocs,omitempty"`
	Example          interface{}            `structs:"example,omitempty"`
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
		var ser map[string]interface{}
		for k, v := range s.Properties {
			ser[k] = v.Map()
		}
		res["properties"] = ser
	}

	if s.Type != nil {
		var value interface{} = s.Type.Multi
		if s.Type.Single != "" && len(s.Type.Multi) == 0 {
			value = s.Type.Single
		}
		res["type"] = value
	}

	if s.Items != nil {
		var value interface{} = s.Type.Multi
		if len(s.Type.Multi) == 0 && s.Type.Single != "" {
			value = s.Type.Single
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
