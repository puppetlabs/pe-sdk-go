// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/models"
)

// PostImportReader is a Reader for the PostImport structure.
type PostImportReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostImportReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostImportOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPostImportDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPostImportOK creates a PostImportOK with default headers values
func NewPostImportOK() *PostImportOK {
	return &PostImportOK{}
}

/* PostImportOK describes a response with status code 200, with default header values.

imported
*/
type PostImportOK struct {
	Payload *PostImportOKBody
}

func (o *PostImportOK) Error() string {
	return fmt.Sprintf("[POST /pdb/admin/v1/archive][%d] postImportOK  %+v", 200, o.Payload)
}
func (o *PostImportOK) GetPayload() *PostImportOKBody {
	return o.Payload
}

func (o *PostImportOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostImportOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostImportDefault creates a PostImportDefault with default headers values
func NewPostImportDefault(code int) *PostImportDefault {
	return &PostImportDefault{
		_statusCode: code,
	}
}

/* PostImportDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type PostImportDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the post import default response
func (o *PostImportDefault) Code() int {
	return o._statusCode
}

func (o *PostImportDefault) Error() string {
	return fmt.Sprintf("[POST /pdb/admin/v1/archive][%d] postImport default  %+v", o._statusCode, o.Payload)
}
func (o *PostImportDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *PostImportDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PostImportOKBody puppetdb import
swagger:model PostImportOKBody
*/
type PostImportOKBody struct {

	// whether the operation succeeded
	Ok bool `json:"ok,omitempty"`
}

// Validate validates this post import o k body
func (o *PostImportOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post import o k body based on context it is used
func (o *PostImportOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostImportOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostImportOKBody) UnmarshalBinary(b []byte) error {
	var res PostImportOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
