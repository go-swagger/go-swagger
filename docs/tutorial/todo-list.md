# Todo List Tutorial

This example walks you through a hypothetical project that is building a todo list.
It specifically uses a todo list because it's a super well-understood application and hopefully this allows you to focus entirely on the new concepts. This example builds a server and then a client.

<!--more-->

When you start an application most likely you think about the functionality it supports.

```shell
swagger init spec \
  --title "A To Do list application" \
  --description "The product of a tutorial on goswagger.io" \
  --version 1.0.0 \
  --scheme http \
  --consumes application/io.goswagger.examples.todo-list.v1+json \
  --produces application/io.goswagger.examples.todo-list.v1+json
```

You can get started with a swagger.yml like this:

```yaml
---
consumes:
- application/io.goswagger.examples.todo-list.v1+json
definitions: {}
info:
  description: The product of a tutorial on goswagger.io
  title: A To Do list application
  version: 1.0.0
paths: {}
produces:
- application/io.goswagger.examples.todo-list.v1+json
schemes:
- http
swagger: "2.0"
```

This doesn't do much but it would validate in the swagger validator step.

```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list
git:(master) ✗ ? » swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

So now you have an empty but valid specification document, time to get to declaring some models and endpoints for the API. You'll probably need a model to represent a todo item, you can define that in the definitions.

```yaml
---
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

In this model definition we say that the model `item` is an _object_ with a required property `description`. This item model has 3 properties: id, description and completed. The `id` property is an int64 value and is marked as _readOnly_, so that means that it will be provided by the API server and it will be ignored when the item is created.
This document also says that the description must be at least 1 char long, this will result in a string property that's [not a pointer](/use/schemas.md#nullability).

At this moment there is enough to get some actual code generated, but let's wait with that and continue defining the rest of the API so that the code generation later on will be more useful. Now you have a model so you probably want to add some endpoints to list the todo's.

```yaml
---
paths:
  /:
    get:
      tags:
        - todos
      responses:
        200:
          description: list the todo operations
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

This snippet of yaml defines a `GET /` operation, and tags it with _todos_. Tagging things is nice because tools do all kinds of fancy things with tags. Tags help UI's group endpoints appropriately, code generators might turn them into 'controllers'. Furthermore there is a response defined with a generic description, about what's in the response. Be aware that some generators think a field like that is a good thing to put in the http status message. And then of course the response defines also the return type of that endpoint. In this case the endpoint will be returning a list of todo items, so the schema is an _array_ and the array will contain items that look like the item definition you declared earlier.

But wait a minute, what if there are 100's of todo items, will we just return all of them for everybody?  It might be best to add a since and limit param here. The ids will have ordered for a since param to work but you're in control of that so that's fine.

```yaml
---
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

With this new version of the operation YAML, there are query params now for the values and they define defaults so people can leave them off and the API will still function as intended.

However  this definition is extremely optimistic and only defines a response for the "happy path". It's very likely that the API will need to return some form of error messages too. So that means you probably have to define a model for the error messages as well as at least one more response definition to cover both the bodies in you contract.

The error definition might look like this:

```yaml
---
definitions:
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

For the extra response you can use the default response, because after all every successful response from your API is defying the odds.

```yaml
---
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

At this point you've got your first endpoint defined completely. To improve the strength of this contract you could define responses for each of the status codes and perhaps return a different error message. In this case the status code will be provided in the error message, and can easily be different from the HTTP status codes, who typically only give you a hint of what went wrong.

Perhaps validate the specification again, having a valid swagger document, is important when using the code generation, there are quite a few factors that contribute to rendering the models for a specification. An invalid swagger document makes it so that the generated code will have unpredictable behavior.

So the completed spec should look like this:

```yaml
---
swagger: "2.0"
info:
  description: The product of a tutorial on goswagger.io
  title: A To Do list application
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

Once you generate a server for this you'll see the following directory listing:

```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-1
git:(master) ✗ !? » swagger generate server -A TodoList -f ./swagger.yml
2016/02/15 12:32:40 building a plan for generation
2016/02/15 12:32:40 planning definitions
2016/02/15 12:32:40 planning operations
2016/02/15 12:32:40 grouping operations into packages
2016/02/15 12:32:40 planning meta data and facades
2016/02/15 12:32:40 rendering 2 models
2016/02/15 12:32:40 rendered model template: error
2016/02/15 12:32:40 rendered model template: item
2016/02/15 12:32:40 rendered handler template: todos.Get
2016/02/15 12:32:40 generated handler todos.Get
2016/02/15 12:32:40 rendered parameters template: todos.GetParameters
2016/02/15 12:32:40 generated parameters todos.GetParameters
2016/02/15 12:32:40 rendered responses template: todos.GetResponses
2016/02/15 12:32:40 generated responses todos.GetResponses
2016/02/15 12:32:40 rendered builder template: operations.TodoList
2016/02/15 12:32:40 rendered embedded Swagger JSON template: server.TodoList
2016/02/15 12:32:40 rendered configure api template: operations.ConfigureTodoList
2016/02/15 12:32:40 rendered doc template: operations.TodoList
2016/02/15 12:32:40 rendered main template: server.TodoList

± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-1
git:(master) ✗ !? » tree
.
├── cmd
│   └── todo-list-server
│       ├── configure_todo_list.go
│       ├── doc.go
│       ├── embedded_spec.go
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   └── operations
│       ├── todo_list_api.go
│       └── todos
│           ├── get.go
│           ├── get_parameters.go
│           └── get_responses.go
└── swagger.yml

6 directories, 11 files
```

In this file tree you notice that there is a cmd/todo-list-server generated. The swagger generator adds -server to the application name (provided to the generated command through the -A argument).

The second major section in this tree is the models package. This package contains go representations for both the defintions from the swagger spec document.

And then the last major section is the rest api, within the rest api there is the code that is generated based on the information from the paths property in the swagger specification. The go swagger generator uses the tags to group the operations into packages.

We skipped over naming operations, you have the ability to name the operations by giving operations an ID in the specification document. For example for the operation defintion with `operationId: findTodos`, the following tree would be generated:

```shellsession
.
├── cmd
│   └── todo-list-server
│       ├── configure_todo_list.go
│       ├── doc.go
│       ├── embedded_spec.go
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   └── operations
│       ├── todo_list_api.go
│       └── todos
│           ├── find_todos.go
│           ├── find_todos_parameters.go
│           └── find_todos_responses.go
└── swagger.yml
```

At this point you're able to start the server, but lets first see what --help gives you. To make this happen you can first install the binary and then run it.

```shellsession
± ivan@aether:~/go/src/.../examples/tutorials/todo-list/server-1
git:(master) ✗ -!? » go install ./cmd/todo-list-server/
± ivan@aether:~/go/src/.../examples/tutorials/todo-list/server-1
git:(master) ✗ -!? » todo-list-server --help
Usage:
  todo-list-server [OPTIONS]

The product of a tutorial on goswagger.io

Application Options:
      --host= the IP to listen on (default: localhost) [$HOST]
      --port= the port to listen on for insecure connections, defaults to a random value [$PORT]

Help Options:
  -h, --help  Show this help message
```

As you can see your application can be run straightaway and it will use a random port value by default. This might not be what you want so you can configure a port with an argument or through an environment variable. Those env vars are chosen because many platforms (like heroku) use those to configure the apps they are running.

```shellsession
git:(master) ✗ !? » todo-list-server
serving todo list at http://127.0.0.1:64637
```

And you can curl it:

```shellsession
git:(master) ✗ !? » curl -i http://127.0.0.1:64637/
```
```http
HTTP/1.1 501 Not Implemented
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Thu, 31 Dec 2015 22:42:10 GMT
Content-Length: 57

"operation todos.FindTodos has not yet been implemented"
```

So immediately after generating, the API has limited usability, but this can serve as a sanity check. Your API will need more endpoints besides listing todo items. Because this tutorial highlights design first, we know we will need a few other things. In addition to listing the todo items the API will need a means to add a new todo item, update the description of one and mark a todo item as completed.

For adding a todo item you probably want to define a POST operation, for our purposes that might look like this:

```yaml
---
paths:
  /:
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
```

So in this YAML snippet there is one new thing: you're defining that your API has a POST body and that that should be the item model defined earlier. Earlier the item schema was defined with a readOnly id, so that means it doesn't need to be included in the POST body. But the response to the POST request will include an id property. The next operation to define is the DELETE operation, where you delete a todo item from the list.

```yaml
---
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

The new concept in the DELETE method is that this time around you're defining a path parameter. When you delete an item you need to provide an ID. This is typically done in the path of a resource, and that's how this operation will know which todo item to delete. Once you've deleted something you can't really return its data anymore, so the success response in this case is `204 No Content`. At this point all you still require is a means to update an item, which combines everything you've just learnt.

