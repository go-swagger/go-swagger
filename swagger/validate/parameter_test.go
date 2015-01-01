package validate

import (
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

func TestParameterValidation(t *testing.T) {
	param := params["name"]
	p, err := Parameter(&param)
	assert.NoError(t, err)
	assert.NotNil(t, p)

}

var params = map[string]swagger.Parameter{
	"name": swagger.Parameter{
		In:   "query",
		Type: "string",
		Name: "name",
	},
}
