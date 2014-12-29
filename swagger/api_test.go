package swagger

import (
	"io"
	"sort"
	"testing"

	"github.com/casualjim/go-swagger/swagger/spec"
	"github.com/stretchr/testify/assert"
)

// type stubAuthHandler struct {
// }

// func (s *stubAuthHandler) Authenticate(_ *http.Request) interface{} {
// 	return nil
// }

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
	api := NewAPI(new(spec.Document))

	// api.RegisterAuth("basic", new(stubAuthHandler))
	api.RegisterConsumer("application/yada", new(stubConsumer))
	api.RegisterProducer("application/yada-2", new(stubProducer))
	api.RegisterOperation("someId", new(stubOperationHandler))

	// assert.NotEmpty(t, api.authHandlers)

	// _, ok := api.authHandlers["BASIC"]
	// assert.True(t, ok)
	_, ok := api.consumers["application/yada"]
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

func TestUntypedAppValidation(t *testing.T) {
	specStr := `{
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "parameters": {
  	"format": {
  		"in": "query",
  		"name": "format",
  		"type": "string"
  	}
  },
  "paths": {
  	"/": {
  		"parameters": [
  			{
  				"name": "limit",
		  		"type": "integer",
		  		"format": "int32",
		  		"x-go-name": "Limit"
		  	}
  		],
  		"get": {
  			"consumes": ["application/x-yaml"],
  			"produces": ["application/x-yaml"],
  			"operationId": "someOperation",
  			"parameters": [
  				{
			  		"name": "skip",
			  		"type": "integer",
			  		"format": "int32"
			  	}
  			]
  		}
  	}
  }
}`
	spec, err := spec.New([]byte(specStr), "")
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	cons := spec.ConsumesFor(spec.AllPaths()["/"].Get)
	assert.Len(t, cons, 2)
	prods := spec.RequiredProduces()
	assert.Len(t, prods, 2)

	api1 := NewAPI(spec)
	err = api1.Validate()
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] consumes registrations", err.Error())
	api1.RegisterConsumer("application/x-yaml", new(stubConsumer))
	err = api1.validate()
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] produces registrations", err.Error())
	api1.RegisterProducer("application/x-yaml", new(stubProducer))
	err = api1.validate()
	assert.Error(t, err)
	assert.Equal(t, "missing [someOperation] operation registrations", err.Error())
	api1.RegisterOperation("someOperation", new(stubOperationHandler))
	err = api1.validate()
	assert.NoError(t, err)
	api1.RegisterConsumer("application/something", new(stubConsumer))
	err = api1.validate()
	assert.Error(t, err)
	assert.Equal(t, "missing from spec file [application/something] consumes", err.Error())

	api2 := NewAPI(spec)
	api2.RegisterConsumer("application/something", new(stubConsumer))
	err = api2.validate()
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] consumes registrations\nmissing from spec file [application/something] consumes", err.Error())

	expected := []string{"application/json", "application/x-yaml"}
	sort.Sort(sort.StringSlice(expected))
	consumes := spec.ConsumesFor(spec.AllPaths()["/"].Get)
	sort.Sort(sort.StringSlice(consumes))
	assert.Equal(t, expected, consumes)
	consumers := api1.ConsumersFor(consumes)
	assert.Len(t, consumers, 2)

	produces := spec.ProducesFor(spec.AllPaths()["/"].Get)
	sort.Sort(sort.StringSlice(produces))
	assert.Equal(t, expected, produces)
	producers := api1.ProducersFor(produces)
	assert.Len(t, producers, 2)

}
