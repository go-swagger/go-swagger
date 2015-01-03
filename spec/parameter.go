package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/util"
)

// QueryParam creates a query parameter
func QueryParam(name string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "query"}}
}

// HeaderParam creates a header parameter, this is always required by default
func HeaderParam(name string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "header", Required: true}}
}

// PathParam creates a path parameter, this is always required
func PathParam(name string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "path", Required: true}}
}

// BodyParam creates a body parameter
func BodyParam(name string, schema *Schema) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "body", Schema: schema}, simpleSchema: simpleSchema{Type: "object"}}
}

// FormDataParam creates a body parameter
func FormDataParam(name string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "formData"}}
}

// FileParam creates a body parameter
func FileParam(name string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name, In: "formData"}, simpleSchema: simpleSchema{Type: "file"}}
}

// SimpleArrayParam creates a param for a simple array (string, int, date etc)
func SimpleArrayParam(name, tpe, fmt string) *Parameter {
	return &Parameter{paramProps: paramProps{Name: name}, simpleSchema: simpleSchema{Type: "array", CollectionFormat: "csv", Items: &Items{simpleSchema: simpleSchema{Type: "string", Format: fmt}}}}
}

type paramProps struct {
	Description string  `json:"description,omitempty"`
	Name        string  `json:"name,omitempty"`
	In          string  `json:"in,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"` // when in == "body"
}

// Parameter a unique parameter is defined by a combination of a [name](#parameterName) and [location](#parameterIn).
//
// There are five possible parameter types.
// * Path - Used together with [Path Templating](#pathTemplating), where the parameter value is actually part of the operation's URL. This does not include the host or base path of the API. For example, in `/items/{itemId}`, the path parameter is `itemId`.
// * Query - Parameters that are appended to the URL. For example, in `/items?id=###`, the query parameter is `id`.
// * Header - Custom headers that are expected as part of the request.
// * Body - The payload that's appended to the HTTP request. Since there can only be one payload, there can only be *one* body parameter. The name of the body parameter has no effect on the parameter itself and is used for documentation purposes only. Since Form parameters are also in the payload, body and form parameters cannot exist together for the same operation.
// * Form - Used to describe the payload of an HTTP request when either `application/x-www-form-urlencoded` or `multipart/form-data` are used as the content type of the request (in Swagger's definition, the [`consumes`](#operationConsumes) property of an operation). This is the only parameter type that can be used to send files, thus supporting the `file` type. Since form parameters are sent in the payload, they cannot be declared together with a body parameter for the same operation. Form parameters have a different format based on the content-type used (for further details, consult http://www.w3.org/TR/html401/interact/forms.html#h-17.13.4):
//   * `application/x-www-form-urlencoded` - Similar to the format of Query parameters but as a payload. For example, `foo=1&bar=swagger` - both `foo` and `bar` are form parameters. This is normally used for simple parameters that are being transferred.
//   * `multipart/form-data` - each parameter takes a section in the payload with an internal header. For example, for the header `Content-Disposition: form-data; name="submit-name"` the name of the parameter is `submit-name`. This type of form parameters is more commonly used for file transfers.
//
// For more information: http://goo.gl/8us55a#parameterObject
type Parameter struct {
	refable
	commonValidations
	simpleSchema
	vendorExtensible
	paramProps
}

// Typed a fluent builder method for the type of parameter
func (p *Parameter) Typed(tpe, format string) *Parameter {
	p.Type = tpe
	p.Format = format
	return p
}

// CollectionOf a fluent builder method for an array parameter
func (p *Parameter) CollectionOf(items *Items, format string) *Parameter {
	p.Type = "array"
	p.Items = items
	p.CollectionFormat = format
	return p
}

// WithDefault sets the default value on this parameter
func (p *Parameter) WithDefault(defaultValue interface{}) *Parameter {
	p.Default = defaultValue
	return p
}

// AsOptional flags this parameter as optional
func (p *Parameter) AsOptional() *Parameter {
	p.Required = false
	return p
}

// AsRequired flags this parameter as required
func (p *Parameter) AsRequired() *Parameter {
	p.Required = true
	return p
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (p *Parameter) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &p.commonValidations); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.refable); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.simpleSchema); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.vendorExtensible); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.paramProps); err != nil {
		return err
	}
	return nil
}

// MarshalJSON converts this items object to JSON
func (p Parameter) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(p.commonValidations)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(p.simpleSchema)
	if err != nil {
		return nil, err
	}
	b3, err := json.Marshal(p.refable)
	if err != nil {
		return nil, err
	}
	b4, err := json.Marshal(p.vendorExtensible)
	if err != nil {
		return nil, err
	}
	b5, err := json.Marshal(p.paramProps)
	if err != nil {
		return nil, err
	}
	return util.ConcatJSON(b3, b1, b2, b4, b5), nil
}
