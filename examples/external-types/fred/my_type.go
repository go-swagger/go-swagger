package fred

import (
	"context"
  "io"

	"github.com/go-openapi/strfmt"
)

// MyAlternateType ...
type MyAlternateType string

// Validate MyAlternateType
func (MyAlternateType) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateInteger ...
type MyAlternateInteger int

// Validate MyAlternateInteger
func (MyAlternateInteger) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateInteger) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateString ...
type MyAlternateString string

// Validate MyAlternateString
func (MyAlternateString) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateString) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateOtherType ...
type MyAlternateOtherType struct{}

// Validate MyAlternateOtherType
func (MyAlternateOtherType) Validate(strfmt.Registry) error                         { return nil }
func (MyAlternateOtherType) ContextValidate(context.Context, strfmt.Registry) error { return nil }

// MyAlternateStreamer ...
type MyAlternateStreamer io.Reader

// MyAlternateInterface ...
type MyAlternateInterface interface{}
