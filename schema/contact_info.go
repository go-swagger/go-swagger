package schema

// ContactInfo contact information for the exposed API.
//
// For more information: http://goo.gl/8us55a#contactObject
type ContactInfo struct {
	Name  string `structs:"name"`
	URL   string `structs:"url"`
	Email string `structs:"email"`
}
