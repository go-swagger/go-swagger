---
title: model
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:model

A **swagger:model** annotation optionally gets a model name as extra data on the line.
when this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

<!--more-->

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

Definitions only appear in the generated spec when they are actually used somewhere in the application (eg. on a params or response struct)

##### Syntax

```go
swagger:model [?model name]
```

##### Properties

Annotation | Description
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
// User represents the user for this application
//
// A user is the security principal for this application.
// It's also used as one of main axes for reporting.
//
// A user can have friends with whom they can share what they like.
//
// swagger:model
type User struct {
	// the id for this user
	//
	// required: true
	// min: 1
	ID int64 `json:"id"`

	// the name for this user
	// required: true
	// min length: 3
	Name string `json:"name"`

	// the email address for this user
	//
	// required: true
	// example: user@provider.net
	Email strfmt.Email `json:"login"`

	// the friends for this user
	//
	// Extensions:
	// ---
	// x-property-value: value
	// x-property-array:
	//   - value1
	//   - value2
	// x-property-array-obj:
	//   - name: obj
	//     value: field
	// ---
	Friends []User `json:"friends"`
}
```

##### Result

```yaml
```
