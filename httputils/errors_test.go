package httputils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := NewParseError("Content-Type", "header", "application(", errors.New("unable to parse"))
	assert.Equal(t, 400, err.Code())
	assert.Equal(t, "parsing Content-Type header from \"application(\" failed, because unable to parse", err.Error())
}
