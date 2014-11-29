package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// XMLObject a metadata object that allows for more fine-tuned XML model definitions.
//
// For more information: http://goo.gl/8us55a#xmlObject
type XMLObject struct {
	Name      string `swagger:"name,omitempty"`
	Namespace string `swagger:"namespace,omitempty"`
	Prefix    string `swagger:"prefix,omitempty"`
	Attribute bool   `swagger:"attribute,omitempty"`
	Wrapped   bool   `swagger:"wrapped,omitempty"`
}

// MarshalJSON converts this spec object to JSON
func (x XMLObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(x))
}

// MarshalYAML converts this spec object to YAML
func (x XMLObject) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(x), nil
}

// UnmarshalJSON hydrates this spec instance with the data from JSON
func (x *XMLObject) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, x)
}

// UnmarshalYAML hydrates this spec instance with the data from YAML
func (x *XMLObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, x)
}
