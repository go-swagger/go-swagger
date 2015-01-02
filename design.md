# Framework design

The goals are to be as unintrusive as possible. The swagger spec is the source of truth for your application.

The reference framework will make use of a swagger API that is based on the denco router.

The general idea is that it is a middleware which you provide with the swagger spec.
This document can be either JSON or YAML as both are required.

In addition to the middleware there are some generator commands that will use the swagger spec to generate models, parameter models, operation interfaces and a mux.

## The middleware

Takes a raw spec document either as a []byte, and it adds the /api-docs route to serve this spec up.

The middleware performs validation, data binding and security as defined in the swagger spec. 
It also uses the API to match request paths to functions of `func(paramsObject) (responseModel, error)`

When a request comes in that doesn't match the /api-docs endpoint it will look for it in the swagger spec routes.
These are provided in the API. There is a tool to generate a statically typed API, based on operation names and 
operation interfaces

### The API

The reference API will use the denco router to register route handlers.
The actual request handler implementation is always the same.  The API must be designed in such a way that other frameworks can use their router implementation and perhaps their own validation infrastructure.

An API is served over http by a router, the default implementation is a router based on denco. This is just an interface implemenation so it can be replaced with another router should you so desire.

The API comes in 2 flavors an untyped one and a typed one. 

#### Untyped API

The untyped API is the main glue. It takes registrations of operation ids to operation handlers.
It also takes the registrations for mime types to consumers and producers. And it links security schemes to authentication handlers.

```go
type OperationHandler func(interface{}) (interface{}, error)
```

The API has methods to register consumers, producers, auth handlers and operation handlers

The register consumer and producer methods are responsible for attaching extra serializers to media types. These are then used during content negotiation phases for look up and binding the data.

When an API is used to initialize a router it goes through a validation step.
This validation step will verify that all the operations in the spec have a handler registered to them. 
It also ensures that for all the mentioned media types there are consumers and producers provided.
And it checks if for each authentication scheme there is a handler present.
If this is not the case it will exit the application with a non-zero exit code.

The register method takes an operation name and a swagger operation handler.  
It will then use that to build a path pattern for the router and it uses the swagger operation handler to produce a result based on the information in an incoming web request. It does this by injecing the handler in the swagger web request handler.

#### Typed API

The typed API uses a swagger spec to generate a typed API. 

For this there is a generator that will take the swagger spec document.
It will then generate an interface for each operation and optionally a default implementation of that interface.
The default implemenation of an interface just returns a not implemented api error.

When all the interfaces and default implementations are generated it will generate a swagger mux implementation.
This swagger mux implemenation links all the interface implementations to operation names. 
The typed API wraps an untyped API to do the actual route registration, it's mostly sugar for providing better compile time type safety.

This is done through integration in the `go generate` command

### The request handler

The request handler does the following things:

1. Authenticate and authorize if necessary
2. Validate the request data
3. Bind the request data to the parameter struct based on the swagger schema
4. Validate the parameter struct based on the swagger schema
5. Produce a model or an error by invoking the operation interface
6. Create a response with status code etc based on the operation interface invocation result

#### Authentication

TODO: put some coherent sentences here that describe how the auth integration is supposed to work.

Does this make it so that we require a context type object or add a pointer param for the principal on each authenticated operation? 

Maybe it's better add a SwaggerPrincipal property to the operation parameter object?

```go
type SecurityHandler func(*http.Request) (interface{}, error)
```

#### Binding

Binding makes use of plain vanilla golang serializers and they are identified by the media type they consume and produce. 

Binding is not only about request bodies but also about values obtained from headers, query string parameters and potentially the route path pattern. So the binding should make use of the full request object to produce a model.

It determines a serializer to use by looking in the the merged consumes values and the `Content-Type` header to determine which deserializer to use.  
When a result is produced it will do the same thing by making use of the `Accept` http header etc and the merged produces clauses for the operation endpoint. 

#### Validation 

When the muxer registers routes it also builds a suite of validation plans, one for each operation. 
Validation allows for adding custom validations for types through implementing a Validatable interface. This interface does not override but extends the validations provided by the swagger schema. 

There is a mapping from validation name to status code, this mapping is also prioritized so that in the event of multiple validation errors that would required different status codes we get a consistent result. 

```go
type Validatable interface {
  Validate() []Error
}

type Error struct {
  Code     int32
  Path     string
  In       string
  Value    interface{}
  Message  string
}
```


