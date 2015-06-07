package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	err := NewParseError("Content-Type", "header", "application(", errors.New("unable to parse"))
	assert.EqualValues(t, 400, err.Code())
	assert.Equal(t, "parsing Content-Type header from \"application(\" failed, because unable to parse", err.Error())

	err = NewParseError("Content-Type", "", "application(", errors.New("unable to parse"))
	assert.EqualValues(t, 400, err.Code())
	assert.Equal(t, "parsing Content-Type from \"application(\" failed, because unable to parse", err.Error())
}
