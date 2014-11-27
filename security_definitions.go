package swagger

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

func (s SecurityDefinitions) Map() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range s {
		res[k] = v.Map()
	}
	return res
}
