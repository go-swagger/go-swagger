package validate

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringValidator_Validate_Panic(t *testing.T) {
	var schemaJSON = `
{
    "properties": {
        "name": {
            "type": "string",
            "pattern": "^[A-Za-z]+$",
            "minLength": 1
        },
        "place": {
            "type": "string",
            "pattern": "^[A-Za-z]+$",
            "minLength": 1
        }
    },
    "required": [
        "name"
    ]
}`
	var inputJSON = `{"name": "Ivan"}`
	schema := new(spec.Schema)
	require.NoError(t, json.Unmarshal([]byte(schemaJSON), schema))
	var input map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(inputJSON), &input))
	input["place"] = json.Number("10")

	assert.Error(t, AgainstSchema(schema, input, strfmt.Default))
}
