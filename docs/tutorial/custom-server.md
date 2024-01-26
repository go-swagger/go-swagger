---
title: Custom server
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Custom Server Tutorial

In this tutorial we'll be building up a custom server.
The core will be generated using a manually written and maintained OpenAPI 2.0 spec.
The cli code will be a thin layer around that, and will simply setup the API and server,
using the parsed configurations and our hand written handlers.

<!--more-->

The server we'll building will be very simple.
In this tutorial we'll assume you are already familiar with defining an API,
using the OpenAPI 2.0 yaml specification format.
Please consult [the official OpenAPI 2.0 specification][OpenAPI2.0]
for more information in case you're new to OpenAPI (Also known as Swagger).

The end product of this tutorial can be found as `./examples/tutorials/custom-server`.

The server we'll be building, will be generated using the following spec:

```yaml
---
swagger: '2.0'
info:
  version: 1.0.0
  title: Greeting Server
paths:
  /hello:
    get:
      produces:
        - text/plain
      parameters:
        - name: name
          required: false
          type: string
          in: query
          description: defaults to World if not given
      operationId: getGreeting
      responses:
        200:
          description: returns a greeting
          schema:
              type: string
              description: contains the actual greeting as plain text
```
