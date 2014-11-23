package schema

// ContactInfo contact information for the exposed API.
type ContactInfo struct {
	Name  string `structs:"name"`
	URL   string `structs:"url"`
	Email string `structs:"email"`
}
