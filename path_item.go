package swagger

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/reflection"
)

// PathItem describes the operations available on a single path.
// A Path Item may be empty, due to [ACL constraints](http://goo.gl/8us55a#securityFiltering).
// The path itself is still exposed to the documentation viewer but they will
// not know which operations and parameters are available.
//
// For more information: http://goo.gl/8us55a#pathItemObject
type PathItem struct {
	Ref        string                 `swagger:"-"`
	Extensions map[string]interface{} `swagger:"-"` // custom extensions, omitted when empty
	Get        *Operation             `swagger:"get,omitempty"`
	Put        *Operation             `swagger:"put,omitempty"`
	Post       *Operation             `swagger:"post,omitempty"`
	Delete     *Operation             `swagger:"delete,omitempty"`
	Options    *Operation             `swagger:"options,omitempty"`
	Head       *Operation             `swagger:"head,omitempty"`
	Patch      *Operation             `swagger:"patch,omitempty"`
	Parameters []Parameter            `swagger:"parameters,omitempty"`
}

// MarshalMap converts this path item to a map
func (p PathItem) MarshalMap() map[string]interface{} {
	if p.Ref != "" {
		return map[string]interface{}{"$ref": p.Ref}
	}

	res := reflection.MarshalMapRecursed(p)
	addExtensions(res, p.Extensions)

	return res
}

// MarshalJSON converts this path item to
func (p PathItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.MarshalMap())
}

// MarshalYAML converts this path item to YAML
func (p PathItem) MarshalYAML() (interface{}, error) {
	return p.MarshalMap(), nil
}
