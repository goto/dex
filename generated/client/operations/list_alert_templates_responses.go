// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/odpf/dex/generated/models"
)

// ListAlertTemplatesReader is a Reader for the ListAlertTemplates structure.
type ListAlertTemplatesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAlertTemplatesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAlertTemplatesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListAlertTemplatesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewListAlertTemplatesOK creates a ListAlertTemplatesOK with default headers values
func NewListAlertTemplatesOK() *ListAlertTemplatesOK {
	return &ListAlertTemplatesOK{}
}

/*
ListAlertTemplatesOK describes a response with status code 200, with default header values.

successful operation
*/
type ListAlertTemplatesOK struct {
	Payload *models.AlertTemplatesArray
}

// IsSuccess returns true when this list alert templates o k response has a 2xx status code
func (o *ListAlertTemplatesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list alert templates o k response has a 3xx status code
func (o *ListAlertTemplatesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list alert templates o k response has a 4xx status code
func (o *ListAlertTemplatesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list alert templates o k response has a 5xx status code
func (o *ListAlertTemplatesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list alert templates o k response a status code equal to that given
func (o *ListAlertTemplatesOK) IsCode(code int) bool {
	return code == 200
}

func (o *ListAlertTemplatesOK) Error() string {
	return fmt.Sprintf("[GET /alertTemplates][%d] listAlertTemplatesOK  %+v", 200, o.Payload)
}

func (o *ListAlertTemplatesOK) String() string {
	return fmt.Sprintf("[GET /alertTemplates][%d] listAlertTemplatesOK  %+v", 200, o.Payload)
}

func (o *ListAlertTemplatesOK) GetPayload() *models.AlertTemplatesArray {
	return o.Payload
}

func (o *ListAlertTemplatesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AlertTemplatesArray)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAlertTemplatesInternalServerError creates a ListAlertTemplatesInternalServerError with default headers values
func NewListAlertTemplatesInternalServerError() *ListAlertTemplatesInternalServerError {
	return &ListAlertTemplatesInternalServerError{}
}

/*
ListAlertTemplatesInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type ListAlertTemplatesInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this list alert templates internal server error response has a 2xx status code
func (o *ListAlertTemplatesInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list alert templates internal server error response has a 3xx status code
func (o *ListAlertTemplatesInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list alert templates internal server error response has a 4xx status code
func (o *ListAlertTemplatesInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this list alert templates internal server error response has a 5xx status code
func (o *ListAlertTemplatesInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this list alert templates internal server error response a status code equal to that given
func (o *ListAlertTemplatesInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *ListAlertTemplatesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /alertTemplates][%d] listAlertTemplatesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListAlertTemplatesInternalServerError) String() string {
	return fmt.Sprintf("[GET /alertTemplates][%d] listAlertTemplatesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListAlertTemplatesInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAlertTemplatesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
