package swagger

import (
	"encoding/json"
	"strings"
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
