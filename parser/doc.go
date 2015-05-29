/*Package parser provides a parser for go files that produces a swagger spec document.

You give it a main file and it will parse all the files that are required by that main
package to produce a swagger specification.  This parser doesn't function standalone though; it needs you to provide it
with a template swagger document. It seems weird to me to have to write doc comments to produce JSON when the only purpose of
those comments is to provide that JSON. So instead of duplicating work it can take a swagger spec to start from and then
synchronize the content of the specs definitions and parameters, with whatever is in your code.
This way this parser can focus on just analyzing go code instead of defining a lossy grammar that will look crappy
when it's displayed on godoc.org.

The parser supports filters to limit the packages to scan for rest operations
There are also filters for models so you can specify which packages
to include when scanning for models.
When a model has a filtered model as field then that filtered model will be
included transitively.

To use you can add a go:generate comment to your main file for example:

		// go:generate swagger generate spec .

The following annotations exist:

+swagger:meta

The +swagger:meta annotation flags a file as source for metadata about the API.
This is typically a doc.go file with your package documentation.

You can specify a Consumes and Produces key which has a new content type on each line
Schemes is a tag that is required and allows for a comma separated string composed of:
http, https, ws or wss

Host and BasePath can be specified but those values will be defaults,
they should get substituted when serving the swagger spec.

Default parameters and responses are not supported at this stage, for those you can edit the template json.

+swagger:strfmt [name]

A +swagger:strfmt annotation names a type as a string formatter. The name is mandatory and that is
what will be used as format name for this particular string format.
String formats should only be used for very well known formats.

+swagger:model [?model name]

A +swagger:model annotation optionally gets a model name as extra data on the line.
when this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

+swagger:route [method] [path pattern] [operation id] [?tag1 tag2 tag3]

A +swagger:route annotation links a path to a method.
This operation gets a unique id, which is used in various places as method name.
One such usage is in method names for client generation for example.

Because there are many routers available, this tool does not try to parse the paths
you provided to your routing library of choice. So you have to specify your path pattern
yourself in valid swagger syntax.

+swagger:params [operationid1 operationid2]

Links a struct to one or more operations. The params in the resulting swagger spec can be composed of several structs.
There are no guarantees given on how property name overlaps are resolved when several structs apply to the same operation.
This tag works very similar to the swagger:model tag except that it produces valid parameter objects instead of schema
objects.

*/
package parser
