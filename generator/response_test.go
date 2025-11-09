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

package generator

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-openapi/spec"
)

func TestSimpleResponseRender(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))

	ff, err := opts.LanguageOpts.FormatContent("update_task_responses.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "o.XErrorCode", string(ff))
	assertInCode(t, "o.Payload", string(ff))
}

func TestDefaultResponseRender(t *testing.T) {
	b, err := opBuilder("getAllParameters", "../fixtures/codegen/todolist.responses.yml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op))

	ff, err := opts.LanguageOpts.FormatContent("get_all_parameters_responses.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "type GetAllParametersDefault struct", res)
	assertInCode(t, `if response.Code()/100 == 2`, res)
	assertNotInCode(t, `switch response.Code()`, res)
	assertNotInCode(t, "o.Payload", res)
}

func TestSimpleResponses(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")
	require.NoError(t, err)

	_, _, op, ok := b.Analyzed.OperationForName("updateTask")
	require.True(t, ok)
	require.NotNil(t, op)
	require.NotNil(t, op.Responses)

	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	require.NotNil(t, op.Responses.Default)

	resp, err := spec.ResolveResponse(b.Doc.Spec(), op.Responses.Default.Ref)
	require.NoError(t, err)

	defCtx := responseTestContext{
		OpID: "updateTask",
		Name: "default",
	}
	res, err := b.MakeResponse("a", defCtx.Name, false, resolver, -1, *resp)
	require.NoError(t, err)

	defCtx.Require(t, *resp, res)

	for code, response := range op.Responses.StatusCodeResponses {
		sucCtx := responseTestContext{
			OpID:      "updateTask",
			Name:      "success",
			IsSuccess: code/100 == 2,
		}
		res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, code, response)
		require.NoError(t, err)
		sucCtx.Require(t, response, res)
	}
}

func TestInlinedSchemaResponses(t *testing.T) {
	b, err := opBuilder("getTasks", "../fixtures/codegen/todolist.responses.yml")
	require.NoError(t, err)

	_, _, op, ok := b.Analyzed.OperationForName("getTasks")
	require.True(t, ok)
	require.NotNil(t, op)
	require.NotNil(t, op.Responses)

	resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
	require.NotNil(t, op.Responses.Default)

	resp := *op.Responses.Default
	defCtx := responseTestContext{
		OpID: "getTasks",
		Name: "default",
	}
	res, err := b.MakeResponse("a", defCtx.Name, false, resolver, -1, resp)
	require.NoError(t, err)

	defCtx.Require(t, resp, res)

	for code, response := range op.Responses.StatusCodeResponses {
		sucCtx := responseTestContext{
			OpID:      "getTasks",
			Name:      "success",
			IsSuccess: code/100 == 2,
		}
		res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, code, response)
		require.NoError(t, err)
		sucCtx.Require(t, response, res)
		assert.Len(t, b.ExtraSchemas, 1)
		// ExtraSchema is not a definition: it is rendered in current operations package
		assert.Equal(t, "[]*SuccessBodyItems0", res.Schema.GoType)
	}
}

