package validate

import (
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateContentType(t *testing.T) {
	data := []struct {
		hdr     string
		allowed []string
		err     *errors.Validation
	}{
		{"application/json", []string{"application/json"}, nil},
		{"application/json", []string{"application/x-yaml", "text/html"}, errors.InvalidContentType("application/json", []string{"application/x-yaml", "text/html"})},
		{"text/html; charset=utf-8", []string{"text/html"}, nil},
		{"text/html;charset=utf-8", []string{"text/html"}, nil},
		{"", []string{"application/json"}, errors.InvalidContentType("", []string{"application/json"})},
		{"text/html;           charset=utf-8", []string{"application/json"}, errors.InvalidContentType("text/html;           charset=utf-8", []string{"application/json"})},
		{"application(", []string{"application/json"}, errors.InvalidContentType("application(", []string{"application/json"})},
		{"application/json;char*", []string{"application/json"}, errors.InvalidContentType("application/json;char*", []string{"application/json"})},
	}

	for _, v := range data {
		err := ContentType(v.allowed, v.hdr)
		if v.err == nil {
			assert.NoError(t, err, "input: %q", v.hdr)
		} else {
			assert.Error(t, err, "input: %q", v.hdr)
			assert.IsType(t, &errors.Validation{}, err, "input: %q", v.hdr)
			assert.Equal(t, v.err.Error(), err.Error(), "input: %q", v.hdr)
			assert.Equal(t, http.StatusUnsupportedMediaType, err.Code())
		}
	}
}
