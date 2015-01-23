# Swagger 2.0

[![Build Status](https://travis-ci.org/casualjim/go-swagger.svg?branch=master)](https://travis-ci.org/casualjim/go-swagger)
[![Coverage Status](https://coveralls.io/repos/casualjim/go-swagger/badge.svg?branch=master)](https://coveralls.io/r/casualjim/go-swagger?branch=master)
[![GoDoc](https://godoc.org/github.com/casualjim/go-swagger?status.svg)](http://godoc.org/github.com/casualjim/go-swagger)
[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE)

Contains an implementation of Swagger 2.0.
It knows how to serialize and deserialize swagger specifications.

Swagger is a simple yet powerful representation of your RESTful API.  
With the largest ecosystem of API tooling on the planet, thousands of developers are supporting Swagger
in almost every modern programming language and deployment environment.   

With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability.
We created Swagger to help fulfill the promise of APIs.   

Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal 
build the best possible services with RESTful APIs.Now in version 2.0, Swagger is more enabling than ever. 
And it's 100% open source software.

## Docs

http://godoc.org/github.com/casualjim/go-swagger

Install:

  go get -u github.com/casualjim/go-swagger/cmd/swagger

The implementation also provides a number of command line tools to help working with swagger.

Currently there is a spec validator tool:

  swagger validate https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json

## What's inside?

- [x] An object model that serializes to swagger yaml or json 
- [x] A tool to work with swagger:
  - [x] validate a swagger spec document
  - [ ] generate stub api based on swagger spec
  - [ ] generate client from a swagger spec
  - [ ] generate spec document based on the code
  - [ ] generate "sensible" random data based on swagger spec
  - [ ] generate tests based on swagger spec for server
  - [ ] generate tests based on swagger spec for client
- [x] Middlewares:
  - [x] serve spec
  - [x] routing
  - [x] validation 
  - [ ] authorization
  - [ ] swagger docs UI
  - [ ] swagger editor UI
- [x] Typed JSON Schema implementation
- [x] extended string formats
  - uuid, uuid3, uuid4, uuid5
  - email
  - uri
  - hostname
  - ipv4
  - ipv6
  - credit card
  - isbn, isbn10, isbn13
  - social security number
  - hexcolor
  - rgbcolor


