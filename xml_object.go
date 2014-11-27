package swagger

// XMLObject a metadata object that allows for more fine-tuned XML model definitions.
//
// For more information: http://goo.gl/8us55a#xmlObject
type XMLObject struct {
	Name      string `structs:"name,omitempty"`
	Namespace string `structs:"namespace,omitempty"`
	Prefix    string `structs:"prefix,omitempty"`
	Attribute bool   `structs:"attribute,omitempty"`
	Wrapped   bool   `structs:"wrapped,omitempty"`
}
