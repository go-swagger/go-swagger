+++
categories = ["tutorial"]
date = "2015-12-30T11:56:40-08:00"
tags = ["todo app", "client", "server"]
title = "Todo List Tutorial"
weight = 7
series = ["tutorials", "home"]
+++

This example walks you through a hypothetical project that is building a todo list.
It specifically uses a todo list because it's a super well-understood application and because of that you can focus entirely on the new concepts. This example builds a server and then a client.

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

But wait a minute, what if there are 100's of todo items, will we just return all of them for everybody?

