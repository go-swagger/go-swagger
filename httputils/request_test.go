package httputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONRequest(t *testing.T) {
	req, err := JSONRequest("GET", "/swagger.json", nil)
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, JSONMime, req.Header.Get(HeaderContentType))
	assert.Equal(t, JSONMime, req.Header.Get(HeaderAccept))

	req, err = JSONRequest("GET", "%2", nil)
	assert.Error(t, err)
	assert.Nil(t, req)
}
