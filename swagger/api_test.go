package swagger

import (
	"io"
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

type stubAuthHandler struct {
}

func (s *stubAuthHandler) Authenticate(_ *http.Request) interface{} {
	return nil
}

type stubConsumer struct {
}

func (s *stubConsumer) Consume(_ io.Reader, _ interface{}) error {
	return nil
}

type stubProducer struct {
}

func (s *stubProducer) Produce(_ io.Writer, _ interface{}) error {
	return nil
}

type stubOperationHandler struct {
}

func (s *stubOperationHandler) ParameterModel() interface{} {
	return nil
}

func (s *stubOperationHandler) Handle(params interface{}) (interface{}, error) {
	return nil, nil
}

func TestUntypedAPIRegistrations(t *testing.T) {
	api := NewAPI(new(swagger.Spec))

	api.RegisterAuth("basic", new(stubAuthHandler))
	api.RegisterConsumer("application/yada", new(stubConsumer))
	api.RegisterProducer("application/yada-2", new(stubProducer))
	api.RegisterOperation("someId", new(stubOperationHandler))

	assert.NotEmpty(t, api.authHandlers)

	_, ok := api.authHandlers["BASIC"]
	assert.True(t, ok)
	_, ok = api.consumers["application/yada"]
	assert.True(t, ok)
	_, ok = api.producers["application/yada-2"]
	assert.True(t, ok)
	_, ok = api.consumers["application/json"]
	assert.True(t, ok)
	_, ok = api.producers["application/json"]
	assert.True(t, ok)
	_, ok = api.operations["someId"]
	assert.True(t, ok)

	h, ok := api.OperationHandlerFor("someId")
	assert.True(t, ok)
	assert.NotNil(t, h)

	_, ok = api.OperationHandlerFor("doesntExist")
	assert.False(t, ok)
}
