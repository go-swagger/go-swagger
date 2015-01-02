package spec

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var infoJSON = `{
	"description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
	"title": "Swagger Sample API",
	"termsOfService": "http://helloreverb.com/terms/",
	"contact": {
		"name": "wordnik api team",
		"url": "http://developer.wordnik.com"
	},
	"license": {
		"name": "Creative Commons 4.0 International",
		"url": "http://creativecommons.org/licenses/by/4.0/"
	},
	"version": "1.0.9-abcd",
	"x-framework": "go-swagger"
}`

var info = Info{
	infoProps: infoProps{
		Version:        "1.0.9-abcd",
		Title:          "Swagger Sample API",
		Description:    "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
		TermsOfService: "http://helloreverb.com/terms/",
		Contact:        &ContactInfo{Name: "wordnik api team", URL: "http://developer.wordnik.com"},
		License:        &License{Name: "Creative Commons 4.0 International", URL: "http://creativecommons.org/licenses/by/4.0/"},
	},
	vendorExtensible: vendorExtensible{map[string]interface{}{"x-framework": "go-swagger"}},
}

func TestIntegrationInfo(t *testing.T) {
	Convey("all fields of info should", t, func() {
		Convey("serialize to JSON", func() {
			b, err := json.MarshalIndent(info, "", "\t")
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, infoJSON)
		})

		Convey("deserialize from JSON", func() {
			actual := Info{}
			err := json.Unmarshal([]byte(infoJSON), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldBeEquivalentTo, info)
		})

	})
}