```yaml
---
paths:
  /{id}:
    parameters:
      - type: integer
        format: int64
        name: id
        in: path
        required: true
    put:
      tags: ["todos"]
      operationId: updateOne
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/item"
      responses:
        '200':
          description: OK
          schema:
            $ref: "#/definitions/item"
        default:
          description: error
          schema:
            $ref: "#/definitions/error"
    delete:
      # elided for brevity
```

For updates there are 2 approaches that people typically take: PUT indicates replace the entity, PATCH indicates only update the fields provided in the request. In this case the "brute force" approach of replacing an entire entity is taken. Another thing you see here is that because the id path parameter is shared between both the put and the delete method, there is an opportunity to DRY the operation definition up a bit. So the path parameter moved to the common parameters section for the path.

At this point you should have a completed specification for the todo list API.

```yaml
---
swagger: "2.0"
info:
  description: The product of a tutorial on goswagger.io
  title: A To Do list application
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

Again this is a good time to sanity check, and run the validator.

```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-2
git:(master) ✗ !? » swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

You're ready to generate the API and start filling out some of the blanks.

```shellsession
git:(master) ✗ !? » swagger generate server -A TodoList -f ./swagger.yml
... elided output ...
2015/12/31 18:16:28 rendered main template: server.TodoList
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-2
git:(master) ✗ !? » tree
.
├── cmd
│   └── todo-list-server
│       ├── configure_todo_list.go
│       ├── doc.go
│       ├── embedded_spec.go
│       └── main.go
├── models
│   ├── error.go
│   └── item.go
├── restapi
│   └── operations
│       ├── todo_list_api.go
│       └── todos
│           ├── add_one.go
│           ├── add_one_parameters.go
│           ├── add_one_responses.go
│           ├── destroy_one.go
│           ├── destroy_one_parameters.go
│           ├── destroy_one_responses.go
│           ├── find_todos.go
│           ├── find_todos_parameters.go
│           ├── find_todos_responses.go
│           ├── update_one.go
│           ├── update_one_parameters.go
│           └── update_one_responses.go
└── swagger.yml

6 directories, 20 files
```

When you want to add functionality to your application, like adding your own code to the generated code, the place to do that is in the `configure_todo_list.go` file. The configure_todo_list file is safe to edit, it won't get overwritten in regeneration passes, to override the file you'd have to delete it yourself and regenerate.

A good first implementation of the todo list api, can be one where everything is stored in a map. This should show that everything works, without the complications of dealing with a database yet.
For this you'll need to track the state for an incrementing id and a structure to store items in.

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

After deleting the item from the store, there is a responder that needs to be created. The code generator has generated responders for each response you defined in the the swagger specification. The other 3 handler implementations are pretty similar to this one, they are provided in the [source for this tutorial](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/server-complete/restapi/configure_todo_list.go).

You're all set now, with a spiffy new todo list api implemented, lets see if it actually works.

```shellsession
git:(master) ✗ -!? » curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:01 GMT
Content-Length: 3

[]
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
```
```http
git:(master) ✗ -!? » curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}"
HTTP/1.1 415 Unsupported Media Type
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:11 GMT
Content-Length: 157

{"code":415,"message":"unsupported media type \"application/x-www-form-urlencoded\", only [application/io.goswagger.examples.todo-list.v1+json] are allowed"}                                                                                                     ± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
```
```shellsession
git:(master) ✗ -!? » curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:20 GMT
Content-Length: 39

{"description":"message 30925","id":1}
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:23 GMT
Content-Length: 37

{"description":"message 104","id":2}
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765 -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:24 GMT
Content-Length: 39

{"description":"message 15225","id":3}
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:26 GMT
Content-Length: 117

[{"description":"message 30925","id":1},{"description":"message 104","id":2},{"description":"message 15225","id":3}]
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765/3 -X PUT -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' -d '{"description":"go shopping"}'
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:32 GMT
Content-Length: 37

{"description":"go shopping","id":3}
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:56:34 GMT
Content-Length: 115

[{"description":"message 30925","id":1},{"description":"message 104","id":2},{"description":"go shopping","id":3}]
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765/1 -X DELETE -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```
```http
HTTP/1.1 204 No Content
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:57:04 GMT
```
```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/server-complete
git:(master) ✗ -!? » curl -i localhost:8765
```
```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Fri, 01 Jan 2016 19:57:06 GMT
Content-Length: 76

[{"description":"message 104","id":2},{"description":"go shopping","id":3}]
```
