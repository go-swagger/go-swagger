# swagger:parameters

The **swagger:parameters** annotation links a struct to one or more operations. The parameters in the resulting swagger spec can be composed of several structs.
There are no guarantees given on how property name overlaps are resolved when several structs apply to the same operation.
This tag works very similarly to the swagger:model tag except that it produces valid parameter objects instead of schema
objects.
<!--more-->
When this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

At this moment the parameters require one or more structs to be defined, it's not possible to annotate plain var
entries at this moment.

##### Syntax:

```
swagger:parameters [operationid1 operationid2]
```

##### Properties:

The fields of this struct can be decorated with a number of annotations. For the field name it uses the struct field
name, it respects the json struct field tag for customizing the name.

Annotation | Format
---------- | ------
**In** | where to find the parameter
**Collection Format** | when a slice the formatter for the collection when serialized on the request
**Maximum** | specifies the maximum a number or integer value can have
**Minimum** | specifies the minimum a number or integer value can have
**Multiple of** | specifies a value a number or integer value must be a multiple of
**Minimum length** | the minimum length for a string value
**Maximum length** | the maximum length for a string value
**Pattern** | a regular expression a string value needs to match
**Minimum items** | the minimum number of items a slice needs to have
**Maximum items** | the maximum number of items a slice can have
**Unique** | when set to true the slice can only contain unique items
**Required** | when set to true this value needs to be present in the request
**Example** | an example value, parsed as the field's type<br/>(objects and slices are parsed as JSON)

For slice properties there are also items to be defined. This might be a nested collection, for indicating nesting
level the value is a 0-based index, so items.minLength is the same as items.0.minLength

Annotation | Format
-----------|--------
**Items.*n*.Maximum** |  specifies the maximum a number or integer value can have at the level *n*
**Items.*n*.Minimum** |  specifies the minimum a number or integer value can have at the level *n*
**Items.*n*.Multiple of** | specifies a value a number or integer value must be a multiple of
**Items.*n*.Minimum length** | the minimum length for a string value at the level *n*
**Items.*n*.Maximum length** | the maximum length for a string value at the level *n*
**Items.*n*.Pattern** | a regular expression a string value needs to match at the level *n*
**Items.*n*.Minimum items** | the minimum number of items a slice needs to have at the level *n*
**Items.*n*.Maximum items** | the maximum number of items a slice can have at the level *n*
**Items.*n*.Unique** | when set to true the slice can only contain unique items at the level *n*

##### Example:

```go
// swagger:parameters listBars addBars
type BarSliceParam struct {
	// a BarSlice has bars which are strings
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minItems: 4
	// items.maxItems: 9
	// items.items.minItems: 5
	// items.items.maxItems: 8
	// items.items.items.minLength: 3
	// items.items.items.maxLength: 10
	// items.items.items.pattern: \w+
	// collection format: pipe
	// in: query
	// example: [[["bar_000"]]]
	BarSlice [][][]string `json:"bar_slice"`
}
```

##### Result:

```yaml
---
operations:
  "/":
    get:
      operationId: listBars
      parameters:
        - name: bar_slice
          in: query
          maxItems: 10
          minItems: 3
          unique: true
          collectionFormat: pipe
          type: array
          example:
            - - - "bar_000"
          items:
            type: array
            maxItems: 9
            minItems: 4
            items:
              type: array
              maxItems: 8
              minItems: 5
              items:
                type: string
                minLength: 3
                maxLength: 10
                pattern: "\\w+"
    post:
      operationId: addBars
      parameters:
        - name: bar_slice
          in: query
          maxItems: 10
          minItems: 3
          unique: true
          collectionFormat: pipe
          type: array
          example:
            - - - "bar_000"
          items:
            type: array
            maxItems: 9
            minItems: 4
            items:
              type: array
              maxItems: 8
              minItems: 5
              items:
                type: string
                minLength: 3
                maxLength: 10
                pattern: "\\w+"
```
