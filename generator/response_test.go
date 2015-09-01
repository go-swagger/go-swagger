package generator

import (
	"fmt"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/kr/pretty"
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
			res, err := b.MakeResponse("a", defCtx.Name, false, resolver, resp)
			if assert.NoError(t, err) {
				if defCtx.Assert(t, resp, res) {
					for code, response := range op.Responses.StatusCodeResponses {
						sucCtx := responseTestContext{
							OpID:      "updateTask",
							Name:      "success",
							IsSuccess: code/100 == 2,
						}
						res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, response)
						if assert.NoError(t, err) {
							sucCtx.Assert(t, response, res)
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
	pretty.Println(res)
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
					ctx.assertHeader(t, k, v, h)
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

func (ctx *respHeaderTestContext) Assert(t testing.TB, name string, header spec.Header, hdr GenHeader) bool {
	if !assert.Equal(t, name, hdr.Name) {
		return false
	}
	if !assert.Equal(t, fmt.Sprintf("%q", name), hdr.Path) {
		return false
	}
	if !assert.Equal(t, "i", hdr.IndexVar) {
		return false
	}
	if !assert.Equal(t, "a", hdr.ReceiverName) {
		return false
	}
	if !assert.Equal(t, "a."+swag.ToGoName(param.Name), hdr.ValueExpression) {
		return false
	}
	if !assert.Equal(t, ctx.Formatter, hdr.Formatter) {
		return false
	}
	if !assert.Equal(t, ctx.Converter, hdr.Converter) {
		return false
	}
	if !assert.Equal(t, param.Description, hdr.Description) {
		return false
	}
	if !assert.Equal(t, param.CollectionFormat, hdr.CollectionFormat) {
		return false
	}
	if !assert.Equal(t, param.Required, hdr.Required) {
		return false
	}
	if !assert.Equal(t, param.Minimum, hdr.Minimum) || !assert.Equal(t, param.ExclusiveMinimum, hdr.ExclusiveMinimum) {
		return false
	}
	if !assert.Equal(t, param.Maximum, hdr.Maximum) || !assert.Equal(t, param.ExclusiveMaximum, hdr.ExclusiveMaximum) {
		return false
	}
	if !assert.Equal(t, param.MinLength, hdr.MinLength) {
		return false
	}
	if !assert.Equal(t, param.MaxLength, hdr.MaxLength) {
		return false
	}
	if !assert.Equal(t, param.Pattern, hdr.Pattern) {
		return false
	}
	if !assert.Equal(t, param.MaxItems, hdr.MaxItems) {
		return false
	}
	if !assert.Equal(t, param.MinItems, hdr.MinItems) {
		return false
	}
	if !assert.Equal(t, param.UniqueItems, hdr.UniqueItems) {
		return false
	}
	if !assert.Equal(t, param.MultipleOf, hdr.MultipleOf) {
		return false
	}
	if !assert.EqualValues(t, param.Enum, hdr.Enum) {
		return false
	}
	if !assert.Equal(t, param.Type, hdr.SwaggerType) {
		return false
	}
	if !assert.Equal(t, param.Format, hdr.SwaggerFormat) {
		return false
	}
	return true
}
