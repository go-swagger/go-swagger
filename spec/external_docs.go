package spec

// ExternalDocumentation allows referencing an external resource for
// extended documentation.
//
// For more information: http://goo.gl/8us55a#externalDocumentationObject
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}
