package swagger

// ContactInfo contact information for the exposed API.
//
// For more information: http://goo.gl/8us55a#contactObject
type ContactInfo struct {
	Name  string `swagger:"name"`
	URL   string `swagger:"url"`
	Email string `swagger:"email"`
}
