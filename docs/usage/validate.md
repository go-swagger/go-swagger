---
title: swagger validate
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 70
---
# Validate a swagger spec

The toolkit has a command to validate swagger specifications for you.
It includes a full json-schema validator and adds some extra validations to ensure the spec is valid.

### Usage

To validate a specification:

```
Usage:
  swagger [OPTIONS] validate [validate-OPTIONS]

validate the provided swagger document against a swagger spec

Application Options:
  -q, --quiet                 silence logs
      --output=LOG-FILE       redirect logs to file

Help Options:
  -h, --help                  Show this help message

[validate command options]
          --skip-warnings     when present will not show up warnings upon validation
          --stop-on-error     when present will not continue validation after critical errors are found
```

### Swagger 2.0 resources

* Specification Documentation: https://github.com/swagger-api/swagger-spec/blob/master/versions/2.0.md
* JSON Schema: https://github.com/swagger-api/swagger-spec/blob/master/schemas/v2.0/schema.json

### Semantic Validation

All the rules the validator tool supports:

*	validate against jsonschema
*	validate extra rules, inspired from [the sway swagger validator](https://github.com/apigee-127/sway/tree/master/docs#semantic-validation)

Rule | Severity
