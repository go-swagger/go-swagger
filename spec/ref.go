package spec

import (
	"encoding/json"
	"errors"

	"github.com/go-swagger/go-swagger/jsonreference"
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
}

// // PointsToRoot returns true when this reference should point to root
// func (r *Ref) PointsToRoot() bool {
// 	return r.RawURL == "#"
// }

// // IsResolved returns true if the reference has been resolved
// func (r Ref) IsResolved() bool {
// 	return r.Resolved != nil
// }

// // NeedsResolving returns true if the reference needs to be resolved
// func (r Ref) NeedsResolving() bool {
// 	return r.Ref.GetURL() != nil && r.Resolved == nil
// }

// Inherits creates a new reference from a parent and a child
// If the child cannot inherit from the parent, an error is returned
func (r *Ref) Inherits(child Ref) (*Ref, error) {
	childURL := child.GetURL()
	parentURL := r.GetURL()
	if childURL == nil {
		return nil, errors.New("child url is nil")
	}
	if parentURL == nil {
		return &child, nil
	}

	ref, err := NewRef(parentURL.ResolveReference(childURL).String())
	if err != nil {
		return nil, err
	}
	return &ref, err
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

// // NewResolvedRef creates a resolved ref
// func NewResolvedRef(refURI string, data interface{}) Ref {
// 	return Ref{
// 		Ref:      jsonreference.MustCreateRef(refURI),
// 		Resolved: data,
// 	}
// }

// MarshalJSON marshals this ref into a JSON object
func (r Ref) MarshalJSON() ([]byte, error) {
	str := r.String()
	if str == "" {
		if r.IsRoot() {
			return []byte(`{"$ref":"#"}`), nil
		}
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
