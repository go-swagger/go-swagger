package validate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateContentType(t *testing.T) {
	data := []struct {
		hdr     string
		allowed []string
		err     *Error
	}{
		{"application/json", []string{"application/json"}, nil},
		{"application/json", []string{"application/x-yaml", "text/html"}, invalidContentType("application/json", []string{"application/x-yaml", "text/html"})},
		{"text/html; charset=utf-8", []string{"text/html"}, nil},
		{"text/html;charset=utf-8", []string{"text/html"}, nil},
		{"", []string{"application/json"}, invalidContentType("", []string{"application/json"})},
		{"text/html;           charset=utf-8", []string{"application/json"}, invalidContentType("text/html;           charset=utf-8", []string{"application/json"})},
		{"application(", []string{"application/json"}, invalidContentType("application(", []string{"application/json"})},
		{"application/json;char*", []string{"application/json"}, invalidContentType("application/json;char*", []string{"application/json"})},
	}

	for _, v := range data {
		err := ContentType(v.allowed, v.hdr)
		if v.err == nil {
			assert.NoError(t, err, "input: %q", v.hdr)
		} else {
			assert.Error(t, err, "input: %q", v.hdr)
			assert.IsType(t, &Error{}, err, "input: %q", v.hdr)
			assert.Equal(t, v.err.Error(), err.Error(), "input: %q", v.hdr)
			assert.Equal(t, http.StatusUnsupportedMediaType, err.Code())
		}
	}
}
