---
title: type
date: 2023-01-01T01:01:01-08:00
draft: true
---
# swagger:type

**swagger:type** can be used to mark a struct with its type. This will override the type inferred by go-swagger.

[These data types](https://swagger.io/docs/specification/data-models/data-types/) are supported by Swagger.

##### Syntax

```go
swagger:type [type]
```

##### Example

```go
// swagger:type string
type NullString struct {
     sql.NullString
}

// swagger:model myString
type MyString struct {
     NS NullString
}
```

##### Result

```yaml
```
