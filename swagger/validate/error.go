package validate

import "fmt"

// Error represents a general swagger error
type Error struct {
	code    int32
	Name    string
	In      string
	Value   interface{}
	message string
	Values  []interface{}
}

func (e *Error) Error() string {
	return e.message
}

// Code the error code
func (e *Error) Code() int32 {
	return e.code
}

var (
	contentTypeFail = `Unsupported media type %q, only %v`
)

func invalidContentType(value string, allowed []string) *Error {
	var values []interface{}
	for _, v := range allowed {
		values = append(values, v)
	}
	return &Error{
		code:    415,
		Name:    "Content-Type",
		In:      "header",
		Value:   value,
		Values:  values,
		message: fmt.Sprintf(contentTypeFail, value, allowed),
	}
}
