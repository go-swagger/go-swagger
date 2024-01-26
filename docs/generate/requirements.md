---
title: Build requirements
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 10
---
# Requirements to build generated code

## First time with golang?

Golang is a powerful and enticing language, but it may sometimes confuse first time users.

Before engaging further with `go-swagger`, please take a while to get comfortable with golang basics 
and conventions. That will save yourself much time and frustration.

## About dependencies

Before generating code, you should make sure your target is going to properly resolve dependencies.

> **NOTE**: generation uses extensively the `goimports` tool and dependencies must be matched at code generation time.

If you have built `go-swagger` locally (e.g. from source), all dependencies are already installed.

The following required dependencies may be fetched by using `go get`:

- [`github.com/go-openapi/errors`](https://github.com/go-openapi/errors)
- [`github.com/go-openapi/loads`](https://github.com/go-openapi/loads)
- [`github.com/go-openapi/runtime`](https://github.com/go-openapi/runtime)
- [`github.com/go-openapi/spec`](https://github.com/go-openapi/spec)
- [`github.com/go-openapi/strfmt`](https://github.com/go-openapi/strfmt)
- [`github.com/go-openapi/swag`](https://github.com/go-openapi/swag)
- [`github.com/go-openapi/validate`](https://github.com/go-openapi/validate)

{{< hint "info" >}}
The code generation process ends with a message indicating the packages required for your generated code.
{{< /hint >}}


### What are the dependencies required by the generated server?

There are some additional packages required by the (default) generated server.
This depends on your generation options, a command line flags handling package:

- [`github.com/jessevdk/go-flags`](https://www.github.com/jessevdk/go-flags), or
- [`github.com/spf13/pflag`](https://www.github.com/spf13/pflag)

### What are the dependencies required by the generated client?

Same as above.

### What are the dependencies required by the generated models?

The generated models only depend on:

- [`github.com/go-openapi/errors`](https://www.github.com/go-openapi/errors)
- [`github.com/go-openapi/strfmt`](https://www.github.com/go-openapi/strfmt)
- [`github.com/go-openapi/swag`](https://www.github.com/go-openapi/swag)
- [`github.com/go-openapi/validate`](https://www.github.com/go-openapi/validate)

## How about generating specs?

The code that is scanned for spec generation _must_ resolve all its dependencies (i.e. must build).
