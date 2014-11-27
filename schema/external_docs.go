package schema

// External Documentation allows referencing an external resource for
// extended documentation.
//
// For more information: http://goo.gl/8us55a#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `structs:"description,omitempty"`
	URL         string `structs:"url"`
}
