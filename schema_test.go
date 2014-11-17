package swagger

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPropertySerialization(t *testing.T) {

	Convey("Properties should serialize", t, func() {
		Convey("a boolean property", func() {
			prop := BooleanProperty()
			So(prop, validateJSON, `{"type":"boolean"}`)
		})
		Convey("a date property", func() {
			prop := DateProperty()
			So(prop, validateJSON, `{"format":"date","type":"string"}`)
		})
		Convey("a date-time property", func() {
			prop := DateTimeProperty()
			So(prop, validateJSON, `{"format":"date-time","type":"string"}`)
		})
		Convey("a float64 property", func() {
			prop := Float64Property()
			So(prop, validateJSON, `{"format":"double","type":"number"}`)
		})
		Convey("a float32 property", func() {
			prop := Float32Property()
			So(prop, validateJSON, `{"format":"float","type":"number"}`)
		})
		Convey("a int32 property", func() {
			prop := Int32Property()
			So(prop, validateJSON, `{"format":"int32","type":"number"}`)
		})
		Convey("a int64 property", func() {
			prop := Int64Property()
			So(prop, validateJSON, `{"format":"int64","type":"number"}`)
		})
		Convey("a string map property", func() {
			prop := MapProperty(StringProperty())
			So(prop, validateJSON, `{"additionalProperties":{"type":"string"},"type":"object"}`)
		})
		Convey("an int32 map property", func() {
			prop := MapProperty(Int32Property())
			So(prop, validateJSON, `{"additionalProperties":{"format":"int32","type":"number"},"type":"object"}`)
		})
		Convey("a ref property", func() {
			prop := RefProperty("Dog")
			So(prop, validateJSON, `{"$ref":"Dog"}`)
		})
		Convey("a string property", func() {
			prop := StringProperty()
			So(prop, validateJSON, `{"type":"string"}`)
		})
		Convey("a string property with enums", func() {
			prop := StringProperty()
			prop.Enum = append(prop.Enum, "a", "b")
			So(prop, validateJSON, `{"enum":["a","b"],"type":"string"}`)
		})
		Convey("a string array property", func() {
			prop := ArrayProperty(StringProperty())

			So(prop, validateJSON, `{"items":{"type":"string"},"type":"array"}`)
		})
	})
}
