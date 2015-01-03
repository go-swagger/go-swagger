package errors

import "fmt"

const (
	typeFail          = `%s in %s must be of type %s`
	typeFailWithData  = `%s in %s must be of type %s: %q`
	typeFailWithError = `%s in %s must be of type %s, because: %s`
)

// InvalidCollectionFormat another flavor of invalid type error
func InvalidCollectionFormat(name, in, format string) *Validation {
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   format,
		message: fmt.Sprintf("the collection format %q is not supported for a %s param", format, in),
	}
}

// InvalidType creates an error for when the type is invalid
func InvalidType(name, in, typeName string, value interface{}) *Validation {
	var message string
	switch value.(type) {
	case string:
		message = fmt.Sprintf(typeFailWithData, name, in, typeName, value)
	case error:
		message = fmt.Sprintf(typeFailWithError, name, in, typeName, value)
	default:
		message = fmt.Sprintf(typeFail, name, in, typeName)
	}
	return &Validation{
		code:    422,
		Name:    name,
		In:      in,
		Value:   value,
		message: message,
	}
}
