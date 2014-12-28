package spec

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/stretchr/testify/assert"
)

func TestUnknownSpecVersion(t *testing.T) {
	_, err := New([]byte{}, "0.9")
	assert.Error(t, err)
}

func TestDefaultsTo20(t *testing.T) {
	d, err := New(PetStoreJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, "2.0", d.Version())
	assert.Equal(t, "2.0", d.data["swagger"].(string))
}

func TestValidatesValidSchema(t *testing.T) {
	d, err := New(PetStoreJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)
	res := d.Validate()
	assert.NotNil(t, res)
	assert.True(t, res.Valid())
	assert.Empty(t, res.Errors())

}

func TestFailsInvalidSchema(t *testing.T) {
	d, err := New(InvalidJSONMessage, "")

	assert.NoError(t, err)
	assert.NotNil(t, d)

	res := d.Validate()
	assert.NotNil(t, res)
	assert.False(t, res.Valid())
	assert.NotEmpty(t, res.Errors())
}

func TestFailsInvalidJSON(t *testing.T) {
	_, err := New(json.RawMessage([]byte("{]")), "")

	assert.Error(t, err)
}

func TestSpecUtilityMethods(t *testing.T) {
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
	js, _ := spec.MarshalJSON()
	d, err := New(json.RawMessage(js), "")
	assert.NoError(t, err)

	expected := []string{"application/json", "application/x-yaml"}
	sort.Sort(sort.StringSlice(expected))
	consumes := d.ConsumesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(consumes))
	assert.Equal(t, expected, consumes)

	produces := d.ProducesFor(spec.Paths.Paths["/"].Get)
	sort.Sort(sort.StringSlice(produces))
	assert.Equal(t, expected, produces)

	parameters := d.ParametersFor(spec.Paths.Paths["/"].Get)
	assert.Len(t, parameters, 3)

	assert.Equal(t, d.spec.Paths.Paths, d.AllPaths())
}
