package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v1"
)

func TestIntegrationLicense(t *testing.T) {
	Convey("all fields of licens should", t, func() {
		Convey("serialize to JSON", func() {
			b, err := json.Marshal(License{"the name", "the url"})
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, `{"name":"the name","url":"the url"}`)
		})

		Convey("serialize to YAML", func() {
			b, err := yaml.Marshal(License{"the name", "the url"})
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, "name: the name\nurl: the url\n")
		})

		Convey("deserialize from JSON", func() {
			actual := License{}
			err := json.Unmarshal([]byte(`{"name":"the name","url":"the url"}`), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, License{"the name", "the url"})
		})

		Convey("deserialize from JSON", func() {
			actual := License{}
			err := yaml.Unmarshal([]byte("name: the name\nurl: the url\n"), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, License{"the name", "the url"})
		})
	})
}
