---
title: About generating a spec document
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about spec generation -->
## Spec generation from source

### Is there an example to generate a swagger spec document from the code?
_Use-Case_: I have read the swagger.json generation and feel confused. Could you please give an example for it?

**Answer**: this folder uses most of the annotations

https://github.com/go-swagger/go-swagger/tree/master/fixtures/goparsing/petstore

>This begs for 3 questions :
> - Q1: Does a struct for Parameter model have to be declared in the SAME .go file where the swagger:route is declared for a router function?
> - Q2: Assume that I have a route "/services/{serviceName}", how would I named the field associated with the path "{serviceName}" above in a struct wrapper for path params?
> - Q3: Is the annotations case sensitive like "Required" vs "required"?. I see mixed examples about this.

**Answers**:

- Q1: nope. It needs to be declared in the same app and it needs to be imported so that goswagger can find it by following imports starting at the main package.
- Q2: you would add all of them in the parameter struct at this moment, in the case of parameters you would add a doc comment: // in: path. Take a look at some of the generated code examples because they contain all the known annotations as well.
- Q3: not case sensitive, didn't want to have debates over casing. Whatever looks good to you in docs is what you can use.

*One more question: Can methods in an interface be annotated?*

**Answer**: only when it's used in a discriminator IIRC. There is code in the scan package that treats nullary methods as properties if certain conditions are met.

>My generated spec is now working but seems to be missing a parameter "description" to indicate to end user of the API URL endpoint of what's its doing. Example below, I wanted the line "disable/enable a compute node EC2 machine with a given IP address" to show up as some sort of description for the parameter... Am I missing something?

```golang
// v2PutXXX disable/enable a compute node EC2 machine with a given IP address
//
// swagger:route PUT /compute/nodes/{nodeIPAddress} v2PutXXX
//
// Disable/enable a compute node machine with a given IP address
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// Schemes: http
//
// Responses:
// default: errorResp
// 200: okResp
//
func v2PutXXX(....)
```

**Answer**: you still need to add enlist a struct as parameters for the operation.
https://goswagger.io/generate/spec/params.html

*Is there an example how to generate example values from the code?*

**Answer**: I don't think that is supported at the moment

