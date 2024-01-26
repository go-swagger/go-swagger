---
title: operation
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:operation

A **swagger:operation** annotation links a path to a method.
This operation gets a unique id, which is used in various places as method name.
One such usage is in method names for client generation for example.

Because there are many routers available, this tool does not try to parse the paths
you provided to your routing library of choice. So you have to specify your path pattern
yourself in valid swagger (YAML) syntax.

<!--more-->

##### Syntax

```go
swagger:operation [method] [path pattern] [?tag1 tag2 tag3] [operation id]
```

##### Properties

Any valid Swagger 2.0 YAML _Operation_ property is valid.
Make sure your indentation is consistent and correct,
as it won't parse correctly otherwise.

You can find all the properties at https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#operationObject

Note that if you specify the _tags_, _summary_, _description_ or _operationId_ as part of the
YAML spec, you will override the _summary_, _descriptions_, _tags_ or _operationId_, specified as part of the regular swagger syntax above.

Also note that you need to start your YAML spec with a triple dash `---`.

##### Example

```go
// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) (err error) {
	// swagger:operation GET /pets getPet
	//
	// Returns all pets from the system that the user has access to
	//
	// Could be any pet
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// parameters:
	// - name: tags
	//   in: query
	//   description: tags to filter by
	//   required: false
	//   type: array
	//   items:
	//     type: string
	//   collectionFormat: csv
	// - name: limit
	//   in: query
	//   description: maximum number of results to return
	//   required: false
	//   type: integer
	//   format: int32
	// responses:
	//   '200':
	//     description: pet response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/pet"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorModel"
	mountItem("GET", basePath+"/pets", nil)

    return
}
```

##### Result

```yaml
```
