package spec

// License information for the exposed API.
//
// For more information: http://goo.gl/8us55a#licenseObject
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
