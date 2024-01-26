---
title: discriminated
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:discriminated

Marks an embedded type as  a member for allOf and sets the x-class value. On interface definitions there is another annotation on methods allowed _swagger:name_

<!--more-->

The swagger:allOf annotation can be followed by a string. This string will be the value for the `x-class` vendor extension. This value is used as constant for the discriminator field.

An interface that is embedded expects to have 1 method that is commented with `Discriminator: true`. That field will be used as discriminator field when marshalling/unmarshalling objects.

Because this is for use with interfaces we can't make use of the json struct tag to allow for spec name overriding. So instead you can annotate method names on an interface with swagger:name and a value this will then provide the json field name.

##### Syntax

```go
swagger:allOf org.example.something.TypeName
```

##### Example

```go
// TeslaCar is a tesla car
//
// swagger:model
type TeslaCar interface {
	// The model of tesla car
	//
	// discriminator: true
	// swagger:name model
	Model() string

	// AutoPilot returns true when it supports autopilot
	// swagger:name autoPilot
	AutoPilot() bool
}

// The ModelS version of the tesla car
//
// swagger:model modelS
type ModelS struct {
	// swagger:allOf com.tesla.models.ModelS
	TeslaCar
	// The edition of this Model S
	Edition string `json:"edition"`
}

// The ModelX version of the tesla car
//
// swagger:model modelX
type ModelX struct {
	// swagger:allOf com.tesla.models.ModelX
	TeslaCar
	// The number of doors on this Model X
	Doors int32 `json:"doors"`
}
```

##### Result

```yaml
```
