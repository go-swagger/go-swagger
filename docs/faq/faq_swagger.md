---
title: About using swagger as your API contract
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about swagger specs -->
## Swagger specification

### Default vs required
_Use-Case_:  spec validation rejects default object without required properties

> If we have the default for the firstName and lastName but didn't set the username, which is required,
> we will have an error message showing that we have to put the default for the required field "username" as well.
>
> The "default" section must contain any "required" properties, which is counterintuitive because the "required" properties
> should always be provided explicitly and shouldn’t have defaults.

Example:

```yaml
paths:
  /person/{id}:
    put:
      summary: Put a person
      description: Put a person.
      consumes:
        - application/json
      parameters:
        - name: id
          in: path
          description: ID of person
          required: true
          type: integer
        - in: body
          name: body
          description: Person
          schema:
            type: object
            required:
              - username
            default:
              firstName: John
              lastName: Smith
            properties:
              firstName:
                type: string
              lastName:
                type: string
              username:
                type: string
      responses:
        200:
          description: successful
```

`swagger generate model -f person.yaml`

I notice that the error message could be prevented if I use `swagger generate model --skip-validation`.
But I think the error message should not show even I just use `swagger generate model -f person.yaml` because
"required" properties should always be provided explicitly and shouldn’t have defaults.

**Answer:** Those defaults need to be put on the individual members.

> The default where you put it is a default for the entire object if it is not provided at all.
> Thus it must pass validation of the object schema.
> As per swagger, defaults MUST validate their schema. This differs from json-Schema spec.

Originally from issue [#1552](https://github.com/go-swagger/go-swagger/issues/1552).

_Use-Case_: `go-swagger` rejects default object that misses required properties.

> My understanding is that properties in default and those listed as required are mutually exclusive, because one is supposed
> to explicitly provide required properties, while properties in default can be filled in.
> See https://swagger.io/blog/unlocking-the-spec-the-default-keyword

Example:

```yaml
paths:
  /person/{id}:
    put:
      summary: Put a person
      description: Put a person.
      consumes:
        - application/json
      parameters:
        - name: id
          in: path
          description: ID of person
          required: true
          type: integer
        - in: body
          name: body
          description: Person
          schema:
            type: object
            required:
              - username
            default:
              firstName: John
              lastName: Smith
            properties:
              firstName:
                type: string
              lastName:
                type: string
              username:
                type: string
      responses:
        200:
          description: successful
```

**Answer:** the spec validation warns you about mixing up default with required.

> If you override this validation (with `swagger generate model --skip-validation`), the model generator would override the
> required clause, and use the default.

**Would --skip-validation skip any other validations?**

**Answer:** Yes. `--skip-validation` skips the call to `swagger validate {spec}`, which is normally carried on before generation.

For some use cases, we have to, since `go-swagger` supports constructs that are not stricly swagger-compliant.
This one is an example: some override silently take place here at generation time.

Originally from issue [#1501](https://github.com/go-swagger/go-swagger/issues/1501).

### type string, format int64 not respected in generator
_Use-Case_:  when generating parameters or models from a swagger file with a definition that specifies type: string and format: int64,
 the **generation fails**.

Example:

```json
{
    "consumes": [
        "application/json"
    ],
    "definitions": {
        "Test": {
            "properties": {
                "id": {
                    "format": "int64",
                    "type": "string"
                }
            },
            "type": "object"
        }
    },
...
```

**Answer**: that's an invalid type. The openapi 2.0 spec says the type should be
```
type: integer
format: int64
```


Originally from issue [#1381](https://github.com/go-swagger/go-swagger/issues/1381).

### Duplicate operationId error
_Use-Case_:  my spec indicates duplicate operationIds but for separate endpoints.
When generating a server and there are multiple endpoints, the `operations` directory in `restapi`
contains subdirectories for each of those endpoints.

Example:

```yaml
paths:
  '/users/{id}':
    parameters:
      - name: id
        in: path
        required: true
        type: string
    get:
      operationId: getById
      summary: Get User By ID
      tags:
        - Users
      responses:
        '200':
          description: ''
          schema:
            $ref: '#/definitions/user-output'
  '/pets/{id}':
    parameters:
      - name: id
        in: path
        required: true
        type: string
    get:
      operationId: getById
      summary: Get Pet By ID
      tags:
        - Pets
      responses:
        '200':
          description: ''
          schema:
            $ref: '#/definitions/pet-output'
```

I am expecting a generated structure like so:
```
└── restapi
    ...
    ├── operations
    │   ├── pets
    │   │   ├── ...
    │   ├── users
    │   │   ├── ...
    │   └── ...
```

However, instead of generating the server I get this error: `"getById" is defined 2 times`

Is it possible to bypass this validation error IF the operation ids will go into separate directories?

**Answer:** operationId is always globally unique per swagger 2.0 specification. The id MUST be unique among all operations described in the API.

https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#operationObject

> Unique string used to identify the operation. The id MUST be unique among all operations described in the API.
> Tools and libraries MAY use the operationId to uniquely identify an operation,
> therefore, it is recommended to follow common programming naming conventions.

Originally from issue [#1143](https://github.com/go-swagger/go-swagger/issues/1143).

