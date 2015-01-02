package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var pathItem = PathItem{
	refable: refable{Ref: "Dog"},
	vendorExtensible: vendorExtensible{
		Extensions: map[string]interface{}{
			"x-framework": "go-swagger",
		},
	},
	pathItemProps: pathItemProps{
		Get: &Operation{
			operationProps: operationProps{Description: "get operation description"},
		},
		Put: &Operation{
			operationProps: operationProps{Description: "put operation description"},
		},
		Post: &Operation{
			operationProps: operationProps{Description: "post operation description"},
		},
		Delete: &Operation{
			operationProps: operationProps{Description: "delete operation description"},
		},
		Options: &Operation{
			operationProps: operationProps{Description: "options operation description"},
		},
		Head: &Operation{
			operationProps: operationProps{Description: "head operation description"},
		},
		Patch: &Operation{
			operationProps: operationProps{Description: "patch operation description"},
		},
		Parameters: []Parameter{
			Parameter{
				paramProps: paramProps{In: "path"},
			},
		},
	},
}

var pathItemJSON = `{
	"$ref": "Dog",
	"x-framework": "go-swagger",
	"get": { "description": "get operation description" },
	"put": { "description": "put operation description" },
	"post": { "description": "post operation description" },
	"delete": { "description": "delete operation description" },
	"options": { "description": "options operation description" },
	"head": { "description": "head operation description" },
	"patch": { "description": "patch operation description" },
	"parameters": [{"in":"path"}]
}`

func TestIntegrationPathItem(t *testing.T) {
	Convey("all fields of a path item should", t, func() {

		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(pathItemJSON), &expected)
			b, err := json.Marshal(pathItem)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual["$ref"], ShouldResemble, expected["$ref"])
			So(actual["x-framework"], ShouldResemble, expected["x-framework"])
			So(actual["get"], ShouldResemble, expected["get"])
			So(actual["put"], ShouldResemble, expected["put"])
			So(actual["post"], ShouldResemble, expected["post"])
			So(actual["delete"], ShouldResemble, expected["delete"])
			So(actual["options"], ShouldResemble, expected["options"])
			So(actual["head"], ShouldResemble, expected["head"])
			So(actual["patch"], ShouldResemble, expected["patch"])
			So(actual["parameters"], ShouldResemble, expected["parameters"])
		})

		Convey("deserialize", func() {
			actual := PathItem{}
			err := json.Unmarshal([]byte(pathItemJSON), &actual)
			So(err, ShouldBeNil)
			So(actual.Ref, ShouldEqual, pathItem.Ref)
			So(actual.Get, ShouldResemble, pathItem.Get)
			So(actual.Put, ShouldResemble, pathItem.Put)
			So(actual.Post, ShouldResemble, pathItem.Post)
			So(actual.Delete, ShouldResemble, pathItem.Delete)
			So(actual.Options, ShouldResemble, pathItem.Options)
			So(actual.Head, ShouldResemble, pathItem.Head)
			So(actual.Patch, ShouldResemble, pathItem.Patch)
			So(actual.Parameters, ShouldResemble, pathItem.Parameters)
		})
	})
}
