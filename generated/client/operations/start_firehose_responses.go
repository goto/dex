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

// StartFirehoseReader is a Reader for the StartFirehose structure.
type StartFirehoseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StartFirehoseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStartFirehoseOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewStartFirehoseBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewStartFirehoseNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewStartFirehoseInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewStartFirehoseOK creates a StartFirehoseOK with default headers values
func NewStartFirehoseOK() *StartFirehoseOK {
	return &StartFirehoseOK{}
}

/*
StartFirehoseOK describes a response with status code 200, with default header values.

Successfully applied update.
*/
type StartFirehoseOK struct {
	Payload *models.Firehose
}

// IsSuccess returns true when this start firehose o k response has a 2xx status code
func (o *StartFirehoseOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this start firehose o k response has a 3xx status code
func (o *StartFirehoseOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this start firehose o k response has a 4xx status code
func (o *StartFirehoseOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this start firehose o k response has a 5xx status code
func (o *StartFirehoseOK) IsServerError() bool {
	return false
}

// IsCode returns true when this start firehose o k response a status code equal to that given
func (o *StartFirehoseOK) IsCode(code int) bool {
	return code == 200
}

func (o *StartFirehoseOK) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseOK  %+v", 200, o.Payload)
}

func (o *StartFirehoseOK) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseOK  %+v", 200, o.Payload)
}

func (o *StartFirehoseOK) GetPayload() *models.Firehose {
	return o.Payload
}

func (o *StartFirehoseOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Firehose)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStartFirehoseBadRequest creates a StartFirehoseBadRequest with default headers values
func NewStartFirehoseBadRequest() *StartFirehoseBadRequest {
	return &StartFirehoseBadRequest{}
}

/*
StartFirehoseBadRequest describes a response with status code 400, with default header values.

Update request is not valid.
*/
type StartFirehoseBadRequest struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this start firehose bad request response has a 2xx status code
func (o *StartFirehoseBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this start firehose bad request response has a 3xx status code
func (o *StartFirehoseBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this start firehose bad request response has a 4xx status code
func (o *StartFirehoseBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this start firehose bad request response has a 5xx status code
func (o *StartFirehoseBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this start firehose bad request response a status code equal to that given
func (o *StartFirehoseBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *StartFirehoseBadRequest) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *StartFirehoseBadRequest) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *StartFirehoseBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StartFirehoseBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStartFirehoseNotFound creates a StartFirehoseNotFound with default headers values
func NewStartFirehoseNotFound() *StartFirehoseNotFound {
	return &StartFirehoseNotFound{}
}

/*
StartFirehoseNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type StartFirehoseNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this start firehose not found response has a 2xx status code
func (o *StartFirehoseNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this start firehose not found response has a 3xx status code
func (o *StartFirehoseNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this start firehose not found response has a 4xx status code
func (o *StartFirehoseNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this start firehose not found response has a 5xx status code
func (o *StartFirehoseNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this start firehose not found response a status code equal to that given
func (o *StartFirehoseNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *StartFirehoseNotFound) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *StartFirehoseNotFound) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *StartFirehoseNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StartFirehoseNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStartFirehoseInternalServerError creates a StartFirehoseInternalServerError with default headers values
func NewStartFirehoseInternalServerError() *StartFirehoseInternalServerError {
	return &StartFirehoseInternalServerError{}
}

/*
StartFirehoseInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type StartFirehoseInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this start firehose internal server error response has a 2xx status code
func (o *StartFirehoseInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this start firehose internal server error response has a 3xx status code
func (o *StartFirehoseInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this start firehose internal server error response has a 4xx status code
func (o *StartFirehoseInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this start firehose internal server error response has a 5xx status code
func (o *StartFirehoseInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this start firehose internal server error response a status code equal to that given
func (o *StartFirehoseInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *StartFirehoseInternalServerError) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *StartFirehoseInternalServerError) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/start][%d] startFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *StartFirehoseInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StartFirehoseInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
