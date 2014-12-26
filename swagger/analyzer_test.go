package swagger

import (
	"sort"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzer(t *testing.T) {
	spec := &swagger.Spec{
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Paths: swagger.Paths{
			Paths: map[string]swagger.PathItem{
				"/": swagger.PathItem{
					Get: &swagger.Operation{
						Consumes: []string{"application/x-yaml"},
						Produces: []string{"application/x-yaml"},
						ID:       "someOperation",
					},
				},
			},
		},
	}
	analyzer := newAnalyzer(spec)

	assert.Len(t, analyzer.consumes, 2)
	assert.Len(t, analyzer.produces, 2)
	assert.Len(t, analyzer.operations, 1)
	assert.Equal(t, analyzer.operations["someOperation"], "/")

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
}
