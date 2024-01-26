---
title: Dynamic API
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Dynamic API definition

The toolkit supports building a swagger specification entirely with go code. It does allow you to serve a spec up quickly. This is one of the building blocks required to serve up stub APIs and to generate a test server with predictable responses, however this is not as bad as it sounds...

<!--more-->

This tutorial uses the todo list application to serve a swagger based API defined entirely in go code.
Because we know what we want the spec to look like, first we'll just build the entire spec with the internal dsl.

## Loading the specification

```go
package main

import (
	"log"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Would be serving:", specDoc.Spec().Info.Title)
}
```

[see source of this code](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/dynamic-1/main.go)

Running this would confirm that we can in fact read a swagger spec from disk. 
The init method enables loading of yaml based specifications. The yaml package for golang used to be licensed as GPL so we made depending on it optional. 

```sh
git:(master) ✗ !? » go run main.go ./swagger.yml  
2016/10/08 20:50:42 loading "./swagger.yml" as contract for the server
2016/10/08 20:50:42 Would be serving: A To Do list application
```

## Setup

Before we can implement our API we'll look at setting up the server for our openapi spec.
Go-swagger wants you to configure your API with an api descriptor so that it knows how to handle requests.

### Validation of requirements

It's probably a good idea to fail starting the server when it can't fulfill all the requests defined in the swagger spec.
So let's start by enabling that validation:

```go
func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	api := untyped.NewAPI(specDoc)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}
	log.Println("serving:", specDoc.Spec().Info.Title)

}
```

[see source of this code](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/dynamic-setup-invalid/main.go)

This code shows how to create an api descriptor and then invoking its verification.
Because our specification contains operations and consumes/produces definitions this program should not run.
When we try to run it, it should exit with a non-zero status.

```sh
git:(master) ✗ -? » go run main.go ./swagger.yml
2016/10/08 21:32:14 loading "./swagger.yml" as contract for the server
2016/10/08 21:32:14 missing [application/io.goswagger.examples.todo-list.v1+json] consumes registrations
missing from spec file [application/json] consumes
exit status 1
```

### Satisfying validation with stubs

For us to be able to start our server we will register the right serializers and we'll stub out the operation handlers with a not implemented handler.

```go
func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// our spec doesn't have application/json in the consumes or produces
	// so we need to clear those settings out
	api := untyped.NewAPI(specDoc).WithoutJSONDefaults()

	// register serializers
	mediaType := "application/io.goswagger.examples.todo-list.v1+json"
	api.DefaultConsumes = mediaType
	api.DefaultProduces = mediaType
	api.RegisterConsumer(mediaType, runtime.JSONConsumer())
	api.RegisterProducer(mediaType, runtime.JSONProducer())
	
  api.RegisterOperation("GET", "/", notImplemented)
	api.RegisterOperation("POST", "/", notImplemented)
	api.RegisterOperation("PUT", "/{id}", notImplemented)
	api.RegisterOperation("DELETE", "/{id}", notImplemented)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}

	// construct the application context for this server
	// use the loaded spec document and the api descriptor with the default router
	app := middleware.NewContext(specDoc, api, nil)

	log.Println("serving", specDoc.Spec().Info.Title, "at http://localhost:8000")
	// serve the api
	if err := http.ListenAndServe(":8000", app.APIHandler(nil)); err != nil {
		log.Fatalln(err)
	}
}

var notImplemented = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	return middleware.NotImplemented("not implemented"), nil
})
```

The untyped API for go-swagger assumes by default you want to serve `application/json` and initializes the descriptor with default values to that effect.
In our spec however we don't serve 'application/json' which means we have to use `WithoutJSONDefaults` when we initialize our api.

The media type we do know is: `application/io.goswagger.examples.todo-list.v1+json`, this is also a json format.
We set it as defaults and register the appropriate consumer and producer functions.

Our specification has 4 methods: findTodos, addOne, updateOne and destroyOne. Because we have no implementation yet, we register a notImplemented handler for all of them.

Our api descriptor validation is now satisfied, so we use the simplest way to start a http server in go on port 8000.

