// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/goto/dex/generated/models"
)

// ListProjectsReader is a Reader for the ListProjects structure.
type ListProjectsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListProjectsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListProjectsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListProjectsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewListProjectsOK creates a ListProjectsOK with default headers values
func NewListProjectsOK() *ListProjectsOK {
	return &ListProjectsOK{}
}

/*
ListProjectsOK describes a response with status code 200, with default header values.

successful operation
*/
type ListProjectsOK struct {
	Payload *models.ProjectArray
}

// IsSuccess returns true when this list projects o k response has a 2xx status code
func (o *ListProjectsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this list projects o k response has a 3xx status code
func (o *ListProjectsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list projects o k response has a 4xx status code
func (o *ListProjectsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this list projects o k response has a 5xx status code
func (o *ListProjectsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this list projects o k response a status code equal to that given
func (o *ListProjectsOK) IsCode(code int) bool {
	return code == 200
}

func (o *ListProjectsOK) Error() string {
	return fmt.Sprintf("[GET /dex/projects][%d] listProjectsOK  %+v", 200, o.Payload)
}

func (o *ListProjectsOK) String() string {
	return fmt.Sprintf("[GET /dex/projects][%d] listProjectsOK  %+v", 200, o.Payload)
}

func (o *ListProjectsOK) GetPayload() *models.ProjectArray {
	return o.Payload
}

func (o *ListProjectsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ProjectArray)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListProjectsInternalServerError creates a ListProjectsInternalServerError with default headers values
func NewListProjectsInternalServerError() *ListProjectsInternalServerError {
	return &ListProjectsInternalServerError{}
}

/*
ListProjectsInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type ListProjectsInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this list projects internal server error response has a 2xx status code
func (o *ListProjectsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this list projects internal server error response has a 3xx status code
func (o *ListProjectsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this list projects internal server error response has a 4xx status code
func (o *ListProjectsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this list projects internal server error response has a 5xx status code
func (o *ListProjectsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this list projects internal server error response a status code equal to that given
func (o *ListProjectsInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *ListProjectsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/projects][%d] listProjectsInternalServerError  %+v", 500, o.Payload)
}

func (o *ListProjectsInternalServerError) String() string {
	return fmt.Sprintf("[GET /dex/projects][%d] listProjectsInternalServerError  %+v", 500, o.Payload)
}

func (o *ListProjectsInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListProjectsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