func TestGenResponses_Issue540(t *testing.T) {
	b, err := opBuilder("postPet", "../fixtures/bugs/540/swagger.yml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))

	ff, err := opts.LanguageOpts.FormatContent("post_pet_responses.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "func (o *PostPetOK) WithPayload(payload models.Pet) *PostPetOK {", string(ff))
	assertInCode(t, "func (o *PostPetOK) SetPayload(payload models.Pet) {", string(ff))
}

func TestGenResponses_Issue718_NotRequired(t *testing.T) {
	b, err := opBuilder("doEmpty", "../fixtures/codegen/todolist.simple.yml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))

	ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "if payload == nil", string(ff))
	assertInCode(t, "payload = make([]models.Foo, 0, 50)", string(ff))
}

func TestGenResponses_Issue718_Required(t *testing.T) {
	t.Run("should prepare operation builder", func(t *testing.T) {
		b, err := opBuilder("doEmpty", "../fixtures/codegen/todolist.simple.yml")
		require.NoError(t, err)

		t.Run("should build operation", func(t *testing.T) {
			op, err := b.MakeOperation()
			require.NoError(t, err)

			t.Run("should generate response", func(t *testing.T) {
				var buf bytes.Buffer
				opts := opts()
				require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))

				t.Run("response code should be properly go-formatted", func(t *testing.T) {
					ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
					if err != nil {
						t.Logf("bad formatting: %v\n%s", err, buf.String())
					}
					require.NoError(t, err)

					t.Run("response code should be guard against nil payload", func(t *testing.T) {
						assertInCode(t, "if payload == nil", string(ff))
					})

					t.Run("response code should preallocate payload", func(t *testing.T) {
						assertInCode(t, "payload = make([]models.Foo, 0, 50)", string(ff))
					})
				})
			})
		})
	})
}

// Issue776 includes references that span multiple files. Flattening or Expanding is required.
func TestGenResponses_Issue776_Spec(t *testing.T) {
	defer discardOutput()()

	b, err := opBuilderWithFlatten("GetItem", "../fixtures/bugs/776/spec.yaml")
	require.NoError(t, err)
	op, err := b.MakeOperation()
	require.NoError(t, err)
	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
	ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
	if err != nil {
		t.Logf("bad formatting: %v\n%s", err, buf.String())
	}
	require.NoError(t, err)
	// This should be models.Item if flat works correctly
	assertInCode(t, "Payload *models.Item", string(ff))
	assertNotInCode(t, "type GetItemOKBody struct", string(ff))
}

func TestGenResponses_Issue776_SwaggerTemplate(t *testing.T) {
	defer discardOutput()()

	b, err := opBuilderWithFlatten("getHealthy", "../fixtures/bugs/776/swagger-template.yml")
	require.NoError(t, err)
	op, err := b.MakeOperation()
	require.NoError(t, err)
	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
	ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
	if err != nil {
		t.Logf("bad formatting: %v\n%s", err, buf.String())
	}
	require.NoError(t, err)
	assertInCode(t, "Payload *models.Error", string(ff))
}

func TestIssue846(t *testing.T) {
	// do it 8 times, to ensure it's always in the same order
	for range 8 {
		b, err := opBuilder("getFoo", "../fixtures/bugs/846/swagger.yml")
		require.NoError(t, err)
		op, err := b.MakeOperation()
		require.NoError(t, err)
		var buf bytes.Buffer
		opts := opts()
		require.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op))
		ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
		if err != nil {
			t.Logf("bad formatting: %v\n%s", err, buf.String())
		}
		require.NoError(t, err)
		// sorted by code
		assert.Regexp(t, "(?s)"+
			"GetFooOK struct.+"+
			"GetFooNotFound struct.+"+
			"GetFooInternalServerError struct", string(ff))
		// sorted by name
		assert.Regexp(t, "(?s)"+
			"GetFooInternalServerErrorBody struct.+"+
			"GetFooNotFoundBody struct.+"+
			"GetFooOKBody struct", string(ff))
	}
}

func TestIssue881(t *testing.T) {
	b, err := opBuilder("getFoo", "../fixtures/bugs/881/swagger.yml")
	require.NoError(t, err)
	op, err := b.MakeOperation()
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
}

func TestIssue881Deep(t *testing.T) {
	b, err := opBuilder("getFoo", "../fixtures/bugs/881/deep.yml")
	require.NoError(t, err)
	op, err := b.MakeOperation()
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
}

func TestGenResponses_XGoName(t *testing.T) {
	b, err := opBuilder("putTesting", "../fixtures/specs/response_name.json")
	require.NoError(t, err)
	op, err := b.MakeOperation()
	require.NoError(t, err)
	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
	ff, err := opts.LanguageOpts.FormatContent("put_testing_responses.go", buf.Bytes())
	if err != nil {
		t.Logf("bad formatting: %v\n%s", err, buf.String())
	}
	require.NoError(t, err)
	assertInCode(t, "const PutTestingAlternateNameCode int =", string(ff))
	assertInCode(t, "type PutTestingAlternateName struct {", string(ff))
	assertInCode(t, "func NewPutTestingAlternateName() *PutTestingAlternateName {", string(ff))
	assertInCode(t, "func (o *PutTestingAlternateName) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {", string(ff))
}

