## Full features list

- [x] An object model that serializes to swagger yaml or json (see: [spec package](https://github.com/go-openapi/spec))
- [x] A tool to work with swagger

  - [x] Serve swagger UI for any swagger spec file
  - [x] Flexible code generation, with customizable templates (package generator)
  - [x] Generate API based on swagger spec
  - [x] Generate go client from a swagger spec
  - [x] Support swagger polymorphism (discriminator with allOf composition)
  - [x] Validate a swagger spec document, with extra rules outlined [here](https://github.com/apigee-127/sway/blob/master/docs/README.md#semantic-validation)

    - [x] definition can't declare a property that's already defined by one of its ancestors (Error)
    - [x] definition's ancestor can't be a descendant of the same model (Error)
    - [x] each api path should be non-verbatim (account for path param names) unique per method (Error)
    - [x] each path parameter should correspond to a parameter placeholder and vice versa (Error)
    - [x] path parameter declarations do not allow empty names _(`/path/{}` is not valid)_ (Error)
    - [x] each definition property listed in the required array must be defined in the properties of the model (Error)
    - [x] each parameter should have a unique `name` and `in` combination (Error)
    - [x] each operation should have at most 1 parameter of type body (Error)
    - [x] each operation cannot have both a body parameter and a formData parameter (Error)
    - [x] each operation must have an unique `operationId` (Error)
    - [x] each reference must point to a valid object (Error)
    - [x] every default value that is specified must validate against the schema for that property (Error)
    - [x] items property is required for all schemas/definitions of type `array` (Error)
    - |x] param in path must have the property required: true
    - [x] every example that is specified should validate against the schema for that property (Warning)
    - [x] $ref should not have siblings (Warning)
    - [x] each referable definition must have references (Warning)
    
  - [x] Validate JSON data against jsonschema (Draft 4), with full $ref support (see: [validate package](https://github.com/go-openapi/validate))
    - [x] Passes current json schema test suite

  - [x] Generate spec document based on annotated code (package scan)

    - [x] generate meta data (top level swagger properties) from package docs
    - [x] generate definition entries for models
    - [x] support composed structs out of several embeds
    - [x] support allOf for composed structs
    - [x] generate path entries for routes
    - [x] generate responses from structs
    - [x] support composed structs out of several embeds
    - [x] generate parameters from structs
    - [x] support composed structs out of several embeds

- [x] Middlewares (see: [runtime package](https://github.com/go-openapi/runtime))

  - [x] serve spec
  - [x] routing
  - [x] validation
  - [x] additional validation through an interface
  - [x] authorization, with auth composition (AND|OR authorization schemes)

    - [x] basic auth
    - [x] api key auth
    - [x] oauth2 bearer auth

  - [x] swagger docs UI (docUI and redoc flavors)

- [x] Typed JSON Schema implementation

  - [x] JSON Pointer that knows about structs
  - [x] JSON Reference that knows about structs
  - [x] Supports most JSON schema features<sup>[1](#footnote1)</sup>

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

- [x] Play nice with golint, go vet etc.

<a name="footnote1">1</a>: currently adds extra support for `additionalItems`(not part of swagger), but not `anyOf`, `oneOf` and `not`.
