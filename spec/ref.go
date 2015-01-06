package spec

import "encoding/json"

// Ref represents a json reference
type Ref string

// MarshalJSON marshal this to JSON
func (r Ref) MarshalJSON() ([]byte, error) {
	if r == "" {
		return []byte("{}"), nil
	}
	v := map[string]interface{}{"$ref": string(r)}
	return json.Marshal(v)
}

// UnmarshalJSON unmarshal this from JSON
func (r *Ref) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	if vv, ok := v["$ref"]; ok {
		if str, ok := vv.(string); ok {
			*r = Ref(str)
		}
	}
	return nil
}

// RefCache is a thread-safe cache of refs that can be optionally passed in
// when the value is nil it will just always fetch
// // Resolve resolves refs
// func (r Ref) Resolve(cache RefCache, root *Schema) (json.RawMessage, error) {
// 	if len(r) == 0 { // bail when we're empty
// 		return nil, nil
// 	}
// 	// check file scheme
// 	// check http scheme
// 	// cache resolved references
// 	return nil, nil
// }
