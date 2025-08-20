package oas

import (
	"context"
)

// Target-neutral writer for building a spec from scanner findings.
type SpecWriter interface {
	BeginDoc(ctx context.Context) error
	SetInfo(title, version, description string)
	AddServer(url, description string)
	EnsurePathOperation(path, method string) OperationBuilder
	DefineSchema(name string, schema *Schema) Ref // returns #/components/schemas/{name}
	DefineParameter(name string, p *Parameter) Ref
	DefineResponse(name string, r *Response) Ref
	DefineRequestBody(name string, rb *RequestBody) Ref
	DefineHeader(name string, h *Header) Ref
	DefineSecurityScheme(name string, s *SecurityScheme) Ref
	FinishDoc(ctx context.Context) (any, error) // returns the concrete document (e.g., *openapi3.T)
}

// OperationBuilder abstracts per-operation authoring (target-specific under the hood).
type OperationBuilder interface {
	SetOperationID(id string)
	SetSummary(summary string)
	SetDescription(desc string)
	SetTags(tags ...string)
	AddParameterRef(ref Ref)
	// Adds/merges requestBody content; mediaType like "application/json".
	SetRequestBody(required bool, mediaType string, schemaRef Ref)
	AddResponse(status string, description string, mediaType string, schemaRef Ref)
	// Optional: headers/examples/links/callbacks per operation.
}

// Minimal portable schema-ish types.
// In production, your OAS31Writer would adapt from internal schema model → OAS 3.1 JSON Schema.
type Schema struct {
	Type        any                 // string or []string (to allow ["string","null"])
	Properties  map[string]*Schema  // nested schemas
	Items       *Schema
	Required    []string
	AllOf       []*Schema
	AnyOf       []*Schema
	OneOf       []*Schema
	Description string
	Format      string
	Enum        []any
	Const       any
	Minimum     *float64
	Maximum     *float64
	MinLength   *int
	MaxLength   *int
	Pattern     string
}

type Parameter struct {
	Name        string
	In          string // path|query|header|cookie
	Required    bool
	Description string
	SchemaRef   Ref
	Style       string
	Explode     *bool
}

type Response struct {
	Description string
	// content map[mediaType]Ref to schema or media type object (simplified here)
}

type RequestBody struct {
	Description string
	Required    bool
	// content map[mediaType]Ref to schema (simplified)
}

type Header struct {
	Description string
	SchemaRef   Ref
	Style       string
	Explode     *bool
}

type SecurityScheme struct {
	Type             string // http|apiKey|oauth2|openIdConnect
	Scheme           string // bearer|basic
	BearerFormat     string
	In               string // header|query|cookie (apiKey)
	Name             string // header/query/cookie key name
	AuthorizationURL string
	TokenURL         string
	OpenIDConnectURL string
	Scopes           map[string]string // oauth2
}

type Ref struct {
	// Abstract reference (e.g., "#/components/schemas/Foo")
	Ref string
}

// -------- Minimal OAS3.1 Writer Skeleton --------

type OpenAPI31Writer struct {
	doc *OAS31Doc // your internal representation; adapt to your chosen backing
}

type OAS31Doc struct {
	OpenAPI         string
	Info            Info
	Servers         []Server
	Paths           map[string]*PathItem
	Components      Components
	JSONSchemaDial  string
}

type Info struct {
	Title       string
	Version     string
	Description string
}
type Server struct {
	URL         string
	Description string
}
type PathItem struct {
	Operations map[string]*Operation // "get","post",...
}
type Operation struct {
	OperationID string
	Summary     string
	Description string
	Tags        []string
	Parameters  []Ref
	RequestBody *RequestBodyHolder
	Responses   map[string]*ResponseHolder
}
type RequestBodyHolder struct {
	Required bool
	Content  map[string]Ref // mediaType → schemaRef
}
type ResponseHolder struct {
	Description string
	Content     map[string]Ref
}
type Components struct {
	Schemas         map[string]*Schema
	Parameters      map[string]*Parameter
	Responses       map[string]*Response
	RequestBodies   map[string]*RequestBody
	Headers         map[string]*Header
	SecuritySchemes map[string]*SecurityScheme
}

