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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestSimpleResponseRender(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("update_task_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "o.XErrorCode", string(ff))
					assertInCode(t, "o.Payload", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestDefaultResponseRender(t *testing.T) {
	b, err := opBuilder("getAllParameters", "../fixtures/codegen/todolist.responses.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("get_all_parameters_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					res := string(ff)
					assertInCode(t, "type GetAllParametersDefault struct", res)
					assertInCode(t, `if response.Code()/100 == 2`, res)
					assertNotInCode(t, `switch response.Code()`, res)
					assertNotInCode(t, "o.Payload", res)
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestSimpleResponses(t *testing.T) {
	b, err := opBuilder("updateTask", "../fixtures/codegen/todolist.responses.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, _, op, ok := b.Analyzed.OperationForName("updateTask")
	if assert.True(t, ok) && assert.NotNil(t, op) && assert.NotNil(t, op.Responses) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		if assert.NotNil(t, op.Responses.Default) {
			resp, err := spec.ResolveResponse(b.Doc.Spec(), op.Responses.Default.Ref)
			if assert.NoError(t, err) {
				defCtx := responseTestContext{
					OpID: "updateTask",
					Name: "default",
				}
				res, err := b.MakeResponse("a", defCtx.Name, false, resolver, -1, *resp)
				if assert.NoError(t, err) {
					if defCtx.Assert(t, *resp, res) {
						for code, response := range op.Responses.StatusCodeResponses {
							sucCtx := responseTestContext{
								OpID:      "updateTask",
								Name:      "success",
								IsSuccess: code/100 == 2,
							}
							res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, code, response)
							if assert.NoError(t, err) {
								if !sucCtx.Assert(t, response, res) {
									return
								}
							}
						}
					}
				}
			}
		}
	}

}
func TestInlinedSchemaResponses(t *testing.T) {
	b, err := opBuilder("getTasks", "../fixtures/codegen/todolist.responses.yml")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, _, op, ok := b.Analyzed.OperationForName("getTasks")
	if assert.True(t, ok) && assert.NotNil(t, op) && assert.NotNil(t, op.Responses) {
		resolver := &typeResolver{ModelsPackage: b.ModelsPackage, Doc: b.Doc}
		if assert.NotNil(t, op.Responses.Default) {
			resp := *op.Responses.Default
			defCtx := responseTestContext{
				OpID: "getTasks",
				Name: "default",
			}
			res, err := b.MakeResponse("a", defCtx.Name, false, resolver, -1, resp)
			if assert.NoError(t, err) {
				if defCtx.Assert(t, resp, res) {
					for code, response := range op.Responses.StatusCodeResponses {
						sucCtx := responseTestContext{
							OpID:      "getTasks",
							Name:      "success",
							IsSuccess: code/100 == 2,
						}
						res, err := b.MakeResponse("a", sucCtx.Name, sucCtx.IsSuccess, resolver, code, response)
						if assert.NoError(t, err) {
							if !sucCtx.Assert(t, response, res) {
								return
							}
						}
						assert.Len(t, b.ExtraSchemas, 1)
						// ExtraSchema is not a definition: it is rendered in current operations package
						assert.Equal(t, "[]*SuccessBodyItems0", res.Schema.GoType)
					}
				}
			}
		}
	}

}

type responseTestContext struct {
	OpID      string
	Name      string
	IsSuccess bool
}

func (ctx *responseTestContext) Assert(t testing.TB, response spec.Response, res GenResponse) bool {
	if !assert.Equal(t, ctx.IsSuccess, res.IsSuccess) {
		return false
	}

	if !assert.Equal(t, ctx.Name, res.Name) {
		return false
	}

	if !assert.Equal(t, response.Description, res.Description) {
		return false
	}

	if len(response.Headers) > 0 {
		for k, v := range response.Headers {
			found := false
			for _, h := range res.Headers {
				if h.Name == k {
					found = true
					if k == "X-Last-Task-Id" {
						hctx := &respHeaderTestContext{k, "swag.FormatInt64", "swag.ConvertInt64"}
						if !hctx.Assert(t, v, h) {
							return false
						}
						break
					}
					if k == "X-Error-Code" {
						hctx := &respHeaderTestContext{k, "", ""}
						if !hctx.Assert(t, v, h) {
							return false
						}
					}

					break
				}
			}
			if !assert.True(t, found) {
				return false
			}
		}
	}

	if response.Schema != nil {
		if !assert.NotNil(t, res.Schema) {
			return false
		}
	} else {
		if !assert.Nil(t, res.Schema) {
			return false
		}
	}

	return true
}

type respHeaderTestContext struct {
	Name      string
	Formatter string
	Converter string
}

func (ctx *respHeaderTestContext) Assert(t testing.TB, header spec.Header, hdr GenHeader) bool {
	if !assert.Equal(t, ctx.Name, hdr.Name) {
		return false
	}
	if !assert.Equal(t, "a", hdr.ReceiverName) {
		return false
	}
	if !assert.Equal(t, ctx.Formatter, hdr.Formatter) {
		return false
	}
	if !assert.Equal(t, ctx.Converter, hdr.Converter) {
		return false
	}
	if !assert.Equal(t, header.Description, hdr.Description) {
		return false
	}
	if !assert.Equal(t, header.Minimum, hdr.Minimum) || !assert.Equal(t, header.ExclusiveMinimum, hdr.ExclusiveMinimum) {
		return false
	}
	if !assert.Equal(t, header.Maximum, hdr.Maximum) || !assert.Equal(t, header.ExclusiveMaximum, hdr.ExclusiveMaximum) {
		return false
	}
	if !assert.Equal(t, header.MinLength, hdr.MinLength) {
		return false
	}
	if !assert.Equal(t, header.MaxLength, hdr.MaxLength) {
		return false
	}
	if !assert.Equal(t, header.Pattern, hdr.Pattern) {
		return false
	}
	if !assert.Equal(t, header.MaxItems, hdr.MaxItems) {
		return false
	}
	if !assert.Equal(t, header.MinItems, hdr.MinItems) {
		return false
	}
	if !assert.Equal(t, header.UniqueItems, hdr.UniqueItems) {
		return false
	}
	if !assert.Equal(t, header.MultipleOf, hdr.MultipleOf) {
		return false
	}
	if !assert.EqualValues(t, header.Enum, hdr.Enum) {
		return false
	}
	if !assert.Equal(t, header.Type, hdr.SwaggerType) {
		return false
	}
	if !assert.Equal(t, header.Format, hdr.SwaggerFormat) {
		return false
	}
	return true
}

func TestGenResponses_Issue540(t *testing.T) {
	b, err := opBuilder("postPet", "../fixtures/bugs/540/swagger.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("post_pet_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "func (o *PostPetOK) WithPayload(payload models.Pet) *PostPetOK {", string(ff))
					assertInCode(t, "func (o *PostPetOK) SetPayload(payload models.Pet) {", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponses_Issue718_NotRequired(t *testing.T) {
	b, err := opBuilder("doEmpty", "../fixtures/codegen/todolist.simple.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "if payload == nil", string(ff))
					assertInCode(t, "payload = make([]models.Foo, 0, 50)", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponses_Issue718_Required(t *testing.T) {
	b, err := opBuilder("doEmpty", "../fixtures/codegen/todolist.simple.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "if payload == nil", string(ff))
					assertInCode(t, "payload = make([]models.Foo, 0, 50)", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

// Issue776 includes references that span multiple files. Flattening or Expanding is required
func TestGenResponses_Issue776_Spec(t *testing.T) {
	//spec.Debug = true
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
		//spec.Debug = false
	}()

	b, err := opBuilderWithFlatten("GetItem", "../fixtures/bugs/776/spec.yaml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					// This should be models.Item if flat works correctly
					assertInCode(t, "Payload *models.Item", string(ff))
					assertNotInCode(t, "type GetItemOKBody struct", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponses_Issue776_SwaggerTemplate(t *testing.T) {
	//spec.Debug = true
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
		//spec.Debug = false
	}()

	b, err := opBuilderWithFlatten("getHealthy", "../fixtures/bugs/776/swagger-template.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "Payload *models.Error", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestIssue846(t *testing.T) {
	// do it 8 times, to ensure it's always in the same order
	for i := 0; i < 8; i++ {
		b, err := opBuilder("getFoo", "../fixtures/bugs/846/swagger.yml")
		if assert.NoError(t, err) {
			op, err := b.MakeOperation()
			if assert.NoError(t, err) {
				var buf bytes.Buffer
				opts := opts()
				if assert.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op)) {
					ff, err := opts.LanguageOpts.FormatContent("do_empty_responses.go", buf.Bytes())
					if assert.NoError(t, err) {
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
					} else {
						fmt.Println(buf.String())
					}
				}
			}
		}
	}
}

func TestIssue881(t *testing.T) {
	b, err := opBuilder("getFoo", "../fixtures/bugs/881/swagger.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
		}
	}
}

func TestIssue881Deep(t *testing.T) {
	b, err := opBuilder("getFoo", "../fixtures/bugs/881/deep.yml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op))
		}
	}
}

func TestGenResponses_XGoName(t *testing.T) {
	b, err := opBuilder("putTesting", "../fixtures/specs/response_name.json")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("put_testing_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "const PutTestingAlternateNameCode int =", string(ff))
					assertInCode(t, "type PutTestingAlternateName struct {", string(ff))
					assertInCode(t, "func NewPutTestingAlternateName() *PutTestingAlternateName {", string(ff))
					assertInCode(t, "func (o *PutTestingAlternateName) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponses_Issue892(t *testing.T) {
	b, err := methodPathOpBuilder("get", "/media/search", "../fixtures/bugs/982/swagger.yaml")
	if assert.NoError(t, err) {
		op, err := b.MakeOperation()
		if assert.NoError(t, err) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("clientResponse").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("get_media_search_responses.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "o.Media = aO0", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponses_Issue1013(t *testing.T) {
	b, erp := methodPathOpBuilder("get", "/test", "../fixtures/bugs/1013/fixture-1013.yaml")
	if assert.NoError(t, erp) {
		op, erm := b.MakeOperation()
		if assert.NoError(t, erm) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "Payload *models.Response `json:\"body,omitempty\"`", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
	b, erp = methodPathOpBuilder("get", "/test2", "../fixtures/bugs/1013/fixture-1013.yaml")
	if assert.NoError(t, erp) {
		op, erm := b.MakeOperation()
		if assert.NoError(t, erm) {
			var buf bytes.Buffer
			opts := opts()
			if assert.NoError(t, templates.MustGet("serverResponses").Execute(&buf, op)) {
				ff, err := opts.LanguageOpts.FormatContent("foo.go", buf.Bytes())
				if assert.NoError(t, err) {
					assertInCode(t, "Payload *models.Response `json:\"body,omitempty\"`", string(ff))
				} else {
					fmt.Println(buf.String())
				}
			}
		}
	}
}

func TestGenResponse_15362_WithExpand(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

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
	log.SetOutput(ioutil.Discard)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	// testing fixture-1572.yaml with minimal flatten
	// edge cases for operations schemas

	/*
	        Run the following test caes and exercise the minimal flatten mode:
	   - [x] nil schema in body param / response
	   - [x] interface{} in body param /response
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
				"	Payload interface{} `json:\"body,omitempty\"`",
				`func NewGetInterfaceOK() *GetInterfaceOK {`,
				`	return &GetInterfaceOK{`,
				`func (o *GetInterfaceOK) WithPayload(payload interface{}) *GetInterfaceOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetInterfaceOK) SetPayload(payload interface{}) {`,
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
				"	Payload interface{} `json:\"body,omitempty\"`",
				`func NewGetNullOK() *GetNullOK {`,
				`	return &GetNullOK{`,
				`func (o *GetNullOK) WithPayload(payload interface{}) *GetNullOK {`,
				`	o.Payload = payload`,
				`	return o`,
				`func (o *GetNullOK) SetPayload(payload interface{}) {`,
				`	o.Payload = payload`,
				`func (o *GetNullOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {`,
				`	rw.WriteHeader(200`,
				`	payload := o.Payload`,
				`	if err := producer.Produce(rw, payload); err != nil {`,
				`const GetNullRequestProcessedCode int = 203`,
				`/*GetNullRequestProcessed OK`,
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
