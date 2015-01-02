package validate

import (
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestParameterValidation(t *testing.T) {
	var param = spec.Parameter{}
	param.In = "query"
	param.Type = "string"
	param.Name = "name"

	p, err := Parameter(&param)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	res := p.Validate(123.0) // needs to be json number so float
	assert.False(t, res.Valid())
}
