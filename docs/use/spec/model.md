# swagger:model

A **swagger:model** annotation optionally gets a model name as extra data on the line.
when this appears anywhere in a comment for a struct, then that struct becomes a schema
in the definitions object of swagger.

<!--more-->

The struct gets analyzed and all the collected models are added to the tree.
The refs are tracked separately so that they can be renamed later on.

Definitions only appear in the generated spec when they are actually used somewhere in the application (eg. on a params or response struct)

##### Syntax:

```
swagger:model [?model name]
```

##### Properties:

Annotation | Description
-----------|------------
**Maximum** | specifies the maximum a number or integer value can have
**Minimum** | specifies the minimum a number or integer value can have
**Multiple of** | specifies a value a number or integer value must be a multiple of
**Minimum length** | the minimum length for a string value
**Maximum length** | the maximum length for a string value
**Pattern** | a regular expression a string value needs to match
**Minimum items** | the minimum number of items a slice needs to have
**Maximum items** | the maximum number of items a slice can have
**Unique** | when set to true the slice can only contain unique items
**Required** | when set to true this value needs to be set on the schema
**Read Only** | when set to true this value will be marked as read-only and is not required in request bodies
**Example** | an example value, parsed as the field's type<br/>(objects and slices are parsed as JSON)
**Extensions** | list of extensions. The field name MUST begin with x-, for example, x-internal-id. The value can be null, a primitive, an array or an object.

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

##### Result:

```yaml
---
definitions:
  User:
    type: object
    title: User represents the user for this application
    description: "A user is the security principal for this application.\nIt's also used as one of the main axes for reporting.\n\nA user can have friends with whom they can share what they like."
    required:
      - id
      - name
      - login
    properties:
      id:
        description: the id for this user
        type: integer
        format: int64
        minimum: 1
      name:
        description: the name for this user
        type: string
        minLength: 3
      login:
        description: the email address for this user
        type: string
        format: email
        x-go-name: Email
        example: user@provider.net
      friends:
        description: the friends for this user
        type: array
        items:
          $ref: "#/definitions/User"
        x-property-value: value
        x-property-array:
          - value1
          - value2
        x-property-array-obj:
          - name: obj
            value: field
```
