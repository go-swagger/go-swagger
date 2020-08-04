# Custom Server Tutorial

In this tutorial we'll be building up a custom server.
The core will be generated using a manually written and maintained OpenAPI 2.0 spec.
The cli code will be a thin layer around that, and will simply setup the API and server,
using the parsed configurations and our hand written handlers.

<!--more-->

The server we'll building will be very simple.
In this tutorial we'll assume you are already familiar with defining an API,
using the OpenAPI 2.0 yaml specification format.
Please consult [the official OpenAPI 2.0 specification][OpenAPI2.0]
for more information in case you're new to OpenAPI (Also known as Swagger).

The end product of this tutorial can be found as `./examples/tutorials/custom-server`.

The server we'll be building, will be generated using the following spec:

```yaml
---
swagger: '2.0'
info:
  version: 1.0.0
  title: Greeting Server
paths:
  /hello:
    get:
      produces:
        - text/plain
      parameters:
        - name: name
          required: false
          type: string
          in: query
          description: defaults to World if not given
      operationId: getGreeting
      responses:
        200:
          description: returns a greeting
          schema:
              type: string
              description: contains the actual greeting as plain text
```

As you can see, there is only 1 operation,
allowing us to focus on how to create a custom server,
without losing track in the details of any specific implementation.

Where you store the specification is not important.
By default the [swagger cli][go-swagger] expects it to be stored as `./swagger.yml`,
but we'll store it as `./swagger/swagger.yml`, to keep our project's root folder clean and tidy.

Once we have our OpenAPI specification ready, it is time to generate our server.
This can be done using the following command, from within our root folder:

```
$ swagger generate server -t gen -f ./swagger/swagger.yml --exclude-main -A greeter
```

In the command above we're specifying the `-t` (target) flag,
specifying that swagger should store all generated code in the given _target_ directory. We're also specifying the `-f` flag, to explicitly define that our spec file can be found at `./swagger/swagger.yml`, rather then the default `./swagger.yml`. As we are writing a custom server, we also don't want the automatically generated `cmd` server, and is excluded using the `--exclude-main` flag. Finally we're also explicitly naming our server using the `-A` flag.
Please consult `swagger generate server --help` for more flags and information.

Once we've executed this command, you should have following file tree:

```
├── gen
│   └── restapi
│       ├── configure_greeter.go
│       ├── doc.go
│       ├── embedded_spec.go
│       ├── operations
│       │   ├── get_greeting.go
│       │   ├── get_greeting_parameters.go
│       │   ├── get_greeting_responses.go
│       │   ├── get_greeting_urlbuilder.go
│       │   └── greeter_api.go
│       └── server.go
└── swagger
    └── swagger.yml
```

After generation we should find only 1 sub directory in our `gen` folder.
`restapi`, which contains the core server. It consists among other things out of the API (`operations/greeter_api.go`), operations (`operations/*`) and parameters (`*_parameters.go`).

Note that in case we also had defined models in the global definitions section,
there would be a `models` folder as well in our `gen` folder,
containing the generated models for those definitions.
But as we don't have any definitions to share, you won't find it in this tutorial.

For more information, read through the generated code.
It might help to keep the `swagger/swagger.yml` definition next to you,
to help you realize what is defined, for what and where.

So now that we have the generated server, it is time to write our actual main file.
In it we'll parse some simple flags that can be used to configure the server,
we'll setup the API and finally start the server. All in all, very minimal and simple.

so let's start writing the `./cmd/greeter/main.go` file:

We start by defining our flags, in our example using the standard `flag` pkg:

```go
var portFlag = flag.Int("port", 3000, "Port to run this service on")
```

Now it's time to write our main logic,
starting by loading our embedded swagger spec.
This is required, as it is used to configure the dynamic server,
the core of our generated server (found under the `./gen/restapi` dir).

```go
swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
if err != nil {
	log.Fatalln(err)
}
```

With the spec loaded in, we can create our API and server:

```go
api := operations.NewGreeterAPI(swaggerSpec)
server := restapi.NewServer(api)
defer server.Shutdown()
```

With the server created, we can overwrite the default port, using our `portFlag`:

```go
flag.Parse()
server.Port = *portFlag
```

After that, we can serve our API and finish our main logic:

```go
if err := server.Serve(); err != nil {
	log.Fatalln(err)
}
```

Putting that all together, we have the following main function:

```go
func main() {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	api := operations.NewGreeterAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// parse flags
	flag.Parse()
	// set the port this service will be run on
	server.Port = *portFlag

	// TODO: Set Handle

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
```

Now that we have our server defined, let's give it a first spin!
You can run it from our root directory using the following command:

```
$ go run ./cmd/greeter/main.go --port 3000
```

Let's now try to call our only defined operation,
using the [httpie][] cli:

