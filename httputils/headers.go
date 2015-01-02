package httputils

import (
	"mime"
	"net/http"
)

const (
	charsetKey = "charset"
	// DefaultMime the default fallback mime type
	DefaultMime = "application/octet-stream"
	// JSONMime the json mime type
	JSONMime = "application/json"
	// YAMLMime the yaml mime type
	YAMLMime = "application/x-yaml"
)

// ContentType parses a content type header
func ContentType(headers http.Header) (string, string, *ParseError) {
	ct := headers.Get(HeaderContentType)
	orig := ct
	if ct == "" {
		ct = DefaultMime
	}

	mt, opts, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", "", NewParseError(HeaderContentType, "header", orig, err)
	}

	if cs, ok := opts[charsetKey]; ok {
		return mt, cs, nil
	}

	return mt, "", nil
}
