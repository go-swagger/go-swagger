package swagger

import (
	"fmt"
	"reflect"

	"github.com/casualjim/go-swagger/reflection"
	// . "github.com/smartystreets/goconvey/convey"
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

// func TestInitializeRouter(t *testing.T) {
// 	Convey("api.Handler should", t, func() {

// 		Convey("for invalid input", func() {
// 			specStr := `{
//   "consumes": ["application/json"],
//   "produces": ["application/json"],
//   "paths": {
// 		"/": {
// 			"get": {
// 				"operationId": "someOperation",
// 				"consumes": ["application/x-yaml"],
// 				"produces": ["application/x-yaml"]
// 			}
// 		}
//   }
// }`
// 			spec, err := spec.New([]byte(specStr), "")
// 			So(err, ShouldBeNil)
// 			So(spec, ShouldNotBeNil)

// 			api := NewAPI(spec)

// 			Convey("return an error when the API registrations are invalid", func() {
// 				h, err := api.Handler(nil)
// 				So(h, ShouldBeNil)
// 				So(err, ShouldNotBeNil)
// 				So(strings.HasPrefix(err.Error(), "missing"), ShouldBeTrue)
// 			})
// 		})

// 		Convey("for valid input", func() {
// 			specStr := `{
// 		  "consumes": ["application/json"],
// 		  "produces": ["application/json"],
// 		  "paths": {
// 				"/": {
// 					"get": {"operationId": "doGet"},
// 					"post": {"operationId": "doNew"},
// 					"options": {"operationId": "doOptions"},
// 					"head": {"operationId": "doHead"}
// 				},
// 				"/{id}": {
// 					"put": {"operationId": "doReplace"},
// 					"patch": {"operationId": "doUpdate"},
// 					"delete": {"operationId": "doDelete"}
// 				}
// 			}
// 		}`
// 			spec, err := spec.New([]byte(specStr), "")
// 			So(err, ShouldBeNil)
// 			So(spec, ShouldNotBeNil)

// 			api := NewAPI(spec)
// 			api.RegisterOperation("doGet", new(stubOperationHandler))
// 			api.RegisterOperation("doNew", new(stubOperationHandler))
// 			api.RegisterOperation("doOptions", new(stubOperationHandler))
// 			api.RegisterOperation("doHead", new(stubOperationHandler))
// 			api.RegisterOperation("doReplace", new(stubOperationHandler))
// 			api.RegisterOperation("doUpdate", new(stubOperationHandler))
// 			api.RegisterOperation("doDelete", new(stubOperationHandler))

// 			router := DefaultRouter().(*defaultRouter)
// 			h, err := api.Handler(router)
// 			So(err, ShouldBeNil)
// 			So(h, ShouldNotBeNil)
// 			So(len(router.handlers), ShouldEqual, 7)

// 			expectedMethods := []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}
// 			sort.Sort(sort.StringSlice(expectedMethods))
// 			var seenMethods []string
// 			for _, h := range router.handlers {
// 				seenMethods = append(seenMethods, h.Method)
// 			}
// 			sort.Sort(sort.StringSlice(seenMethods))
// 			So(seenMethods, ShouldResemble, expectedMethods)
// 		})
// 	})
// }
