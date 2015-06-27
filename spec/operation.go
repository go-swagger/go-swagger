package spec

import (
	"encoding/json"

	"github.com/go-swagger/go-swagger/jsonpointer"
	"github.com/go-swagger/go-swagger/swag"
)

type operationProps struct {
	Description  string                 `json:"description,omitempty"`
	Consumes     []string               `json:"consumes,omitempty"`
	Produces     []string               `json:"produces,omitempty"`
	Schemes      []string               `json:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Tags         []string               `json:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	ID           string                 `json:"operationId,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	Security     []map[string][]string  `json:"security,omitempty"`
	Parameters   []Parameter            `json:"parameters,omitempty"`
	Responses    *Responses             `json:"responses,omitempty"`
}

// Operation describes a single API operation on a path.
//
// For more information: http://goo.gl/8us55a#operationObject
type Operation struct {
	vendorExtensible
	operationProps
}

// SuccessResponse gets a success response model
func (o *Operation) SuccessResponse() (*Response, int, bool) {
	if o.Responses == nil {
		return nil, 0, false
	}

	for k, v := range o.Responses.StatusCodeResponses {
		if k/100 == 2 {
			return &v, k, true
		}
	}

	return o.Responses.Default, 0, false
}

// JSONLookup look up a value by the json property name
func (o Operation) JSONLookup(token string) (interface{}, error) {
	if ex, ok := o.Extensions[token]; ok {
		return &ex, nil
	}
	r, _, err := jsonpointer.GetForToken(o.operationProps, token)
	return r, err
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (o *Operation) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.operationProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &o.vendorExtensible); err != nil {
		return err
	}
	return nil
}

// MarshalJSON converts this items object to JSON
func (o Operation) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(o.operationProps)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(o.vendorExtensible)
	if err != nil {
		return nil, err
	}
	concated := swag.ConcatJSON(b1, b2)
	return concated, nil
}
