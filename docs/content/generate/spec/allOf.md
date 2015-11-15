+++
categories = ["spec"]
tags = ["polymorphism"]
date = "2015-11-14T20:10:58-08:00"
title = "swagger:allOf"
weight = 29
+++

Marks an embedded type as  a member for allOf

<!--more-->

##### Syntax:

```
swagger:allOf
```

##### Example:

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
