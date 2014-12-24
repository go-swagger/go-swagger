# Framework design

The goals are to be as unintrusive as possible.

The reference framework will make use of a swagger muxer that is based on the denco router.

The general idea is that it is a middleware which you provide with the swagger spec.
This document can be either JSON or YAML as both are required.

## The middleware

Takes a raw spec document either as a []byte, and it adds the /api-docs route to serve this spec up.

The middleware performs validation, data binding and security as defined in the swagger spec. 
It also uses the muxer to match request paths to functions of `func(paramsObject) (responseModel, error)`

When a request comes in that doesn't match the /api-docs endpoint it will look for it in the swagger spec routes.
These are provided in the muxer. There is a tool to generate a statically typed muxer, based on operation names and 
operation interfaces

### The muxer

The reference muxer will use the denco router to register route handlers.
The actual request handler implementation is always the same.

The muxer comes in 2 flavors an untyped one and a typed one. 

#### Untyped muxer

The untyped muxer links operation names to operation handlers

```go
type SwaggerOperationHandler func(interface{}) (interface{}, error)
```

The muxer has 2 methods: Register, Validate.

The validate method will verify that all the operations in the spec have a handler registered to them. 
If this is not the case it will exit the application with a non-zero exit code.

The register method takes an operation name and a swagger operation handler. 
It will then use that to build a path pattern for the router and it uses the swagger operation handler to produce a result
based on the information in an incoming web request. It does this by injecing the handler in the swagger web request handler.

#### Typed muxer

The typed muxer uses a swagger spec to generate a typed muxer. 

For this there is a generator that will take the swagger spec document.
It will then generate an interface for each operation and optionally a default implementation of that interface.
The default implemenation of an interface just returns a not implemented api error.

When all the interfaces and default implementations are generated it will generate a swagger mux implementation.
This swagger mux implemenation links all the interface implementations to operation names.

This is done through integration in the `go generate` command

### The request handler

The request handler does the following things:

1. Authenticate and authorize if necessary
2. Bind the request data to the parameter struct based on the swagger schema
3. Validate the parameter struct based on the swagger schema
4. Produce a model or an error by invoking the operation interface
5. Create a response with status code etc based on the operation interface invocation result
