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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var operation = Operation{
	VendorExtensible: VendorExtensible{
		Extensions: map[string]interface{}{
			"x-framework": "go-swagger",
		},
	},
	OperationProps: OperationProps{
		Description: "operation description",
		Consumes:    []string{"application/json", "application/x-yaml"},
		Produces:    []string{"application/json", "application/x-yaml"},
		Schemes:     []string{"http", "https"},
		Tags:        []string{"dogs"},
		Summary:     "the summary of the operation",
		ID:          "sendCat",
		Deprecated:  true,
		Security: []map[string][]string{
			{
				"apiKey": {},
			},
		},
		Parameters: []Parameter{
			{Refable: Refable{Ref: MustCreateRef("Cat")}},
		},
		Responses: &Responses{
			ResponsesProps: ResponsesProps{
				Default: &Response{
					ResponseProps: ResponseProps{
						Description: "void response",
					},
				},
			},
		},
	},
}

const operationJSON = `{
	"description": "operation description",
	"x-framework": "go-swagger",
	"consumes": [ "application/json", "application/x-yaml" ],
	"produces": [ "application/json", "application/x-yaml" ],
	"schemes": ["http", "https"],
	"tags": ["dogs"],
	"summary": "the summary of the operation",
	"operationId": "sendCat",
	"deprecated": true,
	"security": [ { "apiKey": [] } ],
	"parameters": [{"$ref":"Cat"}],
	"responses": {
		"default": {
			"description": "void response"
		}
	}
}`

func TestSuccessResponse(t *testing.T) {
	ope := &Operation{}
	resp, n, f := ope.SuccessResponse()
	assert.Nil(t, resp)
	assert.Equal(t, 0, n)
	assert.Equal(t, false, f)

	resp, n, f = operation.SuccessResponse()
	if assert.NotNil(t, resp) {
		assert.Equal(t, "void response", resp.Description)
	}
	assert.Equal(t, 0, n)
	assert.Equal(t, false, f)

	err := json.Unmarshal([]byte(operationJSON), ope)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	ope = ope.RespondsWith(301, &Response{
		ResponseProps: ResponseProps{
			Description: "failure",
		},
	})
	resp, n, f = ope.SuccessResponse()
	if assert.NotNil(t, resp) {
		assert.Equal(t, "void response", resp.Description)
	}
	assert.Equal(t, 0, n)
	assert.Equal(t, false, f)

	ope = ope.RespondsWith(200, &Response{
		ResponseProps: ResponseProps{
			Description: "success",
		},
	})

	resp, n, f = ope.SuccessResponse()
	if assert.NotNil(t, resp) {
		assert.Equal(t, "success", resp.Description)
	}
	assert.Equal(t, 200, n)
	assert.Equal(t, true, f)
}

func TestOperationBuilder(t *testing.T) {
	ope := NewOperation("").WithID("operationID")
	ope = ope.RespondsWith(200, &Response{
		ResponseProps: ResponseProps{
			Description: "success",
		},
	}).
		WithDefaultResponse(&Response{
			ResponseProps: ResponseProps{
				Description: "default",
			},
		}).
		SecuredWith("scheme-name", "scope1", "scope2").
		WithConsumes("application/json").
		WithProduces("application/json").
		Deprecate().
		WithTags("this", "that").
		AddParam(nil).
		AddParam(QueryParam("myQueryParam").Typed("integer", "int32")).
		AddParam(QueryParam("myQueryParam").Typed("string", "hostname")).
		AddParam(PathParam("myPathParam").Typed("string", "uuid")).
		WithDescription("test operation").
		WithSummary("my summary").
		WithExternalDocs("some doc", "https://www.example.com")

	jazon, _ := json.MarshalIndent(ope, "", " ")
	assert.JSONEq(t, `{
		     "operationId": "operationID",
				 "description": "test operation",
				 "summary": "my summary",
				 "externalDocs": {
					 "description": "some doc",
					 "url": "https://www.example.com"
				 },
	       "security": [
          {
           "scheme-name": [
            "scope1",
            "scope2"
           ]
          }
         ],
         "consumes": [
          "application/json"
         ],
         "produces": [
          "application/json"
         ],
         "tags": [
          "this",
          "that"
         ],
         "deprecated": true,
         "parameters": [
          {
           "type": "string",
           "format": "hostname",
           "name": "myQueryParam",
           "in": "query"
          },
          {
           "type": "string",
           "format": "uuid",
           "name": "myPathParam",
           "in": "path",
           "required": true
          }
         ],
				 "responses": {
          "200": {
           "description": "success"
          },
          "default": {
           "description": "default"
          }
         }
		 }`, string(jazon))

	// check token lookup
	token, err := ope.JSONLookup("responses")
	assert.NoError(t, err)
	jazon, _ = json.MarshalIndent(token, "", " ")
	assert.JSONEq(t, `{
         "200": {
          "description": "success"
         },
         "default": {
          "description": "default"
         }
			 }`, string(jazon))

	// check delete methods
	ope = ope.RespondsWith(200, nil).
		RemoveParam("myQueryParam", "query").
		RemoveParam("myPathParam", "path").
		RemoveParam("fakeParam", "query").
		Undeprecate().
		WithExternalDocs("", "")
	jazon, _ = json.MarshalIndent(ope, "", " ")
	assert.JSONEq(t, `{
         "security": [
          {
           "scheme-name": [
            "scope1",
            "scope2"
           ]
          }
         ],
         "description": "test operation",
         "consumes": [
          "application/json"
         ],
         "produces": [
          "application/json"
         ],
         "tags": [
          "this",
          "that"
         ],
         "summary": "my summary",
         "operationId": "operationID",
         "responses": {
          "default": {
           "description": "default"
          }
         }
			 }`, string(jazon))
}

