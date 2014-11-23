package schema

type Definitions map[string]Schema

func (d Definitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range d {
		res[k] = v.Map()
	}
	return res
}
