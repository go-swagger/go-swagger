package reflection

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type CustomUnmarshaller struct {
	Field string
}

func (c *CustomUnmarshaller) UnmarshalMap(data interface{}) error {
	v := data.(map[string]interface{})
	c.Field = v["A"].(string)
	return nil
}

func TestUnmarshalling(t *testing.T) {
	Convey("Unmarshalling a map should", t, func() {
		Convey("convert a map with interface keys for an interface", func() {
			data := map[string]interface{}{
				"AA": map[interface{}]interface{}{
					"A": []string{"value"},
				},
			}
			actual := new(struct{ AA map[string][]interface{} })
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.AA, ShouldResemble, map[string][]interface{}{"A": []interface{}{"value"}})
		})

		Convey("convert a map with interface keys", func() {
			data := map[string]interface{}{
				"AA": map[interface{}]interface{}{
					"A": "value",
				},
			}
			actual := new(struct{ AA map[string]string })
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.AA, ShouldResemble, map[string]string{"A": "value"})
		})

		Convey("use custom unmarshaller when top level", func() {
			data := map[string]interface{}{
				"A": "value",
			}
			actual := new(CustomUnmarshaller)
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.Field, ShouldEqual, "value")
		})

		Convey("use custom unmarshaller when used as map value", func() {
			data := map[string]interface{}{
				"AA": map[string]interface{}{
					"A": map[string]interface{}{"A": "value"},
				},
			}
			actual := new(struct{ AA map[string]CustomUnmarshaller })
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.AA["A"].Field, ShouldEqual, "value")
		})
		Convey("use custom unmarshaller when used as struct", func() {
			data := map[string]interface{}{
				"AA": map[string]interface{}{
					"A": "value",
				},
			}
			actual := new(struct{ AA CustomUnmarshaller })
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.AA.Field, ShouldEqual, "value")
		})

		Convey("use custom unmarshaller when used as pointer to a struct", func() {
			data := map[string]interface{}{
				"AA": map[string]interface{}{
					"A": "value",
				},
			}
			actual := new(struct{ AA *CustomUnmarshaller })
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.AA.Field, ShouldEqual, "value")
		})

		Convey("convert values to a struct from a map", func() {
			data := map[string]interface{}{
				"A": map[string]string{"AA": "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldResemble, struct{ AA string }{AA: "value"})
		})

		Convey("convert values to a struct from a struct", func() {
			data := map[string]interface{}{
				"A": struct{ AA string }{AA: "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldResemble, data["A"])
		})

		Convey("convert values to a struct from a struct pointer", func() {
			data := map[string]interface{}{
				"A": &struct{ AA string }{AA: "value"},
			}
			actual := &struct{ A struct{ AA string } }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldResemble, *data["A"].(*struct{ AA string }))
		})

		Convey("convert values to a slice", func() {
			data := map[string]interface{}{
				"A": []string{"hello", "world"},
			}
			actual := &struct{ A []string }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldResemble, data["A"])
		})

		Convey("convert values to a map of strings", func() {
			data := map[string]interface{}{
				"A": map[string]string{"AA": "value"},
			}
			actual := &struct{ A map[string]string }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldResemble, map[string]string{"AA": "value"})
		})

		Convey("convert values to interface", func() {
			data := map[string]interface{}{"A": "value"}
			actual := &struct{ A interface{} }{}
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldEqual, "value")
		})

		Convey("convert values to string", func() {
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
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldEqual, "value")
			So(actual.B, ShouldEqual, "true")
			So(actual.C, ShouldEqual, "1")
			So(actual.D, ShouldEqual, "1")
			So(actual.E, ShouldEqual, "1")
		})

		Convey("convert values to bool", func() {
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
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldBeTrue)
			So(actual.B, ShouldBeTrue)
			So(actual.C, ShouldBeTrue)
			So(actual.D, ShouldBeFalse)
			So(actual.E, ShouldBeTrue)
			So(actual.F, ShouldBeFalse)
			So(actual.G, ShouldBeTrue)
			So(actual.H, ShouldBeFalse)
			So(actual.I, ShouldBeTrue)
			So(actual.J, ShouldBeFalse)
		})

		Convey("convert values to int", func() {
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
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldEqual, 1)
			So(actual.B, ShouldEqual, 1)
			So(actual.C, ShouldEqual, 1)
			So(actual.D, ShouldEqual, 1)
			So(actual.E, ShouldEqual, 1)
			So(actual.F, ShouldEqual, 1)
			So(actual.G, ShouldEqual, 1)
			So(actual.H, ShouldEqual, 1)
			So(actual.J, ShouldEqual, 1)
			So(actual.K, ShouldEqual, 1)
			So(actual.L, ShouldEqual, 1)
			So(actual.M, ShouldEqual, 1)
			So(actual.N, ShouldEqual, 1)
			So(actual.O, ShouldEqual, 1)
		})

		Convey("convert values to uint", func() {
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
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldEqual, 1)
			So(actual.B, ShouldEqual, 1)
			So(actual.C, ShouldEqual, 1)
			So(actual.D, ShouldEqual, 1)
			So(actual.E, ShouldEqual, 1)
			So(actual.F, ShouldEqual, 1)
			So(actual.G, ShouldEqual, 1)
			So(actual.H, ShouldEqual, 1)
			So(actual.J, ShouldEqual, 1)
			So(actual.K, ShouldEqual, 1)
			So(actual.L, ShouldEqual, 1)
			So(actual.M, ShouldEqual, 1)
			So(actual.N, ShouldEqual, 1)
			So(actual.O, ShouldEqual, 1)
		})

		Convey("convert values to float", func() {
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
			So(UnmarshalMap(data, actual), ShouldBeNil)
			So(actual.A, ShouldEqual, 1)
			So(actual.B, ShouldEqual, 1)
			So(actual.C, ShouldEqual, 1)
			So(actual.D, ShouldEqual, 1)
			So(actual.E, ShouldEqual, 1)
			So(actual.F, ShouldEqual, 1)
			So(actual.G, ShouldEqual, 1)
			So(actual.H, ShouldEqual, 1)
			So(actual.J, ShouldEqual, 1)
			So(actual.K, ShouldEqual, 1)
			So(actual.L, ShouldEqual, 1)
			So(actual.M, ShouldEqual, 1)
			So(actual.N, ShouldEqual, 1)
			So(actual.O, ShouldEqual, 1)
		})

		Convey("fail when the target is not a pointer", func() {
			So(UnmarshalMap(map[string]interface{}{"A": "B"}, struct{ A string }{}), ShouldNotBeNil)
		})

		Convey("read field with name override", func() {
			obj := &struct {
				A string `swagger:"field"`
			}{}
			err := UnmarshalMap(map[string]interface{}{"field": "value"}, obj)
			So(err, ShouldBeNil)
			So(obj.A, ShouldEqual, "value")
		})

		Convey("skip fields that are set to ignore", func() {
			obj := &struct {
				A string `swagger:"field"`
				B string `swagger:"-"`
			}{}
			err := UnmarshalMap(map[string]interface{}{"field": "value", "B": "value"}, obj)
			So(err, ShouldBeNil)
			So(obj.A, ShouldEqual, "value")
			So(obj.B, ShouldBeEmpty)
		})
	})
}
