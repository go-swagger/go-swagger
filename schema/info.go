package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// Info object provides metadata about the API.
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
