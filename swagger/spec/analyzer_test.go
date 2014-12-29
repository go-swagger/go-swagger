package spec

import (
	"sort"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

func newAnalyzer(spec *swagger.Spec) *specAnalyzer {
	a := &specAnalyzer{
		spec:        spec,
		consumes:    make(map[string]struct{}),
		produces:    make(map[string]struct{}),
		authSchemes: make(map[string]struct{}),
		operations:  make(map[string]map[string]*swagger.Operation),
	}
	a.initialize()
	return a
}

func TestAnalyzer(t *testing.T) {
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
