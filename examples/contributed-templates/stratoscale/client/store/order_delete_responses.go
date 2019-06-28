// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// OrderDeleteReader is a Reader for the OrderDelete structure.
type OrderDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *OrderDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewOrderDeleteNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewOrderDeleteBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewOrderDeleteNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewOrderDeleteNoContent creates a OrderDeleteNoContent with default headers values
func NewOrderDeleteNoContent() *OrderDeleteNoContent {
	return &OrderDeleteNoContent{}
}

/*OrderDeleteNoContent handles this case with default header values.

Deleted successfully
*/
type OrderDeleteNoContent struct {
}

func (o *OrderDeleteNoContent) Error() string {
	return fmt.Sprintf("[DELETE /store/order/{orderId}][%d] orderDeleteNoContent ", 204)
}

func (o *OrderDeleteNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewOrderDeleteBadRequest creates a OrderDeleteBadRequest with default headers values
func NewOrderDeleteBadRequest() *OrderDeleteBadRequest {
	return &OrderDeleteBadRequest{}
}

/*OrderDeleteBadRequest handles this case with default header values.

Invalid ID supplied
*/
type OrderDeleteBadRequest struct {
}

func (o *OrderDeleteBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /store/order/{orderId}][%d] orderDeleteBadRequest ", 400)
}

func (o *OrderDeleteBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewOrderDeleteNotFound creates a OrderDeleteNotFound with default headers values
func NewOrderDeleteNotFound() *OrderDeleteNotFound {
	return &OrderDeleteNotFound{}
}

/*OrderDeleteNotFound handles this case with default header values.

Order not found
*/
type OrderDeleteNotFound struct {
}

func (o *OrderDeleteNotFound) Error() string {
	return fmt.Sprintf("[DELETE /store/order/{orderId}][%d] orderDeleteNotFound ", 404)
}

func (o *OrderDeleteNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
