package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthenticated(t *testing.T) {
	err := Unauthenticated("basic")
	assert.Equal(t, 401, err.Code())
	assert.Equal(t, "unauthenticated for basic", err.Error())
}
