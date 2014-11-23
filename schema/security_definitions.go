package schema

type SecurityDefinitions map[string]*SecurityScheme

func (s SecurityDefinitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range s {
		res[k] = v.Map()
	}
	return res
}