func TestIntegrationOperation(t *testing.T) {
	var actual Operation
	if assert.NoError(t, json.Unmarshal([]byte(operationJSON), &actual)) {
		assert.EqualValues(t, actual, operation)
	}

	assertParsesJSON(t, operationJSON, operation)
}

func TestSecurityProperty(t *testing.T) {
	//Ensure we omit security key when unset
	securityNotSet := OperationProps{}
	jsonResult, err := json.Marshal(securityNotSet)
	if assert.NoError(t, err) {
		assert.NotContains(t, string(jsonResult), "security", "security key should be omitted when unset")
	}

	//Ensure we preseve the security key when it contains an empty (zero length) slice
	securityContainsEmptyArray := OperationProps{
		Security: []map[string][]string{},
	}
	jsonResult, err = json.Marshal(securityContainsEmptyArray)
	if assert.NoError(t, err) {
		var props OperationProps
		if assert.NoError(t, json.Unmarshal(jsonResult, &props)) {
			assert.Equal(t, securityContainsEmptyArray, props)
		}
	}
}

func TestOperationGobEncoding(t *testing.T) {
	// 1. empty scope in security requirements:  "security": [ { "apiKey": [] } ],
	doTestOperationGobEncoding(t, operationJSON)

	// 2. nil security requirements
	doTestOperationGobEncoding(t, `{
	"description": "operation description",
	"x-framework": "go-swagger",
	"consumes": [ "application/json", "application/x-yaml" ],
	"produces": [ "application/json", "application/x-yaml" ],
	"schemes": ["http", "https"],
	"tags": ["dogs"],
	"summary": "the summary of the operation",
	"operationId": "sendCat",
	"deprecated": true,
	"parameters": [{"$ref":"Cat"}],
	"responses": {
		"default": {
			"description": "void response"
		}
	}
}`)

	// 3. empty security requirement
	doTestOperationGobEncoding(t, `{
	"description": "operation description",
	"x-framework": "go-swagger",
	"consumes": [ "application/json", "application/x-yaml" ],
	"produces": [ "application/json", "application/x-yaml" ],
	"schemes": ["http", "https"],
	"tags": ["dogs"],
	"security": [],
	"summary": "the summary of the operation",
	"operationId": "sendCat",
	"deprecated": true,
	"parameters": [{"$ref":"Cat"}],
	"responses": {
		"default": {
			"description": "void response"
		}
	}
}`)

	// 4. non-empty security requirements
	doTestOperationGobEncoding(t, `{
	"description": "operation description",
	"x-framework": "go-swagger",
	"consumes": [ "application/json", "application/x-yaml" ],
	"produces": [ "application/json", "application/x-yaml" ],
	"schemes": ["http", "https"],
	"tags": ["dogs"],
	"summary": "the summary of the operation",
	"security": [ { "scoped-auth": [ "phone", "email" ] , "api-key": []} ],
	"operationId": "sendCat",
	"deprecated": true,
	"parameters": [{"$ref":"Cat"}],
	"responses": {
		"default": {
			"description": "void response"
		}
	}
}`)

}

func doTestOperationGobEncoding(t *testing.T, fixture string) {
	var src, dst Operation

	if !assert.NoError(t, json.Unmarshal([]byte(fixture), &src)) {
		t.FailNow()
	}

	doTestAnyGobEncoding(t, &src, &dst)
}

func doTestAnyGobEncoding(t *testing.T, src, dst interface{}) {
	expectedJSON, _ := json.MarshalIndent(src, "", " ")

	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(src)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	err = gob.NewDecoder(&b).Decode(dst)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	jazon, err := json.MarshalIndent(dst, "", " ")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.JSONEq(t, string(expectedJSON), string(jazon))
}
