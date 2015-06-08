package spec

import (
	"encoding/json"
	"strings"

	"github.com/casualjim/go-swagger/jsonpointer"
	"github.com/casualjim/go-swagger/swag"
)

// Extensions vendor specific extensions
type Extensions map[string]interface{}

// Add adds a value to these extensions
func (e Extensions) Add(key string, value interface{}) {
	realKey := strings.ToLower(key)
	e[realKey] = value
}

// GetString gets a string value from the extensions
func (e Extensions) GetString(key string) (string, bool) {
	if v, ok := e[strings.ToLower(key)]; ok {
		str, ok := v.(string)
		return str, ok
	}
	return "", false
}

type vendorExtensible struct {
	Extensions Extensions
}

func (v *vendorExtensible) AddExtension(key string, value interface{}) {
	if value == nil {
		return
	}
	if v.Extensions == nil {
		v.Extensions = make(map[string]interface{})
	}
	v.Extensions.Add(key, value)
}

func (v vendorExtensible) MarshalJSON() ([]byte, error) {
	toser := make(map[string]interface{})
	for k, v := range v.Extensions {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-") {
			toser[k] = v
		}
	}
	return json.Marshal(toser)
}

func (v *vendorExtensible) UnmarshalJSON(data []byte) error {
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	for k, vv := range d {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-") {
			if v.Extensions == nil {
				v.Extensions = map[string]interface{}{}
			}
			v.Extensions[k] = vv
		}
	}
	return nil
}

type infoProps struct {
	Description    string       `json:"description,omitempty"`
	Title          string       `json:"title,omitempty"`
	TermsOfService string       `json:"termsOfService,omitempty"`
	Contact        *ContactInfo `json:"contact,omitempty"`
	License        *License     `json:"license,omitempty"`
	Version        string       `json:"version,omitempty"`
}

// Info object provides metadata about the API.
// The metadata can be used by the clients if needed, and can be presented in the Swagger-UI for convenience.
//
// For more information: http://goo.gl/8us55a#infoObject
type Info struct {
	vendorExtensible
	infoProps
}

// JSONLookup look up a value by the json property name
func (i Info) JSONLookup(token string) (interface{}, error) {
	if ex, ok := i.Extensions[token]; ok {
		return &ex, nil
	}
	r, _, err := jsonpointer.GetForToken(i.infoProps, token)
	return r, err
}

// MarshalJSON marshal this to JSON
func (i Info) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(i.infoProps)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(i.vendorExtensible)
	if err != nil {
		return nil, err
	}
	return swag.ConcatJSON(b1, b2), nil
}

// UnmarshalJSON marshal this from JSON
func (i *Info) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &i.infoProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &i.vendorExtensible); err != nil {
		return err
	}
	return nil
}
