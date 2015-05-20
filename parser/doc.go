// Package parser provides a parser for go files that produces a swagger spec document
//
// You give it a main file and it will parse all the files that are required by that main
// package to produce a swagger specification.
//
// The parser supports filters to limit the packages to scan for rest operations
// There are also filters for models so you can specify which packages
// to include when scanning for models.
// When a model has a filtered model as field then that filtered model will be
// included transitively.
//
// The following annotations exist:
//
// +swagger:meta
//
// The +swagger:meta annotation flags a file as source for metadata about the API.
// This is typically a doc.go file with your package documentation.
//
// You can specify a Consumes and Produces key which has a new content type on each line
// Schemes is a tag that is required and allows for a comma separated string composed of:
// http, https, ws or wss
//
// Host and BasePath can be specified but those values will be defaults,
// they should get substituted when serving the swagger spec
//
// Default parameters and responses are not supported at this stage
//
// Tags are a mapping of tag to package
//
// +swagger:model [?model name]
//
// A +swagger:model annotation optionally gets a model name as extra data on the line.
// when this appears anywhere in a comment for a struct, then that struct becomes a schema
// in the definitions object of swagger.
//
// The struct gets analyzed and all the collected models are added to the tree.
// The refs are tracked separately so that they can be renamed later on.
//
// +swagger:route [method] [path pattern] [?tag name:tag name:][operation id] [?params model]
//
// A +swagger:route annotation links a path to a method and describes a set of parameters for
// the operation. This operation gets a unique id, which is used in various places as method name.
// One such usage is in client generation for example.
//
// Because there are many routers available, this tool does not try to parse the paths
// you provided to your routing library of choice. So you have to specify your path pattern
// yourself in valid swagger syntax. If you use a struct to represent your request parameters,
// you can provide the name of the struct and that will then be parsed into operation parameters.
//
// A route can have several sub tags, that build up all the required meta data for the paths object
// in a swagger spec document.
package parser
