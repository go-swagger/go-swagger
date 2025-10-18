package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/spec"
)

type responseTestContext struct {
	OpID      string
	Name      string
	IsSuccess bool
}

func (ctx *responseTestContext) Assert(t *testing.T, response spec.Response, res GenResponse) bool {
	if !assert.Equal(t, ctx.IsSuccess, res.IsSuccess) {
		return false
	}

	if !assert.Equal(t, ctx.Name, res.Name) {
		return false
	}

	if !assert.Equal(t, response.Description, res.Description) {
		return false
	}

	if !ctx.assertSchema(t, response, res) {
		return false
	}

	if len(response.Headers) > 0 {
		if !ctx.assertHeaders(t, response, res) {
			return false
		}
	}

	return true
}

func (ctx *responseTestContext) Require(t *testing.T, response spec.Response, res GenResponse) {
	if !ctx.Assert(t, response, res) {
		t.FailNow()
	}
}

func (ctx *responseTestContext) assertSchema(t *testing.T, response spec.Response, res GenResponse) bool {
	if response.Schema != nil {
		if !assert.NotNil(t, res.Schema) {
			return false
		}
	} else {
		if !assert.Nil(t, res.Schema) {
			return false
		}
	}

	return true
}

func (ctx *responseTestContext) assertHeaders(t *testing.T, response spec.Response, res GenResponse) bool {
	for k, v := range response.Headers {
		found := false
		for _, h := range res.Headers {
			if h.Name == k {
				found = true
				if k == "X-Last-Task-Id" {
					hctx := &respHeaderTestContext{k, "swag.FormatInt64", "swag.ConvertInt64"}
					if !hctx.Assert(t, v, h) {
						return false
					}
					break
				}
				if k == "X-Error-Code" {
					hctx := &respHeaderTestContext{k, "", ""}
					if !hctx.Assert(t, v, h) {
						return false
					}
				}

				break
			}
		}

		if !assert.True(t, found) {
			return false
		}
	}

	return true
}

type respHeaderTestContext struct {
	Name      string
	Formatter string
	Converter string
}

func (ctx *respHeaderTestContext) Assert(t *testing.T, header spec.Header, hdr GenHeader) bool {
	if !assert.Equal(t, ctx.Name, hdr.Name) {
		return false
	}
	if !assert.Equal(t, "a", hdr.ReceiverName) {
		return false
	}
	if !assert.Equal(t, ctx.Formatter, hdr.Formatter) {
		return false
	}
	if !assert.Equal(t, ctx.Converter, hdr.Converter) {
		return false
	}
	if !assert.Equal(t, header.Description, hdr.Description) {
		return false
	}
	if !assert.Equal(t, header.Minimum, hdr.Minimum) || !assert.Equal(t, header.ExclusiveMinimum, hdr.ExclusiveMinimum) {
		return false
	}
	if !assert.Equal(t, header.Maximum, hdr.Maximum) || !assert.Equal(t, header.ExclusiveMaximum, hdr.ExclusiveMaximum) {
		return false
	}
	if !assert.Equal(t, header.MinLength, hdr.MinLength) {
		return false
	}
	if !assert.Equal(t, header.MaxLength, hdr.MaxLength) {
		return false
	}
	if !assert.Equal(t, header.Pattern, hdr.Pattern) {
		return false
	}
	if !assert.Equal(t, header.MaxItems, hdr.MaxItems) {
		return false
	}
	if !assert.Equal(t, header.MinItems, hdr.MinItems) {
		return false
	}
	if !assert.Equal(t, header.UniqueItems, hdr.UniqueItems) {
		return false
	}
	if !assert.Equal(t, header.MultipleOf, hdr.MultipleOf) {
		return false
	}
	if !assert.Equal(t, header.Enum, hdr.Enum) {
		return false
	}
	if !assert.Equal(t, header.Type, hdr.SwaggerType) {
		return false
	}
	if !assert.Equal(t, header.Format, hdr.SwaggerFormat) {
		return false
	}
	return true
}
