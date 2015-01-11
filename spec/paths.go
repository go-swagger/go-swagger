package spec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/casualjim/go-swagger/util"
)

// Paths holds the relative paths to the individual endpoints.
// The path is appended to the [`basePath`](http://goo.gl/8us55a#swaggerBasePath) in order
// to construct the full URL.
// The Paths may be empty, due to [ACL constraints](http://goo.gl/8us55a#securityFiltering).
//
// For more information: http://goo.gl/8us55a#pathsObject
type Paths struct {
	vendorExtensible
	Paths map[string]PathItem `json:"-"` // custom serializer to flatten this, each entry must start with "/"
}

func (p Paths) JSONLookup(token string) (interface{}, error) {
	if pi, ok := p.Paths[token]; ok {
		return &pi, nil
	}
	if ex, ok := p.Extensions[token]; ok {
		return &ex, nil
	}
	return nil, fmt.Errorf("object has no field %q", token)
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (p *Paths) UnmarshalJSON(data []byte) error {
	var res map[string]json.RawMessage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	for k, v := range res {
		if strings.HasPrefix(strings.ToLower(k), "x-") {
			if p.Extensions == nil {
				p.Extensions = make(map[string]interface{})
			}
			var d interface{}
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			p.Extensions[k] = d
		}
		if strings.HasPrefix(k, "/") {
			if p.Paths == nil {
				p.Paths = make(map[string]PathItem)
			}
			var pi PathItem
			if err := json.Unmarshal(v, &pi); err != nil {
				return err
			}
			p.Paths[k] = pi
		}
	}
	return nil
}

// MarshalJSON converts this items object to JSON
func (p Paths) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(p.vendorExtensible)
	if err != nil {
		return nil, err
	}
	var pths map[string]PathItem
	for k, v := range p.Paths {
		if strings.HasPrefix(k, "/") {
			if pths == nil {
				pths = make(map[string]PathItem)
			}
			pths[k] = v
		}
	}
	b2, err := json.Marshal(pths)
	if err != nil {
		return nil, err
	}
	concated := util.ConcatJSON(b1, b2)
	return concated, nil
}
