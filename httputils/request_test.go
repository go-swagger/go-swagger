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

func TestCanHaveBody(t *testing.T) {
	assert.True(t, CanHaveBody("put"))
	assert.True(t, CanHaveBody("post"))
	assert.True(t, CanHaveBody("patch"))
	assert.False(t, CanHaveBody(""))
	assert.False(t, CanHaveBody("get"))
	assert.False(t, CanHaveBody("options"))
	assert.False(t, CanHaveBody("head"))
	assert.False(t, CanHaveBody("delete"))
	assert.False(t, CanHaveBody("invalid"))
}
