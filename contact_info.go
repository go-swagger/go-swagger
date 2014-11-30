package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// ContactInfo contact information for the exposed API.
//
// For more information: http://goo.gl/8us55a#contactObject
type ContactInfo struct {
	Name  string `swagger:"name,omitempty"`
	URL   string `swagger:"url,omitempty"`
	Email string `swagger:"email,omitempty"`
}

// MarshalJSON converts this contact info object to JSON
func (c ContactInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.MarshalMap())
}

// MarshalYAML converts this contact info object to YAML
func (c ContactInfo) MarshalYAML() (interface{}, error) {
	return c.MarshalMap(), nil
}

// MarshalMap converts this contact info object to map
func (c ContactInfo) MarshalMap() map[string]interface{} {
	return reflection.MarshalMapRecursed(c)
}

// UnmarshalMap hydrates this contact info instance with the data from a map
func (c *ContactInfo) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	if err := reflection.UnmarshalMapRecursed(dict, c); err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON hydrates this contact info instance with the data from JSON
func (c *ContactInfo) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, c)
}

// UnmarshalYAML hydrates this contact info instance with the data from YAML
func (c *ContactInfo) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(value); err != nil {
		return err
	}
	return reflection.UnmarshalMap(value, c)
}
