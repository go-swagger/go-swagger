// +build external-types

package external

import "github.com/go-openapi/strfmt"

type Error struct {
	code int
}

func (e *Error) Validate(_ strfmt.Registry) error {
	return nil
}
