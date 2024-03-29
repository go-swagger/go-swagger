---
title: TODO List
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Todo List Tutorial

This example walks you through a hypothetical project, building a todo list.

It uses a todo list because this is well-understood application, so you can focus on the go-swagger pieces. In this example we build a server and a client.

<!--more-->

To create your application start with `swagger init`:

```sh
swagger init spec \
  --title "A Todo list application" \
  --description "From the todo list tutorial on goswagger.io" \
  --version 1.0.0 \
  --scheme http \
  --consumes application/io.goswagger.examples.todo-list.v1+json \
  --produces application/io.goswagger.examples.todo-list.v1+json
```

This gives you a skeleton `swagger.yml` file:

```yaml
definitions:
  item:
    type: object
    required:
      - description
    properties:
      id:
        type: integer
        format: int64
        readOnly: true
      description:
        type: string
        minLength: 1
      completed:
        type: boolean
```

In this model definition we say that the model `item` is an _object_ with a required property `description`. This item model has 3 properties: `id`, `description`, and `completed`. The `id` property is an int64 value and is marked as _readOnly_, meaning that it will be provided by the API server and it will be ignored when the item is created.

This document also says that the description must be at least 1 char long, which results in a string property that's [not a pointer](../reference/models/schemas.md#nullability).

At this moment you have enough so that actual code could be generated, but let's continue defining the rest of the API so that the code generation will be more useful. Now that you have a model so you can add some endpoints to list the todo's:

```yaml
paths:
  /:
    get:
      tags:
        - todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

With this new version of the operation you now have query params. These parameters have defaults so users can leave them off and the API will still function as intended.

However, this definition is extremely optimistic and only defines a response for the "happy path". It's very likely that the API will need to return errors too. That means you have to define a model errors, as well as at least one more response definition to cover the error response.

The error definition looks like this:

```yaml
paths:
  /:
    get:
      tags:
        - todos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
```

At this point you've defined your first endpoint completely. To improve the strength of this contract you could define responses for each of the status codes and perhaps return different error messages for different statuses. For now, the status code will be provided in the error message.

Try validating the specification again with `swagger validate ./swagger.yml` to ensure that code generation will work as expected. Generating code from an invalid specification leads to unpredictable results.

Your completed spec should look like this:

```yaml
paths:
  /:
    get:
      tags:
        - todos
      operationId: find_todos
      ...
```

These `operationId` values are used to name the generated files:

```
.
├── cmd
│   └── todo-list-server
│       └── main.go
├── models
│   ├── error.go
│   ├── find_todos_okbody.go
│   ├── get_okbody.go
│   └── item.go
├── restapi
│   ├── configure_todo_list.go
│   ├── doc.go
│   ├── embedded_spec.go
│   ├── operations
│   │   ├── todo_list_api.go
│   │   └── todos
│   │       ├── find_todos.go
│   │       ├── find_todos_parameters.go
│   │       ├── find_todos_responses.go
│   │       └── find_todos_urlbuilder.go
│   └── server.go
└── swagger.yml
```

You can see that the files under `restapi/operations/todos` now use the `operationId` as part of the generated file names.

At this point can start the server, but first let's see what `--help` gives you. First install the server binary and then run it:

```sh
± ~/go/src/.../examples/tutorials/todo-list/server-1
» go install ./cmd/todo-list-server/
± ~/go/src/.../examples/tutorials/todo-list/server-1
» todo-list-server --help
Usage:
  todo-list-server [OPTIONS]

