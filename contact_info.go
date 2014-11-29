package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// ContactInfo contact information for the exposed API.
//
// For more information: http://goo.gl/8us55a#contactObject
type ContactInfo struct {
	Name  string `swagger:"name"`
	URL   string `swagger:"url"`
	Email string `swagger:"email"`
}

// MarshalJSON converts this spec object to JSON
func (c ContactInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflection.MarshalMap(c))
}

// MarshalYAML converts this spec object to YAML
func (c ContactInfo) MarshalYAML() (interface{}, error) {
	return reflection.MarshalMap(c), nil
}

// UnmarshalJSON hydrates this spec instance with the data from JSON
func (c *ContactInfo) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, c)
}

// UnmarshalYAML hydrates this spec instance with the data from YAML
func (c *ContactInfo) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, c)
}
