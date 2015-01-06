package spec

import (
	"encoding/json"
	"reflect"
)

// Swagger this is the root document object for the API specification.
// It combines what previously was the Resource Listing and API Declaration (version 1.2 and earlier) together into one document.
//
// For more information: http://goo.gl/8us55a#swagger-object-
type Swagger struct {
	Consumes            []string               `json:"consumes,omitempty"`
	Produces            []string               `json:"produces,omitempty"`
	Schemes             []string               `json:"schemes,omitempty"` // the scheme, when present must be from [http, https, ws, wss]
	Swagger             string                 `json:"swagger,omitempty"`
	Info                *Info                  `json:"info,omitempty"`
	Host                string                 `json:"host,omitempty"`
	BasePath            string                 `json:"basePath,omitempty"` // must start with a leading "/"
	Paths               *Paths                 `json:"paths,omitempty"`    // required
	Definitions         Definitions            `json:"definitions,omitempty"`
	Parameters          map[string]Parameter   `json:"parameters,omitempty"`
	Responses           map[string]Response    `json:"responses,omitempty"`
	SecurityDefinitions SecurityDefinitions    `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string  `json:"security,omitempty"`
	Tags                []Tag                  `json:"tags,omitempty"`
	ExternalDocs        *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Dependencies represent a dependencies property
type Dependencies map[string]SchemaOrStringArray

// SchemaOrBool represents a schema or boolean value, is biased towards true for the boolean property
type SchemaOrBool struct {
	Allows bool
	Schema *Schema
}

// MarshalJSON convert this object to JSON
func (s SchemaOrBool) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}

	if s.Schema == nil && s.Allows {
		return jsTrue, nil
	}
	return nil, nil
}

var jsTrue = []byte("true")

// UnmarshalJSON converts this bool or schema object from a JSON structure
func (s *SchemaOrBool) UnmarshalJSON(data []byte) error {
	var nw SchemaOrBool
	// nw.Allows = true
	if len(data) < 5 {
		return nil
	}
	if data[0] == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	nw.Allows = !(data[0] == 'f' && data[1] == 'a' && data[2] == 'l' && data[3] == 's' && data[4] == 'e')
	*s = nw

	return nil
}

// SchemaOrStringArray represents a schema or a string array
type SchemaOrStringArray struct {
	Schema   *Schema
	Property []string
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s.Property) > 0 {
		return json.Marshal(s.Property)
	}
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return nil, nil
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrStringArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	first := data[0]
	var nw SchemaOrStringArray
	if first == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Property); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

// Definitions contains the models explicitly defined in this spec
// An object to hold data types that can be consumed and produced by operations.
// These data types can be primitives, arrays or models.
//
// For more information: http://goo.gl/8us55a#definitionsObject
type Definitions map[string]Schema

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

// StringOrArray represents a value that can either be a string
// or an array of strings. Mainly here for serialization purposes
type StringOrArray []string

// Contains returns true when the value is contained in the slice
func (s StringOrArray) Contains(value string) bool {
	for _, str := range s {
		if str == value {
			return true
		}
	}
	return false
}

// UnmarshalJSON unmarshals this string or array object from a JSON array or JSON string
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}

	if data[0] == '[' {
		var parsed []string
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		*s = StringOrArray(parsed)
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	if single != "" {
		*s = StringOrArray([]string{single})
	}
	return nil
}

// MarshalJSON converts this string or array to a JSON array or JSON string
func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal([]string(s)[0])
	}
	return json.Marshal([]string(s))
}

// SchemaOrArray represents a value that can either be a Schema
// or an array of Schema. Mainly here for serialization purposes
type SchemaOrArray struct {
	Schema  *Schema
	Schemas []Schema
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	if len(s.Schemas) > 0 {
		return json.Marshal(s.Schemas)
	}
	return json.Marshal(s.Schema)
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	first := data[0]
	var nw SchemaOrArray
	if first == '{' {
		var sch Schema
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		if reflect.DeepEqual(Schema{}, sch) {
			return nil
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Schemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
	// if len(data) < 3 {
	// 	return nil
	// }
	// if data[0] == '[' {
	// 	var parsed []Schema
	// 	if err := json.Unmarshal(data, &parsed); err != nil {
	// 		return err
	// 	}
	// 	*s = SchemaOrArray{Schemas: parsed}
	// 	return nil
	// }

	// var parsed Schema
	// if err := json.Unmarshal(data, &parsed); err != nil {
	// 	return err
	// }
	// if reflect.DeepEqual(Schema{}, parsed) {
	// 	return nil
	// }
	// *s = SchemaOrArray{Schema: &parsed}
	// return nil
}
