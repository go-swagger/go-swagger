# Framework design

The goals are to be as unintrusive as possible. The swagger spec is the source of truth for your application.

The reference framework will make use of a swagger muxer that is based on the denco router.

The general idea is that it is a middleware which you provide with the swagger spec.
This document can be either JSON or YAML as both are required.

In addition to the middleware there are some generator commands that will use the swagger spec to generate models, parameter models, operation interfaces and a mux.

## The middleware

Takes a raw spec document either as a []byte, and it adds the /api-docs route to serve this spec up.

The middleware performs validation, data binding and security as defined in the swagger spec. 
It also uses the muxer to match request paths to functions of `func(paramsObject) (responseModel, error)`

When a request comes in that doesn't match the /api-docs endpoint it will look for it in the swagger spec routes.
These are provided in the muxer. There is a tool to generate a statically typed muxer, based on operation names and 
operation interfaces

### The muxer

The reference muxer will use the denco router to register route handlers.
The actual request handler implementation is always the same.  The muxer must be designed in such a way that other frameworks can use their router implementation
and perhaps their own validation infrastructure.

The muxer comes in 2 flavors an untyped one and a typed one. 

#### Untyped muxer

The untyped muxer links operation names to operation handlers

```go
type SwaggerOperationHandler func(interface{}) (interface{}, error)
```

The muxer has 4 methods: RegisterSerializer, RegisterSecurityScheme, Register, Validate.

The register serializer is responsible for attaching extra serializers to media types. These are then used during content negotiation phases for looking up 

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
The typed muxer wraps an untyped muxer to do the actual route registration, it's mostly sugar for providing better compile time type safety.

This is done through integration in the `go generate` command

### The request handler

The request handler does the following things:

1. Authenticate and authorize if necessary
2. Bind the request data to the parameter struct based on the swagger schema
3. Validate the parameter struct based on the swagger schema
4. Produce a model or an error by invoking the operation interface
5. Create a response with status code etc based on the operation interface invocation result

#### Authentication

TODO: put some coherent sentences here that describe how the auth integration is supposed to work.

Does this make it so that we require a context type object or add a pointer param for the principal on each authenticated operation? 

Maybe it's better add a SwaggerPrincipal property to the operation parameter object?

```go
type SecurityHandler func(*http.Request) (interface{}, error)
```

#### Serialization

Binding makes use of plain vanilla golang serializers and they are identified by the media type they consume and produce. 

Binding is not only about request bodies but also about values obtained from headers, query string parameters and potentially the route path pattern. So the binding should make use of the full request object to produce a model.

```go
type Binder func(*http.Request, interface{}) error
```

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
  Code int16
  MessageTemplate *template.Template
  FieldName string
}
```


