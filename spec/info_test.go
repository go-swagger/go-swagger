// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
