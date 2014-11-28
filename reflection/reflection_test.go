package reflection

import (
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
)

func TestMarshalling(t *testing.T) {
	c.Convey("marshals a map from", t, func() {

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
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
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
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
			c.So(res["field"], c.ShouldEqual, "the value")
		})

		c.Convey("a simple struct skipping ignored fields", func() {
			type T1 struct {
				A string
				B string `swagger:"-"`
			}
			obj := &T1{"the value", "another value"}
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
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
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
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
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
			c.So(res["field"], c.ShouldEqual, "the value")
			_, ok := res["key"]
			c.So(ok, c.ShouldBeFalse)
		})

		c.Convey("a struct with a field tagged as byValue doesn't expand the struct", func() {
			type T1 struct {
				A time.Time `swagger:"field,byValue"`
			}
			obj := &T1{time.Now()}
			res, err := MarshalMap(obj)
			c.So(err, c.ShouldBeNil)
			c.So(res["field"], c.ShouldHappenOnOrBefore, obj.A)
		})
	})
}
