Swagger 2.0 [![Circle CI](https://circleci.com/gh/go-swagger/go-swagger/tree/master.svg?style=svg)](https://circleci.com/gh/go-swagger/go-swagger/tree/master) [![Slack Status](https://slackin.goswagger.io/badge.svg)](https://slackin.goswagger.io)
========================

[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE) [![GoDoc](https://godoc.org/github.com/go-swagger/go-swagger?status.svg)](http://godoc.org/github.com/go-swagger/go-swagger)

This API is not stable yet, when it is stable it will be distributed over gopkg.in

There is a code coverage report available in the artifacts section of a build. Unfortunately using coveralls made the
build unstable.

Contains an implementation of Swagger 2.0. It knows how to serialize and deserialize swagger specifications.

Swagger is a simple yet powerful representation of your RESTful API.  
With the largest ecosystem of API tooling on the planet, thousands of developers are supporting Swagger in almost every modern programming language and deployment environment.

With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability. We created Swagger to help fulfill the promise of APIs.

Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal build the best possible services with RESTful APIs. Now in version 2.0, Swagger is more enabling than ever. And it's 100% open source software.

Docs
----


https://go-swagger.github.io

Install or update:

    go get -u github.com/go-swagger/go-swagger/cmd/swagger

The implementation also provides a number of command line tools to help working with swagger.

Currently there is a [spec validator tool](http://go-swagger.github.io/usage/validate/):

		swagger validate https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json

To generate a server for a swagger spec document:

		swagger generate server [-f ./swagger.json] -A [application-name [--principal [principal-name]]

To generate a [client for a swagger spec](http://go-swagger.github.io/generate/client/) document:

		swagger generate client [-f ./swagger.json] -A [application-name [--principal [principal-name]]

To generate a [swagger spec document for a go application](http://go-swagger.github.io/generate/spec/):

		swagger generate spec -o ./swagger.json

Much improved documentation is in the works and will actually explain how to use this tool in much more depth.
To learn about which annotations are available and how to use them for generating a spec from any go application
(generating a spec is not opinionated), you can take a look at the files used for [testing the parser](https://github.com/go-swagger/go-swagger/tree/master/fixtures/goparsing/classification).


There are several other sub commands available for the generate command

		Sub command | Description
		------------|----------------------------------------------------------------------------------
		operation   | generates one or more operations specified in the swagger definition
		model       | generates model files for one or more models specified in the swagger definition
		support     | generates the api builder and the main method
		server      | generates an entire server application
		client      | generates a client for a swagger specification
		spec        | generates a swagger spec document based on go code


Design
------

For now what exists of documentation on how all the pieces fit together, is described in this [doc](design.md)


What's inside?
--------------

For a V1 I want to have this feature set completed:

- [ ] Documentation site
- [x] Play nice with golint, go vet etc.
-	[x] An object model that serializes to swagger yaml or json
-	[x] A tool to work with swagger:
	-	[x] validate a swagger spec document:
    -	[x] validate against jsonschema
    -	[ ] validate extra rules outlined [here](https://github.com/apigee-127/sway/blob/master/docs/versions/2.0.md#semantic-validation)
      - [x] definition can't declare a property that's already defined by one of its ancestors (Error)
      - [x] definition's ancestor can't be a descendant of the same model (Error)
      - [x] each api path should be non-verbatim (account for path param names) unique per method (Error)
      - [ ] each security reference should contain only unique scopes (Warning)
      - [ ] each security scope in a security definition should be unique (Warning)
      - [x] each path parameter should correspond to a parameter placeholder and vice versa (Error)
      - [ ] path parameter declarations do not allow empty names _(`/path/{}` is not valid)_ (Error)
      - [x] each definition property listed in the required array must be defined in the properties of the model (Error)
      - [x] each parameter should have a unique `name` and `in` combination (Error)
      - [x] each operation should have at most 1 parameter of type body (Error)
      - [ ] each operation cannot have both a body parameter and a formData parameter (Error)
      - [x] each operation must have an unique `operationId` (Error)
      - [x] each reference must point to a valid object (Error)
      - [ ] each referencable definition must have references (Warning)
      - [x] every default value that is specified must validate against the schema for that property (Error)
      - [x] every example that is specified must validate against the schema for that property (Error)
      - [x] items property is required for all schemas/definitions of type `array` (Error)
	-	[x] serve swagger UI for any swagger spec file
  - [x] code generation
    -	[x] generate api based on swagger spec
    -	[x] generate go client from a swagger spec
  - [x] spec generation
    -	[x] generate spec document based on the code
      - [x] generate meta data (top level swagger properties) from package docs
      - [x] generate definition entries for models
        - [x] support composed structs out of several embeds
        - [x] support allOf for composed structs
      - [x] generate path entries for routes
      - [x] generate responses from structs
        - [x] support composed structs out of several embeds
      - [x] generate parameters from structs
        - [x] support composed structs out of several embeds
-	[x] Middlewares:
	-	[x] serve spec
	-	[x] routing
	-	[x] validation
	-	[x] additional validation through an interface
	-	[x] authorization
		-	[x] basic auth
		-	[x] api key auth
	-	[x] swagger docs UI
-	[x] Typed JSON Schema implementation
	-	[x] JSON Pointer that knows about structs
	-	[x] JSON Reference that knows about structs
	-	[x] Passes current json schema test suite
-	[x] extended string formats
	-	[x] uuid, uuid3, uuid4, uuid5
	-	[x] email
	-	[x] uri (absolute)
	-	[x] hostname
	-	[x] ipv4
	-	[x] ipv6
	-	[x] credit card
	-	[x] isbn, isbn10, isbn13
	-	[x] social security number
	-	[x] hexcolor
	-	[x] rgbcolor
	-	[x] date
	-	[x] date-time
	-	[x] duration
  - [x] password
  -	[x] custom string formats

Later
-----

After the v1 implementation extra transports are on the roadmap.

Many of these fall under the maybe, perhaps, could be nice to have, might not happen bucket:

- Formats:
	- [ ] custom serializer for XML to support namespaces and prefixes
- Tools:
  - Code generation:
    -	[ ] generate "sensible" random data based on swagger spec
    -	[ ] generate tests based on swagger spec for client
    -	[ ] generate tests based on swagger spec for server
    - [ ] generate markdown representation of swagger spec
    -	[ ] watch swagger spec file and regenerate when modified
  - Spec generation:
    -	[ ] watch application folders and regenerate the swagger document
    - [ ] create fluent builder api
- Middlewares:
	- [ ] swagger editor
	- [ ] swagger UI
  - [ ] authorization:
		-	[ ] oauth2
			-	[ ] implicit
			-	[ ] access code
			-	[ ] password
			-	[ ] application
-	Transports:
	-	[ ] swagger socket (swagger over tcp sockets)
	-	[ ] swagger websocket (swagger over websockets)
	- [ ] swagger proxy (assemble several backend apis into a single swagger spec and route the requests)
	- [ ] swagger discovery (repository for swagger specifications)
