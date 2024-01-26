---
title: allOf
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:allOf

Marks an embedded type as  a member for allOf

<!--more-->

##### Syntax

```go
swagger:allOf
```

##### Example

```go
// A SimpleOne is a model with a few simple fields
type SimpleOne struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

// A Something struct is used by other structs
type Something struct {
	DID int64  `json:"did"`
	Cat string `json:"cat"`
}

// Notable is a model in a transitive package.
// it's used for embedding in another model
//
// swagger:model withNotes
type Notable struct {
	Notes string `json:"notes"`

	Extra string `json:"extra"`
}

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

##### Result

```yaml
```
