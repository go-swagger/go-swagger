package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/v1/yaml"
)

var contactInfoJson = `{"email":"some@mailayada.dkdkd","name":"wordnik api team","url":"http://developer.wordnik.com"}`
var contactInfoYaml = `name: wordnik api team
url: http://developer.wordnik.com
email: some@mailayada.dkdkd
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
			So(string(b), ShouldEqual, contactInfoJson)
		})

		Convey("serialize to YAML", func() {
			b, err := yaml.Marshal(contactInfo)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, contactInfoYaml)
		})

		Convey("deserialize from JSON", func() {
			actual := ContactInfo{}
			err := json.Unmarshal([]byte(contactInfoJson), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, contactInfo)
		})

		Convey("deserialize from YAML", func() {
			actual := ContactInfo{}
			err := yaml.Unmarshal([]byte(contactInfoYaml), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, contactInfo)
		})
	})
}
