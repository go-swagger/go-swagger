package validate

import (
	"testing"

	"github.com/casualjim/go-swagger/internal/testing/petstore"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestValidateItems(t *testing.T) {

}

func TestValidateUniqueSecurityScopes(t *testing.T) {
}

func TestValidateUniqueScopesSecurityDefinitions(t *testing.T) {
}

func TestValidateReferenced(t *testing.T) {
}

func TestValidateRequiredDefinitions(t *testing.T) {
}

func TestValidateParameters(t *testing.T) {
	doc, api := petstore.NewAPI(t)
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	res := validator.validateParameters()
	assert.Empty(t, res.Errors)

	sw := doc.Spec()
	sw.Paths.Paths["/pets"].Get.Parameters = append(sw.Paths.Paths["/pets"].Get.Parameters, *spec.QueryParam("limit").Typed("string", ""))
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)

	doc, api = petstore.NewAPI(t)
	sw = doc.Spec()
	sw.Paths.Paths["/pets"].Post.Parameters = append(sw.Paths.Paths["/pets"].Post.Parameters, *spec.BodyParam("fake", spec.RefProperty("#/definitions/Pet")))
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 1)
	assert.Contains(t, res.Errors[0].Error(), "has more than 1 body param")

	doc, api = petstore.NewAPI(t)
	sw = doc.Spec()
	pp := sw.Paths.Paths["/pets/{id}"]
	pp.Delete = nil
	var nameParams []spec.Parameter
	for _, p := range pp.Parameters {
		if p.Name == "id" {
			p.Name = "name"
			nameParams = append(nameParams, p)
		}
	}
	pp.Parameters = nameParams
	sw.Paths.Paths["/pets/{name}"] = pp
	doc.Reload()
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 1)
	assert.Contains(t, res.Errors[0].Error(), "overlaps with")

	doc, api = petstore.NewAPI(t)
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	sw = doc.Spec()
	pp = sw.Paths.Paths["/pets/{id}"]
	pp.Delete = nil
	pp.Get.Parameters = nameParams
	pp.Parameters = nil
	sw.Paths.Paths["/pets/{id}"] = pp
	doc.Reload()
	res = validator.validateParameters()
	assert.NotEmpty(t, res.Errors)
	assert.Len(t, res.Errors, 2)
	assert.Contains(t, res.Errors[1].Error(), "is not present in the path")
	assert.Contains(t, res.Errors[0].Error(), "has no parameter definition")
}

func TestValidateReferencesValid(t *testing.T) {
}

func TestValidateDefaultValueAgainstSchema(t *testing.T) {
}
