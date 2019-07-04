// +build !go1.11

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

package scan

import (
	goparser "go/parser"
	"log"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestOperationsExpression(t *testing.T) {
	assert.Regexp(t, rxOperation, "swagger:operation DELETE /orders/{id} deleteOrder")
	assert.Regexp(t, rxOperation, "swagger:operation GET /v1.2/something deleteOrder")
}

func TestOperationsParser(t *testing.T) {
	docFile := "../fixtures/goparsing/classification/operations_annotation/operations.go"
	fileTree, err := goparser.ParseFile(classificationProg.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	op := newOperationsParser(classificationProg)
	var ops spec.Paths
	err = op.Parse(fileTree, &ops, nil, nil)
	assert.NoError(t, err)

	assert.Len(t, ops.Paths, 3)

	po, ok := ops.Paths["/pets"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertAnnotationOperation(t,
		po.Get,
		"getPet",
		"",
		"List all pets",
		[]string{"pets"},
	)
	if po.Get != nil {
		rsp, k := po.Get.Responses.StatusCodeResponses[200]
		if assert.True(t, k) {
			assert.Equal(t, "An paged array of pets", rsp.Description)
		}
		if assert.NotNil(t, po.Get.Responses.Default) {
			assert.Equal(t, "unexpected error", po.Get.Responses.Default.Description)
		}
	}

	po, ok = ops.Paths["/pets/{id}"]
	assert.True(t, ok)
	assert.NotNil(t, po.Put)
	assertAnnotationOperation(t,
		po.Put,
		"updatePet",
		"Updates the details for a pet.",
		"Some long explanation,\nspanning over multipele lines,\nAKA the description.",
		[]string{"pets"},
	)
	if po.Put != nil {
		rsp, k := po.Put.Responses.StatusCodeResponses[400]
		if assert.True(t, k) {
			assert.Equal(t, "Invalid ID supplied", rsp.Description)
		}
		rsp, k = po.Put.Responses.StatusCodeResponses[404]
		if assert.True(t, k) {
			assert.Equal(t, "Pet not found", rsp.Description)
		}
		rsp, k = po.Put.Responses.StatusCodeResponses[405]
		if assert.True(t, k) {
			assert.Equal(t, "Validation exception", rsp.Description)
		}
	}

	po, ok = ops.Paths["/v1/events"]
	assert.True(t, ok)
	assert.NotNil(t, po.Get)
	assertAnnotationOperation(t,
		po.Get,
		"getEvents",
		"Events",
		"Mitigation Events",
		[]string{"Events"},
	)
	if po.Get != nil {
		rsp, k := po.Get.Responses.StatusCodeResponses[200]
		if assert.True(t, k) {
			assert.Equal(t, "#/definitions/ListResponse", rsp.Schema.Ref.String())
			assert.Equal(t, "200", rsp.Description)
		}
		rsp, k = po.Get.Responses.StatusCodeResponses[400]
		if assert.True(t, k) {
			assert.Equal(t, "#/definitions/ErrorResponse", rsp.Schema.Ref.String())
			assert.Equal(t, "400", rsp.Description)
		}
	}
}

func assertAnnotationOperation(t *testing.T, op *spec.Operation, id, summary, description string, tags []string) {
	assert.NotNil(t, op)
	assert.Equal(t, summary, op.Summary)
	assert.Equal(t, description, op.Description)
	assert.Equal(t, id, op.ID)
	assert.EqualValues(t, tags, op.Tags)
	assert.Contains(t, op.Consumes, "application/json")
	assert.Contains(t, op.Consumes, "application/xml")
	assert.Contains(t, op.Produces, "application/json")
	assert.Contains(t, op.Produces, "application/xml")
	assert.Len(t, op.Security, 1)
	if len(op.Security) > 0 {
		akv, ok := op.Security[0]["petstore_auth"]
		assert.True(t, ok)
		// akv must be defined & not empty
		assert.NotNil(t, akv)
		assert.NotEmpty(t, akv)
	}
}
