package schema

type XMLObject struct {
	Name      string `structs:"name,omitempty"`
	Namespace string `structs:"namespace,omitempty"`
	Prefix    string `structs:"prefix,omitempty"`
	Attribute bool   `structs:"attribute,omitempty"`
	Wrapped   bool   `structs:"wrapped,omitempty"`
}
