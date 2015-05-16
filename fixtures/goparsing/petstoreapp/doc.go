// Package petstoreapp implements a sample petstore api.
// Title: Petstore API
//
// TOS: there are no TOS at this moment, use at your own risk we take no responsibility
//
// Version: 0.0.1
//
// License: MIT http://opensource.org/licenses/MIT
//
// Contact: John Doe<john.doe@example.com> http://john.doe.com
//
// Description:
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Consumes:
// application/json
//
// Produces:
// application/json
//
// Schemes: http, https
//
// Host: localhost:8080/api
package petstoreapp

import "github.com/casualjim/go-swagger/spec"

// +swagger.meta
var (
	APIVersion = "0.0.1"

	License = spec.License{
		Name: "MIT",
		URL:  "http://opensource.org/licenses/MIT",
	}
)
