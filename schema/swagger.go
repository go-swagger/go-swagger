package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

// Swagger this is the root document object for the API specification.
// It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
//
// For more information: http://goo.gl/8us55a#swagger-object-
type Swagger struct {
	Consumes            []string               `structs:"consumes,omitempty"`
	Produces            []string               `structs:"produces,omitempty"`
	Schemes             []string               `structs:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Swagger             string                 `structs:"swagger"`
	Info                Info                   `structs:"-"`
	Host                string                 `structs:"host,omitempty"`
	BasePath            string                 `structs:"basePath,omitempty"` // must start with a leading "/"
	Paths               Paths                  `structs:"-"`                  // required
	Definitions         Definitions            `structs:"-"`
	Parameters          []Parameter            `structs:"-"`
	Responses           ResponsesMap           `structs:"-"`
	SecurityDefinitions SecurityDefinitions    `structs:"-"`
	Security            SecurityRequirements   `structs:"security,omitempty"`
	Tags                []Tag                  `structs:"-"`
	ExternalDocs        *ExternalDocumentation `structs:"externalDocs,omitempty"`
}

func (s Swagger) Map() map[string]interface{} {
	res := structs.Map(s)
	res["info"] = s.Info.Map()
	res["paths"] = s.Paths.Map()
	res["responses"] = s.Responses.Map()
	res["definitions"] = s.Definitions.Map()
	res["securityDefinitions"] = s.SecurityDefinitions.Map()

	var params []map[string]interface{}
	for _, param := range s.Parameters {
		params = append(params, param.Map())
	}
	res["parameters"] = params

	var tags []map[string]interface{}
	for _, t := range s.Tags {
		tags = append(tags, t.Map())
	}
	res["tags"] = tags
	return res
}

func (s Swagger) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s Swagger) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}
