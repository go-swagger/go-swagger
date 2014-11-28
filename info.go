package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Info object provides metadata about the API.
// The metadata can be used by the clients if needed, and can be presented in the Swagger-UI for convenience.
//
// For more information: http://goo.gl/8us55a#infoObject
type Info struct {
	Extensions     map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Description    string                 `swagger:"description,omitempty"`
	Title          string                 `swagger:"title,omitempty"`
	TermsOfService string                 `swagger:"termsOfService,omitempty"`
	Contact        *ContactInfo           `swagger:"contact,omitempty"`
	License        *License               `swagger:"license,omitempty"`
	Version        string                 `swagger:"version,omitempty"`
}

func (i Info) MarshalMap() map[string]interface{} {
	res := reflection.MarshalMapRecursed(i)
	addExtensions(res, i.Extensions)
	return res
}

func (i Info) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.MarshalMap())
}

func (i Info) MarshalYAML() (interface{}, error) {
	return i.MarshalMap(), nil
}
