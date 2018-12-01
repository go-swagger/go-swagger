# swagger:response

Reads a struct decorated with **swagger:response** and uses that information to fill up the headers and the schema for a response.
A swagger:route can specify a response name for a status code and then the matching response will be used for that operation in the swagger definition.

<!--more-->

##### Syntax:

```
swagger:response [?response name]
```

##### Properties:

Annotation | Description
-----------|------------
**In** | where to find the field
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
// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationError struct {
	// The error message
	// in: body
	Body struct {
		// The validation message
		//
		// Required: true
		// Example: Expected type int
		Message string
		// An optional field name to which this validation applies
		FieldName string
	}
}
```

##### Result:

```yaml
---
responses:
  validationError:
    description: A ValidationError is an error that is used when the required input fails validation.
    schema:
      type: object
      description: The error message
      required:
      - Message
      properties:
        Message:
          type: string
          description: The validation message
          example: Expected type int
        FieldName:
          type: string
          description: an optional field name to which this validation applies
```
