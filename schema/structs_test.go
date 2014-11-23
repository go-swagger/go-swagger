package schema

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

func validateJSON(actual interface{}, expected ...interface{}) string {
	ser, err := json.Marshal(actual)
	if err != nil {
		return err.Error()
	}
	exp := expected[0].(string)
	return ShouldEqual(string(ser), exp)
}

func validateYAML(actual interface{}, expected ...interface{}) string {
	ser, err := yaml.Marshal(actual)
	if err != nil {
		return err.Error()
	}
	exp := expected[0].(string)
	return ShouldEqual(string(ser), exp)
}

func TestSerialization(t *testing.T) {
	Convey("Swagger should serialize", t, func() {

		Convey("a string or array property", func() {
			Convey("when string", func() {
				obj := StringOrArray{Single: "hello"}

				Convey("for json returns quoted string", func() {
					So(obj, validateJSON, "\"hello\"")
				})
				Convey("for yaml returns quoted string", func() {
					So(obj, validateYAML, "hello\n")
				})
			})

			Convey("when slice", func() {
				obj := StringOrArray{Multi: []string{"hello", "world", "and", "stuff"}}
				Convey("for json returns an array of strings", func() {
					So(obj, validateJSON, "[\"hello\",\"world\",\"and\",\"stuff\"]")
				})
				Convey("for yaml returns an array of strings", func() {
					So(obj, validateYAML, "- hello\n- world\n- and\n- stuff\n")
				})
			})

			Convey("when empty", func() {
				obj := StringOrArray{}
				Convey("for json returns an empty array", func() {
					So(obj, validateJSON, "null")
				})
				Convey("for yaml returns an emtpy array", func() {
					So(obj, validateYAML, "[]\n")
				})
			})
		})

		Convey("a schema or array property", func() {
			Convey("when string", func() {
				obj := SchemaOrArray{Single: &Schema{Type: &StringOrArray{Single: "string"}}}

				Convey("for json returns quoted string", func() {
					So(obj, validateJSON, "{\"type\":\"string\"}")
				})
				Convey("for yaml returns quoted string", func() {
					So(obj, validateYAML, "type: string\n")
				})
			})

			Convey("when slice", func() {
				obj := SchemaOrArray{
					Multi: []Schema{
						Schema{Type: &StringOrArray{Single: "string"}},
						Schema{Type: &StringOrArray{Single: "string"}},
					},
				}
				Convey("for json returns an array of strings", func() {
					So(obj, validateJSON, "[{\"type\":\"string\"},{\"type\":\"string\"}]")
				})
				Convey("for yaml returns an array of strings", func() {
					So(obj, validateYAML, "- type: string\n- type: string\n")
				})
			})

			Convey("when empty", func() {
				obj := SchemaOrArray{}
				Convey("for json returns an empty array", func() {
					So(obj, validateJSON, "null")
				})
				Convey("for yaml returns an emtpy array", func() {
					So(obj, validateYAML, "[]\n")
				})
			})
		})
	})
}
