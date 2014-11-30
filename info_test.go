package swagger

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/v1/yaml"
)

var infoJson = `{
	"contact": {
		"name": "wordnik api team",
		"url": "http://developer.wordnik.com"
	},
	"description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
	"license": {
		"name": "Creative Commons 4.0 International",
		"url": "http://creativecommons.org/licenses/by/4.0/"
	},
	"termsOfService": "http://helloreverb.com/terms/",
	"title": "Swagger Sample API",
	"version": "1.0.9-abcd",
	"x-framework": "go-swagger"
}`

var infoYaml = `x-framework: go-swagger
description: A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification
title: Swagger Sample API
termsOfService: http://helloreverb.com/terms/
contact:
  name: wordnik api team
	url: http://developer.wordnik.com
license:
	name: Creative Commons 4.0 International
	url: http://creativecommons.org/licenses/by/4.0/
version: 1.0.9-abcd
`

var info = Info{
	Version:        "1.0.9-abcd",
	Title:          "Swagger Sample API",
	Description:    "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
	TermsOfService: "http://helloreverb.com/terms/",
	Extensions:     map[string]interface{}{"x-framework": "go-swagger"},
	Contact:        &ContactInfo{Name: "wordnik api team", URL: "http://developer.wordnik.com"},
	License:        &License{Name: "Creative Commons 4.0 International", URL: "http://creativecommons.org/licenses/by/4.0/"},
}

func TestIntegrationInfo(t *testing.T) {
	Convey("all fields of info should", t, func() {
		Convey("serialize to JSON", func() {
			b, err := json.MarshalIndent(info, "", "\t")
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, infoJson)
		})

		Convey("serialize to YAML", func() {
			b, err := yaml.Marshal(contactInfo)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, contactInfoYaml)
		})

		Convey("deserialize from JSON", func() {
			actual := Info{}
			err := json.Unmarshal([]byte(infoJson), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, info)
		})

		Convey("deserialize from YAML", func() {
			actual := ContactInfo{}
			err := yaml.Unmarshal([]byte(contactInfoYaml), &actual)
			So(err, ShouldBeNil)
			So(actual, ShouldResemble, contactInfo)
		})
	})
}
