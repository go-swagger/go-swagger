package spec

import (
	"sort"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, analyzer.operations["someOperation"], "/")

	expected := []string{"application/json", "application/x-yaml"}
	sort.Sort(sort.StringSlice(expected))
	consumes := analyzer.ConsumesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(consumes))
	assert.Equal(t, expected, consumes)

	produces := analyzer.ProducesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(produces))
	assert.Equal(t, expected, produces)

	parameters := analyzer.ParametersFor(spec.Paths.Paths["/"].Get)
	assert.Len(t, parameters, 3)

	assert.Equal(t, spec.Paths.Paths, analyzer.AllPaths())
}
