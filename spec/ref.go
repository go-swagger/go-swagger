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

// // Resolve resolves refs
// func (r Ref) Resolve() (json.RawMessage, error) {
// 	if len(r) == 0 { // bail when we're empty
// 		return nil, nil
// 	}
// 	// check file scheme
// 	// check http scheme
// 	// cache resolved references
// 	return nil, nil
// }
