package validate

import (
	"testing"

	"github.com/casualjim/go-swagger/internal/testing/petstore"
	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestValidateItems(t *testing.T) {
	doc, api := petstore.NewAPI(t)
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	res := validator.validateItems()
	assert.Empty(t, res.Errors)

	// in operation parameters
	sw := doc.Spec()
	sw.Paths.Paths["/pets"].Get.Parameters[0].Type = "array"
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items = spec.NewItems().Typed("string", "")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items = spec.NewItems().Typed("array", "")
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	sw.Paths.Paths["/pets"].Get.Parameters[0].Items.Items = spec.NewItems().Typed("string", "")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	// in global parameters
	sw.Parameters = make(map[string]spec.Parameter)
	sw.Parameters["other"] = *spec.SimpleArrayParam("other", "array", "csv")
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	pp := spec.SimpleArrayParam("other", "array", "")
	pp.Items = nil
	sw.Parameters["other"] = *pp
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in shared path object parameters
	doc, api = petstore.NewAPI(t)
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	sw = doc.Spec()

	pa := sw.Paths.Paths["/pets"]
	pa.Parameters = []spec.Parameter{*spec.SimpleArrayParam("another", "array", "csv")}
	sw.Paths.Paths["/pets"] = pa
	res = validator.validateItems()
	assert.Empty(t, res.Errors)

	pa = sw.Paths.Paths["/pets"]
	pp = spec.SimpleArrayParam("other", "array", "")
	pp.Items = nil
	pa.Parameters = []spec.Parameter{*pp}
	sw.Paths.Paths["/pets"] = pa
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in body param schema
	doc, api = petstore.NewAPI(t)
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	pa.Post.Parameters[0].Schema = spec.ArrayProperty(nil)
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in response headers
	doc, api = petstore.NewAPI(t)
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	rp := pa.Post.Responses.StatusCodeResponses[200]
	var hdr spec.Header
	hdr.Type = "array"
	rp.Headers = make(map[string]spec.Header)
	rp.Headers["X-YADA"] = hdr
	pa.Post.Responses.StatusCodeResponses[200] = rp
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)

	// in response schema
	doc, api = petstore.NewAPI(t)
	validator = NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	sw = doc.Spec()
	pa = sw.Paths.Paths["/pets"]
	rp = pa.Post.Responses.StatusCodeResponses[200]
	rp.Schema = spec.ArrayProperty(nil)
	pa.Post.Responses.StatusCodeResponses[200] = rp
	res = validator.validateItems()
	assert.NotEmpty(t, res.Errors)
}

func TestValidateUniqueSecurityScopes(t *testing.T) {
}

func TestValidateUniqueScopesSecurityDefinitions(t *testing.T) {
}

func TestValidateReferenced(t *testing.T) {
}

func TestValidateRequiredDefinitions(t *testing.T) {
	doc, api := petstore.NewAPI(t)
	validator := NewSpecValidator(spec.MustLoadSwagger20Schema(), api.Formats())
	validator.spec = doc
	res := validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	// properties
	sw := doc.Spec()
	def := sw.Definitions["Tag"]
	def.Required = append(def.Required, "type")
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)

	// pattern properties
	def.PatternProperties = make(map[string]spec.Schema)
	def.PatternProperties["ty.*"] = *spec.StringProperty()
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	def.PatternProperties = make(map[string]spec.Schema)
	def.PatternProperties["^ty.$"] = *spec.StringProperty()
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)

	// additional properties
	def.PatternProperties = nil
	def.AdditionalProperties = &spec.SchemaOrBool{Allows: true}
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.Empty(t, res.Errors)

	def.AdditionalProperties = &spec.SchemaOrBool{Allows: false}
	sw.Definitions["Tag"] = def
	res = validator.validateRequiredDefinitions()
	assert.NotEmpty(t, res.Errors)
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
