package commands

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test proper validation: missing items in array error
func TestCmd_Validate_Issue1171(t *testing.T) {
	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "1171", "swagger.yaml")
	result := v.Execute([]string{specDoc})
	if assert.Error(t, result) {
		assert.Contains(t, result.Error(), "is invalid against swagger specification 2.0")
		assert.Contains(t, result.Error(), "items in definitions.Zones is required")
	}
}
