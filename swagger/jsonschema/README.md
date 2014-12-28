[![Build Status](https://travis-ci.org/jabley/gojsonschema.svg?branch=master)](https://travis-ci.org/jabley/gojsonschema)

# gojsonschema

## Description
An implementation of JSON Schema, based on IETF's draft v4 - Go language

## Status

Test results: passed 100% of JSON Schema Test Suite (248/248).

## Usage 

### Quick example

```go

package main

import (
    "fmt"
    "github.com/xeipuuv/gojsonschema"
)

func main() {

    // Loads a schema remotely
    schemaDocument, err := gojsonschema.NewJsonSchemaDocument("http://host/schema.json")
    if err != nil {
        panic(err.Error())
    }

    // Loads the JSON to validate from a local file
    jsonDocument, err := gojsonschema.GetFileJson("/home/me/data.json")
    if err != nil {
        panic(err.Error())
    }

	// Try to validate the Json against the schema
    result := schemaDocument.Validate(jsonDocument)

	// Deal with result
    if result.Valid() {
        fmt.Printf("The document is valid\n")
    } else {
        fmt.Printf("The document is not valid. see errors :\n")
        // Loop through errors
        for _, desc := range result.Errors() {
            fmt.Printf("- %s\n", desc)
        }
    }

}


```

#### Loading a schema

Schemas can be loaded remotely from a Http Url:

```go
    schemaDocument, err := gojsonschema.NewJsonSchemaDocument("http://myhost/schema.json")
```

Or a local file, using the file URI scheme:

```go
	schemaDocument, err := gojsonschema.NewJsonSchemaDocument("file:///home/me/schema.json")
```

You may also load the schema from within your code, using a map[string]interface{} variable.

Note that schemas loaded that way are subject to limitations, they need to be standalone schemas; 
Which means references to local files and/or remote files within this schema will not work.

```go
	schemaMap := map[string]interface{}{
		"type": "string"}

	schemaDocument, err := gojsonschema.NewJsonSchemaDocument(schemaMap)
```

#### Loading a JSON

The library virtually accepts any Json since it uses reflection to validate against the schema.

You may use and combine go types like 
* string (JSON string)
* bool (JSON boolean)
* float64 (JSON number)
* nil (JSON null)
* slice (JSON array)
* map[string]interface{} (JSON object)

You may declare your Json from within your code:

```go
	jsonDocument := map[string]interface{}{
		"name": "john"}
```

Helper functions are also available to load from a Http URL:

```go
    jsonDocument, err := gojsonschema.GetHttpJson("http://host/data.json")
```

Or a local file:

```go
	jsonDocument, err := gojsonschema.GetFileJson("/home/me/data.json")
```

#### Validation

Once the schema and the Json to validate are loaded, validation phase becomes easy:

```go
	result := schemaDocument.Validate(jsonDocument)
```

Check the result validity with:

```go
	if result.Valid() {
		// Your Json is valid
	}
```

If not valid, you can loop through the error messages returned by the validation phase:

```go
	for _, desc := range result.Errors() {
    	fmt.Printf("Error: %s\n", desc)
	}
```

## References

###Website
http://json-schema.org

###Schema Core
http://json-schema.org/latest/json-schema-core.html

###Schema Validation
http://json-schema.org/latest/json-schema-validation.html

## Dependencies
https://github.com/xeipuuv/gojsonpointer

https://github.com/xeipuuv/gojsonreference

## Uses

gojsonschema uses the following test suite :

https://github.com/json-schema/JSON-Schema-Test-Suite
