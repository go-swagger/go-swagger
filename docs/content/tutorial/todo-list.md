+++
categories = ["tutorial"]
date = "2015-12-30T11:56:40-08:00"
tags = ["todo app", "client", "server"]
title = "Todo List Tutorial"
weight = 7
series = ["tutorials", "home"]
+++

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
  --scheme https \
  --consumes application/io.goswagger.examples.todo-list.v1+json \
  --produces application/io.goswagger.examples.todo-list.v1+json
```

You can get started with a swagger.yml like this:

```yaml
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
paths: {}
definitions: {}
```

This doesn't do much but it would validate in the swagger validator step.

```shellsession
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list
git:(master) ✗ ? » swagger validate ./swagger.yml
The swagger spec at "./swagger.yml" is valid against swagger specification 2.0
```

So now you have an empty but valid specification document, time to get to declaring some models and endpoints for the
API. You'll probably need a model to represent a todo item, you can define that in the definitions.

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

In this model definition we say that the model `item` is an _object_ with a required property `description`.  This item
model has 3 properties: id, description and completed. The `id` property is an int64 value and is marked as _readOnly_,
so that means that it will be provided by the API server and it will be ignored when the item is created.
This document also says that the description must be at least 1 char long, this will result in a string property that's
[not a pointer](http://goswagger.io/use/schemas/#nullability:176038017a790b96307b48b85dc07885).

At this moment there is enough to get some actual code generated, but let's wait with that and continue defining the
rest of the API so that the code generation later on will be more useful. Now you have a model so you probably want to
add some endpoints to list the todo's.

```yaml
paths:
  /:
    get:
      tags: ["todos"]
      responses:
        '200':
          description: "list the todo operations"
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

This snippet of yaml defines a `GET /` operation, and tags it with _todos_. Tagging things is nice because tools do all
kinds of fancy things with tags. Tags help UI's group endpoints appropriately, code generators might turn them into
'controllers'.  Furthermore there is a response defined with a generic description, about what's in the response.  Be
aware that some generators think a field like that is a good thing to put in the http status message.  And then of
course the response defines also the return type of that endpoint. In this case the endpoint will be returning a list of
todo items, so the schema is an _array_ and the array will contain items that look like the item definition you
declared earlier.

But wait a minute, what if there are 100's of todo items, will we just return all of them for everybody?  It might be
best to add a since and limit param here. The ids will have ordered for a since param to work but you're in control of
that so that's fine.

```yaml
paths:
  /:
    get:
      tags: ["todos"]
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
        '200':
          description: "list the todo operations"
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

With this new version of the operation yaml, there are query params now for the values and they define defaults so
people can leave them off and the API will still function as intented.

However  this definition is extremely optimistic and only defines a response for the "happy path". It's very likely
that the API will need to return some form of error messages too. So that means you probably have to define a model for
the error messages as well as at least one more response definition to cover both the bodies in you contract.

The error definition might look like this:

```yaml
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

For the extra response you can use the default response, because after all every successful response from your API is
defying the odds.

```yaml
paths:
  /:
    get:
      tags: ["todos"]
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
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
        '200':
          description: "list the todo operations"
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
```

At this point you've got your first endpoint defined completely. To improve the strength of this contract you could
define reponses for each of the status codes and perhaps return a different error message. In this case the status code
will be provided in the error message, and can easily be different from the HTTP status codes, who typically only give
you a hint of what went wrong.

Perhaps validate the specification again, having a valid swagger document, is important when using the code generation,
there are quite a few factors that contribute to rendering the models for a specification. An invalid swagger document
makes it so that the generated code will have unpredictable behavior.

So the completed spec should look like this:

```yaml
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
      tags: ["todos"]
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
        default:
          description: generic error response
          schema:
            $ref: "#/definitions/error"
        '200':
          description: "list the todo operations"
          schema:
            type: array
            items:
              $ref: "#/definitions/item"
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
± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/first-generation
git:(master) ✗ !? » swagger generate server -A TodoList -f ./swagger.yml
2015/12/30 19:06:54 rendered model template: error
2015/12/30 19:06:54 generated model error
2015/12/30 19:06:54 rendered model template: item
2015/12/30 19:06:54 generated model item
2015/12/30 19:06:54 rendered handler template: todos.Get
2015/12/30 19:06:54 generated handler todos.Get
2015/12/30 19:06:54 rendered parameters template: todos.GetParameters
2015/12/30 19:06:55 generated parameters todos.GetParameters
2015/12/30 19:06:55 rendered responses template: todos.GetResponses
2015/12/30 19:06:55 generated responses todos.GetResponses
2015/12/30 19:06:55 rendered builder template: operations.TodoList
2015/12/30 19:06:55 rendered embedded Swagger JSON template: server.TodoList
2015/12/30 19:06:55 skipped (already exists) configure api template: operations.ConfigureTodoList
2015/12/30 19:06:55 rendered doc template: operations.TodoList
2015/12/30 19:06:55 rendered main template: server.TodoList

± ivan@aether:~/go/src/github.com/go-swagger/go-swagger/examples/tutorials/todo-list/first-generation
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

In this file tree you notice that there is a cmd/todo-list-server generated. The swagger generator adds -server to the
application name (provided to the generated command through the -A argument).

The second major section in this tree is the models package. This package contains go representations for both the
defintions from the swagger spec document.

And then the last major section is the rest api, within the rest api there is the code that is generated based on the
information from the paths property in the swagger specification. The go swagger generator uses the tags to group the
operations into packages.

We skipped over naming operations, you have the ability to name the operations by giving operations an ID in the
specification document. For example for the operation defintion with `operationId: findTodos`, the following tree would
be generated:

```yaml
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

This shows the basic worflow for using the server generator of swagger.
