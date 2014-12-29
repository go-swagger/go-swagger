package petstore

import (
	"io"
	gotest "testing"

	"github.com/casualjim/go-swagger/swagger"
	"github.com/casualjim/go-swagger/swagger/spec"
	testingutil "github.com/casualjim/go-swagger/swagger/testing"
	"github.com/stretchr/testify/assert"
)

// NewAPI registers a stub api for the pet store
func NewAPI(t *gotest.T) (*spec.Document, *swagger.API) {
	spec, err := spec.New(testingutil.PetStoreJSONMessage, "")
	assert.NoError(t, err)
	api := swagger.NewAPI(spec)

	api.RegisterConsumer("application/json", new(stubConsumer))
	api.RegisterProducer("application/json", new(stubProducer))
	api.RegisterConsumer("application/xml", new(stubConsumer))
	api.RegisterProducer("application/xml", new(stubProducer))
	api.RegisterProducer("text/plain", new(stubProducer))
	api.RegisterProducer("text/html", new(stubProducer))
	api.RegisterConsumer("application/x-yaml", new(stubConsumer))
	api.RegisterProducer("application/x-yaml", new(stubProducer))

	api.RegisterOperation("getAllPets", new(stubOperationHandler))
	api.RegisterOperation("createPet", new(stubOperationHandler))
	api.RegisterOperation("deletePet", new(stubOperationHandler))
	api.RegisterOperation("getPetById", new(stubOperationHandler))

	return spec, api
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
