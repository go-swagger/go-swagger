+++
categories = ["generate"]
date = "2015-11-09T18:58:47-08:00"
title = "swagger.json generation"
weight = 1
series = ["home"]
+++

The toolkit has a command that will let you generate a swagger spec document from your code. 
The command integrates with go doc comments, and makes use of structs when it needs to know of
types.

<!--more-->

Based on the work from https://github.com/yvasiyarov/swagger  
It uses a similar approach but with expanded annotations and it produces a swagger 2.0 spec.

The goal of the syntax is to make it look as a natural part of the documentation for the application code.

The generator is passed a base path (defaults to current) and tries to extract a go package path from that.
Once it has a go package path it will scan the package recursively, skipping the Godeps, files ending in test.go and
directories that start with an underscore, it also skips file system entries that start with a dot.

Once the parser has encountered a comment that matches on of its known tags, the parser will assume that the rest of
the comment block is for swagger.

### Usage

To generate a spec:

```
swagger generate spec -o ./swagger.json
```

You give it a main file and it will parse all the files that are required by that main
package to produce a swagger specification.

To use you can add a go:generate comment to your main file for example:

```
//go:generate swagger generate spec
```

The command requires a main package or file and it wants your code to compile. It uses the go tools loader to load an
application and then scan all the packages that are in use by the code base.
This means that for something to be discoverable it needs to be reachable by a codepath triggered through the main
package.

#### Parsing rules

This command relies heavily on the way godoc works! This means you should be very aware of all the things godoc
supports.

* [godoc documentation](https://godoc.org/golang.org/x/tools/cmd/godoc)
* [godoc documenting go code](http://blog.golang.org/godoc-documenting-go-code)
* [godoc ToHTML](https://golang.org/pkg/go/doc/#ToHTML)

Single page which documents all the currently supported godoc rules:

* [godoc tricks](https://godoc.org/github.com/fluhus/godoc-tricks)


#### Annotation syntax

If you want to exclude something from the spec generation process you can try with the struct tag: `json:"-"`

There are several annotations that mark a comment block as a participant for the swagger spec.

* [swagger:meta](meta) 
* [swagger:route](route)
* [swagger:params](params)
* [swagger:response](response)
* [swagger:model](model)
* [swagger:allOf](allOf)
* [swagger:strfmt](strfmt)
