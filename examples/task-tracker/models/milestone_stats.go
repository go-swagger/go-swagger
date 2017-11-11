// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// MilestoneStats Some counters for this milestone.
//
// This object contains counts for the remaining open issues and the amount of issues that have been closed.
//
// swagger:model milestoneStats
type MilestoneStats struct {

	// The closed issues.
	Closed int32 `json:"closed,omitempty"`

	// The remaining open issues.
	Open int32 `json:"open,omitempty"`

	// The total number of issues for this milestone.
	Total int32 `json:"total,omitempty"`
}

/* polymorph milestoneStats closed false */

/* polymorph milestoneStats open false */

/* polymorph milestoneStats total false */

// Validate validates this milestone stats
func (m *MilestoneStats) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *MilestoneStats) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MilestoneStats) UnmarshalBinary(b []byte) error {
	var res MilestoneStats
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
