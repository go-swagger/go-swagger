package reflection

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

func TestUnmarshalling(t *testing.T) {
	c.Convey("Unmarshalling a map should", t, func() {

		c.Convey("convert values to a struct from a map", func() {
			data := map[string]interface{}{
				"A": map[string]string{"AA": "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldResemble, struct{ AA string }{AA: "value"})
		})

		c.Convey("convert values to a struct from a struct", func() {
			data := map[string]interface{}{
				"A": struct{ AA string }{AA: "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldResemble, data["A"])
		})

		c.Convey("convert values to a struct from a struct pointer", func() {
			data := map[string]interface{}{
				"A": &struct{ AA string }{AA: "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldResemble, *data["A"].(*struct{ AA string }))
		})

		c.Convey("convert values to a slice", func() {
			data := map[string]interface{}{
				"A": []string{"hello", "world"},
			}
			actual := &struct{ A []string }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldResemble, data["A"])
		})

		c.Convey("convert values to a map of strings", func() {
			data := map[string]interface{}{
				"A": map[string]string{"AA": "value"},
			}
			actual := &struct{ A map[string]string }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldResemble, map[string]string{"AA": "value"})
		})

		c.Convey("convert values to interface", func() {
			data := map[string]interface{}{"A": "value"}
			actual := &struct{ A interface{} }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldEqual, "value")
		})

		c.Convey("convert values to string", func() {
			data := map[string]interface{}{
				"A": "value",
				"B": true,
				"C": 1,
				"D": uint(1),
				"E": 1.0,
			}
			actual := &struct {
				A string
				B string
				C string
				D string
				E string
			}{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldEqual, "value")
			c.So(actual.B, c.ShouldEqual, "true")
			c.So(actual.C, c.ShouldEqual, "1")
			c.So(actual.D, c.ShouldEqual, "1")
			c.So(actual.E, c.ShouldEqual, "1")
		})

		c.Convey("convert values to bool", func() {
			data := map[string]interface{}{
				"A": true,
				"B": "1",
				"C": "true",
				"D": "",
				"E": 1,
				"F": 0,
				"G": uint(1),
				"H": uint(0),
				"I": float32(1),
				"J": float32(0),
			}
			actual := &struct{ A, B, C, D, E, F, G, H, I, J bool }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldBeTrue)
			c.So(actual.B, c.ShouldBeTrue)
			c.So(actual.C, c.ShouldBeTrue)
			c.So(actual.D, c.ShouldBeFalse)
			c.So(actual.E, c.ShouldBeTrue)
			c.So(actual.F, c.ShouldBeFalse)
			c.So(actual.G, c.ShouldBeTrue)
			c.So(actual.H, c.ShouldBeFalse)
			c.So(actual.I, c.ShouldBeTrue)
			c.So(actual.J, c.ShouldBeFalse)
		})

		c.Convey("convert values to int", func() {
			data := map[string]interface{}{
				"A": 1,
				"B": int(1),
				"C": int8(1),
				"D": int16(1),
				"E": int32(1),
				"F": int64(1),
				"G": uint(1),
				"H": uint8(1),
				"I": uint16(1),
				"J": uint32(1),
				"K": uint64(1),
				"L": float32(1),
				"M": float64(1),
				"N": "1",
				"O": true,
			}
			actual := &struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O int32 }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldEqual, 1)
			c.So(actual.B, c.ShouldEqual, 1)
			c.So(actual.C, c.ShouldEqual, 1)
			c.So(actual.D, c.ShouldEqual, 1)
			c.So(actual.E, c.ShouldEqual, 1)
			c.So(actual.F, c.ShouldEqual, 1)
			c.So(actual.G, c.ShouldEqual, 1)
			c.So(actual.H, c.ShouldEqual, 1)
			c.So(actual.J, c.ShouldEqual, 1)
			c.So(actual.K, c.ShouldEqual, 1)
			c.So(actual.L, c.ShouldEqual, 1)
			c.So(actual.M, c.ShouldEqual, 1)
			c.So(actual.N, c.ShouldEqual, 1)
			c.So(actual.O, c.ShouldEqual, 1)
		})

		c.Convey("convert values to uint", func() {
			data := map[string]interface{}{
				"A": 1,
				"B": int(1),
				"C": int8(1),
				"D": int16(1),
				"E": int32(1),
				"F": int64(1),
				"G": uint(1),
				"H": uint8(1),
				"I": uint16(1),
				"J": uint32(1),
				"K": uint64(1),
				"L": float32(1),
				"M": float64(1),
				"N": "1",
				"O": true,
			}
			actual := &struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O uint32 }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldEqual, 1)
			c.So(actual.B, c.ShouldEqual, 1)
			c.So(actual.C, c.ShouldEqual, 1)
			c.So(actual.D, c.ShouldEqual, 1)
			c.So(actual.E, c.ShouldEqual, 1)
			c.So(actual.F, c.ShouldEqual, 1)
			c.So(actual.G, c.ShouldEqual, 1)
			c.So(actual.H, c.ShouldEqual, 1)
			c.So(actual.J, c.ShouldEqual, 1)
			c.So(actual.K, c.ShouldEqual, 1)
			c.So(actual.L, c.ShouldEqual, 1)
			c.So(actual.M, c.ShouldEqual, 1)
			c.So(actual.N, c.ShouldEqual, 1)
			c.So(actual.O, c.ShouldEqual, 1)
		})

		c.Convey("convert values to float", func() {
			data := map[string]interface{}{
				"A": 1,
				"B": int(1),
				"C": int8(1),
				"D": int16(1),
				"E": int32(1),
				"F": int64(1),
				"G": uint(1),
				"H": uint8(1),
				"I": uint16(1),
				"J": uint32(1),
				"K": uint64(1),
				"L": float32(1),
				"M": float64(1),
				"N": "1",
				"O": true,
			}
			actual := &struct{ A, B, C, D, E, F, G, H, I, J, K, L, M, N, O float32 }{}
			c.So(UnmarshalMap(data, actual), c.ShouldBeNil)
			c.So(actual.A, c.ShouldEqual, 1)
			c.So(actual.B, c.ShouldEqual, 1)
			c.So(actual.C, c.ShouldEqual, 1)
			c.So(actual.D, c.ShouldEqual, 1)
			c.So(actual.E, c.ShouldEqual, 1)
			c.So(actual.F, c.ShouldEqual, 1)
			c.So(actual.G, c.ShouldEqual, 1)
			c.So(actual.H, c.ShouldEqual, 1)
			c.So(actual.J, c.ShouldEqual, 1)
			c.So(actual.K, c.ShouldEqual, 1)
			c.So(actual.L, c.ShouldEqual, 1)
			c.So(actual.M, c.ShouldEqual, 1)
			c.So(actual.N, c.ShouldEqual, 1)
			c.So(actual.O, c.ShouldEqual, 1)
		})

		c.Convey("fail when the target is not a pointer", func() {
			c.So(UnmarshalMap(map[string]interface{}{"A": "B"}, struct{ A string }{}), c.ShouldNotBeNil)
		})

		c.Convey("read field with name override", func() {
			obj := &struct {
				A string `swagger:"field"`
			}{}
			err := UnmarshalMap(map[string]interface{}{"field": "value"}, obj)
			c.So(err, c.ShouldBeNil)
			c.So(obj.A, c.ShouldEqual, "value")
		})

		c.Convey("skip fields that are set to ignore", func() {
			obj := &struct {
				A string `swagger:"field"`
				B string `swagger:"-"`
			}{}
			err := UnmarshalMap(map[string]interface{}{"field": "value", "B": "value"}, obj)
			c.So(err, c.ShouldBeNil)
			c.So(obj.A, c.ShouldEqual, "value")
			c.So(obj.B, c.ShouldBeEmpty)
		})
	})
}