Server terminal:

```sh
git:(master) ✗ -!? » go run main.go ./swagger.yml
2016/10/08 23:35:18 loading "./swagger.yml" as contract for the server
2016/10/08 23:35:18 serving A To Do list application at http://localhost:8000
```

Client terminal:

```sh
git:(master) ✗ -!? » curl -i localhost:8000
```

```http
HTTP/1.1 501 Not Implemented
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 06:36:11 GMT
Content-Length: 18

"not implemented"
```

> There is a lot more to be done to make this server a production ready server, but for the
> purpose of this tutorial, this is enough.

## Completely untyped

At this point we're ready to actually implement some functionality for our Todo list. We'll create methods to add, update and delete an item.
We'll also render a list of known items. Because http APIs can get concurrent access we need to take care of this as well.

The first thing we'll do is build our "backend", a very simple implementation based on a slice and maps.

```go
var items = []map[string]interface{}{
	map[string]interface{}{"id": int64(1), "description": "feed dog", "completed": true},
	map[string]interface{}{"id": int64(2), "description": "feed cat"},
}

var itemsLock = &sync.Mutex{}
var lastItemID int64 = 2

func newItemID() int64 {
	return atomic.AddInt64(&lastItemID, 1)
}

func addItem(item map[string]interface{}) {
	itemsLock.Lock()
	defer itemsLock.Unlock()
	item["id"] = newItemID()
	items = append(items, item)
}

func updateItem(id int64, body map[string]interface{}) (map[string]interface{}, error) {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	item, err := itemByID(id)
	if err != nil {
		return nil, err
	}
	delete(body, "id")
	for k, v := range body {
		item[k] = v
	}
	return item, nil
}

func removeItem(id int64) {
	itemsLock.Lock()
	defer itemsLock.Unlock()

	var newItems []map[string]interface{}
	for _, item := range items {
		if item["id"].(int64) != id {
			newItems = append(newItems, item)
		}
	}
	items = newItems
}

func itemByID(id int64) (map[string]interface{}, error) {
	for _, item := range items {
		if item["id"].(int64) == id {
			return item, nil
		}
	}
	return nil, errors.NotFound("not found: item %d", id)
}
```

[see source of this code](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/dynamic-untyped/main.go)

The backend code builds a todo-list-item store that's save for concurrent access buy guarding every operation with a lock. This is all in memory so as soon as you quit the process all your changes will be reset.

Because we now have an actual implementation that we can use for testings, lets hook that up in our API:

```go
func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// our spec doesn't have application/json in the consumes or produces
	// so we need to clear those settings out
	api := untyped.NewAPI(specDoc).WithoutJSONDefaults()

	// register serializers
	mediaType := "application/io.goswagger.examples.todo-list.v1+json"
	api.DefaultConsumes = mediaType
	api.DefaultProduces = mediaType
	api.RegisterConsumer(mediaType, runtime.JSONConsumer())
	api.RegisterProducer(mediaType, runtime.JSONProducer())

	// register the operation handlers
	api.RegisterOperation("GET", "/", findTodos)
	api.RegisterOperation("POST", "/", addOne)
	api.RegisterOperation("PUT", "/{id}", updateOne)
	api.RegisterOperation("DELETE", "/{id}", destroyOne)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}

	// construct the application context for this server
	// use the loaded spec document and the api descriptor with the default router
	app := middleware.NewContext(specDoc, api, nil)

	log.Println("serving", specDoc.Spec().Info.Title, "at http://localhost:8000")

	// serve the api with spec and UI
	if err := http.ListenAndServe(":8000", app.APIHandler(nil)); err != nil {
		log.Fatalln(err)
	}
}

var findTodos = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'findTodos'")
	log.Printf("%#v\n", params)

	return items, nil
})

var addOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'addOne'")
	log.Printf("%#v\n", params)

	body := params.(map[string]interface{})["body"].(map[string]interface{})
	addItem(body)
	return body, nil
})

var updateOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'updateOne'")
	log.Printf("%#v\n", params)

	data := params.(map[string]interface{})
	id := data["id"].(int64)
	body := data["body"].(map[string]interface{})
	return updateItem(id, body)
})

var destroyOne = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'destroyOne'")
	log.Printf("%#v\n", params)

	removeItem(params.(map[string]interface{})["id"].(int64))
	return nil, nil
})
```

