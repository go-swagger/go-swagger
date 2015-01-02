package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestXmlObject(t *testing.T) {

	Convey("an xml object should", t, func() {
		Convey("serialize", func() {
			Convey("an empty object", func() {
				obj := XMLObject{}
				expected := "{}"
				actual, err := json.Marshal(obj)
				So(err, ShouldBeNil)
				So(string(actual), ShouldEqual, expected)
			})
			Convey("a completed object", func() {
				obj := XMLObject{
					Name:      "the name",
					Namespace: "the namespace",
					Prefix:    "the prefix",
					Attribute: true,
					Wrapped:   true,
				}
				actual, err := json.Marshal(obj)
				So(err, ShouldBeNil)
				var ad map[string]interface{}
				err = json.Unmarshal(actual, &ad)
				So(err, ShouldBeNil)
				So(ad["name"], ShouldEqual, obj.Name)
				So(ad["namespace"], ShouldEqual, obj.Namespace)
				So(ad["prefix"], ShouldEqual, obj.Prefix)
				So(ad["attribute"], ShouldBeTrue)
				So(ad["wrapped"], ShouldBeTrue)
			})
		})
		Convey("deserialize", func() {
			Convey("an empty object", func() {
				expected := XMLObject{}
				actual := XMLObject{}
				err := json.Unmarshal([]byte("{}"), &actual)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
			Convey("a completed object", func() {
				completed := `{"name":"the name","namespace":"the namespace","prefix":"the prefix","attribute":true,"wrapped":true}`
				expected := XMLObject{"the name", "the namespace", "the prefix", true, true}
				actual := XMLObject{}
				err := json.Unmarshal([]byte(completed), &actual)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})
		})
	})

}
