# swagger:allOf

Marks an embedded type as  a member for allOf

<!--more-->

##### Syntax:

```
swagger:allOf
```

##### Example:

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

##### Result:

```yaml
---
definitions:
  SimpleOne:
    title: A SimpleOne is a model with a few simple fields
    type: object
    properties:
      id:
        type: integer
        format: int64
      name:
        type: string
      age:
        type: integer
        format: int32
  Notable:
    title: "Notable is a model in a transitive package.\nit's used for embedding in another model"
    type: object
    properties:
      notes:
        type: string
      extra:
        type: string
  AllOfModel:
    title: "An AllOfModel is composed out of embedded structs but it should build\nan allOf property"
    allOf: 
      - $ref: "#/definitions/SimpleOne"
      - $ref: "#/definitions/Notable"
      - title: A Something struct is used by other structs
        type: object
        properties:
          did:
            type: integer
            format: int64
          cat:
            type: string
      - type: object
        properties:
          createdAt:
            type: string
            format: date-time
```
