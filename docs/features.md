---
menu:
  - main
title: Features
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 30
---
## Full features list

- [x] An object model that serializes to swagger yaml or json (see the [spec package](https://github.com/go-openapi/spec))

- [x] A tool to work with swagger
  - [x] Serve swagger UI for any swagger spec file
  - [x] Validate a swagger spec document, with extra rules outlined [here](usage/validate.md)

  - [x] Generate API components based on swagger spec
    - [x] Flexible code generation, with customizable templates (package generator)
    - [x] Generate go client from a swagger spec
    - [x] Generate CLI (command line tool) client from a swagger spec
    - [x] Support swagger polymorphism (discriminator with allOf composition)


  - [x] Generate spec document based on annotated code (package `scanner`)
    - generate meta data (top level swagger properties) from package docs
    - generate definition entries for models
    - support composed structs out of several embeds
    - support allOf for composed structs
    - generate path entries for routes
    - generate responses from structs
    - support composed structs out of several embeds
    - generate parameters from structs
    - support composed structs out of several embeds

- [x] Middlewares (see: [runtime package](https://github.com/go-openapi/runtime))
  - [x] serve spec
  - [x] routing
  - [x] validation
  - [x] additional validation through an interface
  - [x] authorization, with auth composition (AND|OR authorization schemes)
    - basic auth
    - api key auth
    - oauth2 bearer auth
  - [x] swagger docs UI (docUI and redoc flavors)

- [x] Typed JSON Schema implementation
  - [x] JSON Pointer that knows about structs
  - [x] JSON Reference that knows about structs
  - [x] Supports most JSON schema features<sup>[1](#footnote1)</sup>
  - [x] Validate JSON data against jsonschema (Draft 4), with full $ref support (see the [validate package](https://github.com/go-openapi/validate))
    - passes current json schema test suite

- [x] extended string and numeric formats (see: [strfmt package](https://github.com/go-openapi/strfmt))
  - [x] JSON-schema draft 4 formats
    - date-time
    - email
    - hostname
    - ipv4
    - ipv6
    - uri
  - [x] swagger 2.0 format extensions
    - binary
    - byte (e.g. base64 encoded string)
    - date (e.g. "1970-01-01")
    - password
  - [x] go-openapi custom format extensions
    - bsonobjectid (BSON objectID)
    - creditcard
    - duration (e.g. "3 weeks", "1ms")
    - hexcolor (e.g. "#FFFFFF")
    - isbn, isbn10, isbn13
    - mac (e.g "01:02:03:04:05:06")
    - rgbcolor (e.g. "rgb(100,100,100)")
    - ssn
    - uuid, uuid3, uuid4, uuid5

- [x] Plays nice with golint, go vet etc.

<a name="footnote1">1</a>: currently adds extra support for `additionalItems`(not part of swagger), but not `anyOf`, `oneOf` and `not`.
