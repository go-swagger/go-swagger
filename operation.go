package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Operation describes a single API operation on a path.
//
// For more information: http://goo.gl/8us55a#operationObject
type Operation struct {
	Description  string                 `swagger:"description,omitempty"`
	Extensions   map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Consumes     []string               `swagger:"consumes,omitempty"`
	Produces     []string               `swagger:"produces,omitempty"`
	Schemes      []string               `swagger:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Tags         []string               `swagger:"tags,omitempty"`
	Summary      string                 `swagger:"summary,omitempty"`
	ExternalDocs *ExternalDocumentation `swagger:"externalDocs,omitempty"`
	ID           string                 `swagger:"operationId"`
	Deprecated   bool                   `swagger:"deprecated,omitempty"`
	Security     []SecurityRequirement  `swagger:"security,omitempty"`
	Parameters   []Parameter            `swagger:"parameters,omitempty"`
	Responses    Responses              `swagger:"responses"`
}

func (o Operation) MarshalMap() map[string]interface{} {
	res := reflection.MarshalMapRecursed(o)
	addExtensions(res, o.Extensions)
	return res
}

func (o Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.MarshalMap())
}

func (o Operation) MarshalYAML() (interface{}, error) {
	return o.MarshalMap(), nil
}
