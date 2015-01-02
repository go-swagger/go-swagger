package httputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
