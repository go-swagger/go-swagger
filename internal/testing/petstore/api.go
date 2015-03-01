package petstore

import (
	"io"
	gotest "testing"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	testingutil "github.com/casualjim/go-swagger/internal/testing"
	"github.com/casualjim/go-swagger/middleware/security"
	"github.com/casualjim/go-swagger/middleware/untyped"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

// NewAPI registers a stub api for the pet store
func NewAPI(t *gotest.T) (*spec.Document, *untyped.API) {
	spec, err := spec.New(testingutil.PetStoreJSONMessage, "")
	assert.NoError(t, err)
	api := untyped.NewAPI(spec)

	api.RegisterConsumer("application/json", swagger.JSONConsumer())
	api.RegisterProducer("application/json", swagger.JSONProducer())
	api.RegisterConsumer("application/xml", new(stubConsumer))
	api.RegisterProducer("application/xml", new(stubProducer))
	api.RegisterProducer("text/plain", new(stubProducer))
	api.RegisterProducer("text/html", new(stubProducer))
	api.RegisterConsumer("application/x-yaml", swagger.YAMLConsumer())
	api.RegisterProducer("application/x-yaml", swagger.YAMLProducer())

	api.RegisterAuth("basic", security.BasicAuth(func(username, password string) (interface{}, error) {
		if username == "admin" && password == "admin" {
			return "admin", nil
		}
		return nil, errors.Unauthenticated("basic")
	}))
	api.RegisterAuth("apiKey", security.APIKeyAuth("X-API-KEY", "header", func(token string) (interface{}, error) {
		if token == "token123" {
			return "admin", nil
		}
		return nil, errors.Unauthenticated("token")
	}))
	api.RegisterOperation("getAllPets", new(stubOperationHandler))
	api.RegisterOperation("createPet", new(stubOperationHandler))
	api.RegisterOperation("deletePet", new(stubOperationHandler))
	api.RegisterOperation("getPetById", new(stubOperationHandler))

	api.Models["pet"] = func() interface{} { return new(Pet) }
	api.Models["newPet"] = func() interface{} { return new(Pet) }
	api.Models["tag"] = func() interface{} { return new(Tag) }

	return spec, api
}

// NewRootAPI registers a stub api for the pet store
func NewRootAPI(t *gotest.T) (*spec.Document, *untyped.API) {
	spec, err := spec.New(testingutil.RootPetStoreJSONMessage, "")
	assert.NoError(t, err)
	api := untyped.NewAPI(spec)

	api.RegisterConsumer("application/json", swagger.JSONConsumer())
	api.RegisterProducer("application/json", swagger.JSONProducer())
	api.RegisterConsumer("application/xml", new(stubConsumer))
	api.RegisterProducer("application/xml", new(stubProducer))
	api.RegisterProducer("text/plain", new(stubProducer))
	api.RegisterProducer("text/html", new(stubProducer))
	api.RegisterConsumer("application/x-yaml", swagger.YAMLConsumer())
	api.RegisterProducer("application/x-yaml", swagger.YAMLProducer())

	api.RegisterAuth("basic", security.BasicAuth(func(username, password string) (interface{}, error) {
		if username == "admin" && password == "admin" {
			return "admin", nil
		}
		return nil, errors.Unauthenticated("basic")
	}))
	api.RegisterAuth("apiKey", security.APIKeyAuth("X-API-KEY", "header", func(token string) (interface{}, error) {
		if token == "token123" {
			return "admin", nil
		}
		return nil, errors.Unauthenticated("token")
	}))
	api.RegisterOperation("getAllPets", new(stubOperationHandler))
	api.RegisterOperation("createPet", new(stubOperationHandler))
	api.RegisterOperation("deletePet", new(stubOperationHandler))
	api.RegisterOperation("getPetById", new(stubOperationHandler))

	api.Models["pet"] = func() interface{} { return new(Pet) }
	api.Models["newPet"] = func() interface{} { return new(Pet) }
	api.Models["tag"] = func() interface{} { return new(Tag) }

	return spec, api
}

// Tag the tag model
type Tag struct {
	ID   int64
	Name string
}

// Pet the pet model
type Pet struct {
	ID        int64
	Name      string
	PhotoURLs []string
	Status    string
	Tags      []Tag
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
