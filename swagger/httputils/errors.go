package httputils

import "fmt"

// ParseError respresents a parsing error
type ParseError struct {
	Name    string
	In      string
	Value   string
	Reason  error
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

const (
	parseErrorTemplContent = `parsing %s %s from %q failed, because %s`
)

// NewParseError creates a new parse error
func NewParseError(name, in, value string, reason error) *ParseError {
	msg := fmt.Sprintf(parseErrorTemplContent, name, in, value, reason)
	e := &ParseError{Name: name, In: in, Value: value, Reason: reason, message: msg}
	return e
}
