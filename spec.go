package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// Spec this is the root document object for the API specification.
// It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
//
// For more information: http://goo.gl/8us55a#swagger-object-
type Spec struct {
	Consumes            []string               `swagger:"consumes,omitempty"`
	Produces            []string               `swagger:"produces,omitempty"`
	Schemes             []string               `swagger:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Swagger             string                 `swagger:"swagger"`
	Info                Info                   `swagger:"info"`
	Host                string                 `swagger:"host,omitempty"`
	BasePath            string                 `swagger:"basePath,omitempty"` // must start with a leading "/"
	Paths               Paths                  `swagger:"paths"`              // required
	Definitions         Definitions            `swagger:"definitions,omitempty"`
	Parameters          []Parameter            `swagger:"parameters,omitempty"`
	Responses           ResponsesMap           `swagger:"responses,omitempty"`
	SecurityDefinitions SecurityDefinitions    `swagger:"securityDefintions"`
	Security            SecurityRequirements   `swagger:"security,omitempty"`
	Tags                []Tag                  `swagger:"tags,omitempty"`
	ExternalDocs        *ExternalDocumentation `swagger:"externalDocs,omitempty"`
}

// MarshalJSON converts this spec object to JSON
func (s Spec) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(s))
}

// MarshalYAML converts this spec object to YAML
func (s Spec) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(s), nil
}

// UnmarshalJSON hydrates this spec instance with the data from JSON
func (s *Spec) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, s)
}

// UnmarshalYAML hydrates this spec instance with the data from YAML
func (s *Spec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, s)
}
