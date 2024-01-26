---
title: route
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:route

A **swagger:route** annotation links a path to a method.
This operation gets a unique id, which is used in various places as method name.
One such usage is in method names for client generation for example.

Because there are many routers available, this tool does not try to parse the paths
you provided to your routing library of choice. So you have to specify your path pattern
yourself in valid swagger syntax.

<!--more-->

##### Syntax:

```go
swagger:route [method] [path pattern] [?tag1 tag2 tag3] [operation id]
```

##### Properties

Annotation | Format
```
paths:
  "/pets":
    get:
      operationId: listPets
      deprecated: true
      summary: Lists pets filtered by some parameters.
      description: "This will show all available pets by default.\nYou can get the pets that are out of stock"
      tags:
      - pets
      - users
      consumes:
      - application/json
      - application/x-protobuf
      produces:
      - application/json
      - application/x-protobuf
      schemes:
      - http
      - https
      - ws
      - wss
      security:
        api_key: []
        oauth:
        - read
        - write
      parameters:
        description: maximum number of results to return
        format: int43
        in: query
        name: limit
        type: integer
      responses:
        default:
          $ref: "#/responses/genericError"
        200:
          $ref: "#/responses/someResponse"
        422:
          $ref: "#/responses/validationError"
      extensions:
        x-example-flag: true
        x-some-list:
        - dog
        - cat
        - bird
```
