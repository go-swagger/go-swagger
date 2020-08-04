// Code generated by go-swagger; DO NOT EDIT.

package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/go-swagger/go-swagger/examples/contributed-templates/stratoscale/models"
)

// PetListReader is a Reader for the PetList structure.
type PetListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PetListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPetListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPetListBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPetListOK creates a PetListOK with default headers values
func NewPetListOK() *PetListOK {
	return &PetListOK{}
}

/*PetListOK handles this case with default header values.

successful operation
*/
type PetListOK struct {
	Payload []*models.Pet
}

func (o *PetListOK) Error() string {
	return fmt.Sprintf("[GET /pet][%d] petListOK  %+v", 200, o.Payload)
}

func (o *PetListOK) GetPayload() []*models.Pet {
	return o.Payload
}

func (o *PetListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPetListBadRequest creates a PetListBadRequest with default headers values
func NewPetListBadRequest() *PetListBadRequest {
	return &PetListBadRequest{}
}

/*PetListBadRequest handles this case with default header values.

Invalid status value
*/
type PetListBadRequest struct {
}

func (o *PetListBadRequest) Error() string {
	return fmt.Sprintf("[GET /pet][%d] petListBadRequest ", 400)
}

func (o *PetListBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
