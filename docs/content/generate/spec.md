+++
date = "2015-11-09T18:58:47-08:00"
title = "spec"
+++

The toolkit has a command that will let you generate a swagger spec document from your code. 
The command integrates with go doc comments, and makes use of structs when it needs to know of
types.

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

### Possible annotations

#### swagger:meta

The **swagger:meta** annotation flags a file as source for metadata about the API.
This is typically a doc.go file with your package documentation.

You can specify a Consumes and Produces key which has a new content type on each line
Schemes is a tag that is required and allows for a comma separated string composed of:
http, https, ws or wss

Host and BasePath can be specified but those values will be defaults,
they should get substituted when serving the swagger spec.

Default parameters and responses are not supported at this stage, for those you can edit the template json.

#### swagger:strfmt

```
swagger:strfmt [name]
```

A **swagger:strfmt** annotation names a type as a string formatter. The name is mandatory and that is
what will be used as format name for this particular string format.
String formats should only be used for very well known formats.

#### swagger:model

```
swagger:model [?model name]
```

A **swagger:model** annotation optionally gets a model name as extra data on the line.
when this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

#### swagger:route

```
swagger:route [method] [path pattern] [operation id] [?tag1 tag2 tag3]
```

A **swagger:route** annotation links a path to a method.
This operation gets a unique id, which is used in various places as method name.
One such usage is in method names for client generation for example.

Because there are many routers available, this tool does not try to parse the paths
you provided to your routing library of choice. So you have to specify your path pattern
yourself in valid swagger syntax.

#### swagger:params

```
swagger:params [operationid1 operationid2]
```

Links a struct to one or more operations. The params in the resulting swagger spec can be composed of several structs.
There are no guarantees given on how property name overlaps are resolved when several structs apply to the same operation.
This tag works very similar to the swagger:model tag except that it produces valid parameter objects instead of schema
objects.

#### swagger:response

```
swagger:response [?response name]
```

Reads a struct decorated with **swagger:response** and uses that information to fill up the headers and the schema for a response.
A swagger:route can specify a response name for a status code and then the matching response will be used for that operation in the swagger definition.

#### swagger:allOf

```
swagger:allOf
```

Marks an embedded type as  a member for allOf

```go
// An AllOfModel is composed out of embedded structs but it should build
// an allOf property
type AllOfModel struct {
	// swagger:allOf
	SimpleOne
	// swagger:allOf
	mods.Notable

	Something // not annotated with anything, so should be included

	CreatedAt strfmt.DateTime `json:"createdAt"`
}
```
