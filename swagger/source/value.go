package source

// Value represents the source of a value
//go:generate stringer -type=Value
type Value int8

const (
	// Query a value obtained from a query string
	Query Value = iota
	// Path a value obtained from a routing param
	Path
	// Header a value obtained from a header
	Header
	// Body a value obtained from a request body
	Body
	// Form a value obtained from a request form parameter
	Form
)
