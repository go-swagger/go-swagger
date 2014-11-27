package schema

// Definitions contains the models explicitly defined in this spec
// An object to hold data types that can be consumed and produced by operations.
// These data types can be primitives, arrays or models.
//
// For more information: http://goo.gl/8us55a#definitionsObject
type Definitions map[string]Schema

// Map generates a map[string]interface{} from the definitions
func (d Definitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range d {
		res[k] = v.Map()
	}
	return res
}
