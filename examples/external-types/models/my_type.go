package models

import (
	"context"
	"io"

	"github.com/go-openapi/strfmt"
)

// MyType is a type manually added to the models package (NOT GENERATED)
type MyType string

// Validate MyType
func (MyType) Validate(strfmt.Registry) error { return nil }

// ContextValidate MyType
func (MyType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyInteger ...
type MyInteger int

// Validate MyInteger
func (MyInteger) Validate(strfmt.Registry) error { return nil }

// ContextValidate MyInteger
func (MyInteger) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyString ...
type MyString string

// Validate MyString
func (MyString) Validate(strfmt.Registry) error { return nil }

// ContextValidate MyInteger
func (MyString) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyOtherType ...
type MyOtherType struct{}

// Validate MyOtherType
func (MyOtherType) Validate(strfmt.Registry) error { return nil }

// ContextValidate MyOtherType
func (MyOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyStreamer ...
type MyStreamer io.Reader
