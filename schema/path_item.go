package schema

import "encoding/json"

// PathItem describes the operations available on a single path.
// A Path Item may be empty, due to [ACL constraints](#securityFiltering).
// The path itself is still exposed to the documentation viewer but they will
// not know which operations and parameters are available.
type PathItem struct {
	Ref        string                 `structs:"-"`
	Extensions map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Get        *Operation             `structs:"get,omitempty"`
	Put        *Operation             `structs:"put,omitempty"`
	Post       *Operation             `structs:"post,omitempty"`
	Delete     *Operation             `structs:"delete,omitempty"`
	Options    *Operation             `structs:"options,omitempty"`
	Head       *Operation             `structs:"head,omitempty"`
	Patch      *Operation             `structs:"patch,omitempty"`
	Parameters []Parameter            `structs:"-"`
}

func (p PathItem) Map() map[string]interface{} {
	if p.Ref != "" {
		return map[string]interface{}{"$ref": p.Ref}
	}

	res := make(map[string]interface{})
	addOp := func(key string, op *Operation) {
		if op != nil {
			res[key] = op.Map()
		}
	}
	addOp("get", p.Get)
	addOp("put", p.Put)
	addOp("post", p.Post)
	addOp("delete", p.Delete)
	addOp("head", p.Head)
	addOp("options", p.Options)
	addOp("patch", p.Patch)

	var params []map[string]interface{}
	for _, param := range p.Parameters {
		params = append(params, param.Map())
	}
	res["parameters"] = params

	addExtensions(res, p.Extensions)

	return res
}

func (p PathItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Map())
}
func (p PathItem) MarshalYAML() (interface{}, error) {
	return p.Map(), nil
}
