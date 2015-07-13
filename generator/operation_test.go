package generator

import (
	"errors"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestMakeResponseHeader(t *testing.T) {
	b, err := opBuilder("getTasks")
	if assert.NoError(t, err) {
		hdr := findResponseHeader(&b.Operation, 200, "X-Rate-Limit")
		gh := b.MakeHeader("a", "X-Rate-Limit", *hdr)
		assert.True(t, gh.IsPrimitive)
		assert.Equal(t, "int32", gh.GoType)
		assert.Equal(t, "X-Rate-Limit", gh.Name)
	}
}

func TestMakeResponse(t *testing.T) {
	b, err := opBuilder("getTasks")
	if assert.NoError(t, err) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		gO, err := b.MakeResponse("a", "getTasksSuccess", true, resolver, b.Operation.Responses.StatusCodeResponses[200])
		if assert.NoError(t, err) {
			assert.Len(t, gO.Headers, 2)
			assert.NotNil(t, gO.Schema)
			assert.True(t, gO.Schema.IsArray)
			assert.NotNil(t, gO.Schema.Items)
			assert.False(t, gO.Schema.IsAnonymous)
			assert.Equal(t, "[]models.Task", gO.Schema.GoType)
		}
	}
}

func opBuilder(name string) (codeGenOpBuilder, error) {

	specDoc, err := spec.Load("../fixtures/codegen/todolist.simple.yml")
	if err != nil {
		return codeGenOpBuilder{}, err
	}

	op, ok := specDoc.OperationForName(name)
	if !ok {
		return codeGenOpBuilder{}, errors.New("No operation could be found for simpleHeaderParams")
	}

	return codeGenOpBuilder{
		Name:          name,
		APIPackage:    "restapi",
		ModelsPackage: "models",
		Principal:     "models.User",
		Target:        ".",
		Operation:     *op,
		Doc:           specDoc,
		Authed:        false,
		ExtraSchemas:  make(map[string]GenSchema),
	}, nil
}

func findResponseHeader(op *spec.Operation, code int, name string) *spec.Header {
	resp := op.Responses.Default
	if code > 0 {
		bb, ok := op.Responses.StatusCodeResponses[code]
		if ok {
			resp = &bb
		}
	}

	if resp == nil {
		return nil
	}

	hdr, ok := resp.Headers[name]
	if !ok {
		return nil
	}

	return &hdr
}
