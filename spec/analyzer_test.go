
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func schemeNames(schemes []SecurityRequirement) []string {
	var names []string
	for _, v := range schemes {
		names = append(names, v.Name)
	}
	sort.Sort(sort.StringSlice(names))
	return names
}

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
	op.Security = []map[string][]string{
		map[string][]string{"oauth2": []string{}},
		map[string][]string{"basic": nil},
	}
	op.ID = "someOperation"
	op.Parameters = []Parameter{*skipParam}
	pi.Get = op

	spec := &Swagger{
		swaggerProps: swaggerProps{
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Security: []map[string][]string{
				map[string][]string{"apikey": nil},
			},
			SecurityDefinitions: map[string]*SecurityScheme{
				"basic":  BasicAuth(),
				"apiKey": APIKeyAuth("api_key", "query"),
				"oauth2": OAuth2AccessToken("http://authorize.com", "http://token.com"),
			},
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

	expectedSchemes := []SecurityRequirement{SecurityRequirement{"oauth2", []string{}}, SecurityRequirement{"basic", nil}}
	schemes := analyzer.SecurityRequirementsFor(spec.Paths.Paths["/"].Get)
	assert.Equal(t, schemeNames(expectedSchemes), schemeNames(schemes))

	securityDefinitions := analyzer.SecurityDefinitionsFor(spec.Paths.Paths["/"].Get)
	assert.Equal(t, securityDefinitions["basic"], *spec.SecurityDefinitions["basic"])
	assert.Equal(t, securityDefinitions["oauth2"], *spec.SecurityDefinitions["oauth2"])

	parameters := analyzer.ParamsFor("GET", "/")
	assert.Len(t, parameters, 2)

	operations := analyzer.OperationIDs()
	assert.Len(t, operations, 1)

	producers := analyzer.RequiredProduces()
	assert.Len(t, producers, 2)
	consumers := analyzer.RequiredConsumes()
	assert.Len(t, consumers, 2)
	authSchemes := analyzer.RequiredSchemes()
	assert.Len(t, authSchemes, 3)

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
