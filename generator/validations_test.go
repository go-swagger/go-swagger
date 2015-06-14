package generator

import (
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestHeaderValidations(t *testing.T) {
	var hdr spec.Header
	hdr.Type = "string"
	var vmin, vmax int64 = 5, 7
	hdr.MinLength = &vmin
	hdr.MaxLength = &vmax

	gh := makeGenHeader("r", "X-Rate-Limited", hdr)
	assert.Equal(t, "XRateLimited", gh.PropertyName)
	assert.Equal(t, "X-Rate-Limited", gh.ParamName)
	assert.Equal(t, "r.XRateLimited", gh.ValueExpression)
	assert.EqualValues(t, *hdr.MinLength, gh.MinLength)
	assert.EqualValues(t, *hdr.MaxLength, gh.MaxLength)
}

func TestParameterValidations(t *testing.T) {
	var hdr spec.Parameter
	hdr.Type = "string"
	var vmin, vmax int64 = 5, 7
	hdr.MinLength = &vmin
	hdr.MaxLength = &vmax
	hdr.Name = "X-Rate-Limited"

	gh, _ := makeCodegenParameter("r", nil, hdr)

	assert.Equal(t, "XRateLimited", gh.PropertyName)
	assert.Equal(t, "xRateLimited", gh.ParamName)
	assert.Equal(t, "r.XRateLimited", gh.ValueExpression)
	assert.EqualValues(t, *hdr.MinLength, gh.MinLength)
	assert.EqualValues(t, *hdr.MaxLength, gh.MaxLength)
}
