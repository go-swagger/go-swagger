package httpkit

const (
	// HeaderContentType represents a http content-type header, it's value is supposed to be a mime type
	HeaderContentType = "Content-Type"
	// HeaderAccept the Accept header
	HeaderAccept = "Accept"

	charsetKey = "charset"

	// DefaultMime the default fallback mime type
	DefaultMime = "application/octet-stream"
	// JSONMime the json mime type
	JSONMime = "application/json"
	// YAMLMime the yaml mime type
	YAMLMime = "application/x-yaml"
)
