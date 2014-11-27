package swagger

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParameterSerialization(t *testing.T) {

	Convey("Parameters should serialize", t, func() {

		Convey("a query parameter", func() {
			param := QueryParam()
			param.Type = "string"
			So(param, ShouldSerializeJSON, `{"in":"query","required":false,"type":"string"}`)
		})

		Convey("a query parameter with array", func() {
			param := QueryParam()
			param.Type = "array"
			param.CollectionFormat = "multi"
			param.Items = &Items{
				Type: "string",
			}
			So(param, ShouldSerializeJSON, `{"collectionFormat":"multi","in":"query","items":{"type":"string"},"required":false,"type":"array"}`)
		})

		Convey("a path parameter", func() {
			param := PathParam()
			param.Type = "string"
			So(param, ShouldSerializeJSON, `{"in":"path","required":true,"type":"string"}`)
		})

		Convey("a path parameter with string array", func() {
			param := PathParam()
			param.Type = "array"
			param.CollectionFormat = "multi"
			param.Items = &Items{
				Type: "string",
			}
			So(param, ShouldSerializeJSON, `{"collectionFormat":"multi","in":"path","items":{"type":"string"},"required":true,"type":"array"}`)
		})

		Convey("a path parameter with an int array", func() {
			param := PathParam()
			param.Type = "string"
			param.Type = "array"
			param.CollectionFormat = "multi"
			param.Items = &Items{
				Type:   "int",
				Format: "int32",
			}
			So(param, ShouldSerializeJSON, `{"collectionFormat":"multi","in":"path","items":{"format":"int32","type":"int"},"required":true,"type":"array"}`)
		})

		Convey("a header parameter", func() {
			param := HeaderParam()
			param.Type = "string"
			So(param, ShouldSerializeJSON, `{"in":"header","required":true,"type":"string"}`)
		})

		Convey("a header parameter with string array", func() {
			param := HeaderParam()
			param.Type = "array"
			param.CollectionFormat = "multi"
			param.Items = &Items{
				Type: "string",
			}
			So(param, ShouldSerializeJSON, `{"collectionFormat":"multi","in":"header","items":{"type":"string"},"required":true,"type":"array"}`)
		})

		Convey("a body parameter", func() {
			param := BodyParam()
			param.Schema = &Schema{
				Properties: map[string]Schema{
					"name": Schema{
						Type: &StringOrArray{Single: "string"},
					},
				},
			}
			So(param, ShouldSerializeJSON, `{"in":"body","required":false,"schema":{"properties":{"name":{"type":"string"}}}}`)
		})

		Convey("a ref body parameter", func() {
			param := BodyParam()
			param.Schema = &Schema{
				Ref: "Cat",
			}
			So(param, ShouldSerializeJSON, `{"in":"body","required":false,"schema":{"$ref":"Cat"}}`)
		})

		Convey("serialize an array body parameter", func() {
			param := BodyParam()
			param.Schema = ArrayProperty(RefProperty("Cat"))
			So(param, ShouldSerializeJSON, `{"in":"body","required":false,"schema":{"items":{"$ref":"Cat"},"type":"array"}}`)
		})
	})
}
