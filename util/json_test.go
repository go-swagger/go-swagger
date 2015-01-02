package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONConcatenation(t *testing.T) {
	Convey("JSON concatenation should", t, func() {

		Convey("return nil when there are args provided", func() {
			res := ConcatJSON()
			So(res, ShouldBeNil)
		})

		Convey("return the first blob when there is only one", func() {
			expected := []byte(`{"id":1}`)
			res := ConcatJSON(expected)
			So(string(res), ShouldEqual, string(expected))
		})

		Convey("concatenate 2 items", func() {
			Convey("with only empty elements", func() {
				Convey("for an object", func() {
					expected := []byte(`{}`)
					res := ConcatJSON([]byte(`{}`), []byte(`{}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[]`)
					res := ConcatJSON([]byte(`[]`), []byte(`[]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with both having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"id":1,"name":"Rachel"}`)
					res := ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"id":1},{"name":"Rachel"}]`)
					res := ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with only the last element having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"name":"Rachel"}`)
					res := ConcatJSON([]byte(`{}`), []byte(`{"name":"Rachel"}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"name":"Rachel"}]`)
					res := ConcatJSON([]byte(`[]`), []byte(`[{"name":"Rachel"}]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with only the first element having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"id":1}`)
					res := ConcatJSON([]byte(`{"id":1}`), []byte(`{}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"id":1}]`)
					res := ConcatJSON([]byte(`[{"id":1}]`), []byte(`[]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})
		})

		Convey("concatenate more than 2 items", func() {
			Convey("with only empty elements", func() {
				Convey("for an object", func() {
					expected := []byte(`{}`)
					res := ConcatJSON([]byte(`{}`), []byte(`{}`), []byte(`{}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[]`)
					res := ConcatJSON([]byte(`[]`), []byte(`[]`), []byte(`[]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with all having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"id":1,"name":"Rachel","age":32}`)
					res := ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`), []byte(`{"age":32}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"id":1},{"name":"Rachel"},{"age":32}]`)
					res := ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`), []byte(`[{"age":32}]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with only the last 2 elements having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"name":"Rachel","age":32}`)
					res := ConcatJSON([]byte(`{}`), []byte(`{"name":"Rachel"}`), []byte(`{"age":32}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"name":"Rachel"},{"age":32}]`)
					res := ConcatJSON([]byte(`[]`), []byte(`[{"name":"Rachel"}]`), []byte(`[{"age":32}]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with only the first and last elements having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"id":1,"age":32}`)
					res := ConcatJSON([]byte(`{"id":1}`), []byte(`{}`), []byte(`{"age":32}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"id":1},{"age":32}]`)
					res := ConcatJSON([]byte(`[{"id":1}]`), []byte(`[]`), []byte(`[{"age":32}]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})

			Convey("with only the first element having data", func() {
				Convey("for an object", func() {
					expected := []byte(`{"id":1,"name":"Rachel"}`)
					res := ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`), []byte(`{}`))
					So(string(res), ShouldEqual, string(expected))
				})
				Convey("for an array", func() {
					expected := []byte(`[{"id":1},{"name":"Rachel"}]`)
					res := ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`), []byte(`[]`))
					So(string(res), ShouldEqual, string(expected))
				})
			})
		})
	})
}
