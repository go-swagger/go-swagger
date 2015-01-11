package spec

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newAnalyzer(spec *Swagger) *specAnalyzer {
	a := &specAnalyzer{
		spec:        spec,
		consumes:    make(map[string]struct{}),
		produces:    make(map[string]struct{}),
		authSchemes: make(map[string]struct{}),
		operations:  make(map[string]map[string]*Operation),
	}
	a.initialize()
	return a
}

func TestAnalyzer(t *testing.T) {
	formatParam := QueryParam("format").Typed("string", "")

	limitParam := QueryParam("limit").Typed("integer", "int32")
	limitParam.Extensions = Extensions(map[string]interface{}{})
	limitParam.Extensions.Add("go-name", "Limit")

	skipParam := QueryParam("skip").Typed("integer", "int32")
	pi := PathItem{}
	pi.Parameters = []Parameter{*limitParam}

	op := &Operation{}
	op.Consumes = []string{"application/x-yaml"}
	op.Produces = []string{"application/x-yaml"}
	op.ID = "someOperation"
	op.Parameters = []Parameter{*skipParam}
	pi.Get = op

	spec := &Swagger{
		swaggerProps: swaggerProps{
			Consumes:   []string{"application/json"},
			Produces:   []string{"application/json"},
			Parameters: map[string]Parameter{"format": *formatParam},
			Paths: &Paths{
				Paths: map[string]PathItem{
					"/": pi,
				},
			},
		},
	}
	analyzer := newAnalyzer(spec)

	assert.Len(t, analyzer.consumes, 2)
	assert.Len(t, analyzer.produces, 2)
	assert.Len(t, analyzer.operations, 1)
	assert.Equal(t, analyzer.operations["GET"]["/"], spec.Paths.Paths["/"].Get)

	expected := []string{"application/json", "application/x-yaml"}
	sort.Sort(sort.StringSlice(expected))
	consumes := analyzer.ConsumesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(consumes))
	assert.Equal(t, expected, consumes)

	produces := analyzer.ProducesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(produces))
	assert.Equal(t, expected, produces)

	parameters := analyzer.ParamsFor("GET", "/")
	assert.Len(t, parameters, 3)

	operations := analyzer.OperationIDs()
	assert.Len(t, operations, 1)

	producers := analyzer.RequiredProduces()
	assert.Len(t, producers, 2)
	consumers := analyzer.RequiredConsumes()
	assert.Len(t, consumers, 2)

	ops := analyzer.Operations()
	assert.Len(t, ops, 1)
	assert.Len(t, ops["GET"], 1)

	op, ok := analyzer.OperationFor("get", "/")
	assert.True(t, ok)
	assert.NotNil(t, op)

	op, ok = analyzer.OperationFor("delete", "/")
	assert.False(t, ok)
	assert.Nil(t, op)
}
