//go:build testintegration

package docs

import "swagger/api"

// swagger:route POST /foobar foobar-tag idOfFoobarEndpoint
// Foobar does some amazing stuff.
// responses:
//   200: foobarResponse

// This text will appear as description of your response body.
// swagger:response foobarResponse
type foobarResponseWrapper struct {
	// in:body
	Body api.FooBarResponse
}

// swagger:parameters idOfFoobarEndpoint
type foobarParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body api.FooBarRequest
}
