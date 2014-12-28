package swagger

import (
	"io"
	"net/http"
	"sort"
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

func TestUntypedAppValidation(t *testing.T) {
	formatParam := swagger.QueryParam()
	formatParam.Name = "format"
	formatParam.Type = "string"

	limitParam := swagger.QueryParam()
	limitParam.Name = "limit"
	limitParam.Type = "integer"
	limitParam.Format = "int32"
	limitParam.Extensions = swagger.Extensions(map[string]interface{}{})
	limitParam.Extensions.Add("go-name", "Limit")

	skipParam := swagger.QueryParam()
	skipParam.Name = "skip"
	skipParam.Type = "integer"
	skipParam.Format = "int32"

	spec := &swagger.Spec{
		Consumes:   []string{"application/json"},
		Produces:   []string{"application/json"},
		Parameters: map[string]swagger.Parameter{"format": *formatParam},
		Paths: swagger.Paths{
			Paths: map[string]swagger.PathItem{
				"/": swagger.PathItem{
					Parameters: []swagger.Parameter{*limitParam},
					Get: &swagger.Operation{
						Consumes:   []string{"application/x-yaml"},
						Produces:   []string{"application/x-yaml"},
						ID:         "someOperation",
						Parameters: []swagger.Parameter{*skipParam},
					},
				},
			},
		},
	}
	analyzer := newAnalyzer(spec)

	api1 := NewAPI(spec)
	err := api1.validateWith(analyzer)
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] consumes registrations", err.Error())

	api1.RegisterConsumer("application/x-yaml", new(stubConsumer))
	err = api1.validateWith(analyzer)
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] produces registrations", err.Error())

	api1.RegisterProducer("application/x-yaml", new(stubProducer))
	err = api1.validateWith(analyzer)
	assert.Error(t, err)
	assert.Equal(t, "missing [someOperation] operation registrations", err.Error())

	api1.RegisterOperation("someOperation", new(stubOperationHandler))
	err = api1.validateWith(analyzer)
	assert.NoError(t, err)

	api1.RegisterConsumer("application/something", new(stubConsumer))
	err = api1.validateWith(analyzer)
	assert.Error(t, err)
	assert.Equal(t, "missing from spec file [application/something] consumes", err.Error())

	api2 := NewAPI(spec)
	api2.RegisterConsumer("application/something", new(stubConsumer))
	err = api2.validateWith(analyzer)
	assert.Error(t, err)
	assert.Equal(t, "missing [application/x-yaml] consumes registrations\nmissing from spec file [application/something] consumes", err.Error())

	expected := []string{"application/json", "application/x-yaml"}
	sort.Sort(sort.StringSlice(expected))
	consumes := analyzer.ConsumesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(consumes))
	assert.Equal(t, expected, consumes)
	consumers := api1.ConsumersFor(consumes)
	assert.Len(t, consumers, 2)

	produces := analyzer.ProducesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(produces))
	assert.Equal(t, expected, produces)
	producers := api1.ProducersFor(produces)
	assert.Len(t, producers, 2)

	parameters := analyzer.ParametersFor(spec.Paths.Paths["/"].Get)
	assert.Len(t, parameters, 3)
}
