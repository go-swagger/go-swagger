package spec

// ContactInfo contact information for the exposed API.
//
// For more information: http://goo.gl/8us55a#contactObject
type ContactInfo struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}
