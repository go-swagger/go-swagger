package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEnum(t *testing.T) {
	enumValues := []string{"aa", "bb", "cc"}

	err := Enum("test", "body", "a", enumValues)
	assert.Error(t, err)
	err = Enum("test", "body", "bb", enumValues)
	assert.NoError(t, err)
}
