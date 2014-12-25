package swagger

import (
	"strings"
	"testing"

	"github.com/casualjim/go-swagger"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServeAPI(t *testing.T) {
	Convey("ServeAPI should", t, func() {

		spec := &swagger.Spec{
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Paths: swagger.Paths{
				Paths: map[string]swagger.PathItem{
					"/": swagger.PathItem{
						Get: &swagger.Operation{
							Consumes: []string{"application/x-yaml"},
							Produces: []string{"application/x-yaml"},
							ID:       "someOperation",
						},
					},
				},
			},
		}

		api := NewAPI(spec)

		Convey("return an error when the passed api is nil", func() {
			h, err := ServeAPI(nil, nil)
			So(h, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("return an error when the API registrations are invalid", func() {
			h, err := ServeAPI(api, nil)
			So(h, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(strings.HasPrefix(err.Error(), "missing"), ShouldBeTrue)
		})
	})
}
