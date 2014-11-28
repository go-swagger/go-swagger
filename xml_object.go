package swagger

// XMLObject a metadata object that allows for more fine-tuned XML model definitions.
//
// For more information: http://goo.gl/8us55a#xmlObject
type XMLObject struct {
	Name      string `swagger:"name,omitempty"`
	Namespace string `swagger:"namespace,omitempty"`
	Prefix    string `swagger:"prefix,omitempty"`
	Attribute bool   `swagger:"attribute,omitempty"`
	Wrapped   bool   `swagger:"wrapped,omitempty"`
}
