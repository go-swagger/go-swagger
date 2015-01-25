package httputils

import (
	"io"
	"net/http"
)

// JSONRequest creates a new http request with json headers set
func JSONRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, JSONMime)
	req.Header.Add(HeaderAccept, JSONMime)
	return req, nil
}
