---
title: Generated server
date: 2023-01-01T01:01:01-08:00
draft: true
---
# Use the generated server

The generated server serves the API but the default implementation returns 501 Not implemented for everything. Let's
look into using the generated code.

Go swagger primarily deals with HTTP and originally only supports the stdlib `net/http` interface. A typical HTTP
request expects a response.  This is reflected in go-swagger where a handler is typically defined as a function of
input parameters to a responder.

```go
type ListTravelsHandler interface {
	Handle(ListTravelsParams) middleware.Responder
}
```

The signature of this handler is one of 2 possible variations. When a handler doesn't use authentication then a handler
interface consists out of input parameters and a responder.

```go
type AddOneAuthenticatedHandler interface {
	Handle(AddOneParams, interface{}) middleware.Responder
}
```

When a handler does use authentication then the second argument to the handler function represents the security
principal for your application. You can specify the type name for this principal at generation time by specifying the
-P or --principal flag.

```
swagger generate server -P models.User
swagger generate client -P models.User
```

The type name can be specified by package path.

```
swagger generate server -P github.com/foobar/models.User
```

See the full list of available options [for server](../generate/server.md) and [for client](../generate/client.md).

When you would execute the generate step with that parameter for the security principal then the
AddOneAuthenticatedHandler would look a bit like this:

```go
type AddOneAuthenticatedHandler interface {
	Handle(AddOneParams, *models.User) middleware.Responder
}
```

## Implement handlers

A handler is an interface/contract that defines a statically typed representation of the input and output parameters of
an operation on your API.
The tool generates handlers that are stubbed with a NotImplemented response when you first generate the server.

### The `not implemented` handler

The not implemented handler is actually a not implemented responder, it returns a responder that will always respond
with status code 501 and a message that lets people know it's not the fault of the client that things don't work.

```go
middleware.NotImplemented("operation todos.AddOne has not yet been implemented")
```

### Your own code

Each HTTP request expects a response of some sort, this response might have no data but it's a response none the less.

Every incoming request is described as a bunch of input parameters which have been validated prior to calling the
handler. So whenever your code is executed, the input parameters are guaranteed to be valid according to what the
swagger specification prescribes.

All the input parameters have been validated, and the request has been authenticated should it have triggered
authentication.

You probably want to return something a bit more useful to the users of your API than a not implemented response.

A possible implementation of the `ListTravelsHandler` interface might look like this:

```go
type PublicListTravelsOK struct {
  Body []models.Travel
}
func (m *PublicListTravelsOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer){
  // generated code here
}

type PublicListTravelsError struct {
  Body models.Error
}
func (m *PublicListTravelsOK) WriteResponse(rw http.ResponseWriter, producer httpkit.Producer){
  // generated code here
}

type PublicListTravelsHandler struct {
  db interface {
    FetchTravels(*PublicListTravelsParams) ([]models.Travel, error)
  }
}

func (m *PublicListTravelsHandler) Handle(params ListTravelsParams) middleware.Responder {
  travels, err := m.db.FetchTravels(&params)
  if err != nil {
    return &PublicListTravelsError{Body: models.Error{Message: err.Error()}}
  }
  return &PublicListTravelsOK{Body: travels}
}
```

In the example above we have a handler implementation with a hypothetical database fetch interface. When the handle
method is executed there are 2 possible responses for the provided parameters. There can either be an error in which
case the PublicListTravelsError will be returned, otherwise the PublicListTravelsOK will be returned.

The code generator has written the remaining code to render that response with the headers etc.

