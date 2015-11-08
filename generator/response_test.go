package generator

import (
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestSimpleResponses(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	op, ok := b.Doc.OperationForName("updateTask")
	if assert.True(t, ok) && assert.NotNil(t, op) && assert.NotNil(t, op.Responses) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		if assert.NotNil(t, op.Responses.Default) {
			resp := *op.Responses.Default
			defCtx := responseTestContext{
				OpID: "updateTask",
				Name: "default",
			}
			res, err := b.MakeResponse("a", defCtx.Name, false, resolver, -1, resp)
			if assert.NoError(t, err) {
				if defCtx.Assert(t, resp, res) {
					for code, response := range op.Responses.StatusCodeResponses {
						sucCtx := responseTestContext{
							OpID:      "updateTask",
							Name:      "success",
							IsSuccess: code/100 == 2,
						}
						res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, code, response)
						if assert.NoError(t, err) {
							if !sucCtx.Assert(t, response, res) {
								return
							}
						}
					}
				}
			}
		}
	}

}

type responseTestContext struct {
	OpID      string
	Name      string
	IsSuccess bool
}

func (ctx *responseTestContext) Assert(t testing.TB, response spec.Response, res GenResponse) bool {
	if !assert.Equal(t, ctx.IsSuccess, res.IsSuccess) {
		return false
	}

	if !assert.Equal(t, ctx.Name, res.Name) {
		return false
	}

	if !assert.Equal(t, response.Description, res.Description) {
		return false
	}

	if len(response.Headers) > 0 {
		for k, v := range response.Headers {
			found := false
			for _, h := range res.Headers {
				if h.Name == k {
					found = true
					hctx := &respHeaderTestContext{k, "swag.FormatInt64", "swag.ConvertInt64"}
					if !hctx.Assert(t, v, h) {
						return false
					}
					break
				}
			}
			if !assert.True(t, found) {
				return false
			}
		}
	}

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

type respHeaderTestContext struct {
	Name      string
	Formatter string
	Converter string
}

func (ctx *respHeaderTestContext) Assert(t testing.TB, header spec.Header, hdr GenHeader) bool {
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
	if !assert.EqualValues(t, header.Enum, hdr.Enum) {
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
