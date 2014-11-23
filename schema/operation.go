package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// Operation describes a single API operation on a path.
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
