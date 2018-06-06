# Generate a server for a swagger spec

The toolkit has a command that will let you generate a docker friendly server with support for TLS.
You can configure it through environment variables that are commonly used on PaaS services.

<!--more-->

A generated server uses _no reflection_ except for an enum validation and the required validation. The server builds all the necessary plans and execution paths at startup time so that at runtime there is only the absolute minimum processing required to respond to requests.

The default router for go-swagger is [naoina's denco](https://github.com/naoina/denco) which is a [**very** fast](https://github.com/julienschmidt/go-http-routing-benchmark#github) ternary search tree based router that allows for much greater flexibility than the trie based router implementation of julienschmidt at almost the same and sometimes lower cost.

You can provide your own router implementation should you so desire it's abstracted through an interface with this use case in mind.

##### Usage

```
swagger [OPTIONS] generate server [server-OPTIONS]

generate all the files for a server application

Help Options:
  -h, --help                                         Show this help message

[server command options]
      -f, --spec=                                    the spec file to use (default swagger.{json,yml,yaml})
      -a, --api-package=                             the package to save the operations (default: operations)
      -m, --model-package=                           the package to save the models (default: models)
      -s, --server-package=                          the package to save the server specific code (default: restapi)
      -c, --client-package=                          the package to save the client specific code (default: client)
      -t, --target=                                  the base directory for generating the files (default: ./)
      -T, --template-dir=                            alternative template override directory
      -C, --config-file=                             configuration file to use for overriding template options
      -A, --name=                                    the name of the application, defaults to a mangled value of info.title
      -O, --operation=                               specify an operation to include, repeat for multiple
          --tags=                                    the tags to include, if not specified defaults to all
      -P, --principal=                               the model to use for the security principal
          --default-scheme=                          the default scheme for this API (default: http)
      -M, --model=                                   specify a model to include, repeat for multiple
          --skip-models                              no models will be generated when this flag is specified
          --skip-operations                          no operations will be generated when this flag is specified
          --skip-support                             no supporting files will be generated when this flag is specified
          --exclude-main                             exclude main function, so just generate the library
          --exclude-spec                             don't embed the swagger specification
          --with-context                             handlers get a context as first arg
          --dump-data                                when present dumps the json for the template generator instead of generating files
          --flag-strategy=[go-flags|pflag]           the strategy to provide flags for the server (default: go-flags)
          --compatibility-mode=[modern|intermediate] the compatibility mode for the tls server (default: modern)
          --skip-validation                          skips validation of spec prior to generation
      -r, --copyright-file=                          the file containing a copyright header for the generated source
          --additional-initialism=                   additional consecutive capitals that should be considered as initialism, repeat for multiple
```

The server application gets generated with all the handlers stubbed out with a not implemented handler. That means that you can start the API server immediately after generating it. It will respond to all valid requests with 501 Not Implemented. When a request is invalid it will most likely respond with an appropriate 4xx response.

The generated server allows for a number of command line parameters to customize it.

```
--host=            the IP to listen on (default: localhost) [$HOST]
--port=            the port to listen on for insecure connections, defaults to a random value [$PORT]
--tls-host=        the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
--tls-port=        the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
--tls-certificate= the certificate to use for secure connections [$TLS_CERTIFICATE]
--tls-key=         the private key to use for secure conections [$TLS_PRIVATE_KEY]
```

The server takes care of a number of things when a request arrives:

* routing
* authentication
* input validation
* content negotiation
* parameter and body binding

To illustrate this with a pseudo handler, this is what happens in a request.

```go
import (
  "net/http"

  "github.com/go-openapi/errors"
  "github.com/go-openapi/runtime/middleware"
  "github.com/gorilla/context"
)

func newCompleteMiddleware(ctx *middleware.Context) http.Handler {
  return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
    defer context.Clear(r)

    // use context to lookup routes
    if matched, ok := ctx.RouteInfo(r); ok {

      if len(matched.Authenticators) > 0 {
        if _, err := ctx.Authorize(r, matched); err != nil {
          ctx.Respond(rw, r, matched.Produces, matched, err)
          return
        }
      }

      bound, validation := ctx.BindAndValidate(r, matched)
      if validation != nil {
        ctx.Respond(rw, r, matched.Produces, matched, validation)
        return
      }

      result, err := matched.Handler.Handle(bound)
      if err != nil {
        ctx.Respond(rw, r, matched.Produces, matched, err)
        return
      }

      ctx.Respond(rw, r, matched.Produces, matched, result)
      return
    }

    // Not found, check if it exists in the other methods first
    if others := ctx.AllowedMethods(r); len(others) > 0 {
      ctx.Respond(rw, r, ctx.spec.RequiredProduces(), nil, errors.MethodNotAllowed(r.Method, others))
      return
    }
    ctx.Respond(rw, r, ctx.spec.RequiredProduces(), nil, errors.NotFound("path %s was not found", r.URL.Path))
  })
}
```

Prior to handling requests however you probably want to configure the API with some actual implementations.  To do that you have to edit the configure_xxx.go file.  That file will only be generated the first time you generate a server application from a swagger spec. So the generated server uses this file to let you fill in the blanks.

For the todolist application that file looks like:

```go
package main

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-openapi/examples/todo-list/restapi/operations"
	"github.com/go-openapi/examples/todo-list/restapi/operations/todos"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureAPI(api *operations.ToDoListAPI) http.Handler {
	// configure the api here
	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.KeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (key) x-petstore-token from header has not yet been implemented")
	}

	api.AddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation addOne has not yet been implemented")
	})
	api.DestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation destroyOne has not yet been implemented")
	})
	api.FindHandler = todos.FindHandlerFunc(func(params todos.FindParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation find has not yet been implemented")
	})
	api.UpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation updateOne has not yet been implemented")
	})

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}
```

When you look at the code for the configureAPI method then you'll notice that the api object has properties for consumers.
A consumer is an object that can marshal things from a wireformat to an object.  Consumers and their counterpart producers who write objects get their names generated from the consumes and produces properties on a swagger specification.

The interface definitions for consumers and producers look like this:

```go
// ConsumerFunc represents a function that can be used as a consumer
type ConsumerFunc func(io.Reader, interface{}) error

// Consume consumes the reader into the data parameter
func (fn ConsumerFunc) Consume(reader io.Reader, data interface{}) error {
	return fn(reader, data)
}

// Consumer implementations know how to bind the values on the provided interface to
// data provided by the request body
type Consumer interface {
	// Consume performs the binding of request values
	Consume(io.Reader, interface{}) error
}

// ProducerFunc represents a function that can be used as a producer
type ProducerFunc func(io.Writer, interface{}) error

// Produce produces the response for the provided data
func (f ProducerFunc) Produce(writer io.Writer, data interface{}) error {
	return f(writer, data)
}

// Producer implementations know how to turn the provided interface into a valid
// HTTP response
type Producer interface {
	// Produce writes to the http response
	Produce(io.Writer, interface{}) error
}
```

So it's something that can turn a reader into a hydrated interface. A producer is the counterpart of a consumer and writes objects to an io.Writer.  When you configure an api with those you make sure it can marshal the types for the supported content types.

Go swagger automatically provides consumers and producers for known media types. To register a new mapping for a media
type or to override an existing mapping, call the corresponding API functions in your configure_xxx.go file:

```go
func configureAPI(api *operations.ToDoListAPI) http.Handler {
	// other setup code here...
	
	api.RegisterConsumer("application/pkcs10", myCustomConsumer)
	api.RegisterProducer("application/pkcs10", myCustomProducer)
}

``` 

The next thing that happens in the configureAPI method is setting up the authentication with a stub handler in this case. This particular swagger specification supports token based authentication and as such it wants you to configure a token auth handler.  Any error for an authentication handler is assumed to be an invalid authentication and will return the 401 status code.

```go
// UserPassAuthentication authentication function
type UserPassAuthentication func(string, string) (interface{}, error)

// TokenAuthentication authentication function
type TokenAuthentication func(string) (interface{}, error)

// AuthenticatorFunc turns a function into an authenticator
type AuthenticatorFunc func(interface{}) (bool, interface{}, error)

// Authenticate authenticates the request with the provided data
func (f AuthenticatorFunc) Authenticate(params interface{}) (bool, interface{}, error) {
	return f(params)
}

// Authenticator represents an authentication strategy
// implementations of Authenticator know how to authenticate the
// request data and translate that into a valid principal object or an error
type Authenticator interface {
	Authenticate(interface{}) (bool, interface{}, error)
}
```

So we finally get to configuring our route handlers. For each operation there exists an interface so that implementations have some freedom to provide alternative implentations. For example mocks in certain tests, automatic stubbing handlers, not implemented handlers. Let's look at the addOne handler in a bit more detail.

```go
// AddOneHandlerFunc turns a function with the right signature into a add one handler
type AddOneHandlerFunc func(AddOneParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn AddOneHandlerFunc) Handle(params AddOneParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// AddOneHandler interface for that can handle valid add one params
type AddOneHandler interface {
	Handle(AddOneParams, interface{}) middleware.Responder
}
```

Because the `addOne` operation requires authentication, this interface definition requires 2 arguments. The first argument is about the request parameters and the second parameter is the security principal for the request.  In this case it is of type `interface{}`, typically that is a type like Account, User, Session, ...

It is your job to provide such a handler. Go swagger guarantees that by the time the request processing ends up at the handler, the parameters and security principal have been bound and validated.  So you can safely proceed with saving the request body to some persistence medium perhaps.

There is a context that gets created where the handlers get wired up into a `http.Handler`. For the add one this looks like this:

```go
// NewAddOne creates a new http.Handler for the add one operation
func NewAddOne(ctx *middleware.Context, handler AddOneHandler) *AddOne {
	return &AddOne{Context: ctx, Handler: handler}
}

/*AddOne swagger:route POST / todos addOne

AddOne add one API

*/
type AddOne struct {
	Context *middleware.Context
	Params  AddOneParams
	Handler AddOneHandler
}

func (o *AddOne) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
```

The `http.Handler` implementation takes care of authentication, binding, user code execution and generating a response. For authentication this request would end up in the `TokenAuthentication` handler that was put on the api context object earlier.  When a request is authenticated it gets bound. This operation eventually requires an object that is an implementation of `RequestBinder`.  The `AddOneParams` are such an implementation:

```go
// RequestBinder is an interface for types to implement
// when they want to be able to bind from a request
type RequestBinder interface {
	BindRequest(*http.Request, *MatchedRoute) error
}

// AddOneParams contains all the bound params for the add one operation
// typically these are obtained from a http.Request
//
// swagger:parameters addOne
type AddOneParams struct {
	/*
	  In: body
	*/
	Body *models.Item
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *AddOneParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	var body models.Item
	if err := route.Consumer.Consume(r.Body, &body); err != nil {
		res = append(res, errors.NewParseError("body", "body", "", err))
	} else {
		if err := body.Validate(route.Formats); err != nil {
			res = append(res, err)
		}

		if len(res) == 0 {
			o.Body = &body
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
```

In this example there is only a body parameter, so we make use of the selected consumer to read the request body and turn it into an instance of models.Item. When the body parameter is bound, it gets validated and when validation passes no error is returned and the body property is set.  After a request is bound and validated the parameters and security principal are passed to the request handler. For this configuration that would return a 501 responder.

Go swagger uses responders which are an interface implementation for things that can write to a response. For the generated server there are status code response and a default response object generated for every entry in the spec. For the `addOne` operation that are 2 objects one for the success case (201) and one for an error (default).

```go
// Responder is an interface for types to implement
// when they want to be considered for writing HTTP responses
type Responder interface {
	WriteResponse(http.ResponseWriter, runtime.Producer)
}

/*AddOneCreated Created

swagger:response addOneCreated
*/
type AddOneCreated struct {

	// In: body
	Payload *models.Item `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *AddOneCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*AddOneDefault error

swagger:response addOneDefault
*/
type AddOneDefault struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// WriteResponse to the client
func (o *AddOneDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
```

So an implementer of the `AddOneHandler` could return one of these 2 objects and go-swagger is able to respect the contract set forward by the spec document.

So to implement the AddOneHandler you could do something like this.

```go
todos.AddOneHandlerFunc(func(params todos.AddOneParams, principal interface{}) middleware.Responder {
  created, err := database.Save(params.Body)
  if err != nil {
    return AddOneDefault{models.Error{500, err.Error()}}
  }
  return AddOneCreated{created}
})
```
