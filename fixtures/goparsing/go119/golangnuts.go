package go119

import (
	"net/http"
)

// The some func was taken from: https://groups.google.com/g/golang-nuts/c/sY7QFzfSMT4/m/Rk6WVRJEAQAJ

// SomeFunc do something
//
// swagger:operation POST /api/v1/somefunc someFunc
//
// Do something
//
//	---
//	x-codeSamples:
//	- lang: 'curl'
//	  source: |
//	    curl -u "${LOGIN}:${PASSWORD}" -d '{"key": "value"}' -X POST   "https://{host}/api/v1/somefunc"
//	    curl -u "${LOGIN}:${PASSWORD}" -d '{"key2": "value2"}' -X POST   "https://{host}/api/v1/somefunc"
//	responses:
//	  '200':
//	    description: "Some func"
//	    examples:
//	      application/json:
//	        key: "value"
//	  '400':
//	    $ref: "#/responses/ErrorResponse"
//	  '503':
//	    $ref: "#/responses/ErrorResponse"
func SomeFunct(rw http.ResponseWriter, req *http.Request) {
	/// do something
}

// SomeFunc do something
//
// swagger:operation POST /api/v1/somefuncTabs someFuncTabs
//
// Do something
//
//	---
//	x-codeSamples:
//	- lang: 'curl'
//		source: |
//			curl -u "${LOGIN}:${PASSWORD}" -d '{"key": "value"}' -X POST   "https://{host}/api/v1/somefunc"
//			curl -u "${LOGIN}:${PASSWORD}" -d '{"key2": "value2"}' -X POST   "https://{host}/api/v1/somefunc"
//	responses:
//		'200':
//			description: "Some func"
//			examples:
//				application/json:
//					key: "value"
//		'400':
//			$ref: "#/responses/ErrorResponse"
//		'503':
//			$ref: "#/responses/ErrorResponse"
func SomeFunctabs(rw http.ResponseWriter, req *http.Request) {
	/// do something
}
