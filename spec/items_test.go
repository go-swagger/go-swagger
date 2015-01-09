package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var items = Items{
	refable: refable{Ref: MustCreateRef("Dog")},
	commonValidations: commonValidations{
		Maximum:          float64Ptr(100),
		ExclusiveMaximum: true,
		ExclusiveMinimum: true,
		Minimum:          float64Ptr(5),
		MaxLength:        int64Ptr(100),
		MinLength:        int64Ptr(5),
		Pattern:          "\\w{1,5}\\w+",
		MaxItems:         int64Ptr(100),
		MinItems:         int64Ptr(5),
		UniqueItems:      true,
		MultipleOf:       float64Ptr(5),
		Enum:             []interface{}{"hello", "world"},
	},
	simpleSchema: simpleSchema{
		Type:   "string",
		Format: "date",
		Items: &Items{
			refable: refable{Ref: MustCreateRef("Cat")},
		},
		CollectionFormat: "csv",
		Default:          "8",
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
			So(actual, ShouldBeEquivalentTo, items)
		})
	})
}
