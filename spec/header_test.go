package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func float64Ptr(f float64) *float64 {
	return &f
}
func int64Ptr(f int64) *int64 {
	return &f
}

var header = Header{
	headerProps: headerProps{Description: "the description of this header"},
	simpleSchema: simpleSchema{
		Items: &Items{
			refable: refable{Ref: MustCreateRef("Cat")},
		},
		Type:    "string",
		Format:  "date",
		Default: "8",
	},
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
			So(actual, ShouldBeEquivalentTo, header)
		})
	})
}
