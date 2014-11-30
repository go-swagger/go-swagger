package swagger

import (
	"encoding/json"
	"strings"

	"github.com/casualjim/go-swagger/reflection"
)

// Paths holds the relative paths to the individual endpoints.
// The path is appended to the [`basePath`](http://goo.gl/8us55a#swaggerBasePath) in order
// to construct the full URL.
// The Paths may be empty, due to [ACL constraints](http://goo.gl/8us55a#securityFiltering).
//
// For more information: http://goo.gl/8us55a#pathsObject
type Paths struct {
	Extensions map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Paths      map[string]PathItem    `swagger:"-"` // custom serializer to flatten this, each entry must start with "/"
}

// UnmarshalMap hydrates this paths instance with the data from the map
func (p *Paths) UnmarshalMap(data interface{}) error {
	dict := reflection.MarshalMap(data)
	var res map[string]PathItem
	for k, v := range dict {
		if strings.HasPrefix(k, "/") {
			if res == nil {
				res = make(map[string]PathItem)
			}
			pathItem := PathItem{}
			if err := reflection.UnmarshalMap(reflection.MarshalMap(v), &pathItem); err != nil {
				return err
			}
			res[k] = pathItem
			delete(dict, k)
		}
	}
	p.Paths = res
	p.Extensions = readExtensions(dict)
	return nil
}

// UnmarshalJSON hydrates this paths instance with the data from JSON
func (p *Paths) UnmarshalJSON(data []byte) error {
	var value map[string]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return p.UnmarshalMap(value)
}

// UnmarshalYAML hydrates this paths instance with the data from YAML
func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]interface{}
	if err := unmarshal(&value); err != nil {
		return err
	}
	return p.UnmarshalMap(value)
}

// MarshalMap converts this paths object to a map
func (p Paths) MarshalMap() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range p.Paths {
		key := k
		if !strings.HasPrefix(key, "/") {
			key = "/" + key
		}
		res[key] = v.MarshalMap()
	}
	addExtensions(res, p.Extensions)
	return res
}

// MarshalJSON converts this paths object to JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.MarshalMap())
}

// MarshalYAML converts this paths object to YAML
func (p Paths) MarshalYAML() (interface{}, error) {
	return p.MarshalMap(), nil
}
