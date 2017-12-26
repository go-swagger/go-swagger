# swagger.json generation

The toolkit has a command that will let you generate a swagger spec document from your code.
The command integrates with go doc comments, and makes use of structs when it needs to know of
types.

<!--more-->

Based on the work from https://github.com/yvasiyarov/swagger
It uses a similar approach but with expanded annotations and it produces a swagger 2.0 spec.

The goal of the syntax is to make it look as a natural part of the documentation for the application code.

The generator is passed a main package and it uses that to discover all the code in use.
To do this it makes use of go's loader package. The same package that is used by tools like goimports to discover which files to format.

Once the parser has encountered a comment that matches on of its known tags, the parser will assume that the rest of the comment block is for swagger.

### Usage

To generate a spec:

```
swagger generate spec -o ./swagger.json
```

You give it a main file and it will parse all the files that are reachable by that main
package to produce a swagger specification.

To use you can add a go:generate comment to your main file for example:

```
//go:generate swagger generate spec
```

The command requires a main package or file and it wants your code to compile. It uses the go tools loader to load an application and then scans all the packages that are in use by the code base.
This means that for something to be discoverable it needs to be reachable by a code path triggered through the main package.

If an annotation is not yet supported or you want to merge with a pre-existing spec, you can use the -i parameter.

```
swagger generate spec -i ./swagger.yml -o ./swagger.json
```

The idea is that there are certain things that are more easily expressed by just using yaml, to

#### Parsing rules

<img class="emoji" title=":warning" alt=":warning" src="https://assets-cdn.github.com/images/icons/emoji/unicode/26a0.png?v5" width="20" height="20" align="absmiddle" />This command relies heavily on the way godoc works. <img class="emoji" title=":warning" alt=":warning" src="https://assets-cdn.github.com/images/icons/emoji/unicode/26a0.png?v5" width="20" height="20" align="absmiddle" />

This means you should be very aware of all the things godoc supports.

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

* [swagger:meta](spec/meta.md)
* [swagger:route](spec/route.md)
* [swagger:parameters](spec/params.md)
* [swagger:response](spec/response.md)
* [swagger:model](spec/model.md)
* [swagger:allOf](spec/allOf.md)
* [swagger:strfmt](spec/strfmt.md)
* [swagger:discriminated](spec/discriminated.md)
* [swagger:ignore](spec/ignore.md)

#### Embedded types

For the embedded schemas there are a set of rules for the spec generator to vary the definition it generates.
When an embedded type isn't decorated with the `swagger:allOf` annotation, then the properties from the embedded value will be included in the generated definition as if they were defined on the definition. But when the embedded type is decorated with the `swagger:allOf` annotation then the all of element will be defined as a "$ref" property instead. For an annotated type there is also the possibility to specify an argument, the value of this argument will be used as the value for the `x-class` extension. This allows for generators that support the
`x-class` extension to reliably build a serializer for a type with a discriminator

#### Known vendor extensions

There are a couple of commonly used vendor extensions that most frameworks support to add functionality to the swagger spec.

For generating a swagger specification document this toolkit supports:

Vendor extension | Description
-----------------|-------------
x-isnullable | makes a property value nullable, for go code that means a pointer
x-nullable | makes a property value nullable, for go code that means a pointer
x-go-name | the go name of a type
x-go-package | the go package of a type
x-class | this is used in conjunction with discriminators to give a full type name
x-omitempty | this is used with arrays to control presence of omitempty tag to be used by JSON Marshaler