From the todo list tutorial on goswagger.io

Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the
                           request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/todo-list.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
                           (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure connections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```

If you run your application now it will start on a random port by default. This might not be what you want, so you can configure a port through a command line argument or a `PORT` env var.

```sh
git:(master) ✗ !? » todo-list-server
serving todo list at http://127.0.0.1:64637
```

You can use `curl` to check your API:

```sh
git:(master) ✗ !? » curl -i http://127.0.0.1:64637/
```
```http
HTTP/1.1 501 Not Implemented
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Thu, 31 Dec 2015 22:42:10 GMT
Content-Length: 57

"operation todos.FindTodos has not yet been implemented"
```

As you can see, the generated API isn't very usable yet, but we know it runs and does something. To make it useful you'll need to implement the actual logic behind those endpoints. And you'll also want to add some more endpoints, like adding a new todo item and updating an existing item to change its description or mark it completed.

To supporting  adding a todo item you should define a `POST` operation:

```yaml
paths:
  /{id}:
    delete:
      tags:
        - todos
      operationId: destroyOne
      parameters:
        - type: integer
          format: int64
          name: id
          in: path
          required: true
      responses:
        204:
          description: Deleted
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
```

This time you're defining a parameter that is part of the `path`. This operation will look in the URI templated path for an id. Since there's nothing to return after a delete, the success response  is `204 No Content`.

Finally, you need to define a way to update an existing item:

```yaml
swagger: "2.0"
info:
  description: From the todo list tutorial on goswagger.io
  title: A Todo list application
  version: 1.0.0
consumes:
- application/io.goswagger.examples.todo-list.v1+json
produces:
- application/io.goswagger.examples.todo-list.v1+json
schemes:
- http
- https
paths:
  /:
    get:
      tags:
        - todos
      operationId: findTodos
      parameters:
        - name: since
          in: query
          type: integer
          format: int64
        - name: limit
          in: query
          type: integer
          format: int32
          default: 20
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
    post:
      tags:
        - todos
      operationId: addOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        201:
          description: Created
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
  /{id}:
    parameters:
      - type: integer
        format: int64
        name: id
        in: path
        required: true
    put:
      tags:
        - todos
      operationId: updateOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        200:
          description: OK
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
    delete:
      tags:
        - todos
      operationId: destroyOne
      responses:
        204:
          description: Deleted
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
definitions:
  item:
    type: object
    required:
      - description
    properties:
      id:
        type: integer
        format: int64
        readOnly: true
      description:
        type: string
        minLength: 1
      completed:
        type: boolean
  error:
    type: object
    required:
      - message
    properties:
      code:
        type: integer
        format: int64
      message:
        type: string
```

This is a good time to sanity check and by validating the schema:

```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-2
git:(master) ✗ !? » swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

Now you're ready to generate the API and start filling in the actual operations:

```sh
git:(master) ✗ !? » swagger generate server -A TodoList -f ./swagger.yml
... elided output ...
2015/12/31 18:16:28 rendered main template: server.TodoList
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-2
git:(master) ✗ !? » tree
.
├── cmd
│   └── todo-list-server
│       └── main.go
├── models
│   ├── error.go
│   ├── find_todos_okbody.go
│   └── item.go
├── restapi
│   ├── configure_todo_list.go
│   ├── doc.go
│   ├── embedded_spec.go
│   ├── operations
│   │   ├── todo_list_api.go
│   │   └── todos
│   │       ├── add_one.go
│   │       ├── add_one_parameters.go
│   │       ├── add_one_responses.go
│   │       ├── add_one_urlbuilder.go
│   │       ├── destroy_one.go
│   │       ├── destroy_one_parameters.go
│   │       ├── destroy_one_responses.go
│   │       ├── destroy_one_urlbuilder.go
│   │       ├── find_todos.go
│   │       ├── find_todos_parameters.go
│   │       ├── find_todos_responses.go
│   │       ├── find_todos_urlbuilder.go
│   │       ├── update_one.go
│   │       ├── update_one_parameters.go
│   │       ├── update_one_responses.go
│   │       └── update_one_urlbuilder.go
│   └── server.go
└── swagger.yml

6 directories, 26 files
```

To implement the core of your application you start by editing `restapi/configure_todo_list.go`. This file is safe to edit. Its content will not be overwritten if you run `swagger generate` again the future.

The simplest way to implement this application is to simply store all the todo items in a golang `map`. This provides a simple way to move forward without bringing in complications like a database or files.

To do this you'll need a map and a counter to track the last assigned id:

```go
// the variables we need throughout our implementation
var items = make(map[int64]*models.Item)
var lastID int64
```

The simplest handler to implement now is the delete handler. Because the store is a map and the id of the item is provided in the request it's a one liner.

```go
api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
  delete(items, params.ID)
  return todos.NewDestroyOneNoContent()
})
```

After deleting the item from the store, you need to provide a response. The code generator created responders for each response you defined in the swagger specification, and you can see how one of those is being used in the example above.

The other 3 handler implementations are similar to this one. They are provided in the [source for this tutorial](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/server-complete/restapi/configure_todo_list.go).

So assuming you go ahead and implement the remainder of the endpoints, you're all set to test it out:

```sh
» curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:01 GMT
Content-Length: 3

[]
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
```
```http
» curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}"
HTTP/1.1 415 Unsupported Media Type
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:11 GMT
Content-Length: 157

{"code":415,"message":"unsupported media type \"application/x-www-form-urlencoded\", only [application/io.goswagger.examples.todo-list.v1+json] are allowed"}                                                                                                     ± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
```
```sh
» curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:20 GMT
Content-Length: 39

{"description":"message 30925","id":1}
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:23 GMT
Content-Length: 37

{"description":"message 104","id":2}
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:24 GMT
Content-Length: 39

{"description":"message 15225","id":3}
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:26 GMT
Content-Length: 117

[{"description":"message 30925","id":1},{"description":"message 104","id":2},{"description":"message 15225","id":3}]
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765/3 -X PUT -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' -d '{"description":"go shopping"}'
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:32 GMT
Content-Length: 37

{"description":"go shopping","id":3}
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:34 GMT
Content-Length: 115

[{"description":"message 30925","id":1},{"description":"message 104","id":2},{"description":"go shopping","id":3}]
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765/1 -X DELETE -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 204 No Content
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:57:04 GMT
```
```sh
± ~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
» curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:57:06 GMT
Content-Length: 76

[{"description":"message 104","id":2},{"description":"go shopping","id":3}]
```
