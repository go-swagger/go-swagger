// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Order order
//
// swagger:model Order
type Order struct {

	// complete
	Complete *bool `json:"complete,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// pet Id
	PetID int64 `json:"petId,omitempty"`

	// quantity
	Quantity int32 `json:"quantity,omitempty"`

	// ship date
	// Format: date-time
	ShipDate strfmt.DateTime `json:"shipDate,omitempty"`

	// Order Status
	// Enum: [placed approved delivered]
	Status string `json:"status,omitempty"`
}

// Validate validates this order
func (m *Order) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateShipDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Order) validateShipDate(formats strfmt.Registry) error {
	if swag.IsZero(m.ShipDate) { // not required
		return nil
	}

	if err := validate.FormatOf("shipDate", "body", "date-time", m.ShipDate.String(), formats); err != nil {
		return err
	}

	return nil
}

var orderTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["placed","approved","delivered"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		orderTypeStatusPropEnum = append(orderTypeStatusPropEnum, v)
	}
}

const (

	// OrderStatusPlaced captures enum value "placed"
	OrderStatusPlaced string = "placed"

	// OrderStatusApproved captures enum value "approved"
	OrderStatusApproved string = "approved"

	// OrderStatusDelivered captures enum value "delivered"
	OrderStatusDelivered string = "delivered"
)

// prop value enum
func (m *Order) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, orderTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Order) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this order based on context it is used
func (m *Order) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Order) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Order) UnmarshalBinary(b []byte) error {
	var res Order
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
