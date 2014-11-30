package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var header = Header{
	Description: "the description of this header",
	Items: &Items{
		Ref: "Cat",
	},
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
	Default:          "8",
}

var headerJSON = `{
	"items": { 
		"$ref": "Cat"
	},
  "description": "the description of this header",
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
	"default": "8"
}`

func TestIntegrationHeader(t *testing.T) {

	Convey("all fields of header should", t, func() {

		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(headerJSON), &expected)
			b, err := json.Marshal(header)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("deserialize", func() {
			actual := Header{}
			err := json.Unmarshal([]byte(headerJSON), &actual)
			So(err, ShouldBeNil)
			So(actual.Items, ShouldResemble, header.Items)
			So(actual.Description, ShouldEqual, header.Description)
			So(actual.Maximum, ShouldEqual, header.Maximum)
			So(actual.Minimum, ShouldEqual, header.Minimum)
			So(actual.ExclusiveMinimum, ShouldEqual, header.ExclusiveMinimum)
			So(actual.ExclusiveMaximum, ShouldEqual, header.ExclusiveMaximum)
			So(actual.MaxLength, ShouldEqual, header.MaxLength)
			So(actual.MinLength, ShouldEqual, header.MinLength)
			So(actual.Pattern, ShouldEqual, header.Pattern)
			So(actual.MaxItems, ShouldEqual, header.MaxItems)
			So(actual.MinItems, ShouldEqual, header.MinItems)
			So(actual.UniqueItems, ShouldBeTrue)
			So(actual.MultipleOf, ShouldEqual, header.MultipleOf)
			So(actual.Enum, ShouldResemble, header.Enum)
			So(actual.Type, ShouldResemble, header.Type)
			So(actual.Format, ShouldEqual, header.Format)
			So(actual.Default, ShouldResemble, header.Default)
		})
	})
}
