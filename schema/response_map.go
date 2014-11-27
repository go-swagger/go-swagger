package schema

// ResponseMap contains the responses by key
type ResponsesMap map[string]Response

func (r ResponsesMap) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range r {
		res[k] = v.Map()
	}
	return res
}
