// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/puppetlabs/pe-sdk-go/app/puppetdb-cli/api/models"
)

// GetQueryReader is a Reader for the GetQuery structure.
type GetQueryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetQueryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetQueryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetQueryBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetQueryForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetQueryInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetQueryDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetQueryOK creates a GetQueryOK with default headers values
func NewGetQueryOK() *GetQueryOK {
	return &GetQueryOK{}
}

/*GetQueryOK handles this case with default header values.

returns query response
*/
type GetQueryOK struct {
	Payload interface{}
}

func (o *GetQueryOK) Error() string {
	return fmt.Sprintf("[GET /pdb/query/v4][%d] getQueryOK  %+v", 200, o.Payload)
}

func (o *GetQueryOK) GetPayload() interface{} {
	return o.Payload
}

func (o *GetQueryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetQueryBadRequest creates a GetQueryBadRequest with default headers values
func NewGetQueryBadRequest() *GetQueryBadRequest {
	return &GetQueryBadRequest{}
}

/*GetQueryBadRequest handles this case with default header values.

PQL parse error response
*/
type GetQueryBadRequest struct {
	Payload string
}

func (o *GetQueryBadRequest) Error() string {
	return fmt.Sprintf("[GET /pdb/query/v4][%d] getQueryBadRequest  %+v", 400, o.Payload)
}

func (o *GetQueryBadRequest) GetPayload() string {
	return o.Payload
}

func (o *GetQueryBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetQueryForbidden creates a GetQueryForbidden with default headers values
func NewGetQueryForbidden() *GetQueryForbidden {
	return &GetQueryForbidden{}
}

/*GetQueryForbidden handles this case with default header values.

Permission denied response
*/
type GetQueryForbidden struct {
	Payload string
}

func (o *GetQueryForbidden) Error() string {
	return fmt.Sprintf("[GET /pdb/query/v4][%d] getQueryForbidden  %+v", 403, o.Payload)
}

func (o *GetQueryForbidden) GetPayload() string {
	return o.Payload
}

func (o *GetQueryForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetQueryInternalServerError creates a GetQueryInternalServerError with default headers values
func NewGetQueryInternalServerError() *GetQueryInternalServerError {
	return &GetQueryInternalServerError{}
}

/*GetQueryInternalServerError handles this case with default header values.

Server error
*/
type GetQueryInternalServerError struct {
	Payload *models.ServerError
}

func (o *GetQueryInternalServerError) Error() string {
	return fmt.Sprintf("[GET /pdb/query/v4][%d] getQueryInternalServerError  %+v", 500, o.Payload)
}

func (o *GetQueryInternalServerError) GetPayload() *models.ServerError {
	return o.Payload
}

func (o *GetQueryInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServerError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetQueryDefault creates a GetQueryDefault with default headers values
func NewGetQueryDefault(code int) *GetQueryDefault {
	return &GetQueryDefault{
		_statusCode: code,
	}
}

/*GetQueryDefault handles this case with default header values.

Unexpected error
*/
type GetQueryDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get query default response
func (o *GetQueryDefault) Code() int {
	return o._statusCode
}

func (o *GetQueryDefault) Error() string {
	return fmt.Sprintf("[GET /pdb/query/v4][%d] getQuery default  %+v", o._statusCode, o.Payload)
}

func (o *GetQueryDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetQueryDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
