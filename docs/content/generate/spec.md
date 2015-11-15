+++
date = "2015-11-09T18:58:47-08:00"
title = "Spec generation reference"
+++

The toolkit has a command that will let you generate a swagger spec document from your code. 
The command integrates with go doc comments, and makes use of structs when it needs to know of
types.

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

#### Annotation syntax

There are several annotations that mark a comment block as a participant for the swagger spec.

* [swagger:meta](meta) 
* [swagger:strfmt](strfmt)
* [swagger:model](model)
* [swagger:route](route)
* [swagger:params](params)
* [swagger:response](response)
* [swagger:allOf](allOf)
