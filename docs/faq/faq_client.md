---
title: About generating a client
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 20
---
<!-- Questions about client generation -->
## Client generation

### Example for dynamic client?
_Use-case_: I have the swagger spec file for an existing 3rd party REST service for their application.
A dynamic client would allow me to load the swagger spec file and provide the ability to formulate requests and parse responses
based on the loaded spec file.

For e.g. there are REST client packages for go, like go-resty, that provide a nice interface
to interact with a REST server, but go-resty doesn't understand swagger specs.

- *Just like the untyped dynamic server example in go-swagger, is there an example for a dynamic client?*
- *Can a REST client be created at runtime by loading a swagger spec file, without going through a code generation and compilation?*
- *Can I get the REST client functionality without first generating a client using go-swagger and then*
compiling it back into the code.
- *Can this be done dynamically like the dynamic server example?*

**Answer**: **you can't currently**.

You'd have to do by hand everything the code generator does for you. Every time the API changes, you would have to do this again.
The use case for the server side is covered but not the client side.

That being said, you might find this test useful as an example:
[runtime test](https://github.com/go-openapi/runtime/blob/master/client/runtime_test.go#L144-L188)

Originally from issue [#996](https://github.com/go-swagger/go-swagger/issues/996).

### Can we set a User-Agent header?
_Use-Case_: we would like to be able to set an arbitrary user-agent header either at client generation time or at compile time.

*Is it possible to do this?*

>The Swagger specification is irrelevant in this case because we are using the same specification to generate many clients.

**Answer**: here is the outline of how to achieve that.

- You can use a custom transport which allows you to set the user agent.
https://github.com/go-openapi/runtime/blob/master/client/runtime.go#L132
- Then you can configure it with this constructor method
https://github.com/go-swagger/go-swagger/blob/master/examples/todo-list/client/todo_list_client.go#L52
- You can also configure that runtime with a `stdlib http.Client`
https://github.com/go-openapi/runtime/blob/master/client/runtime.go#L167
- You can extend intercept a http request with the `http.RoundTripper interface`. https://godoc.org/net/http#RoundTripper
which you can set here: https://github.com/go-openapi/runtime/blob/master/client/runtime.go#L116
- so for the client here:

```golang
var myRoundTripper http.RoundTripper = createRoundTripper()
transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
transport.Transport = myRoundTripper
todoListClient := New(transport, nil)
```

_Other use-Case_: can the same pattern of using an `http.RoundTripper` be used to implement the AWS Signature v4 which requires
reading and modifying the `*http.Request` before its sent?

**Answer**: **yes it can**.

>The roundtripper is the last thing executed before sending the request on the wire.

See also issue [#935](https://github.com/go-swagger/go-swagger/issues/935).

Originally from issue [#911](https://github.com/go-swagger/go-swagger/issues/911).

