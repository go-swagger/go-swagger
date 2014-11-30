package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// License information for the exposed API.
//
// For more information: http://goo.gl/8us55a#licenseObject
type License struct {
	Name string `swagger:"name,omitempty"`
	URL  string `swagger:"url,omitempty"`
}

// MarshalMap converts this license object to map
func (l License) MarshalMap() map[string]interface{} {
	return reflection.MarshalMapRecursed(l)
}

// UnmarshalMap hydrates this license instance with the data from a map
func (l *License) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, l); err != nil {
		return err
	}
	return nil
}

// MarshalJSON converts this spec object to JSON
func (l License) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(l))
}

// MarshalYAML converts this spec object to YAML
func (l License) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(l), nil
}

// UnmarshalJSON hydrates this spec instance with the data from JSON
func (l *License) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, l)
}

// UnmarshalYAML hydrates this spec instance with the data from YAML
func (l *License) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, l)
}
