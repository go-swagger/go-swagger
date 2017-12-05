package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test proper validation: items in object error
func TestCmd_Validate_Issue1238(t *testing.T) {
	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "1238", "swagger.yaml")
	result := v.Execute([]string{specDoc})
	if assert.Error(t, result) {
		/*
			The swagger spec at "../../../fixtures/bugs/1238/swagger.yaml" is invalid against swagger specification 2.0. see errors :
				- definitions.RRSets in body must be of type array
		*/
		assert.Contains(t, result.Error(), "is invalid against swagger specification 2.0")
		assert.Contains(t, result.Error(), "definitions.RRSets in body must be of type array")
	}
}
