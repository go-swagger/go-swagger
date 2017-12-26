## Full features list

For a V1 I want to have this feature set completed:

- [x] An object model that serializes to swagger yaml or json
- [x] A tool to work with swagger

  - [x] Serve swagger UI for any swagger spec file
  - [x] Flexible code generation, with customizable templates
  - [x] Generate API based on swagger spec
  - [x] Generate go client from a swagger spec
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
    - [ ] every example that is specified must validate against the schema for that property (Error)
    - [ ] $ref should not have siblings (Warning)
    - [x] each referable definition must have references (Warning)
    
  - [x] Validate against jsonschema (Draft 4), with full $ref support
  - [x] Generate spec document based on annotated code

    - [x] generate meta data (top level swagger properties) from package docs
    - [x] generate definition entries for models
    - [x] support composed structs out of several embeds
    - [x] support allOf for composed structs
    - [x] generate path entries for routes
    - [x] generate responses from structs
    - [x] support composed structs out of several embeds
    - [x] generate parameters from structs
    - [x] support composed structs out of several embeds

- [x] Middlewares

  - [x] serve spec
  - [x] routing
  - [x] validation
  - [x] additional validation through an interface
  - [x] authorization

    - [x] basic auth
    - [x] api key auth
    - [x] oauth2 bearer auth

  - [x] swagger docs UI

- [x] Typed JSON Schema implementation

  - [x] JSON Pointer that knows about structs
  - [x] JSON Reference that knows about structs
  - [x] Passes current json schema test suite

- [x] extended string and numeric formats

  - [x] uuid, uuid3, uuid4, uuid5
  - [x] email
  - [x] uri (absolute)
  - [x] hostname
  - [x] ipv4
  - [x] ipv6
  - [x] mac
  - [x] credit card
  - [x] isbn, isbn10, isbn13
  - [x] social security number
  - [x] hexcolor
  - [x] rgbcolor
  - [x] date
  - [x] date-time
  - [x] duration
  - [x] password
  - [x] custom string formats
  - [x] int32, int64, float, double

- [x] Project documentation site
- [x] Play nice with golint, go vet etc.