func TestGenResponses_Issue892(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/media/search", "../fixtures/bugs/982/swagger.yaml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opts := opts()
	require.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op))

	ff, err := opts.LanguageOpts.FormatContent("get_media_search_responses.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	assertInCode(t, "o.Media = aO0", string(ff))
}

func TestGenResponses_Issue1013(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/test", "../fixtures/bugs/1013/fixture-1013.yaml")
	require.NoError(t, err)

	op, err := b.MakeOperation()
	require.NoError(t, err)

	var buf bytes.Buffer
	opt := opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))

	ff, err := opt.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())
	assertInCode(t, "Payload *models.Response `json:\"body,omitempty\"`", string(ff))

	buf.Reset()
	b, err = methodPathOpBuilder("get", "/test2", "../fixtures/bugs/1013/fixture-1013.yaml")
	require.NoError(t, err)

	op, err = b.MakeOperation()
	require.NoError(t, err)

	opt = opts()
	require.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
	ff, err = opt.LanguageOpts.FormatContent("foo.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())
	assertInCode(t, "Payload *models.Response `json:\"body,omitempty\"`", string(ff))
}

func TestGenResponse_15362_WithExpand(t *testing.T) {
	defer discardOutput()()

	fixtureConfig := map[string]map[string][]string{
		// load expectations for parameters in operation get_nested_required_responses.go
		"getNestedRequired": { // fixture index
			"serverResponses": { // executed template
				// expected code lines
				`const GetNestedRequiredOKCode int = 200`,
				`type GetNestedRequiredOK struct {`,
				"	Payload [][][][]*GetNestedRequiredOKBodyItems0 `json:\"body,omitempty\"`",
				`func NewGetNestedRequiredOK() *GetNestedRequiredOK {`,
				`	return &GetNestedRequiredOK{`,
				`func (o *GetNestedRequiredOK) WithPayload(payload [][][][]*GetNestedRequiredOKBodyItems0) *GetNestedRequiredOK {`,
				`func (o *GetNestedRequiredOK) SetPayload(payload [][][][]*GetNestedRequiredOKBodyItems0) {`,
				`func (o *GetNestedRequiredOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
			},
			"serverOperation": { // executed template
				// expected code lines
				`type GetNestedRequiredOKBodyItems0 struct {`,
				"	Pkcs *string `json:\"pkcs\"`",
				`func (o *GetNestedRequiredOKBodyItems0) Validate(formats strfmt.Registry) error {`,
				`	if err := o.validatePkcs(formats); err != nil {`,
				`		return errors.CompositeValidationError(res...`,
				`func (o *GetNestedRequiredOKBodyItems0) validatePkcs(formats strfmt.Registry) error {`,
				`	if err := validate.Required("pkcs", "body", o.Pkcs); err != nil {`,
			},
		},
	}

	// assertParams also works for responses
	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1536", "fixture-1536-2-responses.yaml"), true, false)
}

func TestGenResponse_1572(t *testing.T) {
	defer discardOutput()()

	const genNullResponseDocComment = `/*
GetNullRequestProcessed OK`
	// testing fixture-1572.yaml with minimal flatten
	// edge cases for operations schemas

	/*
	        Run the following test caes and exercise the minimal flatten mode:
	   - [x] nil schema in body param / response
	   - [x] any in body param /response
	   - [x] additional schema reused from model (body param and response) (with maps or arrays)
	   - [x] primitive body / response
	   - [x] $ref'ed response and param (check that minimal flatten expands it)

	*/

	fixtureConfig := map[string]map[string][]string{
		// load expectations for responses in operation get_interface_responses.go
		"getInterface": { // fixture index
			"serverResponses": { // executed template
				// expected code lines
				`const GetInterfaceOKCode int = 200`,
				`type GetInterfaceOK struct {`,
				"	Payload any `json:\"body,omitempty\"`",
				`func NewGetInterfaceOK() *GetInterfaceOK {`,
				`	return &GetInterfaceOK{`,
				`func (o *GetInterfaceOK) WithPayload(payload any) *GetInterfaceOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetInterfaceOK) SetPayload(payload any) {`,
				`	o.Payload = payload`,
				`func (o *GetInterfaceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.WriteHeader(200`,
				`	payload := o.Payload`,
				`	if err := producer.Produce(rw, payload); err != nil {`,
			},
		},

		// load expectations for responses in operation get_primitive_responses.go
		"getPrimitive": { // fixture index
			"serverResponses": { // executed template
				// expected code lines
				`const GetPrimitiveOKCode int = 200`,
				`type GetPrimitiveOK struct {`,
				"	Payload float32 `json:\"body,omitempty\"`",
				`func NewGetPrimitiveOK() *GetPrimitiveOK {`,
				`	return &GetPrimitiveOK{`,
				`func (o *GetPrimitiveOK) WithPayload(payload float32) *GetPrimitiveOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetPrimitiveOK) SetPayload(payload float32) {`,
				`	o.Payload = payload`,
				`func (o *GetPrimitiveOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.WriteHeader(200`,
				`	payload := o.Payload`,
				`	if err := producer.Produce(rw, payload); err != nil {`,
			},
		},

		// load expectations for responses in operation get_null_responses.go
		"getNull": { // fixture index
			"serverResponses": { // executed template
				// expected code lines
				`const GetNullOKCode int = 200`,
				`type GetNullOK struct {`,
				"	Payload any `json:\"body,omitempty\"`",
				`func NewGetNullOK() *GetNullOK {`,
				`	return &GetNullOK{`,
				`func (o *GetNullOK) WithPayload(payload any) *GetNullOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetNullOK) SetPayload(payload any) {`,
				`	o.Payload = payload`,
				`func (o *GetNullOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.WriteHeader(200`,
				`	payload := o.Payload`,
				`	if err := producer.Produce(rw, payload); err != nil {`,
				`const GetNullRequestProcessedCode int = 203`,
				genNullResponseDocComment,
				`swagger:response getNullRequestProcessed`,
				`*/`,
				`type GetNullRequestProcessed struct {`,
				`func NewGetNullRequestProcessed() *GetNullRequestProcessed {`,
				`	return &GetNullRequestProcessed{`,
				`func (o *GetNullRequestProcessed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.Header().Del(runtime.HeaderContentType`,
				`	rw.WriteHeader(203`,
			},
		},

		// load expectations for responses in operation get_model_interface_responses.go
		"getModelInterface": { // fixture index
			"serverResponses": { // executed template
				// expected code lines
				`const GetModelInterfaceOKCode int = 200`,
				`type GetModelInterfaceOK struct {`,
				"	Payload []models.ModelInterface `json:\"body,omitempty\"`",
				`func NewGetModelInterfaceOK() *GetModelInterfaceOK {`,
				`	return &GetModelInterfaceOK{`,
				`func (o *GetModelInterfaceOK) WithPayload(payload []models.ModelInterface) *GetModelInterfaceOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetModelInterfaceOK) SetPayload(payload []models.ModelInterface) {`,
				`	o.Payload = payload`,
				`func (o *GetModelInterfaceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.WriteHeader(200`,
				`	payload := o.Payload`,
				`	if payload == nil {`,
				`		payload = make([]models.ModelInterface, 0, 50`,
				`	if err := producer.Produce(rw, payload); err != nil {`,
			},
		},
	}

	// assertParams also works for responses
	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "enhancements", "1572", "fixture-1572.yaml"), true, false)
}

func TestGenResponse_1893(t *testing.T) {
	defer discardOutput()()

	fixtureConfig := map[string]map[string][]string{
		// load expectations for parameters in operation get_nested_required_responses.go
		"getRecords": { // fixture index
			"clientResponse": { // executed template
				// expected code lines
				`type GetRecordsStatus512 struct {`,
				`type GetRecordsStatus515 struct {`,
			},
			"serverResponses": { // executed template
				// expected code lines
				`type GetRecordsStatus512 struct {`,
				`type GetRecordsStatus515 struct {`,
			},
		},
	}
	// assertParams also works for responses
	assertParams(t, fixtureConfig, filepath.Join("..", "fixtures", "bugs", "1893", "fixture-1893.yaml"), true, false)
}
