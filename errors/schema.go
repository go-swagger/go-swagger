package errors

import "fmt"

const (
	typeFail = `%s in %s must be of type %s`
)

// InvalidType creates an error for when the type is invalid
func InvalidType(name, in, typeName string, value interface{}) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   value,
		message: fmt.Sprintf(typeFail, name, in, typeName),
	}
}
