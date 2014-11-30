package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

var contactInfoJSON = `{"email":"some@mailayada.dkdkd","name":"wordnik api team","url":"http://developer.wordnik.com"}`
var contactInfoYAML = `email: some@mailayada.dkdkd
name: wordnik api team
url: http://developer.wordnik.com
`
var contactInfo = ContactInfo{
	Email: "some@mailayada.dkdkd",
	Name:  "wordnik api team",
	URL:   "http://developer.wordnik.com",
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
