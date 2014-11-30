package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var items = Items{
	Ref:              "Dog",
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
	CollectionFormat: "csv",
	Default:          "8",
	Items: &Items{
		Ref: "Cat",
	},
}

var itemsJSON = `{
	"items": { 
		"$ref": "Cat"
	},
  "$ref": "Dog",
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
	"collectionFormat": "csv",
	"default": "8"
}`

func TestIntegrationItems(t *testing.T) {

	Convey("all fields of items should", t, func() {

		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(itemsJSON), &expected)
			b, err := json.Marshal(items)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("deserialize", func() {
			actual := Items{}
			err := json.Unmarshal([]byte(itemsJSON), &actual)
			So(err, ShouldBeNil)
			So(actual.Items, ShouldResemble, items.Items)
			So(actual.Ref, ShouldEqual, items.Ref)
			So(actual.Maximum, ShouldEqual, items.Maximum)
			So(actual.Minimum, ShouldEqual, items.Minimum)
			So(actual.ExclusiveMinimum, ShouldEqual, items.ExclusiveMinimum)
			So(actual.ExclusiveMaximum, ShouldEqual, items.ExclusiveMaximum)
			So(actual.MaxLength, ShouldEqual, items.MaxLength)
			So(actual.MinLength, ShouldEqual, items.MinLength)
			So(actual.Pattern, ShouldEqual, items.Pattern)
			So(actual.MaxItems, ShouldEqual, items.MaxItems)
			So(actual.MinItems, ShouldEqual, items.MinItems)
			So(actual.UniqueItems, ShouldBeTrue)
			So(actual.MultipleOf, ShouldEqual, items.MultipleOf)
			So(actual.Enum, ShouldResemble, items.Enum)
			So(actual.Type, ShouldResemble, items.Type)
			So(actual.Format, ShouldEqual, items.Format)
			So(actual.CollectionFormat, ShouldEqual, items.CollectionFormat)
			So(actual.Default, ShouldResemble, items.Default)
		})
	})
}
