// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// TaskAllOf1 task all of1
// swagger:model taskAllOf1
type TaskAllOf1 struct {

	// attachments
	Attachments TaskAllOf1Attachments `json:"attachments,omitempty"`

	// comments
	Comments TaskAllOf1Comments `json:"comments"`

	// The time at which this issue was last updated.
	//
	// This field is read only so it's only sent as part of the response.
	//
	// Read Only: true
	LastUpdated strfmt.DateTime `json:"lastUpdated,omitempty"`

	// last updated by
	LastUpdatedBy *UserCard `json:"lastUpdatedBy,omitempty"`

	// reported by
	ReportedBy *UserCard `json:"reportedBy,omitempty"`
}

/* polymorph taskAllOf1 attachments false */

/* polymorph taskAllOf1 comments false */

/* polymorph taskAllOf1 lastUpdated false */

/* polymorph taskAllOf1 lastUpdatedBy false */

/* polymorph taskAllOf1 reportedBy false */

// Validate validates this task all of1
func (m *TaskAllOf1) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLastUpdatedBy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateReportedBy(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskAllOf1) validateLastUpdatedBy(formats strfmt.Registry) error {

	if swag.IsZero(m.LastUpdatedBy) { // not required
		return nil
	}

	if m.LastUpdatedBy != nil {

		if err := m.LastUpdatedBy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("lastUpdatedBy")
			}
			return err
		}
	}

	return nil
}

func (m *TaskAllOf1) validateReportedBy(formats strfmt.Registry) error {

	if swag.IsZero(m.ReportedBy) { // not required
		return nil
	}

	if m.ReportedBy != nil {

		if err := m.ReportedBy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("reportedBy")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskAllOf1) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskAllOf1) UnmarshalBinary(b []byte) error {
	var res TaskAllOf1
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
