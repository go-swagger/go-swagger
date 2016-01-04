package stubkit

import "github.com/go-swagger/go-swagger/spec"

// StringGen a type for things that can generate strings
type StringGen func() string

// StringGenerator creates a string gen for a field name or generator hint
func StringGenerator(name string) StringGen {
	return func() string { return "" }
}

// SchemaGen a type for things that can generate data for a schema definition
// composed out of the more primitive data generators
type SchemaGen func() interface{}

// SchemaGenerator creates a generator for a schema
func SchemaGenerator(schema spec.Schema) SchemaGen {
	return func() interface{} { return nil }
}
