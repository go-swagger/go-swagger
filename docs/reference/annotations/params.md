---
title: params
date: 2023-01-01T01:01:01-08:00
draft: true
---
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

##### Syntax

```go
swagger:parameters [operationid1 operationid2]
```

##### Properties

The fields of this struct can be decorated with a number of annotations. For the field name it uses the struct field
name, it respects the json struct field tag for customizing the name.

Annotation | Format
**Items.*n*.Maximum** |  specifies the maximum a number or integer value can have at the level *n*
**Items.*n*.Minimum** |  specifies the minimum a number or integer value can have at the level *n*
**Items.*n*.Multiple of** | specifies a value a number or integer value must be a multiple of
**Items.*n*.Minimum length** | the minimum length for a string value at the level *n*
**Items.*n*.Maximum length** | the maximum length for a string value at the level *n*
**Items.*n*.Pattern** | a regular expression a string value needs to match at the level *n*
**Items.*n*.Minimum items** | the minimum number of items a slice needs to have at the level *n*
**Items.*n*.Maximum items** | the maximum number of items a slice can have at the level *n*
**Items.*n*.Unique** | when set to true the slice can only contain unique items at the level *n*

##### Example

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
  // Extensions:
  //   x-example-flag: true
  //   x-some-list:
  //     - dog
  //     - cat
  //     - bird
	BarSlice [][][]string `json:"bar_slice"`
}
```

##### Result

```yaml
```
