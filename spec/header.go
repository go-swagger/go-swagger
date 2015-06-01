package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/util"
)

type headerProps struct {
	Description string `json:"description,omitempty"`
}

// Header describes a header for a response of the API
//
// For more information: http://goo.gl/8us55a#headerObject
type Header struct {
	commonValidations
	simpleSchema
	headerProps
}

// Typed a fluent builder method for the type of parameter
func (h *Header) Typed(tpe, format string) *Header {
	h.Type = tpe
	h.Format = format
	return h
}

// CollectionOf a fluent builder method for an array item
func (h *Header) CollectionOf(items *Items, format string) *Header {
	h.Type = "array"
	h.Items = items
	h.CollectionFormat = format
	return h
}

// MarshalJSON marshal this to JSON
func (h Header) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(h.commonValidations)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(h.simpleSchema)
	if err != nil {
		return nil, err
	}
	b3, err := json.Marshal(h.headerProps)
	if err != nil {
		return nil, err
	}
	return util.ConcatJSON(b1, b2, b3), nil
}

// UnmarshalJSON marshal this from JSON
func (h *Header) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &h.commonValidations); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &h.simpleSchema); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &h.headerProps); err != nil {
		return err
	}
	return nil
}