[see source of this code](https://github.com/go-swagger/go-swagger/blob/master/examples/tutorials/todo-list/dynamic-untyped/main.go)

With this set up we should be able to start a server, send it some requests and get some meaningful answers.

#### List all

```sh
git:(master) ✗ !? » curl -i localhost:8000
```

```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:50:39 GMT
Content-Length: 87

[{"completed":true,"description":"feed dog","id":1},{"description":"feed cat","id":2}]
```

#### Create new

The default curl POST request should fail because we only allow:  application/io.goswagger.examples.todo-list.v1+json

```
curl -i localhost:8000 -d '{"description":"item for the list"}'
```

```http
HTTP/1.1 415 Unsupported Media Type
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:55:43 GMT
Content-Length: 157

{"code":415,"message":"unsupported media type \"application/x-www-form-urlencoded\", only [application/io.goswagger.examples.todo-list.v1+json] are allowed"}
```

When the content type header is sent, we have a better result:

```
curl -i -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' localhost:8000 -d '{"description":"a new item"}'
```

```http
HTTP/1.1 201 Created
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:56:28 GMT
Content-Length: 36

{"description":"a new item","id":3}
```

#### List again

```sh
git:(master) ✗ !? » curl -i localhost:8000
```

```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:58:06 GMT
Content-Length: 123

[{"completed":true,"description":"feed dog","id":1},{"description":"feed cat","id":2},{"description":"a new item","id":3}]
```

#### Update an item

```sh
curl -i -XPUT -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' localhost:8000/3 -d '{"description":"an updated item"}'
```

```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:58:42 GMT
Content-Length: 41

{"description":"an updated item","id":3}
```

#### List to verify

```sh
git:(master) ✗ !? » curl -i localhost:8000
```

```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 15:58:42 GMT
Content-Length: 41

{"description":"an updated item","id":3}
```

#### Delete an item

```sh
curl -i -XDELETE localhost:8000/3
```

```http
HTTP/1.1 204 No Content
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 16:00:59 GMT
```

#### List to show start state again

```
curl -i localhost:8000
```

```http
HTTP/1.1 200 OK
Content-Type: application/io.goswagger.examples.todo-list.v1+json
Date: Sun, 09 Oct 2016 16:02:19 GMT
Content-Length: 87

[{"completed":true,"description":"feed dog","id":1},{"description":"feed cat","id":2}]
```

At the end of the curl requests the server shows these outputs:

```sh
git:(master) ✗ !? » go run main.go ./swagger.yml
2016/10/09 08:50:34 loading "./swagger.yml" as contract for the server
2016/10/09 08:50:34 serving A To Do list application at http://localhost:8000
2016/10/09 08:50:39 received 'findTodos'
2016/10/09 08:50:39 map[string]interface {}{"since":0, "limit":20}
2016/10/09 08:56:28 received 'addOne'
2016/10/09 08:56:28 map[string]interface {}{"body":map[string]interface {}{"description":"a new item"}}
2016/10/09 08:58:06 received 'findTodos'
2016/10/09 08:58:06 map[string]interface {}{"limit":20, "since":0}
2016/10/09 08:58:42 received 'updateOne'
2016/10/09 08:58:42 map[string]interface {}{"id":3, "body":map[string]interface {}{"description":"an updated item"}}
2016/10/09 09:00:07 received 'findTodos'
2016/10/09 09:00:07 map[string]interface {}{"since":0, "limit":20}
2016/10/09 09:00:59 received 'destroyOne'
2016/10/09 09:00:59 map[string]interface {}{"id":3}
2016/10/09 09:02:19 received 'findTodos'
2016/10/09 09:02:19 map[string]interface {}{"since":0, "limit":20}
```
