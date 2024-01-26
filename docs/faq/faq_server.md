---
title: About generating a server
date: 2023-01-01T01:01:01-08:00
draft: true
---
<!-- Questions about server generation -->
## Server generation and customization

### What are the dependencies required by the generated server?
The code generation process ends with a message indicating package dependencies for your generated code.

Basically, here are the required packages:
- [`github.com/go-openapi/errors`](https://www.github.com/go-openapi/errors)
- [`github.com/go-openapi/loads`](https://www.github.com/go-openapi/loads)
- [`github.com/go-openapi/runtime`](https://www.github.com/go-openapi/runtime)
- [`github.com/go-openapi/spec`](https://www.github.com/go-openapi/spec)
- [`github.com/go-openapi/strfmt`](https://www.github.com/go-openapi/strfmt)
- [`github.com/go-openapi/swag`](https://www.github.com/go-openapi/swag)
- [`github.com/go-openapi/validate`](https://www.github.com/go-openapi/validate)

And depending on your generation options, a command line flags handling package:
- [`github.com/jessevdk/go-flags`](https://www.github.com/jessevdk/go-flags), or
- [`github.com/spf13/pflag`](https://www.github.com/spf13/pflag)
- `flag`

This dependency used to be necessary up to release 0.14:
- [`github.com/tylerb/graceful`](https://www.github.com/tylerb/graceful)

These packages may of course be *vendored* with your own source.

Originally from issue [#1085](https://github.com/go-swagger/go-swagger/issues/1085).

### How to add custom flags?
`go-swagger` ships with an option to select a flag management package: `swagger generate server --flag-strategy=[go-flags|pflag|flag]`

You may of course customize your server to accept arbitrary flags the way you prefer.
This should be normally done with the generated main.go. For customization, you may either skip the generation of the main package (`--skip-main`)
and provide your own, or customize template generation to generate a custom main.

Here's an example: [kv store example](https://github.com/go-openapi/kvstore/blob/master/cmd/kvstored/main.go#L50-L57)

Originally from issue [#1036](https://github.com/go-swagger/go-swagger/issues/1036).

### How do you integrate the flag sets of go-swagger and other packages, in particular, glog?
_Use-case_: logger
>I am trying to integrate a package into a go-swagger generated API that is using the `github.com/golang/glog` logger.
>When I initialize the glog logger it appears to shield the flags defined in the go-swagger runtime.

**Answer**: the generated API has a Logger property that is a function with signature: `func(string, ...interface{})`

You can configure it with any logger that exposes the signature.

eg.: https://github.com/go-swagger/go-swagger/blob/master/examples/authentication/restapi/configure_auth_sample.go#L33

_Use-case_: logger flags
>Still having a problem with how and where to initialize glog so that both sets of flags are honored:
>the runtime flags, such as `--tls-certificate` and the glog flags like `-log_dir` and `-stderrthreshold`.

>If I initialize glog in the config_xxx.go I don't get the go-swagger runtime flags, and if I initialize glog in the engine, I get the error: `logging before flag.Parse`.
>I realize that this question is not so much about logging *per se*, but more about how to merge the flag sets defined by different packages.

**Answer**: you can generate a server with `--flag-strategy pflag`

After that you can use its integration to add goflags. You would do this in the main file.
Subsequently it's probably a good idea to generate code with `--exclude-main` so the update is preserved.

See also: https://github.com/spf13/pflag#supporting-go-flags-when-using-pflag

Example:
```golang
import (
    goflag "flag"
    flag "github.com/spf13/pflag"
)

var ip *int = flag.Int("flagname", 1234, "help message for flagname")

func main() {
    flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
    flag.Parse()
}
```

Originally from issue [#762](https://github.com/go-swagger/go-swagger/issues/762).

### How to serve two or more swagger specs from one server?
_Use-case_: I want a go-swagger generated server to serve 2 swagger specs that have no overlap on paths.
One is a generic metadata api that is implemented by many services and the other is a
service-specific api. The built-in server.go template Server struct, by having exactly
one API & handler, appears to prevent this.

**Suggestions**:
1. Use go-swagger `mixin` command to merge specs into a single one
2. Create yourself a top-level swagger file that just includes the two lower-level ones (using `$ref`).
You may use go-swagger `flatten` to flatten the resulting spec
3. You can also make your own main function and use the code from the generation
of both spec (with `--skip-main`).
This allows for customization like using a different middleware stack, which in turn gives you
the ability to serve 2 swagger specs at different paths.

Originally from issue [#1005](https://github.com/go-swagger/go-swagger/issues/1005). *(comes with a sample main.go for spec composition)*.

### How to access access API struct inside operator handler?
_Use-Case_:
my question is how can I access the generated API interface from within an operation handler function ?
Can i pass it somehow via context or any other way to do that?

**Answer**: **No**. It's not reachable from within the handler anywhere.


>I created a module like apache access module ACL based on IP address for different URLs.
>Instead of URL for lookup key I decided to use Operation.ID.
>Lookup would be faster in that way because each operation has a unique id according to swagger specification.
>The problem comes when I want to check against that ACL of mine...

**Suggestions**:
This is possible in 2 ways.
- first way is by using an authenticator,
- the second way is making a middleware (not global)

Example with Authenticator:
```golang
// Authenticator represents an authentication strategy
// implementations of Authenticator know how to authenticate the
// request data and translate that into a valid principal object or an error
type Authenticator interface {
    Authenticate(interface{}) (bool, interface{}, error)
}
```
You may see the schemes currently supported here: https://github.com/go-openapi/runtime/tree/master/security

Example with Middleware:
```golang
// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
    return handler
}
```
to get to the matched route in one of those you can:
```golang
import "github.com/gorilla/context"
context.Get(3, request)
```
That gets you a matched route.

>Note: it might be worth it to expose `Context` as an exported method on the API builder.
>That would be a simple PR to add the following code in go-swagger
```golang
func (o *LifecycleManagerAPI) Context() *middleware.Context {
    if o.context == nil {
        o.context = middleware.NewRoutableContext(o.spec, o, nil)
    }
    return o.context
}
```
then your middleware could take a context in and use `RouteInfo(request)` like this one: https://github.com/go-openapi/runtime/blob/master/middleware/security.go.

_Similar use-case_:
I have some domain objects I need access to in the handlers (db connection, telemetry client, etc).

*What is the recommended way/place to define these and access them in the handlers?*

>Would I define them in configure_xxx() and make wrapper functions for the handlers to make them accessible?
>Or is there a nice way to add them to the context?
>I was looking for some examples of this but couldn't find anything.

**Hint**: look at this one: https://github.com/go-openapi/kvstore

The important takeaway is that main and the handlers have to be pulled outside of the generated code since only configure_xxx() is protected.
And main() doesn't call configureAPI() ...
that was a little confusing looking through the other examples and not seeing any changes to the vanilla config code.

_Similar use-case_: dependency injection

Wouldn't it be better to have all the handlers automatically be part of a default struct that simply has a Context member variable or empty interface?

>That would save everyone a lot of copy/paste when we need to inject some info.
>I mean, a different context than the one available on params.HTTPRequest.Context(),
>more like an application level context, e.g. something I can stuff a database reference or other business state into.

**Hint**: remember that the generated packages are made of a series of files that you can extend and tailor to your own needs by adding more files.

For instance, a new (non generated) file in the operations package could bring to life the applicative context/dependencies injection and be called from the configure_xxx.
You may alternatively modify the generation by providing your own templates, and possibly extend the interface of the Server struct.

Originally from issue [#661](https://github.com/go-swagger/go-swagger/issues/661).

### Use go-swagger to generate different client or servers
_Use-Case_:
I am using go-swagger to generate some part of a server application automatically
and I would love to reuse our code by transforming the code in go templates.

>It would be  enough to export the `appGenerator` type and have a function that returns it (maybe generator.GenerateServer itself?).
>I would then use `appGenerator` to execute the templates.
>How could I realize this? Is it possible with go-swagger?

**Answer**: you can provide your own templates for go-swagger.

The client and server generators allow you to specify a directory on disk to add custom templates.

Here are some docs: http://goswagger.io/generate/templates/

>In VIC they do this: https://github.com/vmware/vic/tree/master/lib/apiservers/templates
>https://github.com/vmware/vic/blob/master/Makefile#L274-L281

**Hint**: you can also override templates by using the same names:
https://github.com/go-swagger/go-swagger/blob/master/generator/templates.go#L61-L73

*Wouldn't this generate roughly the same structure of the server?*

>I don't want to change minor details, I want to have code that looks totally different
>(but only for the server part, models and clients are more than okay) while using code the parsing and validation from go-swagger.
>This means different number of files and different functionalities.

**Answer**: yes, it would generate a similar structure.

Note: customizing templates already brings many options to the table, including generating artifacts in other languages than go.

There is some documentation on the config file format here: https://github.com/go-swagger/go-swagger/blob/gen-layout-configfile/docs/use/template_layout.md

Also keep in mind that `go-openapi` and `go-swagger` constitute a _toolkit_
and provide you the *tools* to adapt to your own use case.
The `swagger` command line and standard templates only covers general purpose use-cases.

If you think your use-case would benefit to many people, feel free to make the necessary changes for your case to work and submitting a PR.

Example config file for generation:
```YAML
layout:
  application:
    - name: configure
      source: asset:serverConfigureapi
      target: "{{ joinFilePath .Target .ServerPackage }}"
      file_name: "{{ .Name }}_client.go"
      skip_exists: true
    - name: main
      source: asset:serverMain
      target: "{{ joinFilePath .Target \"cmd\" (dasherize (pascalize .Name)) }}-server"
      file_name: "main.go"
    - name: embedded_spec
      source: asset:swaggerJsonEmbed
      target: "{{ joinFilePath .Target .ServerPackage }}"
      file_name: "embedded_spec.go"
    - name: server
      source: asset:serverServer
      target: "{{ joinFilePath .Target .ServerPackage }}"
      file_name: "server.go"
    - name: builder
      source: asset:serverBuilder
      target: "{{ joinFilePath .Target .ServerPackage .Package }}"
      file_name: "{{ snakize (pascalize .Name) }}_api.go"
    - name: doc
      source: asset:serverDoc
      target: "{{ joinFilePath .Target .ServerPackage }}"
      file_name: "doc.go"
  models:
   - name: definition
     source: asset:model
     target: "{{ joinFilePath .Target .ModelPackage }}"
     file_name: "{{ (snakize (pascalize .Name)) }}.go"
  operations:
   - name: parameters
     source: asset:serverParameter
     target: "{{ joinFilePath .Target .ServerPackage .APIPackage .Package }}"
     file_name: "{{ (snakize (pascalize .Name)) }}_parameters.go"
   - name: responses
     source: asset:serverResponses
     target: "{{ joinFilePath .Target .ServerPackage .APIPackage .Package }}"
     file_name: "{{ (snakize (pascalize .Name)) }}_responses.go"
   - name: handler
     source: asset:serverOperation
     target: "{{ joinFilePath .Target .ServerPackage .APIPackage .Package }}"
     file_name: "{{ (snakize (pascalize .Name)) }}.go"
```

### Support streaming responses
_Use-Case_: Docker client expects a stream of JSON structs from daemon to show a progress bar, as in:
```bash
curl -i -X POST http://IP:PORT/images/create?fromImage=alpine&tag=latest
```
*How could I write a server providing a streaming response?*

**Answer**: Operations in the generated server are expected to return a responder.

This interface is defined as:
```golang
// Responder is an interface for types to implement
// when they want to be considered for writing HTTP responses
type Responder interface {
    WriteResponse(http.ResponseWriter, httpkit.Producer)
}
```

With the `middleware.ResponderFunc` helper construct, you can just write a `func(http.ResponseWriter, httpkit.Producer)`
where you want a streaming response.

This should be sufficient. However:

>I've toyed with a channel based stream where you send struct objects to a channel, which then gets streamed to the browser.
>I decided against this because it seemed to just add complexity for little benefit.
>I can be persuaded to implement such a responder though, and should somebody send a PR like that I would not say no to it.

Originally from issue [#305](https://github.com/go-swagger/go-swagger/issues/305).

### OAuth authentication does not redirect to the authorization server
_Use-Case_: oauth2 accessCode flow does not redirect to the authorization server

> In my understanding, if the accessCode flow is used for oauth2 securitydefinition, the generated server could redirect the authentication to the oauth2 server, e.g., https://accounts.google.com/o/oauth2/v2/auth. However, my generated code does not perform the redirection. Could anyone help on this? Thanks.

Like in:
```yaml
```

Back to [all contributions](/faq)
