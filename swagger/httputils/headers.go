package httputils

import (
	"mime"
	"net/http"
)

const (
	charsetKey  = "charset"
	defaultMime = "application/octet-stream"
)

// ContentType parses a content type header
func ContentType(headers http.Header) (string, string, error) {
	ct := headers.Get(HeaderContentType)
	orig := ct
	if ct == "" {
		ct = defaultMime
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
