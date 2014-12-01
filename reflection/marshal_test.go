package reflection

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
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
		if IsZero(expectedValue) && IsZero(reflect.ValueOf(actual)) {
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

type customMarshalling struct {
	A string
}

func (c customMarshalling) MarshalMap() map[string]interface{} {
	return map[string]interface{}{"field": c.A}
}

func TestMarshalling(t *testing.T) {
	c.Convey("marshals a map from", t, func() {

		c.Convey("a simple struct with slices", func() {
			type T1 struct {
				A []string
				B []int
				C []int8
				D []int16
				E []int32
				F []int64
				G []float32
				H []float64
			}
			obj := &T1{
				A: []string{"the value"},
				B: []int{1},
				C: []int8{1},
				D: []int16{1},
				E: []int32{1},
				F: []int64{1},
				G: []float32{1},
				H: []float64{1},
			}
			res := MarshalMap(obj)
			c.So(res["A"], ShouldBeEquivalentTo, []interface{}{"the value"})
			c.So(res["B"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["C"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["D"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["E"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["F"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["G"], ShouldBeEquivalentTo, []interface{}{1})
			c.So(res["H"], ShouldBeEquivalentTo, []interface{}{1})
		})
		c.Convey("a simple struct", func() {
			type T1 struct {
				A string
				B int
				C int8
				D int16
				E int32
				F int64
				G float32
				H float64
			}
			obj := &T1{
				A: "the value",
				B: 1,
				C: 1,
				D: 1,
				E: 1,
				F: 1,
				G: 1,
				H: 1,
			}
			res := MarshalMap(obj)
			c.So(res["A"], c.ShouldEqual, "the value")
			c.So(res["B"], c.ShouldEqual, 1)
			c.So(res["C"], c.ShouldEqual, 1)
			c.So(res["D"], c.ShouldEqual, 1)
			c.So(res["E"], c.ShouldEqual, 1)
			c.So(res["F"], c.ShouldEqual, 1)
			c.So(res["G"], c.ShouldEqual, 1)
			c.So(res["H"], c.ShouldEqual, 1)
		})

		c.Convey("a simple struct with field name override", func() {
			type T1 struct {
				A string `swagger:"field"`
			}
			obj := &T1{"the value"}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldEqual, "the value")
		})

		c.Convey("a simple struct skipping ignored fields", func() {
			type T1 struct {
				A string
				B string `swagger:"-"`
			}
			obj := &T1{"the value", "another value"}
			res := MarshalMap(obj)
			c.So(res["A"], c.ShouldEqual, "the value")
			_, ok := res["B"]
			c.So(ok, c.ShouldBeFalse)
		})

		c.Convey("a simple struct including fields with omitempty if they have a value", func() {
			type T1 struct {
				A string `swagger:"field"`
				B string `swagger:"key,omitempty"`
			}
			obj := &T1{"the value", "another value"}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldEqual, "the value")
			b, ok := res["key"]
			c.So(ok, c.ShouldBeTrue)
			c.So(b, c.ShouldEqual, "another value")
		})

		c.Convey("a simple struct skipping fields with omitempty if they don't have a value", func() {
			type T1 struct {
				A string `swagger:"field"`
				B string `swagger:"key,omitempty"`
			}
			obj := &T1{"the value", ""}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldEqual, "the value")
			_, ok := res["key"]
			c.So(ok, c.ShouldBeFalse)
		})

		c.Convey("a struct with a field tagged as byValue doesn't expand the struct", func() {
			type T1 struct {
				A time.Time `swagger:"field,byValue"`
			}
			obj := &T1{time.Now()}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldHappenOnOrBefore, obj.A)
		})

		c.Convey("a struct with a custom marshaller should use the marshaller", func() {
			obj := &customMarshalling{"a value"}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldEqual, "a value")
		})

		c.Convey("a struct field with a custom marshaller should use the marshaller", func() {
			type T1 struct {
				B *customMarshalling
			}
			obj := &T1{&customMarshalling{"a value"}}
			res := MarshalMap(obj)
			c.So(res["B"].(map[string]interface{})["field"], c.ShouldEqual, "a value")
		})

		c.Convey("a map and convert it to map[string]interface", func() {
			obj := map[string]customMarshalling{"field": customMarshalling{"a value"}}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldResemble, map[string]interface{}{"field": "a value"})
		})

		c.Convey("a map and convert it to map[string]interface", func() {
			obj := map[string]*customMarshalling{"field": &customMarshalling{"a value"}}
			res := MarshalMap(obj)
			c.So(res["field"], c.ShouldResemble, map[string]interface{}{"field": "a value"})
		})

		c.Convey("a struct with a slice property", func() {
			type T1 struct {
				C []string
			}
			obj := &T1{[]string{"first", "second"}}
			res := MarshalMap(obj)
			c.So(res["C"], c.ShouldResemble, []interface{}{"first", "second"})
		})

		c.Convey("a struct with a slice property that implements a map marshaller", func() {
			type T1 struct {
				C []customMarshalling
			}
			obj := &T1{[]customMarshalling{customMarshalling{"first"}, customMarshalling{"second"}}}
			res := MarshalMap(obj)
			c.So(res["C"], c.ShouldResemble, []interface{}{map[string]interface{}{"field": "first"}, map[string]interface{}{"field": "second"}})
		})
	})
}
