package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/jsonreference"
)

type refable struct {
	Ref Ref
}

func (r refable) MarshalJSON() ([]byte, error) {
	return r.Ref.MarshalJSON()
}

func (r *refable) UnmarshalJSON(d []byte) error {
	return json.Unmarshal(d, &r.Ref)
}

// Ref represents a json reference that is potentially resolved
type Ref struct {
	jsonreference.Ref
	Resolved interface{}
}

// IsResolved returns true if the reference has been resolved
func (r *Ref) IsResolved() bool {
	return r.Resolved != nil
}

// NewRef creates a new instance of a ref object
// returns an error when the reference uri is an invalid uri
func NewRef(refURI string) (Ref, error) {
	ref, err := jsonreference.New(refURI)
	if err != nil {
		return Ref{}, err
	}
	return Ref{Ref: ref}, nil
}

// MustCreateRef creates a ref object but
func MustCreateRef(refURI string) Ref {
	return Ref{Ref: jsonreference.MustCreateRef(refURI)}
}

// MarshalJSON marshals this ref into a JSON object
func (r Ref) MarshalJSON() ([]byte, error) {
	str := r.String()
	if str == "" {
		return []byte("{}"), nil
	}
	v := map[string]interface{}{"$ref": str}
	return json.Marshal(v)
}

// UnmarshalJSON unmarshals this ref from a JSON object
func (r *Ref) UnmarshalJSON(d []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(d, &v); err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	if vv, ok := v["$ref"]; ok {
		if str, ok := vv.(string); ok {
			ref, err := jsonreference.New(str)
			if err != nil {
				return err
			}
			*r = Ref{Ref: ref}
		}
	}
	return nil
}

// // Resolve resolves refs
// // RefCache is a thread-safe cache of refs that can be optionally passed in
// // when the value is nil it will just always fetch
// func (r *Ref) resolve(cache *refCache, root interface{}) (interface{}, error) {
// 	if !r.IsCanonical() {
// 		return nil, fmt.Errorf("reference %q must be canonical", str)
// 	}

// 	return nil, nil
// }