func NewOpenAPI31Writer() *OpenAPI31Writer {
	return &OpenAPI31Writer{
		doc: &OAS31Doc{
			OpenAPI: "3.1.0",
			Info:    Info{},
			Servers: []Server{},
			Paths:   map[string]*PathItem{},
			Components: Components{
				Schemas:         map[string]*Schema{},
				Parameters:      map[string]*Parameter{},
				Responses:       map[string]*Response{},
				RequestBodies:   map[string]*RequestBody{},
				Headers:         map[string]*Header{},
				SecuritySchemes: map[string]*SecurityScheme{},
			},
			JSONSchemaDial: "https://json-schema.org/draft/2020-12/schema",
		},
	}
}

func (w *OpenAPI31Writer) BeginDoc(ctx context.Context) error { return nil }
func (w *OpenAPI31Writer) SetInfo(title, version, description string) {
	w.doc.Info = Info{Title: title, Version: version, Description: description}
}
func (w *OpenAPI31Writer) AddServer(url, description string) {
	w.doc.Servers = append(w.doc.Servers, Server{URL: url, Description: description})
}
func (w *OpenAPI31Writer) EnsurePathOperation(path, method string) OperationBuilder {
	pi := w.doc.Paths[path]
	if pi == nil {
		pi = &PathItem{Operations: map[string]*Operation{}}
		w.doc.Paths[path] = pi
	}
	op := pi.Operations[method]
	if op == nil {
		op = &Operation{
			Parameters: []Ref{},
			Responses:  map[string]*ResponseHolder{},
		}
		pi.Operations[method] = op
	}
	return &opBuilder{w: w, op: op}
}
func (w *OpenAPI31Writer) DefineSchema(name string, schema *Schema) Ref {
	w.doc.Components.Schemas[name] = schema
	return Ref{Ref: "#/components/schemas/" + name}
}
func (w *OpenAPI31Writer) DefineParameter(name string, p *Parameter) Ref {
	w.doc.Components.Parameters[name] = p
	return Ref{Ref: "#/components/parameters/" + name}
}
func (w *OpenAPI31Writer) DefineResponse(name string, r *Response) Ref {
	w.doc.Components.Responses[name] = r
	return Ref{Ref: "#/components/responses/" + name}
}
func (w *OpenAPI31Writer) DefineRequestBody(name string, rb *RequestBody) Ref {
	w.doc.Components.RequestBodies[name] = rb
	return Ref{Ref: "#/components/requestBodies/" + name}
}
func (w *OpenAPI31Writer) DefineHeader(name string, h *Header) Ref {
	w.doc.Components.Headers[name] = h
	return Ref{Ref: "#/components/headers/" + name}
}
func (w *OpenAPI31Writer) DefineSecurityScheme(name string, s *SecurityScheme) Ref {
	w.doc.Components.SecuritySchemes[name] = s
	return Ref{Ref: "#/components/securitySchemes/" + name}
}
func (w *OpenAPI31Writer) FinishDoc(ctx context.Context) (any, error) {
	// Return your doc; caller serializes to JSON/YAML.
	return w.doc, nil
}

type opBuilder struct {
	w  *OpenAPI31Writer
	op *Operation
}

func (b *opBuilder) SetOperationID(id string)         { b.op.OperationID = id }
func (b *opBuilder) SetSummary(summary string)        { b.op.Summary = summary }
func (b *opBuilder) SetDescription(desc string)       { b.op.Description = desc }
func (b *opBuilder) SetTags(tags ...string)           { b.op.Tags = append([]string{}, tags...) }
func (b *opBuilder) AddParameterRef(ref Ref)          { b.op.Parameters = append(b.op.Parameters, ref) }
func (b *opBuilder) SetRequestBody(required bool, mediaType string, schemaRef Ref) {
	if b.op.RequestBody == nil {
		b.op.RequestBody = &RequestBodyHolder{Required: required, Content: map[string]Ref{}}
	}
	b.op.RequestBody.Required = required
	b.op.RequestBody.Content[mediaType] = schemaRef
}
func (b *opBuilder) AddResponse(status, description, mediaType string, schemaRef Ref) {
	rh := b.op.Responses[status]
	if rh == nil {
		rh = &ResponseHolder{Description: description, Content: map[string]Ref{}}
		b.op.Responses[status] = rh
	}
	rh.Description = description
	if mediaType != "" && schemaRef.Ref != "" {
		rh.Content[mediaType] = schemaRef
	}
}
