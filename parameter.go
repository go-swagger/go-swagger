package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
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
	Description      string                 `swagger:"description,omitempty"`
	Items            *Items                 `swagger:"items,omitempty"`
	Extensions       map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Ref              string                 `swagger:"-"`
	Maximum          float64                `swagger:"maximum,omitempty"`
	ExclusiveMaximum bool                   `swagger:"exclusiveMaximum,omitempty"`
	Minimum          float64                `swagger:"minimum,omitempty"`
	ExclusiveMinimum bool                   `swagger:"exclusiveMinimum,omitempty"`
	MaxLength        int64                  `swagger:"maxLength,omitempty"`
	MinLength        int64                  `swagger:"minLength,omitempty"`
	Pattern          string                 `swagger:"pattern,omitempty"`
	MaxItems         int64                  `swagger:"maxItems,omitempty"`
	MinItems         int64                  `swagger:"minItems,omitempty"`
	UniqueItems      bool                   `swagger:"uniqueItems,omitempty"`
	MultipleOf       float64                `swagger:"multipleOf,omitempty"`
	Enum             []interface{}          `swagger:"enum,omitempty"`
	Type             string                 `swagger:"type,omitempty"`
	Format           string                 `swagger:"format,omitempty"`
	Name             string                 `swagger:"name,omitempty"`
	In               string                 `swagger:"in,omitempty"`
	Required         bool                   `swagger:"required,omitempty"`
	Schema           *Schema                `swagger:"schema,omitempty"` // when in == "body"
	CollectionFormat string                 `swagger:"collectionFormat,omitempty"`
	Default          interface{}            `swagger:"default,omitempty"`
}

// UnmarshalMap hydrates this parameter instance with the data from the map
func (p *Parameter) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, p); err != nil {
		return err
	}
	if ref, ok := dict["$ref"]; ok {
		p.Ref = ref.(string)
	}
	p.Extensions = readExtensions(dict)
	return nil
}

// UnmarshalJSON hydrates this parameter instance with the data from JSON
func (p *Parameter) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return p.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this parameter instance with the data from YAML
func (p *Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return p.UnmarshalMap(value)
}

// MarshalMap converts this parameter object to a map
func (p Parameter) MarshalMap() map[string]interface{} {
	res := reflection.MarshalMapRecursed(p)
	if p.Ref != "" {
		res["$ref"] = p.Ref
	}
	addExtensions(res, p.Extensions)
	return res
}

// MarshalJSON converts this parameter object to JSON
func (p Parameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.MarshalMap())
}

// MarshalYAML converts this parameter object to YAML
func (p Parameter) MarshalYAML() (interface{}, error) {
	return p.MarshalMap(), nil
}
