package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var paths = Paths{
	vendorExtensible: vendorExtensible{Extensions: map[string]interface{}{"x-framework": "go-swagger"}},
	Paths: map[string]PathItem{
		"/": PathItem{
			refable: refable{Ref: MustCreateRef("cats")},
		},
	},
}

var pathsJSON = `{"x-framework":"go-swagger","/":{"$ref":"cats"}}`

func TestIntegrationPaths(t *testing.T) {
	Convey("all fields of paths should", t, func() {

		Convey("serialize", func() {
			expected := map[string]interface{}{}
			json.Unmarshal([]byte(pathsJSON), &expected)
			b, err := json.Marshal(paths)
			So(err, ShouldBeNil)
			var actual map[string]interface{}
			err = json.Unmarshal(b, &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, expected)
		})

		Convey("deserialize", func() {

			actual := Paths{}
			err := json.Unmarshal([]byte(pathsJSON), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, paths)
		})

	})
}
