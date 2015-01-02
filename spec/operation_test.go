package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var operation = Operation{
	vendorExtensible: vendorExtensible{
		Extensions: map[string]interface{}{
			"x-framework": "go-swagger",
		},
	},
	operationProps: operationProps{
		Description: "operation description",
		Consumes:    []string{"application/json", "application/x-yaml"},
		Produces:    []string{"application/json", "application/x-yaml"},
		Schemes:     []string{"http", "https"},
		Tags:        []string{"dogs"},
		Summary:     "the summary of the operation",
		ID:          "sendCat",
		Deprecated:  true,
		Security: []map[string][]string{
			map[string][]string{
				"apiKey": []string{},
			},
		},
		Parameters: []Parameter{
			Parameter{refable: refable{Ref: "Cat"}},
		},
		Responses: &Responses{
			responsesProps: responsesProps{
				Default: &Response{
					responseProps: responseProps{
						Description: "void response",
					},
				},
			},
		},
	},
}

var operationJSON = `{
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

func TestIntegrationOperation(t *testing.T) {

	Convey("all fields of an operation should", t, func() {

		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(operationJSON), &expected)
			b, err := json.Marshal(operation)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("deserialize", func() {
			actual := Operation{}
			err := json.Unmarshal([]byte(operationJSON), &actual)
			So(err, ShouldBeNil)
			So(actual.Description, ShouldEqual, operation.Description)
			So(actual.Extensions, ShouldResemble, operation.Extensions)
			So(actual.Consumes, ShouldResemble, operation.Consumes)
			So(actual.Produces, ShouldResemble, operation.Produces)
			So(actual.Tags, ShouldResemble, operation.Tags)
			So(actual.Schemes, ShouldResemble, operation.Schemes)
			So(actual.Summary, ShouldEqual, operation.Summary)
			So(actual.ID, ShouldEqual, operation.ID)
			So(actual.Deprecated, ShouldEqual, operation.Deprecated)
			So(actual.Security, ShouldResemble, operation.Security)
			So(actual.Parameters, ShouldResemble, operation.Parameters)
			So(actual.Responses.Default, ShouldResemble, operation.Responses.Default)
		})

	})

}
