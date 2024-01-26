---
title: Generating a swagger spec
date: 2023-01-01T01:01:01-08:00
draft: true
weight: 10
---
# Spec generation

The toolkit has a command that will let you generate a swagger spec document from your code.
The command integrates with go doc comments, and makes use of structs when it needs to know of
types.

Based on the work from https://github.com/yvasiyarov/swagger.

It uses a similar approach but with expanded annotations and it produces a swagger 2.0 spec.

The goal of the syntax is to make it look as a natural part of the documentation for the application code.

The generator is passed a main package and it uses that to discover all the code in use.
To do this it makes use of go's loader package. The same package that is used by tools like goimports to discover which files to format.

Once the parser has encountered a comment that matches one of its known tags, the parser will assume that the rest of the comment block is for swagger.

If it is your first time using this library for generating swagger 
specification, you can also take a look at 
[this getting started guide](https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9?source=friends_link&sk=b402acc563e8d2bfadd1ac02abddc3bb)
on Medium to get the big picture and then return here to read more 
about all annotations and commands provided by this package.

## Usage

To generate a swagger spec:

```sh
swagger generate spec -o ./swagger.json
```

You give it a main file and it will parse all the files that are reachable by that main
package to produce a swagger specification.

You may add a go:generate comment to your main file for example:

```sh
//go:generate swagger generate spec -o swagger.json
```

The command requires a main package or file and it wants your code to compile. It uses the go tools loader to load an application and then scans all the packages that are in use by the code base.
This means that for something to be discoverable it needs to be reachable by a code path triggered through the main package.

If an annotation is not yet supported or you want to merge with a pre-existing spec, you can use the -i parameter.

```sh
swagger generate spec -i ./swagger.yml -o ./swagger.json
```

The idea is that there are certain things that are more easily expressed by just using yaml

To generate spec in yaml format, just name the output file with ".yml" or ".yaml" extension. For example:

```sh
swagger generate spec -o ./swagger.yml
```

If you don't want to generate Go language specific extensions in the spec file, you can disable them by doing

```sh
SWAGGER_GENERATE_EXTENSION=false && swagger generate spec -o ./swagger.yml
```

## Parsing rules

{{< hint warning >}}
This command relies heavily on the way godoc works.

This means you should be very aware of all the things godoc supports.
{{< /hint >}}


* [godoc documenting go code](http://blog.golang.org/godoc-documenting-go-code)
* [godoc ToHTML](https://golang.org/pkg/go/doc/#ToHTML)
* [commenting go effectively](https://golang.org/doc/effective_go.html#commentary)
* [godoc documentation](https://godoc.org/golang.org/x/tools/cmd/godoc)

Single page which documents all the currently supported godoc rules:

* [godoc tricks](https://godoc.org/github.com/fluhus/godoc-tricks)

The generated code tries to avoid golint errors.

* [go lint](https://github.com/golang/lint)
* [go lint style guide](https://github.com/golang/go/wiki/CodeReviewComments)

When an object has a title and a description field, it will use the go rules to parse those. So the first line of the
comment block will become the title, or a header when rendered as godoc. The rest of the comment block will be treated
as description up to either the end of the comment block, or a line that starts with a known annotation.

#### Annotation syntax

If you want to exclude something from the spec generation process you can try with the struct tag: `json:"-"`

There are several annotations that mark a comment block as a participant for the swagger spec.

- [swagger:meta](../reference/annotations/meta.md)
- [swagger:route](../reference/annotations/route.md)
- [swagger:parameters](../reference/annotations/params.md)
- [swagger:response](../reference/annotations/response.md)
- [swagger:operation](../reference/annotations/operation.md)
- [swagger:model](../reference/annotations/model.md)
- [swagger:allOf](../reference/annotations/allOf.md)
- [swagger:strfmt](../reference/annotations/strfmt.md)
- [swagger:discriminated](../reference/annotations/discriminated.md)
- [swagger:ignore](../reference/annotations/ignore.md)

#### Embedded types

For the embedded schemas there are a set of rules for the spec generator to vary the definition it generates.
When an embedded type isn't decorated with the `swagger:allOf` annotation, then the properties from the embedded value will be included in the generated definition as if they were defined on the definition. But when the embedded type is decorated with the `swagger:allOf` annotation then the all of element will be defined as a "$ref" property instead. For an annotated type there is also the possibility to specify an argument, the value of this argument will be used as the value for the `x-class` extension. This allows for generators that support the
`x-class` extension to reliably build a serializer for a type with a discriminator

#### Known vendor extensions

There are a couple of commonly used vendor extensions that most frameworks support to add functionality to the swagger spec.

For generating a swagger specification document this toolkit supports:

Vendor extension | Description
