# Requirements to build generated code

### First time with golang?

Golang is a powerful and enticing language, but it may sometimes confuse first timers.

Before engaging further with `go-swagger`, please take a while to get comfortable with golang basics 
and conventions. That will save yourself much time and frustration.

### Standard golang environment

* version: we support the two latest versions of the go compiler
* `GOPATH` environment variable set: all sources reside under `$GOPATH/src`
* it is recommended, but not mandatory, to use the `dep` tool to manage dependencies
(see [here](https://golang.github.io/dep/docs/introduction.html))

The target directory for your generated code _must_ be under GOPATH/src.

## Getting dependencies

Before generating code, you should make sure your target is going to properly resolve dependencies.

> **NOTE**: generation makes use of the `goimports` tool and dependencies must be matched at code generation time.

If your target is located under the `go-swagger` install directory (when installed from source), dependencies are directly
provided by the `vendor` directory that ships with `go-swagger`.

The following required dependencies may be fetched by using `go get`:

- [`github.com/go-openapi/errors`](https://www.github.com/go-openapi/errors)
- [`github.com/go-openapi/loads`](https://www.github.com/go-openapi/loads)
- [`github.com/go-openapi/runtime`](https://www.github.com/go-openapi/runtime)
- [`github.com/go-openapi/spec`](https://www.github.com/go-openapi/spec)
- [`github.com/go-openapi/strfmt`](https://www.github.com/go-openapi/strfmt)
- [`github.com/go-openapi/swag`](https://www.github.com/go-openapi/swag)
- [`github.com/go-openapi/validate`](https://www.github.com/go-openapi/validate)

You may also build a vendor directory in your planned target: a way to achieve that is to copy there an example from the
`go-swagger/examples` repository then run `dep` - see [how to use dep here](https://github.com/golang/dep).
This will produce `Gopkg.toml` and `Gopkg.lock` files and construct a vendor directory with all required dependencies
(the ones above and all transitive dependencies). Another way is to proceed in two steps, first with `go get`, then generate code, 
then build the vendor tree with `dep`.

> **NOTE** : the code generation process ends with a message indicating the packages required for your generated code.


### What are the dependencies required by the generated server?

Additional packages required by the (default) generated server
depend on your generation options, a command line flags handling package:

- [`github.com/jessevdk/go-flags`](https://www.github.com/jessevdk/go-flags), or
- [`github.com/spf13/pflags`](https://www.github.com/spf13/pflags)

### What are the dependencies required by the generated client?

Same as above, plus:

- `golang.org/x/net/context`

### What are the dependencies required by the generated models?

The generated models package depends only on:

- [`github.com/go-openapi/errors`](https://www.github.com/go-openapi/errors)
- [`github.com/go-openapi/strfmt`](https://www.github.com/go-openapi/strfmt)
- [`github.com/go-openapi/swag`](https://www.github.com/go-openapi/swag)
- [`github.com/go-openapi/validate`](https://www.github.com/go-openapi/validate)

### How about generating specs?

The code that is scanned for spec generation _must_ resolve all its dependencies (i.e. must build).
