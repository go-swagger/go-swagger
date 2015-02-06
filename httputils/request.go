package httputils

import (
	"io"
	"net/http"
	"strings"
)

// CanHaveBody returns true if this method can have a body
func CanHaveBody(method string) bool {
	mn := strings.ToUpper(method)
	return mn == "POST" || mn == "PUT" || mn == "PATCH"
}

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

// Gettable for things with a method Get(string) string
type Gettable interface {
	Get(string) string
}

// ReadSingleValue reads a single value from the source
func ReadSingleValue(values Gettable, name string) string {
	return values.Get(name)
}

// ReadCollectionValue reads a collection value from a string data source
func ReadCollectionValue(values Gettable, name, collectionFormat string) []string {
	v := ReadSingleValue(values, name)
	return split(v, collectionFormat)
}

func split(data, format string) []string {
	if data == "" {
		return nil
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	case "multi":
		return nil
	default:
		sep = ","
	}
	var result []string
	for _, s := range strings.Split(data, sep) {
		if ts := strings.TrimSpace(s); ts != "" {
			result = append(result, ts)
		}
	}
	return result
}
