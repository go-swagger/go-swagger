package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// QueryParam creates a query parameter
func QueryParam() *Parameter {
	return &Parameter{In: "query"}
}

// HeaderParam creates a header parameter, this is always required by default
func HeaderParam() *Parameter {
	return &Parameter{In: "header", Required: true}
}

// PathParam creates a path parameter, this is always required
func PathParam() *Parameter {
	return &Parameter{In: "path", Required: true}
}

// BodyParam creates a body parameter
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
//
// For more information: http://goo.gl/8us55a#parameterObject
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
