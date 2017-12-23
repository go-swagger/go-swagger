package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test proper validation: items in object error
func TestCmd_Validate_Issue1238(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
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

// Test proper validation: missing items in array error
func TestCmd_Validate_Issue1171(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "1171", "swagger.yaml")
	result := v.Execute([]string{specDoc})
	assert.Error(t, result)
}

// Test proper validation: reference to inner property in schema
// NOTE: Open a dedicated issue on message relevance (repeated occurence of message).
func TestCmd_Validate_Issue342_ForbiddenProperty(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "342", "fixture-342.yaml")
	result := v.Execute([]string{specDoc})
	assert.Error(t, result)
}

/*
TODO: fixture 342-2 (a variant of invalid specification), cannot be tested because it calls log.Fatal()

// Test proper validation: reference to shared top level parameter, but with incorrect
// yaml syntax: use map key instead of array item.
// NOTE: this error message is not clear enough. The role of this test
// is to determine that the validation does not panic and correctly states the spec is invalid.
// Open a dedicated issue on message relevance. This test shall be updated with the finalized message.
func TestCmd_Validate_Issue342_CannotUnmarshal(t *testing.T) {
	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "342", "fixture-342-2.yaml")
	assert.NotPanics(t, func() {
		v.Execute([]string{specDoc})
	})
	// TODO: uncomment the following test when validation no more results in a log.Fatal() call
	//result := v.Execute([]string{specDoc})
	//if assert.Error(t, result, "This spec should not pass validation") {
	//	//assert.Contains(t, result.Error(), "is invalid against swagger specification 2.0")
	//	assert.Contains(t, result.Error(), "json: cannot unmarshal object into Go struct field SwaggerProps.paths of type []spec.Parameter")
	//}
}
*/

// This one is a correct version of issue#342 and it validates
func TestCmd_Validate_Issue342_Correct(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	v := ValidateSpec{}
	base := filepath.FromSlash("../../../")
	specDoc := filepath.Join(base, "fixtures", "bugs", "342", "fixture-342-3.yaml")
	result := v.Execute([]string{specDoc})
	assert.NoError(t, result)
}
