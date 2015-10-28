+++
categories = ["vallidate"]
date = "2015-10-22T23:28:09-07:00"
tags = ["usage"]
title = "Validate a swagger spec"

+++

## Validate a swagger spec

The toolkit has a command to validate swagger specifications for you.
It includes a full json-schema validator and adds a bunch of extra validations on top of that to ensure the spec
is as valid as possible, so that there are no surprises.

### Usage

To validate a specification:

```
swagger validate [http-url|filepath]
```

### Supported rules

All the rules the validator tool supports:

-	validate against jsonschema
-	validate extra rules outlined [here](https://github.com/apigee-127/swagger-tools/blob/master/docs/Swagger_Validation.md)
  - definition can't declare a property that's already defined by one of its ancestors (Error)
  - definition's ancestor can't be a descendant of the same model (Error)
  - each api path should be non-verbatim (account for path param names) unique per method (Error)
  - each security reference should contain only unique scopes (Warning)
  - each security scope in a security definition should be unique (Warning)
  - each path parameter should correspond to a parameter placeholder and vice versa (Error)
  - each referencable definition must have references (Warning)
  - each definition property listed in the required array must be defined in the properties of the model (Error)
  - each parameter should have a unique `name` and `type` combination (Error)
  - each operation should have only 1 parameter of type body (Error)
  - each reference must point to a valid object (Error)
  - every default value that is specified must validate against the schema for that property (Error)
  - every example that is specified must validate against the schema for that property (Error)
  - items property is required for all schemas/definitions of type `array` (Error)
