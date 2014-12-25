package swagger

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/reflection"
	. "github.com/smartystreets/goconvey/convey"
)

func ShouldBeEquivalentTo(actual interface{}, expecteds ...interface{}) string {
	expected := expecteds[0]
	if actual == nil || expected == nil {
		return ""
	}

	if reflect.DeepEqual(expected, actual) {
		return ""
	}

	actualType := reflect.TypeOf(actual)
	if reflect.TypeOf(actual).ConvertibleTo(reflect.TypeOf(expected)) {
		expectedValue := reflect.ValueOf(expected)
		if reflection.IsZero(expectedValue) && reflection.IsZero(reflect.ValueOf(actual)) {
			return ""
		}

		// Attempt comparison after type conversion
		if reflect.DeepEqual(actual, expectedValue.Convert(actualType).Interface()) {
			return ""
		}
	}

	// Last ditch effort
	if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
		return ""
	}
	errFmt := "Expected: '%T(%+v)'\nActual:   '%T(%+v)'\n(Should be equivalent)!"
	return fmt.Sprintf(errFmt, expected, expected, actual, actual)

}

func TestInitializeRouter(t *testing.T) {
	Convey("InitializeRouter should", t, func() {

		Convey("for invalid input", func() {

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
				h, err := InitializeRouter(nil, nil)
				So(h, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})

			Convey("return an error when the API registrations are invalid", func() {
				h, err := InitializeRouter(api, nil)
				So(h, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(strings.HasPrefix(err.Error(), "missing"), ShouldBeTrue)
			})
		})

		Convey("for valid input", func() {
			spec := &swagger.Spec{
				Consumes: []string{"application/json"},
				Produces: []string{"application/json"},
				Paths: swagger.Paths{
					Paths: map[string]swagger.PathItem{
						"/": swagger.PathItem{
							Get:     &swagger.Operation{ID: "doGet"},
							Post:    &swagger.Operation{ID: "doNew"},
							Options: &swagger.Operation{ID: "doOptions"},
							Head:    &swagger.Operation{ID: "doHead"},
						},
						"/{id}": swagger.PathItem{
							Put:    &swagger.Operation{ID: "doReplace"},
							Patch:  &swagger.Operation{ID: "doUpdate"},
							Delete: &swagger.Operation{ID: "doDelete"},
						},
					},
				},
			}

			api := NewAPI(spec)
			api.RegisterOperation("doGet", emptyOperationHandler)
			api.RegisterOperation("doNew", emptyOperationHandler)
			api.RegisterOperation("doOptions", emptyOperationHandler)
			api.RegisterOperation("doHead", emptyOperationHandler)
			api.RegisterOperation("doReplace", emptyOperationHandler)
			api.RegisterOperation("doUpdate", emptyOperationHandler)
			api.RegisterOperation("doDelete", emptyOperationHandler)

			router := DefaultRouter().(*defaultRouter)
			h, err := InitializeRouter(api, router)
			So(err, ShouldBeNil)
			So(h, ShouldNotBeNil)
			So(len(router.handlers), ShouldEqual, 7)

			expectedMethods := []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}
			var seenMethods []string
			for _, h := range router.handlers {
				seenMethods = append(seenMethods, h.Method)
			}
			So(seenMethods, ShouldResemble, expectedMethods)
		})
	})
}
