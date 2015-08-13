package generator

import (
	"testing"

	"github.com/go-swagger/go-swagger/spec"
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
			res, err := b.MakeResponse("a", "default", false, resolver, resp)
			if assert.NoError(t, err) {
				if assertResponse(t, resp, res, false) {
					for code, response := range op.Responses.StatusCodeResponses {
						isSuccess := code/100 == 2
						res, err := b.MakeResponse("a", "default", isSuccess, resolver, response)
						if assert.NoError(t, err) {
							pretty.Println(res)
							assertResponse(t, response, res, isSuccess)
						}
					}
				}
			}
		}
	}

	//b, err = opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")

	//if !assert.NoError(t, err) {
	//t.FailNow()
	//}

	//op, ok = b.Doc.OperationForName("updateTask")
	//if assert.True(t, ok) && assert.NotNil(t, op) && assert.NotNil(t, op.Responses) {
	//resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	//if assert.NotNil(t, op.Responses.Default) {
	//resp := *op.Responses.Default
	//res, err := b.MakeResponse("a", "default", false, resolver, resp)
	//if assert.NoError(t, err) {
	//pretty.Println(res)
	//if assertResponse(t, resp, res, false) {
	//for code, response := range op.Responses.StatusCodeResponses {
	//isSuccess := code/100 == 2
	//res, err := b.MakeResponse("a", "default", isSuccess, resolver, response)
	//if assert.NoError(t, err) {
	//assertResponse(t, response, res, isSuccess)
	//}
	//}
	//}
	//}
	//}
	//}
}

type responseTestContext struct {
	OpID      string
	Name      string
	IsSuccess string
}

func assertResponse(t testing.TB, response spec.Response, res GenResponse, isSuccess bool) bool {
	if !assert.Equal(t, isSuccess, res.IsSuccess) {
		return false
	}

	//if !assert.Equal(t, response.Description, res.Description) {
	//return false
	//}

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
