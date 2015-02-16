package errors

import "fmt"

// ParseError respresents a parsing error
type ParseError struct {
	code    int32
	Name    string
	In      string
	Value   string
	Reason  error
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

// Code returns the http status code for this error
func (e *ParseError) Code() int32 {
	return e.code
}

const (
	parseErrorTemplContent     = `parsing %s %s from %q failed, because %s`
	parseErrorTemplContentNoIn = `parsing %s %s failed, because %s`
)

// NewParseError creates a new parse error
func NewParseError(name, in, value string, reason error) *ParseError {
	var msg string
	if in == "" {
		msg = fmt.Sprintf(parseErrorTemplContentNoIn, name, value, reason)
	} else {
		msg = fmt.Sprintf(parseErrorTemplContent, name, in, value, reason)
	}
	return &ParseError{
		code:    400,
		Name:    name,
		In:      in,
		Value:   value,
		Reason:  reason,
		message: msg,
	}
}
