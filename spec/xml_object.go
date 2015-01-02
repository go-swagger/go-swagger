package spec

// XMLObject a metadata object that allows for more fine-tuned XML model definitions.
//
// For more information: http://goo.gl/8us55a#xmlObject
type XMLObject struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}
