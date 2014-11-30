# Swagger 2.0

[![Build Status](https://travis-ci.org/casualjim/go-swagger.svg?branch=master)](https://travis-ci.org/casualjim/go-swagger)
[![Coverage Status](https://img.shields.io/coveralls/casualjim/go-swagger.svg)](https://coveralls.io/r/casualjim/go-swagger)
[![GoDoc](https://godoc.org/github.com/casualjim/go-swagger?status.svg)](http://godoc.org/github.com/casualjim/go-swagger)
[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/swagger-api/swagger-spec/master/LICENSE)

Contains an implementation of Swagger 2.0.
It knows how to serialize and deserialize swagger specifications.

At present it's got the entire object model defined, and I'm writing tests to make that work completely.

Swagger is a simple yet powerful representation of your RESTful API.  
With the largest ecosystem of API tooling on the planet, thousands of developers are supporting Swagger
in almost every modern programming language and deployment environment.   

With a Swagger-enabled API, you get interactive documentation, client SDK generation and discoverability.
We created Swagger to help fulfill the promise of APIs.   

Swagger helps companies like Apigee, Getty Images, Intuit, LivingSocial, McKesson, Microsoft, Morningstar, and PayPal 
build the best possible services with RESTful APIs.Now in version 2.0, Swagger is more enabling than ever. 
And it's 100% open source software.

## What's inside?

* An object model that serializes to swagger yaml or json 
* A tool to work with swagger:
    * validate a swagger spec document

## Planned:
* Generate validations based on the swagger spec
* Later it will also know how to generate those specifications from your source code.
* Generate a stub api based on a swagger spec
* Generate a client from a swagger spec
* Build a full swagger spec by inspecting your source code and embedding it in a go file.
