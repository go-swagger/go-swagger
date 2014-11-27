package schema

// License information for the exposed API.
//
// For more information: http://goo.gl/8us55a#licenseObject
type License struct {
	Name string `structs:"name"`
	URL  string `structs:"url"`
}
