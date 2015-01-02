package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

var contactInfoJSON = `{"name":"wordnik api team","url":"http://developer.wordnik.com","email":"some@mailayada.dkdkd"}`
var contactInfoYAML = `name: wordnik api team
url: http://developer.wordnik.com
email: some@mailayada.dkdkd
`
var contactInfo = ContactInfo{
	Name:  "wordnik api team",
	URL:   "http://developer.wordnik.com",
	Email: "some@mailayada.dkdkd",
}

func TestIntegrationContactInfo(t *testing.T) {
	Convey("all fields of contact info should", t, func() {
		Convey("serialize to JSON", func() {
			b, err := json.Marshal(contactInfo)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, contactInfoJSON)
		})

		Convey("serialize to YAML", func() {
			b, err := yaml.Marshal(contactInfo)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, contactInfoYAML)
		})

		Convey("deserialize from JSON", func() {
			actual := ContactInfo{}
			err := json.Unmarshal([]byte(contactInfoJSON), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, contactInfo)
		})

		Convey("deserialize from YAML", func() {
			actual := ContactInfo{}
			err := yaml.Unmarshal([]byte(contactInfoYAML), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, contactInfo)
		})
	})
}
