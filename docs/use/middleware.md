# BYO middleware

Go-swagger chose the golang `net/http` package as base abstraction. That means that for _any_ supported transport by the toolkit you can reuse _any_ middleware existing middlewares that following the stdlib middleware pattern.

<!--more-->

There are several projects providing middleware libraries for weaving all kinds of functionality into your request handling. None of those things are the job of go-swagger, go-swagger just serves your specs.

The server takes care of a number of things when a request arrives:

* routing
* authentication
* input validation
* content negotiation
* parameter and body binding

If you're unfamiliar with the concept of golang net/http middlewares you can read up on it here:  
[Making and Using HTTP Middleware](http://www.alexedwards.net/blog/making-and-using-middleware)

Besides serving the swagger specification as an API, the toolkit also serves the actual swagger specification document.
The convention is to use the `/swagger.json` location for serving up the specification document, so we serve the
specification at that path.

### Add middleware

The generated server allows for 2 extension points to inject middleware in its middleware chain. These have to do with
the lifecycle of a request. You can find those hooks in the configure_xxx_api.go file.

The first one is to add middleware all the way to the top of the middleware stack. To do this you add them in the
`setupGlobalMiddleware` method. This middleware applies to everything in the go-swagger managed API.

```go
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
```

The second extension point allows for middleware to be injected right before actually handling a matched request.
This excludes the swagger.json document from being affected by this middleware though.  This extension point makes the
middlewares execute right after routing but right before authentication, binding and validation.  You add middlewares
to this point by editing the `setupMiddlewares` method in configure_xxx_api.go

```go
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}
```

The global middleware is an excellent place to do things like panic handling, request logging or adding metrics.  While
the plain middleware allows you to kind of filter this by request path without having to take care of routing. You also
get access to the full context that the go-swagger toolkit uses throughout the lifecycle of a request.

#### Add logging and panic handling

A very common requirement for HTTP APIs is to include some form of logging. Another one is to handle panics from your
API requests.  The example for a possible implementation of this uses [this community provided
middleware](https://github.com/dre1080/recover) to catch panics.

```go
func setupGlobalMiddleware(handler http.Handler) http.Handler {
  recovery := recover.New(&recover.Options{
    Log: log.Print,
  })
  return recovery(handler)
}
```

There are tons of middlewares out there, some are framework specific and some frameworks don't really use the plain
vanilla golang net/http as base abstraction. For those you can use a project like [interpose](https://github.com/carbocation/interpose) that serves as an adapter
layer so you can still reuse middlewares. Of course nobody is stopping you to just implement your own middlewares.

For example using interpose to integrate with [logrus](https://github.com/carbocation/interpose/blob/master/middleware/negronilogrus.go).

```go
import (
  interpose "github.com/carbocation/interpose/middleware"
)
func setupGlobalMiddleware(handler http.Handler) http.Handler {
  logViaLogrus := interpose.NegroniLogrus()
  return logViaLogrus(handler)
}
```

And you can compose these middlewares into a stack using functions.

```go
func setupGlobalMiddleware(handler http.Handler) http.Handler {
  handlePanic := recover.New(&recover.Options{
    Log: log.Print,
  })

  logViaLogrus := interpose.NegroniLogrus()

  return handlePanic(
    logViaLogrus(
      handler
    )
  )
}
```

#### Add rate limiting

You can also add rate limiting in a similar way. Let's say we just want to rate limit the valid requests to our swagger
API. To do so we could use [tollbooth](https://github.com/didip/tollbooth).

```go
func setupMiddlewares(handler http.Handler) http.Handler {
  limiter := tollbooth.NewLimiter(1, time.Second)
  limiter.IPLookups = []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}
	return tollbooth.LimitFuncHandler(handler)
}
```

And with this you've added rate limiting to your application.