```
$ http get :3000/hello
```

Sadly this gives us the following output:

```http
HTTP/1.1 501 Not Implemented
Connection: close
Content-Length: 50
Content-Type: text/plain
Date: Thu, 26 Jan 2017 13:09:52 GMT

operation GetGreeting has not yet been implemented
```

The good news is that our OpenAPI-based Golang service is working.
The bad news is that we haven't implemented our handlers yet.
We'll need one handler per operation, that does the actual logic.

You might wonder why it does give us a sane response, rather then panicking.
After grepping for that error message, or using a recursive search in your favorite editor, you'll find that this error originates from the greeter API constructor (`NewGreeterAPI`) found in `./gen/restapi/operations/greeter_api.go`.
Here we'll see that all our consumers, producers and handlers have sane defaults.

So now that we know that we just got our ass saved by go-swagger,
let's actually start working towards implementing our handlers.

Inspecting the `gen/restapi/operations/get_greeting.go` file,
we'll find the following snippet:

```go
// GetGreetingHandlerFunc turns a function with the right signature into a get greeting handler
type GetGreetingHandlerFunc func(GetGreetingParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetGreetingHandlerFunc) Handle(params GetGreetingParams) middleware.Responder {
	return fn(params)
}
```

Here we can read that there is a function type `GetGreetingHandlerFunc`, defined for the `getGreeting` operation which takes in one parameter of type `GetGreetingParams` and returns a `middleware.Responder`. This is the type alias our Handler has to adhere to.

A bit more down we'll also encounter an interface `GetGreetingHandler`, defined for the `getGreeting` operation:

```go
// GetGreetingHandler interface for that can handle valid get greeting params
type GetGreetingHandler interface {
	Handle(GetGreetingParams) middleware.Responder
}
```

Its only defined method looks very similar to the function type `GetGreetingHandlerFunc` function defined above. This is no coincidence.

Even better, the `GetGreetingHandlerFunc` implements the `Handle` function as defined by the `GetGreetingHandler` interface, meaning that we can use a function respecting the type alias as defined by `GetGreetingHandlerFunc`, where we normally would have to implement a struct adhering to the `GetGreetingHandler` interface.

Implementing our handler as a struct allows us to handle with a certain state in mind. You can check out the [kvstore example][] to see a more elaborate example, where you can see the handlers being implemented using a struct per operation.

Our Greeter API is however simple enough, that we'll opt for just a simple method.
[KISS][] never grows old. So with all of this said, let's implement our one and only handler.

Back to the `./cmd/greeter/main.go` file, we'll define our handler as follows:

```go
api.GetGreetingHandler = operations.GetGreetingHandlerFunc(
	func(params operations.GetGreetingParams) middleware.Responder {
		name := swag.StringValue(params.Name)
		if name == "" {
			name = "World"
		}

		greeting := fmt.Sprintf("Hello, %s!", name)
		return operations.NewGetGreetingOK().WithPayload(greeting)
	})
```

Which replaces the `TODO: Set Handle` comment, originally defined.

Let's break down the snippet above.
First we make use of the [go-openapi/swag][] package, which is full of Goodies.
In this case we use the `StringValue` function which transforms a `*string` into a `string`. The result will be empty in case it was nil. This comes in handy as we know that our parameter can be nil when not given, as it is _not_ required. Finally we form our greeting and return it as our payload with our `200 OK` response.

Let's run our server:

```
$ go run ./cmd/greeter/main.go --port 3000
```

And now we're ready to test our greeter API once again:

```
$ http get :3000/hello
```

```http
HTTP/1.1 200 OK
Connection: close
Content-Length: 13
Content-Type: text/plain
Date: Thu, 26 Jan 2017 13:47:49 GMT

Hello, World!
```

Hurray, let's now greet _Swagger_:

```
$ http get :3000/hello name==Swagger
```

```http
HTTP/1.1 200 OK
Connection: close
Content-Length: 15
Content-Type: text/plain
Date: Thu, 26 Jan 2017 13:48:40 GMT

Hello, Swagger!
```

Great, Swagger will be happy to hear that.

As we just learned, using [go-swagger][] and a manually defined [OpenAPI2.0][] specification file, we can build a Golang service with minimal effort. Please read the other [go-swagger][] docs for more information about how to use it and its different elements.

Also please checkout the [kvstore example][] for a more complex example.
It is the main inspiration for this tutorial and has been built using the exact
same techniques as described in this tutorial.

[OpenAPI2.0]: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md
[httpie]:https://httpie.org
[kvstore example]: https://github.com/go-openapi/kvstore
[KISS]: https://en.wikipedia.org/wiki/KISS_principle
[go-openapi/swag]: https://github.com/go-openapi/swag
[go-swagger]: https://github.com/go-swagger/go-swagger
