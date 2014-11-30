package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var parameter = Parameter{
	Items: &Items{
		Ref: "Cat",
	},
	Extensions: map[string]interface{}{
		"x-framework": "swagger-go",
	},
	Ref:              "Dog",
	Description:      "the description of this parameter",
	Maximum:          100,
	ExclusiveMaximum: true,
	ExclusiveMinimum: true,
	Minimum:          5,
	MaxLength:        100,
	MinLength:        5,
	Pattern:          "\\w{1,5}\\w+",
	MaxItems:         100,
	MinItems:         5,
	UniqueItems:      true,
	MultipleOf:       5,
	Enum:             []interface{}{"hello", "world"},
	Type:             "string",
	Format:           "date",
	Name:             "param-name",
	In:               "header",
	Required:         true,
	Schema:           &Schema{Type: &StringOrArray{Single: "string"}},
	CollectionFormat: "csv",
	Default:          "8",
}

var parameterJson = `{
	"items": { 
		"$ref": "Cat"
	},
	"x-framework": "swagger-go",
  "$ref": "Dog",
  "description": "the description of this parameter",
  "maximum": 100,
  "minimum": 5,
  "exclusiveMaximum": true,
  "exclusiveMinimum": true,
  "maxLength": 100,
  "minLength": 5,
  "pattern": "\\w{1,5}\\w+",
  "maxItems": 100,
  "minItems": 5,
  "uniqueItems": true,
  "multipleOf": 5,
  "enum": ["hello", "world"],
  "type": "string",
  "format": "date",
	"name": "param-name",
	"in": "header",
	"required": true,
	"schema": {
		"type": "string"
	},
	"collectionFormat": "csv",
	"default": "8"
}`

func TestIntegrationParameter(t *testing.T) {
	Convey("for all properties a parameter should", t, func() {
		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(parameterJson), &expected)
			b, err := json.Marshal(parameter)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("deserialize", func() {
			actual := Parameter{}
			err := json.Unmarshal([]byte(parameterJson), &actual)
			So(err, ShouldBeNil)
			So(actual.Items, ShouldResemble, parameter.Items)
			So(actual.Extensions, ShouldResemble, parameter.Extensions)
			So(actual.Ref, ShouldEqual, parameter.Ref)
			So(actual.Description, ShouldEqual, parameter.Description)
			So(actual.Maximum, ShouldEqual, parameter.Maximum)
			So(actual.Minimum, ShouldEqual, parameter.Minimum)
			So(actual.ExclusiveMinimum, ShouldEqual, parameter.ExclusiveMinimum)
			So(actual.ExclusiveMaximum, ShouldEqual, parameter.ExclusiveMaximum)
			So(actual.MaxLength, ShouldEqual, parameter.MaxLength)
			So(actual.MinLength, ShouldEqual, parameter.MinLength)
			So(actual.Pattern, ShouldEqual, parameter.Pattern)
			So(actual.MaxItems, ShouldEqual, parameter.MaxItems)
			So(actual.MinItems, ShouldEqual, parameter.MinItems)
			So(actual.UniqueItems, ShouldBeTrue)
			So(actual.MultipleOf, ShouldEqual, parameter.MultipleOf)
			So(actual.Enum, ShouldResemble, parameter.Enum)
			So(actual.Type, ShouldResemble, parameter.Type)
			So(actual.Format, ShouldEqual, parameter.Format)
			So(actual.Name, ShouldEqual, parameter.Name)
			So(actual.In, ShouldEqual, parameter.In)
			So(actual.Required, ShouldEqual, parameter.Required)
			So(actual.Schema, ShouldResemble, parameter.Schema)
			So(actual.CollectionFormat, ShouldEqual, parameter.CollectionFormat)
			So(actual.Default, ShouldResemble, parameter.Default)
		})
	})
}

func TestParameterSerialization(t *testing.T) {

	Convey("Parameters should serialize", t, func() {

		Convey("a query parameter", func() {
			param := QueryParam()
			param.Type = "string"
			So(param, ShouldSerializeJSON, `{"in":"query","type":"string"}`)
		})

		Convey("a query parameter with array", func() {
			param := QueryParam()
			param.Type = "array"
			param.CollectionFormat = "multi"
			param.Items = &Items{
				Type: "string",
			}
			So(param, ShouldSerializeJSON, `{"collectionFormat":"multi","in":"query","items":{"type":"string"},"type":"array"}`)
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
			So(param, ShouldSerializeJSON, `{"in":"body","schema":{"properties":{"name":{"type":"string"}}}}`)
		})

		Convey("a ref body parameter", func() {
			param := BodyParam()
			param.Schema = &Schema{
				Ref: "Cat",
			}
			So(param, ShouldSerializeJSON, `{"in":"body","schema":{"$ref":"Cat"}}`)
		})

		Convey("serialize an array body parameter", func() {
			param := BodyParam()
			param.Schema = ArrayProperty(RefProperty("Cat"))
			So(param, ShouldSerializeJSON, `{"in":"body","schema":{"items":{"$ref":"Cat"},"type":"array"}}`)
		})
	})
}
