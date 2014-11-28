package swagger

// ExternalDocumentation allows referencing an external resource for
// extended documentation.
//
// For more information: http://goo.gl/8us55a#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `swagger:"description,omitempty"`
	URL         string `swagger:"url"`
}