Originally from issue [#213](https://github.com/go-swagger/go-swagger/issues/213).

### Extra function in example?
In file: `go-swagger/fixtures/goparsing/classification/operations/todo_operation.go`,
`Func: mountItem` looks like extra function. Could you explain?

**Answer**: swagger tool generates a correct specification for proposed routes without calling fake func mountItem.

The main rule is **empty line between `swagger:routes`** like that:

```golang
// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) error {
    // swagger:route GET /pets pets users listPets

    // swagger:route GET /orders orders listOrders

    // swagger:route POST /orders orders createOrder

    // swagger:route GET /orders/{id} orders orderDetails

    // swagger:route PUT /orders/{id} orders updateOrder

    // swagger:route DELETE /orders/{id} orders deleteOrder
return nil
}
```

Originally from issue [#68](https://github.com/go-swagger/go-swagger/issues/68).

### Maps as swagger parameters
_Use-case_: I'm using go-swagger to generate my Swagger docs from code, and I came across a problem with a given parameter.

When I annotate a given struct that has a `map[KeyType]OtherKeyType` with `swagger:parameters`,
it returns `items doesn't support maps`.

**Answer**: **this is not supported**

- In non-body parameters maps are not supported in the swagger spec
- In body parameters, a JSON schema only allows maps with string keys

Originally from issue [#960](https://github.com/go-swagger/go-swagger/issues/960).

### How to define a swagger response that produces a binary file?

_Use-case_: annotating a go struct in order to produce a response as application/octet-stream.

Example:
I would like to get a generated specification like:
```JSON
"fileResponse": {
  "description": "OK",
  "schema": {
    "type": "file"
  }
}
```
>However, I am unable to figure out how to do this with go-swagger response struct and annotations.

**Answer**: you can use `runtime.File` or `os.File` in your struct:
```golang
type fileResponse struct {
    // In: body
    File runtime.File
}
```
Originally from issue [#1003](https://github.com/go-swagger/go-swagger/issues/1003).

### How to use swagger params?
_Use-Case_: I defined a route with!
```golang
// swagger:route GET /services/{serviceName}/version/{version} pets listOneService
```

*How to comment the two params('serviceName' and 'version')?*

**Answer**: `swagger:params` is used to indicate which operations the properties of the operation are included in the struct.

So you'd use something like these:
https://github.com/go-swagger/go-swagger/blob/master/fixtures/goparsing/petstore/rest/handlers/orders.go#L24-L46

or:

```golang
// swagger:params listOneService
type ListOneParams struct {
    // ServiceName description goes here
    //
    // in: path
    // required: true
    ServiceName string `json:"serviceName"`

    // Version description goes here
    //
    // in: path
    // required: true
    Version string `json:"version"`
}
```

Originally from issue [#668](https://github.com/go-swagger/go-swagger/issues/668).

### Empty definitions
_Use-Case_: I don't understand how to deal with model annotation.
When I generate a spec from source, I get empty definitions.

**Answer**: models are discovered through usage in parameter and/or response objects.

If the model isn't used through a parameter or response object it's not part of the API because there is no way that it goes in or out through the API.

Example:

doc.go
```golang
// Schemes: http, https
// Host: localhost
// BasePath: /v1
// Version: 0.0.1
// License: MIT http://opensource.org/licenses/MIT
//
// Consumes:
// - application/json
// - application/xml
//
// Produces:
// - application/json
// - application/xml
//
//
// swagger:meta
package main
...
```
user.go
```golang
// Copyright 2015 go-swagger maintainers
//
// ...

package models

// User represents the user for this application
//
// A user is the security principal for this application.
// It's also used as one of main axis for reporting.
//
// A user can have friends with whom they can share what they like.
//
// swagger:model
type User struct {
    // the id for this user
    //
    // required: true
    // min: 1
    ID int64 `json:"id"`

    // the name for this user
    // required: true
    // min length: 3
    Name string `json:"name"`
}
```

Originally from issue [#561](https://github.com/go-swagger/go-swagger/issues/561).

### Documentation or tutorials on code annotation
_Use-Case_: documentation is scant on how to generate swagger files from annotations.
Is it really all there in http://goswagger.io/generate/spec/?

**Answer**: yes, it's all in there (or directly in the repo: https://github.com/go-swagger/go-swagger/tree/master/docs/generate/spec)

*How about some code examples that show annotations being used?*

**Answer**: there is an "examples" folder in the repo.
All generated code also uses all the annotations that are applicable for it.

https://github.com/go-swagger/go-swagger/tree/master/examples/todo-list

And also: https://github.com/go-swagger/go-swagger/tree/master/fixtures/goparsing/classification
(this is the code used to test parsing the annotations).

Please bear in mind that this is a project (not a product) to which a number of volunteers have
contributed significant amounts of free time to get it to where it is today.

Improvement of documentation is always a good request.
All help we can get is absolutely welcome.

Originally from issue [#599](https://github.com/go-swagger/go-swagger/issues/599).

### Wrong schema in response structure?
I set up this response struct:

```golang
// swagger:response SuccessResponse
type SuccessResponse struct {
    // In: body
    Data ResponseData `json:"data"`
}

type ResponseData struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
}
```
Expected schema:
```JSON
{
  "responses": {
    "ErrorResponse": {},
    "SuccessResponse": {
      "description": "SuccessResponse is success response",
      "schema": {
        "$ref": "#/definitions/SuccessResponse"
      }
    }
  }
}
```
but getting instead:
```JSON
{
  "responses": {
    "ErrorResponse": {},
    "SuccessResponse": {
      "description": "SuccessResponse is success response",
      "schema": {
        "$ref": "#/definitions/ResponseData"
      }
    }
  }
}
```
I know this is expected working behavior, but
I don't want to add additional level of structs just to support pretty output.

**Answer**: you can rename the model in the json with the swagger:model doc tag on the response
data struct, that would get you the expected output.

```golang
// swagger:response SuccessResponse
type SuccessResponse struct {
    // In: body
    Data ResponseData `json:"data"`
}

// swagger:model SuccessResponse
type ResponseData struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
}
```

Originally from issue [#407](https://github.com/go-swagger/go-swagger/issues/407).

### go-swagger not generating model info and showing error on swagger UI

_Use-Case_: when I'm executing : `swagger generate spec -o ./swagger.json` to generate the json spec I'm getting:

```JSON
{
  "consumes": ["application/json", "application/xml"],
  "produces": ["application/json", "application/xml"],
  "schemes": ["http", "https"],
  "swagger": "2.0",
  "info": {
    "description": "the purpose of this application is to provide an application\nthat is using plain go code to define an API\n\nThis should demonstrate all the possible comment annotations\nthat are available to turn go code into a fully compliant swagger 2.0 spec",
    "title": "User API.",
    "termsOfService": "there are no TOS at this moment, use at your own risk we take no responsibility",
    "contact": {
      "name": "John Doe",
      "url": "http://john.doe.com",
      "email": "john.doe@example.com"
    },
    "license": {
      "name": "MIT",
      "url": "http://opensource.org/licenses/MIT"
    },
    "version": "0.0.1"
  },
  "host": "localhost",
  "basePath": "/v2",
  "paths": {
    "/user": {
      "get": {
        "description": "This will show all available pets by default.\nYou can get the pets that are out of stock",
        "consumes": ["application/json", "application/x-protobuf"],
        "produces": ["application/json", "application/x-protobuf"],
        "schemes": ["http", "https", "ws", "wss"],
        "tags": ["listPets", "pets"],
        "summary": "Lists pets filtered by some parameters.",
        "operationId": "users",
        "security": [{
          "api_key": null
         },{
          "oauth": ["read", "write"]
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/someResponse"
          },
          "422": {
            "$ref": "#/responses/validationError"
          },
          "default": {
            "$ref": "#/responses/genericError"
          }
       }
    }
  }
},
"definitions": {}
}
```

Note that my definitions are empty, not sure why. If I paste the same json spec in http://editor.swagger.io/#/ It says
```
Error
Object
message: "options.definition is required"
code: "UNCAUGHT_SWAY_WORKER_ERROR"
```
Any directions on what is the right way to generate swagger documentation would help

**Answer**: can you move the `swagger:model` annotation to be the last line in the doc comments for a struct?

Alternatively, I see some definitions for responses in your specification document, but no matching `swagger:response` definitions structs.
```golang
// swagger:response errorResponse
type ErrorResponse struct {
    // in: body
    Body struct {
        Message string `json:"error,omitempty"`
    }
}

// swagger:response validationError
type ValidationError struct {
    // in: body
    Body struct {
        // required: true
        Message string `json:"error,omitempty"`

        Field string `json:"fieldName,omitempty"`
    }
}

// swagger:response someResponse
type SomeResponse struct {
    // in: body
    Body *User `json:"body,omitempty"`
    }
```
With the `--scan-models` generating option, you should have models picked up, regardless of whether they're in use somewhere else or not.

<!-- Would need a recap/update on that
### Running on google app engine
_Use-Case_: generating a spec for an app built for GoogleApp engine

> App engine apps use some package imports which don't resolve when run with `go build: "appengine","appengine/datastore"`, etc. 
> It seems like `swagger generate spec` fails if it can't first build my app.

*To support app engine you'd need to remove that requirement.*

> I would like to use go-swagger with my app engine project, so please make it parse the comments without first needing to build the app.

> I tried adding the appengine build constraint, and it didn't error out, but it generated an empty spec.
> at this point we make use of the go loader package. This allows us to discover your application and which files to scan for the doc comments.

This application needs to read composed structs and so on, and it's a lot easier to interrogate the application if you know where all the files are and not just the ones you created in this particular folder.
Unfortunately it does require to be able to read

How about a main class that doesn't require appengine imports?
I've personally never used appening so I don't really know what is involved.

For an app engine app, typically the main.go doesn't require any of the special appengine imports... however it uses an init() function instead of a main() function, and in there is where you instantiate the router and bind all the routes to handlers. It is the handlers (usually in their own separate files) which need the appengine imports.

I don't know if that's what you are asking about, but I know app engine fairly well so I can answer more questions if you have any.

I guess what I mean is if your code doesn't compile, how are you running it?
And what I also meant is; if this is important to you, you could look at forking and submitting a pull request.
I, personally, still have a bunch of other things that need fixing in here before I want to look at a niche like appengine.

Fair enough, if you're not interested in supporting app engine that's fine. I did try forking to fix that other bug that I filed, but the code was a bit over my head so I wasn't able to fix it.

To answer your question, app engine apps don't get compiled with "go build", instead they are run on a dev server provided by the app engine SDK, and then they are deployed to app engine and run on the Google Cloud Platform infrastructure. The only reason they don't compile is because some of the packages ("appengine", "appengine/datastore", etc.) are only available in this SDK environment, they are not found in $GOPATH.

Maybe app engine could improve this situation in the future, and then go-swagger wouldn't have to change to support it, but as it stands now this will not work with any app engine apps that use any of the appengine-specific imports.

I will go back to using github.com/yvasiyarov/swagger for now, which doesn't require the app to build to generate the spec, but it is also not generating swagger 2.0, so I hope I can use your package sometime in the future.

but the custom sdk they use also includes a custom go command doesn't it?
I think the problem you're having is related to your GOPATH content and can be fixed there.

afaik go always needs to compile your stuff, whether that's in the SDK env or not, have you tried installing go-swagger in the SDK provided GOPATH?

I'll leave this open so i can track this
Hmm that is an interesting thought. I will experiment more on this today.
The app engine SDK definitely uses the regular system GOPATH to resolve most of the includes, but it maybe has another internal GOPATH also,
I'm not sure. Will post my findings a bit later.

Wow! You're right man! All that was needed to make generate spec work, was to add this to GOPATH:
[go_appengine_sdk_location]/goroot

The appengine includes are in there. It's working now, thanks for the insight!

Originally from issue [#47](https://github.com/go-swagger/go-swagger/issues/47).
-->
<!-- Obsolete / Not helpfu
### Generating spec cannot import dependencies
I'm unable to run the spec generator on my app. What am I doing wrong?
https://gist.github.com/morenoh149/e44be6819bde86f52e7e

I get many errors of the form ... .go:10:2: could not import github.com/facebookgo/stackerr (cannot find package "github.com/ ...

execute $ swagger generate spec -o ./swagger.json in project folder

which go version are you using?
Are you using vendoring?

go version go1.5.3 darwin/amd64

I believe I have vendoring enabled yes. (Not sure I'm new).

for this you require `export GO15VENDOREXPERIMENT=1`

and also does your application compile? because go-swagger makes use of the same code the go compiler/imports/... etc use to discover all the involved packages. So I think it requires your code to be mostly compilable to be able to discover everything.

However the errors seem to be related to it not being able to discover the dependencies

I am running into this issue when trying to generate a swagger spec as well. I believe it is because I installed go with home brew on a mac so my library paths are different.

Command I am running: swagger generate spec
File I am trying to generate with: https://gist.github.com/gsquire/cce277b4bd10ba283f4522e896dc91d6

Error trace:
```
/Users/gsquire/scratch/test-swagger.go:14:2: could not import fmt (cannot find package "fmt" in any of:
/usr/local/go/src/fmt (from $GOROOT)
/Users/gsquire/gopath/src/fmt (from $GOPATH))
/Users/gsquire/scratch/test-swagger.go:15:2: could not import net/http (cannot find package "net/http" in any of:
/usr/local/go/src/net/http (from $GOROOT)
/Users/gsquire/gopath/src/net/http (from $GOPATH))
/Users/gsquire/scratch/test-swagger.go:18:14: undeclared name: http
/Users/gsquire/scratch/test-swagger.go:18:38: undeclared name: http
/Users/gsquire/scratch/test-swagger.go:19:2: undeclared name: fmt
/Users/gsquire/scratch/test-swagger.go:23:2: undeclared name: http
/Users/gsquire/scratch/test-swagger.go:24:2: undeclared name: http
```

Are you on a Linux box? Like I said, I think it's a path issue since home brew installs it in a different location based on my output. I'll have to try on something that isn't mac.

I uninstalled go-swagger from brew, installed it from source and it worked. So strange. And no, I installed it through brew:

Originally from issue [#400](https://github.com/go-swagger/go-swagger/issues/400).
-->

