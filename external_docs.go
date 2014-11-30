package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// ExternalDocumentation allows referencing an external resource for
// extended documentation.
//
// For more information: http://goo.gl/8us55a#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `swagger:"description,omitempty"`
	URL         string `swagger:"url,omitempty"`
}

// MarshalMap converts this external documentation object to map
func (e ExternalDocumentation) MarshalMap() map[string]interface{} {
	return reflection.MarshalMapRecursed(e)
}

// UnmarshalMap hydrates this external documentation instance with the data from a map
func (e *ExternalDocumentation) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, e); err != nil {
		return err
	}
	return nil
}

// MarshalJSON converts this external documentation object to JSON
func (e ExternalDocumentation) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(e))
}

// MarshalYAML converts this external documentation object to YAML
func (e ExternalDocumentation) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(e), nil
}

// UnmarshalJSON hydrates this external documentation instance with the data from JSON
func (e *ExternalDocumentation) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, e)
}

// UnmarshalYAML hydrates this external documentation instance with the data from YAML
func (e *ExternalDocumentation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, e)
}
