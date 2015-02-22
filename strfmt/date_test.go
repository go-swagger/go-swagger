package strfmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	pp := Date{}
	err := pp.UnmarshalText([]byte{})
	assert.NoError(t, err)
	err = pp.UnmarshalText([]byte("yada"))
	assert.Error(t, err)
	orig := "2014-12-15"
	err = pp.UnmarshalText([]byte(orig))
	assert.NoError(t, err)
	txt, err := pp.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, orig, string(txt))
}
